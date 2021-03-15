package loader

import (
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rgaryaev/pspservice/config"
)

func doRequest(method string, url string) (*http.Response, error) {
	httpClient := &http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

// GetFileDate retrievs last modified date of the file.
func getFileDate(url string) (time.Time, error) {

	response, err := doRequest("HEAD", url)
	if err != nil {
		return time.Time{}, err
	}

	lastModified, err := time.Parse(time.RFC1123, response.Header[http.CanonicalHeaderKey("last-modified")][0])
	if err != nil {
		return time.Time{}, err
	}

	// workaround for truncation date  with localization
	lastModified, err = time.Parse(config.ConfigDateFormat, lastModified.Local().Format(config.ConfigDateFormat))
	if err != nil {
		return time.Time{}, err
	}
	return lastModified, nil

}

// DownloadFile dowloads passport data file and save it as csv file
func downloadFile(filepath string, url string) (err error) {

	// Create the destination file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	response, err := doRequest(http.MethodGet, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Check server response
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("response status: %s", response.Status)
	}

	// Writer the body to file
	// using bzip decompression
	bzp := bzip2.NewReader(response.Body)
	_, err = io.Copy(out, bzp)
	if err != nil {
		return err
	}

	return nil
}

// LoadPassportData checks the schedule and load new passport data if need
// return bool if file data was not exists before
func LoadPassportData(cfg *config.Configuration) (bool, error) {

	// check explicitly if the file doesn't exist then  start the force update
	// usually when the app is starting first time
	var isNotExists bool = false
	_, err := os.Stat(cfg.Storage.PassportData)

	if os.IsNotExist(err) {
		isNotExists = true
		// create directory if need
		wd := "./.data/"
		if _, err := os.Stat(wd); os.IsNotExist(err) {
			log.Println("work dir is created")
			os.Mkdir(wd, os.ModeDir)
		}
	} else if err != nil {
		// other sort of error
		return false, errors.New("cannot check existance of the passport data file")
	}

	for {
		// refresh config every iteration in case of changing some parameters
		cfg, err := config.LoadConfiguration()
		if err != nil {
			// other sort of error
			return true, errors.New("cannot read config file")
		}
		// last date of succesfull update
		lastUpdate, err := time.Parse(config.ConfigDateFormat, cfg.Loader.LastUpdate)
		if err != nil {
			return false, err
		}
		//  jast in case - long update period is not allowed
		// from 1 to 30 days
		if !(cfg.Loader.EveryXDay > 0 && cfg.Loader.EveryXDay <= 30) {
			return false, errors.New("update period is incorrect. it can be 1 till 30")
		}
		// to last update date add period of updating
		nextUpdate := lastUpdate.AddDate(0, 0, cfg.Loader.EveryXDay)

		// skip all schedule checks if file just not exists and do force update
		if !isNotExists {

			if nextUpdate.After(time.Now()) {
				// every 120 min
				time.Sleep(120 * time.Minute)
				continue
			}

			// Check version of the source file
			fileLastModified, err := getFileDate(cfg.Loader.SourceURL)

			if err != nil {
				log.Println("passport data source url is unreachable. error :" + err.Error())
				// try to connect next time
				continue
			}
			// if not modified since last update
			if fileLastModified.Before(lastUpdate) {
				// then nothing to do
				continue
			}
			log.Println("new version of the passport data is found")
		}

		log.Println("file dowloading is started...")
		err = downloadFile(cfg.Storage.PassportData, cfg.Loader.SourceURL)
		if err != nil {
			log.Println("passport data source url is unreachable. error :" + err.Error())
			if isNotExists {
				return false, err
			}
			// contunue regular checks for update
			continue
		}
		// successfully dowloaded
		return isNotExists, nil
	}

}

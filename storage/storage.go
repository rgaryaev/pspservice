package storage

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rgaryaev/pspservice/config"
	"github.com/rgaryaev/pspservice/storage/bitmap"
	"github.com/rgaryaev/pspservice/storage/loader"
	pilosadb "github.com/rgaryaev/pspservice/storage/pilosa"
)

//Storage - implementation storage of passport data
type Storage struct {
	engine Engine
}

//  local storage as a package var

func (ms *Storage) openPassportFile(fileName string) (*os.File, error) {
	var err error
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Unable to open file with passport list: "+fileName, err)
		return nil, err
	}

	return file, nil
}

// common checks and parsing series and number
func (ms *Storage) parseSeriesAndNumber(series string, number string) (uint16, uint32, bool) {

	//  if serias not in  [0:9999] or numbers not in [1:999999]
	if !(len(series) == 4 && len(number) == 6) {
		return 0, 0, false
	}
	// if series or number are not digits
	ser, err := strconv.ParseUint(series, 10, 16)
	if err != nil {
		return 0, 0, false
	}
	num, err := strconv.ParseUint(number, 10, 32)
	if err != nil {
		return 0, 0, false
	}

	if !(ser >= 0 && ser <= 9999) || !(num >= 1 && num <= 999999) {
		return 0, 0, false
	}
	return uint16(ser), uint32(num), true
}

func (ms *Storage) readPassportFile(fileName string) (uint32, error) {

	err := ms.engine.ImportData(fileName)
	// nor errors means that data was successfully imported
	if err == nil {
		return 0, nil
	}
	// continue to standart reading if no import

	file, err := ms.openPassportFile(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var total uint32

	// create new reader from csv file
	log.Println("loading passport data from the file...")
	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("unable to read the file with passport list", err)
			return 0, err
		}
		if series, number, ok := ms.parseSeriesAndNumber(record[0], record[1]); ok {
			ms.engine.AddPassport(uint16(series), uint32(number))
			total++
		}

	}

	log.Println("passport data has been loaded")
	return total, nil
}

// TestPassportFile - test loaded data
// for now this method is unused
func (ms *Storage) testPassportFile(fileName string) (uint32, error) {
	var (
		total uint32 = 0
		err   error
	)

	file, err := ms.openPassportFile(fileName)
	if err != nil {
		return 0, err
	}

	// create new reader from csv file
	log.Println("start reading file for testing")
	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("unable to read the file with passport list", err)
			return 0, err
		}
		if series, number, ok := ms.parseSeriesAndNumber(record[0], record[1]); ok {

			if exists, _ := ms.engine.CheckPassport(uint16(series), uint32(number)); !exists {
				total++
			}

		}

	}
	file.Close()
	return total, nil
}

func (ms *Storage) setEngine(engine Engine) {

	ms.engine = engine

}

//Implemenation itefraces for Storage

// StartStorage loads data from exising file
// if file doesnt exists then dowload it via UpdateStorage
func (ms *Storage) StartStorage(cfg *config.Configuration) error {

	var err error
	var engine string

	if cfg.Storage.Engine == strings.ToLower("sparse_bitmap") {

		ms.setEngine(new(bitmap.SparseBitmap))
		engine = "sparse bitmap"

	} else if cfg.Storage.Engine == strings.ToLower("roaring_bitmap") {

		ms.setEngine(new(bitmap.RoarBitmap))
		engine = "roaring bitmap"
	} else if cfg.Storage.Engine == strings.ToLower("pilosa") {

		ms.setEngine(new(pilosadb.PilosaDB))
		engine = "pilosa db"
	} else {
		log.Println("unknown storage engine")
		return err
	}

	ms.engine.Init()
	log.Println("activated storage engine:  " + engine)

	total, err := ms.readPassportFile(cfg.Storage.PassportData)
	if os.IsNotExist(err) {
		err = ms.UpdateStorage(cfg)
	} else if err != nil {
		return err
	} else {
		log.Println("Total records loaded: " + strconv.FormatUint(uint64(total), 10))
	}
	// Testing
	/*
		total, err = ms.testPassportFile(cfg.Storage.PassportData)
		if err != nil {
			return err
		}
		log.Println("Total errors: " + strconv.FormatUint(uint64(total), 10))
	*/
	return nil
}

// CheckStorage - self test memory storage. Usually after update procedure
// the function returns true if everythig is correct
func (ms *Storage) CheckStorage(cfg *config.Configuration) (bool, error) {
	//testFileName := config[testFileName]
	var countTestError uint32
	var err error
	fileName := cfg.Storage.PassportData
	countTestError, err = ms.testPassportFile(fileName)
	if err != nil {
		return false, err
	}
	//  return false if any error of testing
	if countTestError > 0 {
		return false, nil
	}

	return true, nil
}

// IsPassportInList - check is the passport number in the list
func (ms *Storage) IsPassportInList(series string, number string) (bool, error) {

	if ms.engine == nil {
		return false, errors.New("passport storage is not initialized")
	}

	// return true if incorect format or not parsed
	ser, num, ok := ms.parseSeriesAndNumber(series, number)
	if !ok {
		return true, nil
	}

	// it is need to lock in case of updates
	return ms.engine.CheckPassport(ser, num)
}

// UpdateStorage implements regular update passport data
// with simple scheduling
func (ms *Storage) UpdateStorage(cfg *config.Configuration) error {
	for {
		// call and wait for result
		// isNotExists means that file didn't  exist before dowloading
		//
		isNotExists, err := loader.LoadPassportData(cfg)
		if err != nil {
			// other sort of error
			msg := "Error during downloading passport data.  "
			log.Println(msg + err.Error())
			return errors.New(msg)
		}

		// Run update procedure
		log.Println("update is started...")
		// File exists and we are ready to start update
		total, err := ms.readPassportFile(cfg.Storage.PassportData)
		if err != nil {
			log.Println("update procedure of the passport data is failed. error :" + err.Error())
		}
		log.Println("Total records loaded: " + strconv.FormatUint(uint64(total), 10))

		// fix last update date
		lastUpdate := time.Now()
		// Save last update date
		cfg.Loader.LastUpdate = lastUpdate.Format(config.ConfigDateFormat)
		err = config.SaveConfiguration(cfg)
		//  jast in case - long update period is not allowed
		if err != nil {
			return errors.New("last update date is not saved in the cofig file")
		}

		log.Println("updater is finished. passport data is up to date.")
		runtime.GC()

		if isNotExists {
			// there is no need regular update but only download the file
			// return to the caller
			return nil
		}

	}
}

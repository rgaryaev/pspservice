package config

//  Define structure of config file

import (
	"encoding/json"
	"os"
)

//ConfigDateFormat represents format for YYYY-MM-DD
const ConfigDateFormat string = "2006-01-02"

const configFileName string = "config.json"

// Configuration structure of the config file
type Configuration struct {
	Storage struct {
		PassportData     string `json:"passport_data"`
		Engine           string `json:"engine"`
		NumberOfTests    int    `json:"number_of_tests"`
		TestPassportData string `json:"test_passport_data"`
	} `json:"storage"`
	Listener struct {
		Address               string `json:"address"`
		Port                  string `json:"port"`
		MaxPassportPerRequest uint   `json:"max_passport_per_request"`
	} `json:"listener"`
	Loader struct {
		SourceURL  string `json:"source_url"`
		EveryXDay  int    `json:"every_x_day"`
		LastUpdate string `json:"last_update"`
	} `json:"loader"`
}

// LoadConfiguration  load config data from the config file
func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	configFile, err := os.Open(configFileName)
	defer configFile.Close()
	if err != nil {
		return &config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)

	if err != nil {
		return &config, err
	}
	return &config, nil
}

// SaveConfiguration save config data to the config file
func SaveConfiguration(cfg *Configuration) error {

	configFile, err := os.Create(configFileName)
	defer configFile.Close()
	if err != nil {
		return err
	}
	json := json.NewEncoder(configFile)
	json.SetIndent("", "\t")
	err = json.Encode(cfg)
	if err != nil {
		return err
	}
	return nil
}

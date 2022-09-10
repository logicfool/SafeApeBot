package config

import (
	"SafeApeBot/structs"
	"encoding/json"
	"os"
)

func Read_config_from_file(f string) (*structs.SConfig, error) {
	// if config file does not exist create one
	file, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := structs.SConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return &config, nil
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Imgur struct {
		ClientID string `json:"client_id"`
	} `json:"imgur"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	defer file.Close()

	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	config := Config{}
	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}

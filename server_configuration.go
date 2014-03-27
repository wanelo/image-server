package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type ServerConfiguration struct {
	ServerPort            string   `json:"server_port"`
	StatusPort            string   `json:"status_port"`
	SourceDomain          string   `json:"source_domain"`
	WhitelistedExtensions []string `json:"whitelisted_extensions"`
	MaximumWidth          int      `json:"maximum_width"`
	Environment           string
}

func loadServerConfiguration(environment string) (*ServerConfiguration, error) {
	path := "config/" + environment + ".json"
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("configuration error: %v\n", err)
	}

	var config *ServerConfiguration
	json.Unmarshal(configFile, &config)

	log.Printf(" Config: %v", config)
	return config, nil
}

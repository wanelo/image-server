package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ServerConfiguration struct {
	ServerPort            string   `json:"server_port"`
	StatusPort            string   `json:"status_port"`
	SourceDomain          string   `json:"source_domain"`
	WhitelistedExtensions []string `json:"whitelisted_extensions"`
	MaximumWidth          int      `json:"maximum_width"`
	MantaBasePath         string   `json:"manta_base_path"`
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
	config.Environment = environment

	return config, nil
}

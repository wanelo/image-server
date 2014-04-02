package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/marpaia/graphite-golang"
)

type ServerConfiguration struct {
	ServerPort            string   `json:"server_port"`
	StatusPort            string   `json:"status_port"`
	SourceDomain          string   `json:"source_domain"`
	WhitelistedExtensions []string `json:"whitelisted_extensions"`
	MaximumWidth          int      `json:"maximum_width"`
	MantaBasePath         string   `json:"manta_base_path"`
	DefaultQuality        uint     `json:"default_quality"`
	GraphiteEnabled       bool     `json:"graphite_enabled"`
	GraphiteHost          string   `json:"graphite_host"`
	GraphitePort          int      `json:"graphite_port"`
	Graphite              *graphite.Graphite
	Environment           string
	Events                *EventChannels
	DataStore             *MantaAdapter
}

type EventChannels struct {
	ImageProcessed              chan *ImageConfiguration
	ImageProcessedWithErrors    chan *ImageConfiguration
	OriginalDownloaded          chan *ImageConfiguration
	OriginalDownloadUnavailable chan *ImageConfiguration
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
	config.Events = &EventChannels{
		ImageProcessed:     make(chan *ImageConfiguration),
		OriginalDownloaded: make(chan *ImageConfiguration),
	}
	return config, nil
}

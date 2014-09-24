package info

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ImageDetails struct {
	Hash   string `json:"hash"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

// ImageDetailsToJSON returns a string with JSON representation of the ImageDetails
func ImageDetailsToJSON(d *ImageDetails) (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// SaveImageDetail saves ImageDetails in a JSON file on "path"
func SaveImageDetail(d *ImageDetails, path string) error {
	json, err := ImageDetailsToJSON(d)
	if err != nil {
		log.Println(err)
		return err
	}

	d1 := []byte(json)
	err = ioutil.WriteFile(path, d1, 0644)

	if err != nil {
		log.Println(err)
	}
	return err
}

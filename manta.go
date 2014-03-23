package main

import (
	"log"
	"os"
)

var (
	MANTA_USER   string
	MANTA_URL    string
	MANTA_KEY_ID string

	mantaClient *Client
)

func InitializeManta() {
	MANTA_USER = os.Getenv("MANTA_USER")
	MANTA_URL = os.Getenv("MANTA_USER")
	MANTA_KEY_ID = os.Getenv("MANTA_USER")

	mantaClient = DefaultClient()
}

func (ic *ImageConfiguration) CreateMantaDirectory() {

	/*	dir := "../" + ic.DestinationDirectory()*/
	dir := "public/images"

	resp, err := mantaClient.Put(dir, os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

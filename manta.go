package main

import (
	"github.com/joyent/gocommon/client"
	"github.com/joyent/gocommon/jpc"
	"github.com/joyent/gomanta/manta"
	"log"
)

type Manta struct {
	Client *manta.Client
}

var mantaConfig Manta

func InitializeManta() {
	mantaConfig.Client = NewMantaClient()
}

func NewMantaClient() *manta.Client {
	creds, err := jpc.CompleteCredentialsFromEnv("")
	if err != nil {
		log.Fatalf("Error reading credentials for manta: %s", err.Error())
	}

	client := client.NewClient(creds.MantaEndpoint.URL, "", creds, &manta.Logger)
	return manta.New(client)
}

func CreateMantaDirectory() {
	opts := manta.ListDirectoryOpts{100, ""}
	files, err := mantaConfig.Client.ListDirectory("", opts)
	if err != nil {
		log.Fatalf("Error listing directory on manta: %s", err.Error())
	}
	log.Printf("Files: %v", files)

	/*	dir := "test_directory"
		err = mantaClient.PutDirectory(dir)
		if err != nil {
			log.Fatalf("Error creating directory on manta: %s", err.Error())
		}

		dir = "../" + ic.DestinationDirectory()

		log.Printf("Creating directory on manta: %s", dir)

		err = mantaClient.PutDirectory(dir)
		if err != nil {
			log.Fatalf("Error creating directory on manta: %s", err.Error())
		}*/
}

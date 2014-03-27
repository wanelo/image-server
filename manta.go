package main

import (
	"log"

	"github.com/richardiux/gocommon/client"
	"github.com/richardiux/gocommon/jpc"
	"github.com/richardiux/gomanta/manta"
)

type Manta struct {
	Client *manta.Client
}

var mantaConfig Manta

func initializeManta() {
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

func (ic *ImageConfiguration) CreateMantaDirectory() {
	dir := "public/images"
	log.Printf("Creating directory on manta: %s", dir)

	err := mantaConfig.Client.PutDirectory(dir)
	if err != nil {
		log.Fatalf("Error creating directory on manta: %s", err.Error())
	}
}

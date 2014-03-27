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
	mantaConfig.Client = newMantaClient()

	go func() {
		ensureBasePath()
	}()
}

func sendToManta(source string, destination string) {

}

func ensureBasePath() {
	baseDir := serverConfiguration.MantaBasePath
	createMantaDirectory(baseDir)
}

func newMantaClient() *manta.Client {
	creds, err := jpc.CompleteCredentialsFromEnv("")
	if err != nil {
		log.Fatalf("Error reading credentials for manta: %s", err.Error())
	}

	client := client.NewClient(creds.MantaEndpoint.URL, "", creds, &manta.Logger)
	return manta.New(client)
}

func createMantaDirectory(path string) {
	err := mantaConfig.Client.PutDirectory(path)
	if err != nil {
		log.Fatalf("Error creating directory on manta: %s", err.Error())
	}
	log.Printf("Created directory on manta: %s", path)
}

package main

import (
	"log"
	"path/filepath"

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
		ensureManataBasePath()
	}()
}

func sendToManta(source string, destination string) {
	err := ensureMantaImageDirectory(destination)
	if err != nil {
		return
	}
}

func ensureMantaImageDirectory(destination string) err {
	dir := filepath.Dir(destination)
	err := createMantaDirectory(dir)
	if err != nil {
		//  need to create sub directories
		dir2 := filepath.Dir(dir)
		dir3 := filepath.Dir(dir2)
		err = createMantaDirectory(dir3)
		err = createMantaDirectory(dir2)
		err = createMantaDirectory(dir)
	}
}

func ensureMantaBasePath() {
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

func createMantaDirectory(path string) error {
	err := mantaConfig.Client.PutDirectory(path)
	if err != nil {
		log.Printf("Error creating directory on manta: %s", path)
		return err
	}
	log.Printf("Created directory on manta: %s", path)
	return nil
}

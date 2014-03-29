package main

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"

	"github.com/joyent/gocommon/client"
	"github.com/joyent/gocommon/jpc"
	"github.com/richardiux/gomanta/manta"
)

type Manta struct {
	Client *manta.Client
}

var mantaConfig Manta

func initializeManta() {
	mantaConfig.Client = newMantaClient()

	go func() {
		ensureMantaBasePath()
	}()
}

func sendToManta(source string, destination string) {
	path, objectName := path.Split(destination)
	err := ensureMantaImageDirectory(path)
	if err != nil {
		log.Printf("Manta::sentToManta unable to create directory %s", path)
		return
	}
	object, err := ioutil.ReadFile(source)
	if err != nil {
		log.Printf("Manta::sentToManta unable to read file %s", source)
		return
	}
	err = mantaConfig.Client.PutObject(path, objectName, object)
	if err != nil {
		log.Printf("Error uploading image to manta: %s", err)
	}

}

func ensureMantaImageDirectory(dir string) error {
	err := createMantaDirectory(dir)
	if err != nil {
		//  need to create sub directories
		dir2 := filepath.Dir(dir)
		dir3 := filepath.Dir(dir2)
		err = createMantaDirectory(dir3)
		err = createMantaDirectory(dir2)
		err = createMantaDirectory(dir)
		if err != nil {
			return err
		}
	}
	return nil
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

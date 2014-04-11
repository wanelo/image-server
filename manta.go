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

type MantaAdapter struct {
	Client *manta.Client
}

func initializeManta(sc *ServerConfiguration) *MantaAdapter {
	m := &MantaAdapter{Client: newMantaClient()}
	m.ensureBasePath()
	return m
}

func newMantaClient() *manta.Client {
	creds, err := jpc.CompleteCredentialsFromEnv("")
	if err != nil {
		log.Fatalf("Error reading credentials for manta: %s", err.Error())
	}

	client := client.NewClient(creds.MantaEndpoint.URL, "", creds, &manta.Logger)
	return manta.New(client)
}

func (m *MantaAdapter) upload(ic *ImageConfiguration) {
	source := ic.LocalResizedImagePath()
	destination := ic.MantaResizedImagePath()

	path, objectName := path.Split(destination)
	err := m.ensureDirectory(path)
	if err != nil {
		log.Printf("Manta::sentToManta unable to create directory %s", path)
		return
	}
	object, err := ioutil.ReadFile(source)
	if err != nil {
		log.Printf("Manta::sentToManta unable to read file %s", source)
		return
	}
	err = m.Client.PutObject(path, objectName, object)
	if err != nil {
		log.Printf("Error uploading image to manta: %s", err)
	}
}

func (m *MantaAdapter) ensureDirectory(dir string) error {
	err := m.createDirectory(dir)
	if err != nil {
		//  need to create sub directories
		dir2 := filepath.Dir(dir)
		dir3 := filepath.Dir(dir2)
		err = m.createDirectory(dir3)
		if err != nil {
			return err
		}
		err = m.createDirectory(dir2)
		if err != nil {
			return err
		}
		err = m.createDirectory(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MantaAdapter) ensureBasePath() {
	baseDir := serverConfiguration.MantaBasePath
	m.createDirectory(baseDir)
}

func (m *MantaAdapter) createDirectory(path string) error {
	err := m.Client.PutDirectory(path)
	if err != nil {
		log.Printf("Error creating directory on manta: %s", path)
		return err
	}
	log.Printf("Created directory on manta: %s", path)
	return nil
}

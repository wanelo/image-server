package manta

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"

	"github.com/joyent/gocommon/client"
	"github.com/joyent/gocommon/jpc"
	m "github.com/richardiux/gomanta/manta"
	"github.com/wanelo/image-server/core"
)

type MantaAdapter struct {
	Client *m.Client
}

var serverConfiguration *core.ServerConfiguration

func InitializeManta(sc *core.ServerConfiguration) *MantaAdapter {
	serverConfiguration = sc
	ma := &MantaAdapter{Client: newMantaClient()}
	ma.ensureBasePath()
	return ma
}

func newMantaClient() *m.Client {
	creds, err := jpc.CompleteCredentialsFromEnv("")
	if err != nil {
		log.Fatalf("Error reading credentials for manta: %s", err.Error())
	}

	client := client.NewClient(creds.MantaEndpoint.URL, "", creds, &m.Logger)
	return m.New(client)
}

func (ma *MantaAdapter) Upload(ic *core.ImageConfiguration) {
	source := ic.LocalResizedImagePath()
	destination := serverConfiguration.MantaResizedImagePath(ic)

	path, objectName := path.Split(destination)
	err := ma.ensureDirectory(path)
	if err != nil {
		log.Printf("Manta::sentToManta unable to create directory %s", path)
		return
	}
	object, err := ioutil.ReadFile(source)
	if err != nil {
		log.Printf("Manta::sentToManta unable to read file %s", source)
		return
	}
	err = ma.Client.PutObject(path, objectName, object)
	log.Printf("Uploaded file to manta: %s", destination)

	if err != nil {
		log.Printf("Error uploading image to manta: %s", err)
	}
}

func (ma *MantaAdapter) ensureDirectory(dir string) error {
	err := ma.createDirectory(dir)
	if err != nil {
		//  need to create sub directories
		dir2 := filepath.Dir(dir)
		dir3 := filepath.Dir(dir2)
		dir4 := filepath.Dir(dir3)
		dir5 := filepath.Dir(dir4)
		err = ma.createDirectory(dir5)
		if err != nil {
			return err
		}
		err = ma.createDirectory(dir4)
		if err != nil {
			return err
		}
		err = ma.createDirectory(dir3)
		if err != nil {
			return err
		}
		err = ma.createDirectory(dir2)
		if err != nil {
			return err
		}
		err = ma.createDirectory(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ma *MantaAdapter) ensureBasePath() {
	baseDir := serverConfiguration.MantaBasePath
	ma.createDirectory(baseDir)
}

func (ma *MantaAdapter) createDirectory(path string) error {
	err := ma.Client.PutDirectory(path)
	if err != nil {
		log.Printf("Error creating directory on manta: %s", path)
		return err
	}
	log.Printf("Created directory on manta: %s", path)
	return nil
}

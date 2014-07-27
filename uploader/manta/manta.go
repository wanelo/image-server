package manta

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"

	"github.com/joyent/gocommon/client"
	"github.com/joyent/gocommon/jpc"
	m "github.com/richardiux/gomanta/manta"
)

type Uploader struct {
	Client *m.Client
}

func DefaultUploader() *Uploader {
	return &Uploader{
		Client: newMantaClient(),
	}
}

func (u *Uploader) Upload(source string, destination string) error {
	path, objectName := path.Split(destination)
	object, err := ioutil.ReadFile(source)
	if err != nil {
		log.Printf("Manta::sentToManta unable to read file %s", source)
		return err
	}
	err = u.Client.PutObject(path, objectName, object)
	log.Printf("Uploaded file to manta: %s", destination)

	if err != nil {
		log.Printf("Error uploading image to manta: %s", err)
	}
	return err
}

func (u *Uploader) CreateDirectory(dir string) error {
	err := u.createDirectory(dir)
	if err != nil {
		//  need to create sub directories
		dir2 := filepath.Dir(dir)
		dir3 := filepath.Dir(dir2)
		dir4 := filepath.Dir(dir3)
		dir5 := filepath.Dir(dir4)
		err = u.createDirectory(dir5)
		if err != nil {
			return err
		}
		err = u.createDirectory(dir4)
		if err != nil {
			return err
		}
		err = u.createDirectory(dir3)
		if err != nil {
			return err
		}
		err = u.createDirectory(dir2)
		if err != nil {
			return err
		}
		err = u.createDirectory(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func newMantaClient() *m.Client {
	creds, err := jpc.CompleteCredentialsFromEnv("")
	if err != nil {
		log.Fatalf("Error reading credentials for manta: %s", err.Error())
	}

	client := client.NewClient(creds.MantaEndpoint.URL, "", creds, &m.Logger)
	return m.New(client)
}

func (u *Uploader) createDirectory(path string) error {
	err := u.Client.PutDirectory(path)
	if err != nil {
		return err
	}
	log.Printf("Created directory on manta: %s", path)
	return nil
}

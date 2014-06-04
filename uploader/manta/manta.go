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

type Uploader struct {
	Client *m.Client
}

var serverConfiguration *core.ServerConfiguration

func InitializeUploader(sc *core.ServerConfiguration) *Uploader {
	serverConfiguration = sc
	u := &Uploader{Client: newMantaClient()}
	go u.ensureBasePath()
	return u
}

func newMantaClient() *m.Client {
	creds, err := jpc.CompleteCredentialsFromEnv("")
	if err != nil {
		log.Fatalf("Error reading credentials for manta: %s", err.Error())
	}

	client := client.NewClient(creds.MantaEndpoint.URL, "", creds, &m.Logger)
	return m.New(client)
}

func (u *Uploader) Upload(ic *core.ImageConfiguration) {
	source := ic.LocalResizedImagePath()
	destination := serverConfiguration.MantaResizedImagePath(ic)
	u.upload(source, destination)
}

func (u *Uploader) UploadOriginal(ic *core.ImageConfiguration) {
	source := ic.LocalOriginalImagePath()
	destination := serverConfiguration.MantaOriginalImagePath(ic)
	u.upload(source, destination)
}

func (u *Uploader) upload(source string, destination string) error {
	path, objectName := path.Split(destination)
	err := u.ensureDirectory(path)
	if err != nil {
		log.Printf("Manta::sentToManta unable to create directory %s", path)
		return err
	}
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

func (u *Uploader) ensureDirectory(dir string) error {
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

func (u *Uploader) ensureBasePath() {
	baseDir := serverConfiguration.MantaBasePath
	u.createDirectory(baseDir)
}

func (u *Uploader) createDirectory(path string) error {
	err := u.Client.PutDirectory(path)
	if err != nil {
		log.Printf("Error creating directory on manta: %s", path)
		return err
	}
	log.Printf("Created directory on manta: %s", path)
	return nil
}

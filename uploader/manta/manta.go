package manta

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	client "github.com/wanelo/image-server/uploader/manta/client"
)

type MantaClient interface {
	PutObject(path string, objectName string, object io.Reader) error
	PutDirectory(path string) error
}

type Uploader struct {
	Client MantaClient
}

func DefaultUploader() *Uploader {
	c := client.DefaultClient()

	return &Uploader{
		Client: c,
	}
}

func (u *Uploader) Upload(source string, destination string) error {
	path, objectName := path.Split(destination)
	fi, err := os.Open(source)
	if err != nil {
		log.Printf("Manta::sentToManta unable to read file %s, %s", source, err)
		return err
	}
	err = u.Client.PutObject(path, objectName, fi)

	if err != nil {
		log.Printf("Error uploading image to manta: %s", err)
	} else {
		log.Printf("Uploaded file to manta: %s", destination)
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

func (u *Uploader) createDirectory(path string) error {
	if path == "." {
		return nil
	}
	err := u.Client.PutDirectory(path)
	if err != nil {
		return err
	}
	log.Printf("Created directory on manta: %s", path)
	return nil
}

package manta

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	client "github.com/image-server/image-server/uploader/manta/client"
)

type MantaClient interface {
	PutObject(destination string, contentType string, object io.Reader) error
	PutDirectory(path string) error
}

type Uploader struct {
	Client MantaClient
}

var (
	MantaURL    string
	MantaUser   string
	MantaKeyID  string
	SDCIdentity string
)

func DefaultUploader() *Uploader {
	c := client.DefaultClient()

	return &Uploader{
		Client: c,
	}
}

func (u *Uploader) Upload(source string, destination string, contType string) error {
	log.Println("About to Upload:", source)
	fi, err := os.Open(source)

	if err != nil {
		glog.Infof("Manta::sentToManta unable to read file %s, %s", source, err)
		return err
	}

	// content type should be set depending of the type of file uploaded
	contentType := "application/octet-stream"
	err = u.Client.PutObject(destination, contentType, fi)

	if err != nil {
		glog.Infof("Error uploading image to manta: %s", err)
	} else {
		glog.Infof("Uploaded file to manta: %s", destination)
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

func (u *Uploader) ListDirectory(directory string) ([]string, error) {
	var names []string
	return names, nil
}

func Initialize(baseDir string, url string, user string, keyID string, identityPath string) error {
	u := DefaultUploader()
	MantaURL = url
	MantaUser = user
	MantaKeyID = keyID
	SDCIdentity = identityPath

	client.Initialize(MantaURL, MantaUser, MantaKeyID, SDCIdentity)

	return u.CreateDirectory(baseDir)
}

func (u *Uploader) createDirectory(path string) error {
	if path == "." {
		return nil
	}
	err := u.Client.PutDirectory(path)
	if err != nil {
		return err
	}
	glog.Infof("Created directory on manta: %s", path)
	return nil
}

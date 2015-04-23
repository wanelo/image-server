package uploader

import (
	"errors"
	"time"

	"github.com/golang/glog"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/uploader/manta"
	"github.com/wanelo/image-server/uploader/noop"
	"github.com/wanelo/image-server/uploader/s3"
)

type Uploader struct {
	BaseDir  string
	Uploader core.Uploader
}

// DefaultUploader returns an uploader that uses an adapter based on the configuration
// on the server configuration.
func DefaultUploader(sc *core.ServerConfiguration) *Uploader {
	u := &Uploader{
		BaseDir:  sc.RemoteBasePath,
		Uploader: &s3.Uploader{},
	}

	if sc.AWSAccessKeyID != "" {
		u.Uploader = &s3.Uploader{}
	} else if sc.MantaKeyID != "" {
		u.Uploader = manta.DefaultUploader()
	} else {
		u.Uploader = &noop.Uploader{}
	}
	return u
}

func (u *Uploader) Upload(source string, destination string, contType string) error {
	start := time.Now()

	if contType == "" {
		return errors.New("uploader.Upload is missing content type: " + source)
	}

	err := u.Uploader.Upload(source, destination, contType)
	elapsed := time.Since(start)
	glog.Infof("Took %s to upload image: %s", elapsed, destination)
	return err
}

func (u *Uploader) ListDirectory(directory string) ([]string, error) {
	return u.Uploader.ListDirectory(directory)
}

func (u *Uploader) CreateDirectory(path string) error {
	start := time.Now()
	directoryPath := u.Uploader.CreateDirectory(path)
	elapsed := time.Since(start)
	glog.Infof("Took %s to generate remote directory: %s", elapsed, path)
	return directoryPath
}

func Initialize(sc *core.ServerConfiguration) error {
	if sc.AWSAccessKeyID != "" {
		s3.Initialize(sc.AWSAccessKeyID, sc.AWSSecretKey, sc.AWSBucket)
	} else if sc.MantaKeyID != "" {
		manta.Initialize(sc.RemoteBasePath, sc.MantaURL, sc.MantaUser, sc.MantaKeyID, sc.SDCIdentity)
	}
	return nil
}

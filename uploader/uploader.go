package uploader

import (
	"log"
	"time"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/uploader/manta"
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
		u.Uploader = &s3.Uploader{
			AccessKey:  sc.AWSAccessKeyID,
			SecretKey:  sc.AWSSecretKey,
			BucketName: sc.AWSBucket,
		}
	} else {
		u.Uploader = manta.DefaultUploader(sc.RemoteBasePath)
	}
	return u
}

func (u *Uploader) Upload(source string, destination string, contType string) error {
	start := time.Now()
	err := u.Uploader.Upload(source, destination, contType)
	elapsed := time.Since(start)
	log.Printf("Took %s to upload image: %s", elapsed, destination)
	return err
}

func (u *Uploader) CreateDirectory(path string) error {
	start := time.Now()
	directoryPath := u.Uploader.CreateDirectory(path)
	elapsed := time.Since(start)
	log.Printf("Took %s to generate remote directory: %s", elapsed, path)
	return directoryPath
}

func (u *Uploader) Initialize() error {
	return u.Uploader.Initialize()
}

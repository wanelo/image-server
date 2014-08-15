package uploader

import (
	"log"
	"time"

	"github.com/wanelo/image-server/uploader/manta"
)

type Uploader struct {
	BaseDir  string
	Uploader *manta.Uploader
}

func DefaultUploader(baseDir string) *Uploader {
	return &Uploader{
		BaseDir:  baseDir,
		Uploader: manta.DefaultUploader(),
	}
}

func (u *Uploader) Upload(source string, destination string) error {
	start := time.Now()
	err := u.Uploader.Upload(source, destination)
	elapsed := time.Since(start)
	log.Printf("Took %s to upload image: %s", elapsed, destination)
	return err
}

func (u *Uploader) CreateDirectory(path string) error {
	start := time.Now()
	elapsed := time.Since(start)
	directoryPath := u.Uploader.CreateDirectory(path)
	log.Printf("Took %s to generate remote directory: %s", elapsed, path)
	return directoryPath
}

func (u *Uploader) Initialize() error {
	return u.Uploader.CreateDirectory(u.BaseDir)
}

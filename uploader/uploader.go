package uploader

import (
	"log"
	"time"

	"github.com/wanelo/image-server/uploader/manta"
)

type Uploader struct {
	BaseDir string
}

func (u *Uploader) Upload(source string, destination string) error {
	start := time.Now()
	uploader := manta.DefaultUploader()
	err := uploader.Upload(source, destination)
	elapsed := time.Since(start)
	log.Printf("Took %s to upload image: %s", elapsed, destination)
	return err
}

func (u *Uploader) CreateDirectory(path string) error {
	start := time.Now()
	uploader := manta.DefaultUploader()
	elapsed := time.Since(start)
	directoryPath := uploader.CreateDirectory(path)
	log.Printf("Took %s to generate remote directory: %s", elapsed, path)
	return directoryPath
}

func (u *Uploader) Initialize() error {
	uploader := manta.DefaultUploader()
	return uploader.CreateDirectory(u.BaseDir)
}

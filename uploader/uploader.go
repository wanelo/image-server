package uploader

import "github.com/wanelo/image-server/uploader/manta"

type Uploader struct {
	BaseDir string
}

func (u *Uploader) Upload(source string, destination string) error {
	uploader := manta.DefaultUploader()
	err := uploader.Upload(source, destination)
	return err
}

func (u *Uploader) CreateDirectory(path string) error {
	uploader := manta.DefaultUploader()
	return uploader.CreateDirectory(path)
}

func (u *Uploader) Initialize() error {
	uploader := manta.DefaultUploader()
	return uploader.CreateDirectory(u.BaseDir)
}

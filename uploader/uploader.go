package uploader

import "github.com/wanelo/image-server/core"

type Uploader struct {
	Uploader core.Uploader
}

func (u *Uploader) Upload(source string, destination string) error {
	err := u.Uploader.Upload(source, destination)
	return err
}

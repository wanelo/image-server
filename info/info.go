package info

import (
	"crypto/md5"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
)

type Info struct{}

type ImageDetails struct {
	Hash   string
	Height int
	Width  int
}

func (i Info) FileHash(path string) string {
	if contents, err := ioutil.ReadFile(path); err == nil {
		return fmt.Sprintf("%x", md5.Sum(contents))
	}
	return ""
}

func (i Info) ImageDetails(path string) (*ImageDetails, error) {
	if reader, err := os.Open(path); err == nil {
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			return nil, err
		}
		details := &ImageDetails{
			Hash:   i.FileHash(path),
			Height: im.Height,
			Width:  im.Width,
		}
		return details, nil
	} else {
		return nil, err
	}
}

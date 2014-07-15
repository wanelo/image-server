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

type Info struct {
	Path string
}

type ImageDetails struct {
	Hash   string
	Height int
	Width  int
}

func (i Info) FileHash() string {
	if contents, err := ioutil.ReadFile(i.Path); err == nil {
		return fmt.Sprintf("%x", md5.Sum(contents))
	}
	return ""
}

// ImageDetails extracts file hash, height, and width when providing a image path
// it returns an ImageDetails object
func (i Info) ImageDetails() (*ImageDetails, error) {
	if reader, err := os.Open(i.Path); err == nil {
		defer reader.Close()
		im, _, err := image.DecodeConfig(reader)
		if err != nil {
			return nil, err
		}
		details := &ImageDetails{
			Hash:   i.FileHash(),
			Height: im.Height,
			Width:  im.Width,
		}
		return details, nil
	} else {
		return nil, err
	}
}

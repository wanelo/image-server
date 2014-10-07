package info

import (
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	_ "code.google.com/p/go.image/webp"
	"github.com/wanelo/image-server/mime"
)

type Info struct {
	Path string
}

func (i Info) FileHash() (hash string, err error) {
	infile, err := os.Open(i.Path)
	if err != nil {
		return "", err
	}
	defer infile.Close()

	h := md5.New()
	io.Copy(h, infile)

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// ImageDetails extracts file hash, height, and width when providing a image path
// it returns an ImageDetails object
func (i Info) ImageDetails() (*ImageDetails, error) {
	if reader, err := os.Open(i.Path); err == nil {
		defer reader.Close()
		var contentType string
		var details *ImageDetails

		im, format, err := image.DecodeConfig(reader)
		if err != nil {
			return nil, err
		}
		if format != "" {
			contentType, err = getContentTypeFromExtension(format)
			if err != nil {
				return nil, err
			}

			details = &ImageDetails{
				Height:      im.Height,
				Width:       im.Width,
				ContentType: contentType,
			}
		} else {
			return nil, errors.New("Unable to extract the format of the image")
		}

		hash, err := i.FileHash()
		details.Hash = hash
		return details, nil

	} else {
		return nil, err
	}
}

func getContentTypeFromExtension(format string) (string, error) {
	if format == "" {
		return "", errors.New("Can't extract format")
	}

	contentType := mime.ExtToContentType(format)
	if contentType == "" {
		return "", fmt.Errorf("Can't extract content type from format. format=%s, contentType=%s", format, contentType)
	}

	return contentType, nil
}

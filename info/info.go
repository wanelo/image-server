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
	"log"
	"mime"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Info struct {
	Path string
}

func (i Info) FileHash() (hash string, err error) {
	infile, err := os.Open(i.Path)
	if err != nil {
		return "", err
	}
	h := md5.New()
	io.Copy(h, infile)

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// ImageDetails extracts file hash, height, and width when providing a image path
// it returns an ImageDetails object
func (i Info) ImageDetails() (*ImageDetails, error) {
	if reader, err := os.Open(i.Path); err == nil {
		defer reader.Close()
		im, format, err := image.DecodeConfig(reader)
		var contentType string
		var details *ImageDetails

		if err == nil {
			contentType, err = getContentTypeFromExtension(format)

			details = &ImageDetails{
				Height:      im.Height,
				Width:       im.Width,
				ContentType: contentType,
			}
		} else {
			// can't calculate content type, so will use ImageMagick as fallback
			// use fallback
			details, err = i.DetailsFromImageMagick()
			if err != nil {
				return nil, err
			}
		}

		hash, err := i.FileHash()
		details.Hash = hash
		return details, nil

	} else {
		return nil, err
	}
}

func (i Info) DetailsFromImageMagick() (*ImageDetails, error) {
	args := []string{"-format", "%[fx:w]:%[fx:h]:%m", i.Path}
	out, err := exec.Command("identify", args...).Output()
	dimensions := fmt.Sprintf("%s", out)
	dimensions = strings.TrimSpace(dimensions)

	if err != nil {
		return nil, err
	}

	d := strings.Split(dimensions, ":")
	w, err := strconv.Atoi(d[0])
	if err != nil {
		log.Printf("Can't convert width to integer: %s\n", d[0])
		return nil, err
	}

	h, err := strconv.Atoi(d[1])
	if err != nil {
		log.Printf("Can't convert height to integer: %s\n", d[1])
		return nil, err
	}

	contentType, err := getContentTypeFromExtension(d[2])
	if err != nil {
		return nil, err
	}

	return &ImageDetails{
		Height:      h,
		Width:       w,
		ContentType: contentType,
	}, nil
}

func getContentTypeFromExtension(format string) (string, error) {
	if format == "" {
		return "", errors.New("Can't extract format")
	}

	ext := strings.ToLower(fmt.Sprintf(".%s", format))
	return mime.TypeByExtension(ext), nil
}

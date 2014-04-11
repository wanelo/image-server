package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rainycape/magick"
)

func (ic *ImageConfiguration) createImage(sc *ServerConfiguration) (string, error) {
	if ic.width == 0 && ic.height == 0 {
		return createFullSizeImage(ic, sc)
	}

	resizedPath := ic.LocalResizedImagePath()
	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		c := make(chan error)
		go fetchOriginal(c, ic, sc)

		err = <-c
		if err != nil {
			log.Println(err)
			return "", err
		}

		err = createWithMagick(ic)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	return resizedPath, nil
}

func createWithMagick(ic *ImageConfiguration) error {
	start := time.Now()
	fullSizePath := ic.LocalOriginalImagePath()
	im, err := magick.DecodeFile(fullSizePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer im.Dispose()

	im2, err := im.CropResize(ic.width, ic.height, magick.FHamming, magick.CSCenter)
	if err != nil {
		log.Println(err)
		return err
	}

	resizedPath := ic.LocalResizedImagePath()
	dir := filepath.Dir(resizedPath)
	os.MkdirAll(dir, 0700)

	out, err := os.Create(resizedPath)
	defer out.Close()

	info := ic.MagickInfo()
	err = im2.Encode(out, info)

	if err != nil {
		log.Println(err)
		return err
	}
	elapsed := time.Since(start)
	log.Printf("Took %s to generate image: %s", elapsed, resizedPath)

	return nil
}

func createFullSizeImage(ic *ImageConfiguration, sc *ServerConfiguration) (string, error) {
	fullSizePath := ic.LocalOriginalImagePath()
	resizedPath := ic.LocalResizedImagePath()

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		downloadAndSaveOriginal(ic, sc)

		im, err := magick.DecodeFile(fullSizePath)
		if err != nil {
			log.Println(err)
			return "", err
		}
		defer im.Dispose()

		out, err := os.Create(resizedPath)
		defer out.Close()

		info := ic.MagickInfo()
		err = im.Encode(out, info)

		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	return resizedPath, nil
}

func (ic *ImageConfiguration) MagickInfo() *magick.Info {
	info := magick.NewInfo()
	info.SetQuality(ic.quality)
	info.SetFormat(ic.format)
	return info
}

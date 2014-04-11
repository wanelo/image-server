package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rainycape/magick"
)

type ImageProcessingResult struct {
	resizedPath string
	err         error
}

var imageProcessings map[string][]chan ImageProcessingResult

func (ic *ImageConfiguration) createImage(sc *ServerConfiguration) (string, error) {
	c := make(chan ImageProcessingResult)
	go ic.uniqueCreateImage(c, sc)
	ipr := <-c
	return ipr.resizedPath, ipr.err
}

func (ic *ImageConfiguration) uniqueCreateImage(c chan ImageProcessingResult, sc *ServerConfiguration) {
	key := ic.LocalResizedImagePath()
	_, present := imageProcessings[key]

	if present {
		imageProcessings[key] = append(imageProcessings[key], c)
	} else {
		imageProcessings[key] = []chan ImageProcessingResult{c}

		imagePath, err := ic.downloadAndProcessImage(sc)
		for _, cc := range imageProcessings[key] {
			cc <- ImageProcessingResult{imagePath, err}
		}
		delete(imageProcessings, key)
	}
}

func (ic *ImageConfiguration) downloadAndProcessImage(sc *ServerConfiguration) (string, error) {
	if ic.width == 0 && ic.height == 0 {
		return createFullSizeImage(ic, sc)
	}

	resizedPath := ic.LocalResizedImagePath()
	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {

		err = fetchOriginal(ic, sc)
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

		err = fetchOriginal(ic, sc)
		if err != nil {
			log.Println(err)
			return "", err
		}

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

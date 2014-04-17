package magick

import (
	"log"
	"os"
	"path/filepath"
	"time"

	m "github.com/rainycape/magick"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher/http"
)

type Processor struct {
	ServerConfiguration *core.ServerConfiguration
}

func (p *Processor) CreateImage(ic *core.ImageConfiguration) (string, error) {
	c := make(chan ImageProcessingResult)
	go uniqueCreateImage(c, p.ServerConfiguration, ic)
	ipr := <-c
	return ipr.resizedPath, ipr.err
}

type ImageProcessingResult struct {
	resizedPath string
	err         error
}

var ImageProcessings map[string][]chan ImageProcessingResult

func uniqueCreateImage(c chan ImageProcessingResult, sc *core.ServerConfiguration, ic *core.ImageConfiguration) {
	key := ic.LocalResizedImagePath()
	_, present := ImageProcessings[key]

	if present {
		ImageProcessings[key] = append(ImageProcessings[key], c)
	} else {
		ImageProcessings[key] = []chan ImageProcessingResult{c}

		imagePath, err := downloadAndProcessImage(sc, ic)
		for _, cc := range ImageProcessings[key] {
			cc <- ImageProcessingResult{imagePath, err}
		}
		delete(ImageProcessings, key)
	}
}

func downloadAndProcessImage(sc *core.ServerConfiguration, ic *core.ImageConfiguration) (string, error) {
	if ic.Width == 0 && ic.Height == 0 {
		return createFullSizeImage(ic, sc)
	}

	resizedPath := ic.LocalResizedImagePath()
	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {

		err = http.FetchOriginal(ic, sc)
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

func createWithMagick(ic *core.ImageConfiguration) error {
	start := time.Now()
	fullSizePath := ic.LocalOriginalImagePath()
	im, err := m.DecodeFile(fullSizePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer im.Dispose()

	im2, err := im.CropResize(ic.Width, ic.Height, m.FHamming, m.CSCenter)
	if err != nil {
		log.Println(err)
		return err
	}

	resizedPath := ic.LocalResizedImagePath()
	dir := filepath.Dir(resizedPath)
	os.MkdirAll(dir, 0700)

	out, err := os.Create(resizedPath)
	defer out.Close()

	info := magickInfo(ic)
	err = im2.Encode(out, info)

	if err != nil {
		log.Println(err)
		return err
	}
	elapsed := time.Since(start)
	log.Printf("Took %s to generate image: %s", elapsed, resizedPath)

	return nil
}

func createFullSizeImage(ic *core.ImageConfiguration, sc *core.ServerConfiguration) (string, error) {
	fullSizePath := ic.LocalOriginalImagePath()
	resizedPath := ic.LocalResizedImagePath()

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {

		err = http.FetchOriginal(ic, sc)
		if err != nil {
			log.Println(err)
			return "", err
		}

		im, err := m.DecodeFile(fullSizePath)
		if err != nil {
			log.Println(err)
			return "", err
		}
		defer im.Dispose()

		out, err := os.Create(resizedPath)
		defer out.Close()

		info := magickInfo(ic)
		err = im.Encode(out, info)

		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	return resizedPath, nil
}

func magickInfo(ic *core.ImageConfiguration) *m.Info {
	info := m.NewInfo()
	info.SetQuality(ic.Quality)
	info.SetFormat(ic.Format)
	return info
}

package native

import (
	"image"
	"image/jpeg"
	"log"
	"os"

	"code.google.com/p/graphics-go/graphics"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/processor"
)

type Processor struct {
	ServerConfiguration *core.ServerConfiguration
}

func (p *Processor) CreateImage(ic *core.ImageConfiguration) (string, error) {
	c := make(chan processor.ImageProcessingResult)
	go uniqueCreateImage(c, p.ServerConfiguration, ic)
	ipr := <-c
	return ipr.ResizedPath, ipr.Error
}

func uniqueCreateImage(c chan processor.ImageProcessingResult, sc *core.ServerConfiguration, ic *core.ImageConfiguration) {
	key := ic.LocalResizedImagePath()
	_, present := processor.ImageProcessings[key]

	if present {
		processor.ImageProcessings[key] = append(processor.ImageProcessings[key], c)
	} else {
		processor.ImageProcessings[key] = []chan processor.ImageProcessingResult{c}

		imagePath, err := downloadAndProcessImage(sc, ic)
		for _, cc := range processor.ImageProcessings[key] {
			cc <- processor.ImageProcessingResult{imagePath, err}
		}
		delete(processor.ImageProcessings, key)
		go func() {
			ic.ServerConfiguration.Events.ImageProcessed <- ic
		}()
	}
}

func downloadAndProcessImage(sc *core.ServerConfiguration, ic *core.ImageConfiguration) (string, error) {
	resizedPath := ic.LocalResizedImagePath()
	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {

		err = http.FetchOriginal(ic)
		if err != nil {
			log.Println(err)
			return "", err
		}

		err = createResizedImage(ic)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}

	return resizedPath, nil
}

func createResizedImage(ic *core.ImageConfiguration) error {
	fullSizePath := ic.LocalOriginalImagePath()
	resizedPath := ic.LocalResizedImagePath()

	file, err := os.Open(fullSizePath)
	if err != nil {
		return err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		dst := image.NewRGBA(image.Rect(0, 0, ic.Width, ic.Height))
		graphics.Thumbnail(dst, img)

		toimg, err := os.Create(resizedPath)
		if err != nil {
			return err
		}
		defer toimg.Close()

		quality := int(ic.Quality)
		jpeg.Encode(toimg, dst, &jpeg.Options{quality})
		return err
	}
	return nil
}

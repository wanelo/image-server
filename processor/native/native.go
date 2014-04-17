package native

import (
	"image"
	"image/jpeg"
	"os"

	"code.google.com/p/graphics-go/graphics"
	"github.com/wanelo/image-server/core"
)

type Processor struct {
	ServerConfiguration *core.ServerConfiguration
}

func (p *Processor) CreateImage(ic *core.ImageConfiguration) (string, error) {
	err := createResizedImage(p.ServerConfiguration, ic)
	return ic.LocalResizedImagePath(), err
}

func createResizedImage(sc *core.ServerConfiguration, ic *core.ImageConfiguration) error {
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

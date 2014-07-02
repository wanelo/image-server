package native

import (
	"image"
	"image/jpeg"
	"os"

	"code.google.com/p/graphics-go/graphics"
	"github.com/wanelo/image-server/core"
)

type Processor struct{}

func (p *Processor) CreateImage(source string, destination string, ic *core.ImageConfiguration) error {

	file, err := os.Open(source)
	if err != nil {
		return err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	dst := image.NewRGBA(image.Rect(0, 0, ic.Width, ic.Height))
	graphics.Thumbnail(dst, img)

	toimg, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer toimg.Close()

	quality := int(ic.Quality)
	jpeg.Encode(toimg, dst, &jpeg.Options{quality})
	return nil
}

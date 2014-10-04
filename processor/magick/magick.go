// +build !solaris

package magick

import (
	"log"
	"os"
	"path/filepath"

	"github.com/wanelo/image-server/core"
	m "gopkgs.com/magick.v1"
)

const Available = true

type Processor struct{}

func (p *Processor) CreateImage(source string, destination string, ic *core.ImageConfiguration) error {
	if ic.Width == 0 && ic.Height == 0 {
		return createFullSizeImage(source, destination, ic)
	}

	return createResizedImage(source, destination, ic)
}

func createResizedImage(source string, destination string, ic *core.ImageConfiguration) error {
	im, err := m.DecodeFile(source)
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

	dir := filepath.Dir(destination)
	os.MkdirAll(dir, 0700)

	out, err := os.Create(destination)
	defer out.Close()

	info := magickInfo(ic)
	err = im2.Encode(out, info)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func createFullSizeImage(source string, destination string, ic *core.ImageConfiguration) error {
	im, err := m.DecodeFile(source)
	if err != nil {
		log.Println(err)
		return err
	}
	defer im.Dispose()

	out, err := os.Create(destination)
	defer out.Close()

	info := magickInfo(ic)
	err = im.Encode(out, info)

	return err
}

func magickInfo(ic *core.ImageConfiguration) *m.Info {
	info := m.NewInfo()
	info.SetQuality(ic.Quality)
	info.SetFormat(ic.Format)
	return info
}

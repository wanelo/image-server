package cli

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/processor"
)

type Processor struct {
	ServerConfiguration *core.ServerConfiguration
}

func (p *Processor) CreateImage(ic *core.ImageConfiguration) (string, error) {
	c := make(chan processor.ImageProcessingResult)
	go uniqueCreateImage(c, ic)
	ipr := <-c
	return ipr.ResizedPath, ipr.Error
}

func uniqueCreateImage(c chan processor.ImageProcessingResult, ic *core.ImageConfiguration) {
	key := ic.LocalResizedImagePath()
	_, present := processor.ImageProcessings[key]

	if present {
		processor.ImageProcessings[key] = append(processor.ImageProcessings[key], c)
	} else {
		processor.ImageProcessings[key] = []chan processor.ImageProcessingResult{c}

		imagePath, err := downloadAndProcessImage(ic)
		log.Println(imagePath)
		for _, cc := range processor.ImageProcessings[key] {
			cc <- processor.ImageProcessingResult{imagePath, err}
		}
		delete(processor.ImageProcessings, key)
	}
}

func downloadAndProcessImage(ic *core.ImageConfiguration) (string, error) {
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

	cmd := exec.Command("convert", commandArgs(ic)...)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func commandArgs(ic *core.ImageConfiguration) []string {
	args := list.New()

	/*	args.PushBack("convert")*/

	args.PushBack("-format")
	args.PushBack(ic.Format)

	args.PushBack("-flatten")

	args.PushBack("-background")
	args.PushBack("rgba\\(255,255,255,1\\)")

	args.PushBack("-quality")
	args.PushBack(fmt.Sprintf("%d", ic.Quality))

	if ic.Height > 0 && ic.Width > 0 {
		args.PushBack("-extent")
		args.PushBack(fmt.Sprintf("%dx%d", ic.Width, ic.Height))
	} else if ic.Width > 0 {
		args.PushBack("-resize")
		args.PushBack(fmt.Sprintf("%d", ic.Width))
	}

	args.PushBack(ic.LocalOriginalImagePath())
	args.PushBack(ic.LocalResizedImagePath())

	return convertArgumentsToSlice(args)
}

func convertArgumentsToSlice(arguments *list.List) []string {
	argumentSlice := make([]string, 0, arguments.Len())
	for e := arguments.Front(); e != nil; e = e.Next() {
		argumentSlice = append(argumentSlice, e.Value.(string))
	}
	return argumentSlice
}

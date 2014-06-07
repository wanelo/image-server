package cli

import (
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wanelo/image-server/core"
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
		go func() {
			ic.ServerConfiguration.Events.ImageProcessed <- ic
		}()
	}
}

func downloadAndProcessImage(ic *core.ImageConfiguration) (string, error) {
	resizedPath := ic.LocalResizedImagePath()
	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {

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

	args.PushBack("-format")
	args.PushBack(ic.Format)

	args.PushBack("-flatten")

	if ic.Height > 0 && ic.Width > 0 {
		cols, rows, err := originalDimensions(ic)

		if err == nil && (ic.Width != cols || ic.Height != rows) {
			w := float64(ic.Width) / float64(cols)
			h := float64(ic.Height) / float64(rows)
			scale := math.Max(w, h)
			c := scale * (float64(cols) + 0.5)
			c = math.Floor(c + 0.5) // Round
			r := scale * (float64(rows) + 0.5)
			r = math.Floor(r + 0.5) // Round

			resizeVal := fmt.Sprintf("%dx%d", int(c), int(r))

			args.PushBack("-resize")
			args.PushBack(resizeVal)
		}

		args.PushBack("-extent")
		args.PushBack(fmt.Sprintf("%dx%d", ic.Width, ic.Height))

		args.PushBack("-gravity")
		args.PushBack("center")

	} else if ic.Width > 0 {
		args.PushBack("-resize")
		args.PushBack(fmt.Sprintf("%d", ic.Width))
	}

	args.PushBack("-background")
	args.PushBack("rgba\\(255,255,255,1\\)")

	args.PushBack("-quality")
	args.PushBack(fmt.Sprintf("%d", ic.Quality))

	args.PushBack(ic.LocalOriginalImagePath())
	args.PushBack(ic.LocalResizedImagePath())

	return convertArgumentsToSlice(args)
}

func originalDimensions(ic *core.ImageConfiguration) (int, int, error) {
	args := []string{"-format", "%[fx:w]x%[fx:h]", ic.LocalOriginalImagePath()}
	out, err := exec.Command("identify", args...).Output()
	dimensions := fmt.Sprintf("%s", out)

	if err != nil {
		return 0, 0, err
	}

	d := strings.Split(dimensions, "x")
	w, err := strconv.Atoi(d[0])
	if err != nil {
		fmt.Printf("Can't convert width to integer: %s\n", d[0])
		return 0, 0, err
	}

	h, err := strconv.Atoi(d[1])
	if err != nil {
		fmt.Printf("Can't convert height to integer: %s\n", d[1])
		return 0, 0, err
	}

	return w, h, nil
}

func convertArgumentsToSlice(arguments *list.List) []string {
	argumentSlice := make([]string, 0, arguments.Len())
	for e := arguments.Front(); e != nil; e = e.Next() {
		argumentSlice = append(argumentSlice, e.Value.(string))
	}
	return argumentSlice
}

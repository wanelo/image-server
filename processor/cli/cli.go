package cli

import (
	"container/list"
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wanelo/image-server/core"
)

type Processor struct{}

func (p *Processor) CreateImage(source string, destination string, ic *core.ImageConfiguration) error {
	cmd := exec.Command("convert", p.commandArgs(source, destination, ic)...)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) commandArgs(source string, destination string, ic *core.ImageConfiguration) []string {
	args := list.New()

	args.PushBack("-format")
	args.PushBack(ic.Format)

	args.PushBack("-flatten")

	if ic.Height > 0 && ic.Width > 0 {
		cols, rows, err := p.originalDimensions(source, ic)

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
	args.PushBack("rgba(255,255,255,1)")

	args.PushBack("-quality")
	args.PushBack(fmt.Sprintf("%d", ic.Quality))

	args.PushBack(source)
	args.PushBack(destination)

	return p.convertArgumentsToSlice(args)
}

func (p *Processor) originalDimensions(source string, ic *core.ImageConfiguration) (int, int, error) {
	args := []string{"-format", "%[fx:w]x%[fx:h]", source}
	out, err := exec.Command("identify", args...).Output()
	dimensions := fmt.Sprintf("%s", out)
	dimensions = strings.TrimSpace(dimensions)

	if err != nil {
		return 0, 0, err
	}

	d := strings.Split(dimensions, "x")
	w, err := strconv.Atoi(d[0])
	if err != nil {
		log.Printf("Can't convert width to integer: %s\n", d[0])
		return 0, 0, err
	}

	h, err := strconv.Atoi(d[1])
	if err != nil {
		log.Printf("Can't convert height to integer: %s\n", d[1])
		return 0, 0, err
	}

	return w, h, nil
}

func (p *Processor) convertArgumentsToSlice(arguments *list.List) []string {
	argumentSlice := make([]string, 0, arguments.Len())
	for e := arguments.Front(); e != nil; e = e.Next() {
		argumentSlice = append(argumentSlice, e.Value.(string))
	}
	return argumentSlice
}

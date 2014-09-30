package cli

import (
	"container/list"
	"fmt"
	"math"
	"os/exec"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/info"
)

type Processor struct {
	ImageDetails       *info.ImageDetails
	ImageConfiguration *core.ImageConfiguration
	Source             string
	Destination        string
}

func (p *Processor) CreateImage() error {
	cmd := exec.Command("convert", p.commandArgs()...)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) commandArgs() []string {
	ic := p.ImageConfiguration
	source := p.Source
	destination := p.Destination

	args := list.New()

	args.PushBack("-format")
	args.PushBack(ic.Format)

	args.PushBack("-flatten")

	if ic.Height > 0 && ic.Width > 0 {
		cols := p.ImageDetails.Width
		rows := p.ImageDetails.Height

		if ic.Width != cols || ic.Height != rows {
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

func (p *Processor) convertArgumentsToSlice(arguments *list.List) []string {
	argumentSlice := make([]string, 0, arguments.Len())
	for e := arguments.Front(); e != nil; e = e.Next() {
		argumentSlice = append(argumentSlice, e.Value.(string))
	}
	return argumentSlice
}

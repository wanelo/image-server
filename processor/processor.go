package processor

import (
	"log"
	"os"
	"time"

	"github.com/wanelo/image-server/core"
	adapter "github.com/wanelo/image-server/processor/cli"
)

type ProcessorResult struct {
	ResizedPath string
	Error       error
}

var ImageProcessings map[string][]chan ProcessorResult

func init() {
	ImageProcessings = make(map[string][]chan ProcessorResult)
}

type Processor struct {
	Source             string
	Destination        string
	ImageConfiguration *core.ImageConfiguration
	Channels           *ProcessorChannels
}

type ProcessorChannels struct {
	ImageProcessed chan *core.ImageConfiguration
	Skipped        chan string
}

func (p *Processor) CreateImage() (string, error) {
	c := make(chan ProcessorResult)
	go p.uniqueCreateImage(c)
	ipr := <-c
	return ipr.ResizedPath, ipr.Error
}

func (p *Processor) uniqueCreateImage(c chan ProcessorResult) {
	key := p.Destination
	_, present := ImageProcessings[key]

	if present {
		ImageProcessings[key] = append(ImageProcessings[key], c)
		p.notifySkipped(key)
	} else {
		ImageProcessings[key] = []chan ProcessorResult{c}
		err := p.createIfNotAvailable()

		for _, cc := range ImageProcessings[key] {
			cc <- ProcessorResult{p.Destination, err}
		}
		delete(ImageProcessings, key)
		p.notifyProcessed()
	}
}

func (p *Processor) createIfNotAvailable() error {
	if _, err := os.Stat(p.Destination); os.IsNotExist(err) {
		start := time.Now()

		processor := &adapter.Processor{}
		err = processor.CreateImage(p.Source, p.Destination, p.ImageConfiguration)

		if err != nil {
			log.Println(err)
			return err
		}

		elapsed := time.Since(start)
		log.Printf("Took %s to generate image: %s", elapsed, p.Destination)
		p.notifyProcessed()
	} else {
		p.notifySkipped(p.Destination)
	}

	return nil
}

func (p *Processor) notifyProcessed() {
	go func() {
		p.Channels.ImageProcessed <- p.ImageConfiguration
	}()
}

func (p *Processor) notifySkipped(path string) {
	go func() {
		p.Channels.Skipped <- path
	}()
}

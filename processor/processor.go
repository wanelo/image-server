package processor

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/info"
	adapter "github.com/wanelo/image-server/processor/cli"
)

type ProcessorResult struct {
	ResizedPath string
	Error       error
}

var ImageProcessings map[string][]chan ProcessorResult
var processingMutex sync.RWMutex // To protect ImageProcessings

func init() {
	ImageProcessings = make(map[string][]chan ProcessorResult)
}

type Processor struct {
	Source             string
	Destination        string
	ImageConfiguration *core.ImageConfiguration
	ImageDetails       *info.ImageDetails
	Channels           *ProcessorChannels
}

type ProcessorChannels struct {
	ImageProcessed chan *core.ImageConfiguration
	Skipped        chan string
}

func (p *Processor) CreateImage() error {
	if p.ImageDetails == nil {
		log.Panic("ImageDetails is required")
	}

	c := make(chan ProcessorResult)
	go p.uniqueCreateImage(c)
	ipr := <-c
	return ipr.Error
}

func (p *Processor) uniqueCreateImage(c chan ProcessorResult) {
	key := p.Destination
	_, present := ImageProcessings[key]

	processingMutex.Lock()

	if present {
		ImageProcessings[key] = append(ImageProcessings[key], c)
		processingMutex.Unlock()
		p.notifySkipped(key)
	} else {
		ImageProcessings[key] = []chan ProcessorResult{c}
		processingMutex.Unlock()

		err := p.createIfNotAvailable()

		for _, cc := range ImageProcessings[key] {
			cc <- ProcessorResult{p.Destination, err}
		}
		processingMutex.Lock()
		delete(ImageProcessings, key)
		processingMutex.Unlock()
		p.notifyProcessed()
	}

}

func (p *Processor) createIfNotAvailable() error {
	if _, err := os.Stat(p.Destination); os.IsNotExist(err) {
		start := time.Now()

		dir := filepath.Dir(p.Destination)
		os.MkdirAll(dir, 0700)

		processor := &adapter.Processor{
			Source:             p.Source,
			Destination:        p.Destination,
			ImageConfiguration: p.ImageConfiguration,
			ImageDetails:       p.ImageDetails,
		}
		err = processor.CreateImage()

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

package main

import (
	"flag"
	"log"
	"strings"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/processor"
	"github.com/wanelo/image-server/processor/cli"
	sm "github.com/wanelo/image-server/source_mapper/waneloS3"
)

type CliConfiguration struct {
	Namespace           string
	Outputs             []string
	Start               int
	End                 int
	ServerConfiguration *core.ServerConfiguration
	Concurrency         int
}

func extractCliConfiguration() *CliConfiguration {
	start := flag.Int("start", 0, "")
	end := flag.Int("end", 0, "")
	concurrency := flag.Int("concurrency", 20, "")
	environment := flag.String("e", "development", "Specifies the environment to run this server under (test/development/production).")

	flag.Parse()

	path := "config/" + *environment + ".json"
	serverConfiguration, err := core.LoadServerConfiguration(path)

	adapters := &core.Adapters{
		Processor:    &cli.Processor{serverConfiguration},
		SourceMapper: &sm.SourceMapper{serverConfiguration},
	}

	mappings := make(map[string]string)
	mappings["p"] = "product/image"
	serverConfiguration.NamespaceMappings = mappings

	serverConfiguration.Adapters = adapters

	http.ImageDownloads = make(map[string][]chan error)
	processor.ImageProcessings = make(map[string][]chan processor.ImageProcessingResult)

	if err != nil {
		log.Panicln(err)
	}

	return &CliConfiguration{
		Namespace:           *flag.String("namespace", "p", "Namespace of images. i.e. 'p'"),
		Outputs:             strings.Split(*flag.String("outputs", "", ""), ","),
		Start:               *start,
		End:                 *end,
		ServerConfiguration: serverConfiguration,
		Concurrency:         *concurrency,
	}

}

// Returns range of ids
func (c *CliConfiguration) ProductIds() ([]int, error) {

	var ids []int
	for i := c.Start; i <= c.End; i++ {
		ids = append(ids, i)
	}
	return ids, nil
}

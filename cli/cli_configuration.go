package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	config "github.com/wanelo/image-server/config/wanelo"
	"github.com/wanelo/image-server/core"
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
	var (
		start       = flag.Int("start", 0, "")
		end         = flag.Int("end", 0, "")
		concurrency = flag.Int("concurrency", 20, "")
		outputs     = flag.String("outputs", "x200.webp,x200.jpg", "Output files with dimension and compression: 'x300.jpg,x300.webp'")
	)

	flag.Parse()

	if *start == 0 {
		fmt.Println("Enter start range:")
		fmt.Scanf("%d", start)
	}

	if *end == 0 {
		fmt.Println("Enter end range:")
		fmt.Scanf("%d", end)
	}

	serverConfiguration, err := config.ServerConfiguration()

	if err != nil {
		log.Panicln(err)
	}

	return &CliConfiguration{
		Namespace:           *flag.String("namespace", "p", "Namespace of images. i.e. 'p'"),
		Start:               *start,
		End:                 *end,
		ServerConfiguration: serverConfiguration,
		Concurrency:         *concurrency,
		Outputs:             strings.Split(*outputs, ","),
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

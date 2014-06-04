package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/events"
	httpFetcher "github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/processor"
	"github.com/wanelo/image-server/processor/cli"
	sm "github.com/wanelo/image-server/source_mapper/waneloS3"
	"github.com/wanelo/image-server/uploader/manta"
)

func main() {
	port := *flag.String("p", "7000", "Specifies the server port.")
	flag.Parse()

	serverConfiguration, err := core.ServerConfigurationFromFlags()
	if err != nil {
		log.Panicln(err)
	}

	mappings := make(map[string]string)
	mappings["p"] = "product/image"
	mapperConfiguration := &core.MapperConfiguration{mappings}

	adapters := &core.Adapters{
		Processor:    &cli.Processor{serverConfiguration},
		SourceMapper: &sm.SourceMapper{mapperConfiguration},
		Uploader:     manta.InitializeUploader(serverConfiguration),
	}
	serverConfiguration.Adapters = adapters

	httpFetcher.ImageDownloads = make(map[string][]chan error)
	processor.ImageProcessings = make(map[string][]chan processor.ImageProcessingResult)

	go func() {
		events.InitializeEventListeners(serverConfiguration)
	}()

	initializeRouter(serverConfiguration, port)
}

func initializeRouter(sc *core.ServerConfiguration, port string) {
	log.Println("starting server on http://0.0.0.0:" + port)

	m := martini.Classic()
	m.Map(sc)
	m.Get("/:namespace/:id1/:id2/:id3/:filename", genericImageHandler)
	m.Post("/:namespace/:id1/:id2/:id3", multiImageHandler)

	log.Fatal(http.ListenAndServe(":"+port, m))
}

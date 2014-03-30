package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/codegangsta/martini"
)

var (
	serverConfiguration *ServerConfiguration
)

func main() {
	environment := flag.String("e", "development", "Specifies the environment to run this server under (test/development/production).")
	flag.Parse()

	var err error
	serverConfiguration, err = loadServerConfiguration(*environment)
	if err != nil {
		log.Panicln(err)
	}

	imageProcessedChannel := make(chan *ImageConfiguration)
	go initializeManta(imageProcessedChannel)
	initializeRouter(serverConfiguration, imageProcessedChannel)
}

func initializeRouter(sc *ServerConfiguration, ipc chan *ImageConfiguration) {
	log.Println("starting in "+sc.Environment, "on http://0.0.0.0:"+sc.ServerPort)

	m := martini.Classic()
	m.Map(ipc)
	m.Get("/:model/:imageType/:id/:filename", genericImageHandler)

	log.Fatal(http.ListenAndServe(":"+sc.ServerPort, m))
}

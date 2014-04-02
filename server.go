package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
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

	go func() {
		initializeManta(serverConfiguration)
		initializeGraphite(serverConfiguration)
		initializeEventListeners(serverConfiguration)
	}()

	initializeRouter(serverConfiguration)
}

func initializeRouter(sc *ServerConfiguration) {
	log.Println("starting in "+sc.Environment, "on http://0.0.0.0:"+sc.ServerPort)

	m := martini.Classic()
	m.Map(sc)
	m.Get("/:model/:imageType/:id/:filename", genericImageHandler)

	log.Fatal(http.ListenAndServe(":"+sc.ServerPort, m))
}

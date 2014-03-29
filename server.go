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

	initializeManta()
	initializeRouter(serverConfiguration)
}

func initializeRouter(serverConfiguration *ServerConfiguration) {
	log.Println("starting in "+serverConfiguration.Environment, "on http://0.0.0.0:"+serverConfiguration.ServerPort)

	m := martini.Classic()
	m.Get("/:model/:imageType/:id/:filename", genericImageHandler)

	log.Fatal(http.ListenAndServe(":"+serverConfiguration.ServerPort, m))
}

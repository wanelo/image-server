package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/{width:[0-9]+}x{height:[0-9]+}.{format}", rectangleHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/x{width:[0-9]+}.{format}", squareHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/w{width:[0-9]+}.{format}", widthHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/full_size.{format}", fullSizeHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("starting in "+serverConfiguration.Environment, "on http://0.0.0.0:"+serverConfiguration.ServerPort)
	http.ListenAndServe(":"+serverConfiguration.ServerPort, nil)
}

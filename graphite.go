package main

import (
	"log"

	"github.com/marpaia/graphite-golang"
)

func initializeGraphite(sc *ServerConfiguration) {

	var err error
	// try to connect a graphite server
	if sc.GraphiteEnabled {
		sc.Graphite, err = graphite.NewGraphite(sc.GraphiteHost, sc.GraphitePort)
	} else {
		sc.Graphite = graphite.NewGraphiteNop(sc.GraphiteHost, sc.GraphitePort)
	}
	// if you couldn't connect to graphite, use a nop
	if err != nil {
		sc.Graphite = graphite.NewGraphiteNop(sc.GraphiteHost, sc.GraphitePort)
	}

	log.Printf("Loaded Graphite connection: %#v", sc.Graphite)
}

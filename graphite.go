package main

import (
	"log"

	"github.com/marpaia/graphite-golang"
)

func initializeGraphite(sc *ServerConfiguration) *graphite.Graphite {

	var err error
	var g *graphite.Graphite

	// try to connect a graphite server
	if sc.GraphiteEnabled {
		g, err = graphite.NewGraphite(sc.GraphiteHost, sc.GraphitePort)
	} else {
		g = graphite.NewGraphiteNop(sc.GraphiteHost, sc.GraphitePort)
	}
	// if you couldn't connect to graphite, use a nop
	if err != nil {
		g = graphite.NewGraphiteNop(sc.GraphiteHost, sc.GraphitePort)
	}

	log.Println("Loaded Graphite connection")
	return g
}

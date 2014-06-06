package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-martini/martini"
	config "github.com/wanelo/image-server/config/wanelo"
	"github.com/wanelo/image-server/core"
)

func main() {
	port := *flag.String("p", "7000", "Specifies the server port.")
	flag.Parse()

	serverConfiguration, err := config.ServerConfiguration()
	if err != nil {
		log.Panicln(err)
	}

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

package server

import (
	"log"
	"net/http"
	"time"

	"github.com/tylerb/graceful"
	"github.com/wanelo/image-server/processor/cli"
)

// ShuttingDown variable is used to note that the server is about to shut down.
// It is false by default, and set to true when a shutdown signal is received.
var ShuttingDown bool

// StatusOkMsg stores the response body for the server status server.
var StatusOkMsg []byte

func init() {
	ShuttingDown = false
	StatusOkMsg = []byte("OK")
}

// InitializeStatusServer starts a web server that can be used to monitor the health of the application.
// It returns a response with status code 200 if the system is healthy.
func InitializeStatusServer(listen string, port string) {
	log.Printf("starting startus check server on http://%s:%s", listen, port)

	srv := &graceful.Server{
		Timeout: 30 * time.Second,
		Server: &http.Server{
			Addr:    listen + ":" + port,
			Handler: &ServerStatus{},
		},
	}

	srv.ListenAndServe()
}

// ServerStatus implements the http.Handler interface
type ServerStatus struct{}

// ServeHTTP serves the http response for the status page.
// It returns a response code 200 when the image server is available to process images.

// It returns a status code 501 when the server is shutting down, or when a processor is not detected.
// Details are provided in the body of the request.
//
func (f *ServerStatus) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	processorAvailable := cli.Available

	if ShuttingDown {
		w.WriteHeader(501)
		w.Write([]byte("Shutting down"))
	} else if processorAvailable {
		w.WriteHeader(200)
		w.Write(StatusOkMsg)
	} else {
		w.WriteHeader(501)
		w.Write([]byte("There is no processor available. Make sure you have image magick installed."))
	}
}

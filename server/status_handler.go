package server

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/tylerb/graceful"
	"github.com/unrolled/render"
	"github.com/wanelo/image-server/processor/cli"
)

// Counters store atomic values
type Counters struct {
	Posting uint64 `json:"posting"`
}

// Status keeps the current state of the server
type Status struct {
	Message  string
	Counters *Counters
}

// ShuttingDown variable is used to note that the server is about to shut down.
// It is false by default, and set to true when a shutdown signal is received.
var ShuttingDown bool

func init() {
	ShuttingDown = false
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

// IncrCounter increments a given counter by 1
func IncrCounter(addr *uint64) {
	atomic.AddUint64(addr, 1)
}

// DecrCounter decrements a given counter by 1
func DecrCounter(addr *uint64) {
	atomic.AddUint64(addr, ^uint64(0))
}

var status = &Status{
	Counters: &Counters{Posting: 0},
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

	r := render.New(render.Options{
		IndentJSON: true,
	})

	var code int

	if ShuttingDown {
		status.Message = "Shutting down"
		code = 501
	} else if processorAvailable {
		status.Message = "OK"
		code = 200
	} else {
		status.Message = "There is no processor available. Make sure you have image magick installed."
		code = 501
	}

	r.JSON(w, code, status)
}

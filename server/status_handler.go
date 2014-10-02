package server

import (
	"log"
	"net/http"

	"github.com/wanelo/image-server/processor/cli"
)

// ShuttingDown is used to note that a shutdown signal has been sent
var ShuttingDown bool
var OkMsg []byte

func init() {
	ShuttingDown = false
	OkMsg = []byte("OK")
}

func InitializeServerStatus(listen string, port string) {
	log.Printf("starting startus check server on http://%s:%s", listen, port)
	http.ListenAndServe(listen+":"+port, &ServerStatus{})
}

type ServerStatus struct{}

func (f *ServerStatus) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	processorAvailable := cli.Available

	if processorAvailable && !ShuttingDown {
		w.WriteHeader(200)
		w.Write(OkMsg)
	} else {
		w.WriteHeader(501)
	}
}

// StatusHandler returns success when the server is available
// - image processor is available
// - server is not shutting down
func StatusHandler(w http.ResponseWriter, req *http.Request) {
	processorAvailable := cli.Available

	if processorAvailable && !ShuttingDown {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(501)
	}
}

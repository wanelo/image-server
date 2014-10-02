package server

import (
	"net/http"

	"github.com/wanelo/image-server/processor/cli"
)

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

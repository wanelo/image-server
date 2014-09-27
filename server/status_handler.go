package server

import (
	"net/http"

	"github.com/wanelo/image-server/processor/cli"
)

func StatusHandler(w http.ResponseWriter, req *http.Request) {
	processorAvailable := cli.Available

	if processorAvailable {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(501)
	}
}

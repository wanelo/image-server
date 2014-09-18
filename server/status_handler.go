package server

import "net/http"

func StatusHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

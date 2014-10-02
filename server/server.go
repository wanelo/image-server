package server

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/wanelo/image-server/core"
)

// InitializeServer creates a new http server to handle image processing requests
func InitializeServer(sc *core.ServerConfiguration, listen string, port string) {
	go InitializeServerStatus(listen, "7002")
	log.Printf("starting server on http://%s:%s", listen, port)
	router := NewRouter(sc)
	n := negroni.Classic()
	n.UseHandler(router)
	n.Run(listen + ":" + port)
}

// NewRouter creates a mux.Router for use in code or in tests
func NewRouter(sc *core.ServerConfiguration) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/{namespace:[a-z0-9_]+}", func(wr http.ResponseWriter, req *http.Request) {
		NewImageHandler(wr, req, sc)
	}).Methods("POST").Name("newImage")

	router.HandleFunc("/{namespace:[a-z0-9_]+}/{id1:[a-f0-9]{3}}/{id2:[a-f0-9]{3}}/{id3:[a-f0-9]{3}}/{id4:[a-f0-9]{23}}/process", func(wr http.ResponseWriter, req *http.Request) {
		ResizeManyHandler(wr, req, sc)
	}).Methods("POST").Name("resizeMany")

	router.HandleFunc("/{namespace:[a-z0-9_]+}/{id1:[a-f0-9]{3}}/{id2:[a-f0-9]{3}}/{id3:[a-f0-9]{3}}/{id4:[a-f0-9]{23}}/{filename}", func(wr http.ResponseWriter, req *http.Request) {
		ResizeHandler(wr, req, sc)
	}).Methods("GET").Name("resizeImage")

	router.HandleFunc("/{namespace:[a-z0-9_]+}/batch", func(wr http.ResponseWriter, req *http.Request) {
		CreateBatchHandler(wr, req, sc)
	}).Methods("POST").Name("createBatch")

	router.HandleFunc("/{namespace:[a-z0-9_]+}/batch/{uuid:[a-f0-9-]{36}}", func(wr http.ResponseWriter, req *http.Request) {
		BatchHandler(wr, req, sc)
	}).Methods("GET").Name("batch")

	router.HandleFunc("/status_check", StatusHandler)
	return router
}

package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
)

func genericImageHandler(params martini.Params, r *http.Request, w http.ResponseWriter, sc *ServerConfiguration) {
	ic, err := NameToConfiguration(sc, params["filename"])
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, sc, ic)
	}
	qs := r.URL.Query()
	ic.model = params["model"]
	ic.imageType = params["imageType"]
	ic.id = params["id"]
	ic.source = qs.Get("source")
	imageHandler(sc, ic, w, r)

	go func() {
		sc.Events.ImageProcessed <- ic
	}()
}

func imageHandler(sc *ServerConfiguration, ic *ImageConfiguration, w http.ResponseWriter, r *http.Request) {
	if ic.width > sc.MaximumWidth {
		err := fmt.Errorf("Maximum width is: %v\n", sc.MaximumWidth)
		errorHandler(err, w, r, http.StatusNotAcceptable, sc, ic)
		return
	}
	resizedPath, err := ic.createImage(sc)
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, sc, ic)
		return
	}
	http.ServeFile(w, r, resizedPath)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int, sc *ServerConfiguration, ic *ImageConfiguration) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}

	go func() {
		sc.Events.ImageProcessedWithErrors <- ic
	}()
}

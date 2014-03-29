package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/martini"
)

func genericImageHandler(params martini.Params, r *http.Request, w http.ResponseWriter) {
	ic, err := NameToConfiguration(params["filename"])
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound)
	}
	qs := r.URL.Query()
	ic.model = params["model"]
	ic.imageType = params["imageType"]
	ic.id = params["id"]
	ic.source = qs.Get("source")
	ic.quality = serverConfiguration.DefaultCompression
	imageHandler(ic, w, r)
}

func imageHandler(ic *ImageConfiguration, w http.ResponseWriter, r *http.Request) {
	if ic.width > serverConfiguration.MaximumWidth {
		err := fmt.Errorf("Maximum width is: %v\n", serverConfiguration.MaximumWidth)
		errorHandler(err, w, r, http.StatusNotAcceptable)
		return
	}
	resizedPath, err := createImage(ic)
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, resizedPath)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}
}

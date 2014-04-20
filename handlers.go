package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/parser"
)

func genericImageHandler(params martini.Params, r *http.Request, w http.ResponseWriter, sc *core.ServerConfiguration) {
	ic, err := parser.NameToConfiguration(sc, params["filename"])
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, ic)
	}
	qs := r.URL.Query()
	ic.ServerConfiguration = sc
	ic.Model = params["model"]
	ic.ImageType = params["imageType"]
	ic.ID = params["id"]
	ic.Source = qs.Get("source")
	imageHandler(ic, w, r)

	go func() {
		sc.Events.ImageProcessed <- ic
	}()
}

func imageHandler(ic *core.ImageConfiguration, w http.ResponseWriter, r *http.Request) {
	allowed, err := allowedImage(ic)
	if !allowed {
		errorHandler(err, w, r, http.StatusNotAcceptable, ic)
		return
	}

	sc := ic.ServerConfiguration
	resizedPath, err := sc.Adapters.Processor.CreateImage(ic)
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, ic)
		return
	}
	http.ServeFile(w, r, resizedPath)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int, ic *core.ImageConfiguration) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}

	go func() {
		sc := ic.ServerConfiguration
		sc.Events.ImageProcessedWithErrors <- ic
	}()
}

func allowedImage(ic *core.ImageConfiguration) (bool, error) {
	sc := ic.ServerConfiguration
	// verify maximum width
	if ic.Width > sc.MaximumWidth {
		err := fmt.Errorf("maximum width is: %v\n", sc.MaximumWidth)
		return false, err
	}

	// verify image format
	for _, ext := range sc.WhitelistedExtensions {
		if ext == ic.Format {
			return true, nil
		}
	}
	return false, fmt.Errorf("format not allowed %s", ic.Format)
}

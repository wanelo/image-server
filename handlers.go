package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/parser"
	"github.com/wanelo/image-server/processor/magick"
)

func genericImageHandler(params martini.Params, r *http.Request, w http.ResponseWriter, sc *core.ServerConfiguration) {
	ic, err := parser.NameToConfiguration(sc, params["filename"])
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, sc, ic)
	}
	qs := r.URL.Query()
	ic.Model = params["model"]
	ic.ImageType = params["imageType"]
	ic.ID = params["id"]
	ic.Source = qs.Get("source")
	imageHandler(sc, ic, w, r)

	go func() {
		sc.Events.ImageProcessed <- ic
	}()
}

func imageHandler(sc *core.ServerConfiguration, ic *core.ImageConfiguration, w http.ResponseWriter, r *http.Request) {
	allowed, err := allowedImage(sc, ic)
	if !allowed {
		errorHandler(err, w, r, http.StatusNotAcceptable, sc, ic)
		return
	}
	resizedPath, err := magick.CreateImage(sc, ic)
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound, sc, ic)
		return
	}
	http.ServeFile(w, r, resizedPath)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int, sc *core.ServerConfiguration, ic *core.ImageConfiguration) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}

	go func() {
		sc.Events.ImageProcessedWithErrors <- ic
	}()
}

func allowedImage(sc *core.ServerConfiguration, ic *core.ImageConfiguration) (bool, error) {
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

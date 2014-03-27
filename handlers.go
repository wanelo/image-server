package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rainycape/magick"
)

func imageHandler(ic *ImageConfiguration, w http.ResponseWriter, r *http.Request) {
	if ic.width > serverConfiguration.MaximumWidth {
		err := fmt.Errorf("Maximum width is: %v\n", serverConfiguration.MaximumWidth)
		errorHandler(err, w, r, http.StatusNotAcceptable)
		return
	}
	resizedPath, err := createImages(ic)
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, resizedPath)
}

func rectangleHandler(w http.ResponseWriter, r *http.Request) {
	ic := newImageConfiguration(r)
	imageHandler(ic, w, r)
}

func squareHandler(w http.ResponseWriter, r *http.Request) {
	ic := newImageConfiguration(r)
	ic.height = ic.width
	imageHandler(ic, w, r)
}

func widthHandler(w http.ResponseWriter, r *http.Request) {
	ic := newImageConfiguration(r)
	ic.height = 0
	imageHandler(ic, w, r)
}

func fullSizeHandler(w http.ResponseWriter, r *http.Request) {
	ic := newImageConfiguration(r)
	fullSizePath := ic.LocalOriginalImagePath()
	resizedPath := ic.LocalResizedImagePath()

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		downloadAndSaveOriginal(ic)

		im, err := magick.DecodeFile(fullSizePath)
		if err != nil {
			log.Panicln(err)
			return
		}
		defer im.Dispose()

		out, err := os.Create(resizedPath)
		defer out.Close()

		info := ic.MagickInfo()
		err = im.Encode(out, info)

		if err != nil {
			log.Panicln(err)
			return
		}
	}

	http.ServeFile(w, r, resizedPath)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}
}

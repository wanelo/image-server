package main

import (
	"fmt"
	"github.com/rainycape/magick"
	"log"
	"net/http"
	"os"
)

func imageHandler(ic *ImageConfiguration, w http.ResponseWriter, r *http.Request) {
	resizedPath, err := createImages(ic)
	if err != nil {
		errorHandler(err, w, r, http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, resizedPath)
}

func rectangleHandler(w http.ResponseWriter, r *http.Request) {
	ic := buildImageConfiguration(r)
	imageHandler(ic, w, r)
}

func squareHandler(w http.ResponseWriter, r *http.Request) {
	ic := buildImageConfiguration(r)
	ic.height = ic.width
	imageHandler(ic, w, r)
}

func widthHandler(w http.ResponseWriter, r *http.Request) {
	ic := buildImageConfiguration(r)
	ic.height = 0
	imageHandler(ic, w, r)
}

func fullSizeHandler(w http.ResponseWriter, r *http.Request) {
	ic := buildImageConfiguration(r)
	fullSizePath := ic.OriginalImagePath()
	resizedPath := ic.ResizedImagePath()

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

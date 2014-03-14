package main

import (
	"github.com/rainycape/magick"
	"log"
	"net/http"
	"os"
)

func rectangleHandler(w http.ResponseWriter, r *http.Request) {
	ic := buildImageConfiguration(r)
	resizedPath := createImages(ic)
	http.ServeFile(w, r, resizedPath)
}

func squareHandler(w http.ResponseWriter, r *http.Request) {
	ic := buildImageConfiguration(r)
	ic.height = ic.width
	resizedPath := createImages(ic)
	http.ServeFile(w, r, resizedPath)
}

func widthHandler(w http.ResponseWriter, r *http.Request) {
	ic := buildImageConfiguration(r)
	ic.height = 0
	resizedPath := createImages(ic)
	http.ServeFile(w, r, resizedPath)
}

func fullSizeHandler(w http.ResponseWriter, r *http.Request) {
	ic := buildImageConfiguration(r)
	fullSizePath := "public/" + ic.id
	resizedPath := "public/generated/" + ic.id + "_full_size." + ic.format

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

		info := magick.NewInfo()
		info.SetQuality(75)
		info.SetFormat(ic.format)
		err = im.Encode(out, info)

		if err != nil {
			log.Panicln(err)
			return
		}
	}

	http.ServeFile(w, r, resizedPath)
}
package main

import (
	"github.com/gorilla/mux"
	"github.com/rainycape/magick"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	DEFAULT_PORT string = "7000"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/{width:[0-9]+}x{height:[0-9]+}.{format}", rectangleHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/x{width:[0-9]+}.{format}", squareHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/w{width:[0-9]+}.{format}", widthHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/full_size.{format}", fullSizeHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listening on port", DEFAULT_PORT, "...")
	http.ListenAndServe(":"+DEFAULT_PORT, nil)
}

func downloadAndSaveOriginal(ic *ImageConfiguration) {
	path := ic.OriginalImagePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		start := time.Now()

		remoteUrl := ic.RemoteImageUrl()
		resp, err := http.Get(remoteUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		dir := filepath.Dir(path)
		os.MkdirAll(dir, 0700)

		out, err := os.Create(path)
		defer out.Close()

		io.Copy(out, resp.Body)
		log.Printf("Took %s to download image: %s", time.Since(start), path)
	}
}

func createWithMagick(ic *ImageConfiguration) {
	start := time.Now()
	fullSizePath := ic.OriginalImagePath()
	im, err := magick.DecodeFile(fullSizePath)
	if err != nil {
		log.Panicln(err)
		return
	}
	defer im.Dispose()

	im2, err := im.CropResize(ic.width, ic.height, magick.FHamming, magick.CSCenter)
	if err != nil {
		log.Panicln(err)
		return
	}

	resizedPath := ic.ResizedImagePath()
	dir := filepath.Dir(resizedPath)
	os.MkdirAll(dir, 0700)

	out, err := os.Create(resizedPath)
	defer out.Close()

	info := ic.MagickInfo()
	err = im2.Encode(out, info)

	if err != nil {
		log.Panicln(err)
		return
	}
	elapsed := time.Since(start)
	log.Printf("Took %s to generate image: %s", elapsed, resizedPath)
}

func createImages(ic *ImageConfiguration) (path string) {
	resizedPath := ic.ResizedImagePath()

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		downloadAndSaveOriginal(ic)
		createWithMagick(ic)
	}

	return resizedPath
}

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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/product/image/{id:[0-9]+}/{width:[0-9]+}x{height:[0-9]+}.{format}", rectangleHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/x{width:[0-9]+}.{format}", squareHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/w{width:[0-9]+}.{format}", widthHandler).Methods("GET")
	r.HandleFunc("/product/image/{id:[0-9]+}/full_size.{format}", fullSizeHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("Listening on port 7000...")
	http.ListenAndServe(":7000", nil)
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
		os.Mkdir(dir, 0700)

		out, err := os.Create(path)
		defer out.Close()

		io.Copy(out, resp.Body)
		elapsed := time.Since(start)
		log.Printf("Took %s to download image: %s", elapsed, path)
	}
}

func createWithMagick(ic *ImageConfiguration, resizedPath string) {
	fullSizePath := ic.OriginalImagePath()
	start := time.Now()
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

	out, err := os.Create(resizedPath)
	defer out.Close()

	info := magick.NewInfo()
	info.SetQuality(75)
	info.SetFormat(ic.format)
	err = im2.Encode(out, info)

	if err != nil {
		log.Panicln(err)
		return
	}
	elapsed := time.Since(start)
	log.Printf("Took %s to generate image: %s", elapsed, resizedPath)
}

func createImages(ic *ImageConfiguration) (path string) {
	var resizedPath = ic.ResizedImagePath()
	log.Printf("Source specified: %s", ic.source)

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		dir := filepath.Dir(resizedPath)
		os.Mkdir(dir, 0700)

		downloadAndSaveOriginal(ic)
		createWithMagick(ic, resizedPath)
	}

	return resizedPath
}

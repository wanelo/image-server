package main

import (
	/*	"fmt"*/
	"github.com/gorilla/mux"
	"github.com/rainycape/magick"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func createWithMagick(ic *ImageConfiguration, resizedPath string, width string, height string, format string) {
	fullSizePath := ic.OriginalImagePath()
	start := time.Now()
	im, err := magick.DecodeFile(fullSizePath)
	if err != nil {
		log.Panicln(err)
		return
	}
	defer im.Dispose()

	w, _ := strconv.Atoi(width)
	h, _ := strconv.Atoi(height)

	im2, err := im.CropResize(w, h, magick.FHamming, magick.CSCenter)
	if err != nil {
		log.Panicln(err)
		return
	}

	out, err := os.Create(resizedPath)
	defer out.Close()

	info := magick.NewInfo()
	info.SetQuality(75)
	info.SetFormat(format)
	err = im2.Encode(out, info)

	if err != nil {
		log.Panicln(err)
		return
	}
	elapsed := time.Since(start)
	log.Printf("Took %s to generate image: %s", elapsed, resizedPath)
}

func createImages(ic *ImageConfiguration) (path string) {
	var resizedPath string
	if ic.height == "0" {
		resizedPath = "public/generated/" + ic.id + "_x" + ic.width + "." + ic.format
	} else {
		resizedPath = "public/generated/" + ic.id + "_" + ic.width + "x" + ic.height + "." + ic.format
	}

	log.Printf("Source specified: %s", ic.source)
	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		dir := filepath.Dir(resizedPath)
		os.Mkdir(dir, 0700)

		downloadAndSaveOriginal(ic)
		createWithMagick(ic, resizedPath, ic.width, ic.height, ic.format)
	}

	return resizedPath
}

func buildImageConfiguration(r *http.Request) *ImageConfiguration {
	ic := new(ImageConfiguration)
	params := mux.Vars(r)
	qs := r.URL.Query()

	ic.id = params["id"]
	ic.width = params["width"]
	ic.height = params["height"]
	ic.format = params["format"]
	ic.source = qs.Get("source")

	return ic
}

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
	ic.height = "0"
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

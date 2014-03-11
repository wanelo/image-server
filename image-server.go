package main

import (
	/*	"code.google.com/p/graphics-go/graphics"*/
	/*	"fmt"*/
	"github.com/gorilla/mux"
	"github.com/rainycape/magick"
	/*	"image"*/
	/*	"image/jpeg"*/
	"io"
	"log"
	"net/http"
	"os"
	/*	"runtime"*/
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

func downloadAndSaveOriginal(path string, productId string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		start := time.Now()
		resp, err := http.Get("http://cdn-s3-2.wanelo.com/product/image/" + productId + "/original.jpg")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		out, err := os.Create(path)
		defer out.Close()

		/*	imgBody := resp.Body*/
		io.Copy(out, resp.Body)
		elapsed := time.Since(start)
		log.Printf("Took %s to download image: %s", elapsed, path)
	}
}

func createWithMagick(fullSizePath string, resizedPath string, width string, height string, format string) {
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

func createImages(id string, width string, height string, format string) (path string) {
	fullSizePath := "public/" + id
	var resizedPath string
	if height == "0" {
		resizedPath = "public/generated/" + id + "_x" + width + "." + format
	} else {
		resizedPath = "public/generated/" + id + "_" + width + "x" + height + "." + format
	}

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		downloadAndSaveOriginal(fullSizePath, id)
		createWithMagick(fullSizePath, resizedPath, width, height, format)
	}

	return resizedPath
}

func rectangleHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	height := params["height"]
	format := params["format"]

	resizedPath := createImages(id, width, height, format)
	http.ServeFile(w, r, resizedPath)
}

func squareHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	height := params["width"]
	format := params["format"]

	resizedPath := createImages(id, width, height, format)
	http.ServeFile(w, r, resizedPath)
}

func widthHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	width := params["width"]
	height := "0"
	format := params["format"]

	resizedPath := createImages(id, width, height, format)
	http.ServeFile(w, r, resizedPath)
}

func fullSizeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	format := params["format"]

	fullSizePath := "public/" + id
	resizedPath := "public/generated/" + id + "_full_size." + format

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		downloadAndSaveOriginal(fullSizePath, id)

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
		info.SetFormat(format)
		err = im.Encode(out, info)

		if err != nil {
			log.Panicln(err)
			return
		}
	}

	http.ServeFile(w, r, resizedPath)
}

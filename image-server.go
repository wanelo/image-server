package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/rainycape/magick"
)

var (
	serverConfiguration *ServerConfiguration
)

func main() {
	environment := flag.String("e", "development", "Specifies the environment to run this server under (test/development/production).")
	flag.Parse()

	var err error
	serverConfiguration, err = loadServerConfiguration(*environment)
	if err != nil {
		log.Panicln(err)
	}

	InitializeManta()
	InitializeRouter(serverConfiguration)
}

func InitializeRouter(serverConfiguration *ServerConfiguration) {
	r := mux.NewRouter()
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/{width:[0-9]+}x{height:[0-9]+}.{format}", rectangleHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/x{width:[0-9]+}.{format}", squareHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/w{width:[0-9]+}.{format}", widthHandler).Methods("GET")
	r.HandleFunc("/{model}/{imageType}/{id:[0-9]+}/full_size.{format}", fullSizeHandler).Methods("GET")
	http.Handle("/", r)
	log.Println("starting in " + serverConfiguration.Environment, "on http://0.0.0.0:" + serverConfiguration.ServerPort)
	http.ListenAndServe(":"+serverConfiguration.ServerPort, nil)
}

func downloadAndSaveOriginal(ic *ImageConfiguration) error {
	path := ic.OriginalImagePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		start := time.Now()

		remoteUrl := ic.RemoteImageUrl()
		resp, err := http.Get(remoteUrl)

		log.Printf("response code %d", resp.StatusCode)
		if err != nil || resp.StatusCode != 200 {
			log.Printf("Unable to download image: %s, status code: %d", remoteUrl, resp.StatusCode)
			log.Println(err)
			return fmt.Errorf("Unable to download image: %s, status code: %d", remoteUrl, resp.StatusCode)
		}
		defer resp.Body.Close()

		dir := filepath.Dir(path)
		os.MkdirAll(dir, 0700)

		out, err := os.Create(path)
		defer out.Close()
		if err != nil {
			log.Printf("Unable to create file: %s", path)
			log.Println(err)
			return fmt.Errorf("Unable to create file: %s", path)
		}

		io.Copy(out, resp.Body)
		log.Printf("Took %s to download image: %s", time.Since(start), path)
	}
	return nil
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

func createImages(ic *ImageConfiguration) (string, error) {
	resizedPath := ic.ResizedImagePath()

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		err := downloadAndSaveOriginal(ic)
		log.Printf("what errors? %v", err)
		if err != nil {
			log.Printf("--something happened, skipping creation")
			return "", err
		}

		createWithMagick(ic)
	}

	return resizedPath, nil
}

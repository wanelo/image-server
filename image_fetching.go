package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var imageDownloads map[string][]chan error

func fetchOriginal(ic *ImageConfiguration, sc *ServerConfiguration) error {
	c := make(chan error)
	go uniqueFetchOriginal(c, ic, sc)
	return <-c
}

// Even if simultaneous calls request the same image, only the first one will download
// the image, and will then notify all requesters. The channel returns an error object
func uniqueFetchOriginal(c chan error, ic *ImageConfiguration, sc *ServerConfiguration) {
	key := ic.RemoteImageUrl()
	_, present := imageDownloads[key]

	if present {
		imageDownloads[key] = append(imageDownloads[key], c)
	} else {
		imageDownloads[key] = []chan error{c}

		err := downloadAndSaveOriginal(ic, sc)
		for _, cc := range imageDownloads[key] {
			cc <- err
		}
		delete(imageDownloads, key)
	}
}

func downloadAndSaveOriginal(ic *ImageConfiguration, sc *ServerConfiguration) error {
	path := ic.LocalOriginalImagePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		start := time.Now()

		remoteUrl := ic.RemoteImageUrl()
		resp, err := http.Get(remoteUrl)

		log.Printf("response code %d", resp.StatusCode)
		if err != nil || resp.StatusCode != 200 {
			log.Printf("Unable to download image: %s, status code: %d", remoteUrl, resp.StatusCode)
			log.Println(err)
			go func() {
				sc.Events.OriginalDownloadUnavailable <- ic
			}()
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

		go func() {
			sc.Events.OriginalDownloaded <- ic
		}()
	}
	return nil
}

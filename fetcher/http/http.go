package http

import (
	"fmt"
	"io"
	"log"
	gohttp "net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wanelo/image-server/core"
)

var ImageDownloads map[string][]chan error

func FetchOriginal(sc *core.ServerConfiguration, ic *core.ImageConfiguration) error {
	c := make(chan error)
	go uniqueFetchOriginal(c, sc, ic)
	return <-c
}

// Even if simultaneous calls request the same image, only the first one will download
// the image, and will then notify all requesters. The channel returns an error object
func uniqueFetchOriginal(c chan error, sc *core.ServerConfiguration, ic *core.ImageConfiguration) {
	key := sc.RemoteImageURL(ic)
	_, present := ImageDownloads[key]

	if present {
		ImageDownloads[key] = append(ImageDownloads[key], c)
	} else {
		ImageDownloads[key] = []chan error{c}

		err := downloadAndSaveOriginal(ic, sc)
		for _, cc := range ImageDownloads[key] {
			cc <- err
		}
		delete(ImageDownloads, key)
	}
}

func downloadAndSaveOriginal(ic *core.ImageConfiguration, sc *core.ServerConfiguration) error {
	path := ic.LocalOriginalImagePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		start := time.Now()

		remoteURL := sc.RemoteImageURL(ic)
		resp, err := gohttp.Get(remoteURL)

		log.Printf("response code %d", resp.StatusCode)
		if err != nil || resp.StatusCode != 200 {
			log.Printf("Unable to download image: %s, status code: %d", remoteURL, resp.StatusCode)
			log.Println(err)
			go func() {
				sc.Events.OriginalDownloadUnavailable <- ic
			}()
			return fmt.Errorf("unable to download image: %s, status code: %d", remoteURL, resp.StatusCode)
		}
		defer resp.Body.Close()

		dir := filepath.Dir(path)
		os.MkdirAll(dir, 0700)

		out, err := os.Create(path)
		defer out.Close()
		if err != nil {
			log.Printf("Unable to create file: %s", path)
			log.Println(err)
			return fmt.Errorf("unable to create file: %s", path)
		}

		io.Copy(out, resp.Body)
		log.Printf("Took %s to download image: %s", time.Since(start), path)

		go func() {
			sc.Events.OriginalDownloaded <- ic
		}()
	}
	return nil
}

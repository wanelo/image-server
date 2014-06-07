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

func init() {
	ImageDownloads = make(map[string][]chan error)
}

type Fetcher struct {
	SourceMapper core.SourceMapper
}

func (f *Fetcher) FetchOriginal(ic *core.ImageConfiguration) error {
	c := make(chan error)
	go f.uniqueFetchOriginal(c, ic)
	return <-c
}

func (f *Fetcher) remoteImageURL(ic *core.ImageConfiguration) string {
	return f.SourceMapper.RemoteImageURL(ic)
}

// Even if simultaneous calls request the same image, only the first one will download
// the image, and will then notify all requesters. The channel returns an error object
func (f *Fetcher) uniqueFetchOriginal(c chan error, ic *core.ImageConfiguration) {
	key := f.remoteImageURL(ic)

	_, present := ImageDownloads[key]

	if present {
		ImageDownloads[key] = append(ImageDownloads[key], c)
	} else {
		ImageDownloads[key] = []chan error{c}

		err := f.downloadAndSaveOriginal(ic)
		for _, cc := range ImageDownloads[key] {
			cc <- err
		}
		delete(ImageDownloads, key)
	}
}

func (f *Fetcher) downloadAndSaveOriginal(ic *core.ImageConfiguration) error {
	path := ic.LocalOriginalImagePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		start := time.Now()

		remoteURL := f.remoteImageURL(ic)
		resp, err := gohttp.Get(remoteURL)

		if err != nil || resp.StatusCode != 200 {
			log.Printf("Unable to download image: %s, status code: %d", remoteURL, resp.StatusCode)
			log.Println(err)
			go func() {
				ic.ServerConfiguration.Events.OriginalDownloadUnavailable <- ic
			}()
			return fmt.Errorf("Unable to download image: %s, status code: %d", remoteURL, resp.StatusCode)
		}
		log.Printf("Downloaded from %s with code %d", remoteURL, resp.StatusCode)
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
			ic.ServerConfiguration.Events.OriginalDownloaded <- ic
		}()
	}
	return nil
}

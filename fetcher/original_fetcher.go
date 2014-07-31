package fetcher

import (
	"log"
	"os"
	"path/filepath"

	httpFetcher "github.com/wanelo/image-server/fetcher/http"
)

type UniqueFetcher struct {
	Source      string
	Destination string
	Channels    *FetcherChannels
}

func NewUniqueFetcher(source string, destination string) *UniqueFetcher {
	channels := &FetcherChannels{
		make(chan string),
		make(chan string),
		make(chan string),
	}

	return &UniqueFetcher{source, destination, channels}
}

func (f *UniqueFetcher) Fetch() error {
	c := make(chan FetchResult)
	go f.uniqueFetch(c)
	r := <-c
	return r.Error
}

// Even if simultaneous calls request the same image, only the first one will download
// the image, and will then notify all requesters. The channel returns an error object
func (f *UniqueFetcher) uniqueFetch(c chan FetchResult) {
	url := f.Source
	destination := f.Destination
	_, present := ImageDownloads[url]
	var err error

	if present {
		log.Println("Already downloading")
		ImageDownloads[url] = append(ImageDownloads[url], c)
	} else {
		ImageDownloads[url] = []chan FetchResult{c}
		defer delete(ImageDownloads, url)

		// only copy image if does not exist
		if _, err := os.Stat(destination); os.IsNotExist(err) {
			dir := filepath.Dir(destination)
			os.MkdirAll(dir, 0700)

			fetcher := &httpFetcher.Fetcher{}
			fetcher.Fetch(url, destination)

			if err == nil {
				go func() {
					f.Channels.DownloadComplete <- destination
				}()
			}
		} else {
			go func() {
				f.Channels.SkippedDownload <- destination
			}()
		}

		if err == nil {
			log.Printf("Notifying download complete for path %s", destination)
			f.notifyDownloadComplete(url)
		} else {
			go func() {
				f.Channels.DownloadFailed <- destination
			}()
			f.notifyDownloadFailed(url, err)
		}

	}
}

func (f *UniqueFetcher) notifyDownloadComplete(url string) {
	for _, cc := range ImageDownloads[url] {
		fr := FetchResult{nil, nil}
		cc <- fr
	}
}

func (f *UniqueFetcher) notifyDownloadFailed(url string, err error) {
	for _, cc := range ImageDownloads[url] {
		fr := FetchResult{err, nil}
		cc <- fr
	}
}

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
}

func NewUniqueFetcher(source string, destination string) *UniqueFetcher {
	return &UniqueFetcher{source, destination}
}

func (f *UniqueFetcher) Fetch() (bool, error) {
	c := make(chan FetchResult)
	defer close(c)
	go f.uniqueFetch(c)
	r := <-c
	return r.Downloaded, r.Error
}

// Even if simultaneous calls request the same image, only the first one will download
// the image, and will then notify all requesters. The channel returns an error object
func (f *UniqueFetcher) uniqueFetch(c chan FetchResult) {
	url := f.Source
	destination := f.Destination
	_, present := ImageDownloads[url]
	var err error

	if present {
		ImageDownloads[url] = append(ImageDownloads[url], c)
	} else {
		ImageDownloads[url] = []chan FetchResult{c}
		defer delete(ImageDownloads, url)

		// only copy image if does not exist
		if _, err = os.Stat(destination); os.IsNotExist(err) {
			dir := filepath.Dir(destination)
			os.MkdirAll(dir, 0700)

			fetcher := &httpFetcher.Fetcher{}
			err = fetcher.Fetch(url, destination)
		}

		if err == nil {
			log.Printf("Notifying download complete for path %s", destination)
			f.notifyDownloadComplete(url)
		} else {
			f.notifyDownloadFailed(url, err)
		}

	}
}

func (f *UniqueFetcher) notifyDownloadComplete(url string) {
	for i, cc := range ImageDownloads[url] {
		downloaded := i == 0
		fr := FetchResult{nil, nil, downloaded}
		cc <- fr
	}
}

func (f *UniqueFetcher) notifyDownloadFailed(url string, err error) {
	for _, cc := range ImageDownloads[url] {
		fr := FetchResult{err, nil, false}
		cc <- fr
	}
}

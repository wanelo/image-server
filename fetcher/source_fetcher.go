package fetcher

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/info"
)

type SourceFetcher struct {
	Fetcher  core.Fetcher
	Paths    core.Paths
	Channels *FetcherChannels
}

func NewSourceFetcher(paths core.Paths, fetcher core.Fetcher) *SourceFetcher {
	channels := &FetcherChannels{
		make(chan string),
		make(chan string),
		make(chan string),
	}

	return &SourceFetcher{fetcher, paths, channels}
}

func (f *SourceFetcher) Fetch(url string, namespace string) (error, *info.ImageDetails) {
	c := make(chan FetchResult)
	go f.uniqueFetchSource(c, url, namespace)
	r := <-c
	return r.Error, r.ImageDetails
}

// Even if simultaneous calls request the same image, only the first one will download
// the image, and will then notify all requesters. The channel returns an error object
func (f *SourceFetcher) uniqueFetchSource(c chan FetchResult, url string, namespace string) {
	_, present := ImageDownloads[url]
	var i info.Info
	var imageDetails *info.ImageDetails
	var destination string

	if present {
		log.Println("Already downloading")
		ImageDownloads[url] = append(ImageDownloads[url], c)
	} else {
		ImageDownloads[url] = []chan FetchResult{c}
		defer delete(ImageDownloads, url)

		tmp := f.Paths.TempImagePath(url)
		err := f.Fetcher.Fetch(url, tmp)
		i = info.Info{tmp}

		if err == nil {
			md5 := i.FileHash()
			destination = f.Paths.LocalOriginalPath(namespace, md5)

			// only copy image if does not exist
			if _, err = os.Stat(destination); os.IsNotExist(err) {
				dir := filepath.Dir(destination)
				os.MkdirAll(dir, 0700)

				cpCmd := exec.Command("cp", "-rf", tmp, destination)
				err = cpCmd.Run()

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

			i = info.Info{destination}
			imageDetails, err = i.ImageDetails()
		}

		if err == nil {
			log.Printf("Notifying download complete for path %s", destination)
			f.notifyDownloadSourceComplete(url, imageDetails)
		} else {
			go func() {
				f.Channels.DownloadFailed <- destination
			}()
			f.notifyDownloadSourceFailed(url, err)
		}

	}
}

func (f *SourceFetcher) notifyDownloadSourceComplete(url string, d *info.ImageDetails) {
	for _, cc := range ImageDownloads[url] {
		fr := FetchResult{nil, d}
		cc <- fr
	}
}

func (f *SourceFetcher) notifyDownloadSourceFailed(url string, err error) {
	for _, cc := range ImageDownloads[url] {
		fr := FetchResult{err, nil}
		cc <- fr
	}
}

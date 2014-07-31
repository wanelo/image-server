package fetcher

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/info"
)

func NewOriginalFetcher(paths core.Paths, fetcher core.Fetcher) *OriginalFetcher {
	channels := &FetcherChannels{
		make(chan string),
		make(chan string),
		make(chan string),
	}

	return &OriginalFetcher{fetcher, paths, channels}
}

func (f *OriginalFetcher) Fetch(url string, namespace string) (error, *info.ImageDetails) {
	c := make(chan FetchResult)
	go f.uniqueFetchOriginal(c, url, namespace)
	r := <-c
	return r.Error, r.ImageDetails
}

// Even if simultaneous calls request the same image, only the first one will download
// the image, and will then notify all requesters. The channel returns an error object
func (f *OriginalFetcher) uniqueFetchOriginal(c chan FetchResult, url string, namespace string) {
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
			f.notifyDownloadComplete(url, imageDetails)
		} else {
			go func() {
				f.Channels.DownloadFailed <- destination
			}()
			f.notifyDownloadFailed(url, err)
		}

	}
}

func (f *OriginalFetcher) notifyDownloadComplete(url string, d *info.ImageDetails) {

	for _, cc := range ImageDownloads[url] {
		fr := FetchResult{nil, d}
		cc <- fr

		// if i > 0 {
		// 	go func() { f.Channels.SkippedDownload <- url }()
		// }
	}
}

func (f *OriginalFetcher) notifyDownloadFailed(url string, err error) {
	// if err != nil {
	// 	go func() { f.Channels.DownloadFailed <- url }()
	// }

	for _, cc := range ImageDownloads[url] {
		fr := FetchResult{err, nil}
		cc <- fr

		// if i > 0 {
		// 	go func() { f.Channels.SkippedDownload <- url }()
		// }
	}
}

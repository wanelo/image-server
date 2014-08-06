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
	Paths core.Paths
}

func NewSourceFetcher(paths core.Paths) *SourceFetcher {
	return &SourceFetcher{paths}
}

func (f *SourceFetcher) Fetch(url string, namespace string) (*info.ImageDetails, bool, error) {
	c := make(chan FetchResult)
	defer close(c)
	go f.uniqueFetchSource(c, url, namespace)
	r := <-c
	return r.ImageDetails, r.Downloaded, r.Error
}

// Even if simultaneous calls request the same image, only the first one will download
// the image, and will then notify all requesters. The channel returns an error object
func (f *SourceFetcher) uniqueFetchSource(c chan FetchResult, url string, namespace string) {
	tmpOriginalPath, downloaded, err := f.downloadTempSource(url)

	if err != nil {
		f.notifyDownloadSourceFailed(c, err)
		return
	}

	tmpInfo := info.Info{tmpOriginalPath}

	md5, err := tmpInfo.FileHash()
	if err != nil {
		f.notifyDownloadSourceFailed(c, err)
		return
	}

	destination := f.Paths.LocalOriginalPath(namespace, md5)

	if downloaded {
		// only copy image if does not exist
		if _, err = os.Stat(destination); os.IsNotExist(err) {
			dir := filepath.Dir(destination)
			os.MkdirAll(dir, 0700)

			cpCmd := exec.Command("cp", "-rf", tmpOriginalPath, destination)
			err = cpCmd.Run()

			if err != nil {
				f.notifyDownloadSourceFailed(c, err)
				return
			}
		}
	}

	i := info.Info{destination}
	imageDetails, err := i.ImageDetails()

	if err != nil {
		f.notifyDownloadSourceFailed(c, err)
		return
	}

	c <- FetchResult{nil, imageDetails, downloaded}
}

func (f *SourceFetcher) downloadTempSource(url string) (string, bool, error) {
	log.Println("downloadTempSource", url)
	tmpOriginalPath := f.Paths.TempImagePath(url)
	log.Println("tmpOriginalPath", tmpOriginalPath)
	fetcher := NewUniqueFetcher(url, tmpOriginalPath)
	downloaded, err := fetcher.Fetch()
	return tmpOriginalPath, downloaded, err
}

func (f *SourceFetcher) notifyDownloadSourceFailed(c chan FetchResult, err error) {
	c <- FetchResult{err, nil, false}
}

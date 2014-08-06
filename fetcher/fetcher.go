package fetcher

import "github.com/wanelo/image-server/info"

var ImageDownloads map[string][]chan FetchResult

func init() {
	ImageDownloads = make(map[string][]chan FetchResult)
}

type FetcherChannels struct {
	DownloadComplete chan string
	SkippedDownload  chan string
	DownloadFailed   chan string
}

type FetchResult struct {
	Error        error
	ImageDetails *info.ImageDetails
	Downloaded   bool
}

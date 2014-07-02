package fetcher

import (
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/info"
)

var ImageDownloads map[string][]chan FetchResult

func init() {
	ImageDownloads = make(map[string][]chan FetchResult)
}

// Fetcher is initialized with the fetcher adapter
type OriginalFetcher struct {
	Fetcher  core.Fetcher
	Paths    core.Paths
	Channels *FetcherChannels
}

type FetcherChannels struct {
	DownloadComplete chan string
	SkippedDownload  chan string
	DownloadFailed   chan string
}

type FetchResult struct {
	Error        error
	ImageDetails *info.ImageDetails
}

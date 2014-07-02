package core

type Adapters struct {
	Fetcher   Fetcher
	Processor Processor
	Uploader  Uploader
	Paths     Paths
	Logger    Logger
}

type Fetcher interface {
	Fetch(string, string) error
}

type Logger interface {
	ImageProcessed(ic *ImageConfiguration)
	ImageProcessedWithErrors(ic *ImageConfiguration)
	OriginalDownloaded(source string, destination string)
	OriginalDownloadFailed(source string)
	OriginalDownloadSkipped(source string)
}

// Processor
type Processor interface {
	CreateImage(string, string, *ImageConfiguration) error
}

// Paths

type Paths interface {
	OriginalPath(string, string) string
	ImageDirectory(string, string) string
	TempImagePath(string) string
	LocalOriginalPath(string, string) string
	RemoteOriginalPath(string, string) string
}

// SourceMapper

type SourceMapper interface {
	RemoteImageURL(*ImageConfiguration) string
}

// Uploader
type Uploader interface {
	Upload(string, string) error
}

package core

// ServerConfiguration struct
// Most of this configuration comes from json config
type ServerConfiguration struct {
	WhitelistedExtensions []string
	MaximumWidth          int
	LocalBasePath         string
	RemoteBasePath        string
	RemoteBaseURL         string
	DefaultQuality        uint
	UploaderConcurrency   uint
	GraphiteHost          string
	GraphitePort          int
	Adapters              *Adapters
}

// EventChannels struct
// Available image processing/downloading events
// type EventChannels struct {
// ImageProcessed              chan *ImageConfiguration
// ImageProcessedWithErrors    chan *ImageConfiguration
// OriginalDownloaded          chan *ImageConfiguration
// OriginalDownloadFailed chan *ImageConfiguration
// }

package core

// ServerConfiguration struct
// Most of this configuration comes from json config
type ServerConfiguration struct {
	SourceDomain          string
	WhitelistedExtensions []string
	MaximumWidth          int
	LocalBasePath         string
	MantaBasePath         string
	DefaultQuality        uint
	UploaderConcurrency   uint
	GraphiteHost          string
	GraphitePort          int
	Events                *EventChannels
	Adapters              *Adapters
}

// EventChannels struct
// Available image processing/downloading events
type EventChannels struct {
	ImageProcessed              chan *ImageConfiguration
	ImageProcessedWithErrors    chan *ImageConfiguration
	OriginalDownloaded          chan *ImageConfiguration
	OriginalDownloadUnavailable chan *ImageConfiguration
}

package core

// ServerConfiguration struct
type ServerConfiguration struct {
	WhitelistedExtensions []string
	MaximumWidth          int
	LocalBasePath         string
	RemoteBasePath        string
	RemoteBaseURL         string
	DefaultQuality        uint
	UploaderConcurrency   uint
	ProcessorConcurrency  uint
	GraphiteHost          string
	GraphitePort          int
	Adapters              *Adapters
}

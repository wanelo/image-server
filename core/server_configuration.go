package core

import "time"

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
	HTTPTimeout           time.Duration
	GraphiteHost          string
	GraphitePort          int
	Adapters              *Adapters
	Outputs               string
	AWSAccessKeyID        string
	AWSSecretKey          string
	AWSBucket             string
}

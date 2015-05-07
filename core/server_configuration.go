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
	StatsdHost            string
	StatsdPort            int
	Adapters              *Adapters
	Outputs               string
	AWSAccessKeyID        string
	AWSSecretKey          string
	AWSBucket             string
	MantaURL              string
	MantaUser             string
	MantaKeyID            string
	SDCIdentity           string
}

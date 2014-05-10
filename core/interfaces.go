package core

type Adapters struct {
	Processor    Processor
	SourceMapper SourceMapper
}

type Processor interface {
	CreateImage(*ImageConfiguration) (string, error)
}

type SourceMapper interface {
	RemoteImageURL(*ImageConfiguration) string
}

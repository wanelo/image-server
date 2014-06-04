package core

type Adapters struct {
	Processor    Processor
	SourceMapper SourceMapper
	Uploader     Uploader
}

type Processor interface {
	CreateImage(*ImageConfiguration) (string, error)
}

type SourceMapper interface {
	RemoteImageURL(*ImageConfiguration) string
}

type Uploader interface {
	Upload(ic *ImageConfiguration)
	UploadOriginal(ic *ImageConfiguration)
}

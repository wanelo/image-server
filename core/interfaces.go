package core

type Adapters struct {
	Fetcher   Fetcher
	Processor Processor
	Uploader  Uploader
}

type SourceMapper interface {
	RemoteImageURL(*ImageConfiguration) string
}

type Fetcher interface {
	FetchOriginal(*ImageConfiguration) error
}

type Processor interface {
	CreateImage(*ImageConfiguration) (string, error)
}

type Uploader interface {
	Upload(*ImageConfiguration)
	UploadOriginal(*ImageConfiguration)
}

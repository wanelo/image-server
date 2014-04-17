package core

type Adapters struct {
	Processor Processor
}

type Processor interface {
	CreateImage(*ImageConfiguration) (string, error)
}

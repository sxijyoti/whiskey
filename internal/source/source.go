package source

type Metadata struct {
	ETag         string
	LastModified string
}

type Source interface {
	ID() string
	Fetch() ([]byte, error)
	Metadata() (*Metadata, error)
}

type Factory interface {
	Supports(
		ref string,
	) bool

	New(
		ref string,
	) Source
}
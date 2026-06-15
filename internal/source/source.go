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
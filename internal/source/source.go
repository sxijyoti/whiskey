package source

var Offline bool

type Metadata struct {
	ETag         string
	LastModified string
	NotModified  bool
}

type Source interface {
	ID() string
	Fetch() ([]byte, error)
	Metadata() (*Metadata, error)
}

type ConditionalSource interface {
	Source

	ConditionalMetadata(
		old *Metadata,
	) (*Metadata, error)
}

type Materialized struct {
	Workspace   string
	Content     []byte
	ContentHash string
	Metadata    Metadata
}

type Factory interface {
	Supports(
		ref string,
	) bool

	New(
		ref string,
	) Source
}

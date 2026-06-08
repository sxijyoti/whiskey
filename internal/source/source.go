package source

type Source interface {
	ID() string
	Fetch() ([]byte, error)
}


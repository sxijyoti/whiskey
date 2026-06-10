package source

import "os"

type Local struct {
	Path string
}

func (l Local) ID() string {
	return "local:" + l.Path
}

func (l Local) Fetch() ([]byte, error) {
	return os.ReadFile(l.Path)
}

func (l Local) Metadata() (*Metadata, error) {

	info, err := os.Stat(
		l.Path,
	)

	if err != nil {
		return nil, err
	}

	return &Metadata{
		LastModified: info.ModTime().
			UTC().
			String(),
	}, nil
}
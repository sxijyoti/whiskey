package source

import (
	"os"
	"strings"
)

type Local struct {
	Path string
}

type LocalFactory struct{}

func init() {
	Register(
		LocalFactory{},
	)
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

func (LocalFactory) Supports(
	ref string,
) bool {

	return strings.HasPrefix(
		ref,
		"local:",
	)
}

func (LocalFactory) New(
	ref string,
) Source {

	return Local{
		Path: strings.TrimPrefix(
			ref,
			"local:",
		),
	}
}

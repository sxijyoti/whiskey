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
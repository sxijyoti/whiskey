package source

import (
	"io"
	"net/http"
)

type GitHub struct {
	URL string
}

func (g GitHub) ID() string {
	return g.URL
}

func (g GitHub) Fetch() ([]byte, error) {

	resp, err := http.Get(g.URL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
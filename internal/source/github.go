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

	req, err := http.NewRequest(
		http.MethodGet,
		g.URL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Cache-Control",
		"no-cache",
	)

	req.Header.Set(
		"Pragma",
		"no-cache",
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
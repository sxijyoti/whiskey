package source

import (
	"io"
	"net/http"
	"strings"
)

type HTTP struct {
	URL string
}

func (h HTTP) ID() string {
	return h.URL
}

func (h HTTP) Fetch() ([]byte, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		h.URL,
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

func (h HTTP) Metadata() (*Metadata, error) {

	req, err := http.NewRequest(
		http.MethodHead,
		h.URL,
		nil,
	)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return &Metadata{
		ETag: resp.Header.Get(
			"ETag",
		),

		LastModified: resp.Header.Get(
			"Last-Modified",
		),
	}, nil
}

type HTTPFactory struct{}

func init() {
	Register(
		HTTPFactory{},
	)
}

func (HTTPFactory) Supports(
	ref string,
) bool {

	return strings.HasPrefix(
		ref,
		"http://",
	) ||
		strings.HasPrefix(
			ref,
			"https://",
	)
}

func (HTTPFactory) New(
	ref string,
) Source {

	return HTTP{
		URL: ref,
	}
}
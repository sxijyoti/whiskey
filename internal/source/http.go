package source

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sxijyoti/whiskey/internal/transform"
)

type HTTP struct {
	URL string
}

func (h HTTP) ID() string {
	return h.URL
}

var client = &http.Client{}

func (h HTTP) Fetch() ([]byte, error) {

	if err := validateURL(
		h.URL,
	); err != nil {
		return nil, err
	}

	if shouldExtractMarkdown(h.URL) {
		return transform.FetchMarkdown(h.URL)
	}

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

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {

		return nil, fmt.Errorf(
			"%s: %s",
			h.URL,
			resp.Status,
		)
	}

	return io.ReadAll(resp.Body)
}

func (h HTTP) Metadata() (*Metadata, error) {

	if err := validateURL(
		h.URL,
	); err != nil {
		return nil, err
	}

	if shouldExtractMarkdown(
		h.URL,
	) {
		return &Metadata{}, nil
	}

	req, err := http.NewRequest(
		http.MethodHead,
		h.URL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {

		return nil, fmt.Errorf(
			"%s: %s",
			h.URL,
			resp.Status,
		)
	}

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

func validateURL(
	url string,
) error {

	if strings.HasPrefix(
		url,
		"https://github.com/",
	) &&
		strings.Contains(
			url,
			"/blob/",
		) {

		raw := strings.Replace(
			url,
			"https://github.com/",
			"https://raw.githubusercontent.com/",
			1,
		)

		raw = strings.Replace(
			raw,
			"/blob/",
			"/",
			1,
		)

		return fmt.Errorf(
			"github blob URLs are not directly fetchable.\nDid you mean:\n\n%s",
			raw,
		)
	}

	return nil
}

func (h HTTP) ConditionalMetadata(
	old *Metadata,
) (*Metadata, error) {

	if err := validateURL(
		h.URL,
	); err != nil {
		return nil, err
	}

	if shouldExtractMarkdown(
		h.URL,
	) {
		return &Metadata{}, nil
	}

	req, err := http.NewRequest(
		http.MethodHead,
		h.URL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	if old != nil {

		if old.ETag != "" {
			req.Header.Set(
				"If-None-Match",
				old.ETag,
			)
		}

		if old.LastModified != "" {
			req.Header.Set(
				"If-Modified-Since",
				old.LastModified,
			)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {

		return &Metadata{
			ETag:         old.ETag,
			LastModified: old.LastModified,
			NotModified:  true,
		}, nil
	}

	if resp.StatusCode < http.StatusOK ||
		resp.StatusCode >= http.StatusMultipleChoices {

		return nil, fmt.Errorf(
			"%s: %s",
			h.URL,
			resp.Status,
		)
	}

	return &Metadata{
		ETag: resp.Header.Get(
			"ETag",
		),
		LastModified: resp.Header.Get(
			"Last-Modified",
		),
	}, nil
}

func shouldExtractMarkdown(
	url string,
) bool {

	switch {

	case strings.Contains(
		url,
		"raw.githubusercontent.com",
	):
		return false

	case strings.Contains(
		url,
		"githubusercontent.com",
	):
		return false

	case strings.HasSuffix(
		url,
		".md",
	):
		return false

	case strings.HasSuffix(
		url,
		".markdown",
	):
		return false

	case strings.HasSuffix(
		url,
		".txt",
	):
		return false
	}

	return true
}
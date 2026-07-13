package source

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	normalized, err := NormalizeRef(ref)
	if err != nil {
		normalized = ref
	}

	return HTTP{
		URL: normalized,
	}
}

func NormalizeRef(ref string) (string, error) {
	if !strings.HasPrefix(ref, "http://") &&
		!strings.HasPrefix(ref, "https://") {
		return ref, nil
	}

	return normalizeGitHubBlobURL(ref)
}

func normalizeGitHubBlobURL(raw string) (string, error) {
	if strings.HasPrefix(
		raw,
		"https://github.com/",
	) &&
		strings.Contains(
			raw,
			"/blob/",
		) {

		parsed, err := url.Parse(raw)
		if err != nil {
			return "", fmt.Errorf("invalid GitHub blob URL %q: %w", raw, err)
		}

		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) < 5 || parts[2] != "blob" {
			return "", fmt.Errorf("unsupported GitHub blob URL: %s", raw)
		}

		normalized := "https://raw.githubusercontent.com/" +
			parts[0] + "/" +
			parts[1] + "/" +
			parts[3] + "/" +
			strings.Join(parts[4:], "/")

		if parsed.RawQuery != "" {
			normalized += "?" + parsed.RawQuery
		}

		return normalized, nil
	}

	return raw, nil
}

func validateURL(
	raw string,
) error {
	_, err := NormalizeRef(raw)
	return err
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

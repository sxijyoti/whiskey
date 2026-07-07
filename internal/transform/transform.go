package transform

import "strings"

type HTML struct{}

func init() {
	Register(
		HTML{},
	)
}

func (HTML) Supports(
	contentType string,
) bool {

	return strings.HasPrefix(
		contentType,
		"text/html",
	)
}

func (HTML) Transform(
	url string,
	_ []byte,
) ([]byte, error) {

	return FetchMarkdown(
		url,
	)
}

type Transformer interface {
	Supports(
		contentType string,
	) bool

	Transform(
		url string,
		body []byte,
	) ([]byte, error)
}

var transformers []Transformer

func Register(
	t Transformer,
) {
	transformers = append(
		transformers,
		t,
	)
}

func Process(
	url string,
	contentType string,
	body []byte,
) ([]byte, error) {

	for _, t := range transformers {

		if !t.Supports(
			contentType,
		) {
			continue
		}

		return t.Transform(
			url,
			body,
		)
	}

	return body, nil
}

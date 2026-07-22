package parser

import (
	"bytes"

	"github.com/yuin/goldmark"
	goldhtml "github.com/yuin/goldmark/renderer/html"
)

var markdown = goldmark.New(
	goldmark.WithRendererOptions(
		goldhtml.WithUnsafe(),
	),
)

func MdToHTML(
	md string,
) (string, error) {

	var buf bytes.Buffer

	if err := markdown.Convert(
		[]byte(md),
		&buf,
	); err != nil {
		return "", err
	}

	return buf.String(), nil
}

package template

import (
	htmltemplate "html/template"
	"bytes"
)

type PageData struct {
	Title       string
	Description string
	Date        string
	Content     htmltemplate.HTML
}

func RenderPage(data PageData) ([]byte, error) {
	var buf bytes.Buffer

	err := Page.Execute(&buf, data)

	return buf.Bytes(), err
}


package template

import (
	"bytes"
	htmltemplate "html/template"
)

type PageData struct {
	Title       string
	Description string
	Date        string
	Content     htmltemplate.HTML
}

func RenderPage(
	siteRoot string,
	layout string,
	data PageData,
) ([]byte, error) {

	tmpl, err := LoadLayout(
		siteRoot,
		layout,
	)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, data)

	return buf.Bytes(), err
}

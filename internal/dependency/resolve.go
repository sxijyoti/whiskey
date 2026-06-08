package dependency

import (
	"strings"

	"github.com/sxijyoti/whiskey/internal/source"
)

func ResolveIncludes(
	body string,
) (string, error) {

	directives := Extract(body)

	for _, d := range directives {

		src, err := source.Resolve(
			d.Ref,
		)
		if err != nil {
			return "", err
		}

		content, err := src.Fetch()
		if err != nil {
			return "", err
		}

		body = strings.ReplaceAll(
			body,
			"@include "+d.Ref,
			string(content),
		)
	}

	return body, nil
}
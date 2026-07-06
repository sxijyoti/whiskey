package dependency

import (
	"strings"
)

func ResolveIncludes(
	body string,
	resolve func(string) ([]byte, error),
) (string, error) {

	directives := Extract(body)

	for _, d := range directives {

		content, err := resolve(d.Ref)
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

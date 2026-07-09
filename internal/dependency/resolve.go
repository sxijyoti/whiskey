package dependency

import (
	"strings"
	"regexp"
)

var includeRE = regexp.MustCompile(
	`@include\s+([^\n]+)`,
)

func Extract(
	body string,
) []Directive {

	matches := includeRE.FindAllStringSubmatch(
		body,
		-1,
	)

	var directives []Directive

	for _, m := range matches {

		directives = append(
			directives,
			Directive{
				Ref: m[1],
			},
		)
	}

	return directives
}

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

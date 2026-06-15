package dependency

import "regexp"

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
package dependency

import (
	"fmt"
	"regexp"
	"strings"
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
	contextFile string,
	resolve func(string) ([]byte, string, error),
) (string, error) {
	return resolveIncludesRecursive(
		body,
		contextFile,
		contextFile,
		resolve,
		make(map[string]bool),
	)
}

func resolveIncludesRecursive(
	body string,
	contextFile string,
	currentIdentity string,
	resolve func(string) ([]byte, string, error),
	seen map[string]bool,
) (string, error) {

	directives := Extract(body)

	for _, d := range directives {
		identity := d.Ref

		if strings.HasPrefix(d.Ref, "local:") {
			resolved, err := ResolveLocalRef(contextFile, d.Ref)
			if err != nil {
				return "", err
			}
			identity = resolved
		}

		if seen[identity] {
			if identity == currentIdentity {
				continue
			}
			return "", fmt.Errorf("cycle detected: %s", d.Ref)
		}

		seen[identity] = true
		contentBytes, nextContext, err := resolve(d.Ref)
		if err != nil {
			return "", err
		}

		if nextContext == "" {
			nextContext = contextFile
		}

		contentStr, err := resolveIncludesRecursive(
			string(contentBytes),
			nextContext,
			identity,
			resolve,
			seen,
		)
		if err != nil {
			return "", err
		}
		delete(seen, identity)

		body = strings.ReplaceAll(
			body,
			"@include "+d.Ref,
			contentStr,
		)
	}

	return body, nil
}

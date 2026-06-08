package source

import (
	"fmt"
	"strings"
)

func Resolve(
	ref string,
) (Source, error) {

	if strings.HasPrefix(
		ref,
		"http://",
	) || strings.HasPrefix(
		ref,
		"https://",
	) {

		return GitHub{
			URL: ref,
		}, nil
	}

	if strings.HasPrefix(
		ref,
		"local:",
	) {

		return Local{
			Path: strings.TrimPrefix(
				ref,
				"local:",
			),
		}, nil
	}

	return nil,
		fmt.Errorf(
			"unknown source: %s",
			ref,
		)
}
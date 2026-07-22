package source

import "fmt"

var registry []Factory

func Register(
	f Factory,
) {
	registry = append(
		registry,
		f,
	)
}

func Resolve(
	ref string,
) (Source, error) {

	for _, factory := range registry {

		if factory.Supports(ref) {
			return factory.New(ref), nil
		}
	}

	return nil, fmt.Errorf(
		"unsupported source: %s",
		ref,
	)
}

package fingerprint

import (
	"github.com/sxijyoti/whiskey/internal/source"
)

func FingerprintSource(
	src source.Source,
) (string, error) {

	data, err := src.Fetch()

	if err != nil {
		return "", err
	}

	return SHA256(data), nil
}

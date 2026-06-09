package fingerprint

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Save(
	path string,
	store Store,
) error {

	if err := os.MkdirAll(
		filepath.Dir(path),
		0755,
	); err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		store,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		path,
		data,
		0644,
	)
}
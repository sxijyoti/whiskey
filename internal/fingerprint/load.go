package fingerprint

import (
	"encoding/json"
	"os"
)

func Load(
	path string,
) (Store, error) {

	store := Store{}

	raw, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		return store, nil
	}

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(
		raw,
		&store,
	); err != nil {
		return nil, err
	}

	return store, nil
}
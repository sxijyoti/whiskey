package planner

import (
	"encoding/json"
	"os"
	"time"
	"path/filepath"
)

type State struct {
	LastBuild time.Time `json:"last_build"`
}

func LoadState(
	path string,
) (*State, error) {

	state := &State{}

	raw, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		return state, nil
	}

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(
		raw,
		state,
	); err != nil {
		return nil, err
	}

	return state, nil
}

func SaveState(
	path string,
	state *State,
) error {

	if err := os.MkdirAll(
		filepath.Dir(path),
		0755,
	); err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		state,
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
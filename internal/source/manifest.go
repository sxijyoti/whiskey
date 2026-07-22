package source

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ManifestVersion = 1

type Manifest struct {
	Version int `json:"version"`

	Sources map[string]ManifestEntry `json:"sources"`
}

type ManifestEntry struct {
	Workspace string `json:"workspace"`

	ContentHash string `json:"content_hash"`

	State map[string]string `json:"state"`
}

func ManifestPath(
	root string,
) string {

	return filepath.Join(
		root,
		".whiskey",
		"manifest.json",
	)
}

func LoadManifest(
	root string,
) (*Manifest, error) {

	path := ManifestPath(root)

	if _, err := os.Stat(path); os.IsNotExist(err) {

		manifest := &Manifest{
			Version: ManifestVersion,
			Sources: make(
				map[string]ManifestEntry,
			),
		}

		if err := SaveManifest(
			root,
			manifest,
		); err != nil {
			return nil, err
		}

		return manifest, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest

	if err := json.Unmarshal(
		raw,
		&manifest,
	); err != nil {
		return nil, err
	}

	if manifest.Sources == nil {

		manifest.Sources = make(
			map[string]ManifestEntry,
		)
	}

	return &manifest, nil
}

func SaveManifest(
	root string,
	manifest *Manifest,
) error {

	path := ManifestPath(root)

	if err := os.MkdirAll(
		filepath.Dir(path),
		0755,
	); err != nil {
		return err
	}

	raw, err := json.MarshalIndent(
		manifest,
		"",
		"    ",
	)
	if err != nil {
		return err
	}

	return os.WriteFile(
		path,
		raw,
		0644,
	)
}

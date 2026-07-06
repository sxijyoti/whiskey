package source

import (
	"crypto/sha256"
	"encoding/hex"
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
	Provider string `json:"provider"`

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

func WorkspaceDir(
	root string,
) string {

	return filepath.Join(
		root,
		".whiskey",
		"workspace",
	)
}

func WorkspaceName(
	ref string,
) string {

	sum := sha256.Sum256(
		[]byte(ref),
	)

	return hex.EncodeToString(
		sum[:],
	) + ".md"
}

func WorkspacePath(
	root string,
	ref string,
) string {

	return filepath.Join(
		WorkspaceDir(root),
		WorkspaceName(ref),
	)
}

func WorkspaceExists(
	root string,
	ref string,
) bool {

	_, err := os.Stat(
		WorkspacePath(root, ref),
	)

	return err == nil
}

func ReadWorkspace(
	root string,
	ref string,
) ([]byte, error) {

	return os.ReadFile(
		WorkspacePath(root, ref),
	)
}

func WriteWorkspace(
	root string,
	ref string,
	data []byte,
) error {

	path := WorkspacePath(
		root,
		ref,
	)

	if err := os.MkdirAll(
		WorkspaceDir(root),
		0755,
	); err != nil {
		return err
	}

	return os.WriteFile(
		path,
		data,
		0644,
	)
}

func DeleteWorkspace(
	root string,
	ref string,
) error {

	return os.Remove(
		WorkspacePath(
			root,
			ref,
		),
	)
}
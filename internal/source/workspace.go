package source

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

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
package source

import (
	"crypto/sha256"
	"encoding/hex"
)

type Materialized struct {
	Workspace string

	ContentHash string

	Metadata Metadata
}

func Materialize(
	root string,
	src Source,
) (*Materialized, error) {

	meta, err := src.Metadata()
	if err != nil {
		return nil, err
	}

	content, err := src.Fetch()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(
		content,
	)

	contentHash := hex.EncodeToString(
		hash[:],
	)

	if err := WriteWorkspace(
		root,
		src.ID(),
		content,
	); err != nil {
		return nil, err
	}

	return &Materialized{
		Workspace: WorkspaceName(
			src.ID(),
		),

		ContentHash: contentHash,

		Metadata: *meta,
	}, nil
}


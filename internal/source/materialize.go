package source

import (
	"crypto/sha256"
	"encoding/hex"
)

func Materialize(
	root string,
	src Source,
	meta *Metadata,
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

	return &Materialized{
		Workspace: WorkspaceName(
			src.ID(),
		),
		Content: content,

		ContentHash: contentHash,

		Metadata: *meta,
	}, nil
}
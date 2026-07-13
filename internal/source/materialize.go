package source

import (
	"crypto/sha256"
	"encoding/hex"
)

// Materialize fetches the content of a remote source, hashes it, and returns
// the result ready for writing to the workspace.
//
// meta is supplied by the caller (obtained from ConditionalMetadata or
// Metadata) and must not be nil. Materialize does NOT call src.Metadata()
// internally to avoid a redundant round-trip.
func Materialize(
	root string,
	src Source,
	meta *Metadata,
) (*Materialized, error) {

	content, err := src.Fetch()
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(content)

	contentHash := hex.EncodeToString(hash[:])

	return &Materialized{
		Workspace: WorkspaceName(
			src.ID(),
		),
		Content: content,

		ContentHash: contentHash,

		Metadata: *meta,
	}, nil
}
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(
	data []byte,
) string {

	sum := sha256.Sum256(data)

	return hex.EncodeToString(
		sum[:],
	)
}

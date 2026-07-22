package fingerprint

type Entry struct {
	Hash         string `json:"hash"`
	ETag         string `json:"etag"`
	LastModified string `json:"last_modified"`
}

type Store map[string]Entry

package fingerprint

func Changed(
	store Store,
	id string,
	hash string,
) bool {

	old, ok := store[id]

	if !ok {
		return true
	}

	return old.Hash != hash
}

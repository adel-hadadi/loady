package utils

func DiffSlices[T comparable](old, new []T) (added, removed []T) {
	oldMap := make(map[T]struct{}, len(old))
	newMap := make(map[T]struct{}, len(new))

	for _, s := range old {
		oldMap[s] = struct{}{}
	}
	for _, s := range new {
		newMap[s] = struct{}{}
	}

	for _, s := range new {
		if _, exists := oldMap[s]; !exists {
			added = append(added, s)
		}
	}

	for _, s := range old {
		if _, exists := newMap[s]; !exists {
			removed = append(removed, s)
		}
	}

	return
}

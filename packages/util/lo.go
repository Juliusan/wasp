package util

func Take[T any](collection []T, predicate func(item T) bool) (T, []T, bool) {
	newCollection := make([]T, 0, len(collection))
	for i, item := range collection {
		if predicate(item) {
			return item, append(newCollection, collection[i+1:]...), true
		}
		newCollection = append(newCollection, item)
	}

	var result T
	return result, newCollection, false
}

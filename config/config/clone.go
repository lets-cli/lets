package config

func cloneArray[I comparable](a []I) []I {
	if a == nil {
		return nil
	}

	arr := make([]I, len(a))
	for idx, item := range a {
		arr[idx] = item
	}

	return arr
}

func cloneMap[K string, V comparable](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}

	mapping := make(map[K]V, len(m))
	for k, v := range m {
		mapping[k] = v
	}

	return mapping
}

func cloneMapArray[K string, V []string](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}

	mapping := make(map[K]V, len(m))
	for k, v := range m {
		mapping[k] = cloneArray(v)
	}

	return mapping
}

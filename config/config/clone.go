package config

import "golang.org/x/exp/constraints"

func cloneSlice[I any](a []I) []I {
	if a == nil {
		return nil
	}

	arr := make([]I, len(a))
	copy(arr, a)

	return arr
}

func cloneMap[K constraints.Ordered, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}

	mapping := make(map[K]V, len(m))
	for k, v := range m {
		mapping[k] = v
	}

	return mapping
}

func cloneMapSlice[K constraints.Ordered, V []string](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}

	mapping := make(map[K]V, len(m))
	for k, v := range m {
		mapping[k] = cloneSlice(v)
	}

	return mapping
}

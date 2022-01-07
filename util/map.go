package util

// GetMapKeys get keys as array
func GetMapKeys(value interface{}) []string {
	mapping, ok := value.(map[interface{}]interface{})
	
	if !ok {
		return []string{}
	}

	keys := make([]string, len(mapping))
	idx := 0

	for key := range mapping {
		key, _ := key.(string)
		keys[idx] = key
		idx++
	}

	return keys
}
package util

// IsStringInList checks if some string is a list element
func IsStringInList(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
}

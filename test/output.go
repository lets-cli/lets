package test

import (
	"strings"
)

// CompareCmdOutput normalizes two string and makes them comparable
// normalization is just replacing all `\n` and `spaces` and joining with one space
// returns 3 values, result if values are same, and two normalized strings
func CompareCmdOutput(one string, another string) (bool, string, string) {
	normalize := func(str string) string {
		return strings.Join(strings.Fields(str), " ")
	}
	oneNorm := normalize(one)
	anotherNorm := normalize(another)

	return oneNorm == anotherNorm, oneNorm, anotherNorm
}

//go:build windows

package util

// IsForegroundProcess reports whether this process may safely query the
// terminal referred to by fd. Windows has no terminal process groups, so
// querying never suspends the process.
func IsForegroundProcess(_ uintptr) bool {
	return true
}

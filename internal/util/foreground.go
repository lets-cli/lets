package util

import (
	"golang.org/x/sys/unix"
)

// IsForegroundProcess reports whether this process belongs to the foreground
// process group of the terminal referred to by fd. Reading or reconfiguring a
// terminal from a background process group suspends the process with
// SIGTTIN/SIGTTOU — e.g. inside process substitution:
// `source <(lets completion -s zsh)`.
func IsForegroundProcess(fd uintptr) bool {
	pgrp, err := unix.IoctlGetInt(int(fd), unix.TIOCGPGRP)
	if err != nil {
		return false
	}

	return pgrp == unix.Getpgrp()
}

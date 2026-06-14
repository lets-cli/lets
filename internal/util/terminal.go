package util

import (
	"io"

	"github.com/charmbracelet/x/term"
)

// IsTerminalWriter reports whether w is connected to an interactive terminal.
// Writers that are not backed by a file descriptor (e.g. buffers used in tests,
// or pipes) are treated as non-terminals.
func IsTerminalWriter(w io.Writer) bool {
	file, ok := w.(term.File)
	return ok && term.IsTerminal(file.Fd())
}

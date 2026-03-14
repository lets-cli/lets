package executor

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

// DependencyError carries the full dependency chain when a command fails.
// Chain is outermost-first (e.g., ["deploy", "build", "lint"]).
type DependencyError struct {
	Chain []string
	Err   error
}

func (e *DependencyError) Error() string { return e.Err.Error() }

// ExitCode propagates the exit code from the innermost ExecuteError, or returns 1.
func (e *DependencyError) ExitCode() int {
	var exitErr *ExecuteError
	if errors.As(e.Err, &exitErr) {
		return exitErr.ExitCode()
	}

	return 1
}

// prependToChain prepends name to the chain in err if err is already a *DependencyError,
// otherwise wraps err in a new single-element DependencyError.
func prependToChain(name string, err error) error {
	var depErr *DependencyError
	if errors.As(err, &depErr) {
		return &DependencyError{Chain: append([]string{name}, depErr.Chain...), Err: depErr.Err}
	}

	return &DependencyError{Chain: []string{name}, Err: err}
}

// PrintDependencyTree writes an indented tree of the dependency chain to w.
// The failing node (last in chain) is annotated in red.
// Respects NO_COLOR automatically via fatih/color.
func PrintDependencyTree(e *DependencyError, w io.Writer) {
	red := color.New(color.FgRed).SprintFunc()

	for i, name := range e.Chain {
		indent := strings.Repeat("  ", i+1)
		if i == len(e.Chain)-1 {
			fmt.Fprintf(w, "%s%s  %s\n", indent, name, red("<-- failed here"))
		} else {
			fmt.Fprintf(w, "%s%s\n", indent, name)
		}
	}
}

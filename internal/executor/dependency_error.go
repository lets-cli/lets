package executor

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

const dependencyTreeIndent = "  "
const dependencyTreeHeader = "command failed:"
const dependencyTreeJoint = "└─ "

// DependencyError carries the full dependency chain when a command fails.
// Chain is outermost-first (e.g., ["deploy", "build", "lint"]).
type DependencyError struct {
	Chain []string
	Err   error
}

func (e *DependencyError) Error() string { return e.Err.Error() }
func (e *DependencyError) Unwrap() error { return e.Err }

// ExitCode propagates the exit code from the innermost ExecuteError, or returns 1.
func (e *DependencyError) ExitCode() int {
	if exitErr, ok := errors.AsType[*ExecuteError](e.Err); ok {
		return exitErr.ExitCode()
	}

	return 1
}

func (e *DependencyError) FailureMessage() string {
	if executeErr, ok := errors.AsType[*ExecuteError](e.Err); ok {
		return executeErr.Cause().Error()
	}

	return e.Err.Error()
}

func (e *DependencyError) TreeMessage() string {
	red := color.New(color.FgRed).SprintFunc()

	var builder strings.Builder

	builder.WriteString(dependencyTreeHeader)

	for i, name := range e.Chain {
		builder.WriteByte('\n')
		builder.WriteString(strings.Repeat(dependencyTreeIndent, i+1))
		builder.WriteString(dependencyTreeJoint)
		builder.WriteString(name)

		if i == len(e.Chain)-1 {
			builder.WriteString(dependencyTreeIndent)
			builder.WriteString(red("<-- failed here"))
		}
	}

	return builder.String()
}

// prependToChain prepends name to the chain in err if err is already a *DependencyError,
// otherwise wraps err in a new single-element DependencyError.
func prependToChain(name string, err error) error {
	if depErr, ok := errors.AsType[*DependencyError](err); ok {
		return &DependencyError{Chain: append([]string{name}, depErr.Chain...), Err: depErr.Err}
	}

	return &DependencyError{Chain: []string{name}, Err: err}
}

// PrintDependencyTree writes an indented tree of the dependency chain to w.
// The failing node (last in chain) is annotated in red.
// Respects NO_COLOR automatically via fatih/color.
func PrintDependencyTree(e *DependencyError, w io.Writer) {
	fmt.Fprintln(w, e.TreeMessage())
}

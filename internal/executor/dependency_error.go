package executor

import "errors"

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

// prependToChain prepends name to the chain in err if err is already a *DependencyError,
// otherwise wraps err in a new single-element DependencyError.
func prependToChain(name string, err error) error {
	if depErr, ok := errors.AsType[*DependencyError](err); ok {
		return &DependencyError{Chain: append([]string{name}, depErr.Chain...), Err: depErr.Err}
	}

	return &DependencyError{Chain: []string{name}, Err: err}
}

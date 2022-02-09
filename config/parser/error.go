package parser

import (
	"fmt"
	"strings"
)

type ParseError struct {
	Directive   string
	CommandName string
	Err         error
}

func (e *ParseError) Error() string {
	if e.CommandName != "" {
		return fmt.Sprintf("failed to parse '%s' command: %s", e.CommandName, e.Err)
	}

	return fmt.Sprintf("failed to parse directive '%s': %s", e.Directive, e.Err)
}

// env is not proper arg.
// TODO refactor meta arg.
func parseError(msg string, name string, field string, meta string) error {
	fields := []string{field}
	if meta != "" {
		fields = append(fields, meta)
	}

	fullPath := strings.Join(fields, ". ")

	return &ParseError{
		CommandName: name,
		Err:         fmt.Errorf("field %s: %s", fullPath, msg),
	}
}

func parseDirectiveError(directive string, msg string) error {
	return &ParseError{
		Directive:   directive,
		CommandName: "",
		Err:         fmt.Errorf("%s", msg),
	}
}

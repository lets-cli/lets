package cmd

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/lets-cli/fang"
	"github.com/lets-cli/lets/internal/executor"
)

// ErrorHandler renders command execution errors using lets' Fang-based style.
func ErrorHandler(w io.Writer, styles fang.Styles, err error) {
	newErrorRenderer(w, styles, err).Render()
}

// errorRenderer owns the full error rendering flow for one error value.
type errorRenderer struct {
	out    errorOutput
	styles fang.Styles
	err    error
}

// newErrorRenderer wires an error, styles, and output helper together.
func newErrorRenderer(w io.Writer, styles fang.Styles, err error) errorRenderer {
	return errorRenderer{out: newErrorOutput(w, styles), styles: styles, err: err}
}

// errorOutput hides low-level writes and common error-specific styled lines.
type errorOutput struct {
	w      io.Writer
	styles fang.Styles
}

// newErrorOutput creates the error output adapter.
func newErrorOutput(w io.Writer, styles fang.Styles) errorOutput {
	return errorOutput{w: w, styles: styles}
}

// writeln writes one error output line and intentionally ignores write errors.
func (o errorOutput) writeln(v ...any) {
	_, _ = fmt.Fprintln(o.w, v...)
}

// blank writes one blank error output line.
func (o errorOutput) blank() {
	o.writeln()
}

// header writes the styled error heading.
func (o errorOutput) header() {
	o.writeln(o.styles.ErrorHeader.String())
}

// commandTreeTitle writes the dependency tree section heading.
func (o errorOutput) commandTreeTitle() {
	title := o.styles.Title.Margin(0, 0).MarginLeft(2).Padding(0, 0)
	o.writeln(title.Render("command tree:"))
}

// Render prints the complete error view or a plain error for non-terminal output.
func (r errorRenderer) Render() {
	if r.shouldRenderPlain() {
		r.out.writeln(r.err.Error())
		return
	}

	errorText := r.styles.ErrorText

	r.out.header()
	r.renderMessage(errorText)
	r.out.blank()

	if depErr := r.dependencyError(); depErr != nil {
		r.renderDependencyTree(depErr)
		r.out.blank()

		return
	}

	if isUsageError(r.err) {
		r.renderUsageHint(errorText)
	}
}

// shouldRenderPlain reports whether styled terminal output should be skipped.
func (r errorRenderer) shouldRenderPlain() bool {
	w, ok := r.out.w.(term.File)
	return ok && !term.IsTerminal(w.Fd())
}

// dependencyError extracts dependency chain context when the error carries it.
func (r errorRenderer) dependencyError() *executor.DependencyError {
	var depErr *executor.DependencyError
	if errors.As(r.err, &depErr) {
		return depErr
	}

	return nil
}

// renderMessage prints the primary error and, when present, its split cause.
func (r errorRenderer) renderMessage(style lipgloss.Style) {
	message, cause := splitExecuteError(r.err)
	r.out.writeln(style.Render(message + "."))

	if cause != "" {
		r.out.writeln(style.UnsetTransform().Render(capitalizeExitStatus(cause) + "."))
	}
}

// renderUsageHint prints the compact help suggestion for usage errors.
func (r errorRenderer) renderUsageHint(errorText lipgloss.Style) {
	r.out.writeln(lipgloss.JoinHorizontal(
		lipgloss.Left,
		errorText.UnsetWidth().Render("Try"),
		" ",
		r.styles.Program.Flag.Render("--help"),
		" for usage.",
	))
	r.out.blank()
}

// renderDependencyTree prints the failed command chain for dependency errors.
func (r errorRenderer) renderDependencyTree(depErr *executor.DependencyError) {
	joint := r.styles.Program.DimmedArgument.Render("└─ ")
	failed := r.styles.ErrorHeader.UnsetMargins().UnsetString().Render("<-- failed here")

	r.out.commandTreeTitle()

	for i, name := range depErr.Chain {
		line := strings.Repeat("  ", i+2) + joint + r.styles.Program.Command.Render(name)
		if i == len(depErr.Chain)-1 {
			line += "  " + failed
		}

		r.out.writeln(line)
	}
}

// capitalizeExitStatus normalizes Go command exit messages for user-facing output.
func capitalizeExitStatus(text string) string {
	if after, ok := strings.CutPrefix(text, "exit status"); ok {
		return "Exit status" + after
	}

	return text
}

// splitExecuteError separates an executor error's message from its process cause.
func splitExecuteError(err error) (string, string) {
	var executeErr *executor.ExecuteError
	if !errors.As(err, &executeErr) {
		return err.Error(), ""
	}

	cause := executeErr.Cause().Error()
	message := strings.TrimSuffix(err.Error(), ": "+cause)

	return message, cause
}

// isUsageError reports whether an error should include a help hint.
func isUsageError(err error) bool {
	s := err.Error()
	for _, prefix := range []string{
		"flag needs an argument:",
		"unknown flag:",
		"unknown shorthand flag:",
		"unknown command",
		"invalid argument",
	} {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}

	return false
}

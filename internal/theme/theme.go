package theme

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/charmbracelet/x/term"
	"github.com/lets-cli/fang"
	"github.com/lets-cli/lets/internal/util"
)

// Supported theme names.
const (
	DefaultName   = "default"
	ANSIName      = "ansi"
	SynthwaveName = "synthwave"
)

// ValidName reports whether name is a supported theme.
func ValidName(name string) bool {
	switch name {
	case DefaultName, ANSIName, SynthwaveName:
		return true
	default:
		return false
	}
}

// ColorSchemeByName resolves a theme name to a Fang color scheme.
func ColorSchemeByName(name string) fang.ColorSchemeFunc {
	switch name {
	case ANSIName:
		return AnsiColorScheme
	case SynthwaveName:
		return SynthwaveColorScheme
	default:
		return DefaultColorScheme
	}
}

// IsDarkBackground reports whether the terminal behind out has a dark background.
// Detection queries the terminal (writes OSC 11 to out, reads the reply from stdin),
// which suspends a background process group with SIGTTIN/SIGTTOU — so when the
// process is not in the foreground, skip the query and assume dark (lipgloss's
// own fallback).
func IsDarkBackground(out term.File) bool {
	if !util.IsForegroundProcess(out.Fd()) {
		return true
	}

	return lipgloss.HasDarkBackground(os.Stdin, out)
}

// ProgressColorsByName resolves a theme name to fill and empty colors for download progress.
func ProgressColorsByName(name string, out term.File) (color.Color, color.Color) {
	scheme := ColorSchemeByName(name)(lipgloss.LightDark(IsDarkBackground(out)))

	return ProgressColors(scheme)
}

// ProgressColors resolves a Fang color scheme to fill and empty colors for download progress.
func ProgressColors(scheme fang.ColorScheme) (color.Color, color.Color) {
	fill := scheme.Command
	if fill == nil {
		fill = scheme.Flag
	}

	if fill == nil {
		fill = scheme.Base
	}

	empty := scheme.Comment
	if empty == nil {
		empty = scheme.FlagDefault
	}

	if empty == nil {
		empty = scheme.Base
	}

	return fill, empty
}

// DefaultColorScheme is the default colorscheme.
func DefaultColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	baseCyan := charmtone.Turtle
	baseCyanLighter := lipgloss.Color("#0BF4F1")
	baseWhite := charmtone.Ash
	baseGray := charmtone.Oyster

	return fang.ColorScheme{
		Base:           c(charmtone.Charcoal, baseCyanLighter),
		Title:          baseCyanLighter,
		Codeblock:      c(charmtone.Salt, lipgloss.Color("#2F2E36")),
		Program:        c(charmtone.Malibu, baseCyan),
		Command:        c(charmtone.Pony, baseCyan),
		DimmedArgument: c(charmtone.Squid, baseGray),
		Comment:        c(charmtone.Squid, lipgloss.Color("#747282")),
		Flag:           c(lipgloss.Color("#0CB37F"), baseCyan),
		Argument:       c(charmtone.Charcoal, baseWhite),
		Description:    c(charmtone.Charcoal, baseWhite),
		FlagDefault:    c(charmtone.Smoke, charmtone.Squid),
		QuotedString:   c(charmtone.Coral, baseCyan),
		ErrorHeader: [2]color.Color{
			charmtone.Butter,
			charmtone.Sriracha,
		},
	}
}

// AnsiColorScheme is a ANSI colorscheme.
func AnsiColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	base := c(lipgloss.Black, lipgloss.White)

	return fang.ColorScheme{
		Base:         base,
		Title:        lipgloss.White,
		Description:  base,
		Comment:      c(lipgloss.BrightWhite, lipgloss.BrightBlack),
		Flag:         lipgloss.White,
		FlagDefault:  lipgloss.White,
		Command:      lipgloss.White,
		QuotedString: lipgloss.White,
		Argument:     base,
		Help:         base,
		Dash:         base,
		ErrorHeader:  [2]color.Color{lipgloss.Black, lipgloss.Red},
		ErrorDetails: lipgloss.Red,
	}
}

func SynthwaveColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	basePink := lipgloss.Color("#f24db8")

	return fang.ColorScheme{
		Base:           c(charmtone.Charcoal, basePink),
		Title:          charmtone.Grape,
		Codeblock:      c(charmtone.Salt, lipgloss.Color("#2F2E36")),
		Program:        c(charmtone.Malibu, charmtone.Grape),
		Command:        c(charmtone.Pony, basePink),
		DimmedArgument: c(charmtone.Squid, charmtone.Oyster),
		Comment:        c(charmtone.Squid, lipgloss.Color("#747282")),
		Flag:           c(lipgloss.Color("#0CB37F"), basePink),
		Argument:       c(charmtone.Charcoal, charmtone.Ash),
		Description:    c(charmtone.Charcoal, charmtone.Ash),
		FlagDefault:    c(charmtone.Smoke, charmtone.Squid),
		QuotedString:   c(charmtone.Coral, charmtone.Cheeky),
		ErrorHeader: [2]color.Color{
			charmtone.Butter,
			charmtone.Sriracha,
		},
	}
}

package theme

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/lets-cli/fang"
)

// DefaultColorScheme is the default colorscheme.
func DefaultColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	baseCyan := charmtone.Turtle
	baseWhite := charmtone.Ash
	baseGray := charmtone.Oyster
	return fang.ColorScheme{
		Base:           c(charmtone.Charcoal, baseCyan),
		Title:          baseWhite,
		Codeblock:      c(charmtone.Salt, lipgloss.Color("#2F2E36")),
		Program:        c(charmtone.Malibu, baseWhite),
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
	return fang.ColorScheme{
		Base:           c(charmtone.Charcoal, charmtone.Cheeky),
		Title:          charmtone.Grape,
		Codeblock:      c(charmtone.Salt, lipgloss.Color("#2F2E36")),
		Program:        c(charmtone.Malibu, charmtone.Grape),
		Command:        c(charmtone.Pony, charmtone.Cheeky),
		DimmedArgument: c(charmtone.Squid, charmtone.Oyster),
		Comment:        c(charmtone.Squid, lipgloss.Color("#747282")),
		Flag:           c(lipgloss.Color("#0CB37F"), charmtone.Cheeky),
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

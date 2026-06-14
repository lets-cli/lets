package theme

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"
)

func TestValidName(t *testing.T) {
	for _, name := range []string{DefaultName, ANSIName, SynthwaveName} {
		if !ValidName(name) {
			t.Fatalf("expected %q to be valid", name)
		}
	}

	if ValidName("vaporwave") {
		t.Fatal("expected unknown theme to be invalid")
	}
}

func TestColorSchemeByName(t *testing.T) {
	if got := ColorSchemeByName(DefaultName)(lipgloss.LightDark(true)).Title; got != charmtone.Ash {
		t.Fatalf("expected default theme title color %v, got %v", charmtone.Ash, got)
	}

	if got := ColorSchemeByName(ANSIName)(lipgloss.LightDark(true)).ErrorDetails; got != lipgloss.Red {
		t.Fatalf("expected ansi theme error details color %v, got %v", lipgloss.Red, got)
	}

	if got := ColorSchemeByName(SynthwaveName)(lipgloss.LightDark(true)).Title; got != charmtone.Grape {
		t.Fatalf("expected synthwave theme title color %v, got %v", charmtone.Grape, got)
	}

	if got := ColorSchemeByName("unknown")(lipgloss.LightDark(true)).Title; got != charmtone.Ash {
		t.Fatalf("expected unknown theme to fall back to %v, got %v", charmtone.Ash, got)
	}
}

func TestProgressColors(t *testing.T) {
	fill, empty := ProgressColors(ColorSchemeByName(DefaultName)(lipgloss.LightDark(true)))
	if fill != charmtone.Turtle {
		t.Fatalf("expected default fill color %v, got %v", charmtone.Turtle, fill)
	}
	if empty != lipgloss.Color("#747282") {
		t.Fatalf("expected default empty color %v, got %v", lipgloss.Color("#747282"), empty)
	}

	fill, empty = ProgressColors(ColorSchemeByName(ANSIName)(lipgloss.LightDark(true)))
	if fill != lipgloss.White {
		t.Fatalf("expected ansi fill color %v, got %v", lipgloss.White, fill)
	}
	if empty != lipgloss.BrightBlack {
		t.Fatalf("expected ansi empty color %v, got %v", lipgloss.BrightBlack, empty)
	}
}

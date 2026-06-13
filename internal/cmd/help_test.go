package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/lets-cli/fang"
	configpkg "github.com/lets-cli/lets/internal/config/config"
	"github.com/spf13/cobra"
)

func TestHelpRendererShowsSubgroups(t *testing.T) {
	root := CreateRootCommand("v0.0.0-test", "")
	root.InitDefaultHelpFlag()
	root.SetArgs(nil)
	root.SetOut(new(bytes.Buffer))
	root.SetErr(new(bytes.Buffer))

	root.AddCommand(&cobra.Command{
		Use:     "build",
		Short:   "build stuff",
		GroupID: "main",
		RunE: func(*cobra.Command, []string) error {
			return nil
		},
		Annotations: map[string]string{
			annotationSubGroupName: "Development",
		},
	})
	root.AddCommand(&cobra.Command{
		Use:     "deploy",
		Short:   "deploy stuff",
		GroupID: "main",
		RunE: func(*cobra.Command, []string) error {
			return nil
		},
		Annotations: map[string]string{
			annotationSubGroupName: "Operations",
		},
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(new(bytes.Buffer))

	if err := fang.Execute(context.Background(), root, fang.WithHelpRenderer(HelpRenderer)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "Development") {
		t.Fatalf("expected Development subgroup in output: %q", out)
	}
	if !strings.Contains(out, "Operations") {
		t.Fatalf("expected Operations subgroup in output: %q", out)
	}
}

func TestHelpRendererUsesDocoptFlags(t *testing.T) {
	root := CreateRootCommand("v0.0.0-test", "")
	root.InitDefaultHelpFlag()
	root.SetArgs([]string{"help", "release"})

	release := newSubcommand(&configpkg.Command{
		Name:        "release",
		GroupName:   "Common",
		Description: "Create tag and push",
		Docopts: `Usage: lets release <version> --message=<message>

Options:
  <version>                  Set version (e.g. 1.0.0)
  --message=<message>, -m    Release message

Example:
  lets release 1.0.0 -m "Release 1.0.0"`,
	}, nil, false, nil)
	root.AddCommand(release)

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(new(bytes.Buffer))

	if err := fang.Execute(context.Background(), root, fang.WithHelpRenderer(HelpRenderer)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "lets release <version> --message=<message>") {
		t.Fatalf("expected usage in output: %q", out)
	}
	if !strings.Contains(out, "OPTIONS") {
		t.Fatalf("expected options title in output: %q", out)
	}
	argIdx := strings.Index(out, "<version>")
	flagIdx := strings.Index(out, "--message=<message>, -m")
	if argIdx == -1 {
		t.Fatalf("expected docopt argument in output: %q", out)
	}
	if !strings.Contains(out, "--message=<message>, -m") {
		t.Fatalf("expected docopt flag in output: %q", out)
	}
	if argIdx > flagIdx {
		t.Fatalf("expected docopt argument before flag in output: %q", out)
	}
	if strings.Contains(out, "-h --help") {
		t.Fatalf("did not expect help flag in output: %q", out)
	}
}

func TestHelpRendererUsesDocoptUsageInRootCommandList(t *testing.T) {
	root := CreateRootCommand("v0.0.0-test", "")
	root.InitDefaultHelpFlag()
	root.SetArgs(nil)

	release := newSubcommand(&configpkg.Command{
		Name:        "release",
		GroupName:   "Common",
		Description: "Create tag and push",
		Docopts: `Usage: lets release <version> --message=<message>

Options:
  --message=<message>, -m    Release message`,
	}, nil, false, nil)
	root.AddCommand(release)

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(new(bytes.Buffer))

	if err := fang.Execute(context.Background(), root, fang.WithHelpRenderer(HelpRenderer)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(stdout.String(), "release <version> --message=<message>") {
		t.Fatalf("expected command usage in output: %q", stdout.String())
	}
}

func TestHelpRendererDoesNotDuplicateDocoptUsageCommandPath(t *testing.T) {
	root := CreateRootCommand("v0.0.0-test", "")
	root.InitDefaultHelpFlag()
	root.SetArgs([]string{"help", "build"})

	build := newSubcommand(&configpkg.Command{
		Name:        "build",
		GroupName:   "Common",
		Description: "Build lets from source code",
		Docopts: `Usage:
  lets build
  lets build [<bin>]`,
	}, nil, false, nil)
	root.AddCommand(build)

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(new(bytes.Buffer))

	if err := fang.Execute(context.Background(), root, fang.WithHelpRenderer(HelpRenderer)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := stdout.String()
	if strings.Contains(out, "lets build lets build") {
		t.Fatalf("did not expect duplicated command path in output: %q", out)
	}
	if !strings.Contains(out, "lets build [<bin>]") {
		t.Fatalf("expected docopt usage in output: %q", out)
	}
}

func TestErrorHandlerRemovesErrorHeaderLeftPadding(t *testing.T) {
	var stderr bytes.Buffer
	styles := fang.Styles{
		ErrorHeader: lipgloss.NewStyle().Padding(0, 1).SetString("ERROR"),
		ErrorText:   lipgloss.NewStyle(),
		Program: fang.Program{
			Flag: lipgloss.NewStyle(),
		},
	}

	ErrorHandler(&stderr, styles, &unknownCommandError{message: `unknown command "wat" for "lets"`})

	out := stderr.String()
	if !strings.Contains(out, "ERROR") {
		t.Fatalf("expected error header in output: %q", out)
	}
	if !strings.Contains(out, "\nTry --help for usage.") {
		t.Fatalf("expected usage hint in output: %q", out)
	}
}

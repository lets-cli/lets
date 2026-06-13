package cmd

import (
	"encoding/json"
	"testing"

	configpkg "github.com/lets-cli/lets/internal/config/config"
	"github.com/lets-cli/lets/internal/docopt"
	"github.com/spf13/cobra"
)

func TestNewSubcommandUsesDocoptHelpMetadata(t *testing.T) {
	command := &configpkg.Command{
		Name:        "release",
		GroupName:   "Common",
		Description: "Create tag and push",
		Docopts: `Usage: lets release <version> --message=<message>

Options:
  <version>                  Set version (e.g. 1.0.0)
  --message=<message>, -m    Release message

Example:
  lets release 1.0.0 -m "Release 1.0.0"`,
	}

	subCmd := newSubcommand(command, nil, false, nil)
	if subCmd.Use != "release" {
		t.Fatalf("unexpected use: %q", subCmd.Use)
	}
	if subCmd.Name() != "release" {
		t.Fatalf("unexpected name: %q", subCmd.Name())
	}
	if subCmd.Annotations[annotationHelpUsage] != "release <version> --message=<message>" {
		t.Fatalf("unexpected help usage: %q", subCmd.Annotations[annotationHelpUsage])
	}

	if subCmd.Example != "  lets release 1.0.0 -m \"Release 1.0.0\"" {
		t.Fatalf("unexpected example: %q", subCmd.Example)
	}

	payload := subCmd.Annotations[annotationHelpOptions]
	if payload == "" {
		t.Fatal("expected custom flags annotation")
	}

	var helpFlags []docopt.HelpOption
	if err := json.Unmarshal([]byte(payload), &helpFlags); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if len(helpFlags) != 2 {
		t.Fatalf("expected 2 help flags, got %d", len(helpFlags))
	}

	if helpFlags[0].Display != "<version>" {
		t.Fatalf("unexpected first display: %q", helpFlags[0].Display)
	}
	if helpFlags[1].Display != "--message=<message>, -m" {
		t.Fatalf("unexpected second display: %q", helpFlags[1].Display)
	}
}

func TestNewSubcommandKeepsConfiguredNameWhenDocoptUsageDiffers(t *testing.T) {
	command := &configpkg.Command{
		Name:      "options-wrong-usage",
		GroupName: "Common",
		Docopts: `Usage: lets options-wrong-usage-xxx

Options:
  --message=<message>, -m    Release message`,
	}

	root := CreateRootCommand("v0.0.0-test", "")
	subCmd := newSubcommand(command, nil, false, nil)
	root.AddCommand(subCmd)

	found, _, err := root.Find([]string{"options-wrong-usage"})
	if err != nil {
		t.Fatalf("unexpected find error: %v", err)
	}
	if found.Name() != "options-wrong-usage" {
		t.Fatalf("unexpected command name: %q", found.Name())
	}
	if found.Annotations[annotationHelpUsage] != "options-wrong-usage-xxx" {
		t.Fatalf("unexpected help usage: %q", found.Annotations[annotationHelpUsage])
	}
}

func TestSubcommandHelpArg(t *testing.T) {
	command := &configpkg.Command{
		Name:        "release",
		GroupName:   "Common",
		Description: "Create tag and push",
	}

	root := CreateRootCommand("v0.0.0-test", "")
	subCmd := newSubcommand(command, nil, false, nil)
	root.AddCommand(subCmd)
	root.SetArgs([]string{"release", "--help"})

	called := false
	subCmd.SetHelpFunc(func(*cobra.Command, []string) {
		called = true
	})

	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected subcommand help to be called")
	}
}

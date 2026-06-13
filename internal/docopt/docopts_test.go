package docopt

import "testing"

const releaseDocopt = `Usage: lets release <version> --message=<message>

Options:
  <version>                  Set version (e.g. 1.0.0)
  --message=<message>, -m    Release message

Example:
  lets release 1.0.0 -m "Release 1.0.0"
  lets release 1.0.0-rc1 -m "Prerelease 1.0.0-rc1"`

func TestParseDocoptParts(t *testing.T) {
	parts := ParseDocoptParts(releaseDocopt)

	if parts.Usage != "lets release <version> --message=<message>" {
		t.Fatalf("unexpected usage: %q", parts.Usage)
	}

	expectedOptions := "  <version>                  Set version (e.g. 1.0.0)\n  --message=<message>, -m    Release message"
	if parts.Options != expectedOptions {
		t.Fatalf("unexpected options: %q", parts.Options)
	}

	expectedExample := "  lets release 1.0.0 -m \"Release 1.0.0\"\n  lets release 1.0.0-rc1 -m \"Prerelease 1.0.0-rc1\""
	if parts.Example != expectedExample {
		t.Fatalf("unexpected example: %q", parts.Example)
	}
}

func TestParseHelpOptions(t *testing.T) {
	helpOptions := ParseHelpOptions(releaseDocopt, "release")
	if len(helpOptions) != 2 {
		t.Fatalf("expected 2 help options, got %d", len(helpOptions))
	}

	if helpOptions[0].Display != "<version>" {
		t.Fatalf("unexpected first display: %q", helpOptions[0].Display)
	}
	if helpOptions[0].Description != "Set version (e.g. 1.0.0)" {
		t.Fatalf("unexpected first description: %q", helpOptions[0].Description)
	}
	if helpOptions[0].Kind != "arg" {
		t.Fatalf("unexpected first kind: %q", helpOptions[0].Kind)
	}
	if helpOptions[0].Name != "<version>" {
		t.Fatalf("unexpected first name: %q", helpOptions[0].Name)
	}

	if helpOptions[1].Display != "--message=<message>, -m" {
		t.Fatalf("unexpected second display: %q", helpOptions[1].Display)
	}
	if helpOptions[1].Description != "Release message" {
		t.Fatalf("unexpected second description: %q", helpOptions[1].Description)
	}
	if helpOptions[1].Kind != "flag" {
		t.Fatalf("unexpected second kind: %q", helpOptions[1].Kind)
	}
	if helpOptions[1].Name != "--message" {
		t.Fatalf("unexpected second name: %q", helpOptions[1].Name)
	}
	if helpOptions[1].Short != "-m" {
		t.Fatalf("unexpected second short: %q", helpOptions[1].Short)
	}
	if helpOptions[1].Long != "--message" {
		t.Fatalf("unexpected second long: %q", helpOptions[1].Long)
	}
}

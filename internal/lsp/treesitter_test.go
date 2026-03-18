package lsp

import (
	"reflect"
	"testing"

	ts "github.com/odvcencio/gotreesitter"
	"github.com/tliron/commonlog"
	lsp "github.com/tliron/glsp/protocol_3_16"
)

var logger = commonlog.GetLogger("test")

func pos(line, character uint32) lsp.Position {
	return lsp.Position{
		Line:      line,
		Character: character,
	}
}

func TestIsCursorWithinNode(t *testing.T) {
	tests := []struct {
		startPoint ts.Point
		endPoint   ts.Point
		pos        lsp.Position
		want       bool
	}{
		// Single line cases
		{ts.Point{1, 0}, ts.Point{1, 10}, lsp.Position{Line: 0, Character: 0}, false},  // cursor not on line
		{ts.Point{1, 0}, ts.Point{1, 10}, lsp.Position{Line: 1, Character: 0}, true},   // cursor at start
		{ts.Point{1, 0}, ts.Point{1, 10}, lsp.Position{Line: 1, Character: 10}, true},  // cursor at end
		{ts.Point{1, 0}, ts.Point{1, 10}, lsp.Position{Line: 1, Character: 11}, false}, // cursor outside

		// Multiple line cases
		{ts.Point{1, 0}, ts.Point{3, 10}, lsp.Position{Line: 2, Character: 10}, true}, // mid line, len <= end
		{ts.Point{1, 0}, ts.Point{3, 10}, lsp.Position{Line: 2, Character: 15}, true}, // mid line, len > end
		{ts.Point{1, 0}, ts.Point{3, 10}, lsp.Position{Line: 3, Character: 10}, true}, // at last line
		{ts.Point{1, 0}, ts.Point{3, 10}, lsp.Position{Line: 4, Character: 1}, false}, // beyond node
	}

	for i, tt := range tests {
		got := isCursorWithinNodePoints(tt.startPoint, tt.endPoint, tt.pos)
		if got != tt.want {
			t.Errorf("case %d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestDetectMixinsPosition(t *testing.T) {
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test`

	tests := []struct {
		pos  lsp.Position
		want bool
	}{
		{lsp.Position{Line: 0, Character: 0}, false},
		{lsp.Position{Line: 1, Character: 0}, true},
		{lsp.Position{Line: 2, Character: 0}, true},
		{lsp.Position{Line: 2, Character: 15}, true},
		{lsp.Position{Line: 3, Character: 0}, false},
	}

	p := newParser(logger)
	for i, tt := range tests {
		got := p.inMixinsPosition(&doc, tt.pos)
		if got != tt.want {
			t.Errorf("case %d: got %v, want %v", i, got, tt.want)
		}
	}
}

func TestExtractFilenameFromMixins(t *testing.T) {
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test`

	tests := []struct {
		pos      lsp.Position
		expected string
	}{
		{lsp.Position{Line: 1, Character: 0}, ""},
		{lsp.Position{Line: 2, Character: 0}, "lets.my.yaml"},
		{lsp.Position{Line: 2, Character: 15}, "lets.my.yaml"},
	}

	p := newParser(logger)
	for i, tt := range tests {
		result := p.extractFilenameFromMixins(&doc, tt.pos)
		if result != tt.expected {
			t.Errorf("Case %d: expected %v, actual %v", i, tt.expected, result)
		}
	}
}

func TestDetectDependsNodeIsArray(t *testing.T) {
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test

  test2:
    depends: [test]
    cmd: echo Test2`

	tests := []struct {
		pos  lsp.Position
		want bool
	}{
		{lsp.Position{Line: 8, Character: 15}, true},
	}

	p := newParser(logger)
	for i, tt := range tests {
		got := p.inDependsPosition(&doc, tt.pos)
		if got != tt.want {
			t.Errorf("case %d: expected %v, actual %v", i, tt.want, got)
		}
	}
}

func TestDetectDependsNodeIsSequence(t *testing.T) {
	// NOTE: space after '-' in first depends sequence item is importanat,
	// it is a curosor position
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test

  test2:
    depends:
      - 
    cmd: echo Test2`

	tests := []struct {
		pos  lsp.Position
		want bool
	}{
		{lsp.Position{Line: 8, Character: 4}, false},
		{lsp.Position{Line: 9, Character: 0}, true},
		{lsp.Position{Line: 9, Character: 7}, true},
		{lsp.Position{Line: 9, Character: 8}, true},
		{lsp.Position{Line: 10, Character: 0}, false},
	}

	p := newParser(logger)
	for i, tt := range tests {
		got := p.inDependsPosition(&doc, tt.pos)
		if got != tt.want {
			t.Errorf("case %d: expected %v, actual %v", i, tt.want, got)
		}
	}
}

func TestDetectDependsNodeIsSequenceNextLine(t *testing.T) {
	// NOTE: space after '-' in first depends sequence item is importanat,
	// it is a curosor position
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test

  test2:
    depends:
      - test
      - 
    cmd: echo Test2`

	tests := []struct {
		pos  lsp.Position
		want bool
	}{
		{lsp.Position{Line: 8, Character: 4}, false},
		{lsp.Position{Line: 9, Character: 0}, true},
		{lsp.Position{Line: 10, Character: 0}, true},
	}

	p := newParser(logger)
	for i, tt := range tests {
		got := p.inDependsPosition(&doc, tt.pos)
		if got != tt.want {
			t.Errorf("case %d: expected %v, actual %v", i, tt.want, got)
		}
	}
}

func TestGetCommands(t *testing.T) {
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test
  test2:
    cmd: echo Test2`

	p := newParser(logger)
	commands := p.getCommands(&doc)
	if len(commands) != 2 {
		t.Errorf("expected 2 commands, got %d", len(commands))
	}

	expected := []Command{
		{name: "test"},
		{name: "test2"},
	}

	for i, cmd := range commands {
		if cmd.name != expected[i].name {
			t.Errorf("command %d: expected name %q, got %q", i, expected[i].name, cmd.name)
		}
	}
}

func TestGetCurrentCommandInDepends(t *testing.T) {
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test
  test2:
    cmd: echo Test2
  test3:
    depends: [test, ]
    cmd: echo Test3`

	p := newParser(logger)
	command := p.getCurrentCommand(&doc, lsp.Position{Line: 9, Character: 20})
	if command == nil {
		t.Fatal("expected command, got nil")
	}
	if command.name != "test3" {
		t.Errorf("expected command name 'test3', got %q", command.name)
	}
}

func TestExractDependsValuesArray(t *testing.T) {
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test
  test2:
    cmd: echo Test2
  test3:
    depends: [test]
    cmd: echo Test3`

	p := newParser(logger)
	values := p.extractDependsValues(&doc)

	if len(values) == 0 {
		t.Fatal("expected non-empty array")
	}

	if !reflect.DeepEqual(values, []string{"test"}) {
		t.Errorf("expected array [test], got %v", values)
	}
}

func TestExractDependsValuesSequence(t *testing.T) {
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test
  test2:
    cmd: echo Test2
  test3:
    depends:
      - test
    cmd: echo Test3`

	p := newParser(logger)
	values := p.extractDependsValues(&doc)
	if len(values) == 0 {
		t.Fatal("expected non-empty array")
	}

	if !reflect.DeepEqual(values, []string{"test"}) {
		t.Errorf("expected array [test], got %v", values)
	}
}

func TestFindCommand(t *testing.T) {
	doc := `shell: bash
mixins:
  - lets.my.yaml
commands:
  test:
    cmd: echo Test
  test2:
    depends:
      - test
    cmd: echo Test3`

	expected := Command{
		name: "test",
		position: lsp.Position{
			Line:      4,
			Character: 2,
		},
	}

	p := newParser(logger)
	command := p.findCommand(&doc, "test")
	if command == nil {
		t.Fatal("expected non-nil command")
	}

	if command.name != expected.name {
		t.Errorf("expected command name '%s', got %q", expected.name, command.name)
	}

	if command.position.Line != expected.position.Line {
		t.Errorf("expected line %d, got %d", expected.position.Line, command.position.Line)
	}
	if command.position.Character != expected.position.Character {
		t.Errorf("expected character %d, got %d", expected.position.Character, command.position.Character)
	}
}

func TestMixinsHelpersWithMultipleItems(t *testing.T) {
	blockDoc := `shell: bash
mixins:
  - lets.base.yaml
  - lets.extra.yaml
commands:
  build:
    cmd: echo build`

	flowDoc := `shell: bash
mixins: [lets.base.yaml, lets.extra.yaml]
commands:
  build:
    cmd: echo build`

	tests := []struct {
		name         string
		doc          string
		position     lsp.Position
		wantInMixins bool
		wantFilename string
	}{
		{
			name:         "block key line is inside mixins",
			doc:          blockDoc,
			position:     pos(1, 1),
			wantInMixins: true,
		},
		{
			name:         "block first item resolves filename",
			doc:          blockDoc,
			position:     pos(2, 4),
			wantInMixins: true,
			wantFilename: "lets.base.yaml",
		},
		{
			name:         "block second item resolves filename",
			doc:          blockDoc,
			position:     pos(3, 10),
			wantInMixins: true,
			wantFilename: "lets.extra.yaml",
		},
		{
			name:         "outside mixins is false",
			doc:          blockDoc,
			position:     pos(4, 0),
			wantInMixins: false,
		},
		{
			name:         "flow mixins are not matched by query",
			doc:          flowDoc,
			position:     pos(1, 12),
			wantInMixins: false,
		},
	}

	p := newParser(logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInMixins := p.inMixinsPosition(&tt.doc, tt.position)
			if gotInMixins != tt.wantInMixins {
				t.Fatalf("inMixinsPosition() = %v, want %v", gotInMixins, tt.wantInMixins)
			}

			gotFilename := p.extractFilenameFromMixins(&tt.doc, tt.position)
			if gotFilename != tt.wantFilename {
				t.Fatalf("extractFilenameFromMixins() = %q, want %q", gotFilename, tt.wantFilename)
			}
		})
	}
}

func TestDependsHelpersWithBlockAndFlowSequences(t *testing.T) {
	doc := `shell: bash
commands:
  build:
    depends:
      - clean
      - lint
    env:
      GOFLAGS: -mod=mod
    cmd: echo build
  test:
    depends: [build, lint]
    options: |
      Usage: lets test [--watch]
    cmd: echo test
  package:
    cmd: echo package`

	tests := []struct {
		name string
		pos  lsp.Position
		want bool
	}{
		{name: "block key line is not inside depends values", pos: pos(3, 4), want: false},
		{name: "block first item", pos: pos(4, 8), want: true},
		{name: "block second item", pos: pos(5, 8), want: true},
		{name: "nested env is outside depends", pos: pos(6, 4), want: false},
		{name: "flow first item", pos: pos(10, 14), want: true},
		{name: "flow second item", pos: pos(10, 21), want: true},
		{name: "nested options are outside depends", pos: pos(12, 12), want: false},
		{name: "command without depends", pos: pos(15, 10), want: false},
	}

	p := newParser(logger)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.inDependsPosition(&doc, tt.pos)
			if got != tt.want {
				t.Fatalf("inDependsPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandHelpersWithDifferentCommandShapes(t *testing.T) {
	doc := `shell: bash
commands:
  bootstrap:
    cmd: echo bootstrap
  build:
    depends:
      - bootstrap
      - lint
    env:
      GOFLAGS: -mod=mod
    cmd: |
      echo build
      echo done
  test:
    depends: [build, lint]
    options: |
      Usage: lets test [--watch]
    cmd: echo test
  lint:
    cmd: echo lint`

	p := newParser(logger)

	expectedCommands := []Command{
		{name: "bootstrap", position: pos(2, 2)},
		{name: "build", position: pos(4, 2)},
		{name: "test", position: pos(13, 2)},
		{name: "lint", position: pos(18, 2)},
	}

	commands := p.getCommands(&doc)
	if !reflect.DeepEqual(commands, expectedCommands) {
		t.Fatalf("getCommands() = %#v, want %#v", commands, expectedCommands)
	}

	findTests := []struct {
		name    string
		command string
		want    *Command
	}{
		{name: "find bootstrap", command: "bootstrap", want: &Command{name: "bootstrap", position: pos(2, 2)}},
		{name: "find build", command: "build", want: &Command{name: "build", position: pos(4, 2)}},
		{name: "find test", command: "test", want: &Command{name: "test", position: pos(13, 2)}},
		{name: "find lint", command: "lint", want: &Command{name: "lint", position: pos(18, 2)}},
		{name: "missing command", command: "missing", want: nil},
	}

	for _, tt := range findTests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.findCommand(&doc, tt.command)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("findCommand() = %#v, want %#v", got, tt.want)
			}
		})
	}

	currentTests := []struct {
		name     string
		position lsp.Position
		want     *Command
	}{
		{name: "inside bootstrap command body", position: pos(3, 12), want: &Command{name: "bootstrap"}},
		{name: "inside build env block", position: pos(9, 10), want: &Command{name: "build"}},
		{name: "inside build multiline cmd", position: pos(11, 8), want: &Command{name: "build"}},
		{name: "inside test flow depends", position: pos(14, 18), want: &Command{name: "test"}},
		{name: "inside test options block", position: pos(16, 12), want: &Command{name: "test"}},
		{name: "inside lint command body", position: pos(19, 10), want: &Command{name: "lint"}},
		{name: "outside commands tree", position: pos(0, 0), want: nil},
	}

	for _, tt := range currentTests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.getCurrentCommand(&doc, tt.position)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("getCurrentCommand() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExtractDependsValuesFromMixedCommands(t *testing.T) {
	doc := `shell: bash
commands:
  build:
    depends:
      - bootstrap
      - lint
    cmd: echo build
  test:
    depends: [build, lint]
    cmd: echo test
  release:
    env:
      TARGET: prod
    depends:
      - test
    cmd: echo release
  lint:
    cmd: echo lint`

	p := newParser(logger)
	got := p.extractDependsValues(&doc)
	want := []string{"bootstrap", "lint", "build", "lint", "test"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("extractDependsValues() = %#v, want %#v", got, want)
	}
}

func TestWordUnderCursor(t *testing.T) {
	tests := []struct {
		line     string
		position lsp.Position
		want     string
	}{
		{"test word here", lsp.Position{Character: 0}, "test"},
		{"test word here", lsp.Position{Character: 2}, "test"},
		{"test word here", lsp.Position{Character: 5}, "word"},
		{"test word here", lsp.Position{Character: 10}, "here"},
		{"test-word_123", lsp.Position{Character: 5}, "test-word_123"},
		{"", lsp.Position{Character: 0}, ""},
		{"test", lsp.Position{Character: 10}, ""},
		{"test word", lsp.Position{Character: 4}, ""},
		{"  test  ", lsp.Position{Character: 3}, "test"},
	}

	for i, tt := range tests {
		got := wordUnderCursor(tt.line, &tt.position)
		if got != tt.want {
			t.Errorf("case %d: got %q, want %q", i, got, tt.want)
		}
	}
}

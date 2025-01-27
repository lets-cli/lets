package lsp

import (
	"reflect"
	"testing"

	"github.com/tliron/commonlog"
	lsp "github.com/tliron/glsp/protocol_3_16"
	ts "github.com/tree-sitter/go-tree-sitter"
)

var logger = commonlog.GetLogger("test")

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

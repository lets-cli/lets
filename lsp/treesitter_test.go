package lsp

import (
	"testing"

	ts "github.com/tree-sitter/go-tree-sitter"
	lsp "github.com/tliron/glsp/protocol_3_16"
)


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
		{ts.Point{1, 0}, ts.Point{3, 10}, lsp.Position{Line: 2, Character: 10}, true},  // mid line, len <= end
		{ts.Point{1, 0}, ts.Point{3, 10}, lsp.Position{Line: 2, Character: 15}, true},  // mid line, len > end
		{ts.Point{1, 0}, ts.Point{3, 10}, lsp.Position{Line: 3, Character: 10}, true},  // at last line
		{ts.Point{1, 0}, ts.Point{3, 10}, lsp.Position{Line: 4, Character: 1}, false},  // beyond node
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

	for i, tt := range tests {
		got := inMixinsPosition(&doc, tt.pos)
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

	for i, tt := range tests {
		result := extractFilenameFromMixins(&doc, tt.pos)
		if result != tt.expected {
			t.Errorf("Case %d: expected %v, actual %v", i, tt.expected, result)
		}
	}
}

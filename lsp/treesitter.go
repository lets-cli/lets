package lsp

import (
	"context"

	ts "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/yaml"
	lsp "github.com/tliron/glsp/protocol_3_16"
)

type PositionType int

const (
    PositionTypeMixins PositionType = iota
    PositionTypeNone
)

func isCursorWithinNode(node *ts.Node, pos lsp.Position) bool {
	return isCursorWithinNodePoints(node.StartPoint(), node.EndPoint(), pos)
}

func isCursorWithinNodePoints(startPoint, endPoint ts.Point, pos lsp.Position) bool {
	if uint32(pos.Line) < startPoint.Row || uint32(pos.Line) > endPoint.Row {
		return false
	}

	if uint32(pos.Line) == startPoint.Row && uint32(pos.Character) < startPoint.Column {
		return false
	}

	if uint32(pos.Line) == endPoint.Row && uint32(pos.Character) > endPoint.Column {
		return false
	}

	return true
}

func isCursorAtLine(node *ts.Node, pos lsp.Position) bool {
	startPoint := node.StartPoint()
	endPoint := node.EndPoint()
	return uint32(pos.Line) == startPoint.Row && uint32(pos.Line) == endPoint.Row
}

func getPositionType(document *string, position lsp.Position) PositionType {
	if inMixinsPosition(document, position) {
		return PositionTypeMixins
	}
	return PositionTypeNone
}

// TODO: handle errors ?
func inMixinsPosition(document *string, position lsp.Position) bool {
	parser := ts.NewParser()
	parser.SetLanguage(yaml.GetLanguage())

	query, err := ts.NewQuery([]byte(`
		(block_mapping_pair
			key: (flow_node) @key
			value: (block_node
				(block_sequence
					(block_sequence_item
						(flow_node) @value)))
			(#eq? @key "mixins")
		)
	`), yaml.GetLanguage())
	if err != nil {
		return false
	}

	docBytes := []byte(*document)

	tree, err := parser.ParseCtx(context.Background(), nil, docBytes)
	if err != nil {
		return false
	}
	root := tree.RootNode()

	cursor := ts.NewQueryCursor()
	cursor.Exec(query, root)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if parent := capture.Node.Parent(); parent != nil {
				nodeText := capture.Node.Content(docBytes)
				if parent.Type() == "block_mapping_pair" &&
				   string(nodeText) == "mixins" &&
				   isCursorWithinNode(parent, position) {
					return true
				}
			}
		}
	}
	return false
}

func extractFilenameFromMixins(document *string, position lsp.Position) string {
	parser := ts.NewParser()
	parser.SetLanguage(yaml.GetLanguage())

	query, err := ts.NewQuery([]byte(`
		(block_mapping_pair
			key: (flow_node) @key
			value: (block_node
				(block_sequence
					(block_sequence_item
						(flow_node) @value)))
			(#eq? @key "mixins")
		)
	`), yaml.GetLanguage())
	if err != nil {
		return ""
	}

	docBytes := []byte(*document)

	tree, err := parser.ParseCtx(context.Background(), nil, docBytes)
	if err != nil {
		return ""
	}
	root := tree.RootNode()

	cursor := ts.NewQueryCursor()
	cursor.Exec(query, root)

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if parent := capture.Node.Parent(); parent != nil {
				if parent.Type() == "block_sequence_item" && isCursorAtLine(capture.Node, position) {
					return capture.Node.Content(docBytes)
				}
			}
		}
	}
	return ""
}
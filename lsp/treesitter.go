package lsp

import (
	// "fmt"

	tree_sitter_yaml "github.com/tree-sitter-grammars/tree-sitter-yaml/bindings/go"
	ts "github.com/tree-sitter/go-tree-sitter"

	lsp "github.com/tliron/glsp/protocol_3_16"
)

type PositionType int

const (
    PositionTypeMixins PositionType = iota
    PositionTypeNone
)

func isCursorWithinNode(node *ts.Node, pos lsp.Position) bool {
	return isCursorWithinNodePoints(node.StartPosition(), node.EndPosition(), pos)
}

func isCursorWithinNodePoints(startPoint, endPoint ts.Point, pos lsp.Position) bool {
	if uint(pos.Line) < startPoint.Row || uint(pos.Line) > endPoint.Row {
		return false
	}

	if uint(pos.Line) == startPoint.Row && uint(pos.Character) < startPoint.Column {
		return false
	}

	if uint(pos.Line) == endPoint.Row && uint(pos.Character) > endPoint.Column {
		return false
	}

	return true
}

func isCursorAtLine(node *ts.Node, pos lsp.Position) bool {
	startPoint := node.StartPosition()
	endPoint := node.EndPosition()
	return uint(pos.Line) == startPoint.Row && uint(pos.Line) == endPoint.Row
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
    defer parser.Close()
	lang := ts.NewLanguage(tree_sitter_yaml.Language())
    parser.SetLanguage(lang)

	docBytes := []byte(*document)

    tree := parser.Parse(docBytes, nil)
    defer tree.Close()

	query, err := ts.NewQuery(lang, `
		(block_mapping_pair
			key: (flow_node) @key
			value: (block_node
				(block_sequence
					(block_sequence_item
						(flow_node) @value)))
			(#eq? @key "mixins")
		)
	`)
	if err != nil {
		return false
	}

	defer query.Close()

    root := tree.RootNode()
    // fmt.Println(root.ToSexp())

	cursor := ts.NewQueryCursor()
	defer cursor.Close()

	matches := cursor.Matches(query, root, docBytes)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, capture := range match.Captures {
			if parent := capture.Node.Parent(); parent != nil {
				nodeText := capture.Node.Utf8Text(docBytes)
				if parent.Kind() == "block_mapping_pair" &&
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
    defer parser.Close()
	lang := ts.NewLanguage(tree_sitter_yaml.Language())
    parser.SetLanguage(lang)

	docBytes := []byte(*document)

    tree := parser.Parse(docBytes, nil)
    defer tree.Close()

	query, err := ts.NewQuery(lang, `
		(block_mapping_pair
			key: (flow_node) @key
			value: (block_node
				(block_sequence
					(block_sequence_item
						(flow_node) @value)))
			(#eq? @key "mixins")
		)
	`)
	if err != nil {
		return ""
	}
	defer query.Close()

    root := tree.RootNode()
    // fmt.Println(root.ToSexp())

	cursor := ts.NewQueryCursor()
	defer cursor.Close()
	matches := cursor.Matches(query, root, docBytes)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, capture := range match.Captures {
			if parent := capture.Node.Parent(); parent != nil {
				if parent.Kind() == "block_sequence_item" && isCursorAtLine(&capture.Node, position) {
					return capture.Node.Utf8Text(docBytes)
				}
			}
		}
	}

	return ""
}
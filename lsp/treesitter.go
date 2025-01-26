package lsp

import (
	tree_sitter_yaml "github.com/tree-sitter-grammars/tree-sitter-yaml/bindings/go"
	ts "github.com/tree-sitter/go-tree-sitter"

	"github.com/tliron/commonlog"
	lsp "github.com/tliron/glsp/protocol_3_16"
)

type PositionType int

const (
    PositionTypeMixins PositionType = iota
	PositionTypeDepends
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

type parser struct {
	log commonlog.Logger
}

func newParser(log commonlog.Logger) *parser {
	return &parser{
		log: log,
	}
}

func (p *parser) getPositionType(document *string, position lsp.Position) PositionType {
	if p.inMixinsPosition(document, position) {
		return PositionTypeMixins
	} else if p.inDependsPosition(document, position) {
		return PositionTypeDepends
	}
	return PositionTypeNone
}

func (p *parser) inMixinsPosition(document *string, position lsp.Position) bool {
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

func (p *parser) inDependsPosition(document *string, position lsp.Position) bool {
	parser := ts.NewParser()
	defer parser.Close()
	lang := ts.NewLanguage(tree_sitter_yaml.Language())
	parser.SetLanguage(lang)

	docBytes := []byte(*document)

	tree := parser.Parse(docBytes, nil)
	defer tree.Close()

	query, err := ts.NewQuery(lang, `
		(block_mapping_pair
			key: (flow_node) @keydepends
			value: [
				(flow_node(flow_sequence)) @depends
				(flow_node(flow_sequence(flow_node(plain_scalar(string_scalar))))) @depends
				(block_node(block_sequence(block_sequence_item) @depends))
			]
			(#eq? @keydepends "depends")
		)
	`)
	if err != nil {
		return false
	}
	defer query.Close()

	root := tree.RootNode()
	cursor := ts.NewQueryCursor()
	defer cursor.Close()
	matches := cursor.Matches(query, root, docBytes)

	dependsIndex, _ := query.CaptureIndexForName("depends")

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, capture := range match.Captures {
			nodeKind := capture.Node.Kind()

			if capture.Index != uint32(dependsIndex) {
				continue
			}

			// if is a sequence
			if nodeKind == "block_sequence_item" || nodeKind == "block_sequence" {
				if isCursorWithinNode(&capture.Node, position) || isCursorAtLine(&capture.Node, position) {
					return true
				}
			// if is an array
			} else if nodeKind == "flow_sequence" || nodeKind == "flow_node" {
				if isCursorWithinNode(&capture.Node, position) {
					return true
				}
			}
		}
	}

	return false
}

func (p *parser) extractFilenameFromMixins(document *string, position lsp.Position) string {
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

type Command struct {
	name string
}

func (p *parser) getCommands(document *string) []Command {
	parser := ts.NewParser()
	defer parser.Close()
	lang := ts.NewLanguage(tree_sitter_yaml.Language())
	parser.SetLanguage(lang)

	docBytes := []byte(*document)

	tree := parser.Parse(docBytes, nil)
	defer tree.Close()

	query, err := ts.NewQuery(lang, `
		(block_mapping_pair
			key: (flow_node(plain_scalar(string_scalar)) @parent)
			value: (block_node
				(block_mapping
					(block_mapping_pair
						key: (flow_node
							(plain_scalar
								(string_scalar)) @cmd_key)
						value: (block_node) @cmd) @values))
			(#eq? @parent "commands")
		)
	`)
	if err != nil {
		return nil
	}
	defer query.Close()

	root := tree.RootNode()
	cursor := ts.NewQueryCursor()
	defer cursor.Close()
	matches := cursor.Matches(query, root, docBytes)

	var commands []Command
	cmdKeyIndex, _ := query.CaptureIndexForName("cmd_key")

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, capture := range match.Captures {
			if capture.Index == uint32(cmdKeyIndex) {
				commands = append(commands, Command{
					name: capture.Node.Utf8Text(docBytes),
				})
			}
		}
	}

	return commands
}

func (p *parser) getCurrentCommand(document *string, position lsp.Position) *Command {
	parser := ts.NewParser()
	defer parser.Close()
	lang := ts.NewLanguage(tree_sitter_yaml.Language())
	parser.SetLanguage(lang)

	docBytes := []byte(*document)

	tree := parser.Parse(docBytes, nil)
	defer tree.Close()

	query, err := ts.NewQuery(lang, `
		(block_mapping_pair
			key: (flow_node(plain_scalar(string_scalar)) @commands)
			value: (block_node
				(block_mapping
					(block_mapping_pair) @cmd))
			(#eq? @commands "commands")
		)
	`)
	if err != nil {
		return nil
	}
	defer query.Close()

	root := tree.RootNode()
	cursor := ts.NewQueryCursor()
	defer cursor.Close()
	matches := cursor.Matches(query, root, docBytes)

	cmdIndex, _ := query.CaptureIndexForName("cmd")

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, capture := range match.Captures {
			if capture.Index != uint32(cmdIndex) {
				continue
			}
			if !isCursorWithinNode(&capture.Node, position) {
				continue
			}
			if key := capture.Node.ChildByFieldName("key"); key != nil {
				return &Command{
					name: key.Utf8Text(docBytes),
				}
			}
		}
	}

	return nil
}

func (p *parser) extractDependsValues(document *string) []string {
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
			value: (flow_node
				(flow_sequence
					(flow_node
						(plain_scalar
							(string_scalar) @value))))
			(#eq? @key "depends")
		)
	`)
	if err != nil {
		return nil
	}
	defer query.Close()

	root := tree.RootNode()
	cursor := ts.NewQueryCursor()
	defer cursor.Close()
	matches := cursor.Matches(query, root, docBytes)

	var values []string
	valueIndex, _ := query.CaptureIndexForName("value")

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, capture := range match.Captures {
			if capture.Index == uint32(valueIndex) {
				values = append(values, capture.Node.Utf8Text(docBytes))
			}
		}
	}

	return values
}
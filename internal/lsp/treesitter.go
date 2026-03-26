package lsp

import (
	"strings"

	ts "github.com/odvcencio/gotreesitter"
	"github.com/odvcencio/gotreesitter/grammars"
	"github.com/tliron/commonlog"
	lsp "github.com/tliron/glsp/protocol_3_16"
)

type PositionType int

const (
	PositionTypeMixins PositionType = iota
	PositionTypeDepends
	PositionTypeCommandAlias
	PositionTypeNone
)

var yamlLanguage = grammars.YamlLanguage()

func isCursorWithinNode(node *ts.Node, pos lsp.Position) bool {
	return isCursorWithinNodePoints(node.StartPoint(), node.EndPoint(), pos)
}

func isCursorWithinNodePoints(startPoint, endPoint ts.Point, pos lsp.Position) bool {
	if pos.Line < startPoint.Row || pos.Line > endPoint.Row {
		return false
	}

	if pos.Line == startPoint.Row && pos.Character < startPoint.Column {
		return false
	}

	if pos.Line == endPoint.Row && pos.Character > endPoint.Column {
		return false
	}

	return true
}

func isCursorAtLine(node *ts.Node, pos lsp.Position) bool {
	startPoint := node.StartPoint()
	endPoint := node.EndPoint()

	return pos.Line == startPoint.Row && pos.Line == endPoint.Row
}

func parseYAMLDocument(document *string) (*ts.Tree, []byte, error) {
	docBytes := []byte(*document)

	tree, err := ts.NewParser(yamlLanguage).Parse(docBytes)
	if err != nil {
		return nil, nil, err
	}

	return tree, docBytes, nil
}

func getLine(document *string, line uint32) string {
	lines := strings.Split(*document, "\n")
	if line >= uint32(len(lines)) {
		return ""
	}

	return lines[line]
}

// position.
func wordUnderCursor(text string, position *lsp.Position) string {
	if len(text) == 0 {
		return ""
	}

	character := position.Character

	if character >= uint32(len(text)) {
		return ""
	}

	if text[character] == ' ' {
		return ""
	}

	// Find word boundaries
	start := position.Character
	for start > 0 && isWordChar(text[start-1]) {
		start--
	}

	end := position.Character
	for end < uint32(len(text)) && isWordChar(text[end]) {
		end++
	}

	return text[start:end]
}

func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-'
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
	} else if p.inCommandAliasPosition(document, position) {
		return PositionTypeCommandAlias
	}

	return PositionTypeNone
}

func (p *parser) inMixinsPosition(document *string, position lsp.Position) bool {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return false
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node) @key
			value: (block_node
				(block_sequence
					(block_sequence_item
						(flow_node) @value)))
			(#eq? @key "mixins")
		)
	`, yamlLanguage)
	if err != nil {
		return false
	}

	root := tree.RootNode()
	if root == nil {
		return false
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if parent := capture.Node.Parent(); parent != nil {
				nodeText := capture.Node.Text(docBytes)
				if parent.Type(yamlLanguage) == "block_mapping_pair" &&
					nodeText == "mixins" &&
					isCursorWithinNode(parent, position) {
					return true
				}
			}
		}
	}

	return false
}

func (p *parser) inDependsPosition(document *string, position lsp.Position) bool {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return false
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node) @keydepends
			value: [
				(flow_node(flow_sequence)) @depends
				(flow_node(flow_sequence(flow_node(plain_scalar(string_scalar))))) @depends
				(block_node(block_sequence(block_sequence_item) @depends))
			]
			(#eq? @keydepends "depends")
		)
	`, yamlLanguage)
	if err != nil {
		return false
	}

	root := tree.RootNode()
	if root == nil {
		return false
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Name != "depends" {
				continue
			}

			nodeKind := capture.Node.Type(yamlLanguage)

			// if is a sequence
			switch nodeKind {
			case "block_sequence_item", "block_sequence":
				if isCursorWithinNode(capture.Node, position) || isCursorAtLine(capture.Node, position) {
					return true
				}
				// if is an array
			case "flow_sequence", "flow_node":
				if isCursorWithinNode(capture.Node, position) {
					return true
				}
			}
		}
	}

	return false
}

func (p *parser) inCommandAliasPosition(document *string, position lsp.Position) bool {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return false
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node) @keymerge
			value: (flow_node(alias) @alias)
			(#eq? @keymerge "<<")
		)
	`, yamlLanguage)
	if err != nil {
		return false
	}

	root := tree.RootNode()
	if root == nil {
		return false
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Name == "alias" && isCursorWithinNode(capture.Node, position) {
				return true
			}
		}
	}

	return false
}

func (p *parser) extractFilenameFromMixins(document *string, position lsp.Position) string {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return ""
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node) @key
			value: (block_node
				(block_sequence
					(block_sequence_item
						(flow_node) @value)))
			(#eq? @key "mixins")
		)
	`, yamlLanguage)
	if err != nil {
		return ""
	}

	root := tree.RootNode()
	if root == nil {
		return ""
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if parent := capture.Node.Parent(); parent != nil {
				if parent.Type(yamlLanguage) == "block_sequence_item" && isCursorAtLine(capture.Node, position) {
					return capture.Node.Text(docBytes)
				}
			}
		}
	}

	return ""
}

func (p *parser) extractCommandReference(document *string, position lsp.Position) string {
	if commandName := p.extractDependsCommandReference(document, position); commandName != "" {
		return commandName
	}

	return p.extractAliasCommandReference(document, position)
}

func (p *parser) extractDependsCommandReference(document *string, position lsp.Position) string {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return ""
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node) @keydepends
			value: [
				(flow_node
					(flow_sequence
						(flow_node
							(plain_scalar
								(string_scalar)) @reference)))
				(block_node
					(block_sequence
						(block_sequence_item
							(flow_node
								(plain_scalar
									(string_scalar)) @reference))))
			]
			(#eq? @keydepends "depends")
		)
	`, yamlLanguage)
	if err != nil {
		return ""
	}

	root := tree.RootNode()
	if root == nil {
		return ""
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Name == "reference" && isCursorWithinNode(capture.Node, position) {
				return capture.Node.Text(docBytes)
			}
		}
	}

	return ""
}

func (p *parser) extractAliasCommandReference(document *string, position lsp.Position) string {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return ""
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node) @keymerge
			value: (flow_node(alias) @reference)
			(#eq? @keymerge "<<")
		)
	`, yamlLanguage)
	if err != nil {
		return ""
	}

	root := tree.RootNode()
	if root == nil {
		return ""
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Name != "reference" || !isCursorWithinNode(capture.Node, position) {
				continue
			}

			return strings.TrimPrefix(capture.Node.Text(docBytes), "*")
		}
	}

	return ""
}

type Command struct {
	name string
	// TODO: maybe range will be more appropriate
	position lsp.Position
}

func (p *parser) getCommands(document *string) []Command {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return nil
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
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
	`, yamlLanguage)
	if err != nil {
		return nil
	}

	root := tree.RootNode()
	if root == nil {
		return nil
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	var commands []Command

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Name == "cmd_key" {
				commands = append(commands, Command{
					name: capture.Node.Text(docBytes),
					position: lsp.Position{
						Line:      capture.Node.StartPoint().Row,
						Character: capture.Node.StartPoint().Column,
					},
				})
			}
		}
	}

	return commands
}

func (p *parser) getCurrentCommand(document *string, position lsp.Position) *Command {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return nil
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node(plain_scalar(string_scalar)) @commands)
			value: (block_node
				(block_mapping
					(block_mapping_pair) @cmd))
			(#eq? @commands "commands")
		)
	`, yamlLanguage)
	if err != nil {
		return nil
	}

	root := tree.RootNode()
	if root == nil {
		return nil
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Name != "cmd" {
				continue
			}

			if !isCursorWithinNode(capture.Node, position) {
				continue
			}

			if key := capture.Node.ChildByFieldName("key", yamlLanguage); key != nil {
				return &Command{
					name: key.Text(docBytes),
				}
			}
		}
	}

	return nil
}

func (p *parser) findCommand(document *string, commandName string) *Command {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return nil
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node(plain_scalar(string_scalar)) @commands)
			value: (block_node
				(block_mapping
					(block_mapping_pair
						key: (flow_node
							(plain_scalar
								(string_scalar)) @cmd_key)
						value: (block_node) @cmd_value)) @values)
			(#eq? @commands "commands")
		)
	`, yamlLanguage)
	if err != nil {
		return nil
	}

	root := tree.RootNode()
	if root == nil {
		return nil
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Name != "cmd_key" {
				continue
			}

			if capture.Node.Text(docBytes) == commandName {
				return &Command{
					name: commandName,
					position: lsp.Position{
						Line:      capture.Node.StartPoint().Row,
						Character: capture.Node.StartPoint().Column,
					},
				}
			}
		}
	}

	return nil
}

func (p *parser) extractDependsValues(document *string) []string {
	tree, docBytes, err := parseYAMLDocument(document)
	if err != nil {
		return nil
	}
	defer tree.Release()

	query, err := ts.NewQuery(`
		(block_mapping_pair
			key: (flow_node) @key
			value: [
				(flow_node
					(flow_sequence
						(flow_node
							(plain_scalar
								(string_scalar) @value))))
				(block_node
					(block_sequence
						(block_sequence_item
							(flow_node
								(plain_scalar
									(string_scalar) @value)))))
			]
			(#eq? @key "depends")
		)
	`, yamlLanguage)
	if err != nil {
		return nil
	}

	root := tree.RootNode()
	if root == nil {
		return nil
	}

	matches := query.Exec(root, yamlLanguage, docBytes)

	var values []string

	for {
		match, ok := matches.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Name == "value" {
				values = append(values, capture.Node.Text(docBytes))
			}
		}
	}

	return values
}

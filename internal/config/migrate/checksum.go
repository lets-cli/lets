package migrate

import "gopkg.in/yaml.v3"

type ChecksumMigration struct{}

func (ChecksumMigration) Name() string {
	return "checksum"
}

func (ChecksumMigration) Apply(root *yaml.Node) (bool, error) {
	commands := mappingValue(document(root), "commands")
	if commands == nil || commands.Kind != yaml.MappingNode {
		return false, nil
	}

	changed := false

	for i := 0; i < len(commands.Content); i += 2 {
		command := commands.Content[i+1]
		if command.Kind != yaml.MappingNode {
			continue
		}

		commandChanged, err := migrateCommandChecksum(command)
		if err != nil {
			return false, err
		}

		changed = changed || commandChanged
	}

	return changed, nil
}

func migrateCommandChecksum(command *yaml.Node) (bool, error) {
	checksumIdx := mappingIndex(command, "checksum")
	persistIdx := mappingIndex(command, "persist_checksum")

	if checksumIdx == -1 {
		return false, nil
	}

	checksumNode := command.Content[checksumIdx+1]
	changed := false

	if isNewChecksumNode(checksumNode) {
		if persistIdx != -1 && mappingIndex(checksumNode, "persist") == -1 {
			appendMapping(checksumNode, scalar("persist"), cloneNode(command.Content[persistIdx+1]))

			changed = true
		}
	} else if checksumNode.Kind == yaml.SequenceNode || checksumNode.Kind == yaml.MappingNode {
		filesNode := cloneNode(checksumNode)
		checksumNode.Kind = yaml.MappingNode
		checksumNode.Tag = "!!map"
		checksumNode.Content = []*yaml.Node{scalar("files"), filesNode}

		if persistIdx != -1 {
			appendMapping(checksumNode, scalar("persist"), cloneNode(command.Content[persistIdx+1]))
		}

		changed = true
	}

	if persistIdx != -1 {
		removeMappingIndex(command, persistIdx)

		changed = true
	}

	return changed, nil
}

func isNewChecksumNode(node *yaml.Node) bool {
	if node == nil || node.Kind != yaml.MappingNode {
		return false
	}

	return mappingIndex(node, "files") != -1 || mappingIndex(node, "sh") != -1 || mappingIndex(node, "persist") != -1
}

func mappingIndex(node *yaml.Node, key string) int {
	if node == nil || node.Kind != yaml.MappingNode {
		return -1
	}

	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return i
		}
	}

	return -1
}

func appendMapping(node *yaml.Node, key *yaml.Node, value *yaml.Node) {
	node.Content = append(node.Content, key, value)
}

func removeMappingIndex(node *yaml.Node, idx int) {
	node.Content = append(node.Content[:idx], node.Content[idx+2:]...)
}

func scalar(value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}
}

func cloneNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}

	clone := *node

	clone.Content = make([]*yaml.Node, len(node.Content))
	for idx, child := range node.Content {
		clone.Content[idx] = cloneNode(child)
	}

	return &clone
}

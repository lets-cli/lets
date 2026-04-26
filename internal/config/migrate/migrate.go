package migrate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/lets-cli/lets/internal/config"
	configpath "github.com/lets-cli/lets/internal/config/path"
	"gopkg.in/yaml.v3"
)

type Migration interface {
	Name() string
	Apply(root *yaml.Node) (bool, error)
}

type Result struct {
	ChangedFiles []string
	RemoteMixins []string
	Applied      []string
	Changed      bool
	DryRun       bool
	Previews     []string
}

func DefaultMigrations() []Migration {
	return []Migration{
		ChecksumMigration{},
	}
}

func Fix(configName string, configDir string, dryRun bool, out io.Writer) (Result, error) {
	pathInfo, err := config.FindConfig(configName, configDir)
	if err != nil {
		return Result{}, err
	}

	paths, remoteMixins, err := collectConfigPaths(pathInfo.AbsPath, pathInfo.WorkDir)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		DryRun:       dryRun,
		RemoteMixins: remoteMixins,
	}

	for _, path := range paths {
		fileChanged, applied, preview, err := fixFile(path, dryRun, DefaultMigrations())
		if err != nil {
			return Result{}, err
		}

		if !fileChanged {
			continue
		}

		result.Changed = true
		result.ChangedFiles = append(result.ChangedFiles, path)

		result.Applied = append(result.Applied, applied...)
		if preview != "" {
			result.Previews = append(result.Previews, preview)
		}
	}

	writeResult(out, result)

	return result, nil
}

func fixFile(path string, dryRun bool, migrations []Migration) (bool, []string, string, error) {
	original, err := os.ReadFile(path)
	if err != nil {
		return false, nil, "", fmt.Errorf("can not read config %s: %w", path, err)
	}

	root, err := decodeYAML(original)
	if err != nil {
		return false, nil, "", fmt.Errorf("can not parse config %s: %w", path, err)
	}

	applied := []string{}

	for _, migration := range migrations {
		changed, err := migration.Apply(root)
		if err != nil {
			return false, nil, "", fmt.Errorf("can not apply migration %s to %s: %w", migration.Name(), path, err)
		}

		if changed {
			applied = append(applied, migration.Name())
		}
	}

	if len(applied) == 0 {
		return false, nil, "", nil
	}

	updated, err := encodeYAML(root)
	if err != nil {
		return false, nil, "", fmt.Errorf("can not render config %s: %w", path, err)
	}

	if bytes.Equal(original, updated) {
		return false, nil, "", nil
	}

	preview := ""

	if !dryRun {
		if err := os.WriteFile(path, updated, 0o644); err != nil {
			return false, nil, "", fmt.Errorf("can not write config %s: %w", path, err)
		}
	} else {
		preview = string(updated)
	}

	return true, applied, preview, nil
}

func decodeYAML(data []byte) (*yaml.Node, error) {
	root := &yaml.Node{}
	if err := yaml.Unmarshal(data, root); err != nil {
		return nil, err
	}

	return root, nil
}

func encodeYAML(root *yaml.Node) ([]byte, error) {
	var buf bytes.Buffer

	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(root); err != nil {
		return nil, err
	}

	if err := encoder.Close(); err != nil {
		return nil, err
	}

	return formatCommandSpacing(buf.Bytes()), nil
}

func formatCommandSpacing(data []byte) []byte {
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(lines) == 0 {
		return data
	}

	formatted := make([]string, 0, len(lines))
	inCommands := false
	commandSeen := false

	for _, line := range lines {
		if line == "commands:" {
			inCommands = true
			commandSeen = false
			formatted = append(formatted, line)

			continue
		}

		if inCommands && line != "" && !strings.HasPrefix(line, " ") {
			inCommands = false
		}

		if inCommands && isCommandEntryLine(line) {
			if commandSeen && len(formatted) > 0 && formatted[len(formatted)-1] != "" {
				formatted = append(formatted, "")
			}

			commandSeen = true
		}

		formatted = append(formatted, line)
	}

	return []byte(strings.Join(formatted, "\n") + "\n")
}

func isCommandEntryLine(line string) bool {
	if !strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "    ") {
		return false
	}

	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "-") {
		return false
	}

	return strings.Contains(trimmed, ":")
}

func collectConfigPaths(rootPath string, workDir string) ([]string, []string, error) {
	paths := []string{rootPath}
	seen := map[string]struct{}{rootPath: {}}

	data, err := os.ReadFile(rootPath)
	if err != nil {
		return nil, nil, err
	}

	root, err := decodeYAML(data)
	if err != nil {
		return nil, nil, err
	}

	mixins := mappingValue(document(root), "mixins")
	if mixins == nil || mixins.Kind != yaml.SequenceNode {
		return paths, nil, nil
	}

	remoteMixins := []string{}

	for _, mixin := range mixins.Content {
		if mixin.Kind == yaml.ScalarNode {
			ignored := strings.HasPrefix(mixin.Value, "-")
			mixinPath := strings.TrimPrefix(mixin.Value, "-")

			absPath, err := configpath.GetFullConfigPath(mixinPath, workDir)
			if err != nil {
				if ignored {
					continue
				}

				return nil, nil, err
			}

			if _, ok := seen[absPath]; !ok {
				seen[absPath] = struct{}{}
				paths = append(paths, absPath)
			}

			continue
		}

		if mixin.Kind != yaml.MappingNode {
			continue
		}

		if url := mappingValue(mixin, "url"); url != nil && url.Value != "" {
			remoteMixins = append(remoteMixins, url.Value)
		}
	}

	return paths, remoteMixins, nil
}

func document(root *yaml.Node) *yaml.Node {
	if root.Kind == yaml.DocumentNode && len(root.Content) > 0 {
		return root.Content[0]
	}

	return root
}

func mappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}

	return nil
}

func writeResult(out io.Writer, result Result) {
	for _, preview := range result.Previews {
		fmt.Fprint(out, preview)
	}

	if result.Changed && !result.DryRun {
		applied := slices.Compact(slices.Sorted(slices.Values(result.Applied)))
		for _, migration := range applied {
			fmt.Fprintf(out, "Migration '%s' applied successfully\n", migration)
		}
	}

	for _, remote := range result.RemoteMixins {
		fmt.Fprintf(out, "remote mixin not updated: %s\n", remote)
	}

	if !result.Changed && len(result.RemoteMixins) == 0 {
		fmt.Fprintln(out, "No config migrations needed.")
	}
}

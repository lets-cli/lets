package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	skillpkg "github.com/lets-cli/lets/internal/skills"
	"github.com/spf13/cobra"
)

func initSkillsCommand() *cobra.Command {
	skillsCmd := &cobra.Command{
		Use:   "skills <command>",
		Short: "Manage lets agent skills. (EXPERIMENTAL)",
		Long: strings.TrimSpace(`Install the bundled lets agent skill so that AI agents can discover
and use lets effectively.

Skills follow the Agent Skills specification and work with any compatible agent,
including Claude Code, Codex, and Gemini CLI, PI, Open Code, etc.

This feature is an experiment and is not ready for production use. It might be
unstable or removed at any time.`),
		Args: validateCommandArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	skillsCmd.AddCommand(initSkillsShowCommand())
	skillsCmd.AddCommand(initSkillsInstallCommand())
	skillsCmd.AddCommand(initSkillsUpdateCommand())

	return skillsCmd
}

func initSkillsShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show lets bundled agent skill. (EXPERIMENTAL)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cmd.OutOrStdout().Write(skillpkg.LetsSkill())
			return err
		},
	}
}

type skillsInstallOptions struct {
	global bool
	local  bool
	path   string
	force  bool
}

func initSkillsInstallCommand() *cobra.Command {
	opts := &skillsInstallOptions{}

	installCmd := &cobra.Command{
		Use:   "install [name]",
		Short: "Install lets' bundled agent skill. (EXPERIMENTAL)",
		Long: strings.TrimSpace(`Install the bundled lets SKILL.md file to .agents/skills/, the
cross-agent standard defined by the Agent Skills specification.

By default, install prompts for local project scope or global user scope. Local
scope installs to .agents/skills/ at the current Git repository root. Global
scope installs to ~/.agents/skills/.`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateSkillNameArg(args); err != nil {
				return err
			}

			targetDir, err := opts.targetDir(cmd)
			if err != nil {
				return err
			}

			return installLetsSkill(cmd.OutOrStdout(), targetDir, opts.force)
		},
	}

	installCmd.Flags().BoolVarP(&opts.global, "global", "g", false, "Install the skill at user scope (~/.agents/skills/).")
	installCmd.Flags().BoolVarP(&opts.local, "local", "l", false, "Install the skill at local project scope (.agents/skills/).")
	installCmd.Flags().StringVar(&opts.path, "path", "", "Install the skill to the directory at <path>.")
	installCmd.Flags().BoolVarP(&opts.force, "force", "f", false, "Overwrite an existing skill file.")
	installCmd.MarkFlagsMutuallyExclusive("global", "local", "path")

	return installCmd
}

func (o skillsInstallOptions) targetDir(cmd *cobra.Command) (string, error) {
	if o.path != "" {
		return o.path, nil
	}

	if o.global {
		return globalSkillsDir()
	}

	if o.local {
		return localSkillsDir()
	}

	localDir, err := localSkillsDir()
	if err != nil {
		return "", err
	}

	globalDir, err := globalSkillsDir()
	if err != nil {
		return "", err
	}

	scope, err := promptSkillScope(cmd.InOrStdin(), cmd.ErrOrStderr(), localDir, globalDir)
	if err != nil {
		return "", err
	}

	if scope == "global" {
		return globalDir, nil
	}

	return localDir, nil
}

func promptSkillScope(in io.Reader, out io.Writer, localDir string, globalDir string) (string, error) {
	_, _ = fmt.Fprintf(out, "Install lets skill:\n")
	_, _ = fmt.Fprintf(out, "  1. Local  %s\n", filepath.Join(localDir, skillpkg.LetsName))
	_, _ = fmt.Fprintf(out, "  2. Global %s\n", filepath.Join(globalDir, skillpkg.LetsName))
	_, _ = fmt.Fprint(out, "Select [1/2]: ")

	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && (!errors.Is(err, io.EOF) || line == "") {
		return "", fmt.Errorf("reading install scope: %w", err)
	}

	switch strings.ToLower(strings.TrimSpace(line)) {
	case "l", "local", "1":
		return "local", nil
	case "g", "global", "2":
		return "global", nil
	default:
		return "", fmt.Errorf("invalid install scope %q; enter local or global", strings.TrimSpace(line))
	}
}

func initSkillsUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update [name]",
		Short: "Update installed agent skills to the current shipped version. (EXPERIMENTAL)",
		Long: strings.TrimSpace(`Update installed lets agent skills in known locations to the current
version bundled in this lets binary.

Known locations are the current project's .agents/skills/ directory and the
user-scope ~/.agents/skills/ directory.`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateSkillNameArg(args); err != nil {
				return err
			}

			return updateLetsSkill(cmd.OutOrStdout())
		},
	}
}

func validateSkillNameArg(args []string) error {
	if len(args) == 0 || args[0] == skillpkg.LetsName {
		return nil
	}

	return fmt.Errorf("unknown bundled skill %q", args[0])
}

func installLetsSkill(out io.Writer, targetDir string, force bool) error {
	skillDir := filepath.Join(targetDir, skillpkg.LetsName)

	skillPath := filepath.Join(skillDir, skillpkg.SkillFile)
	if _, err := os.Stat(skillPath); err == nil && !force {
		_, _ = fmt.Fprintf(out, "%s already exists. Use --force to overwrite.\n", skillPath)
		return nil
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("checking %s: %w", skillPath, err)
	}

	if err := writeLetsSkill(skillDir); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(out, "Installed %s\n", skillDir)

	return nil
}

func updateLetsSkill(out io.Writer) error {
	dirs, err := knownLetsSkillDirs()
	if err != nil {
		return err
	}

	updated := 0

	for _, skillDir := range dirs {
		skillPath := filepath.Join(skillDir, skillpkg.SkillFile)

		current, err := os.ReadFile(skillPath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		if err != nil {
			return fmt.Errorf("reading %s: %w", skillPath, err)
		}

		if bytes.Equal(current, skillpkg.LetsSkill()) {
			_, _ = fmt.Fprintf(out, "%s already up to date.\n", skillDir)
			updated++

			continue
		}

		if err := writeLetsSkill(skillDir); err != nil {
			return err
		}

		_, _ = fmt.Fprintf(out, "Updated %s\n", skillDir)
		updated++
	}

	if updated == 0 {
		return errors.New("lets skill is not installed in any known location. Run 'lets self skills install' first")
	}

	return nil
}

func writeLetsSkill(skillDir string) error {
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", skillDir, err)
	}

	skillPath := filepath.Join(skillDir, skillpkg.SkillFile)
	if err := os.WriteFile(skillPath, skillpkg.LetsSkill(), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", skillPath, err)
	}

	return nil
}

func knownLetsSkillDirs() ([]string, error) {
	dirs := []string{}
	if localDir, err := localSkillsDir(); err == nil {
		dirs = append(dirs, filepath.Join(localDir, skillpkg.LetsName))
	}

	globalDir, err := globalSkillsDir()
	if err != nil {
		return nil, err
	}

	dirs = append(dirs, filepath.Join(globalDir, skillpkg.LetsName))

	return dirs, nil
}

func globalSkillsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	return filepath.Join(home, skillpkg.SkillsRelDir), nil
}

func localSkillsDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not determine current directory: %w", err)
	}

	root, err := findGitRoot(wd)
	if err != nil {
		return "", err
	}

	return filepath.Join(root, skillpkg.SkillsRelDir), nil
}

func findGitRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("resolving current directory: %w", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("checking Git root: %w", err)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("not in a Git repository. Use --global or --path to specify a target")
		}

		dir = parent
	}
}

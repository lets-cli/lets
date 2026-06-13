package cmd

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/ansi"
	"github.com/lets-cli/fang"
	"github.com/lets-cli/lets/internal/docopt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// helpItem is a pre-rendered key/description row in a help section.
type helpItem struct {
	key  string
	help string
}

// commandHelpItem carries command metadata needed for grouping and row rendering.
type commandHelpItem struct {
	name     string
	subgroup string
	key      string
	help     string
}

var (
	optionalArgsRe = regexp.MustCompile(`(\[.*\])`)
	paddedLeft2    = lipgloss.NewStyle().PaddingLeft(2)
)

// HelpRenderer renders Cobra help using lets' Fang-based layout.
//
// It coordinates the full output, for example:
//
//	USAGE
//	  lets build [<bin>] [--flags]
//
//	COMMANDS
//	  test        Run all tests
//
//	OPTIONS
//	  --help, -h  Show help
func HelpRenderer(cmd *cobra.Command, ctx fang.HelpContext) {
	newHelpRenderer(cmd, ctx).Render()
}

// helpRenderer owns the full help rendering flow for one command.
type helpRenderer struct {
	cmd *cobra.Command
	ctx fang.HelpContext
	out helpOutput
}

// newHelpRenderer wires a command, Fang context, and output helper together.
func newHelpRenderer(cmd *cobra.Command, ctx fang.HelpContext) helpRenderer {
	return helpRenderer{cmd: cmd, ctx: ctx, out: newHelpOutput(ctx)}
}

// helpOutput hides low-level writes and common help-specific styled lines.
type helpOutput struct {
	w      io.Writer
	styles fang.Styles
}

// newHelpOutput creates the help output adapter from Fang's help context.
func newHelpOutput(ctx fang.HelpContext) helpOutput {
	return helpOutput{w: ctx.Writer, styles: ctx.Styles}
}

// println writes one help output line and intentionally ignores write errors.
func (o helpOutput) println(v ...any) {
	_, _ = fmt.Fprintln(o.w, v...)
}

// blank writes one blank help output line.
func (o helpOutput) blank() {
	o.println()
}

// sectionTitle writes a top-level help section title.
func (o helpOutput) sectionTitle(title string) {
	o.println(compactTitleStyle(o.styles).Render(title))
}

// subgroupTitle writes a nested command subgroup title.
func (o helpOutput) subgroupTitle(title string) {
	o.println(subgroupTitleStyle(o.styles).Render(title))
}

// Render prints the complete help view for the renderer's command.
//
// It owns the section order: description, usage/examples, command groups, options.
func (r helpRenderer) Render() {
	r.renderLongShort(cmp.Or(r.cmd.Long, r.cmd.Short))
	r.renderUsageAndExamples()

	content := r.collectContent()
	r.renderCommandGroups(content)
	r.renderOptions(content)

	r.out.blank()
}

// helpContent is the collected command/options model used by the render phase.
type helpContent struct {
	groups       map[string]string
	groupKeys    []string
	commands     map[string][]commandHelpItem
	options      []helpItem
	optionsTitle string
	hasSubgroups bool
	space        int
}

// renderLongShort prints the command description before structured sections.
//
// Example output:
//
//	Run build tasks defined in lets.yaml.
func (r helpRenderer) renderLongShort(longShort string) {
	if longShort == "" {
		return
	}

	longShort = strings.TrimRight(longShort, "\n")
	r.out.blank()
	r.out.println(r.ctx.Styles.Text.Width(r.ctx.Width).Render(longShort))
}

// renderUsageAndExamples prints usage and example code blocks.
//
// Example output:
//
//	USAGE
//	  lets release <version> [--flags]
//
//	EXAMPLES
//	  lets release 1.2.3 --message "Release"
func (r helpRenderer) renderUsageAndExamples() {
	usage := styleHelpUsage(r.cmd, r.ctx.Styles.Codeblock.Program, true)
	examples := fang.StyleExamples(r.cmd, r.ctx.Styles)
	blockStyle := compactCodeBlockStyle(r.ctx, append([]string{usage}, examples...)...)

	r.out.sectionTitle("usage")
	r.out.blank()
	r.out.println(blockStyle.Render(usage))

	if len(examples) > 0 {
		cw := blockStyle.GetWidth() - blockStyle.GetHorizontalPadding()
		r.out.sectionTitle("examples")
		for i, example := range examples {
			if lipgloss.Width(example) > cw {
				examples[i] = ansi.Truncate(example, cw, "…")
			}
		}
		r.out.println(blockStyle.Render(strings.Join(examples, "\n")))
	}
}

// collectContent builds the command and option model before rendering rows.
func (r helpRenderer) collectContent() helpContent {
	groups, groupKeys := r.groups()
	commands := r.commandGroups()
	options, optionsTitle := r.optionItems()
	hasSubgroups := hasMultipleSubgroups(commands)
	space := r.helpSpace(commands, options, hasSubgroups)

	return helpContent{
		groups:       groups,
		groupKeys:    groupKeys,
		commands:     commands,
		options:      options,
		optionsTitle: optionsTitle,
		hasSubgroups: hasSubgroups,
		space:        space,
	}
}

// renderCommandGroups prints all non-empty command groups in Cobra group order.
//
// Example output:
//
//	COMMANDS
//	  build       Build lets
//	  test        Run all tests
//
//	INTERNAL COMMANDS
//	  completion  Generate completion scripts
func (r helpRenderer) renderCommandGroups(content helpContent) {
	for _, groupID := range content.groupKeys {
		items := content.commands[groupID]
		if len(items) == 0 {
			continue
		}
		r.renderCommandGroup(content.space, content.groups[groupID], items, content.hasSubgroups)
	}
}

// renderOptions prints flags or docopt options when the command exposes any.
//
// Example output:
//
//	OPTIONS
//	  <version>                Set version
//	  --message=<message>, -m  Release message
func (r helpRenderer) renderOptions(content helpContent) {
	if len(content.options) > 0 {
		r.renderHelpGroup(content.space, content.optionsTitle, content.options)
	}
}

// compactTitleStyle removes title spacing that Fang's defaults add for broader layouts.
func compactTitleStyle(styles fang.Styles) lipgloss.Style {
	return styles.Title.Margin(0, 0).PaddingBottom(0)
}

// subgroupTitleStyle derives the nested subgroup title from the compact title style.
func subgroupTitleStyle(styles fang.Styles) lipgloss.Style {
	return compactTitleStyle(styles).PaddingTop(0).PaddingLeft(2)
}

// compactCodeBlockStyle sizes code blocks to their content while respecting terminal width.
func compactCodeBlockStyle(ctx fang.HelpContext, blocks ...string) lipgloss.Style {
	base := ctx.Styles.Codeblock.Base.Padding(1, 2)
	padding := base.GetHorizontalPadding()
	blockWidth := 0
	for _, block := range blocks {
		blockWidth = max(blockWidth, lipgloss.Width(block))
	}
	blockWidth = min(ctx.Width-padding, blockWidth+padding)
	blockStyle := base.Width(blockWidth)

	if ctx.Writer.Profile <= colorprofile.Ascii || reflect.DeepEqual(blockStyle.GetBackground(), lipgloss.NoColor{}) {
		blockStyle = blockStyle.PaddingTop(0).PaddingBottom(0)
	}

	return blockStyle
}

// styleHelpUsage renders custom docopt usage annotations or falls back to Fang's usage renderer.
func styleHelpUsage(cmd *cobra.Command, styles fang.Program, complete bool) string {
	usage := cmd.Annotations[annotationHelpUsage]
	if usage == "" {
		return fang.StyleUsage(cmd, styles, complete)
	}

	var lines []string
	for _, line := range strings.Split(usage, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if complete {
			line = completeHelpUsage(cmd, line)
		}
		lines = append(lines, styleUsageText(cmd, styles, line, complete))
	}

	if len(lines) == 0 {
		return fang.StyleUsage(cmd, styles, complete)
	}

	return strings.Join(lines, "\n")
}

// completeHelpUsage prefixes custom usage lines with the parent command path when needed.
func completeHelpUsage(cmd *cobra.Command, usage string) string {
	parent := cmd.Parent()
	if parent == nil {
		return usage
	}

	parentPath := parent.CommandPath()
	if parentPath == "" || strings.HasPrefix(usage, parentPath+" ") {
		return usage
	}

	return parentPath + " " + usage
}

// styleUsageText applies Fang program styles to one normalized usage line.
func styleUsageText(cmd *cobra.Command, styles fang.Program, usage string, complete bool) string {
	hasArgs := strings.Contains(usage, "[args]")
	hasFlags := strings.Contains(usage, "[flags]") ||
		strings.Contains(usage, "[--flags]") ||
		cmd.HasFlags() ||
		cmd.HasPersistentFlags() ||
		cmd.HasAvailableFlags()
	hasCommands := strings.Contains(usage, "[command]") || cmd.HasAvailableSubCommands()
	for _, marker := range []string{
		"[args]",
		"[flags]", "[--flags]",
		"[command]",
	} {
		usage = strings.ReplaceAll(usage, marker, "")
	}

	var optionalArgs []string //nolint:prealloc
	for _, arg := range optionalArgsRe.FindAllString(usage, -1) {
		usage = strings.ReplaceAll(usage, arg, "")
		optionalArgs = append(optionalArgs, arg)
	}

	usage = strings.TrimSpace(usage)

	var useLine []string
	if complete {
		parts := strings.Fields(usage)
		if len(parts) > 0 {
			useLine = append(useLine, styles.Name.Render(parts[0]))
		}
		if len(parts) > 1 {
			useLine = append(useLine, styles.Command.Render(" "+strings.Join(parts[1:], " ")))
		}
	} else {
		useLine = append(useLine, styles.Command.Render(usage))
	}
	if hasCommands {
		useLine = append(useLine, styles.DimmedArgument.Render(" [command]"))
	}
	if hasArgs {
		useLine = append(useLine, styles.DimmedArgument.Render(" [args]"))
	}
	for _, arg := range optionalArgs {
		useLine = append(useLine, styles.DimmedArgument.Render(" "+arg))
	}
	if hasFlags {
		useLine = append(useLine, styles.DimmedArgument.Render(" [--flags]"))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, useLine...)
}

// groups returns Cobra command group titles in their render order.
func (r helpRenderer) groups() (map[string]string, []string) {
	ids := []string{""}
	groups := map[string]string{"": "commands"}

	for _, group := range r.cmd.Groups() {
		ids = append(ids, group.ID)
		groups[group.ID] = group.Title
	}

	return groups, ids
}

// commandGroups collects available subcommands grouped by Cobra group ID.
func (r helpRenderer) commandGroups() map[string][]commandHelpItem {
	commands := map[string][]commandHelpItem{}

	for _, subCmd := range r.cmd.Commands() {
		if !subCmd.IsAvailableCommand() && subCmd.Name() != "help" {
			continue
		}

		commands[subCmd.GroupID] = append(commands[subCmd.GroupID], commandHelpItem{
			name:     subCmd.Name(),
			subgroup: subCmd.Annotations[annotationSubGroupName],
			key:      r.ctx.Styles.Program.Command.Render(subCmd.Name()),
			help:     renderHelpDescription(r.ctx.Styles, subCmd.Short),
		})
	}

	for groupID := range commands {
		slices.SortFunc(commands[groupID], func(a, b commandHelpItem) int {
			return strings.Compare(a.name, b.name)
		})
	}

	return commands
}

// optionItems merges docopt options and Cobra flags into one rendered option list.
func (r helpRenderer) optionItems() ([]helpItem, string) {
	items := make([]helpItem, 0)
	docoptOptions := commandHelpOptions(r.cmd)

	for _, option := range docoptOptions {
		items = append(items, helpItem{
			key:  renderDocoptFlag(r.ctx.Styles.Program, option.Display),
			help: renderHelpDescription(r.ctx.Styles, option.Description),
		})
	}

	r.cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden || shouldSkipHelpFlag(r.cmd, flag) {
			return
		}

		help := renderHelpDescription(r.ctx.Styles, flag.Usage)
		if flag.DefValue != "" && flag.DefValue != "false" && flag.DefValue != "0" && flag.DefValue != "[]" {
			help += r.ctx.Styles.FlagDefault.Render(" (" + flag.DefValue + ")")
		}

		items = append(items, helpItem{
			key:  renderCobraFlag(r.ctx.Styles.Program, flag),
			help: help,
		})
	})

	if len(docoptOptions) > 0 {
		return items, "options"
	}

	return items, "flags"
}

// commandHelpOptions decodes docopt help options stored on the Cobra command.
func commandHelpOptions(cmd *cobra.Command) []docopt.HelpOption {
	payload := cmd.Annotations[annotationHelpOptions]
	if payload == "" {
		return nil
	}

	var options []docopt.HelpOption
	if err := json.Unmarshal([]byte(payload), &options); err != nil {
		return nil
	}

	return options
}

// shouldSkipHelpFlag hides inherited help flags on subcommands unless explicitly changed.
func shouldSkipHelpFlag(cmd *cobra.Command, flag *pflag.Flag) bool {
	return flag.Name == "help" && cmd != cmd.Root() && !flag.Changed
}

// renderCobraFlag renders one Cobra flag key with lets' program styles.
func renderCobraFlag(styles fang.Program, flag *pflag.Flag) string {
	if flag.Shorthand == "" {
		return styles.Flag.Render("--" + flag.Name)
	}

	return styles.Flag.Render("-" + flag.Shorthand + " --" + flag.Name)
}

// renderDocoptFlag renders a docopt option display string, including aliases.
func renderDocoptFlag(styles fang.Program, display string) string {
	parts := strings.Split(display, ", ")
	rendered := make([]string, 0, len(parts))

	for _, part := range parts {
		rendered = append(rendered, renderDocoptFlagPart(styles, part))
	}

	return strings.Join(rendered, styles.DimmedArgument.Render(", "))
}

// renderDocoptFlagPart styles one docopt option or argument fragment.
func renderDocoptFlagPart(styles fang.Program, part string) string {
	if left, right, ok := strings.Cut(part, "="); ok {
		return styles.Flag.Render(left+"=") + styles.Flag.Render(right)
	}

	return styles.Flag.Render(part)
}

// renderHelpDescription styles single and multi-line help descriptions.
func renderHelpDescription(styles fang.Styles, usage string) string {
	if !strings.Contains(usage, "\n") {
		return styles.FlagDescription.Render(usage)
	}

	noTransform := styles.FlagDescription.UnsetTransform()
	parts := strings.Split(usage, "\n")
	lines := make([]string, 0, len(parts))

	for i, line := range parts {
		if line == "" {
			lines = append(lines, "")
			continue
		}
		if i == 0 {
			lines = append(lines, styles.FlagDescription.Render(line))
			continue
		}
		lines = append(lines, noTransform.Render(line))
	}

	return strings.Join(lines, "\n")
}

// renderHelpGroup prints a generic key/description help section.
//
// Example output:
//
//	FLAGS
//	  --debug  Enable debug logging
func (r helpRenderer) renderHelpGroup(space int, title string, items []helpItem) {
	r.out.sectionTitle(title)
	for _, item := range items {
		r.renderHelpItem(space, item.key, item.help)
	}
}

// renderCommandGroup prints one command section, including subgroup headings when useful.
//
// Example output:
//
//	COMMANDS
//	  release  Create a release
//
//	  CI
//	    lint   Run lint checks
func (r helpRenderer) renderCommandGroup(space int, title string, items []commandHelpItem, hasSubgroups bool) {
	r.out.sectionTitle(title)

	names := subgroupNames(items)
	showSubgroupTitles := len(names) > 1

	bySubgroup := make(map[string][]commandHelpItem, len(names))
	var ungrouped []commandHelpItem
	for _, item := range items {
		if item.subgroup == "" {
			ungrouped = append(ungrouped, item)
		} else {
			bySubgroup[item.subgroup] = append(bySubgroup[item.subgroup], item)
		}
	}

	for _, subgroup := range names {
		if showSubgroupTitles {
			r.out.subgroupTitle(subgroup)
		}
		for _, item := range bySubgroup[subgroup] {
			r.renderHelpItem(space, displayCommandKey(item, hasSubgroups), item.help)
		}
	}

	for _, item := range ungrouped {
		r.renderHelpItem(space, displayCommandKey(item, hasSubgroups), item.help)
	}
}

// renderHelpItem prints one aligned key/description row.
//
// Example output:
//
//	--config  Path to lets config
func (r helpRenderer) renderHelpItem(space int, key string, help string) {
	r.out.println(lipgloss.JoinHorizontal(
		lipgloss.Left,
		paddedLeft2.Render(key),
		strings.Repeat(" ", max(space-lipgloss.Width(key), 0)),
		help,
	))
}

// subgroupNames returns sorted unique subgroup names present in command items.
func subgroupNames(items []commandHelpItem) []string {
	seen := map[string]struct{}{}
	var names []string

	for _, item := range items {
		if item.subgroup == "" {
			continue
		}
		if _, ok := seen[item.subgroup]; ok {
			continue
		}
		seen[item.subgroup] = struct{}{}
		names = append(names, item.subgroup)
	}

	slices.Sort(names)

	return names
}

// hasMultipleSubgroups reports whether command rendering needs subgroup-aware alignment.
func hasMultipleSubgroups(commands map[string][]commandHelpItem) bool {
	seen := map[string]struct{}{}

	for _, items := range commands {
		for _, item := range items {
			if item.subgroup == "" {
				continue
			}
			seen[item.subgroup] = struct{}{}
			if len(seen) > 1 {
				return true
			}
		}
	}

	return false
}

// displayCommandKey adjusts command key spacing for mixed grouped and ungrouped commands.
func displayCommandKey(item commandHelpItem, hasSubgroups bool) string {
	if !hasSubgroups {
		return item.key
	}
	if item.subgroup != "" {
		return "  " + item.key
	}

	return item.key + "  "
}

// helpSpace calculates the aligned key column width for commands and options.
func (r helpRenderer) helpSpace(commands map[string][]commandHelpItem, flags []helpItem, hasSubgroups bool) int {
	space := 10

	for _, items := range commands {
		for _, item := range items {
			space = max(space, lipgloss.Width(displayCommandKey(item, hasSubgroups))+2)
		}
	}

	for _, item := range flags {
		space = max(space, lipgloss.Width(item.key)+2)
	}

	return space
}

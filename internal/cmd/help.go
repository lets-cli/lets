package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
	"github.com/lets-cli/fang"
	"github.com/lets-cli/lets/internal/docopt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type helpItem struct {
	key  string
	help string
}

type commandHelpItem struct {
	name     string
	subgroup string
	key      string
	help     string
}

var optionalArgsRe = regexp.MustCompile(`(\[.*\])`)

func HelpRenderer(cmd *cobra.Command, ctx fang.HelpContext) {
	renderLongShort(ctx.Writer, ctx.Styles, ctx.Width, cmpOr(cmd.Long, cmd.Short))

	usage := styleHelpUsage(cmd, ctx.Styles.Codeblock.Program, true)
	examples := fang.StyleExamples(cmd, ctx.Styles)
	blockStyle := compactCodeBlockStyle(ctx, append([]string{usage}, examples...)...)
	usageTitle := ctx.Styles.Title.Margin(0, 0)
	sectionTitle := compactTitleStyle(ctx.Styles)

	_, _ = fmt.Fprintln(ctx.Writer, usageTitle.Render("usage"))
	_, _ = fmt.Fprintln(ctx.Writer, blockStyle.Render(usage))

	if len(examples) > 0 {
		cw := blockStyle.GetWidth() - blockStyle.GetHorizontalPadding()
		_, _ = fmt.Fprintln(ctx.Writer, sectionTitle.Render("examples"))
		for i, example := range examples {
			if lipgloss.Width(example) > cw {
				examples[i] = ansi.Truncate(example, cw, "…")
			}
		}
		_, _ = fmt.Fprintln(ctx.Writer, blockStyle.Render(strings.Join(examples, "\n")))
	}

	groups, groupKeys := helpGroups(cmd)
	commands := helpCommands(cmd, ctx.Styles)
	options, optionsTitle := helpOptions(cmd, ctx.Styles)
	hasSubgroups := hasMultipleSubgroups(commands)
	space := helpSpace(commands, options, hasSubgroups)

	for _, groupID := range groupKeys {
		items := commands[groupID]
		if len(items) == 0 {
			continue
		}
		renderCommandGroup(ctx.Writer, ctx.Styles, space, groups[groupID], items, hasSubgroups)
	}

	if len(options) > 0 {
		renderHelpGroup(ctx.Writer, ctx.Styles, space, optionsTitle, options)
	}

	_, _ = fmt.Fprintln(ctx.Writer)
}

func ErrorHandler(w io.Writer, styles fang.Styles, err error) {
	if w, ok := w.(term.File); ok {
		if !term.IsTerminal(w.Fd()) {
			_, _ = fmt.Fprintln(w, err.Error())
			return
		}
	}

	errorHeader := styles.ErrorHeader
	errorText := styles.ErrorText

	_, _ = fmt.Fprintln(w, errorHeader.String())
	_, _ = fmt.Fprintln(w, errorText.Render(err.Error()+"."))
	_, _ = fmt.Fprintln(w)
	if isUsageError(err) {
		_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
			lipgloss.Left,
			errorText.UnsetWidth().Render("Try"),
			" ",
			styles.Program.Flag.Render("--help"),
			" for usage.",
		))
		_, _ = fmt.Fprintln(w)
	}
}

func isUsageError(err error) bool {
	s := err.Error()
	for _, prefix := range []string{
		"flag needs an argument:",
		"unknown flag:",
		"unknown shorthand flag:",
		"unknown command",
		"invalid argument",
	} {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}

	return false
}

func cmpOr(v1 string, v2 string) string {
	if v1 != "" {
		return v1
	}

	return v2
}

func compactTitleStyle(styles fang.Styles) lipgloss.Style {
	return styles.Title.Margin(0, 0).MarginBottom(0).PaddingBottom(0)
}

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

func styleHelpUsage(cmd *cobra.Command, styles fang.Program, complete bool) string {
	usage := cmd.Annotations[annotationHelpUsage]
	if usage == "" {
		return fang.StyleUsage(cmd, styles, complete)
	}

	lines := make([]string, 0)
	for _, line := range strings.Split(usage, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		lines = append(lines, styleHelpUsageLine(cmd, styles, line, complete))
	}

	if len(lines) == 0 {
		return fang.StyleUsage(cmd, styles, complete)
	}

	return strings.Join(lines, "\n")
}

func styleHelpUsageLine(cmd *cobra.Command, styles fang.Program, usage string, complete bool) string {
	if complete {
		usage = completeHelpUsage(cmd, usage)
	}

	return styleUsageText(cmd, styles, usage, complete)
}

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

	useLine := []string{}
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

func renderLongShort(w io.Writer, styles fang.Styles, width int, longShort string) {
	if longShort == "" {
		return
	}

	longShort = strings.TrimRight(longShort, "\n")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, styles.Text.Width(width).Render(longShort))
}

func helpGroups(cmd *cobra.Command) (map[string]string, []string) {
	ids := []string{""}
	groups := map[string]string{"": "commands"}

	for _, group := range cmd.Groups() {
		ids = append(ids, group.ID)
		groups[group.ID] = group.Title
	}

	return groups, ids
}

func helpCommands(cmd *cobra.Command, styles fang.Styles) map[string][]commandHelpItem {
	commands := map[string][]commandHelpItem{}

	for _, subCmd := range cmd.Commands() {
		if !subCmd.IsAvailableCommand() && subCmd.Name() != "help" {
			continue
		}

		commands[subCmd.GroupID] = append(commands[subCmd.GroupID], commandHelpItem{
			name:     subCmd.Name(),
			subgroup: subCmd.Annotations[annotationSubGroupName],
			key:      styleHelpUsage(subCmd, styles.Program, false),
			help:     renderHelpDescription(styles, subCmd.Short),
		})
	}

	for groupID := range commands {
		sort.Slice(commands[groupID], func(i, j int) bool {
			return commands[groupID][i].name < commands[groupID][j].name
		})
	}

	return commands
}

func helpOptions(cmd *cobra.Command, styles fang.Styles) ([]helpItem, string) {
	items := make([]helpItem, 0)
	docoptOptions := commandHelpOptions(cmd)

	for _, option := range docoptOptions {
		items = append(items, helpItem{
			key:  renderDocoptFlag(styles.Program, option.Display),
			help: renderHelpDescription(styles, option.Description),
		})
	}

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden || shouldSkipHelpFlag(cmd, flag) {
			return
		}

		help := renderHelpDescription(styles, flag.Usage)
		if flag.DefValue != "" && flag.DefValue != "false" && flag.DefValue != "0" && flag.DefValue != "[]" {
			help += styles.FlagDefault.Render(" (" + flag.DefValue + ")")
		}

		items = append(items, helpItem{
			key:  renderCobraFlag(styles.Program, flag),
			help: help,
		})
	})

	if len(docoptOptions) > 0 {
		return items, "options"
	}

	return items, "flags"
}

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

func shouldSkipHelpFlag(cmd *cobra.Command, flag *pflag.Flag) bool {
	return flag.Name == "help" && cmd != cmd.Root() && !flag.Changed
}

func renderCobraFlag(styles fang.Program, flag *pflag.Flag) string {
	if flag.Shorthand == "" {
		return styles.Flag.Render("--" + flag.Name)
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		styles.Flag.Render("-"+flag.Shorthand+" --"+flag.Name),
	)
}

func renderDocoptFlag(styles fang.Program, display string) string {
	parts := strings.Split(display, ", ")
	rendered := make([]string, 0, len(parts))

	for _, part := range parts {
		rendered = append(rendered, renderDocoptFlagPart(styles, part))
	}

	return strings.Join(rendered, styles.DimmedArgument.Render(", "))
}

func renderDocoptFlagPart(styles fang.Program, part string) string {
	if left, right, ok := strings.Cut(part, "="); ok {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			styles.Flag.Render(left+"="),
			styles.Flag.Render(right),
		)
	}

	return styles.Flag.Render(part)
}

func renderHelpDescription(styles fang.Styles, usage string) string {
	noTransform := styles.FlagDescription.UnsetTransform()
	lines := make([]string, 0, 1)

	for i, line := range strings.Split(usage, "\n") {
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

func renderHelpGroup(w io.Writer, styles fang.Styles, space int, title string, items []helpItem) {
	_, _ = fmt.Fprintln(w, compactTitleStyle(styles).Render(title))
	for _, item := range items {
		renderHelpItem(w, space, item.key, item.help)
	}
}

func renderCommandGroup(w io.Writer, styles fang.Styles, space int, title string, items []commandHelpItem, hasSubgroups bool) {
	_, _ = fmt.Fprintln(w, compactTitleStyle(styles).Render(title))

	subgroupNames := subgroupNames(items)
	showSubgroupTitles := len(subgroupNames) > 1

	for _, subgroup := range subgroupNames {
		if showSubgroupTitles {
			_, _ = fmt.Fprintln(w)
			_, _ = fmt.Fprintln(w, lipgloss.NewStyle().PaddingLeft(2).Render(styles.Text.Render(subgroup)))
		}

		for _, item := range items {
			if item.subgroup != subgroup {
				continue
			}
			renderHelpItem(w, space, displayCommandKey(item, hasSubgroups), item.help)
		}
	}

	for _, item := range items {
		if item.subgroup != "" {
			continue
		}
		renderHelpItem(w, space, displayCommandKey(item, hasSubgroups), item.help)
	}
}

func renderHelpItem(w io.Writer, space int, key string, help string) {
	_, _ = fmt.Fprintln(w, lipgloss.JoinHorizontal(
		lipgloss.Left,
		lipgloss.NewStyle().PaddingLeft(2).Render(key),
		strings.Repeat(" ", max(space-lipgloss.Width(key), 0)),
		help,
	))
}

func subgroupNames(items []commandHelpItem) []string {
	seen := map[string]struct{}{}
	names := make([]string, 0, len(items))

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

	sort.Strings(names)

	return names
}

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

func displayCommandKey(item commandHelpItem, hasSubgroups bool) string {
	if !hasSubgroups {
		return item.key
	}
	if item.subgroup != "" {
		return "  " + item.key
	}

	return item.key + "  "
}

func helpSpace(commands map[string][]commandHelpItem, flags []helpItem, hasSubgroups bool) int {
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

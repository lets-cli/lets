package docopt

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	dopt "github.com/kindermax/docopt.go"
)

// aliases for docopt types.
type (
	Opts   = dopt.Opts
	Option = dopt.Option
)

var docoptParser = &dopt.Parser{
	HelpHandler:   dopt.NoHelpHandler,
	OptionsFirst:  false,
	SkipHelpFlags: false,
}

// Parse parses docopts for command options with args from os.Args.
func Parse(cmdName string, args []string, docopts string) (Opts, error) {
	// no options at all
	if docopts == "" {
		return Opts{}, nil
	}

	return docoptParser.ParseArgs(docopts, append([]string{cmdName}, args...), "")
}

// ParseOptions parses docopts only to get all available options for a command.
func ParseOptions(docopts string, cmdName string) ([]Option, error) {
	return docoptParser.ParseOptions(docopts, []string{cmdName})
}

func OptsToLetsOpt(opts Opts) map[string]string {
	envMap := make(map[string]string, len(opts))

	for origKey, value := range opts {
		if !isOptKey(origKey) {
			continue
		}

		key := normalizeKey(origKey)
		envKey := "LETSOPT_" + key

		var strValue string

		switch value := value.(type) {
		case string:
			strValue = value
		case bool:
			if value {
				strValue = strconv.FormatBool(value)
			} else {
				strValue = ""
			}
		case []string:
			strValue = strings.Join(value, " ")
		case nil:
			strValue = ""
		default:
			strValue = ""
		}

		envMap[envKey] = strValue
	}

	return envMap
}

func OptsToLetsCli(opts Opts) map[string]string {
	cliMap := make(map[string]string, len(opts))
	formatVal := func(k, v string) string {
		return fmt.Sprintf("%s %s", k, v)
	}

	for origKey, value := range opts {
		if !isOptKey(origKey) {
			continue
		}

		key := normalizeKey(origKey)
		cliKey := "LETSCLI_" + key

		var strValue string

		switch value := value.(type) {
		case string:
			if value != "" {
				strValue = formatVal(origKey, value)
			}
		case bool:
			if value {
				strValue = origKey
			}
		case []string:
			if len(value) == 0 {
				strValue = ""
			} else {
				values := value
				if strings.HasPrefix(origKey, "-") {
					values = append([]string{origKey}, values...)
				}
				// TODO maybe we should escape each value
				strValue = strings.Join(values, " ")
			}
		case nil:
			strValue = ""
		default:
			strValue = ""
		}

		cliMap[cliKey] = strValue
	}

	return cliMap
}

func isOptKey(key string) bool {
	if key == "--" {
		return false
	}

	if strings.HasPrefix(key, "--") {
		return true
	}

	if strings.HasPrefix(key, "-") {
		return true
	}

	if strings.HasPrefix(key, "<") && strings.HasSuffix(key, ">") {
		return true
	}

	return false
}

func normalizeKey(origKey string) string {
	key := strings.TrimLeft(origKey, "-")
	key = strings.TrimLeft(key, "<")
	key = strings.TrimRight(key, ">")
	key = strings.ReplaceAll(key, "-", "_")
	key = strings.ToUpper(key)

	return key
}

type docoptParts struct {
	Usage   string
	Options string
	Example string
}

type HelpOption struct {
	Display     string `json:"display"`
	Description string `json:"description"`
	Name        string `json:"name,omitempty"`
	Short       string `json:"short,omitempty"`
	Long        string `json:"long,omitempty"`
	Kind        string `json:"kind,omitempty"`
}

var helpOptionSeparator = regexp.MustCompile(`\s{2,}`)

func ParseDocoptParts(docopt string) docoptParts {
	sections := map[string]*strings.Builder{
		"usage":   {},
		"options": {},
		"example": {},
	}

	section := ""

	for line := range strings.SplitSeq(docopt, "\n") {
		switch {
		case strings.HasPrefix(line, "Usage:"):
			section = "usage"
			line = strings.TrimSpace(strings.TrimPrefix(line, "Usage:"))
		case strings.HasPrefix(line, "Options:"):
			section = "options"
			line = strings.TrimSpace(strings.TrimPrefix(line, "Options:"))
		case strings.HasPrefix(line, "Example:"):
			section = "example"
			line = strings.TrimSpace(strings.TrimPrefix(line, "Example:"))
		}

		if section == "" || line == "" {
			continue
		}

		text := sections[section]
		if text.Len() > 0 {
			text.WriteByte('\n')
		}

		text.WriteString(line)
	}

	return docoptParts{
		Usage:   sections["usage"].String(),
		Options: sections["options"].String(),
		Example: sections["example"].String(),
	}
}

func ParseHelpOptions(docopt string, cmdName string) []HelpOption {
	parts := ParseDocoptParts(docopt)
	if parts.Options == "" {
		return nil
	}

	rawOpts, err := ParseOptions(docopt, cmdName)
	if err != nil {
		rawOpts = nil
	}

	var options []HelpOption

	for line := range strings.SplitSeq(parts.Options, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}

		parts := helpOptionSeparator.Split(trimmed, 2)
		if len(parts) == 0 {
			continue
		}

		display := strings.TrimSpace(parts[0])
		if display == "" {
			continue
		}

		if !strings.HasPrefix(display, "-") && !strings.HasPrefix(display, "<") {
			if len(options) == 0 {
				continue
			}

			description := strings.TrimSpace(trimmed)
			if description == "" {
				continue
			}

			if options[len(options)-1].Description != "" {
				options[len(options)-1].Description += "\n"
			}

			options[len(options)-1].Description += description

			continue
		}

		option := HelpOption{Display: display}
		if len(parts) > 1 {
			option.Description = strings.TrimSpace(parts[1])
		}

		if strings.HasPrefix(display, "<") {
			option.Kind = "arg"
		} else {
			option.Kind = "flag"
		}

		for _, rawOpt := range rawOpts {
			if rawOpt.Name == cmdName {
				continue
			}

			if !strings.Contains(display, rawOpt.Name) && (rawOpt.Short == "" || !strings.Contains(display, rawOpt.Short)) {
				continue
			}

			option.Name = rawOpt.Name
			option.Short = rawOpt.Short
			option.Long = rawOpt.Long

			break
		}

		options = append(options, option)
	}

	return options
}

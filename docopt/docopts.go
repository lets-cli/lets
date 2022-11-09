package docopt

import (
	"fmt"
	"strconv"
	"strings"

	dopt "github.com/docopt/docopt-go"
)

// aliases for docopt types.
type Opts = dopt.Opts
type Option = dopt.Option

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
		envKey := fmt.Sprintf("LETSOPT_%s", key)

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
		cliKey := fmt.Sprintf("LETSCLI_%s", key)

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

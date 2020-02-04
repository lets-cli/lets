package command

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
)

var DocoptParser = &docopt.Parser{
	HelpHandler:   docopt.NoHelpHandler,
	OptionsFirst:  false,
	SkipHelpFlags: false,
}

// ParseDocopts parses docopts for command options with args from os.Args
// TODO maybe this must be a struct method
func ParseDocopts(cmd Command) (map[string]string, error) {
	// just command name in args
	if len(os.Args[1:]) == 1 && os.Args[1] == cmd.Name {
		return make(map[string]string), nil
	}
	// no options at all
	if cmd.RawOptions == "" {
		return make(map[string]string), nil
	}
	opts, err := DocoptParser.ParseArgs(cmd.RawOptions, os.Args[1:], "")

	if err != nil {
		return nil, err
	}
	return normalizeOpts(opts), nil
}

func normalizeOpts(opts map[string]interface{}) map[string]string {
	// TODO
	// non-passed flags (counted 0)
	// passed flags
	// passed several times
	// non-passed positional args
	// passed positional args
	// list (still not get it)
	envMap := make(map[string]string, len(opts))
	for origKey, value := range opts {
		key := normalizeKey(origKey)
		envKey := fmt.Sprintf("LETSOPT_%s", key)
		var strValue string
		switch value.(type) {
		case string:
			strValue = value.(string)
		case bool:
			strValue = strconv.FormatBool(value.(bool))
		case []string:
			strValue = strings.Join(value.([]string), " ")
		default:
			strValue = ""
		}
		envMap[envKey] = strValue
	}
	return envMap
}

func normalizeKey(origKey string) string {
	key := strings.TrimLeft(origKey, "-")
	key = strings.TrimLeft(key, "<")
	key = strings.TrimRight(key, ">")
	key = strings.ReplaceAll(key, "-", "_")
	key = strings.ToUpper(key)
	return key
}

package env

import (
	"os"
	"strconv"
)

const MAX_DEBUG_LEVEL = 2

// IsDebug checks env or --debug flag
// - LETS_DEBUG env var - if present and is a trythy value - we in debug mode.
// - flag --debug - if present - we in debug mode
func IsDebug() bool {
	return DebugLevel() > 0
}

type debugLevel struct {
	value int
	ready bool
}

func (d *debugLevel) set(value int) {
	d.value = value
	d.ready = true
}

var debugLvl = &debugLevel{}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// DebugLevel determines verbosity level of debug logs.
// If LETS_DEBUG set to int - then verbosity is 1 or 2
// If --debug or -d used multiple times - then verbosity is 1 or 2
// If -dd used - then verbosity is 2
func DebugLevel() int {
	if debugLvl.ready {
		return debugLvl.value
	}

	level := 0

	envValue := os.Getenv("LETS_DEBUG")

	level, err := strconv.Atoi(envValue)

	if err != nil {
		// probably not integer, try just determine bool value
		debug, err := strconv.ParseBool(envValue)
		if err != nil || !debug {
			level = 0
		} else {
			level = 1
		}
	}

	if level == 0 {
		for _, arg := range os.Args {
			if arg == "--debug" || arg == "-d" {
				level += 1
			} else if arg == "-dd" {
				level += 2
			}
		}
	}

	debugLvl.set(min(level, MAX_DEBUG_LEVEL))

	return debugLvl.value
}

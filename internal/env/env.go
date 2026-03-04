package env

import (
	"os"
	"strconv"
)

const MaxDebugLevel = 2

func IsDebug() bool { return DebugLevel() > 0 }

type debug struct {
	level int
	ready bool
}

func (d *debug) set(level int) {
	d.level = level
	d.ready = true
}

var debugLevel = &debug{}

// DebugLevel determines verbosity level of debug logs.
// If LETS_DEBUG set to int - then verbosity is 1 or 2
// If --debug or -d used multiple times - then verbosity is 1 or 2
// If -dd used - then verbosity is 2.
// When determined - set debug level globally.
func SetDebugLevel(level int) int {
	if level == 0 {
		envValue := os.Getenv("LETS_DEBUG")

		envLevel, err := strconv.Atoi(envValue)
		level = envLevel

		if err != nil {
			// probably not integer, try just determine bool value
			debug, err := strconv.ParseBool(envValue)
			if err != nil || !debug {
				level = 0
			} else {
				level = 1
			}
		}
	}

	level = min(level, MaxDebugLevel)

	debugLevel.set(level)
	return level
}

func DebugLevel() int {
	if !debugLevel.ready {
		panic("must run SetDebugLevel first")
	}

	return debugLevel.level
}

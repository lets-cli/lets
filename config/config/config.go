package config

var (
	// COMMANDS is a top-level directive. Includes all commands to run.
	COMMANDS = "commands"
	SHELL    = "shell"
	ENV      = "env"
	EvalEnv  = "eval_env"
	MIXINS   = "mixins"
	VERSION  = "version"
	BEFORE   = "before"
)

var (
	ValidConfigFields      = []string{COMMANDS, SHELL, ENV, EvalEnv, MIXINS, VERSION, BEFORE}
	ValidMixinConfigFields = []string{COMMANDS, ENV, EvalEnv, BEFORE}
)

// Config is a struct for loaded config file.
type Config struct {
	// absolute path to work dir - where config is placed
	WorkDir  string
	FilePath string
	Commands map[string]Command
	Shell    string
	// before is a script which will be included before every cmd
	Before  string
	Env     map[string]string
	Version string
	isMixin bool // if true, we consider config as mixin and apply different parsing and validation
	// absolute path to .lets
	DotLetsDir string
}

func NewConfig(workDir string, configAbsPath string, dotLetsDir string) *Config {
	return &Config{
		Commands: make(map[string]Command),
		Env:      make(map[string]string),
		WorkDir:  workDir,
		FilePath: configAbsPath,
		DotLetsDir: dotLetsDir,
	}
}

func NewMixinConfig(workDir string, configAbsPath string, dotLetsDir string) *Config {
	cfg := NewConfig(workDir, configAbsPath, dotLetsDir)
	cfg.isMixin = true

	return cfg
}

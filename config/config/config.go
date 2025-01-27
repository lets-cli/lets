package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lets-cli/lets/config/path"
	"github.com/lets-cli/lets/util"
	"gopkg.in/yaml.v3"
)

// Config is a struct for loaded config file.
type Config struct {
	// absolute path to work dir - where config is placed
	WorkDir string
	// absolute path for lets config file
	FilePath string
	Commands Commands
	Shell    string
	// before is a script which will be included before every cmd
	Before string
	// init is a script which will be called exactly once before any command calls
	Init    string
	Env     *Envs
	Version string
	isMixin bool // if true, we consider config as mixin and apply different parsing and validation
	// absolute path to .lets
	DotLetsDir string
	// absolute path to .lets/checksums
	ChecksumsDir string
	// absolute path to .lets/mixins
	MixinsDir string
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var config struct {
		Version  Version
		Mixins   []*Mixin
		Commands Commands
		Shell    string
		Before   string
		Init     string
		Env      *Envs
		EvalEnv  *Envs `yaml:"eval_env"`
	}

	if err := unmarshal(&config); err != nil {
		return err
	}

	c.Version = string(config.Version)
	c.Init = config.Init
	c.Commands = config.Commands
	if c.Commands == nil {
		c.Commands = make(Commands, 0)
	}

	c.Shell = config.Shell

	if c.Shell == "" && !c.isMixin {
		return errors.New("'shell' is required")
	}

	c.Before = config.Before
	c.Env = config.Env
	if c.Env == nil {
		c.Env = &Envs{}
	}

	// support for deprecated eval_env
	_ = config.EvalEnv.Range(func(name string, value Env) error {
		c.Env.Set(name, Env{Name: name, Sh: value.Value})

		return nil
	})

	for name, cmd := range c.Commands {
		cmd.Name = name
	}

	if err := c.readMixins(config.Mixins); err != nil {
		return err
	}

	if !c.isMixin {
		if err := c.resolveRefs(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) resolveRefs() error {
	commandsFromRef := []*Command{}
	for _, cmd := range c.Commands {
		// resolve command by ref
		if ref := cmd.ref; ref != nil {
			command, exists := c.Commands[ref.Name]
			if !exists {
				return fmt.Errorf("ref '%s' points to command '%s' which is not exist", cmd.Name, ref.Name)
			}

			command = command.Clone()
			command.Name = cmd.Name
			command.Args = append(command.Args, ref.Args...)
			// fixing docopt string
			if command.Docopts != "" {
				command.Docopts = strings.Replace(
					command.Docopts,
					"lets "+ref.Name,
					"lets "+command.Name,
					1,
				)
			}
			commandsFromRef = append(commandsFromRef, command)
		}
	}

	for _, cmd := range commandsFromRef {
		c.Commands[cmd.Name] = cmd
	}

	return nil
}

func joinBeforeScripts(beforeScripts ...string) string {
	buf := new(bytes.Buffer)

	for _, script := range beforeScripts {
		if script == "" {
			continue
		}
		buf.WriteString(script)
		buf.WriteString("\n")
	}

	return buf.String()
}

// Merge main and mixin configs. If there is a conflict - return error as we do not override values.
func (c *Config) mergeMixin(mixin *Config) error {
	for cmdName, cmd := range mixin.Commands {
		if _, conflict := c.Commands[cmdName]; conflict {
			return fmt.Errorf(
				"command '%s' from mixin '%s' is already declared in main config's commands",
				cmdName, mixin.FilePath,
			)
		}

		cmd.Name = cmdName
		c.Commands[cmdName] = cmd
	}

	err := mixin.Env.Range(func(key string, value Env) error {
		if c.Env.Has(key) {
			return fmt.Errorf("env '%s' from mixin '%s' is already declared in main config's env", key, mixin.FilePath)
		}

		c.Env.Set(key, value)

		return nil
	})
	if err != nil {
		return err
	}

	c.Before = joinBeforeScripts(
		c.Before,
		mixin.Before,
	)

	return nil
}

func (c *Config) readMixin(mixin *Mixin) error {
	if mixin.IsRemote() {
		mixin.Remote.mixinsDir = c.MixinsDir

		rm := mixin.Remote

		data, err := rm.tryRead()
		if err != nil {
			return err
		}

		if data == nil {
			data, err = rm.download()
			if err != nil {
				return err
			}
		}

		// TODO: what if multiple mixins have same commands
		//  1 option - fail and suggest use to namespace all commands in remote mixin
		//  2 option - namespace it (this may require specifying namespace in mixin config or in main config mixin section)

		mixinCfg := NewMixinConfig(c, rm.Filename())
		reader := bytes.NewReader(data)
		if err := yaml.NewDecoder(reader).Decode(mixinCfg); err != nil {
			return fmt.Errorf("failed to parse remote mixin config '%s': %w", rm.URL, err)
		}

		if err := c.mergeMixin(mixinCfg); err != nil {
			return fmt.Errorf("failed to merge remote mixin config '%s' with main config: %w", rm.URL, err)
		}

		if err := rm.persist(data); err != nil {
			return fmt.Errorf("failed to persist remote mixin config %s: %w", rm.URL, err)
		}
	} else {
		mixinAbsPath, err := path.GetFullConfigPath(mixin.FileName, c.WorkDir)
		if err != nil {
			if mixin.Ignored && errors.Is(err, path.ErrFileNotExists) {
				return nil
			}

			// complain non-existed mixin only if its filename does not start with dash `-`
			return err
		}

		file, err := os.Open(mixinAbsPath)
		if err != nil {
			return fmt.Errorf("failed to read mixin config %s: %w", mixin.FileName, err)
		}

		// TODO(maybe bug): probably not filename but mixinAbsPath
		mixinCfg := NewMixinConfig(c, mixin.FileName)
		if err := yaml.NewDecoder(file).Decode(mixinCfg); err != nil {
			return fmt.Errorf("can not parse mixin config %s:\n%w", mixin.FileName, err)
		}

		if err := c.mergeMixin(mixinCfg); err != nil {
			return fmt.Errorf("failed to merge mixin config '%s' with main config: %w", mixin.FileName, err)
		}
	}

	return nil
}

func (c *Config) readMixins(mixins []*Mixin) error {
	if c.isMixin {
		// disallow recursive mixins
		return nil
	}

	if len(mixins) == 0 {
		return nil
	}

	if err := c.createMixinsDir(); err != nil {
		return err
	}

	for _, mixin := range mixins {
		if err := c.readMixin(mixin); err != nil {
			// TODO: check if error is correct, concise and for humans
			return err
		}
	}

	return nil
}

func (c *Config) GetEnv() map[string]string {
	return c.Env.Dump()
}

// SetupEnv must be called once. It is not intended to be called
// multiple times hence does not have mutex.
func (c *Config) SetupEnv() error {
	if err := c.Env.Execute(*c); err != nil {
		return err
	}

	// expand env for args
	for _, cmd := range c.Commands {
		for idx, arg := range cmd.Args {
			// we have to expand env here on our own, since this args not came from users tty, and not expanded before lets
			cmd.Args[idx] = os.Expand(arg, func(key string) string {
				return c.Env.Mapping[key].Value
			})
		}
	}

	return nil
}

func NewConfig(workDir string, configAbsPath string, dotLetsDir string) *Config {
	return &Config{
		WorkDir:      workDir,
		FilePath:     configAbsPath,
		DotLetsDir:   dotLetsDir,
		ChecksumsDir: filepath.Join(dotLetsDir, "checksums"),
		MixinsDir:    filepath.Join(dotLetsDir, "mixins"),
	}
}

func NewMixinConfig(cfg *Config, configAbsPath string) *Config {
	mixin := NewConfig(cfg.WorkDir, configAbsPath, cfg.DotLetsDir)
	mixin.isMixin = true

	return mixin
}

func (c *Config) CreateChecksumsDir() error {
	if err := util.SafeCreateDir(c.ChecksumsDir); err != nil {
		return fmt.Errorf("can not create %s: %w", c.ChecksumsDir, err)
	}

	return nil
}

func (c *Config) createMixinsDir() error {
	if err := util.SafeCreateDir(c.MixinsDir); err != nil {
		return fmt.Errorf("can not create %s: %w", c.MixinsDir, err)
	}

	return nil
}

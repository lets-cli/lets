package config

import (
	"fmt"
	"strings"
)

type Cmds struct {
	Commands []*Cmd
	Append   bool
	Parallel bool
}

type Cmd struct {
	Name   string
	Script string // list will be joined
}

// A workaround function which helps to prevent breaking
// strings with special symbols (' ', '*', '$', '#'...)
// When you run a command with an argument containing one of these, you put it into quotation marks:
// lets alembic -n dev revision --autogenerate -m "revision message"
// which makes shell understand that "revision message" is a single argument, but not two args
// The problem is, lets constructs a script string
// and then passes it to an appropriate interpreter (sh -c $SCRIPT)
// so we need to wrap args with quotation marks to prevent breaking
// This also solves problem with json params: --key='{"value": 1}' => '--key={"value": 1}'.
func escapeArgs(args []string) []string {
	var escapedArgs []string

	for _, arg := range args {
		// wraps every argument with quotation marks to avoid ambiguity
		// TODO: maybe use some kind of blacklist symbols to wrap only necessary args
		escapedArg := fmt.Sprintf("'%s'", arg)
		escapedArgs = append(escapedArgs, escapedArg)
	}

	return escapedArgs
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (c *Cmds) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var script string
	if err := unmarshal(&script); err == nil {
		c.Commands = []*Cmd{{Name: "", Script: script}}

		return nil
	}

	var cmdList []string
	if err := unmarshal(&cmdList); err == nil {
		script := strings.TrimSpace(strings.Join(cmdList, " "))
		c.Commands = []*Cmd{{Name: "", Script: script}}
		c.Append = true

		return nil
	}

	var cmdMap map[string]string
	if err := unmarshal(&cmdMap); err == nil {
		for name, script := range cmdMap {
			c.Commands = append(c.Commands, &Cmd{Name: name, Script: script})
		}
		c.Parallel = true

		return nil
	}

	return nil
}

func (c Cmds) Clone() Cmds {
	commands := make([]*Cmd, len(c.Commands))

	for idx, cmd := range c.Commands {
		commands[idx] = &Cmd{
			Name:   cmd.Name,
			Script: cmd.Script,
		}
	}

	cmds := Cmds{
		Commands: commands,
		Parallel: c.Parallel,
	}

	return cmds
}

// SingleCommand returns Cmd only if there is only one command.
func (c Cmds) SingleCommand() *Cmd {
	if len(c.Commands) == 1 {
		return c.Commands[0]
	}

	return nil
}

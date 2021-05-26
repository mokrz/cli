package cli

import (
	"flag"
)

type HandleFn func(cmd *Command, args []string) error

type Command struct {
	Name     string
	Handler  HandleFn
	Flags    *flag.FlagSet
	children map[string]*Command
	parent *Command
}

func NewCommand(name string, handler HandleFn) *Command {
	return &Command{
		Name:    name,
		Handler: handler,
		Flags: flag.NewFlagSet(name, flag.ExitOnError),
		children: make(map[string]*Command),
	}
}

func (c *Command) AddCommand(cmd *Command) {
	c.children[cmd.Name] = cmd
	c.children[cmd.Name].parent = c
}

func (c *Command) GetSubcommand(name string) (*Command) {
	return c.children[name]
}

func (c *Command) Execute(args []string) error {

	if len(args) > 1 {

		if subCmd := c.GetSubcommand(args[1]); subCmd != nil {
			return subCmd.Execute(args[1:])
		}
	}

	return c.Handler(c, args)
}

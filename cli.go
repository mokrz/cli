package cli

import (
	"flag"
	"io"
	"os"
)

type HandleFn func(cmd *Command, args []string) error

type Command struct {
	Name           string
	Handler        HandleFn
	Flags          *flag.FlagSet
	children       map[string]*Command
	parent         *Command
	stdout, stderr io.Writer
}

func NewCommand(name string, handler HandleFn) *Command {
	return &Command{
		Name:     name,
		Handler:  handler,
		Flags:    flag.NewFlagSet(name, flag.ExitOnError),
		children: make(map[string]*Command),
	}
}

func (c *Command) AddCommand(cmd *Command) {
	c.children[cmd.Name] = cmd
	c.children[cmd.Name].parent = c
}

func (c *Command) Subcommand(name string) *Command {
	return c.children[name]
}

func (c *Command) Execute(args []string) error {

	if len(args) > 1 {

		if subCmd := c.Subcommand(args[1]); subCmd != nil {
			return subCmd.Execute(args[1:])
		}
	}

	return c.Handler(c, args)
}

func (c *Command) SetStdout(out io.Writer) {
	c.stdout = out
}

func (c *Command) SetStderr(err io.Writer) {
	c.stderr = err
}

func (c *Command) Stdout() io.Writer {

	if c.stdout != nil {
		return c.stdout
	} else if c.parent != nil {
		return c.parent.Stdout()
	}

	return os.Stdout
}

func (c *Command) Stderr() io.Writer {

	if c.stderr != nil {
		return c.stderr
	} else if c.parent != nil {
		return c.parent.Stderr()
	}

	return os.Stderr
}

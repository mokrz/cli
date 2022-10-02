package cli

import (
	"flag"
	"io"
	"os"
)

type HandleFn func(cmd *Command, args []string) error

type Command struct {
	Name                 string
	Handler              HandleFn
	Flags                *flag.FlagSet
	children             map[string]*Command
	parent               *Command
	outStream, errStream io.Writer
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

func (c *Command) SubCommand(name string) *Command {
	return c.children[name]
}

func (c *Command) Execute(args []string) error {

	if len(args) > 1 {

		if subCmd := c.SubCommand(args[1]); subCmd != nil {
			return subCmd.Execute(args[1:])
		}
	}

	return c.Handler(c, args)
}

func (c *Command) SetOutputStream(out io.Writer) {
	c.outStream = out
}

func (c *Command) SetErrorStream(err io.Writer) {
	c.errStream = err
}

func (c *Command) OutputStream() io.Writer {

	if c.outStream != nil {
		return c.outStream
	} else if c.parent != nil {
		return c.parent.OutputStream()
	}

	return os.Stdout
}

func (c *Command) ErrorStream() io.Writer {

	if c.errStream != nil {
		return c.errStream
	} else if c.parent != nil {
		return c.parent.ErrorStream()
	}

	return os.Stderr
}

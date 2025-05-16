package commands

import (
	"errors"
	"fmt"
	"os"
	"slices"
)

var helpFlags []string = []string{
	"-help",
	"--help",
	"-h",
}

var (
	ErrCmdNotFound = errors.New("command not found")
	ErrNotReady    = errors.New("not ready yet")
	ErrNoSubOrFunc = errors.New("command has no saubcomands and execution func")
)

type CommandFunc func(ctx CommandContext) error

type CommandInterface interface {
	GetName() string
	GetDesc() string
	Help() error
	Execute(args ...string) error
}

func checkHelp(args ...string) bool {
	for _, flag := range helpFlags {
		if slices.Contains(args, flag) {
			return true
		}
	}
	return false
}

type Command struct {
	name     string
	desc     string
	commands map[string]CommandInterface
	cmdFunc  CommandFunc
	flags    []FlagDefinitionOpt
}

func (c *Command) GetName() string {
	return c.name
}

func (c *Command) GetDesc() string {
	return c.desc
}

// Prints help message to screen based on defined flags
func (c *Command) Help() error {
	fmt.Println(c.name)
	fmt.Printf("---\n\t%s\n---\n\n", c.desc)
	for _, cmd := range c.commands {
		fmt.Printf("\t\t%s\t%s\n", cmd.GetName(), cmd.GetDesc())
	}
	for _, flo := range c.flags {
		fl, err := flo()
		if err != nil {
			return err
		}
		fmt.Printf("\t\t-%s\t%s\n\t\t  %s\n\n", fl.GetName(), fl.GetDescription(), fl.GetShortList())
	}
	fmt.Println("\n---")

	return nil
}

// check for subcommand and executes it
func (c *Command) Execute(args ...string) error {
	// If doesnt have sub commands execute it registered function
	if len(c.commands) < 1 {
		if checkHelp(args...) {
			return c.Help()
		}
		// Create context first for execution if not help
		if c.cmdFunc == nil {
			return ErrNoSubOrFunc
		}
		return c.executeCmdFunc(args...)
	}

	if len(args) < 1 {
		return c.Help()
	}
	// serach for proper subcommand
	commandName := args[0]
	cmd, ok := c.commands[commandName]
	if !ok {
		return fmt.Errorf("%s: %w", commandName, c.Help())
	}

	// if subcommand found execute it
	return cmd.Execute(args[1:]...)
}

// This just starts command resolution from the beginning of arguments
func (c *Command) ExecuteAsRootCommand() error {
	if len(os.Args) < 2 {
		return c.Help()
	}
	return c.Execute(os.Args[1:]...)
}

func (c *Command) executeCmdFunc(args ...string) error {
	ctx, err := NewCommandLineContext(c.flags...)
	if err != nil {
		return err
	}
	ctx.Parse(args...)

	return c.cmdFunc(ctx)
}

type WithCommand func() (*Command, error)

func NewCommandWithSubcommands(name string, description string, commands ...*Command) *Command {
	cmdMap := map[string]CommandInterface{}
	for _, cmd := range commands {
		cmdMap[cmd.GetName()] = cmd
	}
	return &Command{
		commands: cmdMap,
		name:     name,
		desc:     description,
	}
}

func NewCommandWithFunc(name string, description string, cmdFunc CommandFunc, flags ...FlagDefinitionOpt) *Command {
	return &Command{
		name:    name,
		desc:    description,
		cmdFunc: cmdFunc,
		flags:   flags,
	}
}

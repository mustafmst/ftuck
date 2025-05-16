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
)

type CommandContext struct {
	Flags map[string]string
}

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
	commands map[string]CommandInterface
	cmdFunc  CommandFunc
}

func (cr *Command) GetName() string {
	return os.Args[0]
}

// Prints help message to screen based on defined flags
func (c *Command) Help() error {
	return ErrNotReady
}

// check for subcommand and executes it
func (c *Command) Execute(args ...string) error {
	// If doesnt have sub commands execute it registered function
	if len(c.commands) < 1 {
		if checkHelp(args...) {
			return c.Help()
		}
		// Create context first for execution if not help
		panic("not implemented")
	}
	// serach for proper subcommand
	commandName := args[0]
	cmd, ok := c.commands[commandName]
	if !ok {
		return fmt.Errorf("%s: %w", commandName, ErrCmdNotFound)
	}
	return cmd.Execute(args[1:]...)
}

func NewCommandWithSubcommands(name string, commands ...CommandInterface) *Command {
	cmdMap := map[string]CommandInterface{}
	for _, cmd := range commands {
		cmdMap[cmd.GetName()] = cmd
	}
	return &Command{
		commands: cmdMap,
	}
}

func NewCommandWithFunc(name string, cmdFunc CommandFunc) *Command {
}

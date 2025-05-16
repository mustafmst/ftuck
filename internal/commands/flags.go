package commands

import (
	"errors"
	"flag"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrUnsupportedFlagType = errors.New("unsupported flag type")
	ErrWrongArgType        = errors.New("wrong type of arg")
	ErrFlagNotFound        = errors.New("flag not found")
)

type FlagType string

var (
	StringArg FlagType = "STRING"
	IntArg    FlagType = "INT"
	BoolArg   FlagType = "BOOL"
)

type CommandContext interface {
	GetString(key string) (string, error)
	GetInt(key string) (int, error)
	GetBool(key string) (bool, error)
	Parse(args ...string)
}

type FlagDefinition struct {
	stringVal   string
	intVal      int
	boolVal     bool
	argType     FlagType
	flags       []string
	description string
}

func (a *FlagDefinition) GetShortList() string {
	if len(a.flags) < 2 {
		return ""
	}
	l := []string{}
	for _, f := range a.flags[1:] {
		l = append(l, "-"+f)
	}
	return strings.Join(l, " ")
}

func (a *FlagDefinition) GetName() string {
	return a.flags[0]
}

func (a *FlagDefinition) GetDescription() string {
	return a.description
}

func (a *FlagDefinition) registerFlags() {
	switch a.argType {
	case StringArg:
		for _, argFlag := range a.flags {
			flag.StringVar(&a.stringVal, argFlag, "", a.description)
		}
	case IntArg:
		for _, argFlag := range a.flags {
			flag.IntVar(&a.intVal, argFlag, -1, a.description)
		}
	case BoolArg:
		for _, argFlag := range a.flags {
			flag.BoolVar(&a.boolVal, argFlag, false, a.description)
		}
	}
}

type FlagDefinitionOpt func() (*FlagDefinition, error)

func RegisterFlag(name string, description string, argType FlagType, replacementFlags ...string) FlagDefinitionOpt {
	if !slices.Contains([]FlagType{StringArg, IntArg, BoolArg}, argType) {
		return func() (*FlagDefinition, error) {
			return nil, fmt.Errorf("(flag name: %s) %w", name, ErrUnsupportedFlagType)
		}
	}

	return func() (*FlagDefinition, error) {
		return &FlagDefinition{
			flags:       append([]string{name}, replacementFlags...),
			argType:     argType,
			description: description,
		}, nil
	}
}

type CommandLineContext struct {
	flags     map[string]*FlagDefinition
	wasParsed bool
}

// GetBool implements CommandContext.
func (c *CommandLineContext) GetBool(key string) (bool, error) {
	fl, ok := c.flags[key]

	if !ok {
		return false, fmt.Errorf("(key: %s) %w", key, ErrFlagNotFound)
	}

	if fl.argType != BoolArg {
		return false, fmt.Errorf("(key: %s) %w", key, ErrWrongArgType)
	}

	return fl.boolVal, nil
}

// GetInt implements CommandContext.
func (c *CommandLineContext) GetInt(key string) (int, error) {
	fl, ok := c.flags[key]

	if !ok {
		return -1, fmt.Errorf("(key: %s) %w", key, ErrFlagNotFound)
	}

	if fl.argType != IntArg {
		return -1, fmt.Errorf("(key: %s) %w", key, ErrWrongArgType)
	}

	return fl.intVal, nil
}

// GetString implements CommandContext.
func (c *CommandLineContext) GetString(key string) (string, error) {
	fl, ok := c.flags[key]

	if !ok {
		return "", fmt.Errorf("(key: %s) %w", key, ErrFlagNotFound)
	}

	if fl.argType != StringArg {
		return "", fmt.Errorf("(key: %s) %w", key, ErrWrongArgType)
	}

	return fl.stringVal, nil
}

// Parse implements CommandContext.
func (c *CommandLineContext) Parse(args ...string) {
	if c.wasParsed {
		return
	}
	defer func() {
		c.wasParsed = true
	}()
	for _, fl := range c.flags {
		fl.registerFlags()
	}
	flag.CommandLine.Parse(args)
}

func NewCommandLineContext(flags ...FlagDefinitionOpt) (*CommandLineContext, error) {
	flagMap := map[string]*FlagDefinition{}

	for _, flo := range flags {
		fl, err := flo()
		if err != nil {
			return nil, fmt.Errorf("flag creation error: %w", err)
		}
		flagMap[fl.GetName()] = fl
	}

	return &CommandLineContext{
		flags:     flagMap,
		wasParsed: false,
	}, nil
}

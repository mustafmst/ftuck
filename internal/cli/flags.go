package cli

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
	ErrDefValueType        = errors.New("wrong data type for default value")
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

type flagDefinition struct {
	stringVal   string
	intVal      int
	boolVal     bool
	argType     FlagType
	flags       []string
	description string
	defaultVal  any
}

func (a *flagDefinition) GetShortList() string {
	if len(a.flags) < 2 {
		return ""
	}
	l := []string{}
	for _, f := range a.flags[1:] {
		l = append(l, "-"+f)
	}
	return strings.Join(l, " ")
}

func (a *flagDefinition) GetName() string {
	return a.flags[0]
}

func (a *flagDefinition) GetDescription() string {
	return a.description
}

func (a *flagDefinition) registerFlags() {
	switch a.argType {
	case StringArg:
		dv, _ := a.defaultVal.(string)
		for _, argFlag := range a.flags {
			flag.StringVar(&a.stringVal, argFlag, dv, a.description)
		}
	case IntArg:
		dv, _ := a.defaultVal.(int)
		for _, argFlag := range a.flags {
			flag.IntVar(&a.intVal, argFlag, dv, a.description)
		}
	case BoolArg:
		dv, _ := a.defaultVal.(bool)
		for _, argFlag := range a.flags {
			flag.BoolVar(&a.boolVal, argFlag, dv, a.description)
		}
	}
}

type FlagDefinitionOpt func() (*flagDefinition, error)

func RegisterFlag(name string, description string, argType FlagType, defaultValue any, replacementFlags ...string) FlagDefinitionOpt {
	if !slices.Contains([]FlagType{StringArg, IntArg, BoolArg}, argType) {
		return func() (*flagDefinition, error) {
			return nil, fmt.Errorf("(flag name: %s) %w", name, ErrUnsupportedFlagType)
		}
	}

	var err error

	switch argType {
	case StringArg:
		_, ok := defaultValue.(string)
		if !ok {
			err = ErrDefValueType
		}
	case IntArg:
		_, ok := defaultValue.(int)
		if !ok {
			err = ErrDefValueType
		}
	case BoolArg:
		_, ok := defaultValue.(bool)
		if !ok {
			err = ErrDefValueType
		}
	}

	if err != nil {
		return func() (*flagDefinition, error) {
			return nil, fmt.Errorf("(flag name: %s) %w", name, err)
		}
	}

	return func() (*flagDefinition, error) {
		return &flagDefinition{
			flags:       append([]string{name}, replacementFlags...),
			argType:     argType,
			description: description,
			defaultVal:  defaultValue,
		}, nil
	}
}

type CommandLineContext struct {
	flags     map[string]*flagDefinition
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
	flagMap := map[string]*flagDefinition{}

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

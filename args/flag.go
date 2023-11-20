package args

import (
	"fmt"
	"strings"
)

func NewArgs(app *App, args ...*Flag) *Args {
	a := &Args{
		App: app,
	}

	for _, arg := range args {
		a.Flags = append(a.Flags, arg)
	}
	return a
}

func NewApp(name, usage, version string) *App {
	return &App{
		Name:       name,
		Usage:      usage,
		Copyright:  "",
		Version:    &version,
		Changelog:  nil,
		Banner:     nil,
		CommitHash: nil,
		Date:       nil,
	}
}

func NewFlag(name string, def string, desc string, fn func(int, []string) error) *Flag {
	return &Flag{
		Name:        []string{name},
		Description: desc,
		Default:     []string{def},
		Required:    true,
		Env:         []string{strings.ToUpper(fmt.Sprintf("A-%s", name))},
		SingleValue: true,
		HookFunc:    fn,
	}
}

func NewDefaultFlag(name string, def []string, desc string, fn func(int, []string) error) *Flag {
	return &Flag{
		Name:        []string{name},
		Description: desc,
		Default:     def,
		Required:    true,
		Env:         []string{strings.ToUpper(fmt.Sprintf("A-%s", name))},
		SingleValue: true,
		HookFunc:    fn,
	}
}

func NewNoValueFlag(name string, desc string) *Flag {
	return &Flag{
		Name:        []string{name},
		Description: desc,
		NoValue:     true,
		Env:         []string{strings.ToUpper(fmt.Sprintf("A-%s", name))},
	}
}

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

func NewRequiredFlag(name string, desc string, fn FlagHookFunc) *Flag {
	return &Flag{
		Name:        []string{name},
		Description: desc,
		Required:    true,
		Env:         []string{strings.ToUpper(fmt.Sprintf("A%s", strings.TrimPrefix(name, "--")))},
		SingleValue: true,
		HookFunc:    fn,
	}
}

func NewFlag(name string, def string, desc string, fn FlagHookFunc) *Flag {
	a := NewRequiredFlag(name, desc, fn)
	a.Default = []string{def}
	return a
}

func NewOptionFlag(name string, def string, desc string, fn FlagHookFunc) *Flag {
	a := NewRequiredFlag(name, desc, fn)
	a.Default = []string{def}
	a.Required = false
	return a
}

func NewDefaultFlag(name string, def []string, desc string, fn FlagHookFunc) *Flag {
	a := NewRequiredFlag(name, desc, fn)
	a.Default = def
	a.Required = true
	return a
}

func NewNoValueFlag(name string, desc string, required bool) *Flag {
	a := NewRequiredFlag(name, desc, nil)
	a.NoValue = true
	a.Required = required
	return a
}

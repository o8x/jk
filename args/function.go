package args

import "github.com/o8x/jk/v2/args/flag"

func New(app *App, args ...*flag.Flag) *Args {
	a := &Args{
		App: app,
	}

	a.Flags = append(a.Flags, args...)
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

package option

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/o8x/jk/v2/args/ini"

	"github.com/o8x/jk/v2/args/flag"
)

func Merge(list []*flag.Flag, flags ...*flag.Flag) []*flag.Flag {
	return append(list, flags...)
}

func Option(name string, def string, desc string, fn ...flag.HookFunc) *flag.Flag {
	a := Required(name, desc, fn...)
	a.Default = []string{def}
	a.Required = false
	return a
}

func Required(name string, desc string, fn ...flag.HookFunc) *flag.Flag {
	return &flag.Flag{
		Name:        []string{name},
		Description: desc,
		Required:    true,
		Env:         []string{strings.ToUpper(fmt.Sprintf("A%s", strings.TrimPrefix(name, "--")))},
		SingleValue: true,
		HookFuncs:   fn,
	}
}

func Default(name string, def string, desc string, fn ...flag.HookFunc) *flag.Flag {
	a := Required(name, desc, fn...)
	a.Default = []string{def}
	return a
}

func Defaults(name string, def []string, desc string, fn ...flag.HookFunc) *flag.Flag {
	a := Required(name, desc, fn...)
	a.Default = def
	a.Required = true
	return a
}

func NoValue(name string, desc string, fn ...flag.HookFunc) *flag.Flag {
	a := Required(name, desc, fn...)
	a.NoValue = true
	a.Required = false
	return a
}

func BindInt64(v *int64) flag.HookFunc {
	return func(f *flag.Flag) error {
		val, ok := f.GetInt64()
		if ok {
			*v = val
			return nil
		}

		return fmt.Errorf("unable to convert '%s' to int64", f.JoinName())
	}
}

func BindFloat64(v *float64) flag.HookFunc {
	return func(f *flag.Flag) error {
		val, ok := f.Get()
		if ok {
			float, err := strconv.ParseFloat(val, 0)
			if err != nil {
				return err
			}

			*v = float
			return nil
		}

		return fmt.Errorf("unable to convert '%s' to float64", f.JoinName())
	}
}

func BindBool(v *bool) flag.HookFunc {
	return func(f *flag.Flag) error {
		val, ok := f.GetBool()
		if ok {
			*v = val
			return nil
		}

		if !f.Required {
			return nil
		}

		return fmt.Errorf("unable to convert '%s' to bool", f.JoinName())
	}
}

func BindString(v *string) flag.HookFunc {
	return func(f *flag.Flag) error {
		val, ok := f.Get()
		if ok {
			*v = val
			return nil
		}

		return fmt.Errorf("unable to convert '%s' to string", f.JoinName())
	}
}

func BindIni(key string) flag.HookFunc {
	return func(f *flag.Flag) error {
		if ini.Get() == nil {
			return nil
		}

		s := ""
		k := ""

		split := strings.Split(key, ".")
		if len(split) == 1 {
			k = key
		} else {
			s = split[0]
			k = split[1]
		}

		f.Values = []string{ini.Get().Section(s).Key(k).String()}
		return nil
	}
}

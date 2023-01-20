package args

import (
	"encoding/base64"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/o8x/jk/signal"
)

type Properties map[string]any

func (p Properties) GetInt(name string) (int, bool) {
	v, ok := p.GetInt64(name)
	if !ok {
		return 0, false
	}

	return int(v), true
}

func (p Properties) GetInt64(name string) (int64, bool) {
	v, ok := p.Get(name)
	if !ok {
		return 0, false
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (p Properties) IsSet(name string) bool {
	_, ok := p[name]
	return ok
}

func (p Properties) Get(name string) (string, bool) {
	a, ok := p[name]
	if !ok {
		return "", false
	}

	s, ok := a.(string)
	if !ok {
		return "", false
	}

	return s, true
}

func (p Properties) GetBool(name string) (bool, bool) {
	v, ok := p.Get(name)
	if !ok {
		return false, false
	}

	if v == "true" {
		return true, true
	}

	if v == "false" {
		return false, true
	}

	return false, false
}

type Flag struct {
	Name             []string `json:"name"`
	Description      string   `json:"description"`
	PropertyMode     bool     `json:"property_mode"`
	Default          []string `json:"default"`
	Required         bool     `json:"required"`
	Env              []string `json:"env"`
	Error            error    `json:"error_message"`
	NoValue          bool     `json:"no_value"`
	ValuesOnlyInEnum []string `json:"value_only_in_enum"`
	SingleValue      bool     `json:"single_value"`
	values           []string
	properties       Properties
	exist            bool
}

func (a Flag) JoinName() string {
	return strings.Join(a.Name, "|")
}

func (a Flag) JoinDefault() string {
	return strings.Join(a.Default, ",")
}

func (a Flag) JoinEnum() string {
	return strings.Join(a.ValuesOnlyInEnum, "|")
}

type App struct {
	Name       string  `json:"name"`
	Usage      string  `json:"usage"`
	Copyright  string  `json:"copyright"`
	Version    *string `json:"version"`
	Changelog  *string `json:"changelog"`
	Banner     *string `json:"banner"`
	CommitHash *string `json:"hash"`
	Date       *string `json:"date"`
}

func (a App) AppFullVersion() string {
	b := strings.Builder{}
	b.WriteString(a.Name)
	if a.Version != nil {
		b.WriteString(fmt.Sprintf(" v%s", *a.Version))
	}
	b.WriteString(fmt.Sprintf(" (%s/%s) %s", runtime.GOOS, runtime.GOARCH, runtime.Version()))

	return b.String()
}

type Args struct {
	Executable string   `json:"executable"`
	App        *App     `json:"app"`
	Source     []string `json:"source"`
	Flags      []*Flag  `json:"args"`
	cmdline    string
	HelpFunc   func() string
}

func (a *Args) init() {
	if a.Source == nil {
		a.Executable = os.Args[0]
		a.Source = os.Args[1:]
	}

	a.Flags = append(a.Flags, &Flag{
		Name:        []string{"-help", "-h"},
		Description: "print this help and exit",
	})

	a.Flags = append(a.Flags, &Flag{
		Name:        []string{"-version", "-v"},
		Description: "print version info and exit",
	})
}

func isArgName(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func (a *Args) PrintVersionExit() {
	if a.App == nil {
		fmt.Println("no version info")
		os.Exit(1)
	}

	b := strings.Builder{}
	prepareDataPtr := func(v *string) {
		s := strings.TrimPrefix(*v, "base64://")
		if s != *v {
			if d, err := base64.StdEncoding.DecodeString(s); err == nil {
				*v = string(d)
			}
		}

		*v = strings.Trim(*v, "\n")
	}

	if a.App.Banner != nil {
		prepareDataPtr(a.App.Banner)
		b.WriteString(*a.App.Banner)
		b.WriteString("\n\n")
	}

	b.WriteString(a.App.AppFullVersion() + "\n\n")

	if a.App.Date != nil {
		prepareDataPtr(a.App.Date)
		b.WriteString(fmt.Sprintf("Release-Date: %s\n", *a.App.Date))
	}

	if a.App.CommitHash != nil {
		prepareDataPtr(a.App.CommitHash)
		b.WriteString(fmt.Sprintf("Commit Hash: %s\n", *a.App.CommitHash))
	}

	if a.App.Changelog != nil {
		prepareDataPtr(a.App.Changelog)

		b.WriteString("Changelog:")
		b.WriteString(fmt.Sprintf(" %s\n", *a.App.Changelog))
	}

	if a.App.Copyright != "" {
		b.WriteString(fmt.Sprintf("%s\n", a.App.Copyright))
	}

	fmt.Print(b.String())
	os.Exit(0)
}

func (a *Args) PrintHelpExit(err error) {
	fmt.Print(a.Help(err))
	if err == nil {
		os.Exit(0)
	}
	os.Exit(1)
}

func (a *Args) PrintErrorExit(err error) {
	fmt.Print(fmt.Sprintf("error: %v\n", err))
	os.Exit(1)
}

func (a *Args) Help(err error) string {
	if a.HelpFunc != nil {
		return a.HelpFunc()
	}

	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("Usage of %s\n", a.Executable))

	if err != nil {
		b.WriteString(fmt.Sprintf("error: %v\n", err))
		b.WriteString("\n")
	}

	if a.App != nil {
		b.WriteString(a.App.AppFullVersion())
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("usage: %s\n", a.App.Usage))
		b.WriteString("\n")
	}

	b.WriteString("commands: \n")
	for _, it := range a.Flags {
		b.WriteString(fmt.Sprintf("    %s", it.JoinName()))

		var properties []string
		if it.Required {
			properties = append(properties, "required")
		}

		if it.NoValue {
			properties = append(properties, "no value")
		}

		if it.SingleValue {
			properties = append(properties, "single")
		}

		if properties != nil {
			b.WriteString(": " + strings.Join(properties, ","))
		}
		b.WriteString("\n")

		b.WriteString(fmt.Sprintf("        %s", it.Description))
		if it.Default != nil {
			b.WriteString(fmt.Sprintf(" (default: %s)\n", it.JoinDefault()))
		} else {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (a *Args) Parse() error {
	a.init()
	a.cmdline = strings.Join(os.Args[1:], " ")

	for i := 0; i < len(a.Source); i++ {
		arg := a.Source[i]
		if arg == "-h" || arg == "-help" {
			a.PrintHelpExit(nil)
		}

		if arg == "-v" || arg == "-version" {
			a.PrintVersionExit()
		}

		if isArgName(arg) {
			k, v, found := strings.Cut(arg, "=")
			if found {
				a, err := a.findArg(k)
				if err != nil {
					return fmt.Errorf("flag provided but not defined: %s", arg)
				}

				a.exist = true
				a.values = append(a.values, v)
				continue
			}

			p, err := a.findArg(k)
			if err != nil {
				return err
			}

			p.exist = true
			if p.NoValue {
				continue
			}

			if i+1 >= len(a.Source) || strings.HasPrefix(a.Source[i+1], "-") {
				return fmt.Errorf("flag: %s need to provide a value", arg)
			}

			p.values = append(p.values, a.Source[i+1])
		}
	}

	for _, arg := range a.Flags {
		// 有默认值，并且没有传值
		if arg.values == nil && arg.Default != nil {
			arg.values = arg.Default
		}

		// 从默认值也没有获取到值，但是提供了环境变量名
		if arg.values == nil && arg.Env != nil {
			for _, name := range arg.Env {
				if val, found := os.LookupEnv(name); found {
					arg.values = append(arg.values, val)
				}
			}
		}

		// 必填但没有填并且也没有默认值
		if arg.Required && !arg.exist && arg.Default == nil {
			if arg.Error == nil {
				return fmt.Errorf("flag '%s' is required", arg.JoinName())
			}

			msg := strings.ReplaceAll(arg.Error.Error(), "{{name}}", arg.JoinName())
			if a.App != nil {
				msg = strings.ReplaceAll(msg, "{{app_name}}", a.App.Name)
				if a.App.Version != nil {
					msg = strings.ReplaceAll(msg, "{{app_version}}", *a.App.Version)
				}
			}
			return fmt.Errorf(msg)
		}

		if arg.SingleValue {
			if len(arg.values) > 1 || len(arg.values) == 0 {
				return fmt.Errorf("flag %s only allow one value", arg.JoinName())
			}
		}

		if arg.ValuesOnlyInEnum != nil {
			in := 0
			for _, v := range arg.values {
				for _, it := range arg.ValuesOnlyInEnum {
					if v == it {
						in++
					}
				}
			}

			if in != len(arg.values) {
				return fmt.Errorf("flag %s only one value in '%s' can be selected", arg.JoinName(), arg.JoinEnum())
			}
		}

		if arg.PropertyMode {
			arg.properties = Properties{}
			for _, value := range arg.values {
				k, v, found := strings.Cut(value, "=")
				if found {
					arg.properties[k] = v
					continue
				}

				arg.properties[value] = ""
			}
		}
	}

	return nil
}

func (a *Args) findArg(arg string) (*Flag, error) {
	if !strings.HasPrefix(arg, "-") {
		arg = "-" + arg
	}

	for _, it := range a.Flags {
		for _, name := range it.Name {
			if name == arg {
				return it, nil
			}
		}
	}

	return nil, fmt.Errorf("flag %s is not defined", arg)
}

func (a *Args) ParseCmdline(cmdline string) error {
	a.cmdline = cmdline
	a.Source = strings.Fields(cmdline)
	return a.Parse()
}

func (a *Args) IsSet(name string) bool {
	arg, err := a.findArg(name)
	if err != nil {
		return false
	}

	return arg.exist
}

func (a *Args) GetInt64(name string) (int64, bool) {
	arg, err := a.findArg(name)
	if err != nil {
		return 0, false
	}

	if arg.values == nil {
		return 0, false
	}

	i, err := strconv.ParseInt(arg.values[0], 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (a *Args) GetInt(name string) (int, bool) {
	v, ok := a.GetInt64(name)
	return int(v), ok
}

func (a *Args) GetInts(name string) []int {
	s, err := a.GetInt64s(name)
	if err != nil {
		return nil
	}

	var list []int
	for _, it := range s {
		list = append(list, int(it))
	}

	return list
}

func (a *Args) GetInt64s(name string) ([]int64, error) {
	arg, err := a.findArg(name)
	if err != nil {
		return nil, err
	}

	var result []int64
	for _, v := range arg.values {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}

		result = append(result, i)
	}

	return result, nil
}

func (a *Args) GetProperties(name string) Properties {
	arg, err := a.findArg(name)
	if err != nil || !arg.PropertyMode {
		return nil
	}

	return arg.properties
}

func (a *Args) GetProperty(name string, property string) (string, bool) {
	properties := a.GetProperties(name)
	if properties != nil {
		return properties.Get(property)
	}
	return "", false
}

func (a *Args) GetPropertyX(name string, property string) string {
	properties := a.GetProperties(name)
	if properties != nil {
		v, ok := properties.Get(property)
		if ok {
			return v
		}
	}

	panic(fmt.Sprintf("property %s.%s not found", name, property))
}

func (a *Args) Get(name string) (string, bool) {
	arg, err := a.findArg(name)
	if err != nil {
		return "", false
	}

	if arg.values == nil {
		return "", false
	}

	return arg.values[0], true
}

func (a *Args) GetX(name string) string {
	arg, err := a.findArg(name)
	if err != nil {
		panic(err)
	}

	if arg.values == nil {
		panic(fmt.Errorf("flag %s values is nil", name))
	}

	return arg.values[0]
}

func (a *Args) Gets(name string) []string {
	arg, err := a.findArg(name)
	if err != nil {
		return nil
	}

	return arg.values
}

func (a *Args) GetBool(name string) (bool, bool) {
	arg, err := a.findArg(name)
	if err != nil {
		return false, false
	}

	if arg.values == nil {
		return false, false
	}

	v := arg.values[0]
	if v == "true" {
		return true, true
	}

	if v == "false" {
		return false, true
	}

	return false, false
}

func (a *Args) Exit(code int) {
	os.Exit(code)
}

func (a *Args) WaitSignal() os.Signal {
	return signal.Wait()
}

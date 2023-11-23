package args

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/o8x/jk/v2/args/flag"
	"github.com/o8x/jk/v2/signal"
)

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
	Executable string       `json:"executable"`
	App        *App         `json:"app"`
	Source     []string     `json:"source"`
	Flags      []*flag.Flag `json:"args"`
	cmdline    string
	HelpFunc   func() string
	cacheMap   map[string]any
}

func (a *Args) init() {
	if a.Source == nil {
		a.Executable = os.Args[0]
		a.Source = os.Args[1:]
	}

	a.Flags = append(a.Flags, &flag.Flag{
		Name:        []string{"-help", "-h"},
		Description: "print this help and exit",
	})

	a.Flags = append(a.Flags, &flag.Flag{
		Name:        []string{"-version", "-v"},
		Description: "print version info and exit",
	})

	a.cacheMap = map[string]any{}
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

		s = strings.TrimPrefix(*v, "hex://")
		if s != *v {
			if d, err := hex.DecodeString(s); err == nil {
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

func (a *Args) DumpExit() {
	b := strings.Builder{}
	if err := a.Parse(); err != nil {
		b.WriteString(fmt.Sprintf("error: %s", err))
		b.WriteString("\n\n")
	}

	b.WriteString("commands: ")
	b.WriteString("\n")

	for _, arg := range a.Flags {
		b.WriteString(fmt.Sprintf("  %s", strings.Join(arg.Name, "|")))
		b.WriteString("\n")

		if arg.Values != nil {
			b.WriteString(fmt.Sprintf("\tvalues: %s (from cmdline)", strings.Join(arg.Values, ", ")))
			b.WriteString("\n")
		}

		if arg.Values == nil && arg.Default != nil {
			arg.Values = arg.Default
			b.WriteString(fmt.Sprintf("\tvalues: %s (use default)", strings.Join(arg.Values, ", ")))
			b.WriteString("\n")
		}

		if arg.Values == nil && arg.Env != nil {
			for _, name := range arg.Env {
				if val, found := os.LookupEnv(name); found {
					arg.Values = append(arg.Values, val)
				}
			}

			b.WriteString(fmt.Sprintf("\tvalues: %s (use environment)", strings.Join(arg.Values, ", ")))
			b.WriteString("\n")
		}

		// 必填但没有填并且也没有默认值
		if arg.Required && !arg.Exist && arg.Default == nil {
			b.WriteString("\terror: required, but no provide value.")
			b.WriteString("\n")
		}

		if arg.SingleValue {
			if len(arg.Values) > 1 || len(arg.Values) == 0 {
				b.WriteString(fmt.Sprintf("\terror: only one value allowed, provide values: %s", strings.Join(arg.Values, ", ")))
				b.WriteString("\n")
			}
		}

		if arg.HookFuncs != nil {
			for i, fn := range arg.HookFuncs {
				if err := fn(arg); err != nil {
					b.WriteString(fmt.Sprintf("%s hook %d error: %v", arg.Name[0], i, err))
					b.WriteString("\n")
				}
			}
		}
	}

	fmt.Println(b.String())
	a.Exit(0)
}

func (a *Args) PrintHelpExit(err error) {
	fmt.Print(a.Help(err))
	if err == nil {
		os.Exit(0)
	}
	os.Exit(1)
}

func (a *Args) PrintErrorExit(err error) {
	fmt.Printf("error: %v\n", err)
	os.Exit(1)
}

func (a *Args) GenerateUsage() string {
	var s []string
	var m = map[string]struct{}{}
	for _, it := range a.Flags {
		if !it.Required {
			continue
		}

		if it.NoValue {
			s = append(s, it.Name[0])
			m[it.Name[0]] = struct{}{}
			continue
		}

		if it.Default != nil {
			for _, d := range it.Default {
				// 显式声明空字符串默认值
				if d == "" {
					d = `""`
				}

				s = append(s, fmt.Sprintf("%s %s", it.Name[0], d))
				m[it.Name[0]] = struct{}{}
			}
			continue
		}

		if it.Env != nil {
			for _, e := range it.Env {
				val, ok := os.LookupEnv(e)
				if !ok {
					continue
				}

				m[it.Name[0]] = struct{}{}
				s = append(s, fmt.Sprintf("%s ${%s:=%s}", it.Name[0], e, val))
			}
		}

		if _, ok := m[it.Name[0]]; !ok {
			s = append(s, fmt.Sprintf(`%s ""`, it.Name[0]))
		}
	}

	return strings.Join(s, " ")
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
		a.App.Usage = strings.ReplaceAll(a.App.Usage, "{{auto}}", a.GenerateUsage())
		a.App.Usage = strings.ReplaceAll(a.App.Usage, "{{executable}}", a.Executable)

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

				a.Exist = true
				a.Values = append(a.Values, v)
				continue
			}

			p, err := a.findArg(k)
			if err != nil {
				return err
			}

			p.Exist = true
			if p.NoValue {
				continue
			}

			if i+1 >= len(a.Source) || strings.HasPrefix(a.Source[i+1], "-") {
				if p.Default != nil {
					p.Values = p.Default
					continue
				}

				return fmt.Errorf("flag: %s need to provide a value", arg)
			}

			p.Values = append(p.Values, a.Source[i+1])
		}
	}

	for _, arg := range a.Flags {
		// 有默认值，并且没有传值
		if arg.Values == nil && arg.Default != nil {
			arg.Values = arg.Default
		}

		// 从默认值也没有获取到值，但是提供了环境变量名
		if arg.Values == nil && arg.Env != nil {
			for _, name := range arg.Env {
				if val, found := os.LookupEnv(name); found {
					arg.Values = append(arg.Values, val)
				}
			}
		}

		// 必填但没有填并且也没有默认值
		if arg.Required && !arg.Exist && arg.Default == nil {
			return fmt.Errorf("flag '%s' is required", arg.JoinName())
		}

		if arg.SingleValue && !arg.NoValue {
			if len(arg.Values) > 1 || len(arg.Values) == 0 {
				return fmt.Errorf("flag %s only allow one value", arg.JoinName())
			}
		}

		// 生成缓存
		for _, n := range arg.Name {
			if len(arg.Values) == 1 {
				a.cacheMap[n] = arg.Values[0]
			} else {
				a.cacheMap[n] = arg.Values
			}
		}

		if arg.HookFuncs != nil {
			for _, fn := range arg.HookFuncs {
				if err := fn(arg); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (a *Args) findArg(arg string) (*flag.Flag, error) {
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

func (a *Args) Bytes() []byte {
	marshal, _ := json.Marshal(a.cacheMap)
	return marshal
}

func (a *Args) Copy() Readonly {
	return Readonly{
		parent: *a,
	}
}

func (a *Args) Bind(v any) error {
	if a.cacheMap == nil {
		if err := a.Parse(); err != nil {
			return err
		}
	}

	if v == nil {
		return nil
	}

	d := Readonly{
		parent: *a,
	}
	return d.Unmarshal(v)
}

func (a *Args) Exit(code int) {
	os.Exit(code)
}

func (a *Args) WaitSignal() os.Signal {
	return signal.Wait()
}

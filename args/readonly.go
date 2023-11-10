package args

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Readonly struct {
	parent Args
}

func (a *Readonly) findArg(arg string) (*Flag, error) {
	return a.parent.findArg(arg)
}

func (a *Readonly) Unmarshal(v any) error {
	return json.Unmarshal(a.parent.Bytes(), v)
}

func (a *Args) ParseCmdline(cmdline string) error {
	a.cmdline = cmdline
	a.Source = strings.Fields(cmdline)
	return a.Parse()
}

func (a *Readonly) IsSet(name string) bool {
	arg, err := a.findArg(name)
	if err != nil {
		return false
	}

	return arg.exist
}

func (a *Readonly) GetInt64(name string) (int64, bool) {
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

func (a *Readonly) GetInt(name string) (int, bool) {
	v, ok := a.GetInt64(name)
	return int(v), ok
}

func (a *Readonly) GetInts(name string) []int {
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

func (a *Readonly) GetInt64s(name string) ([]int64, error) {
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

func (a *Readonly) GetProperties(name string) Properties {
	arg, err := a.findArg(name)
	if err != nil || !arg.PropertyMode {
		return nil
	}

	return arg.properties
}

func (a *Readonly) GetProperty(name string, property string) (string, bool) {
	properties := a.GetProperties(name)
	if properties != nil {
		return properties.Get(property)
	}
	return "", false
}

func (a *Readonly) GetPropertyX(name string, property string) string {
	properties := a.GetProperties(name)
	if properties != nil {
		v, ok := properties.Get(property)
		if ok {
			return v
		}
	}

	panic(fmt.Sprintf("property %s.%s not found", name, property))
}

func (a *Readonly) Get(name string) (string, bool) {
	arg, err := a.findArg(name)
	if err != nil {
		return "", false
	}

	if arg.values == nil {
		return "", arg.NoValue
	}

	return arg.values[0], true
}

func (a *Readonly) GetX(name string) string {
	arg, err := a.findArg(name)
	if err != nil {
		panic(err)
	}

	if arg.values == nil {
		panic(fmt.Errorf("flag %s values is nil", name))
	}

	return arg.values[0]
}

func (a *Readonly) Gets(name string) []string {
	arg, err := a.findArg(name)
	if err != nil {
		return nil
	}

	return arg.values
}

func (a *Readonly) GetBool(name string) (bool, bool) {
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

type Properties map[string]any

func (p Properties) GetInt(name string) (int, bool) {
	v, ok := p.GetInt64(name)
	if !ok {
		return 0, false
	}

	return int(v), true
}

func (p Properties) GetIntDefault(name string, def int) int {
	get, ok := p.GetInt(name)
	if ok {
		return get
	}

	return def
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

func (p Properties) GetInt64Default(name string, def int64) int64 {
	get, ok := p.GetInt64(name)
	if ok {
		return get
	}

	return def
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

func (p Properties) GetDefault(name, def string) string {
	get, ok := p.Get(name)
	if ok {
		return get
	}

	return def
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

func (p Properties) GetBoolDefault(name string, def bool) bool {
	get, ok := p.GetBool(name)
	if ok {
		return get
	}

	return def
}

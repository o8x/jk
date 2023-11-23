package args

import (
	"fmt"
	"strconv"
	"strings"
)

type HookFunc func(int, []string) error

type Flag struct {
	Name        []string `json:"name"`
	Description string   `json:"description"`
	Default     []string `json:"default"`
	Required    bool     `json:"required"`
	Env         []string `json:"env"`
	NoValue     bool     `json:"no_value"`
	SingleValue bool     `json:"single_value"`
	HookFunc    HookFunc
	values      []string
	exist       bool
}

func (a *Flag) JoinName() string {
	return strings.Join(a.Name, "|")
}

func (a *Flag) JoinDefault() string {
	return strings.Join(a.Default, ",")
}

func (a *Flag) BindInt64(v *int64) *Flag {
	a.HookFunc = func(i int, i2 []string) error {
		val, ok := a.GetInt64()
		if ok {
			*v = val
			return nil
		}

		return fmt.Errorf("unable to convert '%s' to int64", a.JoinName())
	}

	return a
}

func (a *Flag) BindBool(v *bool) *Flag {
	a.HookFunc = func(i int, i2 []string) error {
		val, ok := a.GetBool()
		if ok {
			*v = val
			return nil
		}

		if !a.Required {
			return nil
		}

		return fmt.Errorf("unable to convert '%s' to bool", a.JoinName())
	}

	return a
}

func (a *Flag) BindString(v *string) *Flag {
	a.HookFunc = func(i int, i2 []string) error {
		val, ok := a.Get()
		if ok {
			*v = val
			return nil
		}

		return fmt.Errorf("unable to convert '%s' to string", a.JoinName())
	}

	return a
}

func (a *Flag) GetInt64() (int64, bool) {
	if a.values == nil {
		return 0, false
	}

	i, err := strconv.ParseInt(a.values[0], 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (a *Flag) GetInt() (int, bool) {
	v, ok := a.GetInt64()
	return int(v), ok
}

func (a *Flag) GetInts() []int {
	s, err := a.GetInt64s()
	if err != nil {
		return nil
	}

	var list []int
	for _, it := range s {
		list = append(list, int(it))
	}

	return list
}

func (a *Flag) GetInt64s() ([]int64, error) {
	var result []int64
	for _, v := range a.values {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}

		result = append(result, i)
	}

	return result, nil
}

func (a *Flag) Get() (string, bool) {
	if a.values == nil {
		return "", a.NoValue
	}

	return a.values[0], true
}

func (a *Flag) GetX() string {
	if a.values == nil {
		panic(fmt.Errorf("flag %s values is nil", a.JoinName()))
	}

	return a.values[0]
}

func (a *Flag) Gets() []string {
	return a.values
}

func (a *Flag) GetBool() (bool, bool) {
	if a.values == nil {
		return false, false
	}

	// 存在即为 true
	if a.NoValue {
		return true, true
	}

	v := a.values[0]
	if v == "true" {
		return true, true
	}

	if v == "false" {
		return false, true
	}

	return false, false
}

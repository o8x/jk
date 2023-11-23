package flag

import (
	"fmt"
	"strconv"
	"strings"
)

type HookFunc func(*Flag) error

type Flag struct {
	Name        []string `json:"name"`
	Description string   `json:"description"`
	Default     []string `json:"default"`
	Required    bool     `json:"required"`
	Env         []string `json:"env"`
	NoValue     bool     `json:"no_value"`
	SingleValue bool     `json:"single_value"`
	HookFunc    HookFunc
	Values      []string
	Exist       bool
}

func (a *Flag) ValsLen() int {
	return len(a.Values)
}

func (a *Flag) JoinName() string {
	return strings.Join(a.Name, "|")
}

func (a *Flag) JoinDefault() string {
	return strings.Join(a.Default, ",")
}

func (a *Flag) GetInt64() (int64, bool) {
	if a.Values == nil {
		return 0, false
	}

	i, err := strconv.ParseInt(a.Values[0], 10, 64)
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
	for _, v := range a.Values {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}

		result = append(result, i)
	}

	return result, nil
}

func (a *Flag) Get() (string, bool) {
	if a.Values == nil {
		return "", a.NoValue
	}

	return a.Values[0], true
}

func (a *Flag) GetX() string {
	if a.Values == nil {
		panic(fmt.Errorf("flag %s Values is nil", a.JoinName()))
	}

	return a.Values[0]
}

func (a *Flag) Gets() []string {
	return a.Values
}

func (a *Flag) GetBool() (bool, bool) {
	if a.Values == nil {
		return false, false
	}

	// 存在即为 true
	if a.NoValue {
		return true, true
	}

	v := a.Values[0]
	if v == "true" {
		return true, true
	}

	if v == "false" {
		return false, true
	}

	return false, false
}

package args

import (
	"encoding/json"
	"strconv"

	"github.com/o8x/jk/v2/args/flag"
)

type Readonly struct {
	parent Args
}

func (a *Readonly) findArg(arg string) (*flag.Flag, error) {
	return a.parent.findArg(arg)
}

func (a *Readonly) Unmarshal(v any) error {
	return json.Unmarshal(a.parent.Bytes(), v)
}

func (a *Readonly) IsSet(name string) bool {
	arg, err := a.findArg(name)
	if err != nil {
		return false
	}

	return arg.Exist
}

func (a *Readonly) GetInt64(name string) (int64, bool) {
	arg, err := a.findArg(name)
	if err != nil {
		return 0, false
	}

	if arg.Values == nil {
		return 0, false
	}

	i, err := strconv.ParseInt(arg.Values[0], 10, 64)
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
	for _, v := range arg.Values {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}

		result = append(result, i)
	}

	return result, nil
}

func (a *Readonly) Get(name string) (string, bool) {
	arg, err := a.findArg(name)
	if err != nil {
		return "", false
	}

	return arg.Get()
}

func (a *Readonly) GetX(name string) string {
	arg, err := a.findArg(name)
	if err != nil {
		panic(err)
	}

	return arg.GetX()
}

func (a *Readonly) Gets(name string) []string {
	arg, err := a.findArg(name)
	if err != nil {
		return nil
	}

	return arg.Gets()
}

func (a *Readonly) GetBool(name string) (bool, bool) {
	arg, err := a.findArg(name)
	if err != nil {
		return false, false
	}

	return arg.GetBool()
}

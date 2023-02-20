package json

import (
	"bytes"
	"encoding/json"
)

type Type interface {
	string | []byte
}

func Unmarshal[T Type](data T, v any) error {
	return json.Unmarshal([]byte(data), v)
}

func Marshal[T Type](a any) T {
	marshal, _ := json.Marshal(a)
	return T(marshal)
}

func Prettify[T Type](v any) T {
	marshal := Marshal[[]byte](v)
	buf := bytes.NewBuffer(nil)
	if err := json.Indent(buf, marshal, "", "    "); err != nil {
		return T(marshal)
	}

	return T(buf.Bytes())
}

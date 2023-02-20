Manual
===

## API

### 泛型 Type

    type Type interface {
        string | []byte
    }

### json 解码

输入类型支持泛型，golang 可自动推导，无需显式声明

    func Unmarshal[T Type](data T, v any) error

### json 编码

忽略了错误的 json 编码器，输出类型支持泛型

    func Marshal[T Type](a any) T

### 美化 json

忽略了错误的 json 编码器，并自动对返回值进行美化，输出类型支持泛型

    func Prettify[T Type](v any) T

## 示例

```go
package main

import (
	"fmt"
	"time"

	"github.com/o8x/jk/json"
)

type d struct {
	Name     string `json:"name"`
	DateTime string `json:"date_time"`
}

func main() {
	var data = d{
		Name:     "Json Package Test",
		DateTime: time.Now().Format(time.DateTime),
	}

	jsonStr := json.Marshal[string](data)
	jsonPrettifyString := json.Prettify[string](data)
	jsonBytes := json.Marshal[[]byte](data)
	jsonPrettifyBytes := json.Prettify[[]byte](data)

	fmt.Println(jsonStr, "\n", jsonPrettifyString, "\n", jsonBytes, "\n", jsonPrettifyBytes)

	var d1 d
	if err := json.Unmarshal(jsonStr, &d1); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(jsonBytes, &d1); err != nil {
		panic(err)
	}

	fmt.Println(d1.Name, d1.DateTime)
}
```

运行它

```shell
> go run .
{"name":"Json Package Test","date_time":"2023-02-20 17:20:07"} 
 {
    "name": "Json Package Test",
    "date_time": "2023-02-20 17:20:07"
} 
 [123 34 ...] 
 [123 10 ...]
Json Package Test 2023-02-20 17:20:07
```

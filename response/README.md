## API

## Response

```go
type Response struct {
    StatusCode int    `json:"status_code"`
    Message    string `json:"message,omitempty"`
    Body       any    `json:"body,omitempty"`
}
```

### 判断当前是否为异常响应

    func (r Response) IsError() (string, bool)

### 判断当前是否为正常响应

    func (r Response) IsNormal() bool

### 将当前响应体格式化为 json

    func (r Response) Dump() []byte

## API

### 生成正常响应

此时 status_code 为 200，message 无值

    func OK(body any) *Response

### 生成警告响应

此时 status_code 为 500，message 有值

    func Warn(msg string) *Response

### 生成错误响应

status_code 为 500，message 有值

    func Error(err error) *Response

### 生成响应

    func Build(code int, msg string, body any) *Response

### NoContent

**变量，不是方法**

status_code 为 204，message 无值

### BadRequest

**变量，不是方法**

status_code 为 400，message 无值


## 示例

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/o8x/jk/response"
)

func main() {
	fmt.Println(string(response.OK([]string{"hello"}).Dump()))
	fmt.Println(string(response.Warn("warn").Dump()))
	fmt.Println(string(response.Error(fmt.Errorf("error")).Dump()))
	fmt.Println(string(response.Build(http.StatusOK, "", "hello").Dump()))

	fmt.Println(string(response.NoContent.Dump()))
	fmt.Println(response.NoContent.IsError())
	fmt.Println(response.NoContent.IsNormal())

	fmt.Println(string(response.BadRequest.Dump()))
	fmt.Println(response.BadRequest.IsError())
	fmt.Println(response.BadRequest.IsNormal())
} 
```

运行它
```shell
> go run .
{"status_code":200,"body":["hello"]}
{"status_code":500,"message":"warn"}
{"status_code":500,"message":"error"}
{"status_code":200,"body":"hello"}
{"status_code":204}
 false
true
{"status_code":400}
 false
false 
```


Manual
===

捕获运行时 panic

## API

### 捕获 panic 并打印日志

    func Recover(message string)

## 示例

```go
package main

import (
	"github.com/o8x/jk/crash"
	"github.com/o8x/jk/signal"
)

func foo() {
	panic("panic in foo")
}

func main() {
	go func() {
		defer crash.Recover("recover by main")

		foo()
	}()

	signal.Wait()
}
```

运行它

```shell
> go run .
ERRO[0000] recover by main       recover="panic in foo" stack="goroutine .. [running]:....."
```

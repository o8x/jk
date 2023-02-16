Manual
===

捕获运行时 panic

## API

### Crash

    type Crash struct {
        Logger  *logrus.Logger `json:"logger"`
        Message string         `json:"message"`
    }

### RecoverFunc

    type RecoverFunc func(*Crash)

### 捕获 panic 并打印日志

    func Recover(message string, fns ...RecoverFunc)

## 示例

```go
package main

import (
	"github.com/sirupsen/logrus"

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

	go func() {
		l := logrus.New()
		l.SetLevel(logrus.FatalLevel)
		defer crash.Recover("crash", crash.Logger(l))

		foo()
	}()

	signal.Wait()
}
```

运行它

第二次注入的 logger 是 fatal 级别，所以 Error 级别的日志不会被打印

```shell
> go run .
ERRO[0000] recover by main       recover="panic in foo" stack="goroutine .. [running]:....."
```

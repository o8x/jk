Manual
===

## 阻塞主 goroutine

会阻塞主 goroutine 直到接收到 USR1、USR2、INT、TERM、QUIT 其中之一的信号

    func Wait() os.Signal

示例

```go
package main

import (
	"fmt"

	"github.com/o8x/jk/signal"
)

func main() {
	signal.Wait()
	
	fmt.Println("app shutdown")
}
```

## 阻塞主 goroutine

会阻塞主 goroutine 直到接收到提供的信号

    func WaitSignal(s ...os.Signal) os.Signal

示例

```go
package main

import (
	"fmt"
	"syscall"

	"github.com/o8x/jk/signal"
)

func main() {
	signal.WaitSignal(syscall.SIGINT)

	fmt.Println("app shutdown")
}
```

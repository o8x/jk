Manual
===

简单的 udp server

## API

### 监听并阻塞

    func ListenAndServe(ctx context.Context, addr string, fn func([]byte, net.Addr, error)) error 

### 向 udp server 写入一个包

    func WritePacket(addr string, packet []byte) (int, error)

## 示例

```go
package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/o8x/jk/logger"
	"github.com/o8x/jk/udp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * 5)
		logger.Info("canceled")
		cancel()
	}()

	err := udp.ListenAndServe(ctx, ":64458", func(bytes []byte, addr net.Addr, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(addr.String(), string(bytes))
		logger.Info("released")
	})

	if err != nil {
		panic(err)
	}

	logger.Info("application exit")
}
```

运行它

```shell
> go run .
INFO[0000] udp server listen on udp://[::]:64458        
127.0.0.1:53966 Hello

INFO[0002] released
INFO[0005] canceled
INFO[0005] udp server closed
INFO[0005] application exit
```

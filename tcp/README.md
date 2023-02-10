Manual
===

简单的 tcp server

## API

### 监听并阻塞

    func ListenAndServe(ctx context.Context, addr string, fn func(net.Conn, error)) error

### 向 tcp server 写入一个包

    func WritePacket(addr string, packet []byte) (int, error)

## 示例

```go
package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/net/context"

	"github.com/o8x/jk/logger"
	"github.com/o8x/jk/signal"
	"github.com/o8x/jk/tcp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(time.Second * 5)
		logger.Info("canceled")
		cancel()
	}()

	err := tcp.ListenAndServe(ctx, ":64458", func(conn net.Conn, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}

		all, err := io.ReadAll(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(conn.RemoteAddr(), string(all))
		logger.Info("connect released")
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
INFO[0000] tcp server listen on tcp://[::]:64458        
INFO[0005] canceled                                     
INFO[0005] tcp server closed                            
127.0.0.1:65037 Hello

INFO[0011] connect released                             
INFO[0011] application exit  
```

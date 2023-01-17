Manual
===

## 初始化 sqlite

    func Init(file string) error

## 获取 sqlite 实例

    func Get() *sql.DB

## 获取 sqlite 实例

*该方法并不安全，可能引发 panic*

生成一个在内存中的 sqlite 实例

    func Default() *sql.DB

## 示例

```go
package main

import (
	"github.com/o8x/jk/logger"
	"github.com/o8x/jk/sqlite"
)

func main() {
	err := sqlite.Init("file:s?mode=memory")
	if err != nil {
		logger.Get().WithError(err).Fatal("init sqlite error")
	}

	logger.Get().Info("init success")
}
```

运行它

```shell
> go run .
time="2023-01-17T11:40:37+08:00" level=info msg="init success"
```

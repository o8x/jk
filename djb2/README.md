## API

### 生成 djb2 Hash（uint64）

    func Sum(s string) uint64

### 生成 djb2 Hash（字符串）

    func Make(s string) string

### 校验 djb2 Hash

    func Check(s string, hash string) bool

## 示例

```go 
package main

import (
	"fmt"

	"github.com/o8x/jk/djb2"
)

func main() {
	sum := djb2.Make("hello world")
	
	fmt.Println("sum", sum)
	fmt.Println("check sum", djb2.Check("hello world", sum))
}

```

运行它

```shell
> go run .
sum 13876786532495509697
check sum true
```

Manual
===

## API

### 判断文件是否存在

    func FileExist(name string) bool

### 判断对文件是否具有读权限

    func HasRDPermission(name string) bool

### 判断对文件是否具有写权限

    func HasWRPermission(name string) bool

### 判断对文件是否具有读写权限

    func HasRWPermission(name string) bool

## 示例

```go
package main

import (
	"fmt"
	"os"
	"path"

	"github.com/o8x/jk/fs"
)

func main() {
	dir, _ := os.UserHomeDir()
	file := path.Join(dir, ".bashrc")

	fmt.Println("file name", file)
	if fs.FileExist(file) {
		fmt.Println("file exists")
	}
	fmt.Println("\thas read permission", fs.HasRDPermission(file))
	fmt.Println("\thas write permission", fs.HasWRPermission(file))
	fmt.Println("\thas read/write permission", fs.HasRWPermission(file))
	fmt.Println()

	file = "/etc/hosts"
	fmt.Println("file name", file)
	if fs.FileExist(file) {
		fmt.Println("file exists")
	}
	fmt.Println("\thas read permission", fs.HasRDPermission(file))
	fmt.Println("\thas write permission", fs.HasWRPermission(file))
	fmt.Println("\thas read/write permission", fs.HasRWPermission(file))
}
```

运行它

```shell
> go run .
file name /Users/alex/.bashrc
file exists
        has read permission true
        has write permission true
        has read/write permission true

file name /etc/hosts
file exists
        has read permission true
        has write permission false
        has read/write permission false
```

Manual
===

支持泛型的 sync.Map

## 创建一个 map 实例

    func New[Tk any, Tv any]() *Map[Tk, Tv] 

## API

用法等同 sync.Map 下同名方法，不赘述

```go     
func (r *Map[Tk, Tv]) Store(name Tk, value Tv)
func (r *Map[Tk, Tv]) LoadOrStore(name Tk, value Tv) (v Tv, ok bool)
func (r *Map[Tk, Tv]) Load(name Tk) (v Tv)
func (r *Map[Tk, Tv]) Range(fn func (Tk, Tv) bool)
func (r *Map[Tk, Tv]) LoadAndDelete(name Tk) (v Tv)
func (r *Map[Tk, Tv]) Delete(name Tk)
```

### 判断 Key 是否存在

补充方法，本质上是判断是否能进行 Load

    func (r *Map[Tk, Tv]) Exist(name Tk) bool

## 示例

```go
package main

import (
	"fmt"

	"github.com/o8x/jk/syncmap"
)

func main() {
	m := syncmap.New[string, int]()
	m.Store("k", 1)
	fmt.Println(m.Load("k"))

	m.Store("k", 9)
	fmt.Println(m.Load("k"))

	m.Store("k2", 2)
	fmt.Println(m.Load("k2"))

	m.LoadOrStore("k3", 3)
	fmt.Println(m.Load("k3"))
	m.LoadOrStore("k3", 5)
	fmt.Println(m.Load("k3"))

	fmt.Println(m.Exist("k3"))
	m.Delete("k3")
	fmt.Println(m.Exist("k3"))
}
```

运行它

```shell
> go run .
1
9
2
3
3
true
false
```

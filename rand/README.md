Manual
===

## 随机数

是 std rand 包同名方法的 alias

    func Intn(n int) int
    func Int() int
    func Int63() int64
    func Int63n(n int64) int64

## 加权随机

### 随机种子类型泛型

    type SeedType interface {
        any
    }

### 随机种子项

    type Item[T SeedType] struct {
        item   T
        weight int
    }

### 加权随机实现

    type WeightRand[T SeedType] struct {
        seed   []Item[T]
        weight int
        lock   *sync.Mutex
    }

### 初始化 WeightRand

如果提供了 items，则会自动执行 Add(item)

    func NewWeightRand[T SeedType](items ...T) *WeightRand[T]

### 增加一个随机种子

本质上是 w.AddWeight(it, 1)

    func (w *WeightRand[T]) Add(it T)

### 增加一个随机种子，并附带权重

每次增加种子都会将其 weight 值累加到 w.weight 中

    func (w *WeightRand[T]) AddWeight(it T, weight int)

### 从提供的 item 中，随机获取一个

将会按加入时的权重进行随机取值，算法为从 0 到 w.weight 之间取出随机数 R，遍历所有 item 并累加其 weight，当 item.weight 的和大于 R 时返回对应的 item

    func (w *WeightRand[T]) Get() (t T)

## 示例

```go
package main

import (
	"fmt"

	"github.com/o8x/jk/v2/json"
	"github.com/o8x/jk/v2/rand"
)

func main() {
	w := rand.NewWeightRand[int](0, 1, 2, 3)
	w.Add(4)
	w.AddWeight(5, 2)

	m := map[int]int{}
	for i := 0; i < 1000000; i++ {
		m[w.Get()]++
	}

	fmt.Println(json.Prettify[string](m))
}
```

运行它

```shell
> go run .
{
    "0": 142731,
    "1": 142644,
    "2": 142807,
    "3": 142435,
    "4": 143444,
    "5": 285939
}
```

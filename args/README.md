Args Manual
===

无外部依赖的命令行参数解析器

<!-- vscode-markdown-toc -->
- [Args Manual](#args-manual)
	- [ 开始](#-开始)
	- [示例](#示例)
		- [Flags](#flags)
		- [未定义](#未定义)
	- [Flag](#flag)
		- [Name](#name)
		- [Description](#description)
		- [PropertyMode](#propertymode)
		- [SingleValue](#singlevalue)
		- [Default](#default)
		- [Required](#required)
		- [Env](#env)
		- [Error](#error)
		- [NoValue](#novalue)
		- [ValuesOnlyInEnum](#valuesonlyinenum)
		- [优先级](#优先级)
	- [取值](#取值)
		- [判断 Flag 是否存在](#判断-flag-是否存在)
		- [获取 Int64 Flag](#获取-int64-flag)
		- [获取 Int64 数组 Flag](#获取-int64-数组-flag)
		- [获取 Int Flag](#获取-int-flag)
		- [获取 Int 数组 Flag](#获取-int-数组-flag)
		- [获取 字符串 Flag](#获取-字符串-flag)
		- [获取不到有效 字符串 Flag 时 panic](#获取不到有效-字符串-flag-时-panic)
		- [获取 字符串 数组 Flag](#获取-字符串-数组-flag)
		- [获取 Bool Flag](#获取-bool-flag)
	- [ Property 模式](#-property-模式)
		- [获取 Property Flag 的 Properties](#获取-property-flag-的-properties)
		- [获取 Flag 的 Property 的值](#获取-flag-的-property-的值)
		- [在获取不到有效 Property 时 panic](#在获取不到有效-property-时-panic)
		- [取值方法](#取值方法)
	- [帮助](#帮助)
	- [版本](#版本)
		- [注入版本信息](#注入版本信息)
	- [其他 API](#其他-api)
		- [ 阻塞主 goroutine](#-阻塞主-goroutine)
		- [ 退出](#-退出)
		- [打印调试日志并退出](#打印调试日志并退出)
	- [脚手架](#脚手架)
		- [ArgFunc](#argfunc)
		- [生成基本的 Args 实例](#生成基本的-args-实例)
		- [注入应用程序名称和使用示例](#注入应用程序名称和使用示例)
		- [增加基本的 Flag](#增加基本的-flag)
		- [增加带有默认值的 Flag](#增加带有默认值的-flag)
		- [增加 Enum Flag](#增加-enum-flag)
		- [增加没有值的 Flag](#增加没有值的-flag)
		- [手动扩展 Option](#手动扩展-option)
		- [调用示例](#调用示例)

<!-- vscode-markdown-toc-config
	numbering=false
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

无外部依赖的命令行参数解析器

## <a name=''></a> 开始

此应用程序将正确运行并显示 Hello World.

```go
package main

import (
	"fmt"

	"github.com/o8x/jk/args"
)

func main() {
	a := args.Args{}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}

	fmt.Println("Hello World.")
}
```

也可以使用既定的 cmdline 进行解析，此命令的原理是使用 strings.Fields 对 cmdline 实参进行解析，得到的结果类似 os.Args

```go
if err := a.ParseCmdline("-h -v"); err != nil {
	a.PrintErrorExit(err)
}
```

在错误时同时显示帮助文本

```go
if err := a.Parse(); err != nil {
	a.PrintHelpExit(err)
}
```

解析 os.Args 可能不像你想象的那么简单，但是也不会像你想象的那么复杂。有时候，可能只需要一个纯粹的 `flag.Parse()`，但更多的时候我们并不喜欢这样的虽然纯粹但无比复杂的方式。

Args 提供了大量的功能来让解析 os.Args 变得简单，例如对 Flag 和 PropertyMode 等内容的支持，下面将介这些内容。

## <a name='-1'></a>示例

首先我们要创建一个名为 greet 的目录，并在其中创建一个名为 main.go 的文件

```shell
mkdir greet
touch greet/main.go
```

main.go 中将包含以下代码

```go
package main

import (
	"fmt"

	"github.com/o8x/jk/args"
)

func main() {
	a := args.Args{
		App: &args.App{
			Name:  "Greet",
			Usage: "Greet -h",
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}

	fmt.Println("Hello Friend.")
}
```

将我们的命令安装到 `$GOPATH/bin` 中

```shell 
go install
```

运行我们的新命令

```shell
> greet
Hello Friend.
```

### <a name='Flags'></a>Flags

通过使用 Args 的实例来设置和查询 Flag 都很容易，使用内置方法取值时可以省略 Flag 前导的 -。

```go
package main

import (
	"fmt"

	"github.com/o8x/jk/args"
)

func main() {
	a := args.Args{
		App: &args.App{
			Name:  "Greet",
			Usage: "Greet -h",
		},
		Flags: []*args.Flag{
			{
				Name:        []string{"-language", "-lang", "-l"},
				Description: "Greet 的语言参数",
				Default:     []string{"zh-CN"},
			},
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}

	name := "o8x"

	switch a.Get("lang") {
	case "zh-CN":
		fmt.Printf("你好，")
	case "english":
		fmt.Printf("Hello,")
	case "spanish":
		fmt.Printf("Hola,")
	}

	fmt.Printf("%s\n", name)
}
```

运行它将会得到如下结果

```shell
> greet
你好，o8x
> greet -l spanish
Hola,o8x
> greet -lang english
Hello,o8x
```

### <a name='-1'></a>未定义

如果输入一个未定义的参数，将会的到一个错误

```shell
> greet -s
error: flag -s is not defined
```

## <a name='Flag'></a>Flag

定义

```go
type Properties map[string]any

type Flag struct {
	Name             []string `json:"name"`
	Description      string   `json:"description"`
	PropertyMode     bool     `json:"property_mode"`
	Default          []string `json:"default"`
	Required         bool     `json:"required"`
	Env              []string `json:"env"`
	Error            error    `json:"error_message"`
	NoValue          bool     `json:"no_value"`
	ValuesOnlyInEnum []string `json:"value_only_in_enum"`
	SingleValue      bool     `json:"single_value"`
	values           []string
	properties       Properties
	exist            bool
}
```

### <a name='Name'></a>Name

	Flag.Name []string

Flag 的名称和别名，必须以 - 或 -- 开头，否则将不会被作为 Flag 进行解析。

### <a name='Description'></a>Description

	Flag.Description string

Flag 的描述信息或 Usage，被用于输出 Help

### <a name='PropertyMode'></a>PropertyMode

	Flag.PropertyMode bool

为该 Flag 启用 [ Property 模式](#-property-模式)

### <a name='SingleValue'></a>SingleValue

	Flag.SingleValue bool

默认允许重复使用 Flag，可以通过设置 SingleValue 为 true 禁止这一行为

### <a name='Default'></a>Default

	Flag.Default []string

为 Flag 设置默认值，允许设置多个值，但当 SingleValue 为 true 时只会使用第一个值。

### <a name='Required'></a>Required

	Flag.Required bool

声明 Flag 是必填的，在没有设置 Flag、Env、Default 时将会得到错误。

### <a name='Env'></a>Env

	Flag.Env []string

该 Flag 可以从环境变量中取值，优先级最低。如果设置了多个变量名，将会依次取值，结果集形式上类似 Default。

### <a name='Error'></a>Error

	Flag.Error error

当 Flag 不符合属性规定时，使用自定义错误进行报错，暂时仅支持 Required。

### <a name='NoValue'></a>NoValue

	Flag.NoValue bool

该 Flag 仅需要存在，无需设置值

### <a name='ValuesOnlyInEnum'></a>ValuesOnlyInEnum

	Flag.ValuesOnlyInEnum []string

Flag 的实际参数必须在这个数组中取值，无论 SingleValue 是否为 true

### <a name='-1'></a>优先级

cmdline Value -> Default -> Env

## <a name='-1'></a>取值

Args 提供了一系列取值 API

### <a name='Flag-1'></a>判断 Flag 是否存在

    func (a *Args) IsSet(name string) bool

示例

```go
if a.IsSet("flag") {
	fmt.Println("已设置 flag 属性")
}
```

### <a name='Int64Flag'></a>获取 Int64 Flag

在未设置 SingleValue 时，会取出第一个 Flag 设置的实际参数

    func (a *Args) GetInt64(name string) (int64, bool)

示例

```go
package main

import (
	"fmt"
	"github.com/o8x/jk/args"
)

func main() {
	a := args.Args{
		Flags: []*args.Flag{
			{
				Name: []string{"-int", "-i"},
			},
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}

	fmt.Println(a.GetInt("int"))
}
```

运行它

```shell
> go run . -int 3
3 true
```

```shell
> go run . -int 3 -i 6
3 true
```

### <a name='Int64Flag-1'></a>获取 Int64 数组 Flag

    func (a *Args) GetInt64s(name string) ([]int64, error)

示例

```go
fmt.Println(a.GetInt64s("int"))
```

运行它

```shell
> go run . -int 3 -i 6
[3 6] <nil>
```

### <a name='IntFlag'></a>获取 Int Flag

    func (a *Args) GetInt(name string) (int, bool)

等同 GetInt64

### <a name='IntFlag-1'></a>获取 Int 数组 Flag

    func (a *Args) GetInts(name string) ([]int)

等同 GetInt64s

### <a name='Flag-1'></a>获取 字符串 Flag

	func (a *Args) Get(name string) (string, bool) 

示例

```go
package main

import (
	"fmt"
	"github.com/o8x/jk/args"
)

func main() {
	a := args.Args{
		Flags: []*args.Flag{
			{
				Name: []string{"-str", "-s"},
			},
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}

	fmt.Println(a.Get("str"))
}
```

运行它

```shell
> go run . -s jk
jk true
```

```shell
> go run . -s jk -s ef -s abc
jk true
```

### <a name='Flagpanic'></a>获取不到有效 字符串 Flag 时 panic

    func (a *Args) GetX(name string) string

示例

```go
fmt.Println(a.GetX("undefined"))
```

运行它

```
go run .
panic: flag -undefined is not defined
```

### <a name='Flag-1'></a>获取 字符串 数组 Flag

    func (a *Args) Gets(name string) []string

示例

```go
fmt.Println(a.Gets("str"))
```

运行它

```shell
> go run . -s jk -s ef -s abc
[jk ef abc]
```

### <a name='BoolFlag'></a>获取 Bool Flag

第一个返回值为参数的实际值，如果参数存在且实际值为 true 或 false，则第二个返回值为 true，否则为 false

    func (a *Args) GetBool(name string) (bool, bool)

示例

```go
fmt.Println(a.Gets("bool"))
```

运行它

```shell
> go run . -b true
true true
```

```shell
> go run . -b false
false true
```

```shell
> go run . -b 1
false false
```

## <a name='Property'></a> Property 模式

即将形式是 K=V 的 Flag 值自动格式化为 map 类型，在设置了 Flag.PropertyMode 为 true 时将会自动为该 Flag 开启 Property 模式。

示例

```go
package main

import (
	"fmt"
	"github.com/o8x/jk/args"
)

func main() {
	a := args.Args{
		Flags: []*args.Flag{
			{
				Name:         []string{"-property", "-p"},
				PropertyMode: true,
			},
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}

	fmt.Println(a.GetProperties("property"))
}
```

运行它

```
> go run . -property app.name=app -property app.port=8080 -property app.workdir=/sbin
map[app.name:app app.port:8080 app.workdir:/sbin]
```

### <a name='PropertyFlagProperties'></a>获取 Property Flag 的 Properties

返回值为 Properties 类型，内部类型为 map[string]string

    func (a *Args) GetProperties(name string) Properties

### <a name='FlagProperty'></a>获取 Flag 的 Property 的值

    func (a *Args) GetProperty(name string, property string) (string, bool)

示例

```shell
fmt.Println(a.GetProperty("property", "app.workdir"))
```

运行它

```
> go run . -property app.workdir=/sbin
/sbin true
```

### <a name='Propertypanic'></a>在获取不到有效 Property 时 panic

    func (a *Args) GetPropertyX(name string, property string) string

示例

```shell
fmt.Println(a.GetPropertyX("property", "app.name"))
```

运行它

```
> go run . -p app.port=:8080
panic: property property.app.name not found
```

### <a name='-1'></a>取值方法

用法参考 Args，不再赘述

	func (p Properties) GetInt(name string) (int, bool)
	func (p Properties) GetInt64(name string) (int64, bool)
	func (p Properties) IsSet(name string) bool
	func (p Properties) Get(name string) (string, bool)
	func (p Properties) GetBool(name string) (bool, bool)

## <a name='-1'></a>帮助

args 会生成整洁的帮助文本，未提供 Flag 时只有 help 和 version 两个 command

```shell
package main

import (
	"fmt"
	"github.com/o8x/jk/args"
)

func main() {
	a := args.Args{
		App: &args.App{
			Name:  "Greet",
			Usage: "Greet -h",
		},
		Flags: []*args.Flag{
			{
				Name:        []string{"-language", "-lang", "-l"},
				Description: "Greet 的语言参数",
				Default:     []string{"zh-CN"},
				Required:    true,
				Error:       fmt.Errorf("{{name}} is required"),
			},
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}
}
```

运行它

```shell
> greet -h
Usage of greet

Greet v0.0.1 (darwin/arm64) go1.19
usage: Greet -h

commands: 
    -language|-lang|-l: required
        Greet 的语言参数 (default: zh-CN)
    -help|-h
        print this help and exit
    -version|-v
        print version info and exit
```

## <a name='-1'></a>版本

args 会生成整洁的版本信息

```go
package main

import (
	"fmt"
	"time"

	"github.com/o8x/jk/args"
)

var Version = "0.0.1"
var Changelog = "版本更新日志"
var Banner = ` ________                      __   
 /  _____/______   ____   _____/  |_ 
/   \  __\_  __ \_/ __ \_/ __ \   __\
\    \_\  \  | \/\  ___/\  ___/|  |  
 \______  /__|    \___  >\___  >__|  
        \/            \/     \/`
var CommitHash = "138152a4247d2cedd930e8c437917d8ad4c2c674"
var Date = time.Now().Local().String()

func main() {
	a := args.Args{
		App: &args.App{
			Name:       "Greet",
			Usage:      "Greet -h",
			Copyright:  "Copyright © 2023 Alex.",
			Version:    &Version,
			Changelog:  &Changelog,
			Banner:     &Banner,
			CommitHash: &CommitHash,
			Date:       &Date,
		}
	}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}
}
```

运行它

```shell
> greet -v
 ________                      __   
 /  _____/______   ____   _____/  |_ 
/   \  __\_  __ \_/ __ \_/ __ \   __\
\    \_\  \  | \/\  ___/\  ___/|  |  
 \______  /__|    \___  >\___  >__|  
        \/            \/     \/

Greet v0.0.1 (darwin/arm64) go1.19

Release-Date: 2023-01-16 11:40:46.2418 +0800 CST
Commit Hash: 138152a4247d2cedd930e8c437917d8ad4c2c674
Changelog: 版本更新日志
Copyright © 2023 Alex.
```

### <a name='-1'></a>注入版本信息

可能你已经注意到 Version、Changelog、Banner、CommitHash、Date 等属性是指针类型。这么做的原因是，通常我们并不会直接将版本信息作为常量写入代码，而是从外部获取或在编译期使用
ldflags -X 注入。

因为 ldflags 不接受带有空格的数据，所以 Args 对于上述的指针属性提供了 base64 支持，带有 base64:// 前缀的数据会被自动进行
decode，同时也会删除数据中前导和末尾的 \n。

```go
package main

import (
	"github.com/o8x/jk/args"
)

var (
	Version    = ""
	Changelog  = ""
	Banner     = ""
	CommitHash = ""
	Date       = ""
)

func main() {
	a := args.Args{
		App: &args.App{
			Name:       "Greet",
			Usage:      "Greet -h",
			Copyright:  "Copyright © 2023 Alex.",
			Version:    &Version,
			Changelog:  &Changelog,
			Banner:     &Banner,
			CommitHash: &CommitHash,
			Date:       &Date,
		},
	}

	if err := a.Parse(); err != nil {
		a.PrintErrorExit(err)
	}
}
```

重新安装

```shell
set -v

Banner=$(echo '  ________                      __   
 /  _____/______   ____   _____/  |_ 
/   \  __\_  __ \_/ __ \_/ __ \   __\
\    \_\  \  | \/\  ___/\  ___/|  |  
 \______  /__|    \___  >\___  >__|  
        \/            \/     \/' | gbase64 -w 0)
Version=0.0.1 
CommitHash=7fee7dfb2d17d9ac2bb8a30cc7b7605d3f94f543
Changelog=$(echo '版本更新日志' | gbase64 -w 0)
Date=$(date "+%Y-%m-%d %H:%M:%S" | gbase64 -w 0)

go install -ldflags "
    -X main.Version=$Version
    -X main.Banner=base64://$Banner
    -X main.Changelog=base64://$Changelog
    -X main.CommitHash=$CommitHash
    -X main.Date=base64://$Date" .
```

运行它

```shell
> greet -v             
  ________                      __   
 /  _____/______   ____   _____/  |_ 
/   \  __\_  __ \_/ __ \_/ __ \   __\
\    \_\  \  | \/\  ___/\  ___/|  |  
 \______  /__|    \___  >\___  >__|  
        \/            \/     \/

Greet v0.0.1 (darwin/arm64) go1.19

Release-Date: 2023-01-16 14:19:05
Commit Hash: 7fee7dfb2d17d9ac2bb8a30cc7b7605d3f94f543
Changelog: 版本更新日志
Copyright © 2023 Alex.
```

## <a name='API'></a>其他 API

### <a name='goroutine'></a> 阻塞主 goroutine

该方法会阻塞主 goroutine 直到 USR1、USR2、INT、TERM、QUIT 其中之一的信号被触发

    func (a *Args) WaitSignal() os.Signal

示例

```go
a.WaitSignal()
fmt.Println("app shutdown")
```

### <a name='-1'></a> 退出

该方法会退出应用程序，用法与 os.Exit() 一致

	func (a *Args) Exit(code int)

示例

```go
a.Exit(0)
```

### <a name='-1'></a>打印调试日志并退出

	func (a *Args) DumpExit()

示例

```go
a.DumpExit()
```

## <a name='-1'></a>脚手架

如果你只想从命令行获取一些简单的参数，不想大动干戈构造 Args 实例，Args 也提供了 Option API

### <a name='ArgFunc'></a>ArgFunc

注入方法类型，可以按此类型限定自由扩展 Option Function

    type ArgFunc func(a *Args)

### <a name='Args'></a>生成基本的 Args 实例

    func New(fns ...ArgFunc) *Args

### <a name='-1'></a>注入应用程序名称和使用示例

    func WithApp(name, usage string) ArgFunc

### <a name='Flag-1'></a>增加基本的 Flag

    func AddFlag(name string) ArgFunc

### <a name='Flag-1'></a>增加带有默认值的 Flag

    func AddDefaultFlag(name, def string) ArgFunc

### <a name='EnumFlag'></a>增加 Enum Flag

    func AddEnumFlag(name, def string, enum ...string) ArgFunc

### <a name='Flag-1'></a>增加没有值的 Flag

    func AddNoValueFlag(name string) ArgFu

### <a name='Option'></a>手动扩展 Option

任何符合 args.ArgFunc 类型的方法都可以作为 Args 的 Option

```go
func WithAppCopyright(c string) args.ArgFunc {
    return func(a *args.Args) {
        if a.App != nil {
            a.App.Copyright = c
            return
        }

        a.App = &args.App{
            Copyright: c,
        }
    }
}
```

### <a name='-1'></a>调用示例

```go
package main

import (
	"github.com/o8x/jk/args"
)

func main() {
	a := args.New(
		args.WithApp("user", "./app -name Alex -dryrun -time 2"),
		args.AddFlag("-name"),
		args.AddDefaultFlag("-def", "Def"),
		args.AddNoValueFlag("-dryrun"),
		args.AddEnumFlag("-ts", "now", "nextDay", "now"),
		WithAppCopyright("Alex"),
	)

	if err := a.Parse(); err != nil {
		a.PrintHelpExit(err)
	}

	a.DumpExit()
}
```

Manual
===

## 初始化 logrus

初始化 logrus 并设置输出目标

    func New(level logrus.Level, out io.Writer) *logrus.Logger

初始化 logrus 并将输出目标设置为使用 lumberjack 实现的自动切割日志写入器

	func NewFile(level string, out string) *logrus.Logger

## 使用默认日志写入器

无需调用任何方法，日志将会自动被输出到 stdout，实际行为类似 `New(logrus.DebugLevel, os.Stdout)`

## 使用默认文件日志写入器

相当于 NewFile("info", "log/jk.log")

	func UseDefault()

## lumberjack

生成一个支持自动切割日志的写入器

	func NewRotater(file string) io.Writer

## 获取 logrus 实例

在没有进行 Init 时调用该方法，将会自动进行对 `/dev/stdout` 进行 info 级别的初始化

    func Get() *logrus.Logger

## 注入 logrus 实例

使用已有的 logrus 实例替换 std logger

    func UseLogger(l *logrus.Logger)

## API

提供了以下的API，一般情况下不需要使用 Get() 获取 logger 实例

### 写入 Info 级别的日志

    func Info(format string, args ...any)

### 写入 Fatal 级别的日志

将会退出程序

    func Fatal(format string, args ...any)

### 写入 Warn 级别的日志

    func Warn(format string, args ...any)

### 写入 Error 级别的日志

    func Error(format string, args ...any)

### 写入 Debug 级别的日志

    func Debug(format string, args ...any)

### Group

补充功能，本质上是 WithField("group", strings.ToUpper(group))

    func Group(group string) *logrus.Entry

### WithField

等同 logrus.WithField

    func WithField(key string, value any) *logrus.Entry

### WithError

等同 logrus.WithError

    func WithError(err error) *logrus.Entry

## 示例

```go
package main

import (
	"fmt"

	"github.com/o8x/jk/logger"
)

func main() {
	logger.UseLogger(logger.NewFile("info", "/dev/stdout"))

	logger.Info("Info log")
	logger.Fatal("Fatal log")
	logger.Warn("Warn log")
	logger.Error("Error log")
	logger.Group("group").Info("log with group")
	logger.WithError(fmt.Errorf("with error")).Info("log with error")
	logger.WithField("n", 1).Info("log with field")

	logger.Get().Error("use Get() print log")
}
```

运行它

```shell
> go run .
time="2023-01-17T11:25:56+08:00" level=info msg="Info log"
time="2023-01-17T11:25:56+08:00" level=fatal msg="Fatal log"
time="2023-01-17T11:25:56+08:00" level=warning msg="Warn log"
time="2023-01-17T11:25:56+08:00" level=error msg="Error log"
time="2023-01-17T11:25:56+08:00" level=info msg="log with group" group=GROUP
time="2023-01-17T11:25:56+08:00" level=info msg="log with error" error="with error"
time="2023-01-17T11:25:56+08:00" level=info msg="log with field" n=1
time="2023-01-17T11:25:56+08:00" level=error msg="use Get() print log"
```

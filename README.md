JK
=======

[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)
[![LICENSE](https://img.shields.io/badge/license-Anti%20996-blue.svg)](https://github.com/996icu/996.ICU/blob/master/LICENSE)

Go Development Kit

## Kits

* [args](./args) : 无外部依赖的命令行参数解析器
* [signal](signal) : 应用程序信号处理
* [fs](fs) : 文件系统辅助函数
* [djb2](djb2): djb2 哈希算法
* [response](response): JSON 响应生成器
* [logger](logger): 带有日志切割的 Logrus 工具
* [sqlite](sqlite): sqlite 工具
* [syncmap](syncmap): 支持泛型的 sync.Map
* [rand](rand): 随机数 lib
* [http2](http2): http2 lib
* [xor](xor): 简单的 XOR 加密和计算 lib
* [gzip](gzip): gzip 编解码器
* [base58](base58): base58 编解码器
* [hash](hash): sha1、sha2、md5 等 hash lib
* [size](size): 格式化 int 类型的 size
* [uniqid](uniqid): 基于纳秒的唯一ID生成器
* [crash](crash): crash recovery function
* [tcp](tcp): tcp 服务器
* [udp](udp): udp 服务器
* [http](http): 一次性 http Get/Post 请求构造器
* [json](json): 带有美化功能简单 json 编解码器
* [puresqlite](puresqlite): 纯 go sqlite lib
* [context](context): context 的封装，实现类似面向对象的用法
* [utils](x): 一些实用方法

## Tools

* [cmd/cert](cmd/cert): tls 自签证书生成工具

## Example

```go
package main

import (
	"github.com/o8x/jk"
)

func main() {
	jk.Hello()
}
```

```shell
> go run .
github.com/o8x/jk say hello to you
```

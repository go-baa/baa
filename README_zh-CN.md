# [Baa](http://go-baa.github.io/baa) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-baa/baa) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) [![Build Status](http://img.shields.io/travis/go-baa/baa.svg?style=flat-square)](https://travis-ci.org/go-baa/baa) [![Coverage Status](http://img.shields.io/coveralls/go-baa/baa.svg?style=flat-square)](https://coveralls.io/r/go-baa/baa)

一个简单高效的Go web开发框架。主要有路由、中间件，依赖注入和HTTP上下文构成。

Baa 不使用 ``反射``和``正则``，没有魔法的实现。

## 文档

* [简体中文](https://github.com/go-baa/doc/tree/master/zh-CN)
* [English](https://github.com/go-baa/doc/tree/master/en-US)
* [godoc](https://godoc.org/github.com/go-baa/baa)

## 快速上手

安装：

```
go get -u gopkg.in/baa.v1
```

示例：

```
// baa.go
package main

import (
    "gopkg.in/baa.v1"
)

func main() {
    app := baa.New()
    app.Get("/", func(c *baa.Context) {
        c.String(200, "Hello, 世界")
    })
    app.Run(":1323")
}
```

编译:

Baa use encoding/json as default json package but you can change to [jsoniter](https://github.com/json-iterator/go) by build from other tags

```
go build -tags=jsoniter .
```

运行:

```
go run baa.go
```

浏览:

```
http://127.0.0.1:1323/
```

## 特性

* 支持静态路由、参数路由、组路由（前缀路由/命名空间）和路由命名
* 路由支持链式操作
* 路由支持文件/目录服务
* 中间件支持链式操作
* 支持依赖注入*
* 支持JSON/JSONP/XML/HTML格式输出
* 统一的HTTP错误处理
* 统一的日志处理
* 支持任意更换模板引擎（实现baa.Renderer接口即可）

## 示例

https://github.com/go-baa/example

* [blog](https://github.com/go-baa/example/tree/master/blog)
* [api](https://github.com/go-baa/example/tree/master/api)
* [websocket](https://github.com/go-baa/example/tree/master/websocket)

## 中间件

* [gzip](https://github.com/baa-middleware/gzip)
* [accesslog](https://github.com/baa-middleware/accesslog)
* [recovery](https://github.com/baa-middleware/recovery)
* [session](https://github.com/baa-middleware/session)
* [static](https://github.com/baa-middleware/static)
* [requestcache](https://github.com/baa-middleware/requestcache)
* [nocache](https://github.com/baa-middleware/nocache)
* [jwt](https://github.com/baa-middleware/jwt)
* [cors](https://github.com/baa-middleware/cors)
* [authz](https://github.com/baa-middleware/authz)

## 组件

* [cache](https://github.com/go-baa/cache)
* [render](https://github.com/go-baa/render)
* [pongo2](https://github.com/go-baa/pongo2)
* [router](https://github.com/go-baa/router)
* [pool](https://github.com/go-baa/pool)
* [bat](https://github.com/go-baa/bat)
* [log](https://github.com/go-baa/log)
* [setting](https://github.com/go-baa/setting)

## 性能测试

### 路由测试

使用 [go-http-routing-benchmark] (https://github.com/safeie/go-http-routing-benchmark) 测试, 2016-02-27 更新.

##### [GitHub API](http://developer.github.com/v3)

> Baa的路由性能非常接近 Echo.

```
BenchmarkBaa_GithubAll          	   30000	     50984 ns/op	       0 B/op	       0 allocs/op
BenchmarkBeego_GithubAll        	    3000	    478556 ns/op	    6496 B/op	     203 allocs/op
BenchmarkEcho_GithubAll         	   30000	     47121 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubAll          	   30000	     41004 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubAll   	    3000	    450709 ns/op	  131656 B/op	    1686 allocs/op
BenchmarkGorillaMux_GithubAll   	     200	   6591485 ns/op	  154880 B/op	    2469 allocs/op
BenchmarkMacaron_GithubAll      	    2000	    679559 ns/op	  201140 B/op	    1803 allocs/op
BenchmarkMartini_GithubAll      	     300	   5680389 ns/op	  228216 B/op	    2483 allocs/op
BenchmarkRevel_GithubAll        	    1000	   1413894 ns/op	  337424 B/op	    5512 allocs/op
```

### HTTP测试

#### 代码

Baa:

```
package main

import (
	"gopkg.in/baa.v1"
)

func main() {
	app := baa.New()
    app.Get("/", func(c *baa.Context) {
        c.String(200, "Hello, 世界")
    })
    app.Run(":1323")
}
```

#### 测试结果:

```
$ wrk -t 10 -c 100 -d 30 http://127.0.0.1:1323/
Running 30s test @ http://127.0.0.1:1323/
  10 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.64ms  299.23us   8.25ms   66.84%
    Req/Sec     6.11k   579.08     8.72k    68.74%
  1827365 requests in 30.10s, 228.30MB read
Requests/sec:  60704.90
Transfer/sec:      7.58MB
```

## 案例

目前使用在 健康一线 的私有项目中。

## 贡献

Baa的灵感来自 [beego](https://github.com/astaxie/beego) [echo](https://github.com/labstack/echo) [macaron](https://github.com/go-macaron/macaron)

- [safeie](https://github.com/safeie)、[micate](https://github.com/micate) - Author
- [betty](https://github.com/betty3039) - Language Consultant
- [Contributors](https://github.com/go-baa/baa/graphs/contributors)

## License

This project is under the MIT License (MIT) See the [LICENSE](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) file for the full license text.

# [Baa](http://go-baa.github.io/baa) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-baa/baa) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) [![Build Status](http://img.shields.io/travis/go-baa/baa.svg?style=flat-square)](https://travis-ci.org/go-baa/baa) [![Coverage Status](http://img.shields.io/coveralls/go-baa/baa.svg?style=flat-square)](https://coveralls.io/r/go-baa/baa)

一个简单高效的Go web开发框架。主要有路由、中间件，依赖注入和HTTP上下文构成。

Baa 不使用 ``反射``和``正则``，没有魔法的实现。

## 快速上手

安装：

```
go get -u gopkg.in/baa.v1
```

示例：

```
package main

import (
    "gopkg.in/baa.v1"
)

func main() {
    app := baa.New()
    app.Get("/", func(c *baa.Context) {
        c.String(200, "Hello World!")
    })
    app.Run(":1323")
}
```

## 特性

* 支持静态路由、参数路由、组路由（前缀路由/命名空间）和路由命名
* 路由支持链式操作
* 路由支持文件/目录服务
* 支持中间件和链式操作
* 支持依赖注入*
* 支持JSON/JSONP/XML/HTML格式输出
* 统一的HTTP错误处理
* 统一的日志处理
* 支持任意更换模板引擎（实现baa.Renderer接口即可）

## 中间件

* [gzip](https://github.com/baa-middleware/gzip)
* [logger](https://github.com/baa-middleware/logger)
* [recovery](https://github.com/baa-middleware/recovery)
* [session](https://github.com/baa-middleware/session)

## 组件(DI)

* [cache](https://github.com/go-baa/cache)
* [render](https://github.com/go-baa/render)

## 性能测试

和快速的Echo框架对比 [Echo](https://github.com/labstack/echo)

> 注意：

[Echo](https://github.com/labstack/echo) 在V2版本中使用了fasthttp，我们这里使用 [Echo V1](https://github.com/labstack/echo/releases/tag/v1.4) 测试。

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
	"github.com/baa-middleware/logger"
	"github.com/baa-middleware/recovery"
	"gopkg.in/baa.v1"
)

func hello(c *baa.Context) {
	c.String(200, "Hello, World!\n")
}

func main() {
	b := baa.New()
	b.Use(logger.Logger())
	b.Use(recovery.Recovery())

	b.Get("/", hello)

	b.Run(":8001")
}
```

Echo:

```
package main

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

// Handler
func hello(c *echo.Context) error {
	return c.String(200, "Hello, World!\n")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())

	// Routes
	e.Get("/", hello)

	// Start server
	e.Run(":8001")
}
```

#### 测试结果:

> Baa 在http中的表现还稍稍比 Echo 好一些。

Baa:

```
$ wrk -t 10 -c 100 -d 30 http://127.0.0.1:8001/
Running 30s test @ http://127.0.0.1:8001/
  10 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.92ms    1.43ms  55.26ms   90.86%
    Req/Sec     5.46k   257.26     6.08k    88.30%
  1629324 requests in 30.00s, 203.55MB read
Requests/sec:  54304.14
Transfer/sec:      6.78MB
```

Echo:

```
$ wrk -t 10 -c 100 -d 30 http://127.0.0.1:8001/
Running 30s test @ http://127.0.0.1:8001/
  10 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.83ms    3.76ms  98.38ms   90.20%
    Req/Sec     4.79k     0.88k   45.22k    96.27%
  1431144 requests in 30.10s, 178.79MB read
Requests/sec:  47548.11
Transfer/sec:      5.94MB
```


## 案例

目前使用在 健康一线 的私有项目中。

## 手册

[godoc](http://godoc.org/github.com/go-baa/baa)

[document](#)


## 贡献

Baa的灵感来自 [beego](https://github.com/astaxie/beego) [echo](https://github.com/labstack/echo) [macaron](https://github.com/go-macaron/macaron)

- [safeie](https://github.com/safeie)、[micate](https://github.com/micate) - Author
- [betty](https://github.com/betty3039) - Language Consultant
- [Contributors](https://github.com/go-baa/baa/graphs/contributors)

## License

This project is under the MIT License (MIT) See the [LICENSE](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) file for the full license text.

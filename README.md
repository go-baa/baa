# [Baa](http://go-baa.github.io/baa) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-baa/baa) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) [![Build Status](http://img.shields.io/travis/go-baa/baa.svg?style=flat-square)](https://travis-ci.org/go-baa/baa) [![Coverage Status](http://img.shields.io/coveralls/go-baa/baa.svg?style=flat-square)](https://coveralls.io/r/go-baa/baa)

an express Go web framework with routing, middleware, dependency injection, http context. 

Baa is ``no reflect``, ``no regexp``.

## document

* [简体中文](https://github.com/go-baa/doc/tree/master/zh-CN)
* [English](https://github.com/go-baa/doc/tree/master/en-US)
* [godoc](https://godoc.org/github.com/go-baa/baa)

## Getting Started

Install:

```
go get -u gopkg.in/baa.v1
```

Example:

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

Build:

Baa use encoding/json as default json package but you can change to [jsoniter](https://github.com/json-iterator/go) by build from other tags

```
go build -tags=jsoniter .
```

Run:

```
go run baa.go
```

Explore:

```
http://127.0.0.1:1323/
```

## Features

* route support static, param, group
* route support handler chain
* route support static file serve
* middleware supoort handle chain
* dependency injection support*
* context support JSON/JSONP/XML/HTML response
* centralized HTTP error handling
* centralized log handling
* whichever template engine support(emplement baa.Renderer)

## Examples

https://github.com/go-baa/example

* [blog](https://github.com/go-baa/example/tree/master/blog)
* [api](https://github.com/go-baa/example/tree/master/api)
* [websocket](https://github.com/go-baa/example/tree/master/websocket)

## Middlewares

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

## Components

* [cache](https://github.com/go-baa/cache)
* [render](https://github.com/go-baa/render)
* [pongo2](https://github.com/go-baa/pongo2)
* [router](https://github.com/go-baa/router)
* [pool](https://github.com/go-baa/pool)
* [bat](https://github.com/go-baa/bat)
* [log](https://github.com/go-baa/log)
* [setting](https://github.com/go-baa/setting)

## Performance

### Route Test

Based on [go-http-routing-benchmark] (https://github.com/safeie/go-http-routing-benchmark), Feb 27, 2016.

##### [GitHub API](http://developer.github.com/v3)

> Baa route test is very close to Echo.

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

### HTTP Test

#### Code

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

#### Result:

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

## Use Cases

vodjk private projects.

## Credits

Get inspirations from [beego](https://github.com/astaxie/beego) [echo](https://github.com/labstack/echo) [macaron](https://github.com/go-macaron/macaron)

- [safeie](https://github.com/safeie)、[micate](https://github.com/micate) - Author
- [betty](https://github.com/betty3039) - Language Consultant
- [Contributors](https://github.com/go-baa/baa/graphs/contributors)

## License

This project is under the MIT License (MIT) See the [LICENSE](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) file for the full license text.

# [Baa](http://go-baa.github.io/baa) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-baa/baa) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) [![Build Status](http://img.shields.io/travis/go-baa/baa.svg?style=flat-square)](https://travis-ci.org/go-baa/baa) [![Coverage Status](http://img.shields.io/coveralls/go-baa/baa.svg?style=flat-square)](https://coveralls.io/r/go-baa/baa)

an express Go web framework with routing, middleware, dependency injection, http context. 

baa is ``no reflect``, ``no regexp``.

## Getting Started

```
package main

import (
    "github.com/go-baa/baa"
)

func main() {
    app := baa.New()
    app.Get("/", func(c *baa.Context) {
        c.String(200, "Hello World!")
    })
    app.Run(":1323")
}
```

## Features

* route support static, param, group
* route support handler chain
* route support static file serve
* middleware supoort handle chain
* dependency injection support
* context support JSON/JSONP/XML/HTML response
* centralized HTTP error handling
* centralized log handling
* whichever template engine support(emplement baa.Renderer)


## Middleware

* [gzip](https://github.com/baa-middleware/gzip)
* [logger](https://github.com/baa-middleware/logger)
* [recovery](https://github.com/baa-middleware/recovery)
* [render](https://github.com/baa-middleware/render)
* [session](https://github.com/baa-middleware/session)


## Performance

Based on [go-http-routing-benchmark] (https://github.com/julienschmidt/go-http-routing-benchmark), Feb 27, 2016.

##### [GitHub API](http://developer.github.com/v3)

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


## Use Cases

NONE

## Guide

[godoc](http://godoc.org/github.com/go-baa/baa)

[document](#)


## Credits

Get inspirations from [beego](https://github.com/astaxie/beego) [echo](https://github.com/labstack/echo) [macaron](https://github.com/go-macaron/macaron)

- [safeie](https://github.com/safeie)、[micate](https://github.com/micate) - Author
- [betty](https://github.com/betty3039) - Language Consultant
- [Contributors](https://github.com/go-baa/baa/graphs/contributors)

## License

This project is under the MIT License (MIT) See the [LICENSE](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) file for the full license text.
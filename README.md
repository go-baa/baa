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

Based on [go-http-routing-benchmark] (https://github.com/julienschmidt/go-http-routing-benchmark), Feb 24, 2016.

##### [GitHub API](http://developer.github.com/v3)

```
> Baa:         68317 ns/op,       0 B/op,     0 allocs/op
> Beego:     1109160 ns/op,  146272 B/op,  2092 allocs/op
> Echo:        38662 ns/op,       0 B/op,     0 allocs/op
> Macaron:    709932 ns/op,  224960 B/op,  2315 allocs/op
> Martini:  10261331 ns/op,  237953 B/op,  2686 allocs/op
```

![Performance](http://i.imgur.com/hB2qdRS.png)

```
BenchmarkAce_GithubAll              20000             93675 ns/op           13792 B/op       167 allocs/op
BenchmarkBaa_GithubAll              20000             68317 ns/op               0 B/op         0 allocs/op
BenchmarkBear_GithubAll             10000            264194 ns/op           79952 B/op       943 allocs/op
BenchmarkBeego_GithubAll             2000           1109160 ns/op          146272 B/op      2092 allocs/op
BenchmarkBone_GithubAll              1000           2063973 ns/op          648016 B/op      8119 allocs/op
BenchmarkDenco_GithubAll            20000             83114 ns/op           20224 B/op       167 allocs/op
BenchmarkEcho_GithubAll             30000             38662 ns/op               0 B/op         0 allocs/op
BenchmarkGin_GithubAll              30000             43467 ns/op               0 B/op         0 allocs/op
BenchmarkGocraftWeb_GithubAll        5000            386829 ns/op          133280 B/op      1889 allocs/op
BenchmarkGoji_GithubAll              3000            561131 ns/op           56113 B/op       334 allocs/op
BenchmarkGoJsonRest_GithubAll        3000            490789 ns/op          135995 B/op      2940 allocs/op
BenchmarkGoRestful_GithubAll          100          15569513 ns/op          797239 B/op      7725 allocs/op
BenchmarkGorillaMux_GithubAll         200           7431130 ns/op          153137 B/op      1791 allocs/op
BenchmarkHttpRouter_GithubAll       30000             51192 ns/op           13792 B/op       167 allocs/op
BenchmarkHttpTreeMux_GithubAll      10000            138164 ns/op           56112 B/op       334 allocs/op
BenchmarkKocha_GithubAll            10000            139625 ns/op           23304 B/op       843 allocs/op
BenchmarkMacaron_GithubAll           2000            709932 ns/op          224960 B/op      2315 allocs/op
BenchmarkMartini_GithubAll            100          10261331 ns/op          237953 B/op      2686 allocs/op
BenchmarkPat_GithubAll                500           3989686 ns/op         1504104 B/op     32222 allocs/op
BenchmarkPossum_GithubAll            5000            259165 ns/op           97441 B/op       812 allocs/op
BenchmarkR2router_GithubAll         10000            240345 ns/op           77328 B/op      1182 allocs/op
BenchmarkRevel_GithubAll             2000           1203336 ns/op          345554 B/op      5918 allocs/op
BenchmarkRivet_GithubAll            10000            247213 ns/op           84272 B/op      1079 allocs/op
BenchmarkTango_GithubAll             5000            379960 ns/op           87081 B/op      2470 allocs/op
BenchmarkTigerTonic_GithubAll        2000            931401 ns/op          241089 B/op      6052 allocs/op
BenchmarkTraffic_GithubAll            200           7292170 ns/op         2664770 B/op     22390 allocs/op
BenchmarkVulcan_GithubAll            5000            271682 ns/op           19894 B/op       609 allocs/op
BenchmarkZeus_GithubAll              2000            748827 ns/op          300688 B/op      2648 allocs/op
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
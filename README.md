# [Baa](http://go-baa.github.io/baa) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-baa/baa) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/go-baa/baa/master/LICENSE) [![Build Status](http://img.shields.io/travis/go-baa/baa.svg?style=flat-square)](https://travis-ci.org/go-baa/baa) [![Coverage Status](http://img.shields.io/coveralls/go-baa/baa.svg?style=flat-square)](https://coveralls.io/r/go-baa/baa)

a fast &amp; simple Go web framework, routing, middleware, dependency injection, http context.

## Features

* no reflect
* no regexp
* route support static, param, group
* route support handler chain
* route support static file serve
* middleware supoort handle chain
* dependency injection support
* context support JSON/JSONP/XML/HTML reponse
* Centralized HTTP error handling
* Centralized log handling (use baa.Logger interface)

## Performance

## Middleware

* [logger](https://github.com/baa-middleware/logger)
* [recovery](https://github.com/baa-middleware/recovery)
* [session](https://github.com/baa-middleware/session)
* [render](https://github.com/baa-middleware/render)
* [gzip](https://github.com/baa-middleware/gzip)


## TODO

* Middleware: session, gzip, csrf, render

# Install

```
$ go get github.com/go-baa/baa
```

# Quick Start

classic

```
package main

import (
    "github.com/go-baa/baa"
)

func main() {
    app := baa.Classic()
    app.Get("/", func(c *baa.Context) {
        c.String(200, "Hello World!")
    })
    app.Run(":1323")
}
```

use middleware

```
package main

import (
    "github.com/go-baa/baa"
    "github.com/baa-middleware/logger"
    "github.com/baa-middleware/recovery"
)

func main() {
    app := baa.New()
    app.Use(logger.Logger())
    app.Use(recovery.Recovery())
    app.Get("/", func(c *baa.Context) {
        c.String(200, "Hello World!")
    })
    app.Run(":1323")
}
```
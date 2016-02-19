# [Baa](http://go-baa.github.io/baa) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/go-baa/baa) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE) [![Build Status](http://img.shields.io/travis/go-baa/baa.svg?style=flat-square)](https://travis-ci.org/go-baa/baa) [![Coverage Status](http://img.shields.io/coveralls/go-baa/baa.svg?style=flat-square)]

a fast &amp; simple Go web framework, routing, middleware, dependency injection, http context.

## Features

* no reflect
* no regexp
* route support static and param
* route support handle chain
* route support static file serve
* middleware supoort handle chain
* dependency injection support
* context support JSON/JSONP/XML/HTML reponse
* Centralized HTTP error handling

## Performance

## Middleware

* logger
* recovery
* session
* render


## TODO

* group router
* static file serve
* context json/jsonp/xml support
* context cookie
* render and template
* Middleware: session, gzip, csrf

# Install

```
$ go get github.com/go-baa/baa
```

# Quick Start

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
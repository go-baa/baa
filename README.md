# baa
a fast &amp; simple Go web framework, routing, middleware, dependency injection, http context.

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
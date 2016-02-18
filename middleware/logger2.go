package middleware

import (
	"time"

	"github.com/go-baa/baa"
)

// Logger2 returns a Middleware that logs requests.
func Logger2() baa.MiddlewareFunc {
	return func(h baa.HandlerFunc) baa.HandlerFunc {
		return func(c *baa.Context) error {
			start := time.Now()

			l := c.Baa().GetLogger()
            l.Println("in log2")

			if err := h(c); err != nil {
				c.Error(err)
			}
            
			l.Printf("log2 %s %s %v", c.Req.Method, c.Req.RequestURI, time.Since(start))

            l.Println("out log2")
            
			return nil
		}
	}
}

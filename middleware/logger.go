package middleware

import (
	"time"

	"github.com/go-baa/baa"
)

// Logger returns a Middleware that logs requests.
func Logger() baa.MiddlewareFunc {
	return func(h baa.HandlerFunc) baa.HandlerFunc {
		return func(c *baa.Context) error {
			start := time.Now()

			l := c.Baa().GetLogger()
            //l.Println("in log1")

			if err := h(c); err != nil {
				c.Error(err)
			}
            
			l.Printf("log1 %s %s %v", c.Req.Method, c.Req.RequestURI, time.Since(start))

            //l.Println("out log1")
            
			return nil
		}
	}
}

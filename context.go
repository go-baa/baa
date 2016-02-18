package baa

import (
	"net/http"
)

// Context provlider a HTTP context for baa
// context contains reqest, response, header, cookie and some content type.
type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
	Data map[string]interface{}
	baa  *Baa
}

// NewContext create a http context
func NewContext(w http.ResponseWriter, r *http.Request, b *Baa) *Context {
	c := new(Context)
	c.reset(w, r, b)
	return c
}

// reset ...
func (c *Context) reset(w http.ResponseWriter, r *http.Request, b *Baa) {
	c.Req = r
	c.Resp = w
	c.baa = b
	c.Data = make(map[string]interface{})
}

// String write text by string
func (c *Context) String(code int, s string) error {
	c.Resp.Header().Set("Content-Type", "charset=utf-8")
	c.Resp.WriteHeader(code)
	c.Resp.Write([]byte(s))
	return nil
}

// Text write text by []byte
func (c *Context) Text(code int, s []byte) error {
	c.Resp.Header().Set("Content-Type", "charset=utf-8")
	c.Resp.WriteHeader(code)
	c.Resp.Write(s)
	return nil
}

// Redirect redirects the request using http.Redirect with status code.
func (c *Context) Redirect(code int, url string) error {
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return ErrInvalidRedirectCode
	}
	http.Redirect(c.Resp, c.Req, url, code)
	return nil
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *Context) Error(err error) {
	c.baa.httpErrorHandler(err, c)
}

// Baa ...
func (c *Context) Baa() *Baa {
	return c.baa
}

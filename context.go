package baa

import (
	"net/http"
	"strconv"
)

// Context provlider a HTTP context for baa
// context contains reqest, response, header, cookie and some content type.
type Context struct {
	Req          *http.Request
	Resp         *Response
	baa          *Baa
	Data         map[string]interface{}
	params       map[string]string
	routeHandler HandlerFunc // route match handler
	mi           int         // middleware order
}

// newContext create a http context
func newContext(w http.ResponseWriter, r *http.Request, b *Baa) *Context {
	c := new(Context)
	c.Resp = NewResponse(w, b)
	c.baa = b
	c.Data = make(map[string]interface{})
	c.params = make(map[string]string)
	c.reset(w, r)
	return c
}

// reset ...
func (c *Context) reset(w http.ResponseWriter, r *http.Request) {
	c.Resp.reset(w)
	c.Req = r
	c.mi = 0
	var k string
	for k = range c.Data {
		delete(c.Data, k)
	}
	for k = range c.params {
		delete(c.params, k)
	}
}

// SetParam read route param value from uri
func (c *Context) SetParam(key, value string) {
	c.params[key] = value
}

// Param get route param from context
func (c *Context) Param(name string) string {
	return c.params[name]
}

// ParamInt get route param from context and format to int
func (c *Context) ParamInt(name string) int {
	i, _ := strconv.Atoi(c.params[name])
	return i
}

// ParamInt64 get route param from context and format to int64
func (c *Context) ParamInt64(name string) int64 {
	i, _ := strconv.ParseInt(c.params[name], 10, 64)
	return i
}

// ParamFloat get route param from context and format to float64
func (c *Context) ParamFloat(name string) float64 {
	f, _ := strconv.ParseFloat(c.params[name], 64)
	return f
}

// ParamBool get route param from context and format to bool
func (c *Context) ParamBool(name string) bool {
	b, _ := strconv.ParseBool(c.params[name])
	return b
}

// Query get a param from http.Request.Form
func (c *Context) Query(name string) string {
	return ""
}

// QueryInt get a param from http.Request.Form and format to int
func (c *Context) QueryInt(name string) int {
	return 0
}

// QueryInt64 get a param from http.Request.Form and format to int64
func (c *Context) QueryInt64(name string) int64 {
	return 0
}

// QueryFloat get a param from http.Request.Form and format to float64
func (c *Context) QueryFloat(name string) float64 {
	return 0.0
}

// QueryBool get a param from http.Request.Form and format to bool
func (c *Context) QueryBool(name string) bool {
	return false
}

// QueryStrings get a group param from http.Request.Form and format to string slice
func (c *Context) QueryStrings(name string) []string {
	return nil
}

// Gets return http.Request.URL queryString params
func (c *Context) Gets() map[string]interface{} {
	return nil
}

// Posts return http.Request form data
func (c *Context) Posts() map[string]interface{} {
	return nil
}

// String write text by string
func (c *Context) String(code int, s string) {
	c.Resp.Header().Set("Content-Type", "charset=utf-8")
	c.Resp.WriteHeader(code)
	c.Resp.Write([]byte(s))
}

// Text write text by []byte
func (c *Context) Text(code int, s []byte) {
	c.Resp.Header().Set("Content-Type", "charset=utf-8")
	c.Resp.WriteHeader(code)
	c.Resp.Write(s)
}

// JSON write data by json format
func (c *Context) JSON(code int, d interface{}) {
}

// JSONP write data by jsonp format
func (c *Context) JSONP(code int, d interface{}) {
}

// XML write data by XML format
func (c *Context) XML(code int, d interface{}) {
}

// HTML write data by html template engine, use context.Data
func (c *Context) HTML(code int, tpl string) {
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
	c.baa.Error(err, c)
}

// Baa ...
func (c *Context) Baa() *Baa {
	return c.baa
}

// Next execute next middleware
// if something wrote to http, break chain and return
// handle middleware first
// last execute route handler
func (c *Context) Next() {
	if c.Resp.Wrote() {
		return
	}
	if c.mi > len(c.baa.middleware) {
		return
	}
	if c.mi < len(c.baa.middleware) {
		c.baa.middleware[c.mi](c)
	} else {
		c.routeHandler(c)
	}
	c.mi++
}

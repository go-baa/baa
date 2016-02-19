package baa

import (
	"net/http"
	"strconv"
)

// Context provlider a HTTP context for baa
// context contains reqest, response, header, cookie and some content type.
type Context struct {
	Req    *http.Request
	Resp   http.ResponseWriter
	Data   map[string]interface{}
	baa    *Baa
	params map[string]string
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
	c.params = make(map[string]string)
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

// JSON write data by json format
func (c *Context) JSON(code int, d interface{}) error {
	return nil
}

// JSONP write data by jsonp format
func (c *Context) JSONP(code int, d interface{}) error {
	return nil
}

// XML write data by XML format
func (c *Context) XML(code int, d interface{}) error {
	return nil
}

// HTML write data by html template engine, use context.Data
func (c *Context) HTML(code int, tpl string) error {
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

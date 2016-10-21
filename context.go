package baa

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	// defaultMaxMemory Maximum amount of memory to use when parsing a multipart form.
	// Set this to whatever value you prefer; default is 32 MB.
	defaultMaxMemory = 32 << 20 // 32 MB

	// CharsetUTF8 ...
	CharsetUTF8 = "charset=utf-8"

	// MediaTypes
	ApplicationJSON                  = "application/json"
	ApplicationJSONCharsetUTF8       = ApplicationJSON + "; " + CharsetUTF8
	ApplicationJavaScript            = "application/javascript"
	ApplicationJavaScriptCharsetUTF8 = ApplicationJavaScript + "; " + CharsetUTF8
	ApplicationXML                   = "application/xml"
	ApplicationXMLCharsetUTF8        = ApplicationXML + "; " + CharsetUTF8
	ApplicationForm                  = "application/x-www-form-urlencoded"
	ApplicationProtobuf              = "application/protobuf"
	TextHTML                         = "text/html"
	TextHTMLCharsetUTF8              = TextHTML + "; " + CharsetUTF8
	TextPlain                        = "text/plain"
	TextPlainCharsetUTF8             = TextPlain + "; " + CharsetUTF8
	MultipartForm                    = "multipart/form-data"
)

// Context provlider a HTTP context for baa
// context contains reqest, response, header, cookie and some content type.
type Context struct {
	Req      *http.Request
	Resp     *Response
	baa      *Baa
	store    map[string]interface{}
	pNames   []string      // route params names
	pValues  []string      // route params values
	handlers []HandlerFunc // middleware handler and route match handler
	hi       int           // handlers execute position
}

// NewContext create a http context
func NewContext(w http.ResponseWriter, r *http.Request, b *Baa) *Context {
	c := new(Context)
	c.Resp = NewResponse(w, b)
	c.baa = b
	c.pNames = make([]string, 0, 16)
	c.pValues = make([]string, 0, 16)
	c.handlers = make([]HandlerFunc, len(b.middleware), len(b.middleware)+3)
	copy(c.handlers, b.middleware)
	c.Reset(w, r)
	return c
}

// Reset ...
func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	c.Resp.reset(w)
	c.Req = r
	c.hi = 0
	c.handlers = c.handlers[:len(c.baa.middleware)]
	c.pNames = c.pNames[:0]
	c.pValues = c.pValues[:0]
	c.store = nil
}

// Set store data in context
func (c *Context) Set(key string, v interface{}) {
	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	c.store[key] = v
}

// Get returns data from context
func (c *Context) Get(key string) interface{} {
	if c.store == nil {
		return nil
	}
	return c.store[key]
}

// Gets returns data map from content store
func (c *Context) Gets() map[string]interface{} {
	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	return c.store
}

// SetParam read route param value from uri
func (c *Context) SetParam(name, value string) {
	c.pNames = append(c.pNames, name)
	c.pValues = append(c.pValues, value)
}

// Param get route param from context
func (c *Context) Param(name string) string {
	for i := len(c.pNames) - 1; i >= 0; i-- {
		if c.pNames[i] == name {
			return c.pValues[i]
		}
	}
	return ""
}

// Params returns route params from context
func (c *Context) Params() map[string]string {
	m := make(map[string]string)
	for i := 0; i < len(c.pNames); i++ {
		m[c.pNames[i]] = c.pValues[i]
	}
	return m
}

// ParamInt get route param from context and format to int
func (c *Context) ParamInt(name string) int {
	v, _ := strconv.Atoi(c.Param(name))
	return v
}

// ParamInt32 get route param from context and format to int32
func (c *Context) ParamInt32(name string) int32 {
	return int32(c.ParamInt64(name))
}

// ParamInt64 get route param from context and format to int64
func (c *Context) ParamInt64(name string) int64 {
	v, _ := strconv.ParseInt(c.Param(name), 10, 64)
	return v
}

// ParamFloat get route param from context and format to float64
func (c *Context) ParamFloat(name string) float64 {
	v, _ := strconv.ParseFloat(c.Param(name), 64)
	return v
}

// ParamBool get route param from context and format to bool
func (c *Context) ParamBool(name string) bool {
	v, _ := strconv.ParseBool(c.Param(name))
	return v
}

// Query get a param from http.Request.Form
func (c *Context) Query(name string) string {
	c.ParseForm(0)
	return c.Req.Form.Get(name)
}

// QueryTrim querys and trims spaces form parameter.
func (c *Context) QueryTrim(name string) string {
	c.ParseForm(0)
	return strings.TrimSpace(c.Req.Form.Get(name))
}

// QueryStrings get a group param from http.Request.Form and format to string slice
func (c *Context) QueryStrings(name string) []string {
	c.ParseForm(0)
	if v, ok := c.Req.Form[name]; ok {
		return v
	}
	return []string{}
}

// QueryEscape returns escapred query result.
func (c *Context) QueryEscape(name string) string {
	c.ParseForm(0)
	return template.HTMLEscapeString(c.Req.Form.Get(name))
}

// QueryInt get a param from http.Request.Form and format to int
func (c *Context) QueryInt(name string) int {
	c.ParseForm(0)
	v, _ := strconv.Atoi(c.Req.Form.Get(name))
	return v
}

// QueryInt32 get a param from http.Request.Form and format to int32
func (c *Context) QueryInt32(name string) int32 {
	return int32(c.QueryInt64(name))
}

// QueryInt64 get a param from http.Request.Form and format to int64
func (c *Context) QueryInt64(name string) int64 {
	c.ParseForm(0)
	v, _ := strconv.ParseInt(c.Req.Form.Get(name), 10, 64)
	return v
}

// QueryFloat get a param from http.Request.Form and format to float64
func (c *Context) QueryFloat(name string) float64 {
	c.ParseForm(0)
	v, _ := strconv.ParseFloat(c.Req.Form.Get(name), 64)
	return v
}

// QueryBool get a param from http.Request.Form and format to bool
func (c *Context) QueryBool(name string) bool {
	c.ParseForm(0)
	v, _ := strconv.ParseBool(c.Req.Form.Get(name))
	return v
}

// Querys return http.Request.URL queryString data
func (c *Context) Querys() map[string]interface{} {
	params := make(map[string]interface{})
	var newValues url.Values
	if c.Req.URL != nil {
		newValues, _ = url.ParseQuery(c.Req.URL.RawQuery)
	}
	for k, v := range newValues {
		if len(v) > 1 {
			params[k] = v
		} else {
			params[k] = v[0]
		}
	}
	return params
}

// Posts return http.Request form data
func (c *Context) Posts() map[string]interface{} {
	if err := c.ParseForm(0); err != nil {
		return nil
	}
	params := make(map[string]interface{})
	data := c.Req.PostForm
	if len(data) == 0 && len(c.Req.Form) > 0 {
		data = c.Req.Form
	}
	for k, v := range data {
		if len(v) > 1 {
			params[k] = v
		} else {
			params[k] = v[0]
		}
	}
	return params
}

// GetFile returns information about user upload file by given form field name.
func (c *Context) GetFile(name string) (multipart.File, *multipart.FileHeader, error) {
	if err := c.ParseForm(0); err != nil {
		return nil, nil, err
	}
	return c.Req.FormFile(name)
}

// SaveToFile reads a file from request by field name and saves to given path.
func (c *Context) SaveToFile(name, savePath string) error {
	fr, _, err := c.GetFile(name)
	if err != nil {
		return err
	}
	defer fr.Close()

	fw, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fw.Close()

	_, err = io.Copy(fw, fr)
	return err
}

// Body get raw request body and return RequestBody
func (c *Context) Body() *RequestBody {
	return NewRequestBody(c.Req.Body)
}

// SetCookie sets given cookie value to response header.
// full params example:
// SetCookie(<name>, <value>, <max age>, <path>, <domain>, <secure>, <http only>)
func (c *Context) SetCookie(name string, value string, others ...interface{}) {
	cookie := http.Cookie{}
	cookie.Name = name
	cookie.Value = url.QueryEscape(value)

	if len(others) > 0 {
		switch v := others[0].(type) {
		case int:
			cookie.MaxAge = v
		case int64:
			cookie.MaxAge = int(v)
		case int32:
			cookie.MaxAge = int(v)
		}
	}

	cookie.Path = "/"
	if len(others) > 1 {
		if v, ok := others[1].(string); ok && len(v) > 0 {
			cookie.Path = v
		}
	}

	if len(others) > 2 {
		if v, ok := others[2].(string); ok && len(v) > 0 {
			cookie.Domain = v
		}
	}

	if len(others) > 3 {
		switch v := others[3].(type) {
		case bool:
			cookie.Secure = v
		default:
			if others[3] != nil {
				cookie.Secure = true
			}
		}
	}

	if len(others) > 4 {
		if v, ok := others[4].(bool); ok && v {
			cookie.HttpOnly = true
		}
	}

	c.Resp.Header().Add("Set-Cookie", cookie.String())
}

// GetCookie returns given cookie value from request header.
func (c *Context) GetCookie(name string) string {
	cookie, err := c.Req.Cookie(name)
	if err != nil {
		return ""
	}
	v, _ := url.QueryUnescape(cookie.Value)
	return v
}

// GetCookieInt returns cookie result in int type.
func (c *Context) GetCookieInt(name string) int {
	v, _ := strconv.Atoi(c.GetCookie(name))
	return v
}

// GetCookieInt32 returns cookie result in int32 type.
func (c *Context) GetCookieInt32(name string) int32 {
	return int32(c.GetCookieInt64(name))
}

// GetCookieInt64 returns cookie result in int64 type.
func (c *Context) GetCookieInt64(name string) int64 {
	v, _ := strconv.ParseInt(c.GetCookie(name), 10, 64)
	return v
}

// GetCookieFloat64 returns cookie result in float64 type.
func (c *Context) GetCookieFloat64(name string) float64 {
	v, _ := strconv.ParseFloat(c.GetCookie(name), 64)
	return v
}

// GetCookieBool returns cookie result in float64 type.
func (c *Context) GetCookieBool(name string) bool {
	v, _ := strconv.ParseBool(c.GetCookie(name))
	return v
}

// String write text by string
func (c *Context) String(code int, s string) {
	c.Resp.Header().Set("Content-Type", TextPlainCharsetUTF8)
	c.Resp.WriteHeader(code)
	c.Resp.Write([]byte(s))
}

// Text write text by []byte
func (c *Context) Text(code int, s []byte) {
	c.Resp.Header().Set("Content-Type", TextHTMLCharsetUTF8)
	c.Resp.WriteHeader(code)
	c.Resp.Write(s)
}

// JSON write data by json format
func (c *Context) JSON(code int, v interface{}) {
	var re []byte
	var err error
	if c.baa.debug {
		re, err = json.MarshalIndent(v, "", "  ")
	} else {
		re, err = json.Marshal(v)
	}
	if err != nil {
		c.Error(err)
		return
	}

	c.Resp.Header().Set("Content-Type", ApplicationJSONCharsetUTF8)
	c.Resp.WriteHeader(code)
	c.Resp.Write(re)
}

// JSONString return string by Marshal interface
func (c *Context) JSONString(v interface{}) (string, error) {
	var re []byte
	var err error
	if c.baa.debug {
		re, err = json.MarshalIndent(v, "", "  ")
	} else {
		re, err = json.Marshal(v)
	}
	if err != nil {
		return "", err
	}
	return string(re), nil
}

// JSONP write data by jsonp format
func (c *Context) JSONP(code int, callback string, v interface{}) {
	re, err := json.Marshal(v)
	if err != nil {
		c.Error(err)
		return
	}

	c.Resp.Header().Set("Content-Type", ApplicationJavaScriptCharsetUTF8)
	c.Resp.WriteHeader(code)
	c.Resp.Write([]byte(callback + "("))
	c.Resp.Write(re)
	c.Resp.Write([]byte(");"))
}

// XML sends an XML response with status code.
func (c *Context) XML(code int, v interface{}) {
	var re []byte
	var err error
	if c.baa.debug {
		re, err = xml.MarshalIndent(v, "", "  ")
	} else {
		re, err = xml.Marshal(v)
	}
	if err != nil {
		c.Error(err)
		return
	}

	c.Resp.Header().Set("Content-Type", ApplicationXMLCharsetUTF8)
	c.Resp.WriteHeader(code)
	c.Resp.Write([]byte(xml.Header))
	c.Resp.Write(re)
}

// HTML write render data by html template engine use context.store
// it is a alias of c.Render
func (c *Context) HTML(code int, tpl string) {
	c.Render(code, tpl)
}

// Render write render data by html template engine use context.store
func (c *Context) Render(code int, tpl string) {
	re, err := c.Fetch(tpl)
	if err != nil {
		c.Error(err)
		return
	}
	c.Resp.Header().Set("Content-Type", TextHTMLCharsetUTF8)
	c.Resp.WriteHeader(code)
	c.Resp.Write(re)
}

// Fetch render data by html template engine use context.store and returns data
func (c *Context) Fetch(tpl string) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := c.baa.Render().Render(buf, tpl, c.store); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Redirect redirects the request using http.Redirect with status code.
func (c *Context) Redirect(code int, url string) error {
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return fmt.Errorf("invalid redirect status code")
	}
	http.Redirect(c.Resp, c.Req, url, code)
	return nil
}

// RemoteAddr returns more real IP address.
func (c *Context) RemoteAddr() string {
	var addr string
	var key string
	key = "__ctx_remoteAddr"
	if addr, ok := c.Get(key).(string); ok {
		return addr
	}
	addr = c.Req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = c.Req.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = c.Req.RemoteAddr
			addr, _, _ = net.SplitHostPort(addr)
		}
	}
	c.Set(key, addr)
	return addr
}

// Referer returns http request Referer
func (c *Context) Referer() string {
	return c.Req.Header.Get("Referer")
}

// UserAgent returns http request UserAgent
func (c *Context) UserAgent() string {
	return c.Req.Header.Get("User-Agent")
}

// URL returns http request full url
func (c *Context) URL(hasQuery bool) string {
	scheme := c.Req.URL.Scheme
	host := c.Req.URL.Host
	if scheme == "" {
		if c.Req.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	if host == "" {
		host = c.Req.Host
	}
	if len(host) > 0 {
		if host[0] == ':' {
			//
		} else if host[0] == '/' {
			scheme += ":"
		} else {
			scheme += "://"
		}
	} else {
		scheme = ""
	}
	if hasQuery {
		url := c.Req.RequestURI
		if url == "" {
			url = c.Req.URL.Path
			if len(c.Req.URL.RawQuery) > 0 {
				url += "?" + c.Req.URL.RawQuery
			}
		}
		return scheme + host + url
	}
	return scheme + host + c.Req.URL.Path
}

// IsMobile returns if it is a mobile phone device request
func (c *Context) IsMobile() bool {
	userAgent := c.UserAgent()
	for _, v := range []string{"iPhone", "iPod", "Android"} {
		if strings.Contains(userAgent, v) {
			return true
		}
	}
	return false
}

// IsAJAX returns if it is a ajax request
func (c *Context) IsAJAX() bool {
	return c.Req.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// ParseForm parses a request body as multipart/form-data or
//	parses the raw query from the URL and updates r.Form.
func (c *Context) ParseForm(maxSize int64) error {
	if c.Req.Form != nil {
		return nil
	}
	contentType := c.Req.Header.Get("Content-Type")
	if (c.Req.Method == "POST" || c.Req.Method == "PUT") &&
		len(contentType) > 0 && strings.Contains(contentType, MultipartForm) {
		if maxSize == 0 {
			maxSize = defaultMaxMemory
		}
		return c.Req.ParseMultipartForm(maxSize)
	}
	return c.Req.ParseForm()
}

// Next execute next handler
// handle middleware first, last execute route handler
// if something wrote to http, break chain and return
func (c *Context) Next() {
	if c.hi >= len(c.handlers) {
		return
	}
	if c.Resp.Wrote() {
		if c.baa.Debug() {
			c.baa.Logger().Println("Warning: content has been written, handle chain break.")
		}
		return
	}
	i := c.hi
	c.hi++
	c.handlers[i](c)
}

// Break break the handles chain and Immediate return
func (c *Context) Break() {
	c.hi = len(c.handlers)
}

// Error invokes the registered HTTP error handler.
func (c *Context) Error(err error) {
	c.baa.Error(err, c)
}

// NotFound invokes the registered HTTP NotFound handler.
func (c *Context) NotFound() {
	c.baa.NotFound(c)
}

// Baa get app instance
func (c *Context) Baa() *Baa {
	return c.baa
}

// DI get registered dependency injection service
func (c *Context) DI(name string) interface{} {
	return c.baa.GetDI(name)
}

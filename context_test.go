package baa

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContextStore1(t *testing.T) {
	Convey("context store", t, func() {
		b.Get("/context", func(c *Context) {
			c.Get("name")
			c.Gets()
			c.Set("name", "Baa")
			c.Get("name")
			c.Gets()
		})

		req, _ := http.NewRequest("GET", "/context", nil)
		w := httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestContextParam1(t *testing.T) {
	Convey("context route param", t, func() {
		Convey("param", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.Param("id")
				So(id, ShouldEqual, "123")
			})

			req, _ := http.NewRequest("GET", "/context/123", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("param int", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.ParamInt("id")
				So(id, ShouldEqual, 123)
			})

			req, _ := http.NewRequest("GET", "/context/123", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("param int64", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.ParamInt64("id")
				So(id, ShouldEqual, 123)
			})

			req, _ := http.NewRequest("GET", "/context/123", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("param float", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.ParamFloat("id")
				So(id, ShouldEqual, 123.1)
			})

			req, _ := http.NewRequest("GET", "/context/123.1", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("param bool", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.ParamBool("id")
				So(id, ShouldEqual, true)
			})

			req, _ := http.NewRequest("GET", "/context/1", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func TestContextQuery1(t *testing.T) {
	Convey("context query param", t, func() {
		Convey("query string param", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.Query("p")
				So(id, ShouldEqual, "123")
			})

			req, _ := http.NewRequest("GET", "/context/1?p=123", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("form string param", func() {
			b.Post("/context/:id", func(c *Context) {
				id := c.Query("p")
				So(id, ShouldEqual, "123")
			})
			data := url.Values{}
			data.Add("p", "123")
			req, _ := http.NewRequest("POST", "/context/1", strings.NewReader(data.Encode()))
			req.Header.Set(ContentType, ApplicationForm)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("type param", func() {
			b.Post("/context/:id", func(c *Context) {
				var p interface{}
				p = c.QueryInt("int")
				So(p, ShouldEqual, 123)

				p = c.QueryInt64("int64")
				So(p, ShouldEqual, 123)

				p = c.QueryFloat("float")
				So(p, ShouldEqual, 123.4)

				p = c.QueryBool("bool")
				So(p, ShouldEqual, true)

				p = c.QueryBool("bool2")
				So(p, ShouldEqual, false)

				p = c.QueryTrim("trim")
				So(p, ShouldEqual, "abc")

				p = c.QueryStrings("strings")
				So(fmt.Sprintf("%s", p.([]string)), ShouldEqual, "[abc1 abc2]")

				p = c.QueryStrings("strings2")
				So(fmt.Sprintf("%s", p.([]string)), ShouldEqual, "[]")

				p = c.QueryEscape("escape")
				So(p, ShouldEqual, "&lt;a href&gt;string&lt;/a&gt;")
			})
			data := url.Values{}
			data.Add("int", "123")
			data.Add("int64", "123")
			data.Add("float", "123.4")
			data.Add("bool", "1")
			data.Add("bool2", "0")
			data.Add("trim", "abc ")
			data.Add("strings", "abc1")
			data.Add("strings", "abc2")
			data.Add("escape", "<a href>string</a>")
			req, _ := http.NewRequest("POST", "/context/1", strings.NewReader(data.Encode()))
			req.Header.Set(ContentType, ApplicationForm)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("querys/gets, not contains form data", func() {
			b.Post("/context/:id", func(c *Context) {
				querys := c.Querys()
				So(querys, ShouldNotBeNil)
				p := querys["a"].(string)
				So(p, ShouldEqual, "1")
				p = querys["b"].(string)
				So(p, ShouldEqual, "1")
				ps := querys["d"].([]string)
				So(fmt.Sprintf("%s", ps), ShouldEqual, "[1 2]")
			})
			data := url.Values{}
			data.Add("a", "2")
			data.Add("b", "2")
			data.Add("d", "2")
			req, _ := http.NewRequest("POST", "/context/1?a=1&b=1&d=1&d=2", strings.NewReader(data.Encode()))
			req.Header.Set(ContentType, ApplicationForm)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("posts, not contains get params", func() {
			b.Post("/contextp/:id", func(c *Context) {
				querys := c.Posts()
				So(querys, ShouldNotBeNil)
				p := querys["a"].(string)
				So(p, ShouldEqual, "2")
				p = querys["b"].(string)
				So(p, ShouldEqual, "2")
				ps := querys["d"].([]string)
				So(fmt.Sprintf("%s", ps), ShouldEqual, "[2 3]")
			})
			data := url.Values{}
			data.Add("a", "2")
			data.Add("b", "2")
			data.Add("d", "2")
			data.Add("d", "3")
			req, _ := http.NewRequest("POST", "/contextp/1?a=1&b=1&d=1", strings.NewReader(data.Encode()))
			req.Header.Set(ContentType, ApplicationForm)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func TestContextCookie1(t *testing.T) {
	Convey("context cookie", t, func() {
		Convey("cookie get", func() {
			b.Get("/cookie", func(c *Context) {
				var p interface{}
				p = c.GetCookie("s")
				So(p, ShouldEqual, "123")
				p = c.GetCookieInt("int")
				So(p, ShouldEqual, 123)
				p = c.GetCookieInt64("int64")
				So(p, ShouldEqual, 123)
				p = c.GetCookieFloat64("float")
				So(p, ShouldEqual, 123.4)
				p = c.GetCookieBool("bool")
				So(p, ShouldEqual, true)
				p = c.GetCookieBool("bool2")
				So(p, ShouldEqual, false)
				p = c.GetCookie("not")
				So(p, ShouldEqual, "")
			})
			req, _ := http.NewRequest("GET", "/cookie", nil)
			req.Header.Set("Cookie", "s=123; int=123; int64=123; float=123.4; bool=1; boo2=0;")
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("cookie set", func() {
			b.Get("/cookie", func(c *Context) {
				c.SetCookie("name", "baa")
				c.SetCookie("name", "baa", 10)
				c.SetCookie("name", "baa", int32(10))
				c.SetCookie("name", "baa", int64(10))
				c.SetCookie("name", "baa", 10, "/")
				c.SetCookie("name", "baa", 10, "/", "localhost")
				c.SetCookie("name", "baa", 10, "/", "localhost", "1")
				c.SetCookie("name", "baa", 10, "/", "localhost", true, true)
			})
			req, _ := http.NewRequest("GET", "/cookie", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Header().Get("set-cookie"), ShouldEqual, "name=baa; Path=/")
		})
	})
}

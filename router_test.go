package baa

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouteAdd1(t *testing.T) {
	Convey("add static route", t, func() {
		r.add("GET", "/", []HandlerFunc{f})
		r.add("GET", "/bcd", []HandlerFunc{f})
		r.add("GET", "/abcd", []HandlerFunc{f})
		r.add("GET", "/abc", []HandlerFunc{f})
		r.add("GET", "/abd", []HandlerFunc{f})
		r.add("GET", "/abcdef", []HandlerFunc{f})
		r.add("GET", "/bcdefg", []HandlerFunc{f})
		r.add("GET", "/abc/123", []HandlerFunc{f})
		r.add("GET", "/abc/234", []HandlerFunc{f})
		r.add("GET", "/abc/125", []HandlerFunc{f})
		r.add("GET", "/abc/235", []HandlerFunc{f})
		r.add("GET", "/cbd/123", []HandlerFunc{f})
		r.add("GET", "/cbd/234", []HandlerFunc{f})
		r.add("GET", "/cbd/345", []HandlerFunc{f})
		r.add("GET", "/cbd/456", []HandlerFunc{f})
		r.add("GET", "/cbd/346", []HandlerFunc{f})
	})
}

func TestRouteAdd2(t *testing.T) {
	Convey("add param route", t, func() {
		r.add("GET", "/", []HandlerFunc{f})
		r.add("GET", "/a/:id/id", []HandlerFunc{f})
		r.add("GET", "/a/:id/name", []HandlerFunc{f})
		r.add("GET", "/a", []HandlerFunc{f})
		r.add("GET", "/a/:id/", []HandlerFunc{f})
		r.add("GET", "/a/", []HandlerFunc{f})
		r.add("GET", "/a/*/xxx", []HandlerFunc{f})
		r.add("GET", "/p/:project/file/:name", []HandlerFunc{f})
		r.add("GET", "/cbd/:id", []HandlerFunc{f})

		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.add("GET", "/p/:/a", []HandlerFunc{f})
	})
}

func TestRouteAdd3(t *testing.T) {
	Convey("add param route with two different param", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.add("GET", "/a/:id", []HandlerFunc{f})
		r.add("GET", "/a/:name", []HandlerFunc{f})
	})
}

func TestRouteAdd4(t *testing.T) {
	Convey("add route by group", t, func() {
		b.Group("/user", func() {
			b.Get("/info", f)
			b.Get("/info2", f)
			b.Group("/group", func() {
				b.Get("/info", f)
				b.Get("/info2", f)
			})
		})
		b.Group("/user", func() {
			b.Get("/pass", f)
			b.Get("/pass2", f)
		}, f)
	})
}

func TestRouteAdd5(t *testing.T) {
	Convey("add route then set name, URLFor", t, func() {
		b.Get("/article/:id/show", f).Name("articleShow")
		b.Get("/article/:id/detail", f).Name("")
		url := b.URLFor("articleShow", 123)
		So(url, ShouldEqual, "/article/123/show")
		url = b.URLFor("", nil)
		So(url, ShouldEqual, "")
		url = b.URLFor("not exits", "no")
		So(url, ShouldEqual, "")
	})
}

func TestRouteAdd6(t *testing.T) {
	Convey("add route with not support method", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.add("TRACE", "/", []HandlerFunc{f})
	})
}

func TestRouteAdd7(t *testing.T) {
	Convey("add route with empty pattern", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.add("GET", "", []HandlerFunc{f})
	})
}

func TestRouteAdd8(t *testing.T) {
	Convey("add route with pattern that not begin with /", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.add("GET", "abc", []HandlerFunc{f})
	})
}

func TestRouteAdd9(t *testing.T) {
	Convey("other route method", t, func() {
		b2 := New()
		Convey("set auto head route", func() {
			b2.SetAutoHead(true)
			b2.Get("/head", func(c *Context) {
				So(c.Req.Method, ShouldEqual, "HEAD")
			})
			req, _ := http.NewRequest("HEAD", "/head", nil)
			w := httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("set multi method", func() {
			b2.Route("/mul", "GET,HEAD,POST", func(c *Context) {
				c.String(200, "mul")
			})
			req, _ := http.NewRequest("HEAD", "/mul", nil)
			w := httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)

			req, _ = http.NewRequest("GET", "/mul", nil)
			w = httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			req, _ = http.NewRequest("POST", "/mul", nil)
			w = httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("methods", func() {
			b2.Get("/methods", f)
			b2.Patch("/methods", f)
			b2.Post("/methods", f)
			b2.Put("/methods", f)
			b2.Delete("/methods", f)
			b2.Options("/methods", f)
			b2.Head("/methods", f)
			b2.Any("/any", f)
			b2.NotFound(func(c *Context) {
				c.String(404, "baa not found")
			})
		})
	})
}

func TestRouteMatch1(t *testing.T) {
	Convey("match route", t, func() {

		ru := r.match("GET", "/", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/abc/1234", c)
		So(ru, ShouldBeNil)

		ru = r.match("GET", "xxx", c)
		So(ru, ShouldBeNil)

		ru = r.match("GET", "/a/123/id", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/p/yst/file/a.jpg", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/user/info", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/user/pass", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/user/pass32", c)
		So(ru, ShouldBeNil)

		ru = r.match("GET", "/user/xxx", c)
		So(ru, ShouldBeNil)

		ru = r.match("GET", "/xxxx", c)
		So(ru, ShouldBeNil)
	})
}

func TestRoutePrint1(t *testing.T) {
	Convey("print route table", t, func() {
		r.print("", nil)
	})
}

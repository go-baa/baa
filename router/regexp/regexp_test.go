package regexp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-baa/baa"
	. "github.com/smartystreets/goconvey/convey"
)

var b = baa.New()
var r = New(b)
var f = func(c *baa.Context) {}
var c = baa.NewContext(nil, nil, b)

func init() {
	b.SetDI("router", r)
}

// print the route map
func (t *Router) print() {
	for _, m := range t.nodes {
		for _, n := range m {
			fmt.Println(n.pattern)
		}
	}
}
func TestRegexpRouteAdd1(t *testing.T) {
	Convey("add static route", t, func() {
		r.Add("GET", "/", []baa.HandlerFunc{f})
		r.Add("GET", "/bcd", []baa.HandlerFunc{f})
		r.Add("GET", "/abcd", []baa.HandlerFunc{f})
		r.Add("GET", "/abc", []baa.HandlerFunc{f})
		r.Add("GET", "/abd", []baa.HandlerFunc{f})
		r.Add("GET", "/abcdef", []baa.HandlerFunc{f})
		r.Add("GET", "/bcdefg", []baa.HandlerFunc{f})
		r.Add("GET", "/abc/123", []baa.HandlerFunc{f})
		r.Add("GET", "/abc/234", []baa.HandlerFunc{f})
		r.Add("GET", "/abc/125", []baa.HandlerFunc{f})
		r.Add("GET", "/abc/235", []baa.HandlerFunc{f})
		r.Add("GET", "/cbd/123", []baa.HandlerFunc{f})
		r.Add("GET", "/cbd/234", []baa.HandlerFunc{f})
		r.Add("GET", "/cbd/345", []baa.HandlerFunc{f})
		r.Add("GET", "/cbd/456", []baa.HandlerFunc{f})
		r.Add("GET", "/cbd/346", []baa.HandlerFunc{f})
	})
}

func TestRegexpRouteAdd2(t *testing.T) {
	Convey("add param route", t, func() {
		r.Add("GET", "/a/:id/id", []baa.HandlerFunc{f})
		r.Add("GET", "/a/:id/name", []baa.HandlerFunc{f})
		r.Add("GET", "/a", []baa.HandlerFunc{f})
		r.Add("GET", "/a/:id/", []baa.HandlerFunc{f})
		r.Add("GET", "/a/", []baa.HandlerFunc{f})
		r.Add("GET", "/a/*/xxx", []baa.HandlerFunc{f})
		r.Add("GET", "/p/:project/file/:name", []baa.HandlerFunc{f})
		r.Add("GET", "/cbd/:id", []baa.HandlerFunc{f})

		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.Add("GET", "/p/:/a", []baa.HandlerFunc{f})
	})
}

func TestRegexpRouteAdd3(t *testing.T) {
	Convey("add param route with two different param", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.Add("GET", "/a/:id", []baa.HandlerFunc{f})
		r.Add("GET", "/a/:name", []baa.HandlerFunc{f})
	})
}

func TestRegexpRouteAdd4(t *testing.T) {
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
			b.Get("/", f)
			b.Get("/pass", f)
			b.Get("/pass2", f)
		}, f)
	})
}

func TestRegexpRouteAdd5(t *testing.T) {
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

func TestRegexpRouteAdd6(t *testing.T) {
	Convey("add route with not support method", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.Add("TRACE", "/", []baa.HandlerFunc{f})
	})
}

func TestRegexpRouteAdd7(t *testing.T) {
	Convey("add route with empty pattern", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.Add("GET", "", []baa.HandlerFunc{f})
	})
}

func TestRegexpRouteAdd8(t *testing.T) {
	Convey("add route with pattern that not begin with /", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		r.Add("GET", "abc", []baa.HandlerFunc{f})
	})
}

func TestRegexpRouteAdd9(t *testing.T) {
	Convey("other route method", t, func() {
		b2 := baa.New()
		Convey("set auto head route", func() {
			b2.SetAutoHead(true)
			b2.Get("/head", func(c *baa.Context) {
				So(c.Req.Method, ShouldEqual, "HEAD")
			})
			req, _ := http.NewRequest("HEAD", "/head", nil)
			w := httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("set auto training slash", func() {
			b2.SetAutoTrailingSlash(true)
			b2.Get("/slash", func(c *baa.Context) {})
			b2.Group("/slash2", func() {
				b2.Get("/", func(c *baa.Context) {})
				b2.Get("/exist", func(c *baa.Context) {})
			})
			req, _ := http.NewRequest("GET", "/slash", nil)
			w := httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			req, _ = http.NewRequest("GET", "/slash/", nil)
			w = httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			req, _ = http.NewRequest("GET", "/slash2", nil)
			w = httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			req, _ = http.NewRequest("GET", "/slash2/", nil)
			w = httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
			req, _ = http.NewRequest("GET", "/slash2/exist/", nil)
			w = httptest.NewRecorder()
			b2.ServeHTTP(w, req)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("set multi method", func() {
			b2.Route("/mul", "*", func(c *baa.Context) {
				c.String(200, "mul")
			})
			b2.Route("/mul", "GET,HEAD,POST", func(c *baa.Context) {
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
			b2.SetNotFound(func(c *baa.Context) {
				c.String(404, "baa not found")
			})
		})
	})
}

func TestRegexpRouteMatch1(t *testing.T) {
	Convey("match route", t, func() {
		ru := r.Match("GET", "/", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/abc/1234", c)
		So(ru, ShouldBeNil)

		ru = r.Match("GET", "xxx", c)
		So(ru, ShouldBeNil)

		ru = r.Match("GET", "/a/123/id", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/p/yst/file/a.jpg", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/user/info", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/user/pass", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/user/pass32", c)
		So(ru, ShouldBeNil)

		ru = r.Match("GET", "/user/xxx", c)
		So(ru, ShouldBeNil)

		ru = r.Match("GET", "/xxxx", c)
		So(ru, ShouldBeNil)
	})
}

func TestRegexpRouteRegexp1(t *testing.T) {
	Convey("regexp match", t, func() {
		b.Get("/jk/:id/:page.html", f)
		b.Get("/yy/:py([a-zA-Z0-9_-]+)-:id(int)/:page.html", f)
		ru := r.Match("GET", "/jk/123/1.html", c)
		So(ru, ShouldNotBeNil)
		So(c.ParamInt("id"), ShouldEqual, 123)
		So(c.ParamInt("page"), ShouldEqual, 1)
		ru = r.Match("GET", "/yy/bei-yi-san-yuan-123/video.html", c)
		So(ru, ShouldNotBeNil)
		So(c.Param("py"), ShouldEqual, "bei-yi-san-yuan")
		So(c.ParamInt("id"), ShouldEqual, 123)
		So(c.Param("page"), ShouldEqual, "video")
	})
}

func TestRegexpRoutePrint1(t *testing.T) {
	Convey("print route table", t, func() {
		r.(*Router).print()
	})
}

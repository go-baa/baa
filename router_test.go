package baa

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var r = newRouter()
var c = newContext(nil, nil, nil)
var f = func(c *Context) {}

func TestRouteAdd1(t *testing.T) {
	Convey("测试路由添加", t, func() {
		r.add("GET", "/", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/abc", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/bcd", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/abcd", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/abd", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/abcdefg", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/bcdefg", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

	})
}

func TestRouteAdd2(t *testing.T) {
	Convey("测试参数路由添加", t, func() {
		r.add("GET", "/", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/:id", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/:ibb", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/:id/id", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/:ibb/name", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/:project/file/:name", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/*/xxx", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/ab/:name", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/applications/:client_id/tokens/:access_token", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

	})
}

func TestRouteAdd3(t *testing.T) {
	Convey("测试组路由添加", t, func() {
		app := New()
		app.SetRouter(r)
		app.Group("/user", func() {
			app.Get("/info", f)
			app.Get("/info2", f)
		})
		app.Group("/user", func() {
			app.Get("/pass", f)
			app.Get("/pass2", f)
		})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")
	})
}

func TestRoutematch2(t *testing.T) {
	Convey("测试参数路由获取", t, func() {
		ru := r.match("GET", "/", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/a/123/id", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/a/123/name", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/a/yst/file/a.jpg", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/applications/123/tokens/a8sadkfas87jas", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/user/info", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/user/pass", c)
		So(ru, ShouldNotBeNil)

		ru = r.match("GET", "/user/xxx", c)
		So(ru, ShouldBeNil)

		ru = r.match("GET", "/xxxx", c)
		So(ru, ShouldBeNil)
	})
}

package baa

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var r = NewRouter()
var c = NewContext(nil, nil, nil)
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
func TestRouteMatch2(t *testing.T) {
	Convey("测试参数路由获取", t, func() {
		ru := r.Match("GET", "/", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/a/123/id", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/a/123/name", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/a/yst/file/a.jpg", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/applications/123/tokens/a8sadkfas87jas", c)
		So(ru, ShouldNotBeNil)

		ru = r.Match("GET", "/xxxx", c)
		So(ru, ShouldBeNil)
	})
}

package baa

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var r = NewRouter()
var c = NewContext(nil, nil, nil)

func TestRouteAdd1(t *testing.T) {
	Convey("测试路由添加", t, func() {
		f := func(c *Context) error {
			return nil
		}
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
		f := func(c *Context) error {
			return nil
		}
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

	})
}
func TestRouteMatch2(t *testing.T) {
	Convey("测试参数路由获取", t, func() {
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		ru := r.Match("GET", "/", c)
		fmt.Printf("%#v\n", ru)
		fmt.Println(".")

		ru = r.Match("GET", "/a/123/id", c)
		fmt.Printf("%#v\n", ru)
		fmt.Println(".")

		ru = r.Match("GET", "/a/123/name", c)
		fmt.Printf("%#v\n", ru)
		fmt.Println(".")

		ru = r.Match("GET", "/a/yst/file/a.jpg", c)
		fmt.Printf("%#v\n", ru)
		fmt.Println(".")
	})
}

package baa

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouteAdd1(t *testing.T) {
	Convey("测试路由添加", t, func() {
		r := NewRouter()
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
		r := NewRouter()
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

		r.add("GET", "/a/:name", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/:id/id", []HandlerFunc{f})
		r.print("", r.routeMap["GET"])
		fmt.Println(".")

		r.add("GET", "/a/:name/name", []HandlerFunc{f})
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

	})
}

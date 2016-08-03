package router

import (
	"testing"

	"github.com/go-baa/baa"
	. "github.com/smartystreets/goconvey/convey"
)

var b = baa.New()
var r = NewRegexp()
var f = func(c *baa.Context) {}

func init() {
	b.SetDI("router", r)
}

func TestRegexpRouteAdd1(t *testing.T) {
	Convey("test regexp route add", t, func() {
		b.Get("/yiyuan/:py(\\w+)-:id(\\d+).html", f)
		b.Get("/yiyuan/:py.html", f)
		ru := r.Match("GET", "/user/info", nil)
		So(ru, ShouldBeNil)
	})
}

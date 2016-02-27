package baa

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var b = New()
var r = b.router
var c = newContext(nil, nil, b)
var f = func(c *Context) {}

func TestNew1(t *testing.T) {
	Convey("new baa app", t, func() {
		b2 := New()
		So(b2, ShouldNotBeNil)
	})
}

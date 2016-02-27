package baa

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"os"
	"testing"
)

func TestDISetLogger1(t *testing.T) {
	Convey("register error logger", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		b := New()
		b.SetDI("logger", nil)
	})
}

func TestDISetLogger2(t *testing.T) {
	Convey("register correct logger", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldBeNil)
		}()
		b := New()
		b.SetDI("logger", log.New(os.Stderr, "BaaTest ", log.LstdFlags))
	})
}
func TestDISetRender1(t *testing.T) {
	Convey("register error render", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldNotBeNil)
		}()
		b := New()
		b.SetDI("render", nil)
	})
}

func TestDISetRender2(t *testing.T) {
	Convey("register correct render", t, func() {
		defer func() {
			e := recover()
			So(e, ShouldBeNil)
		}()
		b := New()
		b.SetDI("render", newRender())
	})
}

func TestDISet1(t *testing.T) {
	Convey("register di", t, func() {
		b.SetDI("test", "hiDI")
	})
}

func TestDIGet1(t *testing.T) {
	Convey("get registerd di", t, func() {
		var v interface{}
		v = b.GetDI("")
		So(v, ShouldBeNil)
		v = b.GetDI("test")
		So(v, ShouldNotBeNil)
		So(v.(string), ShouldEqual, "hiDI")
	})
}

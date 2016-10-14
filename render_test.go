package baa

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRender1(t *testing.T) {
	Convey("response method", t, func() {
		b.Get("/render", func(c *Context) {
			c.Set("name", "Baa")
			c.HTML(200, "_fixture/index1.html")
		})

		b.Get("/render2", func(c *Context) {
			c.Set("name", "Baa")
			c.HTML(200, "_fixture/index2.html")
		})

		b.Get("/render3", func(c *Context) {
			c.Set("name", "Baa")
			c.HTML(200, "_fixture/index3.html")
		})

		req, _ := http.NewRequest("GET", "/render", nil)
		w := httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusOK)

		req, _ = http.NewRequest("GET", "/render2", nil)
		w = httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)

		req, _ = http.NewRequest("GET", "/render3", nil)
		w = httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusInternalServerError)
	})
}

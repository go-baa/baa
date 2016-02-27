package baa

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestContext1(t *testing.T) {
	Convey("context store", t, func() {
		b.Get("/context", func(c *Context) {
			c.Get("name")
			c.Gets()
			c.Set("name", "Baa")
			c.Get("name")
			c.Gets()
		})

		req, _ := http.NewRequest("GET", "/context", nil)
		w := httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

func TestContext2(t *testing.T) {
	Convey("context route param", t, func() {
		Convey("param", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.Param("id")
				So(id, ShouldEqual, "123")
			})

			req, _ := http.NewRequest("GET", "/context/123", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("param int", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.ParamInt("id")
				So(id, ShouldEqual, 123)
			})

			req, _ := http.NewRequest("GET", "/context/123", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("param int64", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.ParamInt64("id")
				So(id, ShouldEqual, 123)
			})

			req, _ := http.NewRequest("GET", "/context/123", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("param float", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.ParamFloat("id")
				So(id, ShouldEqual, 123.1)
			})

			req, _ := http.NewRequest("GET", "/context/123.1", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("param bool", func() {
			b.Get("/context/:id", func(c *Context) {
				id := c.ParamBool("id")
				So(id, ShouldEqual, true)
			})

			req, _ := http.NewRequest("GET", "/context/1", nil)
			w := httptest.NewRecorder()
			b.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}

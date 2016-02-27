package baa

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStaticRoute1(t *testing.T) {
	Convey("static serve register", t, func() {
		b.Static("/static", "./_fixture", false, f)
		b.Static("/assets/", "./_fixture/", true, f)
	})
}

func TestStaticServe(t *testing.T) {
	Convey("static serve request", t, func() {
		req, _ := http.NewRequest("GET", "/static", nil)
		w := httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusForbidden)

		req, _ = http.NewRequest("GET", "/assets/", nil)
		w = httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusOK)

		req, _ = http.NewRequest("GET", "/assets/img", nil)
		w = httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusFound)

		req, _ = http.NewRequest("GET", "/static/index.html", nil)
		w = httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusMovedPermanently)

		req, _ = http.NewRequest("GET", "/static/favicon.ico", nil)
		w = httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusOK)

		req, _ = http.NewRequest("GET", "/static/notfound", nil)
		w = httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusNotFound)
	})
}

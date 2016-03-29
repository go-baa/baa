package baa

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRequest1(t *testing.T) {
	Convey("request body", t, func() {

		b.Get("/body", func(c *Context) {
			s, err := c.Body().String()
			So(err, ShouldBeNil)
			So(s, ShouldEqual, "{}")

			// because one req can only read once
			b, err := c.Body().Bytes()
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, "")

			r := c.Body().ReadCloser()
			b, err = ioutil.ReadAll(r)
			So(err, ShouldBeNil)
			So(string(b), ShouldEqual, "")
		})
		body := bytes.NewBuffer(nil)
		body.Write([]byte("{}"))
		req, _ := http.NewRequest("GET", "/body", body)
		w := httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

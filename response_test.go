package baa

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestResponse1(t *testing.T) {
	Convey("response method", t, func() {

		b.Get("/response", func(c *Context) {
			c.Resp.Size()
			c.Resp.Write([]byte("1"))
			c.Resp.Write([]byte("2"))
			c.Resp.WriteHeader(200)
			c.Resp.Flush()
			c.Resp.Status()

			func() {
				defer func() {
					e := recover()
					So(e, ShouldNotBeNil)
				}()
				c.Resp.Hijack()
			}()

			func() {
				defer func() {
					e := recover()
					So(e, ShouldNotBeNil)
				}()
				c.Resp.CloseNotify()
			}()
		})
		req, _ := http.NewRequest("GET", "/response", nil)
		w := httptest.NewRecorder()
		b.ServeHTTP(w, req)
		So(w.Code, ShouldEqual, http.StatusOK)
	})
}

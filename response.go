package baa

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// Response implement ResponseWriter
type Response struct {
	wroteHeader bool  // reply header has been (logically) written
	written     int64 // number of bytes written in body
	status      int   // status code passed to WriteHeader
	resp        http.ResponseWriter
	writer      io.Writer
	baa         *Baa
}

// NewResponse ...
func NewResponse(w http.ResponseWriter, b *Baa) *Response {
	r := new(Response)
	r.resp = w
	r.writer = w
	r.baa = b
	return r
}

// Header returns the header map that will be sent by
// WriteHeader. Changing the header after a call to
// WriteHeader (or Write) has no effect unless the modified
// headers were declared as trailers by setting the
// "Trailer" header before the call to WriteHeader (see example).
// To suppress implicit response headers, set their value to nil.
func (r *Response) Header() http.Header {
	return r.resp.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (r *Response) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.writer.Write(b)
	r.written += int64(n)
	return n, err
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (r *Response) WriteHeader(code int) {
	if r.wroteHeader {
		r.baa.Logger().Println("http: multiple response.WriteHeader calls")
		return
	}
	r.wroteHeader = true
	r.status = code
	r.resp.WriteHeader(code)
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Response) Flush() {
	r.resp.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.resp.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotifier interface to allow detecting
// when the underlying connection has gone away.
// This mechanism can be used to cancel long operations on the server if the
// client has disconnected before the response is ready.
// See [http.CloseNotifier](https://golang.org/pkg/net/http/#CloseNotifier)
func (r *Response) CloseNotify() <-chan bool {
	return r.resp.(http.CloseNotifier).CloseNotify()
}

// reset reuse response
func (r *Response) reset(w http.ResponseWriter) {
	r.resp = w
	r.writer = w
	r.wroteHeader = false
	r.written = 0
	r.status = http.StatusOK
}

// Status returns status code
func (r *Response) Status() int {
	return r.status
}

// Size returns body size
func (r *Response) Size() int64 {
	return r.written
}

// Wrote returns if writes something
func (r *Response) Wrote() bool {
	return r.wroteHeader
}

// GetWriter returns response io writer
func (r *Response) GetWriter() io.Writer {
	return r.writer
}

// SetWriter set response io writer
func (r *Response) SetWriter(w io.Writer) {
	r.writer = w
}

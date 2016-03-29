package baa

import (
	"io"
	"io/ioutil"
)

// RequestBody represents a request body.
type RequestBody struct {
	reader io.ReadCloser
}

// NewRequestBody ...
func NewRequestBody(r io.ReadCloser) *RequestBody {
	return &RequestBody{
		reader: r,
	}
}

// Bytes reads and returns content of request body in bytes.
func (rb *RequestBody) Bytes() ([]byte, error) {
	return ioutil.ReadAll(rb.reader)
}

// String reads and returns content of request body in string.
func (rb *RequestBody) String() (string, error) {
	data, err := rb.Bytes()
	return string(data), err
}

// ReadCloser returns a ReadCloser for request body.
func (rb *RequestBody) ReadCloser() io.ReadCloser {
	return rb.reader
}

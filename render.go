package baa

import (
	"io"
)

// Renderer is the interface that wraps the Render method.
type Renderer interface {
	Render(w io.Writer, name string, data interface{}) error
}

// Render default baa template engine
type Render struct {
}

// Render ...
func (r *Render) Render(w io.Writer, name string, data interface{}) error {
	return nil
}

// NewRender create a render instance
func NewRender() *Render {
	r := new(Render)
	return r
}

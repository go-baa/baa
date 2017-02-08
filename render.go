package baa

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
)

// Renderer is the interface that wraps the Render method.
type Renderer interface {
	Render(w io.Writer, tpl string, data interface{}) error
}

// Render default baa template engine
type Render struct {
}

// Render ...
func (r *Render) Render(w io.Writer, tpl string, data interface{}) error {
	t, err := parseFile(tpl)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

// parseFile ...
func parseFile(filename string) (*template.Template, error) {
	var t *template.Template
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s := string(b)
	name := filepath.Base(filename)
	t = template.New(name)
	_, err = t.Parse(s)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// newRender create a render instance
func newRender() *Render {
	r := new(Render)
	return r
}

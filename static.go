package baa

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
)

// compatible with go net standard indexPage
const indexPage = "/index.html"

// Static provider static file serve for baa.
type static struct {
	handler HandlerFunc
	prefix  string
	dir     string
	index   bool
}

// newStatic returns a route handler with static file serve
func newStatic(prefix, dir string, index bool, h HandlerFunc) HandlerFunc {
	if len(prefix) > 1 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	if len(dir) > 1 && dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}
	s := &static{
		dir:     dir,
		index:   index,
		prefix:  prefix,
		handler: h,
	}

	return func(c *Context) {
		file := c.Param("")
		if len(file) > 0 && file[0] == '/' {
			file = file[1:]
		}
		file = s.dir + "/" + file

		if s.handler != nil {
			s.handler(c)
		}

		// directory index
		if f, err := os.Stat(file); err == nil {
			if f.IsDir() {
				if s.index {
					// if no end slash, add slah and redriect
					if c.Req.URL.Path[len(c.Req.URL.Path)-1] != '/' {
						c.Redirect(302, c.Req.URL.Path+"/")
						return
					}
					listDir(file, s, c)
				} else {
					// check index
					if err := serveFile(file+indexPage, c); err != nil {
						c.Resp.WriteHeader(http.StatusForbidden)
					}
				}
				return
			}
		}

		if len(file) >= len(indexPage) && file[len(file)-len(indexPage):] == indexPage {
			if err := serveFile(file, c); err != nil {
				c.Error(err)
			}
		} else {
			http.ServeFile(c.Resp, c.Req, file)
		}
	}
}

// listDir list given dir files
func listDir(dir string, s *static, c *Context) {
	f, err := os.Open(dir)
	if err != nil {
		c.baa.Error(fmt.Errorf("baa.Static listDir Error: %s", err), c)
	}
	defer f.Close()
	fl, err := f.Readdir(-1)
	if err != nil {
		c.baa.Error(fmt.Errorf("baa.Static listDir Error: %s", err), c)
	}

	dirName := f.Name()
	dirName = dirName[len(s.dir):]
	c.Resp.Header().Add("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(c.Resp, "<h3 style=\"padding-bottom:5px;border-bottom:1px solid #ccc;\">%s</h3>\n", dirName)
	fmt.Fprintf(c.Resp, "<pre>\n")
	var color, name string
	for _, v := range fl {
		name = v.Name()
		color = "#333333"
		if v.IsDir() {
			name += "/"
			color = "#3F89C8"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: name}
		fmt.Fprintf(c.Resp, "<a style=\"color:%s\" href=\"%s\">%s</a>\n", color, url.String(), template.HTMLEscapeString(name))
	}
	fmt.Fprintf(c.Resp, "</pre>\n")
}

func serveFile(file string, c *Context) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	fs, err := f.Stat()
	if err != nil {
		return err
	}
	if fs.IsDir() {
		return fmt.Errorf("given path is dir, not file")
	}
	http.ServeContent(c.Resp, c.Req, f.Name(), fs.ModTime(), f)
	return nil
}

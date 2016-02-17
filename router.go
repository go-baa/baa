package baa

import "strings"

var (
	// _HTTP_METHODS Known HTTP methods.
	_HTTP_METHODS = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"OPTIONS": true,
		"HEAD":    true,
	}
)

// Router provlider router for baa
type Router struct {
	autoHead bool
}

// Route is a single route
type Route struct {
}

// Handler provlider handle function
type Handler func(*Context)

// Name set name of route
func (r *Route) Name(name string) {

}

// SetAutoHead sets the value who determines whether add HEAD method automatically
// when GET method is added. Combo router will not be affected by this value.
func (r *Router) SetAutoHead(v bool) {
	r.autoHead = v
}

// Handle registers a new request handle with the given pattern, method and handlers.
func (r *Router) Handle(method string, pattern string, handlers []Handler) *Route {
	return nil
}

// Group registers a list of same prefix route
func (r *Router) Group(pattern string, fn func(), h ...Handler) {

}

// Get is a shortcut for r.Handle("GET", pattern, handlers)
func (r *Router) Get(pattern string, h ...Handler) *Route {
	rs := r.Handle("GET", pattern, h)
	if r.autoHead {
		r.Head(pattern, h...)
	}
	return rs
}

// Patch is a shortcut for r.Handle("PATCH", pattern, handlers)
func (r *Router) Patch(pattern string, h ...Handler) *Route {
	return r.Handle("PATCH", pattern, h)
}

// Post is a shortcut for r.Handle("POST", pattern, handlers)
func (r *Router) Post(pattern string, h ...Handler) *Route {
	return r.Handle("POST", pattern, h)
}

// Put is a shortcut for r.Handle("PUT", pattern, handlers)
func (r *Router) Put(pattern string, h ...Handler) *Route {
	return r.Handle("PUT", pattern, h)
}

// Delete is a shortcut for r.Handle("DELETE", pattern, handlers)
func (r *Router) Delete(pattern string, h ...Handler) *Route {
	return r.Handle("DELETE", pattern, h)
}

// Options is a shortcut for r.Handle("OPTIONS", pattern, handlers)
func (r *Router) Options(pattern string, h ...Handler) *Route {
	return r.Handle("OPTIONS", pattern, h)
}

// Head is a shortcut for r.Handle("HEAD", pattern, handlers)
func (r *Router) Head(pattern string, h ...Handler) *Route {
	return r.Handle("HEAD", pattern, h)
}

// Any is a shortcut for r.Handle("*", pattern, handlers)
func (r *Router) Any(pattern string, h ...Handler) *Route {
	return r.Handle("*", pattern, h)
}

// Route is a shortcut for same handlers but different HTTP methods.
//
// Example:
// 		m.Route("/", "GET,POST", h)
func (r *Router) Route(pattern, methods string, h ...Handler) (route *Route) {
	for _, m := range strings.Split(methods, ",") {
		route = r.Handle(strings.TrimSpace(m), pattern, h)
	}
	return route
}

// NotFound set the route not match result.
// Configurable http.HandlerFunc which is called when no matching route is
// found. If it is not set, http.NotFound is used.
// Be sure to set 404 response code in your handler.
func (r *Router) NotFound(handlers ...Handler) {

}

// InternalServerError set the application accer panic.
// Configurable handler which is called when route handler returns
// error. If it is not set, default handler is used.
// Be sure to set 500 response code in your handler.
func (r *Router) InternalServerError(handlers ...Handler) {

}

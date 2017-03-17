package baa

const (
	GET int = iota
	POST
	PUT
	DELETE
	PATCH
	OPTIONS
	HEAD
	// RouteLength route table length
	RouteLength
)

// RouterMethods declare method key in routeMap
var RouterMethods = map[string]int{
	"GET":     GET,
	"POST":    POST,
	"PUT":     PUT,
	"DELETE":  DELETE,
	"PATCH":   PATCH,
	"OPTIONS": OPTIONS,
	"HEAD":    HEAD,
}

// Router is an router interface for baa
type Router interface {
	// SetAutoHead sets the value who determines whether add HEAD method automatically
	// when GET method is added. Combo router will not be affected by this value.
	SetAutoHead(v bool)
	// SetAutoTrailingSlash optional trailing slash.
	SetAutoTrailingSlash(v bool)
	// Match find matched route then returns handlers and name
	Match(method, uri string, c *Context) ([]HandlerFunc, string)
	// URLFor use named route return format url
	URLFor(name string, args ...interface{}) string
	// Add registers a new handle with the given method, pattern and handlers.
	Add(method, pattern string, handlers []HandlerFunc) RouteNode
	// GroupAdd registers a list of same prefix route
	GroupAdd(pattern string, f func(), handlers []HandlerFunc)
}

// RouteNode is an router node
type RouteNode interface {
	Name(name string)
}

// IsParamChar check the char can used for route params
// a-z->65:90, A-Z->97:122, 0-9->48->57, _->95
func IsParamChar(c byte) bool {
	if (c >= 65 && c <= 90) || (c >= 97 && c <= 122) || (c >= 48 && c <= 57) || c == 95 {
		return true
	}
	return false
}

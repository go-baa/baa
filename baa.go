// Package baa provider a fast & simple Go web framework, routing, middleware, dependency injection, http context.
package baa

const (
	DEV  = "development"
	PROD = "production"
	TEST = "test"
)

// Baa provlider an application
type Baa struct {
}

// Classic create a baa application with default config.
func Classic() *Baa {
	return new(Baa)
}

// New create a baa application without any config.
func New() *Baa {
	return new(Baa)
}

// Run begin a baa application
func (b *Baa) Run(args ...interface{}) {
}

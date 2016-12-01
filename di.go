package baa

import (
	"sync"
)

// DIer is an interface for baa dependency injection
type DIer interface {
	Set(name string, v interface{})
	Get(name string) interface{}
}

// DI provlider a dependency injection service for baa
type DI struct {
	store map[string]interface{}
	rw    sync.RWMutex
}

// NewDI create a DI instance
func NewDI() DIer {
	d := new(DI)
	d.store = make(map[string]interface{})
	return d
}

// Set register a di
// baa dependency injection must be the special interface
func (d *DI) Set(name string, v interface{}) {
	d.rw.Lock()
	d.store[name] = v
	d.rw.Unlock()
}

// Get fetch a di by name, return nil when name not set.
func (d *DI) Get(name string) interface{} {
	d.rw.RLock()
	v := d.store[name]
	d.rw.RUnlock()
	return v
}

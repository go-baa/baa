package baa

// DI provlider a dependency injection service for baa
type DI struct {
	store map[string]interface{}
}

// newDI create a DI instance
func newDI() *DI {
	d := new(DI)
	d.store = make(map[string]interface{})
	return d
}

// set register a di
// baa dependency injection must be the special interface
func (d *DI) set(name string, v interface{}) {
	switch name {
	case "logger":
		if _, ok := v.(Logger); !ok {
			panic("DI logger must be implement interface baa.Logger")
		}
	case "render":
		if _, ok := v.(Renderer); !ok {
			panic("DI render must be implement interface baa.Renderer")
		}
	}
	d.store[name] = v
}

// get fetch a di by name, return nil when name not set.
func (d *DI) get(name string) interface{} {
	return d.store[name]
}

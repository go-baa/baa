package baa

// DI provlider a dependency injection service for baa
type DI struct {
	data map[string]interface{}
}

// newDI create a DI instance
func newDI() *DI {
	d := new(DI)
	d.data = make(map[string]interface{})
	return d
}

// set register a di
func (d *DI) set(name string, v interface{}) {
	if name == "" {
		return
	}
	d.data[name] = v
}

// get fetch a di by name, return nil when name not set.
func (d *DI) get(name string) interface{} {
	if name == "" {
		return nil
	}
	return d.data[name]
}

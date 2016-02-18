package baa

// DI provlider a dependency injection service for baa
type DI struct {
	data map[string]interface{}
}

// NewDI create a DI instance
func NewDI() *DI {
	d := new(DI)
	d.data = make(map[string]interface{})
	return d
}

// Set register a di
func (d *DI) Set(name string, v interface{}) {
	if name == "" {
		return
	}
	d.data[name] = v
}

// Get fetch a di by name, return nil when name not set.
func (d *DI) Get(name string) interface{} {
	if name == "" {
		return nil
	}
	return d.data[name]
}

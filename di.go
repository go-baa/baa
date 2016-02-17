package baa

// DI provlider a dependency injection service for baa
type DI struct {
}

// Set register a di
func (d *DI) Set(name string, v interface{}) {

}

// Get fetch a di by name, return nil when name not set.
func (d *DI) Get(name string) interface{} {
	return nil
}

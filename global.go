package venom

var v *Venom

func init() {
	v = New()
}

// SetDefault sets the provided key and value into the global venom instance at
// the default level
func SetDefault(key string, value interface{}) {
	v.SetDefault(key, value)
}

// SetOverride sets the provided key and value into the global venom instance at
// the override level
func SetOverride(key string, value interface{}) {
	v.SetOverride(key, value)
}

// Get retrievs the requested key from the global venom instance
func Get(key string) interface{} {
	return v.Get(key)
}

// Find searches for the given key, returning the discovered value and a
// boolean indicating whether or not the key was found
func Find(key string) (interface{}, bool) {
	return v.Find(key)
}

// Clear removes all data from the ConfigMap and resets the heap of config
// levels
func Clear() {
	v.Clear()
}

// Debug returns the current venom ConfigMap as a pretty-printed JSON string
func Debug() string {
	return v.Debug()
}

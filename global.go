package venom

var v *Venom

func init() {
	v = Default()
}

// RegisterResolver registers a custom config resolver into the global venom
// instance.
func RegisterResolver(level ConfigLevel, r Resolver) {
	v.RegisterResolver(level, r)
}

// Alias registers an alias for a given key in the global venom instance.
func Alias(from, to string) {
	v.Alias(from, to)
}

// SetLevel is a generic key/value setter method. It sets the provided k/v at
// the specified level inside the global venom instance.
func SetLevel(level ConfigLevel, key string, value interface{}) {
	v.SetLevel(level, key, value)
}

// SetDefault sets the provided key and value into the global venom instance at
// the default level.
func SetDefault(key string, value interface{}) {
	v.SetDefault(key, value)
}

// SetOverride sets the provided key and value into the global venom instance at
// the override level.
func SetOverride(key string, value interface{}) {
	v.SetOverride(key, value)
}

// Get retrieves the requested key from the global venom instance.
func Get(key string) interface{} {
	return v.Get(key)
}

// Find searches for the given key, returning the discovered value and a
// boolean indicating whether or not the key was found.
func Find(key string) (interface{}, bool) {
	return v.Find(key)
}

// LoadFile loads the file from the provided path into Venoms configs. If the
// file can't be opened, if no loader for the files extension exists, or if
// loading the file fails, an error is returned.
func LoadFile(name string) error {
	return v.LoadFile(name)
}

// LoadDirectory loads any config files found in the provided directory,
// optionally recursing into any sub-directories.
func LoadDirectory(dir string, recurse bool) error {
	return v.LoadDirectory(dir, recurse)
}

// Clear removes all data from the ConfigMap and resets the heap of config
// levels.
func Clear() {
	v.Clear()
}

// Debug returns the current venom ConfigMap as a pretty-printed JSON string.
func Debug() string {
	return v.Debug()
}

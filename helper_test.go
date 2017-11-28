package venom

// kv is a test struct containing a (k)ey and a (v)alue
type kv struct {
	k string
	v interface{}
}

// lkv is a test struct containing a config (l)evel, a (k)ey, and a (v)alue
type lkv struct {
	l ConfigLevel
	k string
	v interface{}
}

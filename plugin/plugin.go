package sgc7plugin

// IPlugin - plugin
type IPlugin interface {
	// Random - return [0, r)
	Random(r int) (int, error)
}

package sgc7plugin

import sgc7utils "github.com/zhs007/slotsgamecore7/utils"

// IPlugin - plugin
type IPlugin interface {
	// Random - return [0, r)
	Random(r int) (int, error)
	// GetUsedRngs - get used rngs
	GetUsedRngs() []*sgc7utils.RngInfo
	// ClearUsedRngs - clear used rngs
	ClearUsedRngs()
}

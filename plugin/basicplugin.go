package sgc7plugin

import (
	"math/rand"
	"time"

	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

var isBasicPluginInited = false

// BasicPlugin - basic plugin
type BasicPlugin struct {
	RngUsed []*sgc7utils.RngInfo
	Cache   []int
	Tag     int
}

// NewBasicPlugin - new a BasicPlugin
func NewBasicPlugin() *BasicPlugin {
	bp := &BasicPlugin{
		Tag: -1,
	}

	bp.Init()

	return bp
}

// Random - return [0, r)
func (bp *BasicPlugin) Random(r int) (int, error) {
	var ci int
	if len(bp.Cache) > 0 {
		ci = bp.Cache[0]
		bp.Cache = bp.Cache[1:]
	} else {
		ci = rand.Int()
	}

	cr := ci % r

	bp.AddRngUsed(&sgc7utils.RngInfo{
		Bits:  cr,
		Range: r,
		Value: cr,
	})

	return cr, nil
}

// GetUsedRngs - get used rngs
func (bp *BasicPlugin) GetUsedRngs() []*sgc7utils.RngInfo {
	return bp.RngUsed
}

// ClearUsedRngs - clear used rngs
func (bp *BasicPlugin) ClearUsedRngs() {
	bp.Tag = -1
	bp.RngUsed = nil
}

// AddRngUsed - added used rngs
func (bp *BasicPlugin) AddRngUsed(ri *sgc7utils.RngInfo) {
	bp.RngUsed = append(bp.RngUsed, ri)
}

// SetCache - set cache
func (bp *BasicPlugin) SetCache(arr []int) {
	bp.Cache = arr

	// if arr == nil {
	// 	bp.Cache = nil
	// } else {
	// 	bp.Cache = make([]int, len(arr))
	// 	copy(bp.Cache, arr)
	// }
}

// ClearCache - clear cached rngs
func (bp *BasicPlugin) ClearCache() {
	bp.Cache = nil
}

// Init - initial
func (bp *BasicPlugin) Init() {
	if !isBasicPluginInited {
		rand.Seed(time.Now().UnixNano())

		isBasicPluginInited = true
	}
}

// TagUsedRngs - new a tag for current UsedRngs
func (bp *BasicPlugin) TagUsedRngs() {
	bp.Tag = len(bp.RngUsed)
}

// RollbackUsedRngs - rollback UsedRngs with a tag
func (bp *BasicPlugin) RollbackUsedRngs() error {
	if bp.Tag >= 0 && bp.Tag <= len(bp.RngUsed) {
		bp.RngUsed = bp.RngUsed[0:bp.Tag]

		return nil
	}

	return ErrInvalidTag
}

package sgc7plugin

import (
	"math/rand"
	"time"

	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// BasicPlugin - basic plugin
type BasicPlugin struct {
	RngUsed []*sgc7utils.RngInfo
}

// NewBasicPlugin - new a BasicPlugin
func NewBasicPlugin() *BasicPlugin {
	return &BasicPlugin{}
}

// Random - return [0, r)
func (bp *BasicPlugin) Random(r int) (int, error) {
	ci := rand.Int()
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
	bp.RngUsed = nil
}

// AddRngUsed - added used rngs
func (bp *BasicPlugin) AddRngUsed(ri *sgc7utils.RngInfo) {
	bp.RngUsed = append(bp.RngUsed, ri)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

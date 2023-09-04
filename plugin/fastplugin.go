package sgc7plugin

import (
	"context"

	"github.com/valyala/fastrand"

	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// BasicPlugin - basic plugin
type FastPlugin struct {
	RngUsed []*sgc7utils.RngInfo
	Cache   []int
	Tag     int
	RNG     *fastrand.RNG
}

// NewBasicPlugin - new a BasicPlugin
func NewFastPlugin() *FastPlugin {
	fp := &FastPlugin{
		Tag: -1,
		RNG: &fastrand.RNG{},
	}

	return fp
}

// Random - return [0, r)
func (fp *FastPlugin) Random(ctx context.Context, r int) (int, error) {
	var ci int
	if len(fp.Cache) > 0 {
		ci = fp.Cache[0]
		fp.Cache = fp.Cache[1:]
	} else {
		ci = int(fp.RNG.Uint32())
	}

	cr := ci % r

	fp.AddRngUsed(&sgc7utils.RngInfo{
		Bits:  cr,
		Range: r,
		Value: cr,
	})

	return cr, nil
}

// GetUsedRngs - get used rngs
func (fp *FastPlugin) GetUsedRngs() []*sgc7utils.RngInfo {
	return fp.RngUsed
}

// ClearUsedRngs - clear used rngs
func (fp *FastPlugin) ClearUsedRngs() {
	fp.Tag = -1
	fp.RngUsed = nil
}

// AddRngUsed - added used rngs
func (fp *FastPlugin) AddRngUsed(ri *sgc7utils.RngInfo) {
	fp.RngUsed = append(fp.RngUsed, ri)
}

// SetCache - set cache
func (fp *FastPlugin) SetCache(arr []int) {
	fp.Cache = arr
}

// ClearCache - clear cached rngs
func (fp *FastPlugin) ClearCache() {
	fp.Cache = nil
}

// TagUsedRngs - new a tag for current UsedRngs
func (fp *FastPlugin) TagUsedRngs() {
	fp.Tag = len(fp.RngUsed)
}

// RollbackUsedRngs - rollback UsedRngs with a tag
func (fp *FastPlugin) RollbackUsedRngs() error {
	if fp.Tag >= 0 && fp.Tag <= len(fp.RngUsed) {
		fp.RngUsed = fp.RngUsed[0:fp.Tag]

		return nil
	}

	return ErrInvalidTag
}

package sgc7plugin

import (
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// PluginBase - base
type PluginBase struct {
	RngUsed   []*sgc7utils.RngInfo
	Cache     []int
	Tag       int
	ScenePool any
}

// GetUsedRngs - get used rngs
func (bp *PluginBase) GetUsedRngs() []*sgc7utils.RngInfo {
	return bp.RngUsed
}

// ClearUsedRngs - clear used rngs
func (bp *PluginBase) ClearUsedRngs() {
	bp.Tag = -1
	bp.RngUsed = nil
}

// AddRngUsed - added used rngs
func (bp *PluginBase) AddRngUsed(ri *sgc7utils.RngInfo) {
	bp.RngUsed = append(bp.RngUsed, ri)
}

// SetCache - set cache
func (bp *PluginBase) SetCache(arr []int) {
	bp.Cache = arr
}

// ClearCache - clear cached rngs
func (bp *PluginBase) ClearCache() {
	bp.Cache = nil
}

// TagUsedRngs - new a tag for current UsedRngs
func (bp *PluginBase) TagUsedRngs() {
	bp.Tag = len(bp.RngUsed)
}

// RollbackUsedRngs - rollback UsedRngs with a tag
func (bp *PluginBase) RollbackUsedRngs() error {
	if bp.Tag >= 0 && bp.Tag <= len(bp.RngUsed) {
		bp.RngUsed = bp.RngUsed[0:bp.Tag]

		return nil
	}

	return ErrInvalidTag
}

// SetScenePool - set scene pool
func (bp *PluginBase) SetScenePool(pool any) {
	bp.ScenePool = pool
}

// GetScenePool - get scene pool
func (bp *PluginBase) GetScenePool() any {
	return bp.ScenePool
}

func NewPluginBase() PluginBase {
	return PluginBase{
		Tag: -1,
	}
}

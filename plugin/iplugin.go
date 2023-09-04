package sgc7plugin

import (
	"context"

	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// IPlugin - plugin
type IPlugin interface {
	// Random - return [0, r)
	Random(ctx context.Context, r int) (int, error)
	// GetUsedRngs - get used rngs
	GetUsedRngs() []*sgc7utils.RngInfo
	// ClearUsedRngs - clear used rngs
	ClearUsedRngs()
	// TagUsedRngs - new a tag for current UsedRngs
	TagUsedRngs()
	// RollbackUsedRngs - rollback UsedRngs with the latest tag
	RollbackUsedRngs() error
	// SetCache - set cache
	SetCache(arr []int)
	// ClearCache - clear cached rngs
	ClearCache()
	// Init - initial
	Init()
	// SetScenePool - set scene pool
	SetScenePool(any)
	// GetScenePool - get scene pool
	GetScenePool() any
}

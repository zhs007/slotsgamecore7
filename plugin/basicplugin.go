package sgc7plugin

import (
	"context"
	"math/rand"
	"time"

	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

var isBasicPluginInited = false

// BasicPlugin - basic plugin
type BasicPlugin struct {
	PluginBase
}

// NewBasicPlugin - new a BasicPlugin
func NewBasicPlugin() *BasicPlugin {
	bp := &BasicPlugin{
		PluginBase: NewPluginBase(),
	}

	bp.Init()

	return bp
}

// Random - return [0, r)
func (bp *BasicPlugin) Random(ctx context.Context, r int) (int, error) {
	if IsNoRNGCache {
		ci := rand.Int()

		cr := ci % r

		return cr, nil
	}

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

// Init - initial
func (bp *BasicPlugin) Init() {
	if !isBasicPluginInited {
		rand.Seed(time.Now().UnixNano())

		isBasicPluginInited = true
	}
}

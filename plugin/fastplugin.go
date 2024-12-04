package sgc7plugin

import (
	"context"
	"time"

	"github.com/valyala/fastrand"

	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// BasicPlugin - basic plugin
type FastPlugin struct {
	PluginBase
	RNG *fastrand.RNG
}

// NewFastPlugin - new a BasicPlugin
func NewFastPlugin() *FastPlugin {
	fp := &FastPlugin{
		PluginBase: NewPluginBase(),
		RNG:        &fastrand.RNG{},
	}

	fp.Init()

	return fp
}

// Random - return [0, r)
func (fp *FastPlugin) Random(ctx context.Context, r int) (int, error) {
	if IsNoRNGCache {
		ci := int(fp.RNG.Uint32())

		cr := ci % r

		fp.AddRngUsed(&sgc7utils.RngInfo{
			Bits:  cr,
			Range: r,
			Value: cr,
		})

		return cr, nil
	}

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

// Init - initial
func (fp *FastPlugin) Init() {
	fp.RNG.Seed(uint32(time.Now().UnixNano()))
}

// SetSeed - set a seed
func (fp *FastPlugin) SetSeed(seed int) {

}

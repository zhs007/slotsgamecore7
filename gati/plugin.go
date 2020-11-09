package gati

import (
	"context"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// PluginGATI - plugin for GATI
type PluginGATI struct {
	sgc7plugin.BasicPlugin

	Cfg  *Config
	Rngs []int
}

// NewPluginGATI - new PluginGATI (IPlugin)
func NewPluginGATI(cfg *Config) *PluginGATI {
	return &PluginGATI{
		Cfg: cfg,
	}
}

// Random - return [0, r)
func (plugin *PluginGATI) Random(ctx context.Context, r int) (int, error) {
	if len(plugin.Rngs) == 0 {
		rngs, err := GetRngs(plugin.Cfg.RNGURL, plugin.Cfg.GameID, plugin.Cfg.RngNums)
		if err != nil {
			return -1, err
		}

		plugin.Rngs = rngs
	}

	cv := plugin.Rngs[0]
	plugin.Rngs = plugin.Rngs[1:]

	v := cv % r

	plugin.AddRngUsed(&sgc7utils.RngInfo{
		Bits:  cv,
		Range: r,
		Value: v,
	})

	return v, nil
}

// SetCache - set cache
func (plugin *PluginGATI) SetCache(arr []int) {
	plugin.Rngs = arr
}

// ClearCache - clear cached rngs
func (plugin *PluginGATI) ClearCache() {
	plugin.Rngs = nil
}

// Init - initial
func (plugin *PluginGATI) Init() {
}

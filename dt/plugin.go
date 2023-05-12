package dt

import (
	"context"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// PluginDT - plugin for DT
type PluginDT struct {
	sgc7plugin.BasicPlugin

	Rngs      []uint32
	RngClient *RngClient
}

// NewPluginDT - new PluginDT (IPlugin)
func NewPluginDT(rngServAddr string, gameCode string, useOpenTelemetry bool) *PluginDT {
	return &PluginDT{
		RngClient: NewRngClient(rngServAddr, gameCode, useOpenTelemetry),
	}
}

// Random - return [0, r)
func (plugin *PluginDT) Random(ctx context.Context, r int) (int, error) {
	if len(plugin.Rngs) == 0 {
		rngs, err := plugin.RngClient.GetRngs(ctx, 0)
		if err != nil {
			return -1, err
		}

		plugin.Rngs = rngs
	}

	maxval := int64(r)

	cr := 0
	MAX_RANGE := int64(1) << 32
	limit := MAX_RANGE - (MAX_RANGE % maxval)

	for {
		if len(plugin.Rngs) == 0 {
			rngs, err := plugin.RngClient.GetRngs(ctx, 0)
			if err != nil {
				return -1, err
			}

			plugin.Rngs = rngs
		}

		cr = int(plugin.Rngs[0])
		plugin.Rngs = plugin.Rngs[1:]

		if int64(cr) < limit {
			break
		}
	}

	v := cr % r

	plugin.AddRngUsed(&sgc7utils.RngInfo{
		Bits:  cr,
		Range: r,
		Value: v,
	})

	return v, nil
}

// SetCache - set cache
func (plugin *PluginDT) SetCache(arr []int) {
	plugin.Rngs = nil

	for _, v := range arr {
		plugin.Rngs = append(plugin.Rngs, uint32(v))
	}
}

// ClearCache - clear cached rngs
func (plugin *PluginDT) ClearCache() {
	plugin.Rngs = nil
}

// Init - initial
func (plugin *PluginDT) Init() {
}

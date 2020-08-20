package gati

// PluginGATI - plugin for GATI
type PluginGATI struct {
	Cfg     *Config
	Rngs    []int
	RngUsed []*RngInfo
}

// NewPluginGATI - new PluginGATI (IPlugin)
func NewPluginGATI(cfg *Config) *PluginGATI {
	return &PluginGATI{
		Cfg: cfg,
	}
}

// Random - return [0, r)
func (plugin *PluginGATI) Random(r int) (int, error) {
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

	plugin.RngUsed = append(plugin.RngUsed, &RngInfo{
		Bits:  cv,
		Range: r,
		Value: v,
	})

	return v, nil
}

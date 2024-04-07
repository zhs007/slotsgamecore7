package lowcode

import (
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func NewBasicRNG() IRNG {
	return &BasicRNG{}
}

type BasicRNG struct {
	curPlugin sgc7plugin.IPlugin
}

// OnNewGame -
func (rng *BasicRNG) OnNewGame(plugin sgc7plugin.IPlugin) error {
	rng.curPlugin = plugin

	return nil
}

// GetCurRNG -
func (rng *BasicRNG) GetCurRNG(componentName string) sgc7plugin.IPlugin {
	return rng.curPlugin
}

package lowcode

import (
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type SimpleRNG struct {
	IterateComponent string
	curPlugin        sgc7plugin.IPlugin
	plugin           *sgc7plugin.FastPlugin
	weights          []int
	curIndex         int
}

func (rng *SimpleRNG) Clone() *SimpleRNG {
	return &SimpleRNG{
		IterateComponent: rng.IterateComponent,
		plugin:           rng.plugin,
	}
}

// OnNewGame -
func (rng *SimpleRNG) OnNewGame(plugin sgc7plugin.IPlugin) error {
	rng.curPlugin = plugin

	return nil
}

// GetCurRNG -
func (rng *SimpleRNG) GetCurRNG(gameProp *GameProperty, curComponent IComponent, cd IComponentData, fl IFeatureLevel) (bool, int, sgc7plugin.IPlugin, string) {
	if curComponent.GetName() == rng.IterateComponent {
		if curComponent.GetBranchNum() > 0 {
			if len(rng.weights) == 0 {
				rng.weights = curComponent.GetBranchWeights()
				rng.curIndex = 0
			} else {
				if rng.curIndex < len(rng.weights) {
					cd.ForceBranch(rng.curIndex)

					rng.curIndex++
				}
			}

			return false, -1, rng.plugin, ""
		}
	}

	return false, -1, rng.curPlugin, ""
}

func (rng *SimpleRNG) IsNeedIterate() bool {
	return len(rng.weights) > 0
}

func (rng *SimpleRNG) IsIterateEnding() bool {
	return rng.curIndex < len(rng.weights)
}

func NewSimpleRNG(iterateComponent string) IRNG {
	return &SimpleRNG{
		IterateComponent: iterateComponent,
		plugin:           sgc7plugin.NewFastPlugin(),
	}
}

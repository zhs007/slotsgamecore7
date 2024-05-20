package lowcode

import (
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type FuncNewRNG func() IRNG

type IRNG interface {
	// OnNewGame -
	OnNewGame(plugin sgc7plugin.IPlugin) error
	// GetCurRNG -
	GetCurRNG(componentName string) (bool, int, sgc7plugin.IPlugin, string)
}

package lowcode

import (
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
)

type FuncNewComponent func(name string) IComponent

type IComponent interface {
	// Init -
	Init(fn string, gameProp *GameProperty) error
	// OnNewGame -
	OnNewGame(gameProp *GameProperty) error
	// OnNewStep -
	OnNewStep(gameProp *GameProperty) error
	// OnPlayGame - on playgame
	OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
		cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error
	// OnAsciiGame - outpur to asciigame
	OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error
	// OnStats
	OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64)
}

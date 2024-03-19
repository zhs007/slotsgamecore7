package lowcode

import (
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/types/known/anypb"
)

type FuncNewComponent func(name string) IComponent

type ForeachSymbolData struct {
	SymbolCode int
	Index      int
}

type IComponent interface {
	// Init -
	Init(fn string, pool *GamePropertyPool) error
	// InitEx -
	InitEx(cfg any, pool *GamePropertyPool) error
	// OnGameInited - on game inited
	OnGameInited(components *ComponentList) error

	// OnPlayGame - on playgame
	OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
		cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error)
	// OnAsciiGame - outpur to asciigame
	OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error
	// OnStats -
	OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64)
	// NewComponentData -
	NewComponentData() IComponentData

	// EachUsedResults -
	EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult)
	// OnPlayGame - on playgame
	OnPlayGameEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
		cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) error
	// GetName -
	GetName() string

	// IsRespin -
	IsRespin() bool
	// IsForeach -
	IsForeach() bool

	// NewStats2 -
	NewStats2(parent string) *stats2.Feature
	// OnStats2
	OnStats2(icd IComponentData, s2 *stats2.Cache)

	// GetAllLinkComponents - get all link components
	GetAllLinkComponents() []string

	// GetNextLinkComponents - get next link components
	GetNextLinkComponents() []string
	// GetChildLinkComponents - get child link components
	GetChildLinkComponents() []string

	// CanTriggerWithScene -
	CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result)

	//----------------------------
	// for mask

	// IsMask -
	IsMask() bool
	// SetMask -
	SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error
	// SetMaskVal -
	SetMaskVal(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, index int, mask bool) error
	// SetMaskOnlyTrue -
	SetMaskOnlyTrue(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error

	//----------------------------
	// for foreach symbols

	// EachSymbols - each symbols
	EachSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState, stake *sgc7game.Stake,
		prs []*sgc7game.PlayResult, cd IComponentData) error

	//----------------------------
	// PositionCollection

	// AddPos -
	AddPos(cd IComponentData, x int, y int)
}

package lowcode

import (
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
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
	// // OnStats -
	// OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64)
	// NewComponentData -
	NewComponentData() IComponentData

	// EachUsedResults -
	EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult)
	// ProcRespinOnStepEnd - 现在只有respin需要特殊处理结束，如果多层respin嵌套时，只要新的有next，就不会继续结束respin
	ProcRespinOnStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, canRemove bool) (string, error)
	// GetName -
	GetName() string

	// IsRespin -
	IsRespin() bool
	// IsForeach -
	IsForeach() bool

	// NewStats2 -
	NewStats2(parent string) *stats2.Feature
	// OnStats2 - 除respin外，其它component都是在onPlayGame后调用；respin会在onStepEnd这个环节调用，而且是遍历Respin队列
	OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool)
	// IsNeedOnStepEndStats2 - 除respin外，如果也有component也需要在stepEnd调用的话，这里需要返回true
	IsNeedOnStepEndStats2() bool

	// GetAllLinkComponents - get all link components
	GetAllLinkComponents() []string

	// GetNextLinkComponents - get next link components
	GetNextLinkComponents() []string
	// GetChildLinkComponents - get child link components
	GetChildLinkComponents() []string

	// CanTriggerWithScene -
	CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, icd IComponentData) (bool, []*sgc7game.Result)

	// ProcControllers -
	ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string)

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

	//----------------------------
	// for Set

	// OnPlayGameWithSet - on playgame with a set
	OnPlayGameWithSet(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
		cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData, set int) (string, error)

	//----------------------------
	// Branch

	// GetBranchNum -
	GetBranchNum() int
	// GetBranchWeights -
	GetBranchWeights() []int

	//----------------------------
	// IComponentData

	// ClearData -
	ClearData(icd IComponentData, bForceNow bool)

	//----------------------------
	// PlayerState

	// InitPlayerState -
	// 2 种调用时机，一个是玩家第一次初始化时，这时 bet 为 0，gameProp 和 plugin 为 nil
	// 另外一个是玩家下注时，这时 bet、gameProp、plugin 都有效
	InitPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin, ps *PlayerState, betMethod int, bet int) error

	// NewPlayerState - new IComponentPS
	NewPlayerState() IComponentPS
}

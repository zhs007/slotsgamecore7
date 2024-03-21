package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"google.golang.org/protobuf/proto"
)

type FuncOnEachUsedResult func(*sgc7game.Result)

type IComponentData interface {
	// OnNewGame -
	OnNewGame(gameProp *GameProperty, component IComponent)
	// // OnNewStep -
	// OnNewStep(gameProp *GameProperty, component IComponent)
	// BuildPBComponentData
	BuildPBComponentData() proto.Message

	// GetVal -
	GetVal(key string) int
	// SetVal -
	SetVal(key string, val int)

	// GetConfigVal -
	GetConfigVal(key string) string
	// SetConfigVal -
	SetConfigVal(key string, val string)
	// GetConfigIntVal -
	GetConfigIntVal(key string) (int, bool)
	// SetConfigIntVal -
	SetConfigIntVal(key string, val int)
	// ChgConfigIntVal -
	ChgConfigIntVal(key string, off int)
	// ClearConfigIntVal -
	ClearConfigIntVal(key string)

	// GetResults -
	GetResults() []int
	// GetOutput -
	GetOutput() int
	// GetStringOutput -
	GetStringOutput() string

	//----------------------------
	// SymbolCollection

	// GetSymbols -
	GetSymbols() []int
	// AddSymbol -
	AddSymbol(symbolCode int)

	//----------------------------
	// PositionCollection

	// GetPos -
	GetPos() []int
	// HasPos -
	HasPos(x int, y int) bool
	// AddPos -
	AddPos(x int, y int)

	//----------------------------
	// Respin

	// GetLastRespinNum -
	GetLastRespinNum() int
	// IsRespinEnding -
	IsRespinEnding() bool
	// IsRespinStarted -
	IsRespinStarted() bool
	// AddTriggerRespinAward -
	AddTriggerRespinAward(award *Award)
	// AddRespinTimes -
	AddRespinTimes(num int)
	// TriggerRespin
	TriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams)
	// PushTriggerRespin -
	PushTriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int)

	//----------------------------
	// Mask

	// GetMask -
	GetMask() []bool
	// ChgMask -
	ChgMask(curMask int, val bool) bool

	//----------------------------
	// PiggyBank

	PutInMoney(coins int)
}

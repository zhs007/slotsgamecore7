package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"google.golang.org/protobuf/proto"
)

type FuncOnEachUsedResult func(*sgc7game.Result)
type FuncNewComponentData func() IComponentData

type GetComponentValType int

const (
	GCVTypeNormal GetComponentValType = 0
	GCVTypeMin    GetComponentValType = 1
	GCVTypeMax    GetComponentValType = 2
)

type IComponentData interface {
	// OnNewGame -
	OnNewGame(gameProp *GameProperty, component IComponent)
	// BuildPBComponentData
	BuildPBComponentData() proto.Message
	// // LoadPB
	// LoadPB(pb *anypb.Any) error
	// Clone
	Clone() IComponentData

	// GetValEx -
	GetValEx(key string, getType GetComponentValType) (int, bool)

	// GetStrVal -
	GetStrVal(key string) (string, bool)

	// GetConfigVal -
	GetConfigVal(key string) string
	// SetConfigVal -
	SetConfigVal(key string, val string)
	// GetConfigIntVal -
	GetConfigIntVal(key string) (int, bool)
	// SetConfigIntVal -
	SetConfigIntVal(key string, val int)
	// ChgConfigIntVal -
	ChgConfigIntVal(key string, off int) int
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

	//----------------------------
	// ReelsCollector

	ChgReelsCollector(reelsData []int)

	// //----------------------------
	// // Branch

	// // ForceBranch -
	// ForceBranch(branchIndex int)
}

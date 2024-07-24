package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type BasicComponentData struct {
	UsedScenes            []int
	UsedOtherScenes       []int
	UsedResults           []int
	UsedPrizeScenes       []int
	CashWin               int64
	CoinWin               int
	TargetSceneIndex      int
	TargetOtherSceneIndex int
	RNG                   []int
	MapConfigVals         map[string]string
	MapConfigIntVals      map[string]int
	SrcScenes             []int
	Output                int
	StrOutput             string
	ForceBranchIndex      int
}

// Clone
func (basicComponentData *BasicComponentData) CloneBasicComponentData() BasicComponentData {
	target := BasicComponentData{
		CashWin:               basicComponentData.CashWin,
		CoinWin:               basicComponentData.CoinWin,
		TargetSceneIndex:      basicComponentData.TargetSceneIndex,
		TargetOtherSceneIndex: basicComponentData.TargetOtherSceneIndex,
		MapConfigVals:         make(map[string]string),
		MapConfigIntVals:      make(map[string]int),
		Output:                basicComponentData.Output,
		StrOutput:             basicComponentData.StrOutput,
	}

	target.UsedScenes = make([]int, len(basicComponentData.UsedScenes))
	copy(target.UsedScenes, basicComponentData.UsedScenes)

	target.UsedOtherScenes = make([]int, len(basicComponentData.UsedOtherScenes))
	copy(target.UsedOtherScenes, basicComponentData.UsedOtherScenes)

	target.UsedResults = make([]int, len(basicComponentData.UsedResults))
	copy(target.UsedResults, basicComponentData.UsedResults)

	target.UsedPrizeScenes = make([]int, len(basicComponentData.UsedPrizeScenes))
	copy(target.UsedPrizeScenes, basicComponentData.UsedPrizeScenes)

	target.RNG = make([]int, len(basicComponentData.RNG))
	copy(target.RNG, basicComponentData.RNG)

	for k, v := range basicComponentData.MapConfigVals {
		target.MapConfigVals[k] = v
	}

	for k, v := range basicComponentData.MapConfigIntVals {
		target.MapConfigIntVals[k] = v
	}

	target.SrcScenes = make([]int, len(basicComponentData.SrcScenes))
	copy(target.SrcScenes, basicComponentData.SrcScenes)

	if !gIsReleaseMode {
		target.ForceBranchIndex = basicComponentData.ForceBranchIndex
	}

	return target
}

// Clone
func (basicComponentData *BasicComponentData) Clone() IComponentData {
	target := basicComponentData.CloneBasicComponentData()

	return &target
}

// OnNewGame -
func (basicComponentData *BasicComponentData) OnNewGame(gameProp *GameProperty, component IComponent) {
	basicComponentData.MapConfigVals = make(map[string]string)
	basicComponentData.MapConfigIntVals = make(map[string]int)
}

// // OnNewStep -
// func (basicComponentData *BasicComponentData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	basicComponentData.UsedScenes = nil
// 	basicComponentData.UsedOtherScenes = nil
// 	basicComponentData.UsedResults = nil
// 	basicComponentData.UsedPrizeScenes = nil
// 	basicComponentData.CashWin = 0
// 	basicComponentData.CoinWin = 0
// 	basicComponentData.TargetSceneIndex = -1
// 	basicComponentData.TargetOtherSceneIndex = -1
// 	basicComponentData.RNG = nil

// 	basicComponentData.initSrcScenes()
// }

// GetVal -
func (basicComponentData *BasicComponentData) GetVal(key string) (int, bool) {
	return 0, false
}

// SetVal -
func (basicComponentData *BasicComponentData) SetVal(key string, val int) {

}

// GetConfigVal -
func (basicComponentData *BasicComponentData) GetConfigVal(key string) string {
	return basicComponentData.MapConfigVals[key]
}

// SetConfigVal -
func (basicComponentData *BasicComponentData) SetConfigVal(key string, val string) {
	basicComponentData.MapConfigVals[key] = val
}

// GetConfigIntVal -
func (basicComponentData *BasicComponentData) GetConfigIntVal(key string) (int, bool) {
	ival, isok := basicComponentData.MapConfigIntVals[key]
	return ival, isok
}

// SetConfigIntVal -
func (basicComponentData *BasicComponentData) SetConfigIntVal(key string, val int) {
	basicComponentData.MapConfigIntVals[key] = val
}

// ChgConfigIntVal -
func (basicComponentData *BasicComponentData) ChgConfigIntVal(key string, off int) {
	basicComponentData.MapConfigIntVals[key] += off
}

// ClearConfigIntVal -
func (basicComponentData *BasicComponentData) ClearConfigIntVal(key string) {
	delete(basicComponentData.MapConfigIntVals, key)
}

// // InitSrcScenes -
// func (basicComponentData *BasicComponentData) initSrcScenes() {
// 	for i := range basicComponentData.SrcScenes {
// 		basicComponentData.SrcScenes[i] = -1
// 	}
// }

// BuildPBComponentData
func (basicComponentData *BasicComponentData) BuildPBComponentData() proto.Message {
	return &sgc7pb.BasicComponentData{
		BasicComponentData: basicComponentData.BuildPBBasicComponentData(),
	}
}

// LoadPB
func (basicComponentData *BasicComponentData) LoadPB(pb *anypb.Any) error {
	if pb.TypeUrl == "type.googleapis.com/sgc7pb.ComponentData" {
		var msg sgc7pb.ComponentData

		err := anypb.UnmarshalTo(pb, &msg, proto.UnmarshalOptions{})
		if err != nil {
			goutils.Error("BasicComponentData.LoadPB:UnmarshalTo:ComponentData",
				goutils.Err(err))

			return err
		}

		return nil
	}

	goutils.Error("BasicComponentData.LoadPB",
		goutils.Err(ErrInvalidPBComponentData))

	return ErrInvalidPBComponentData
}

// LoadPB
func (basicComponentData *BasicComponentData) LoadPBComponentData(pb *sgc7pb.ComponentData) error {
	basicComponentData.CashWin = pb.CashWin
	basicComponentData.CoinWin = int(pb.CoinWin)
	basicComponentData.TargetSceneIndex = int(pb.TargetScene)
	basicComponentData.Output = int(pb.Output)
	basicComponentData.StrOutput = pb.StrOutput

	basicComponentData.UsedOtherScenes = nil
	for _, v := range pb.UsedOtherScenes {
		basicComponentData.UsedOtherScenes = append(basicComponentData.UsedOtherScenes, int(v))
	}

	basicComponentData.UsedScenes = nil
	for _, v := range pb.UsedScenes {
		basicComponentData.UsedScenes = append(basicComponentData.UsedScenes, int(v))
	}

	basicComponentData.UsedResults = nil
	for _, v := range pb.UsedResults {
		basicComponentData.UsedResults = append(basicComponentData.UsedResults, int(v))
	}

	basicComponentData.UsedPrizeScenes = nil
	for _, v := range pb.UsedPrizeScenes {
		basicComponentData.UsedPrizeScenes = append(basicComponentData.UsedPrizeScenes, int(v))
	}

	basicComponentData.SrcScenes = nil
	for _, v := range pb.SrcScenes {
		basicComponentData.SrcScenes = append(basicComponentData.SrcScenes, int(v))
	}

	return nil
}

// BuildPBBasicComponentData
func (basicComponentData *BasicComponentData) BuildPBBasicComponentData() *sgc7pb.ComponentData {
	pbcd := &sgc7pb.ComponentData{}

	pbcd.CashWin = basicComponentData.CashWin
	pbcd.CoinWin = int32(basicComponentData.CoinWin)
	pbcd.TargetScene = int32(basicComponentData.TargetSceneIndex)
	pbcd.Output = int32(basicComponentData.Output)
	pbcd.StrOutput = basicComponentData.StrOutput

	for _, v := range basicComponentData.UsedOtherScenes {
		pbcd.UsedOtherScenes = append(pbcd.UsedOtherScenes, int32(v))
	}

	for _, v := range basicComponentData.UsedScenes {
		pbcd.UsedScenes = append(pbcd.UsedScenes, int32(v))
	}

	for _, v := range basicComponentData.UsedResults {
		pbcd.UsedResults = append(pbcd.UsedResults, int32(v))
	}

	for _, v := range basicComponentData.UsedPrizeScenes {
		pbcd.UsedPrizeScenes = append(pbcd.UsedPrizeScenes, int32(v))
	}

	for _, v := range basicComponentData.SrcScenes {
		pbcd.SrcScenes = append(pbcd.SrcScenes, int32(v))
	}

	return pbcd
}

// GetResults -
func (basicComponentData *BasicComponentData) GetResults() []int {
	return basicComponentData.UsedResults
}

// GetOutput -
func (basicComponentData *BasicComponentData) GetOutput() int {
	return basicComponentData.Output
}

// GetStringOutput -
func (basicComponentData *BasicComponentData) GetStringOutput() string {
	return basicComponentData.StrOutput
}

// GetSymbols -
func (basicComponentData *BasicComponentData) GetSymbols() []int {
	return nil
}

// AddSymbol -
func (basicComponentData *BasicComponentData) AddSymbol(symbolCode int) {

}

// GetPos -
func (basicComponentData *BasicComponentData) GetPos() []int {
	return nil
}

// HasPos -
func (basicComponentData *BasicComponentData) HasPos(x int, y int) bool {
	return false
}

// AddPos -
func (basicComponentData *BasicComponentData) AddPos(x int, y int) {
}

// GetLastRespinNum -
func (basicComponentData *BasicComponentData) GetLastRespinNum() int {
	return 0
}

// IsRespinEnding -
func (basicComponentData *BasicComponentData) IsRespinEnding() bool {
	return false
}

// IsRespinStarted -
func (basicComponentData *BasicComponentData) IsRespinStarted() bool {
	return false
}

// // AddRetriggerRespinNum -
// func (basicComponentData *BasicComponentData) AddRetriggerRespinNum(num int) {

// }

// AddTriggerRespinAward -
func (basicComponentData *BasicComponentData) AddTriggerRespinAward(award *Award) {

}

// AddRespinTimes -
func (basicComponentData *BasicComponentData) AddRespinTimes(num int) {

}

// TriggerRespin
func (basicComponentData *BasicComponentData) TriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {

}

// PushTrigger -
func (basicComponentData *BasicComponentData) PushTriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {

}

// // SaveRetriggerRespinNum -
// func (basicComponentData *BasicComponentData) SaveRetriggerRespinNum()

// GetMask -
func (basicComponentData *BasicComponentData) GetMask() []bool {
	return nil
}

// ChgMask -
func (basicComponentData *BasicComponentData) ChgMask(curMask int, val bool) bool {
	return false
}

func (basicComponentData *BasicComponentData) PutInMoney(coins int) {

}

// ForceBranch -
func (basicComponentData *BasicComponentData) ForceBranch(branchIndex int) {
	if !gIsReleaseMode {
		basicComponentData.ForceBranchIndex = branchIndex
	}
}

// GetStrVal -
func (basicComponentData *BasicComponentData) GetStrVal(key string) (string, bool) {
	return "", false
}

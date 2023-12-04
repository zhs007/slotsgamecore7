package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const RespinTypeName = "respin"

type RespinData struct {
	BasicComponentData
	LastRespinNum         int
	CurRespinNum          int
	CurAddRespinNum       int
	RetriggerAddRespinNum int // 再次触发时增加的次数
	TotalCoinWin          int64
	TotalCashWin          int64
	LastTriggerNum        int      // 剩余的触发次数，respin有2种模式，一种是直接增加免费次数，一种是累积整体触发次数
	CurTriggerNum         int      // 当前已经触发次数
	Awards                []*Award // 当前已经触发次数
	TriggerRespinNum      []int    // 配合LastTriggerNum用的respin次数，-1表示用当前的RetriggerAddRespinNum，否则就是具体值
}

// OnNewGame -
func (respinData *RespinData) OnNewGame() {
	respinData.BasicComponentData.OnNewGame()

	respinData.LastRespinNum = 0
	respinData.CurRespinNum = 0
	respinData.CurAddRespinNum = 0
	respinData.TotalCoinWin = 0
	respinData.TotalCashWin = 0
	respinData.RetriggerAddRespinNum = 0
	respinData.LastTriggerNum = 0
	respinData.CurTriggerNum = 0
	respinData.Awards = nil
}

// OnNewStep -
func (respinData *RespinData) OnNewStep() {
	respinData.BasicComponentData.OnNewStep()

	respinData.CurAddRespinNum = 0
}

// BuildPBComponentData
func (respinData *RespinData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RespinData{
		BasicComponentData:    respinData.BuildPBBasicComponentData(),
		LastRespinNum:         int32(respinData.LastRespinNum),
		CurRespinNum:          int32(respinData.CurRespinNum),
		CurAddRespinNum:       int32(respinData.CurAddRespinNum),
		TotalCoinWin:          respinData.TotalCoinWin,
		TotalCashWin:          respinData.TotalCashWin,
		RetriggerAddRespinNum: int32(respinData.RetriggerAddRespinNum),
		LastTriggerNum:        int32(respinData.LastTriggerNum),
		CurTriggerNum:         int32(respinData.CurTriggerNum),
	}

	return pbcd
}

// RespinLevelConfig - configuration for Respin Level
type RespinLevelConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	LastRespinNum        int    `yaml:"lastRespinNum" json:"lastRespinNum"` // 倒数第几局开始
	MaxCoinWins          int    `yaml:"maxCoinWins" json:"maxCoinWins"`     // 如果最大获奖低于这个
	JumpComponent        string `yaml:"jumpComponent" json:"jumpComponent"` // 跳转到这个component
}

// RespinConfig - configuration for Respin
type RespinConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	InitRespinNum        int                  `yaml:"initRespinNum" json:"initRespinNum"`
	MainComponent        string               `yaml:"mainComponent" json:"mainComponent"`
	IsWinBreak           bool                 `yaml:"isWinBreak" json:"isWinBreak"`
	Levels               []*RespinLevelConfig `yaml:"levels" json:"levels"`
}

type Respin struct {
	*BasicComponent `json:"-"`
	Config          *RespinConfig `json:"config"`
}

// // OnNewGame -
// func (respin *Respin) OnNewGame(gameProp *GameProperty) error {
// 	cd := gameProp.MapComponentData[respin.Name]

// 	cd.OnNewGame()

// 	return nil
// }

// OnPlayGame - on playgame
func (respin *Respin) procLevel(level *RespinLevelConfig, respinData *RespinData, gameProp *GameProperty) bool {
	if respinData.LastRespinNum <= level.LastRespinNum && respinData.CoinWin < level.MaxCoinWins {
		return true
	}

	return false
}

// OnPlayGame - on playgame
func (respin *Respin) AddRespinTimes(gameProp *GameProperty, num int) {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	cd.LastRespinNum += num
	cd.CurAddRespinNum += num
}

// Init -
func (respin *Respin) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Respin.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &RespinConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Respin.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return respin.InitEx(cfg, pool)
}

// InitEx -
func (respin *Respin) InitEx(cfg any, pool *GamePropertyPool) error {
	respin.Config = cfg.(*RespinConfig)
	respin.Config.ComponentType = RespinTypeName

	respin.onInit(&respin.Config.BasicComponentConfig)

	return nil
}

// playgame
func (respin *Respin) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	respin.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	if cd.CurRespinNum == 0 && cd.LastRespinNum == 0 && respin.Config.InitRespinNum > 0 {
		cd.LastRespinNum = respin.Config.InitRespinNum
	}

recheck:
	if cd.LastRespinNum == 0 {
		if cd.LastTriggerNum > 0 {
			respin.Trigger(gameProp, plugin, curpr, gp)

			goto recheck
		}

		respin.onStepEnd(gameProp, curpr, gp, respin.Config.DefaultNextComponent)
	} else {
		nextComponent := respin.Config.MainComponent

		for _, v := range respin.Config.Levels {
			if respin.procLevel(v, cd, gameProp) {
				nextComponent = v.JumpComponent

				break
			}
		}

		if cd.LastRespinNum > 0 {
			cd.LastRespinNum--
		}

		cd.CurRespinNum++

		respin.onStepEnd(gameProp, curpr, gp, nextComponent)
	}

	// gp.AddComponentData(respin.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (respin *Respin) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	if cd.CurAddRespinNum > 0 {
		fmt.Printf("%v last %v, current %v, retrigger %v\n", respin.Name, cd.LastRespinNum, cd.CurRespinNum, cd.CurAddRespinNum)
	} else {
		fmt.Printf("%v last %v, current %v\n", respin.Name, cd.LastRespinNum, cd.CurRespinNum)
	}

	return nil
}

// OnStats
func (respin *Respin) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	if feature != nil && len(lst) > 0 {

		if feature.RespinNumStatus != nil ||
			feature.RespinWinStatus != nil {
			pbcd, lastpr := findLastPBComponentData(lst, respin.Name)
			if pbcd != nil {
				respin.onStatsWithPBEnding(feature, pbcd, lastpr)
			}
		}

		if feature.RespinStartNumStatus != nil {
			pbcd, firstpr := findFirstPBComponentData(lst, respin.Name)
			if pbcd != nil {
				respin.onStatsWithPBStart(feature, pbcd, firstpr)
			}
		}
	}

	return false, 0, 0
}

// onStatsWithPBEnding -
func (respin *Respin) onStatsWithPBEnding(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) error {
	pbcd, isok := pbComponentData.(*sgc7pb.RespinData)
	if !isok {
		goutils.Error("Respin.onStatsWithPBEnding",
			zap.Error(ErrIvalidProto))

		return ErrIvalidProto
	}

	if feature.RespinNumStatus != nil {
		feature.RespinNumStatus.AddStatus(int(pbcd.CurRespinNum))
	}

	if feature.RespinWinStatus != nil {
		feature.RespinWinStatus.AddStatus(int(pbcd.TotalCoinWin))
	}

	return nil
}

// onStatsWithPBEnding -
func (respin *Respin) onStatsWithPBStart(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) error {
	pbcd, isok := pbComponentData.(*sgc7pb.RespinData)
	if !isok {
		goutils.Error("Respin.onStatsWithPBStart",
			zap.Error(ErrIvalidProto))

		return ErrIvalidProto
	}

	if feature.RespinStartNumStatus != nil {
		feature.RespinStartNumStatus.AddStatus(int(pbcd.LastRespinNum))
	}

	return nil
}

// NewComponentData -
func (respin *Respin) NewComponentData() IComponentData {
	return &RespinData{}
}

// EachUsedResults -
func (respin *Respin) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.RespinData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("Respin.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

// OnPlayGame - on playgame
func (respin *Respin) OnPlayGameEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	cd.TotalCashWin += curpr.CashWin
	cd.TotalCoinWin += int64(curpr.CoinWin)

	if respin.Config.IsWinBreak && cd.TotalCoinWin > 0 {
		cd.LastRespinNum = 0
	}

	if cd.LastRespinNum == 0 && cd.LastTriggerNum == 0 {
		gameProp.onRespinEnding(respin.Name)
	}

	return nil
}

// IsRespin -
func (respin *Respin) IsRespin() bool {
	return true
}

// SaveRetriggerRespinNum -
func (respin *Respin) SaveRetriggerRespinNum(gameProp *GameProperty) {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	cd.RetriggerAddRespinNum = cd.LastRespinNum
}

// // Retrigger -
// func (respin *Respin) Retrigger(gameProp *GameProperty) {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	cd.LastRespinNum += cd.RetriggerAddRespinNum
// 	cd.CurAddRespinNum += cd.RetriggerAddRespinNum

// 	cd.CurTriggerNum++

// 	if cd.LastTriggerNum > 0 {
// 		cd.LastTriggerNum--
// 	}
// }

// Trigger -
func (respin *Respin) Trigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	n := cd.TriggerRespinNum[cd.CurTriggerNum]
	if n <= 0 {
		n = cd.RetriggerAddRespinNum

		cd.TriggerRespinNum[cd.CurTriggerNum] = n
	}

	cd.LastRespinNum += n
	cd.CurAddRespinNum += n

	cd.CurTriggerNum++

	if cd.LastTriggerNum > 0 {
		cd.LastTriggerNum--
	}

	for _, v := range cd.Awards {
		if v.TriggerIndex == cd.CurTriggerNum {
			gameProp.procAward(plugin, v, curpr, gp, true)
		}
	}
}

// AddRetriggerRespinNum -
func (respin *Respin) AddRetriggerRespinNum(gameProp *GameProperty, num int) {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	cd.RetriggerAddRespinNum += num
}

// AddTriggerAward -
func (respin *Respin) AddTriggerAward(gameProp *GameProperty, award *Award) {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	award.TriggerIndex = cd.CurTriggerNum + cd.LastTriggerNum

	cd.Awards = append(cd.Awards, award)
}

// PushTrigger -
func (respin *Respin) PushTrigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	cd.LastTriggerNum++

	cd.TriggerRespinNum = append(cd.TriggerRespinNum, num)

	if cd.LastRespinNum == 0 {
		respin.Trigger(gameProp, plugin, curpr, gp)
	}
}

func NewRespin(name string) IComponent {
	return &Respin{
		BasicComponent: NewBasicComponent(name),
	}
}

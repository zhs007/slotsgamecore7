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

type RespinData struct {
	BasicComponentData
	LastRespinNum   int
	TotalRespinNum  int
	CurRespinNum    int
	CurAddRespinNum int
	TotalCoinWin    int64
	TotalCashWin    int64
}

// OnNewGame -
func (respinData *RespinData) OnNewGame() {
	respinData.BasicComponentData.OnNewGame()

	respinData.LastRespinNum = 0
	respinData.TotalRespinNum = 0
	respinData.CurRespinNum = 0
	respinData.CurAddRespinNum = 0
	respinData.TotalCoinWin = 0
	respinData.TotalCashWin = 0
}

// OnNewGame -
func (respinData *RespinData) OnNewStep() {
	respinData.BasicComponentData.OnNewStep()

	respinData.CurAddRespinNum = 0
}

// BuildPBComponentData
func (respinData *RespinData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RespinData{
		BasicComponentData: respinData.BuildPBBasicComponentData(),
		LastRespinNum:      int32(respinData.LastRespinNum),
		TotalRespinNum:     int32(respinData.TotalRespinNum),
		CurRespinNum:       int32(respinData.CurRespinNum),
		CurAddRespinNum:    int32(respinData.CurAddRespinNum),
		TotalCoinWin:       respinData.TotalCoinWin,
		TotalCashWin:       respinData.TotalCashWin,
	}

	return pbcd
}

// RespinLevelConfig - configuration for Respin Level
type RespinLevelConfig struct {
	BasicComponentConfig `yaml:",inline"`
	LastRespinNum        int    `yaml:"lastRespinNum"` // 倒数第几局开始
	MaxCoinWins          int    `yaml:"maxCoinWins"`   // 如果最大获奖低于这个
	JumpComponent        string `yaml:"jumpComponent"` // 跳转到这个component
}

// RespinConfig - configuration for Respin
type RespinConfig struct {
	BasicComponentConfig `yaml:",inline"`
	DefaultRespinNum     int                  `yaml:"defaultRespinNum"`
	MainComponent        string               `yaml:"mainComponent"`
	Levels               []*RespinLevelConfig `yaml:"levels"`
}

type Respin struct {
	*BasicComponent
	Config *RespinConfig
}

// OnPlayGame - on playgame
func (respin *Respin) procLevel(level *RespinLevelConfig, respinData *RespinData, gameProp *GameProperty) bool {
	if respinData.LastRespinNum <= level.LastRespinNum && respinData.CoinWin < level.MaxCoinWins {
		return true
	}

	return false
}

// OnPlayGame - on playgame
func (respin *Respin) AddRespinTimes(gameProp *GameProperty, num int) error {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	cd.LastRespinNum += num
	cd.CurAddRespinNum += num

	return nil
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

	respin.Config = cfg

	respin.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (respin *Respin) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

	if cd.LastRespinNum == 0 {
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
		cd.TotalRespinNum++

		respin.onStepEnd(gameProp, curpr, gp, nextComponent)
	}

	gp.AddComponentData(respin.Name, cd)

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
	if feature != nil && feature.Status != nil && len(lst) > 0 {
		lastpr := lst[len(lst)-1]
		gp := lastpr.CurGameModParams.(*GameParams)
		if gp != nil {
			pbcd := gp.MapComponents[respin.Name]

			if pbcd != nil {
				respin.OnStatsWithPB(feature, pbcd, lastpr)
			}
		}
	}

	return false, 0, 0
}

// OnStatsWithPB -
func (respin *Respin) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.RespinData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("Respin.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	feature.Status.AddStatus(int(pbcd.TotalRespinNum))

	return respin.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
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

	if cd.LastRespinNum == 0 {
		gameProp.onRespinEnding(respin.Name)
	}

	return nil
}

func NewRespin(name string) IComponent {
	return &Respin{
		BasicComponent: NewBasicComponent(name),
	}
}

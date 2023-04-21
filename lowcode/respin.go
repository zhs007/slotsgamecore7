package lowcode

import (
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

// RespinConfig - configuration for Respin
type RespinConfig struct {
	BasicComponentConfig `yaml:",inline"`
	DefaultRespinNum     int    `yaml:"defaultRespinNum"`
	MainComponent        string `yaml:"mainComponent"`
}

type Respin struct {
	*BasicComponent
	Config *RespinConfig
}

// OnPlayGame - on playgame
func (respin *Respin) AddRespinTimes(gameProp *GameProperty, num int) error {
	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

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
		if cd.LastRespinNum > 0 {
			cd.LastRespinNum--
		}

		cd.CurRespinNum++

		respin.onStepEnd(gameProp, curpr, gp, respin.Config.MainComponent)
	}

	gp.AddComponentData(respin.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (respin *Respin) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

// OnStats
func (respin *Respin) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
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

	if cd.CurAddRespinNum > 0 {
		cd.LastRespinNum += cd.CurAddRespinNum
	}

	cd.TotalCashWin += curpr.CashWin
	cd.TotalCoinWin += int64(curpr.CoinWin)

	return nil
}

func NewRespin(name string) IComponent {
	return &Respin{
		BasicComponent: NewBasicComponent(name),
	}
}

package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
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
	LastTriggerNum        int      // 剩余的触发次数，respin有2种模式，一种是直接增加免费次数，一种是累积整体触发次数
	CurTriggerNum         int      // 当前已经触发次数
	Awards                []*Award // 当前已经触发次数
	TriggerRespinNum      []int    // 配合LastTriggerNum用的respin次数，-1表示用当前的RetriggerAddRespinNum，否则就是具体值
}

// OnNewGame -
func (respinData *RespinData) OnNewGame(gameProp *GameProperty, component IComponent) {
	respinData.BasicComponentData.OnNewGame(gameProp, component)

	respinData.LastRespinNum = 0
	respinData.CurRespinNum = 0
	respinData.CurAddRespinNum = 0
	respinData.TotalCoinWin = 0
	// respinData.TotalCashWin = 0
	respinData.RetriggerAddRespinNum = 0
	respinData.LastTriggerNum = 0
	respinData.CurTriggerNum = 0
	respinData.Awards = nil
}

// onNewStep -
func (respinData *RespinData) onNewStep() {
	// respinData.BasicComponentData.OnNewStep(gameProp, component)

	respinData.CurAddRespinNum = 0
}

// Clone
func (respinData *RespinData) Clone() IComponentData {
	target := &RespinData{
		BasicComponentData:    respinData.CloneBasicComponentData(),
		LastRespinNum:         respinData.LastRespinNum,
		CurRespinNum:          respinData.CurRespinNum,
		CurAddRespinNum:       respinData.CurAddRespinNum,
		TotalCoinWin:          respinData.TotalCoinWin,
		RetriggerAddRespinNum: respinData.RetriggerAddRespinNum,
		LastTriggerNum:        respinData.LastTriggerNum,
		CurTriggerNum:         respinData.CurTriggerNum,
	}

	target.TriggerRespinNum = make([]int, len(respinData.TriggerRespinNum))
	copy(target.TriggerRespinNum, respinData.TriggerRespinNum)

	return target
}

// BuildPBComponentData
func (respinData *RespinData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RespinData{
		BasicComponentData: respinData.BuildPBBasicComponentData(),
		LastRespinNum:      int32(respinData.LastRespinNum),
		CurRespinNum:       int32(respinData.CurRespinNum),
		CurAddRespinNum:    int32(respinData.CurAddRespinNum),
		TotalCoinWin:       respinData.TotalCoinWin,
		// TotalCashWin:          respinData.TotalCashWin,
		RetriggerAddRespinNum: int32(respinData.RetriggerAddRespinNum),
		LastTriggerNum:        int32(respinData.LastTriggerNum),
		CurTriggerNum:         int32(respinData.CurTriggerNum),
	}

	for _, v := range respinData.TriggerRespinNum {
		pbcd.TriggerRespinNum = append(pbcd.TriggerRespinNum, int32(v))
	}

	return pbcd
}

// GetValEx -
func (respinData *RespinData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVCurRespinNum {
		return respinData.CurRespinNum, true
	}

	return 0, false
}

// GetLastRespinNum -
func (respinData *RespinData) GetLastRespinNum() int {
	return respinData.LastRespinNum
}

// IsEnding -
func (respinData *RespinData) IsRespinEnding() bool {
	return respinData.LastRespinNum == 0 && respinData.LastTriggerNum == 0
}

// IsStarted -
func (respinData *RespinData) IsRespinStarted() bool {
	return respinData.CurRespinNum > 0
}

// ChgConfigIntVal -
func (respinData *RespinData) ChgConfigIntVal(key string, off int) {
	if key == CCVLastRespinNum {
		respinData.AddRespinTimes(off)
	} else if key == CCVRetriggerAddRespinNum {
		respinData.RetriggerAddRespinNum += off
	} else {
		respinData.BasicComponentData.ChgConfigIntVal(key, off)
	}
}

// SetConfigIntVal -
func (respinData *RespinData) SetConfigIntVal(key string, val int) {
	if key == CCVLastRespinNum {
		respinData.ResetRespinTimes(val)
	} else if key == CCVRetriggerAddRespinNum {
		respinData.RetriggerAddRespinNum = val
	} else {
		respinData.BasicComponentData.ChgConfigIntVal(key, val)
	}
}

// // AddRetriggerRespinNum -
// func (respinData *RespinData) AddRetriggerRespinNum(num int) {
// 	respinData.RetriggerAddRespinNum += num
// }

// AddTriggerRespinAward -
func (respinData *RespinData) AddTriggerRespinAward(award *Award) {
	award.TriggerIndex = respinData.CurTriggerNum + respinData.LastTriggerNum

	respinData.Awards = append(respinData.Awards, award)
}

// AddRespinTimes -
func (respinData *RespinData) AddRespinTimes(num int) {
	respinData.LastRespinNum += num
	respinData.CurAddRespinNum += num
}

// ResetRespinTimes -
func (respinData *RespinData) ResetRespinTimes(num int) {
	if respinData.LastRespinNum >= num {
		respinData.LastRespinNum = num
	} else {
		off := num - respinData.LastRespinNum
		respinData.AddRespinTimes(off)
	}
}

// OnTriggerRespin
func (respinData *RespinData) TriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {
	n := respinData.TriggerRespinNum[respinData.CurTriggerNum]
	if n <= 0 {
		n = respinData.RetriggerAddRespinNum

		respinData.TriggerRespinNum[respinData.CurTriggerNum] = n
	}

	respinData.LastRespinNum += n
	respinData.CurAddRespinNum += n

	respinData.CurTriggerNum++

	if respinData.LastTriggerNum > 0 {
		respinData.LastTriggerNum--
	}

	for _, v := range respinData.Awards {
		if v.TriggerIndex == respinData.CurTriggerNum {
			gameProp.procAward(plugin, v, curpr, gp, true)
		}
	}
}

// PushTrigger -
func (respinData *RespinData) PushTriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {
	respinData.LastTriggerNum++

	respinData.TriggerRespinNum = append(respinData.TriggerRespinNum, num)

	// 第一次trigger时，需要直接
	if respinData.LastRespinNum == 0 && respinData.CurRespinNum == 0 {
		respinData.TriggerRespin(gameProp, plugin, curpr, gp)
	}
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
	// InitRespinNum        int    `yaml:"initRespinNum" json:"initRespinNum"`
	MainComponent string `yaml:"mainComponent" json:"mainComponent"`
	// IsWinBreak           bool                 `yaml:"isWinBreak" json:"isWinBreak"`
	// Levels []*RespinLevelConfig `yaml:"levels" json:"levels"`
}

// SetLinkComponent
func (cfg *RespinConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "loop" {
		cfg.MainComponent = componentName
	}
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
func (respin *Respin) procLevel(level *RespinLevelConfig, respinData *RespinData, _ *GameProperty) bool {
	if respinData.LastRespinNum <= level.LastRespinNum && respinData.CoinWin < level.MaxCoinWins {
		return true
	}

	return false
}

// // OnPlayGame - on playgame
// func (respin *Respin) AddRespinTimes(gameProp *GameProperty, num int) {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	cd.LastRespinNum += num
// 	cd.CurAddRespinNum += num
// }

// Init -
func (respin *Respin) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Respin.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RespinConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Respin.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// respin.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*RespinData)

	cd.onNewStep()

	// if cd.CurRespinNum == 0 && cd.LastRespinNum == 0 && respin.Config.InitRespinNum > 0 {
	// 	cd.LastRespinNum = respin.Config.InitRespinNum
	// }

recheck:
	if cd.LastRespinNum == 0 {
		if cd.LastTriggerNum > 0 {
			cd.TriggerRespin(gameProp, plugin, curpr, gp)

			goto recheck
		}

		if respin.Config.DefaultNextComponent == "" {
			nc := respin.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		nc := respin.onStepEnd(gameProp, curpr, gp, respin.Config.DefaultNextComponent)

		return nc, nil
	}

	nextComponent := respin.Config.MainComponent

	// for _, v := range respin.Config.Levels {
	// 	if respin.procLevel(v, cd, gameProp) {
	// 		nextComponent = v.JumpComponent

	// 		break
	// 	}
	// }

	if cd.LastRespinNum > 0 {
		cd.LastRespinNum--
	}

	cd.CurRespinNum++

	nc := respin.onStepEnd(gameProp, curpr, gp, nextComponent)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (respin *Respin) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*RespinData)

	if cd.CurAddRespinNum > 0 {
		fmt.Printf("%v last %v, current %v, retrigger %v\n", respin.Name, cd.LastRespinNum, cd.CurRespinNum, cd.CurAddRespinNum)
	} else {
		fmt.Printf("%v last %v, current %v\n", respin.Name, cd.LastRespinNum, cd.CurRespinNum)
	}

	return nil
}

// // OnStats
// func (respin *Respin) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	if feature != nil && len(lst) > 0 {

// 		if feature.RespinNumStatus != nil ||
// 			feature.RespinWinStatus != nil {
// 			pbcd, lastpr := findLastPBComponentData(lst, respin.Name)
// 			if pbcd != nil {
// 				respin.onStatsWithPBEnding(feature, pbcd, lastpr)
// 			}
// 		}

// 		if feature.RespinStartNumStatus != nil {
// 			pbcd, firstpr := findFirstPBComponentData(lst, respin.Name)
// 			if pbcd != nil {
// 				respin.onStatsWithPBStart(feature, pbcd, firstpr)
// 			}
// 		}
// 	}

// 	return false, 0, 0
// }

// // onStatsWithPBEnding -
// func (respin *Respin) onStatsWithPBEnding(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) error {
// 	pbcd, isok := pbComponentData.(*sgc7pb.RespinData)
// 	if !isok {
// 		goutils.Error("Respin.onStatsWithPBEnding",
// 			goutils.Err(ErrIvalidProto))

// 		return ErrIvalidProto
// 	}

// 	if feature.RespinNumStatus != nil {
// 		feature.RespinNumStatus.AddStatus(int(pbcd.CurRespinNum))
// 	}

// 	if feature.RespinWinStatus != nil {
// 		feature.RespinWinStatus.AddStatus(int(pbcd.TotalCoinWin))
// 	}

// 	return nil
// }

// // onStatsWithPBEnding -
// func (respin *Respin) onStatsWithPBStart(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) error {
// 	pbcd, isok := pbComponentData.(*sgc7pb.RespinData)
// 	if !isok {
// 		goutils.Error("Respin.onStatsWithPBStart",
// 			goutils.Err(ErrIvalidProto))

// 		return ErrIvalidProto
// 	}

// 	if feature.RespinStartNumStatus != nil {
// 		feature.RespinStartNumStatus.AddStatus(int(pbcd.LastRespinNum))
// 	}

// 	return nil
// }

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
			goutils.Err(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

// ProcRespinOnStepEnd - 现在只有respin需要特殊处理结束，如果多层respin嵌套时，只要新的有next，就不会继续结束respin
func (respin *Respin) ProcRespinOnStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, canRemove bool) (string, error) {

	rcd := cd.(*RespinData)

	// rcd.TotalCashWin += curpr.CashWin
	rcd.TotalCoinWin += int64(curpr.CoinWin)

	if canRemove && rcd.LastRespinNum == 0 && rcd.LastTriggerNum == 0 {
		gameProp.removeRespin(respin.Name)

		if respin.Config.DefaultNextComponent != "" {
			return respin.Config.DefaultNextComponent, nil
		}
	}

	return "", nil
}

// IsRespin -
func (respin *Respin) IsRespin() bool {
	return true
}

// NewStats2 -
func (respin *Respin) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptRootTrigger})
}

// OnStats2
func (respin *Respin) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult) {
	isRunning := false
	isEnding := false

	if goutils.IndexOfStringSlice(gp.HistoryComponents, respin.Name, 0) >= 0 {
		isRunning = true
	}

	if goutils.IndexOfStringSlice(gp.RespinComponents, respin.Name, 0) < 0 {
		isEnding = true
	}

	s2.ProcStatsRespinTrigger(respin.Name, isRunning, int64(pr.CoinWin), isEnding)

	// rcd := icd.(*RespinData)

	// if rcd.LastRespinNum == 0 && rcd.LastTriggerNum == 0 {
	// 	s2.ProcStatsRootTrigger(respin.Name, rcd.TotalCoinWin, true)
	// } else {
	// 	s2.ProcStatsRootTrigger(respin.Name, 0, false)
	// }
}

// // IsTriggerRespin -
// func (respin *Respin) IsTriggerRespin() bool {
// 	rcd := cd.(*RespinData)

// 	return false
// }

// // SaveRetriggerRespinNum -
// func (respin *Respin) SaveRetriggerRespinNum(gameProp *GameProperty) {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	cd.RetriggerAddRespinNum = cd.LastRespinNum
// }

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

// // Trigger -
// func (respin *Respin) Trigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {
// 	cd := gameProp.GetGlobalComponentData(respin) //gameProp.MapComponentData[respin.Name].(*RespinData)

// 	n := cd.TriggerRespinNum[cd.CurTriggerNum]
// 	if n <= 0 {
// 		n = cd.RetriggerAddRespinNum

// 		cd.TriggerRespinNum[cd.CurTriggerNum] = n
// 	}

// 	cd.LastRespinNum += n
// 	cd.CurAddRespinNum += n

// 	cd.CurTriggerNum++

// 	if cd.LastTriggerNum > 0 {
// 		cd.LastTriggerNum--
// 	}

// 	for _, v := range cd.Awards {
// 		if v.TriggerIndex == cd.CurTriggerNum {
// 			gameProp.procAward(plugin, v, curpr, gp, true)
// 		}
// 	}
// }

// // AddRetriggerRespinNum -
// func (respin *Respin) AddRetriggerRespinNum(gameProp *GameProperty, num int) {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	cd.RetriggerAddRespinNum += num
// }

// // AddTriggerAward -
// func (respin *Respin) AddTriggerAward(gameProp *GameProperty, award *Award) {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	award.TriggerIndex = cd.CurTriggerNum + cd.LastTriggerNum

// 	cd.Awards = append(cd.Awards, award)
// }

// // PushTrigger -
// func (respin *Respin) PushTrigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	cd.LastTriggerNum++

// 	cd.TriggerRespinNum = append(cd.TriggerRespinNum, num)

// 	// 第一次trigger时，需要直接
// 	if cd.LastRespinNum == 0 && cd.CurRespinNum == 0 {
// 		respin.Trigger(gameProp, plugin, curpr, gp)
// 	}
// }

func (respin *Respin) getRetriggerRespinNum(basicCD *BasicComponentData) int {
	val, isok := basicCD.GetConfigIntVal(CCVReelSet)
	if isok {
		return val
	}

	return 0
}

// // GetLastRespinNum -
// func (respin *Respin) GetLastRespinNum(gameProp *GameProperty) int {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	return cd.LastRespinNum
// }

// // IsEnding -
// func (respin *Respin) IsEnding(gameProp *GameProperty) bool {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	return cd.LastRespinNum == 0 && cd.LastTriggerNum == 0
// }

// // IsStarted -
// func (respin *Respin) IsStarted(gameProp *GameProperty) bool {
// 	cd := gameProp.MapComponentData[respin.Name].(*RespinData)

// 	return cd.CurRespinNum > 0
// }

// GetAllLinkComponents - get all link components
func (respin *Respin) GetAllLinkComponents() []string {
	return []string{respin.Config.DefaultNextComponent, respin.Config.MainComponent}
}

// GetChildLinkComponents - get child link components
func (respin *Respin) GetChildLinkComponents() []string {
	return []string{respin.Config.MainComponent}
}

func NewRespin(name string) IComponent {
	return &Respin{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

//	"configuration": {
//		"isWinBreak": "false"
//	},
type jsonRespin struct {
	// IsWinBreak string `json:"isWinBreak"`
}

func (jr *jsonRespin) build() *RespinConfig {
	cfg := &RespinConfig{
		// IsWinBreak: jr.IsWinBreak == "true",
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseRespin(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	_, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRespin2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	// buf, err := cfg.MarshalJSON()
	// if err != nil {
	// 	goutils.Error("parseRespin2:MarshalJSON",
	// 		goutils.Err(err))

	// 	return "", err
	// }

	// data := &jsonRespin{}

	// err = sonic.Unmarshal(buf, data)
	// if err != nil {
	// 	goutils.Error("parseRespin2:Unmarshal",
	// 		goutils.Err(err))

	// 	return "", err
	// }

	cfgd := &RespinConfig{}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RespinTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

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

const MaskTypeName = "mask"

func boolArr2Int(vals []bool) int {
	iv := 0

	for _, v := range vals {
		if v {
			iv = iv*10 + 1
		} else {
			iv *= 10
		}
	}

	return iv
}

const (
	MaskTypeNone         int = 0
	MaskTypeSymbolInReel int = 1
)

func parserMaskType(str string) int {
	if str == "symbolInReel" {
		return MaskTypeSymbolInReel
	}

	return MaskTypeNone
}

type MaskData struct {
	BasicComponentData
	Num      int
	Vals     []bool
	NewChged int
	NewVals  []bool
}

// OnNewGame -
func (maskData *MaskData) OnNewGame() {
	maskData.BasicComponentData.OnNewGame()

	maskData.Vals = make([]bool, maskData.Num)
	maskData.NewVals = make([]bool, maskData.Num)
	maskData.NewChged = 0
}

// OnNewStep -
func (maskData *MaskData) OnNewStep() {
	maskData.BasicComponentData.OnNewStep()

	if maskData.NewChged > 0 {
		maskData.NewVals = make([]bool, maskData.Num)
		maskData.NewChged = 0
	}
}

// BuildPBComponentData
func (maskData *MaskData) BuildPBComponentData() proto.Message {
	pb := &sgc7pb.MaskData{
		Num:      int32(maskData.Num),
		NewChged: int32(maskData.NewChged),
		Vals:     make([]bool, len(maskData.Vals)),
		NewVals:  make([]bool, len(maskData.NewVals)),
	}

	copy(pb.Vals, maskData.Vals)
	copy(pb.NewVals, maskData.NewVals)

	return pb
}

// IsFull -
func (maskData *MaskData) IsFull() bool {
	for _, v := range maskData.Vals {
		if !v {
			return false
		}
	}

	return true
}

func newMaskData(num int) *MaskData {
	return &MaskData{
		Num:      num,
		Vals:     make([]bool, num),
		NewChged: 0,
		NewVals:  make([]bool, num),
	}
}

// MaskConfig - configuration for Mask
type MaskConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MaskType             string           `yaml:"maskType" json:"maskType"`
	Symbol               string           `yaml:"symbol" json:"symbol"`
	Num                  int              `yaml:"num" json:"num"`
	PerMaskAwards        []*Award         `yaml:"perMaskAwards" json:"perMaskAwards"`
	MapSPMaskAwards      map[int][]*Award `yaml:"mapSPMaskAwards" json:"mapSPMaskAwards"` // -1表示全满的奖励
	EndingSPAward        string           `yaml:"endingSPAward" json:"endingSPAward"`
}

type Mask struct {
	*BasicComponent `json:"-"`
	Config          *MaskConfig `json:"config"`
	MaskType        int         `json:"-"`
	SymbolCode      int         `json:"-"`
}

// Init -
func (mask *Mask) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Mask.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MaskConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Mask.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return mask.InitEx(cfg, pool)
}

// InitEx -
func (mask *Mask) InitEx(cfg any, pool *GamePropertyPool) error {
	mask.Config = cfg.(*MaskConfig)
	mask.Config.ComponentType = MaskTypeName

	mask.MaskType = parserMaskType(mask.Config.MaskType)
	mask.SymbolCode = pool.DefaultPaytables.MapSymbols[mask.Config.Symbol]

	if mask.Config.PerMaskAwards != nil {
		for _, v := range mask.Config.PerMaskAwards {
			v.Init()
		}
	}

	if mask.Config.MapSPMaskAwards != nil {
		for _, lst := range mask.Config.MapSPMaskAwards {
			for _, v := range lst {
				v.Init()
			}
		}
	}

	mask.onInit(&mask.Config.BasicComponentConfig)

	return nil
}

// OnNewGame - 因为 BasicComponent 考虑到效率，没有执行ComponentData的OnNewGame，所以这里需要特殊处理
func (mask *Mask) OnNewGame(gameProp *GameProperty) error {
	cd := gameProp.MapComponentData[mask.Name]

	cd.OnNewGame()

	return nil
}

// onMaskChg -
func (mask *Mask) ChgMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, md *MaskData, curpr *sgc7game.PlayResult, gp *GameParams, curMask int, val bool, noProcSPLevel bool) {
	if md.Vals[curMask] != val {
		md.Vals[curMask] = val
		md.NewVals[curMask] = val
		md.NewChged++

		mask.onMaskChg(plugin, gameProp, curpr, gp, curMask, noProcSPLevel)
	}
}

// onMaskChg -
func (mask *Mask) onMaskChg(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, curMask int, noProcSPLevel bool) {
	if mask.Config.PerMaskAwards != nil {
		for _, v := range mask.Config.PerMaskAwards {
			gameProp.procAward(plugin, v, curpr, gp)
		}
	}

	if noProcSPLevel {
		return
	}

	sp, isok := mask.Config.MapSPMaskAwards[curMask-1]
	if isok {
		for _, v := range sp {
			gameProp.procAward(plugin, v, curpr, gp)
		}
	}
}

// onMaskChg -
func (mask *Mask) ProcMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, targetScene string) {
	if mask.MaskType == MaskTypeSymbolInReel {
		cd := gameProp.MapComponentData[mask.Name].(*MaskData)

		gs := mask.GetTargetScene(gameProp, curpr, &cd.BasicComponentData, targetScene)

		for x, v := range cd.Vals {
			if !v {
				for _, s := range gs.Arr[x] {
					if s == mask.SymbolCode {
						mask.ChgMask(plugin, gameProp, cd, curpr, gp, x, true, mask.Config.EndingSPAward != "")

						break
					}
				}
			}
		}
	}
}

// playgame
func (mask *Mask) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	mask.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	mask.ProcMask(plugin, gameProp, curpr, gp, "")

	mask.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(mask.Name, gameProp.MapComponentData[mask.Name])

	return nil
}

// OnAsciiGame - outpur to asciigame
func (mask *Mask) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[mask.Name].(*MaskData)

	if cd.NewChged <= 0 {
		fmt.Printf("%v dose not collect new value, the mask value is %v\n", mask.Name, cd.Vals)
	} else {
		fmt.Printf("%v collect %v. the mask value is %v\n", mask.Name, cd.NewChged, cd.Vals)
	}

	return nil
}

// OnStats
func (mask *Mask) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	if feature != nil && len(lst) > 0 {
		if feature.RespinEndingStatus != nil {
			pbcd, lastpr := findLastPBComponentDataEx(lst, feature.RespinEndingName, mask.Name)

			if pbcd != nil {
				mask.OnStatsWithPB(feature, pbcd, lastpr)
			}
		}

		if feature.RespinStartStatus != nil {
			pbcd, lastpr := findFirstPBComponentDataEx(lst, feature.RespinStartName, mask.Name)

			if pbcd != nil {
				mask.OnStatsWithPB(feature, pbcd, lastpr)
			}
		}

		if feature.RespinStartStatusEx != nil {
			pbs, prs := findAllPBComponentDataEx(lst, feature.RespinStartNameEx, mask.Name)

			if len(pbs) > 0 {
				for i, v := range pbs {
					mask.OnStatsWithPB(feature, v, prs[i])
				}
			}
		}
	}

	return false, 0, 0
}

// OnStatsWithPB -
func (mask *Mask) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.MaskData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("Mask.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	if feature.RespinEndingStatus != nil {
		feature.RespinEndingStatus.AddStatus(boolArr2Int(pbcd.Vals))
	}

	if feature.RespinStartStatus != nil {
		feature.RespinStartStatus.AddStatus(boolArr2Int(pbcd.Vals))
	}

	if feature.RespinStartStatusEx != nil {
		feature.RespinStartStatusEx.AddStatus(boolArr2Int(pbcd.Vals))
	}

	return 0, nil
}

// NewComponentData -
func (mask *Mask) NewComponentData() IComponentData {
	return newMaskData(mask.Config.Num)
}

// EachUsedResults -
func (mask *Mask) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

// OnPlayGame - on playgame
func (mask *Mask) OnPlayGameEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	// 因为respin一定在最前面触发，所以可以在这里判断是否结束
	if mask.Config.EndingSPAward != "" {
		icd := gameProp.MapComponentData[mask.Config.EndingSPAward]
		if icd != nil {
			cd := icd.(*RespinData)
			if cd.LastRespinNum == 0 {
				md := gameProp.MapComponentData[mask.Name].(*MaskData)

				fullAward := mask.Config.MapSPMaskAwards[-1]
				if fullAward != nil {
					if md.IsFull() {
						for _, v := range fullAward {
							gameProp.procAward(plugin, v, curpr, gp)
						}
					}
				}
			}
		}
	}

	return nil
}

func NewMask(name string) IComponent {
	mask := &Mask{
		BasicComponent: NewBasicComponent(name),
	}

	return mask
}

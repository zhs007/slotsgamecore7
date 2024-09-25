package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const MaskTypeName = "mask"

type MaskData struct {
	BasicComponentData
	Num      int
	Vals     []bool
	NewChged int
	NewVals  []bool
}

// OnNewGame -
func (maskData *MaskData) OnNewGame(gameProp *GameProperty, component IComponent) {
	maskData.BasicComponentData.OnNewGame(gameProp, component)

	maskData.Vals = make([]bool, maskData.Num)
	maskData.NewVals = make([]bool, maskData.Num)
	maskData.NewChged = 0
}

// onNewStep -
func (maskData *MaskData) onNewStep() {
	if maskData.NewChged > 0 {
		maskData.NewVals = make([]bool, maskData.Num)
		maskData.NewChged = 0
	}
}

// Clone
func (maskData *MaskData) Clone() IComponentData {
	target := &MaskData{
		BasicComponentData: maskData.CloneBasicComponentData(),
		Num:                maskData.Num,
		NewChged:           maskData.NewChged,
	}

	target.Vals = make([]bool, len(maskData.Vals))
	copy(target.Vals, maskData.Vals)

	target.NewVals = make([]bool, len(maskData.NewVals))
	copy(target.NewVals, maskData.NewVals)

	return target
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

// GetMask -
func (maskData *MaskData) GetMask() []bool {
	return maskData.Vals
}

// ChgMask -
func (maskData *MaskData) ChgMask(curMask int, val bool) bool {
	if maskData.Vals[curMask] != val {
		maskData.Vals[curMask] = val
		maskData.NewVals[curMask] = val
		maskData.NewChged++

		return true
	}

	return false
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
	Num                  int              `yaml:"num" json:"num"`
	IgnoreFalse          bool             `yaml:"ignoreFalse" json:"ignoreFalse"`
	PerMaskAwards        []*Award         `yaml:"perMaskAwards" json:"perMaskAwards"`
	MapSPMaskAwards      map[int][]*Award `yaml:"mapSPMaskAwards" json:"mapSPMaskAwards"` // -1表示全满的奖励
}

// SetLinkComponent
func (cfg *MaskConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type Mask struct {
	*BasicComponent `json:"-"`
	Config          *MaskConfig `json:"config"`
}

// Init -
func (mask *Mask) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Mask.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &MaskConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Mask.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return mask.InitEx(cfg, pool)
}

// InitEx -
func (mask *Mask) InitEx(cfg any, pool *GamePropertyPool) error {
	mask.Config = cfg.(*MaskConfig)
	mask.Config.ComponentType = MaskTypeName

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

// onMaskChg -
func (mask *Mask) onMaskChg(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, curMask int) {
	if mask.Config.PerMaskAwards != nil {
		for _, v := range mask.Config.PerMaskAwards {
			gameProp.procAward(plugin, v, curpr, gp, false)
		}
	}

	sp, isok := mask.Config.MapSPMaskAwards[curMask]
	if isok {
		for _, v := range sp {
			gameProp.procAward(plugin, v, curpr, gp, false)
		}
	}
}

// onMaskChg -
func (mask *Mask) ProcMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult, gp *GameParams, targetScene string) {
	// if mask.MaskType == MaskTypeSymbolInReel {
	// 	cd := gameProp.MapComponentData[mask.Name].(*MaskData)

	// 	gs := mask.GetTargetScene3(gameProp, curpr, prs, &cd.BasicComponentData, mask.Name, targetScene, 0)

	// 	for x, v := range cd.Vals {
	// 		if !v {
	// 			for _, s := range gs.Arr[x] {
	// 				if s == mask.SymbolCode {
	// 					mask.ChgMask(plugin, gameProp, cd, curpr, gp, x, true, mask.Config.EndingSPAward != "")

	// 					break
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

// playgame
func (mask *Mask) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	mcd := cd.(*MaskData)
	mcd.onNewStep()

	nc := mask.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (mask *Mask) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	mcd := cd.(*MaskData)

	if mcd.NewChged <= 0 {
		fmt.Printf("%v dose not collect new value, the mask value is %v\n", mask.Name, mcd.Vals)
	} else {
		fmt.Printf("%v collect %v. the mask value is %v\n", mask.Name, mcd.NewChged, mcd.Vals)
	}

	return nil
}

// NewComponentData -
func (mask *Mask) NewComponentData() IComponentData {
	return newMaskData(mask.Config.Num)
}

// EachUsedResults -
func (mask *Mask) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

// IsMask -
func (mask *Mask) IsMask() bool {
	return true
}

// SetMaskVal -
func (mask *Mask) SetMaskVal(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, index int, val bool) error {
	mcd := cd.(*MaskData)

	if mcd.ChgMask(index, val) {
		mask.onMaskChg(plugin, gameProp, curpr, gp, index)
	}

	return nil
}

// SetMask -
func (mask *Mask) SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, arrMask []bool) error {
	mcd := cd.(*MaskData)

	if mask.Config.IgnoreFalse {
		for x, v := range arrMask {
			if v && mcd.ChgMask(x, v) {
				mask.onMaskChg(plugin, gameProp, curpr, gp, x)
			}
		}
	} else {
		for x, v := range arrMask {
			if mcd.ChgMask(x, v) {
				mask.onMaskChg(plugin, gameProp, curpr, gp, x)
			}
		}
	}

	if mcd.NewChged > 0 && mcd.IsFull() {
		sp, isok := mask.Config.MapSPMaskAwards[-1]
		if isok {
			for _, v := range sp {
				gameProp.procAward(plugin, v, curpr, gp, false)
			}
		}
	}

	return nil
}

// SetMaskOnlyTrue -
func (mask *Mask) SetMaskOnlyTrue(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, arrMask []bool) error {
	mcd := cd.(*MaskData)

	for x, v := range arrMask {
		if v && mcd.ChgMask(x, v) {
			mask.onMaskChg(plugin, gameProp, curpr, gp, x)
		}
	}

	return nil
}

func NewMask(name string) IComponent {
	mask := &Mask{
		BasicComponent: NewBasicComponent(name, 1),
	}

	return mask
}

//	"configuration": {
//		"length": 5
//	},
type jsonMask struct {
	Length      int  `json:"length"`
	IgnoreFalse bool `json:"ignoreFalse"`
}

func (jcfg *jsonMask) build() *MaskConfig {
	cfg := &MaskConfig{
		Num:         jcfg.Length,
		IgnoreFalse: jcfg.IgnoreFalse,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseMask(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseMask:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseMask:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonMask{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseMask:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, mapAwards, err := parseMaskControllers(ctrls)
		if err != nil {
			goutils.Error("parseMask:parseMaskControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.PerMaskAwards = awards
		cfgd.MapSPMaskAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: MaskTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

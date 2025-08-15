package lowcode

import (
	"log/slog"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const HoldAndRespinReelsTypeName = "holdAndRespinReels"

type HoldAndRespinReelsType int

const (
	HARRTypeKeepReels   HoldAndRespinReelsType = 0 // keep reels
	HARRTypeResetReels  HoldAndRespinReelsType = 1 // reset reels
	HARRTypeRerollReels HoldAndRespinReelsType = 2 // reroll reels
)

func parseHoldAndRespinReelType(str string) HoldAndRespinReelsType {
	switch str {
	case "resetreels":
		return HARRTypeResetReels
	case "rerollreels":
		return HARRTypeRerollReels
	}

	return HARRTypeKeepReels
}

type HoldAndRespinReelsData struct {
	BasicComponentData
	ReelSetIndex int // The index of the currently selected reelset
}

// OnNewGame -
func (harData *HoldAndRespinReelsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	harData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (harData *HoldAndRespinReelsData) onNewStep() {
	harData.UsedScenes = nil
	harData.ReelSetIndex = -1
}

// GetValEx -
func (harData *HoldAndRespinReelsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSelectedIndex {
		return harData.ReelSetIndex, true
	}

	return 0, false
}

// Clone
func (harData *HoldAndRespinReelsData) Clone() IComponentData {
	target := &HoldAndRespinReelsData{
		BasicComponentData: harData.CloneBasicComponentData(),
		ReelSetIndex:       harData.ReelSetIndex,
	}

	return target
}

// BuildPBComponentData
func (harData *HoldAndRespinReelsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WeightReelsData{
		BasicComponentData: harData.BuildPBBasicComponentData(),
		ReelSetIndex:       int32(harData.ReelSetIndex),
	}

	return pbcd
}

// BasicReelsConfig - configuration for HoldAndRespinReels
type HoldAndRespinReelsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                 `yaml:"type" json:"type"`
	Type                 HoldAndRespinReelsType `yaml:"-" json:"-"`
	ReelSetsWeight       string                 `yaml:"reelSetWeight" json:"reelSetWeight"`
	ReelSetsWeightVW     *sgc7game.ValWeights2  `json:"-"`
	HoldReels            []bool                 `yaml:"holdReels" json:"holdReels"`
	ReelSet              string                 `yaml:"reelSet" json:"reelSet"`
	MapAwards            map[string][]*Award    `yaml:"controllers" json:"controllers"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *HoldAndRespinReelsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type HoldAndRespinReels struct {
	*BasicComponent `json:"-"`
	Config          *HoldAndRespinReelsConfig `json:"config"`
}

// Init -
func (har *HoldAndRespinReels) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("HoldAndRespinReels.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &HoldAndRespinReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("HoldAndRespinReels.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return har.InitEx(cfg, pool)
}

// InitEx -
func (har *HoldAndRespinReels) InitEx(cfg any, pool *GamePropertyPool) error {
	har.Config = cfg.(*HoldAndRespinReelsConfig)
	har.Config.ComponentType = HoldAndRespinReelsTypeName

	har.Config.Type = parseHoldAndRespinReelType(har.Config.StrType)

	if har.Config.ReelSetsWeight != "" {
		vw2, err := pool.LoadStrWeights(har.Config.ReelSetsWeight, har.Config.UseFileMapping)
		if err != nil {
			goutils.Error("HoldAndRespinReels.Init:LoadValWeights",
				slog.String("ReelSetsWeight", har.Config.ReelSetsWeight),
				goutils.Err(err))

			return err
		}

		har.Config.ReelSetsWeightVW = vw2
	}

	for _, awards := range har.Config.MapAwards {
		for _, award := range awards {
			award.Init()
		}
	}

	har.onInit(&har.Config.BasicComponentConfig)

	return nil
}

func (har *HoldAndRespinReels) getReelSetWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVReelSetWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadStrWeights(str, har.Config.UseFileMapping)

		return vw2
	}

	return har.Config.ReelSetsWeightVW
}

func (har *HoldAndRespinReels) getReelSet(basicCD *BasicComponentData) string {
	str := basicCD.GetConfigVal(CCVReelSet)
	if str != "" {
		return str
	}

	return har.Config.ReelSet
}

// OnProcControllers -
func (har *HoldAndRespinReels) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if strVal != "" {
		awards, isok := har.Config.MapAwards[strVal]
		if isok {
			gameProp.procAwards(plugin, awards, curpr, gp)
		}
	}
}

// playgame
func (har *HoldAndRespinReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	hrd := icd.(*HoldAndRespinReelsData)

	hrd.onNewStep()

	reelname := ""
	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)

	switch har.Config.Type {
	case HARRTypeKeepReels:
		reelname = gs.ReelName

		rd, isok := gameProp.Pool.Config.MapReels[reelname]
		if !isok {
			goutils.Error("HoldAndRespinReels.OnPlayGame:MapReels",
				goutils.Err(ErrInvalidReels))

			return "", ErrInvalidReels
		}

		gameProp.CurReels = rd

		sc := gs.CloneEx(gameProp.PoolScene)

		sc.RandMaskReelsWithReelData(gameProp.CurReels, plugin, har.Config.HoldReels, true)

		har.AddScene(gameProp, curpr, sc, &hrd.BasicComponentData)
	case HARRTypeResetReels:
		reelname = har.getReelSet(&hrd.BasicComponentData)

		rd, isok := gameProp.Pool.Config.MapReels[reelname]
		if !isok {
			goutils.Error("HoldAndRespinReels.OnPlayGame:MapReels",
				goutils.Err(ErrInvalidReels))

			return "", ErrInvalidReels
		}

		gameProp.CurReels = rd

		sc := gs.CloneEx(gameProp.PoolScene)

		sc.RandMaskReelsWithReelData(gameProp.CurReels, plugin, har.Config.HoldReels, true)
		sc.ReelName = reelname

		har.AddScene(gameProp, curpr, sc, &hrd.BasicComponentData)
	case HARRTypeRerollReels:
		vw2 := har.getReelSetWeight(gameProp, &hrd.BasicComponentData)
		if vw2 != nil {
			val, si, err := vw2.RandValEx(plugin)
			if err != nil {
				goutils.Error("HoldAndRespinReels.OnPlayGame:ReelSetWeights.RandVal",
					goutils.Err(err))

				return "", err
			}

			hrd.ReelSetIndex = si

			curreels := val.String()

			rd, isok := gameProp.Pool.Config.MapReels[curreels]
			if !isok {
				goutils.Error("HoldAndRespinReels.OnPlayGame:MapReels",
					goutils.Err(ErrInvalidReels))

				return "", ErrInvalidReels
			}

			gameProp.CurReels = rd
			reelname = curreels
		}

		sc := gs.CloneEx(gameProp.PoolScene)

		sc.RandMaskReelsWithReelData(gameProp.CurReels, plugin, har.Config.HoldReels, true)
		sc.ReelName = reelname

		har.AddScene(gameProp, curpr, sc, &hrd.BasicComponentData)
	}

	har.ProcControllers(gameProp, plugin, curpr, gp, -1, "<trigger>")
	har.ProcControllers(gameProp, plugin, curpr, gp, -1, reelname)

	nc := har.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (har *HoldAndRespinReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	hrd := icd.(*HoldAndRespinReelsData)

	if len(hrd.UsedScenes) > 0 {
		asciigame.OutputScene("hold and respin symbols", pr.Scenes[hrd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// NewComponentData -
func (har *HoldAndRespinReels) NewComponentData() IComponentData {
	return &HoldAndRespinReelsData{}
}

func NewHoldAndRespinReels(name string) IComponent {
	har := &HoldAndRespinReels{
		BasicComponent: NewBasicComponent(name, 0),
	}

	return har
}

// "type": "resetReels",
// "holdReels": [
//
//	1,
//	0,
//	0,
//	0,
//	0,
//	1
//
// ],
// "reelSet": "bg-reel01"
type jsonHoldAndRespinReels struct {
	Type          string `json:"type"`
	HoldReels     []int  `json:"holdReels"`
	ReelSet       string `json:"reelSet"` // The reel set to use for the hold and respin
	ReelSetWeight string `json:"reelSetWeight"`
}

func (jcfg *jsonHoldAndRespinReels) build() *HoldAndRespinReelsConfig {
	cfg := &HoldAndRespinReelsConfig{
		StrType:        strings.ToLower(jcfg.Type),
		HoldReels:      make([]bool, len(jcfg.HoldReels)),
		ReelSet:        jcfg.ReelSet,
		ReelSetsWeight: jcfg.ReelSetWeight,
	}

	for i, v := range jcfg.HoldReels {
		if v == 1 {
			cfg.HoldReels[i] = true
		} else {
			cfg.HoldReels[i] = false
		}
	}

	return cfg
}

func parseHoldAndRespinReels(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseHoldAndRespinReels:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseHoldAndRespinReels:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonHoldAndRespinReels{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseHoldAndRespinReels:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseHoldAndRespinReels:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: HoldAndRespinReelsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

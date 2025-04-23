package lowcode

import (
	"context"
	"log/slog"
	"os"
	"slices"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const FeaturePickTypeName = "featurePick"

type FeaturePickData struct {
	BasicComponentData
	Selected    []string
	UnSelected  []string
	CurSelected []string
}

// OnNewGame -
func (featurePickData *FeaturePickData) OnNewGame(gameProp *GameProperty, component IComponent) {
	featurePickData.BasicComponentData.OnNewGame(gameProp, component)

	featurePickData.Selected = nil
	featurePickData.UnSelected = nil
}

// onNewStep -
func (featurePickData *FeaturePickData) onNewStep() {
	if featurePickData.CurSelected != nil {
		featurePickData.Selected = append(featurePickData.Selected, featurePickData.CurSelected...)

		featurePickData.CurSelected = nil
	}
}

// Clone
func (featurePickData *FeaturePickData) Clone() IComponentData {
	target := &FeaturePickData{
		BasicComponentData: featurePickData.CloneBasicComponentData(),
		Selected:           slices.Clone(featurePickData.Selected),
		UnSelected:         slices.Clone(featurePickData.UnSelected),
		CurSelected:        slices.Clone(featurePickData.CurSelected),
	}

	return target
}

// BuildPBComponentData
func (featurePickData *FeaturePickData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.FeaturePickData{
		BasicComponentData: featurePickData.BuildPBBasicComponentData(),
		Selected:           slices.Clone(featurePickData.Selected),
		UnSelected:         slices.Clone(featurePickData.UnSelected),
		CurSelected:        slices.Clone(featurePickData.CurSelected),
	}

	return pbcd
}

// FeaturePickConfig - configuration for FeaturePick
type FeaturePickConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrWeight            string                `yaml:"weight" json:"weight"` // weight
	Weight               *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	PoolSize             int                   `json:"poolSize"`
	PickNum              int                   `json:"pickNum"`
	MapControllers       map[string][]*Award   `yaml:"mapControllers" json:"mapControllers"`
}

// SetLinkComponent
func (cfg *FeaturePickConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type FeaturePick struct {
	*BasicComponent `json:"-"`
	Config          *FeaturePickConfig `json:"config"`
}

// Init -
func (featurePick *FeaturePick) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("FeaturePick.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &FeaturePickConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("FeaturePick.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return featurePick.InitEx(cfg, pool)
}

// InitEx -
func (featurePick *FeaturePick) InitEx(cfg any, pool *GamePropertyPool) error {
	featurePick.Config = cfg.(*FeaturePickConfig)
	featurePick.Config.ComponentType = FeaturePickTypeName

	vw2, err := pool.LoadStrWeights(featurePick.Config.StrWeight, true)
	if err != nil {
		goutils.Error("FeaturePick.InitEx:LoadStrWeights",
			slog.String("weight", featurePick.Config.StrWeight),
			goutils.Err(err))

		return err
	}

	featurePick.Config.Weight = vw2

	for _, awards := range featurePick.Config.MapControllers {
		for _, award := range awards {
			award.Init()
		}
	}

	featurePick.onInit(&featurePick.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (featurePick *FeaturePick) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	controllers, isok := featurePick.Config.MapControllers[strVal]
	if isok {
		if len(controllers) > 0 {
			gameProp.procAwards(plugin, controllers, curpr, gp)
		}
	}
}

// playgame
func (featurePick *FeaturePick) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*FeaturePickData)
	cd.onNewStep()

	if cd.UnSelected == nil && cd.Selected == nil && featurePick.Config.PoolSize > 0 {
		vw := featurePick.getWeight(gameProp, &cd.BasicComponentData).Clone()
		cd.UnSelected = make([]string, featurePick.Config.PoolSize)

		for i := range featurePick.Config.PoolSize {
			cv, err := vw.RandVal(plugin)
			if err != nil {
				goutils.Error("FeaturePick.OnPlayGame:RandVal",
					goutils.Err(err))

				return "", err
			}

			cd.UnSelected[i] = cv.String()
			vw.RemoveVal(cv)
		}
	}

	pickNum := featurePick.getPickNum(gameProp, &cd.BasicComponentData)

	if pickNum > 0 {
		curPickNum := pickNum
		if pickNum > len(cd.UnSelected) {
			curPickNum = len(cd.UnSelected)
		}

		cd.CurSelected = make([]string, curPickNum)

		for i := range curPickNum {
			ci, err := plugin.Random(context.Background(), curPickNum-i)
			if err != nil {
				goutils.Error("FeaturePick.OnPlayGame:Random",
					goutils.Err(err))

				return "", err
			}

			cd.CurSelected[i] = cd.UnSelected[ci]
			featurePick.ProcControllers(gameProp, plugin, curpr, gp, -1, cd.UnSelected[ci])

			cd.UnSelected = slices.Delete(cd.UnSelected, ci, ci+1)
		}

		if pickNum > curPickNum {
			for i := curPickNum; i < pickNum; i++ {
				featurePick.ProcControllers(gameProp, plugin, curpr, gp, -1, "<extra>")
			}
		}

		nc := featurePick.onStepEnd(gameProp, curpr, gp, "")
		return nc, nil
	}

	nc := featurePick.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (featurePick *FeaturePick) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewComponentData -
func (featurePick *FeaturePick) NewComponentData() IComponentData {
	return &FeaturePickData{}
}

func (featurePick *FeaturePick) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return featurePick.Config.Weight
}

func (featurePick *FeaturePick) getPickNum(gameProp *GameProperty, basicCD *BasicComponentData) int {
	ival, isok := basicCD.GetConfigIntVal(CCVPickNum)
	if isok {
		return ival
	}

	return featurePick.Config.PickNum
}

// OnStats2
func (featurePick *FeaturePick) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	featurePick.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*FeaturePickData)

	for _, v := range cd.CurSelected {
		s2.ProcStatsStrVal(featurePick.GetName(), v)
	}
}

// NewStats2 -
func (featurePick *FeaturePick) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, []stats2.Option{stats2.OptStrVal})
}

func NewFeaturePick(name string) IComponent {
	return &FeaturePick{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "poolSize": 15,
// "weight": "fgpickweight",
// "pickNum": 6
type jsonFeaturePick struct {
	PoolSize int    `json:"poolSize"`
	PickNum  int    `json:"pickNum"`
	Weight   string `json:"weight"`
}

func (jcfg *jsonFeaturePick) build() *FeaturePickConfig {
	cfg := &FeaturePickConfig{
		StrWeight: jcfg.Weight,
		PoolSize:  jcfg.PoolSize,
		PickNum:   jcfg.PickNum,
	}

	return cfg
}

func parseFeaturePick(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseFeaturePick:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseFeaturePick:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonFeaturePick{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseFeaturePick:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseMapStringAndAllControllers(ctrls)
		if err != nil {
			goutils.Error("parseFeaturePick:parseMapStringAndAllControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapControllers = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: FeaturePickTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

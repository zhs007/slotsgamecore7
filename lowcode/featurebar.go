package lowcode

import (
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
	"gopkg.in/yaml.v2"
)

const FeatureBarTypeName = "featureBar"

type FeatureBarData struct {
	BasicComponentData
	Features []int
	cfg      *FeatureBarConfig
}

// OnNewGame -
func (featureBarData *FeatureBarData) OnNewGame(gameProp *GameProperty, component IComponent) {
	featureBarData.BasicComponentData.OnNewGame(gameProp, component)

	featureBarData.Features = make([]int, 0, featureBarData.cfg.Length)
}

// BuildPBComponentData
func (featureBarData *FeatureBarData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.FeatureBarData{
		BasicComponentData: featureBarData.BuildPBBasicComponentData(),
		Features:           make([]int32, len(featureBarData.Features)),
	}

	for i, f := range featureBarData.Features {
		pbcd.Features[i] = int32(f)
	}

	return pbcd
}

// Clone
func (featureBarData *FeatureBarData) Clone() IComponentData {
	target := &FeatureBarData{
		BasicComponentData: featureBarData.CloneBasicComponentData(),
		cfg:                featureBarData.cfg,
	}

	return target
}

// FeatureBarConfig - configuration for FeatureBar
type FeatureBarConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Length               int                   `yaml:"length" json:"length"`               // bar 的长度
	StrFeatureWeight     string                `yaml:"featureWeight" json:"featureWeight"` // feature权重
	FeatureWeight        *sgc7game.ValWeights2 `yaml:"-" json:"-"`                         // feature权重
	MapAwards            map[int][]*Award      `yaml:"awards" json:"awards"`               // 新的奖励系统
}

// SetLinkComponent
func (cfg *FeatureBarConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type FeatureBar struct {
	*BasicComponent `json:"-"`
	Config          *FeatureBarConfig `json:"config"`
}

// Init -
func (featureBar *FeatureBar) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("FeatureBar.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &FeatureBarConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("FeatureBar.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return featureBar.InitEx(cfg, pool)
}

// InitEx -
func (featureBar *FeatureBar) InitEx(cfg any, pool *GamePropertyPool) error {
	featureBar.Config = cfg.(*FeatureBarConfig)
	featureBar.Config.ComponentType = FeatureBarTypeName

	if featureBar.Config.StrFeatureWeight != "" {
		vw2, err := pool.LoadIntWeights(featureBar.Config.StrFeatureWeight, true)
		if err != nil {
			goutils.Error("FeatureBar.InitEx:LoadIntWeights",
				slog.String("FeatureWeight", featureBar.Config.StrFeatureWeight),
				goutils.Err(err))

			return err
		}

		featureBar.Config.FeatureWeight = vw2
	}

	for _, awards := range featureBar.Config.MapAwards {
		for _, award := range awards {
			award.Init()
		}
	}

	featureBar.onInit(&featureBar.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (featureBar *FeatureBar) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(featureBar.Config.MapAwards) > 0 {
		awards, isok := featureBar.Config.MapAwards[val]
		if isok {
			gameProp.procAwards(plugin, awards, curpr, gp)
		}
	}
}

// getWeight -
func (featureBar *FeatureBar) getWeight(gameProp *GameProperty, cd *FeatureBarData) *sgc7game.ValWeights2 {
	str := cd.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return featureBar.Config.FeatureWeight
}

// playgame
func (featureBar *FeatureBar) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*FeatureBarData)

	vw := featureBar.getWeight(gameProp, cd)

	if len(cd.Features) == 0 {
		for i := 0; i < featureBar.Config.Length; i++ {
			feature, err := vw.RandVal(plugin)
			if err != nil {
				goutils.Error("FeatureBar.OnPlayGame:RandVal",
					goutils.Err(err))

				return "", err
			}

			cd.Features = append(cd.Features, feature.Int())
		}
	} else {
		curFeature := cd.Features[0]

		cd.Features = cd.Features[1:]

		feature, err := vw.RandVal(plugin)
		if err != nil {
			goutils.Error("FeatureBar.OnPlayGame:RandVal",
				goutils.Err(err))

			return "", err
		}

		cd.Features = append(cd.Features, feature.Int())

		featureBar.ProcControllers(gameProp, plugin, curpr, gp, curFeature, "")
	}

	nc := featureBar.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (featureBar *FeatureBar) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewComponentData -
func (featureBar *FeatureBar) NewComponentData() IComponentData {
	return &FeatureBarData{
		cfg: featureBar.Config,
	}
}

func NewFeatureBar(name string) IComponent {
	return &FeatureBar{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "configuration": {
// "Length": 5,
// "FeatureWeights": "bgfeature"
// }
type jsonFeatureBar struct {
	Length           int    `json:"Length"`         // bar 的长度
	StrFeatureWeight string `json:"FeatureWeights"` // feature权重
}

func (jcfg *jsonFeatureBar) build() *FeatureBarConfig {
	cfg := &FeatureBarConfig{
		Length:           jcfg.Length,
		StrFeatureWeight: jcfg.StrFeatureWeight,
	}

	return cfg
}

func parseFeatureBar(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseFeatureBar:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseControllerWorker:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonFeatureBar{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseFeatureBar:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseFeatureBarControllers(ctrls)
		if err != nil {
			goutils.Error("parseFeatureBar:parseFeatureBarControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: FeatureBarTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

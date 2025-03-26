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
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const FeatureBar2TypeName = "featureBar2"

type FeatureBar2PS struct {
	Features []string `json:"features"` // features
}

// SetPublicJson
func (ps *FeatureBar2PS) SetPublicJson(str string) error {
	err := sonic.UnmarshalString(str, ps)
	if err != nil {
		goutils.Error("FeatureBar2PS.SetPublicJson:UnmarshalString",
			goutils.Err(err))

		return err
	}

	return nil
}

// SetPrivateJson
func (ps *FeatureBar2PS) SetPrivateJson(str string) error {
	return nil
}

// GetPublicJson
func (ps *FeatureBar2PS) GetPublicJson() string {
	str, err := sonic.MarshalString(ps)
	if err != nil {
		goutils.Error("FeatureBar2PS.GetPublicJson:MarshalString",
			goutils.Err(err))

		return ""
	}

	return str
}

// GetPrivateJson
func (ps *FeatureBar2PS) GetPrivateJson() string {
	return ""
}

// Clone
func (ps *FeatureBar2PS) Clone() IComponentPS {
	return &FeatureBar2PS{
		Features: slices.Clone(ps.Features),
	}
}

type FeatureBar2Data struct {
	BasicComponentData
	Features      []string
	UsedFeatures  []string
	CacheFeatures []string
	CurFeature    string
	cfg           *FeatureBar2Config
}

// OnNewGame -
func (featureBar2Data *FeatureBar2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	featureBar2Data.BasicComponentData.OnNewGame(gameProp, component)

	featureBar2Data.Features = make([]string, 0, featureBar2Data.cfg.Length)
	featureBar2Data.UsedFeatures = nil
	featureBar2Data.CacheFeatures = nil
	featureBar2Data.CurFeature = featureBar2Data.cfg.EmptyFeature
}

// GetStrVal -
func (featureBar2Data *FeatureBar2Data) GetStrVal(key string) (string, bool) {
	if key == CSVValue {
		return featureBar2Data.CurFeature, true
	}

	return "", false
}

// BuildPBComponentData
func (featureBar2Data *FeatureBar2Data) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.FeatureBar2Data{
		BasicComponentData: featureBar2Data.BuildPBBasicComponentData(),
		Features:           slices.Clone(featureBar2Data.Features),
		UsedFeatures:       slices.Clone(featureBar2Data.UsedFeatures),
		CacheFeatures:      slices.Clone(featureBar2Data.CacheFeatures),
		CurFeature:         featureBar2Data.CurFeature,
	}

	return pbcd
}

// Clone
func (featureBar2Data *FeatureBar2Data) Clone() IComponentData {
	target := &FeatureBar2Data{
		BasicComponentData: featureBar2Data.CloneBasicComponentData(),
		cfg:                featureBar2Data.cfg,
		Features:           slices.Clone(featureBar2Data.Features),
		UsedFeatures:       slices.Clone(featureBar2Data.UsedFeatures),
		CacheFeatures:      slices.Clone(featureBar2Data.CacheFeatures),
		CurFeature:         featureBar2Data.CurFeature,
	}

	return target
}

// FeatureBar2Config - configuration for FeatureBar2
type FeatureBar2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Length               int                   `yaml:"length" json:"length"`                   // bar 的长度
	StrFeatureWeight     string                `yaml:"featureWeight" json:"featureWeight"`     // feature权重
	FeatureWeight        *sgc7game.ValWeights2 `yaml:"-" json:"-"`                             // feature权重
	FirstJumpWeight      int                   `yaml:"firstJumpWeight" json:"firstJumpWeight"` // firstJump 权重
	EmptyFeature         string                `yaml:"emptyFeature" json:"emptyFeature"`       // emptyFeature
	IsPlayerState        bool                  `yaml:"IsPlayerState" json:"IsPlayerState"`     // IsPlayerState
	IsMergeData          bool                  `yaml:"IsMergeData" json:"IsMergeData"`         // IsMergeData
	MapAwards            map[string][]*Award   `yaml:"mapAwards" json:"mapAwards"`             // 新的奖励系统
	Awards               []*Award              `yaml:"awards" json:"awards"`                   // 新的奖励系统
	MapBranch            map[string]string     `yaml:"mapBranch" json:"mapBranch"`             // mapBranch
}

// SetLinkComponent
func (cfg *FeatureBar2Config) SetLinkComponent(link string, componentName string) {
	if cfg.MapBranch == nil {
		cfg.MapBranch = make(map[string]string)
		cfg.MapBranch[link] = componentName
	} else {
		cfg.MapBranch[link] = componentName
	}
}

type FeatureBar2 struct {
	*BasicComponent `json:"-"`
	Config          *FeatureBar2Config `json:"config"`
}

// Init -
func (featureBar2 *FeatureBar2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("FeatureBar2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &FeatureBar2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("FeatureBar2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return featureBar2.InitEx(cfg, pool)
}

// InitEx -
func (featureBar2 *FeatureBar2) InitEx(cfg any, pool *GamePropertyPool) error {
	featureBar2.Config = cfg.(*FeatureBar2Config)
	featureBar2.Config.ComponentType = FeatureBar2TypeName

	if featureBar2.Config.StrFeatureWeight != "" {
		vw2, err := pool.LoadStrWeights(featureBar2.Config.StrFeatureWeight, true)
		if err != nil {
			goutils.Error("FeatureBar2.InitEx:LoadStrWeights",
				slog.String("FeatureWeight", featureBar2.Config.StrFeatureWeight),
				goutils.Err(err))

			return err
		}

		featureBar2.Config.FeatureWeight = vw2
	} else {
		goutils.Error("FeatureBar2.InitEx:StrFeatureWeight",
			goutils.Err(ErrInvalidWeightVal))

		return ErrInvalidWeightVal
	}

	if featureBar2.Config.EmptyFeature == "" {
		goutils.Error("FeatureBar2.InitEx:EmptyFeature",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if featureBar2.Config.FirstJumpWeight < 0 || featureBar2.Config.FirstJumpWeight > 100 {
		goutils.Error("FeatureBar2.InitEx:FirstJumpWeight",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	for _, award := range featureBar2.Config.Awards {
		award.Init()
	}

	for _, awards := range featureBar2.Config.MapAwards {
		for _, award := range awards {
			award.Init()
		}
	}

	featureBar2.onInit(&featureBar2.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (featureBar2 *FeatureBar2) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	gameProp.procAwards(plugin, featureBar2.Config.Awards, curpr, gp)

	if len(featureBar2.Config.MapAwards) > 0 {
		awards, isok := featureBar2.Config.MapAwards[strVal]
		if isok {
			gameProp.procAwards(plugin, awards, curpr, gp)
		}
	}
}

// getWeight -
func (featureBar2 *FeatureBar2) getWeight(gameProp *GameProperty, cd *FeatureBar2Data) *sgc7game.ValWeights2 {
	str := cd.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadStrWeights(str, true)

		return vw2
	}

	return featureBar2.Config.FeatureWeight
}

// randFirstJump -
func (featureBar2 *FeatureBar2) randFirstJump(_ *GameProperty, _ *FeatureBar2Data, plugin sgc7plugin.IPlugin) (bool, error) {
	if featureBar2.Config.FirstJumpWeight <= 0 {
		return false, nil
	}

	if featureBar2.Config.FirstJumpWeight >= 100 {
		return true, nil
	}

	cr, err := plugin.Random(context.Background(), 100)
	if err != nil {
		goutils.Error("FeatureBar2.randFirstJump:Random",
			goutils.Err(err))

		return false, err
	}

	return cr < featureBar2.Config.FirstJumpWeight, nil
}

// procFeature -
func (featureBar2 *FeatureBar2) procFeature(gameProp *GameProperty, cd *FeatureBar2Data, curpr *sgc7game.PlayResult,
	gp *GameParams, plugin sgc7plugin.IPlugin, featureVW *sgc7game.ValWeights2) error {

	cd.CurFeature = cd.Features[0]
	cd.UsedFeatures = append(cd.UsedFeatures, cd.CurFeature)

	cd.Features = cd.Features[1:]

	feature, err := featureVW.RandVal(plugin)
	if err != nil {
		goutils.Error("FeatureBar2.procFeature:RandVal",
			goutils.Err(err))

		return err
	}

	cd.Features = append(cd.Features, feature.String())

	featureBar2.ProcControllers(gameProp, plugin, curpr, gp, -1, cd.CurFeature)

	return nil
}

func (featureBar2 *FeatureBar2) isClear(basicCD *BasicComponentData) bool {
	clear, isok := basicCD.GetConfigIntVal(CCVClear)
	if isok {
		return clear != 0
	}

	return false
}

// playgame
func (featureBar2 *FeatureBar2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*FeatureBar2Data)

	if featureBar2.isClear(&cd.BasicComponentData) {
		cd.Features = nil
		cd.SetConfigIntVal(CCVClear, 0)
	}

	vw := featureBar2.getWeight(gameProp, cd)

	if len(cd.Features) == 0 {
		cd.CurFeature = featureBar2.Config.EmptyFeature

		for i := 0; i < featureBar2.Config.Length; i++ {
			feature, err := vw.RandVal(plugin)
			if err != nil {
				goutils.Error("FeatureBar2.OnPlayGame:RandVal",
					goutils.Err(err))

				return "", err
			}

			cd.Features = append(cd.Features, feature.String())
		}

		isFirstJump, err := featureBar2.randFirstJump(gameProp, cd, plugin)
		if err != nil {
			goutils.Error("FeatureBar2.OnPlayGame:randFirstJump",
				goutils.Err(err))

			return "", err
		}

		if isFirstJump {
			err = featureBar2.procFeature(gameProp, cd, curpr, gp, plugin, vw)
			if err != nil {
				goutils.Error("FeatureBar2.OnPlayGame:isFirstJump:procFeature",
					goutils.Err(err))

				return "", err
			}
		}
	} else {
		err := featureBar2.procFeature(gameProp, cd, curpr, gp, plugin, vw)
		if err != nil {
			goutils.Error("FeatureBar2.OnPlayGame:procFeature",
				goutils.Err(err))

			return "", err
		}
	}

	nc := featureBar2.onStepEnd(gameProp, curpr, gp, featureBar2.Config.MapBranch[cd.CurFeature])

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (featureBar2 *FeatureBar2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult,
	mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	return nil
}

// NewComponentData -
func (featureBar2 *FeatureBar2) NewComponentData() IComponentData {
	return &FeatureBar2Data{
		cfg: featureBar2.Config,
	}
}

// InitPlayerState -
func (featureBar2 *FeatureBar2) InitPlayerState(pool *GamePropertyPool, gameProp *GameProperty, plugin sgc7plugin.IPlugin,
	ps *PlayerState, betMethod int, bet int) error {

	if featureBar2.Config.IsPlayerState {
		bmd := ps.GetBetMethodPub(betMethod)
		if bet <= 0 {
			return nil
		}

		bps := bmd.GetBetPS(bet)

		_, isok := bps.MapComponentData[featureBar2.GetName()]
		if !isok {
			cps := &FeatureBar2PS{}

			vw := featureBar2.Config.FeatureWeight

			for range featureBar2.Config.Length {
				val, err := vw.RandVal(plugin)
				if err != nil {
					goutils.Error("FeatureBar2.InitPlayerState:RandVal",
						goutils.Err(err))

					return err
				}

				cps.Features = append(cps.Features, val.String())
			}

			bps.MapComponentData[featureBar2.GetName()] = cps
		}
	}

	return nil
}

func NewFeatureBar2(name string) IComponent {
	return &FeatureBar2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "Length": 5,
// "FeatureWeights": "bgfeature",
// "FirstJumpWeight": 5
type jsonFeatureBar2 struct {
	Length           int    `json:"length"`          // bar 的长度
	StrFeatureWeight string `json:"weight"`          // feature权重
	EmptyFeature     string `json:"emptyFeature"`    // emptyFeature
	FirstJumpWeight  int    `json:"FirstJumpWeight"` // firstJump 权重
	IsPlayerState    bool   `json:"IsPlayerState"`   // IsPlayerState
	IsMergeData      bool   `json:"IsMergeData"`     // IsMergeData
}

func (jcfg *jsonFeatureBar2) build() *FeatureBar2Config {
	cfg := &FeatureBar2Config{
		Length:           jcfg.Length,
		StrFeatureWeight: jcfg.StrFeatureWeight,
		FirstJumpWeight:  jcfg.FirstJumpWeight,
		EmptyFeature:     jcfg.EmptyFeature,
		IsPlayerState:    jcfg.IsPlayerState,
		IsMergeData:      jcfg.IsMergeData,
	}

	return cfg
}

func parseFeatureBar2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseFeatureBar2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseFeatureBar2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonFeatureBar2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseFeatureBar2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, mapAwards, err := parseAllAndStrMapControllers(ctrls)
		if err != nil {
			goutils.Error("parseFeatureBar2:parseAllAndStrMapControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: FeatureBar2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

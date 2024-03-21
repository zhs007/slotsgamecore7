package lowcode

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const WeightReelsTypeName = "weightReels"

type WeightReelsData struct {
	BasicComponentData
	ReelSetIndex int // The index of the currently selected reelset
}

// OnNewGame -
func (weightReelsData *WeightReelsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	weightReelsData.BasicComponentData.OnNewGame(gameProp, component)
}

// // OnNewStep -
// func (weightReelsData *WeightReelsData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	weightReelsData.BasicComponentData.OnNewStep(gameProp, component)

// 	weightReelsData.ReelSetIndex = -1
// }

// BuildPBComponentData
func (weightReelsData *WeightReelsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WeightReelsData{
		BasicComponentData: weightReelsData.BuildPBBasicComponentData(),
		ReelSetIndex:       int32(weightReelsData.ReelSetIndex),
	}

	return pbcd
}

// BasicReelsConfig - configuration for WeightReels
type WeightReelsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ReelSetsWeight       string                `yaml:"reelSetWeight" json:"reelSetWeight"`
	ReelSetsWeightVW     *sgc7game.ValWeights2 `json:"-"`
	IsExpandReel         bool                  `yaml:"isExpandReel" json:"isExpandReel"`
}

// SetLinkComponent
func (cfg *WeightReelsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type WeightReels struct {
	*BasicComponent `json:"-"`
	Config          *WeightReelsConfig `json:"config"`
}

// Init -
func (weightReels *WeightReels) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WeightReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &WeightReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return weightReels.InitEx(cfg, pool)
}

// InitEx -
func (weightReels *WeightReels) InitEx(cfg any, pool *GamePropertyPool) error {
	weightReels.Config = cfg.(*WeightReelsConfig)
	weightReels.Config.ComponentType = WeightReelsTypeName

	if weightReels.Config.ReelSetsWeight != "" {
		vw2, err := pool.LoadStrWeights(weightReels.Config.ReelSetsWeight, weightReels.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightReels.Init:LoadValWeights",
				zap.String("ReelSetsWeight", weightReels.Config.ReelSetsWeight),
				zap.Error(err))

			return err
		}

		weightReels.Config.ReelSetsWeightVW = vw2
	}

	weightReels.onInit(&weightReels.Config.BasicComponentConfig)

	return nil
}

// func (weightReels *WeightReels) GetReelSet(basicCD *BasicComponentData) string {
// 	str := basicCD.GetConfigVal(BRCVReelSet)
// 	if str != "" {
// 		return str
// 	}

// 	return weightReels.Config.ReelSetsWeightVW.Vals[]
// }

// playgame
func (weightReels *WeightReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// weightReels.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	wrd := icd.(*WeightReelsData)

	wrd.UsedScenes = nil
	wrd.ReelSetIndex = -1

	reelname := ""
	if weightReels.Config.ReelSetsWeightVW != nil {
		val, si, err := weightReels.Config.ReelSetsWeightVW.RandValEx(plugin)
		if err != nil {
			goutils.Error("WeightReels.OnPlayGame:ReelSetWeights.RandVal",
				zap.Error(err))

			return "", err
		}

		wrd.ReelSetIndex = si

		// weightReels.AddRNG(gameProp, si, &wrd.BasicComponentData)

		curreels := val.String()
		// gameProp.TagStr(TagCurReels, curreels)

		rd, isok := gameProp.Pool.Config.MapReels[curreels]
		if !isok {
			goutils.Error("WeightReels.OnPlayGame:MapReels",
				zap.Error(ErrInvalidReels))

			return "", ErrInvalidReels
		}

		gameProp.CurReels = rd
		reelname = curreels
	}
	// else {
	// 	reelname = weightReels.GetReelSet(cd)
	// 	rd, isok := gameProp.Pool.Config.MapReels[reelname]
	// 	if !isok {
	// 		goutils.Error("BasicReels.OnPlayGame:MapReels",
	// 			zap.Error(ErrInvalidReels))

	// 		return ErrInvalidReels
	// 	}

	// 	gameProp.TagStr(TagCurReels, reelname)

	// 	gameProp.CurReels = rd
	// 	// reelname = basicReels.Config.ReelSet
	// }

	sc := gameProp.PoolScene.New(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	sc.ReelName = reelname
	// sc, err := sgc7game.NewGameScene(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	// if err != nil {
	// 	goutils.Error("BasicReels.OnPlayGame:NewGameScene",
	// 		zap.Error(err))

	// 	return err
	// }

	if weightReels.Config.IsExpandReel {
		sc.RandExpandReelsWithReelData(gameProp.CurReels, plugin)
	} else {
		sc.RandReelsWithReelData(gameProp.CurReels, plugin)
	}

	weightReels.AddScene(gameProp, curpr, sc, &wrd.BasicComponentData)

	nc := weightReels.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(basicReels.Name, cd)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (weightReels *WeightReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	wrd := icd.(*WeightReelsData)

	if len(wrd.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[wrd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (weightReels *WeightReels) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// NewComponentData -
func (weightReels *WeightReels) NewComponentData() IComponentData {
	return &WeightReelsData{}
}

func NewWeightReels(name string) IComponent {
	weightReels := &WeightReels{
		BasicComponent: NewBasicComponent(name, 0),
	}

	return weightReels
}

//	"configuration": {
//		"isExpandReel": "false",
//		"reelSetWeight": "bgreelweight"
//	}
type jsonWeightReels struct {
	ReelSetWeight string `json:"reelSetWeight"`
	IsExpandReel  string `json:"isExpandReel"`
}

func (jwr *jsonWeightReels) build() *WeightReelsConfig {
	cfg := &WeightReelsConfig{
		ReelSetsWeight: jwr.ReelSetWeight,
		IsExpandReel:   jwr.IsExpandReel == "true",
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseWeightReels(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWeightReels:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWeightReels:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonWeightReels{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWeightReels:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WeightReelsTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}

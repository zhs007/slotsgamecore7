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

const WeightReelsTypeName = "weightReels"

type WeightReelsData struct {
	BasicComponentData
	ReelSetIndex int // The index of the currently selected reelset
}

// OnNewGame -
func (weightReelsData *WeightReelsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	weightReelsData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (weightReelsData *WeightReelsData) onNewStep() {
	weightReelsData.UsedScenes = nil
	weightReelsData.ReelSetIndex = -1
}

// GetValEx -
func (weightReelsData *WeightReelsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSelectedIndex {
		return weightReelsData.ReelSetIndex, true
	}

	return 0, false
}

// Clone
func (weightReelsData *WeightReelsData) Clone() IComponentData {
	target := &WeightReelsData{
		BasicComponentData: weightReelsData.CloneBasicComponentData(),
		ReelSetIndex:       weightReelsData.ReelSetIndex,
	}

	return target
}

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
	Awards               []*Award              `yaml:"awards" json:"awards"` // 新的奖励系统
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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WeightReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightReels.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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
				slog.String("ReelSetsWeight", weightReels.Config.ReelSetsWeight),
				goutils.Err(err))

			return err
		}

		weightReels.Config.ReelSetsWeightVW = vw2
	}

	for _, award := range weightReels.Config.Awards {
		award.Init()
	}

	weightReels.onInit(&weightReels.Config.BasicComponentConfig)

	return nil
}

func (weightReels *WeightReels) GetReelSetWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVReelSetWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadStrWeights(str, weightReels.Config.UseFileMapping)

		return vw2
	}

	return weightReels.Config.ReelSetsWeightVW
}

// func (weightReels *WeightReels) GetReelSet(basicCD *BasicComponentData) string {
// 	str := basicCD.GetConfigVal(BRCVReelSet)
// 	if str != "" {
// 		return str
// 	}

// 	return weightReels.Config.ReelSetsWeightVW.Vals[]
// }

// OnProcControllers -
func (weightReels *WeightReels) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(weightReels.Config.Awards) > 0 {
		gameProp.procAwards(plugin, weightReels.Config.Awards, curpr, gp)
	}
}

// playgame
func (weightReels *WeightReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// weightReels.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	wrd := icd.(*WeightReelsData)

	wrd.onNewStep()

	reelname := ""
	vw2 := weightReels.GetReelSetWeight(gameProp, &wrd.BasicComponentData)
	if vw2 != nil {
		val, si, err := vw2.RandValEx(plugin)
		if err != nil {
			goutils.Error("WeightReels.OnPlayGame:ReelSetWeights.RandVal",
				goutils.Err(err))

			return "", err
		}

		wrd.ReelSetIndex = si

		// weightReels.AddRNG(gameProp, si, &wrd.BasicComponentData)

		curreels := val.String()
		// gameProp.TagStr(TagCurReels, curreels)

		rd, isok := gameProp.Pool.Config.MapReels[curreels]
		if !isok {
			goutils.Error("WeightReels.OnPlayGame:MapReels",
				goutils.Err(ErrInvalidReels))

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
	// 			goutils.Err(ErrInvalidReels))

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
	// 		goutils.Err(err))

	// 	return err
	// }

	if weightReels.Config.IsExpandReel {
		sc.RandExpandReelsWithReelData(gameProp.CurReels, plugin)
	} else {
		sc.RandReelsWithReelData(gameProp.CurReels, plugin)
	}

	weightReels.AddScene(gameProp, curpr, sc, &wrd.BasicComponentData)

	weightReels.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

	nc := weightReels.onStepEnd(gameProp, curpr, gp, "")

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
	IsExpandReel  bool   `json:"isExpandReel"`
}

func (jwr *jsonWeightReels) build() *WeightReelsConfig {
	cfg := &WeightReelsConfig{
		ReelSetsWeight: jwr.ReelSetWeight,
		IsExpandReel:   jwr.IsExpandReel,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

type jsonWeightReelsT struct {
	ReelSetWeight string `json:"reelSetWeight"`
	IsExpandReel  string `json:"isExpandReel"`
}

func (jwr *jsonWeightReelsT) build() *WeightReelsConfig {
	cfg := &WeightReelsConfig{
		ReelSetsWeight: jwr.ReelSetWeight,
		IsExpandReel:   jwr.IsExpandReel == "true",
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseWeightReels(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWeightReels:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWeightReels:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWeightReels{}
	var cfgd *WeightReelsConfig

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		data2 := &jsonWeightReelsT{}

		err = sonic.Unmarshal(buf, data2)
		if err != nil {
			goutils.Error("parseWeightReels:Unmarshal",
				goutils.Err(err))

			return "", err
		}

		cfgd = data2.build()
	} else {
		cfgd = data.build()
	}

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseWeightReels:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WeightReelsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

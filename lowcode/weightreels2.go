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

const WeightReels2TypeName = "weightReels2"

type WeightReels2Data struct {
	BasicComponentData
	ReelSetIndex int // The index of the currently selected reelset
}

// OnNewGame -
func (weightReels2Data *WeightReels2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	weightReels2Data.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (weightReels2Data *WeightReels2Data) onNewStep() {
	weightReels2Data.UsedScenes = nil
	weightReels2Data.ReelSetIndex = -1
}

// GetValEx -
func (weightReels2Data *WeightReels2Data) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSelectedIndex {
		return weightReels2Data.ReelSetIndex, true
	}

	return 0, false
}

// Clone
func (weightReels2Data *WeightReels2Data) Clone() IComponentData {
	target := &WeightReels2Data{
		BasicComponentData: weightReels2Data.CloneBasicComponentData(),
		ReelSetIndex:       weightReels2Data.ReelSetIndex,
	}

	return target
}

// BuildPBComponentData
func (weightReels2Data *WeightReels2Data) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WeightReelsData{
		BasicComponentData: weightReels2Data.BuildPBBasicComponentData(),
		ReelSetIndex:       int32(weightReels2Data.ReelSetIndex),
	}

	return pbcd
}

// BasicReelsConfig - configuration for WeightReels
type WeightReels2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ReelSetsWeight       string                `yaml:"reelSetWeight" json:"reelSetWeight"`
	ReelSetsWeightVW     *sgc7game.ValWeights2 `json:"-"`
	IsExpandReel         bool                  `yaml:"isExpandReel" json:"isExpandReel"`
	MapAwards            map[string][]*Award   `yaml:"mapAwards" json:"mapAwards"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *WeightReels2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type WeightReels2 struct {
	*BasicComponent `json:"-"`
	Config          *WeightReels2Config `json:"config"`
}

// Init -
func (weightReels2 *WeightReels2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WeightReels2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WeightReels2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightReels2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return weightReels2.InitEx(cfg, pool)
}

// InitEx -
func (weightReels2 *WeightReels2) InitEx(cfg any, pool *GamePropertyPool) error {
	weightReels2.Config = cfg.(*WeightReels2Config)
	weightReels2.Config.ComponentType = WeightReels2TypeName

	if weightReels2.Config.ReelSetsWeight != "" {
		vw2, err := pool.LoadStrWeights(weightReels2.Config.ReelSetsWeight, weightReels2.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightReels2.Init:LoadValWeights",
				slog.String("ReelSetsWeight", weightReels2.Config.ReelSetsWeight),
				goutils.Err(err))

			return err
		}

		weightReels2.Config.ReelSetsWeightVW = vw2
	}

	for _, awards := range weightReels2.Config.MapAwards {
		for _, award := range awards {
			award.Init()
		}
	}

	weightReels2.onInit(&weightReels2.Config.BasicComponentConfig)

	return nil
}

func (weightReels2 *WeightReels2) GetReelSetWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVReelSetWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadStrWeights(str, weightReels2.Config.UseFileMapping)

		return vw2
	}

	return weightReels2.Config.ReelSetsWeightVW
}

// OnProcControllers -
func (weightReels2 *WeightReels2) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if strVal != "" {
		awards, isok := weightReels2.Config.MapAwards[strVal]
		if isok {
			gameProp.procAwards(plugin, awards, curpr, gp)
		}
	}
}

// playgame
func (weightReels2 *WeightReels2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	wrd := icd.(*WeightReels2Data)

	wrd.onNewStep()

	reelname := ""
	vw2 := weightReels2.GetReelSetWeight(gameProp, &wrd.BasicComponentData)
	if vw2 != nil {
		val, si, err := vw2.RandValEx(plugin)
		if err != nil {
			goutils.Error("WeightReels2.OnPlayGame:ReelSetWeights.RandVal",
				goutils.Err(err))

			return "", err
		}

		wrd.ReelSetIndex = si

		curreels := val.String()

		rd, isok := gameProp.Pool.Config.MapReels[curreels]
		if !isok {
			goutils.Error("WeightReels2.OnPlayGame:MapReels",
				goutils.Err(ErrInvalidReels))

			return "", ErrInvalidReels
		}

		gameProp.CurReels = rd
		reelname = curreels
	}

	sc := gameProp.PoolScene.New(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	sc.ReelName = reelname

	if weightReels2.Config.IsExpandReel {
		sc.RandExpandReelsWithReelData(gameProp.CurReels, plugin)
	} else {
		sc.RandReelsWithReelData(gameProp.CurReels, plugin)
	}

	weightReels2.AddScene(gameProp, curpr, sc, &wrd.BasicComponentData)

	weightReels2.ProcControllers(gameProp, plugin, curpr, gp, -1, "<any>")
	weightReels2.ProcControllers(gameProp, plugin, curpr, gp, -1, reelname)

	nc := weightReels2.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (weightReels2 *WeightReels2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	wrd := icd.(*WeightReels2Data)

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
func (weightReels2 *WeightReels2) NewComponentData() IComponentData {
	return &WeightReels2Data{}
}

func NewWeightReels2(name string) IComponent {
	weightReels := &WeightReels2{
		BasicComponent: NewBasicComponent(name, 0),
	}

	return weightReels
}

//	"configuration": {
//		"isExpandReel": "false",
//		"reelSetWeight": "bgreelweight"
//	}
type jsonWeightReels2 struct {
	ReelSetWeight string `json:"reelSetWeight"`
	IsExpandReel  bool   `json:"isExpandReel"`
}

func (jwr *jsonWeightReels2) build() *WeightReels2Config {
	cfg := &WeightReels2Config{
		ReelSetsWeight: jwr.ReelSetWeight,
		IsExpandReel:   jwr.IsExpandReel,
	}

	return cfg
}

func parseWeightReels2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWeightReels2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWeightReels2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWeightReels2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWeightReels2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseWeightReels2:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WeightReels2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

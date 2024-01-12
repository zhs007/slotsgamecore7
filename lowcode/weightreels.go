package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
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
func (weightReelsData *WeightReelsData) OnNewGame() {
	weightReelsData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (weightReelsData *WeightReelsData) OnNewStep() {
	weightReelsData.BasicComponentData.OnNewStep()

	weightReelsData.ReelSetIndex = -1
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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	weightReels.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	wrd := gameProp.MapComponentData[weightReels.Name].(*WeightReelsData)

	reelname := ""
	if weightReels.Config.ReelSetsWeightVW != nil {
		val, si, err := weightReels.Config.ReelSetsWeightVW.RandValEx(plugin)
		if err != nil {
			goutils.Error("WeightReels.OnPlayGame:ReelSetWeights.RandVal",
				zap.Error(err))

			return err
		}

		wrd.ReelSetIndex = si

		// weightReels.AddRNG(gameProp, si, &wrd.BasicComponentData)

		curreels := val.String()
		// gameProp.TagStr(TagCurReels, curreels)

		rd, isok := gameProp.Pool.Config.MapReels[curreels]
		if !isok {
			goutils.Error("WeightReels.OnPlayGame:MapReels",
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
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

	sc := gameProp.PoolScene.New(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), false)
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

	weightReels.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(basicReels.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (weightReels *WeightReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	wrd := gameProp.MapComponentData[weightReels.Name].(*WeightReelsData)

	if len(wrd.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[wrd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (weightReels *WeightReels) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

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

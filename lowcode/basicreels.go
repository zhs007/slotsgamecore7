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
	"gopkg.in/yaml.v2"
)

const BasicReelsTypeName = "basicReels"

// const (
// 	BRCVReelSet string = "reelSet" // 可以修改配置项里的ReelSet
// )

// BasicReelsConfig - configuration for BasicReels
type BasicReelsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ReelSet              string   `yaml:"reelSet" json:"reelSet"`
	IsExpandReel         bool     `yaml:"isExpandReel" json:"isExpandReel"`
	Awards               []*Award `yaml:"awards" json:"awards"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *BasicReelsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type BasicReels struct {
	*BasicComponent `json:"-"`
	Config          *BasicReelsConfig `json:"config"`
}

// Init -
func (basicReels *BasicReels) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &BasicReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return basicReels.InitEx(cfg, pool)
}

// InitEx -
func (basicReels *BasicReels) InitEx(cfg any, pool *GamePropertyPool) error {
	basicReels.Config = cfg.(*BasicReelsConfig)
	basicReels.Config.ComponentType = BasicReelsTypeName

	// if basicReels.Config.ReelSetsWeight != "" {
	// 	vw2, err := pool.LoadStrWeights(basicReels.Config.ReelSetsWeight, basicReels.Config.UseFileMapping)
	// 	if err != nil {
	// 		goutils.Error("BasicReels.Init:LoadValWeights",
	// 			slog.String("ReelSetsWeight", basicReels.Config.ReelSetsWeight),
	// 			goutils.Err(err))

	// 		return err
	// 	}

	// 	basicReels.ReelSetWeights = vw2
	// }

	for _, award := range basicReels.Config.Awards {
		award.Init()
	}

	basicReels.onInit(&basicReels.Config.BasicComponentConfig)

	return nil
}

func (basicReels *BasicReels) getReelSet(basicCD *BasicComponentData) string {
	str := basicCD.GetConfigVal(CCVReelSet)
	if str != "" {
		return str
	}

	return basicReels.Config.ReelSet
}

// OnProcControllers -
func (basicReels *BasicReels) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(basicReels.Config.Awards) > 0 {
		gameProp.procAwards(plugin, basicReels.Config.Awards, curpr, gp)
	}
}

// playgame
func (basicReels *BasicReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// basicReels.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	bcd.UsedScenes = nil

	// bcd.OnNewStep()

	// reelname := ""
	// if basicReels.ReelSetWeights != nil {
	// 	val, _, err := basicReels.ReelSetWeights.RandValEx(plugin)
	// 	if err != nil {
	// 		goutils.Error("BasicReels.OnPlayGame:ReelSetWeights.RandVal",
	// 			goutils.Err(err))

	// 		return "", err
	// 	}

	// 	// basicReels.AddRNG(gameProp, si, bcd)

	// 	curreels := val.String()
	// 	gameProp.TagStr(TagCurReels, curreels)

	// 	rd, isok := gameProp.Pool.Config.MapReels[curreels]
	// 	if !isok {
	// 		goutils.Error("BasicReels.OnPlayGame:MapReels",
	// 			goutils.Err(ErrInvalidReels))

	// 		return "", ErrInvalidReels
	// 	}

	// 	gameProp.CurReels = rd
	// 	reelname = curreels
	// } else {
	reelname := basicReels.getReelSet(bcd)
	rd, isok := gameProp.Pool.Config.MapReels[reelname]
	if !isok {
		goutils.Error("BasicReels.OnPlayGame:MapReels",
			goutils.Err(ErrInvalidReels))

		return "", ErrInvalidReels
	}

	gameProp.TagStr(TagCurReels, reelname)

	gameProp.CurReels = rd
	// reelname = basicReels.Config.ReelSet
	// }

	sc := gameProp.PoolScene.New(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	sc.ReelName = reelname
	// sc, err := sgc7game.NewGameScene(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	// if err != nil {
	// 	goutils.Error("BasicReels.OnPlayGame:NewGameScene",
	// 		goutils.Err(err))

	// 	return err
	// }

	if basicReels.Config.IsExpandReel {
		sc.RandExpandReelsWithReelData(gameProp.CurReels, plugin)
	} else {
		sc.RandReelsWithReelData(gameProp.CurReels, plugin)
	}

	basicReels.AddScene(gameProp, curpr, sc, bcd)

	basicReels.ProcControllers(gameProp, plugin, curpr, gp, -1, "")
	// if len(basicReels.Config.Awards) > 0 {
	// 	gameProp.procAwards(plugin, basicReels.Config.Awards, curpr, gp)
	// }

	nc := basicReels.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (basicReels *BasicReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (basicReels *BasicReels) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewBasicReels(name string) IComponent {
	basicReels := &BasicReels{
		BasicComponent: NewBasicComponent(name, 0),
	}

	return basicReels
}

//	"configuration": {
//		"isExpandReel": "false",
//		"reelSet": "bgreelweight"
//	}
type jsonBasicReels struct {
	ReelSet      string `json:"reelSet"`
	IsExpandReel bool   `json:"isExpandReel"`
}

func (jbr *jsonBasicReels) build() *BasicReelsConfig {
	cfg := &BasicReelsConfig{
		ReelSet:      jbr.ReelSet,
		IsExpandReel: jbr.IsExpandReel,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

type jsonBasicReels2 struct {
	ReelSet      string `json:"reelSet"`
	IsExpandReel string `json:"isExpandReel"`
}

func (jbr *jsonBasicReels2) build() *BasicReelsConfig {
	cfg := &BasicReelsConfig{
		ReelSet:      jbr.ReelSet,
		IsExpandReel: jbr.IsExpandReel == "true",
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseBasicReels(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseBasicReels:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseBasicReels:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonBasicReels{}
	var cfgd *BasicReelsConfig

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		data2 := &jsonBasicReels2{}

		err = sonic.Unmarshal(buf, data2)
		if err != nil {
			goutils.Error("parseBasicReels:Unmarshal",
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
			goutils.Error("parseBasicReels:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: BasicReelsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

package lowcode

import (
	"log/slog"
	"os"

	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const ControllerWorkerTypeName = "controllerWorker"

// ControllerWorkerConfig - configuration for ControllerWorker
type ControllerWorkerConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Awards               []*Award `yaml:"awards" json:"awards"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *ControllerWorkerConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type ControllerWorker struct {
	*BasicComponent `json:"-"`
	Config          *ControllerWorkerConfig `json:"config"`
}

// Init -
func (controllerWorker *ControllerWorker) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ControllerWorker.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ReRollReelConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ControllerWorker.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return controllerWorker.InitEx(cfg, pool)
}

// InitEx -
func (controllerWorker *ControllerWorker) InitEx(cfg any, pool *GamePropertyPool) error {
	controllerWorker.Config = cfg.(*ControllerWorkerConfig)
	controllerWorker.Config.ComponentType = ControllerWorkerTypeName

	for _, award := range controllerWorker.Config.Awards {
		award.Init()
	}

	controllerWorker.onInit(&controllerWorker.Config.BasicComponentConfig)

	return nil
}

// playgame
func (controllerWorker *ControllerWorker) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// reRollReel.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	if len(controllerWorker.Config.Awards) > 0 {
		gameProp.procAwards(plugin, controllerWorker.Config.Awards, curpr, gp)
	}

	nc := controllerWorker.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (controllerWorker *ControllerWorker) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	// cd := icd.(*BasicComponentData)

	// if len(cd.UsedScenes) > 0 {
	// 	asciigame.OutputScene("after reRollReel", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	// }

	return nil
}

// // OnStats
// func (reRollReel *ReRollReel) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewControllerWorker(name string) IComponent {
	return &ControllerWorker{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "configuration": {
// }
type jsonControllerWorker struct {
}

func (jbr *jsonControllerWorker) build() *ControllerWorkerConfig {
	cfg := &ControllerWorkerConfig{}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseControllerWorker(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	_, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseControllerWorker:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	// buf, err := cfg.MarshalJSON()
	// if err != nil {
	// 	goutils.Error("parseControllerWorker:MarshalJSON",
	// 		goutils.Err(err))

	// 	return "", err
	// }

	// data := &jsonControllerWorker{}

	// err = sonic.Unmarshal(buf, data)
	// if err != nil {
	// 	goutils.Error("parseControllerWorker:Unmarshal",
	// 		goutils.Err(err))

	// 	return "", err
	// }

	// cfgd := data.build()
	cfgd := &ControllerWorkerConfig{}

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseControllerWorker:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ControllerWorkerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

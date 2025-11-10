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

const BasicReels2TypeName = "basicReels2"

// BasicReels2Config - configuration for BasicReels2
type BasicReels2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ReelSet              string   `yaml:"reelSet" json:"reelSet"`
	IsExpandReel         bool     `yaml:"isExpandReel" json:"isExpandReel"`
	Height               int      `yaml:"height" json:"height"`
	MaskX                string   `yaml:"maskX" json:"maskX"`
	MaskY                string   `yaml:"maskY" json:"maskY"`
	Controllers          []*Award `yaml:"controllers" json:"controllers"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *BasicReels2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type BasicReels2 struct {
	*BasicComponent `json:"-"`
	Config          *BasicReels2Config `json:"config"`
}

// Init -
func (basicReels2 *BasicReels2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &BasicReels2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return basicReels2.InitEx(cfg, pool)
}

// InitEx -
func (basicReels2 *BasicReels2) InitEx(cfg any, pool *GamePropertyPool) error {
	basicReels2.Config = cfg.(*BasicReels2Config)
	basicReels2.Config.ComponentType = BasicReels2TypeName

	for _, ctrl := range basicReels2.Config.Controllers {
		ctrl.Init()
	}

	basicReels2.onInit(&basicReels2.Config.BasicComponentConfig)

	return nil
}

func (basicReels2 *BasicReels2) getReelSet(basicCD *BasicComponentData) string {
	str := basicCD.GetConfigVal(CCVReelSet)
	if str != "" {
		return str
	}

	return basicReels2.Config.ReelSet
}

func (basicReels2 *BasicReels2) getHeight(basicCD *BasicComponentData) int {
	v, isok := basicCD.GetConfigIntVal(CCVHeight)
	if isok {
		return v
	}

	return basicReels2.Config.Height
}

// OnProcControllers -
func (basicReels2 *BasicReels2) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(basicReels2.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, basicReels2.Config.Controllers, curpr, gp)
	}
}

func (basicReels2 *BasicReels2) getMaskX(gameProp *GameProperty, basicCD *BasicComponentData) string {
	str := basicCD.GetConfigVal(CCVMaskX)
	if str != "" {
		if str == "<empty>" {
			return ""
		}

		return str
	}

	return basicReels2.Config.MaskX
}

// playgame
func (basicReels2 *BasicReels2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	bcd := cd.(*BasicComponentData)

	bcd.UsedScenes = nil

	reelname := basicReels2.getReelSet(bcd)
	rd, isok := gameProp.Pool.Config.MapReels[reelname]
	if !isok {
		goutils.Error("BasicReels2.OnPlayGame:MapReels",
			goutils.Err(ErrInvalidReels))

		return "", ErrInvalidReels
	}

	gameProp.TagStr(TagCurReels, reelname)

	gameProp.CurReels = rd

	sc := gameProp.PoolScene.New(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	sc.ReelName = reelname

	height := basicReels2.getHeight(bcd)

	if basicReels2.Config.IsExpandReel {
		sc.RandExpandReelsWithReelData(gameProp.CurReels, plugin)
	} else {
		if height <= 0 && height > sc.Height {
			goutils.Error("BasicReels2.OnPlayGame:MapReels",
				goutils.Err(ErrInvalidReels))

			return "", ErrInvalidReels
		}

		maskx := basicReels2.getMaskX(gameProp, bcd)

		if maskx != "" {
			imaskd := gameProp.GetComponentDataWithName(basicReels2.Config.MaskX)
			if imaskd != nil {
				arr := imaskd.GetMask()
				if len(arr) != sc.Width {
					goutils.Error("BasicReels2.OnPlayGame:MaskX:len(arr)!=gs.Width",
						goutils.Err(ErrInvalidComponentConfig))

					return "", ErrInvalidComponentConfig
				}

				sc.RandReelsWithReelDataMaskAndHeight(gameProp.CurReels, height, arr, plugin)

			} else {
				goutils.Error("BasicReels2.OnPlayGame:MaskX",
					slog.String("maskX", maskx),
					goutils.Err(ErrInvalidComponentConfig))

				return "", ErrInvalidComponentConfig
			}
		} else {
			sc.RandReelsWithReelDataAndHeight(gameProp.CurReels, height, plugin)
		}
	}

	basicReels2.AddScene(gameProp, curpr, sc, bcd)

	basicReels2.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

	nc := basicReels2.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (basicReels2 *BasicReels2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

func NewBasicReels2(name string) IComponent {
	basicReels2 := &BasicReels2{
		BasicComponent: NewBasicComponent(name, 0),
	}

	return basicReels2
}

// "isExpandReel": false,
// "height": 6,
// "reelSet": "bg-reel01",
// "maskX": "mask-6"
type jsonBasicReels2 struct {
	ReelSet      string `json:"reelSet"`
	IsExpandReel bool   `json:"isExpandReel"`
	Height       int    `json:"height"`
	MaskX        string `json:"maskX"`
	MaskY        string `json:"maskY"`
}

func (jbr *jsonBasicReels2) build() *BasicReels2Config {
	cfg := &BasicReels2Config{
		ReelSet:      jbr.ReelSet,
		IsExpandReel: jbr.IsExpandReel,
		Height:       jbr.Height,
		MaskX:        jbr.MaskX,
		MaskY:        jbr.MaskY,
	}

	return cfg
}

func parseBasicReels2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseBasicReels2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseBasicReels2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonBasicReels2{}
	var cfgd *BasicReels2Config

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseBasicReels2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd = data.build()

	if ctrls != nil {
		controllers, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseBasicReels2:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = controllers
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: BasicReels2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

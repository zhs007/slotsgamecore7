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
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const RefillSymbolsTypeName = "refillSymbols"

// RefillSymbolsConfig - configuration for RefillSymbols
type RefillSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	IsNeedProcSymbolVals bool   `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int    `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
	DefaultSymbolVal     int    `yaml:"defaultSymbolVal" json:"defaultSymbolVal"`         // 重新填充的symbolVal是什么
	OutputToComponent    string `yaml:"outputToComponent" json:"outputToComponent"`       // 输出到哪个组件
}

// SetLinkComponent
func (cfg *RefillSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RefillSymbols struct {
	*BasicComponent `json:"-"`
	Config          *RefillSymbolsConfig `json:"config"`
}

// Init -
func (refillSymbols *RefillSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RefillSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RefillSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RefillSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return refillSymbols.InitEx(cfg, pool)
}

// InitEx -
func (refillSymbols *RefillSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	refillSymbols.Config = cfg.(*RefillSymbolsConfig)
	refillSymbols.Config.ComponentType = RefillSymbolsTypeName

	refillSymbols.onInit(&refillSymbols.Config.BasicComponentConfig)

	return nil
}

func (refillSymbols *RefillSymbols) getSymbol(rd *sgc7game.ReelsData, x int, index int) int {
	index--

	for ; index < 0; index += len(rd.Reels[x]) {
	}

	return rd.Reels[x][index]
}

// playgame
func (refillSymbols *RefillSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	bcd := cd.(*BasicComponentData)

	bcd.UsedScenes = nil
	bcd.UsedOtherScenes = nil

	gs := refillSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	ngs := gs

	var os *sgc7game.GameScene
	if refillSymbols.Config.IsNeedProcSymbolVals {
		os = refillSymbols.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	var outputCD IComponentData
	if refillSymbols.Config.OutputToComponent != "" {
		outputCD = gameProp.GetComponentDataWithName(refillSymbols.Config.OutputToComponent)
		if outputCD == nil {
			goutils.Error("RefillSymbols.OnPlayGame:OutputToComponent",
				slog.String("outputToComponent", refillSymbols.Config.OutputToComponent),
				goutils.Err(ErrInvalidComponent))

			return "", ErrInvalidComponent
		}

		outputCD.ClearPos()
	}

	if os != nil {
		nos := os

		for x := 0; x < gs.Width; x++ {
			for y := gs.Height - 1; y >= 0; y-- {
				if ngs.Arr[x][y] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
						nos = os.CloneEx(gameProp.PoolScene)
					}

					cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

					ngs.Arr[x][y] = refillSymbols.getSymbol(cr, x, ngs.Indexes[x])
					ngs.Indexes[x]--

					nos.Arr[x][y] = refillSymbols.Config.DefaultSymbolVal

					if outputCD != nil {
						outputCD.AddPos(x, y)
					}
				}
			}
		}

		if ngs == gs {
			nc := refillSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		refillSymbols.AddOtherScene(gameProp, curpr, nos, bcd)
	} else {
		for x := 0; x < gs.Width; x++ {
			for y := gs.Height - 1; y >= 0; y-- {
				if ngs.Arr[x][y] == -1 {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

					ngs.Arr[x][y] = refillSymbols.getSymbol(cr, x, ngs.Indexes[x])
					ngs.Indexes[x]--

					if outputCD != nil {
						outputCD.AddPos(x, y)
					}
				}
			}
		}

		if ngs == gs {
			nc := refillSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	}

	refillSymbols.AddScene(gameProp, curpr, ngs, bcd)

	nc := refillSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (refillSymbols *RefillSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after refillSymbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// EachUsedResults -
func (refillSymbols *RefillSymbols) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewRefillSymbols(name string) IComponent {
	return &RefillSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "isNeedProcSymbolVals": true,
// "defaultSymbolVal": 0,
// "outputToComponent": "bg-pos-collect"
type jsonRefillSymbols struct {
	IsNeedProcSymbolVals bool   `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int    `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
	DefaultSymbolVal     int    `yaml:"defaultSymbolVal" json:"defaultSymbolVal"`         // 重新填充的symbolVal是什么
	OutputToComponent    string `yaml:"outputToComponent" json:"outputToComponent"`       // 输出到哪个组件
}

func (jcfg *jsonRefillSymbols) build() *RefillSymbolsConfig {
	cfg := &RefillSymbolsConfig{
		IsNeedProcSymbolVals: jcfg.IsNeedProcSymbolVals,
		EmptySymbolVal:       jcfg.EmptySymbolVal,
		DefaultSymbolVal:     jcfg.DefaultSymbolVal,
		OutputToComponent:    jcfg.OutputToComponent,
	}

	return cfg
}

func parseRefillSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRefillSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRefillSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRefillSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRefillSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RefillSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

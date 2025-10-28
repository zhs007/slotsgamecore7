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

const RefillSymbols2TypeName = "refillSymbols2"

// RefillSymbols2Config - configuration for RefillSymbols2
type RefillSymbols2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	IsNeedProcSymbolVals bool   `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int    `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
	DefaultSymbolVal     int    `yaml:"defaultSymbolVal" json:"defaultSymbolVal"`         // 重新填充的symbolVal是什么
	OutputToComponent    string `yaml:"outputToComponent" json:"outputToComponent"`       // 输出到哪个组件
	Height               int    `yaml:"height" json:"height"`                             // 重新填充的symbolVal是什么
}

// SetLinkComponent
func (cfg *RefillSymbols2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RefillSymbols2 struct {
	*BasicComponent `json:"-"`
	Config          *RefillSymbols2Config `json:"config"`
}

// Init -
func (refillSymbols2 *RefillSymbols2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RefillSymbols2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RefillSymbols2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RefillSymbols2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return refillSymbols2.InitEx(cfg, pool)
}

// InitEx -
func (refillSymbols2 *RefillSymbols2) InitEx(cfg any, pool *GamePropertyPool) error {
	refillSymbols2.Config = cfg.(*RefillSymbols2Config)
	refillSymbols2.Config.ComponentType = RefillSymbols2TypeName

	refillSymbols2.onInit(&refillSymbols2.Config.BasicComponentConfig)

	return nil
}

func (refillSymbols2 *RefillSymbols2) getSymbol(rd *sgc7game.ReelsData, x int, index int) int {
	index--

	for ; index < 0; index += len(rd.Reels[x]) {
	}

	return rd.Reels[x][index]
}

func (refillSymbols2 *RefillSymbols2) GetHeight(basicCD *BasicComponentData) int {
	height, isok := basicCD.GetConfigIntVal(CCVHeight)
	if isok {
		return height
	}

	return refillSymbols2.Config.Height
}

// playgame
func (refillSymbols2 *RefillSymbols2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	bcd := cd.(*BasicComponentData)

	bcd.UsedScenes = nil
	bcd.UsedOtherScenes = nil

	gs := refillSymbols2.GetTargetScene3(gameProp, curpr, prs, 0)
	ngs := gs

	var os *sgc7game.GameScene
	if refillSymbols2.Config.IsNeedProcSymbolVals {
		os = refillSymbols2.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	var outputCD IComponentData
	if refillSymbols2.Config.OutputToComponent != "" {
		outputCD = gameProp.GetComponentDataWithName(refillSymbols2.Config.OutputToComponent)
		if outputCD == nil {
			goutils.Error("RefillSymbols2.OnPlayGame:OutputToComponent",
				slog.String("outputToComponent", refillSymbols2.Config.OutputToComponent),
				goutils.Err(ErrInvalidComponent))

			return "", ErrInvalidComponent
		}

		outputCD.ClearPos()
	}

	height := refillSymbols2.GetHeight(bcd)

	if height > 0 && height < gs.Height {
		if os != nil {
			nos := os

			for x := 0; x < gs.Width; x++ {
				for y := gs.Height - 1; y >= 0; y-- {
					if y < gs.Height-height {
						continue
					}

					if ngs.Arr[x][y] == -1 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
							nos = os.CloneEx(gameProp.PoolScene)
						}

						cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

						ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
						ngs.Indexes[x]--

						nos.Arr[x][y] = refillSymbols2.Config.DefaultSymbolVal

						if outputCD != nil {
							outputCD.AddPos(x, y)
						}
					}
				}
			}

			if ngs == gs {
				nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

				return nc, ErrComponentDoNothing
			}

			refillSymbols2.AddOtherScene(gameProp, curpr, nos, bcd)
		} else {
			for x := 0; x < gs.Width; x++ {
				for y := gs.Height - 1; y >= 0; y-- {
					if y < gs.Height-height {
						continue
					}

					if ngs.Arr[x][y] == -1 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

						ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
						ngs.Indexes[x]--

						if outputCD != nil {
							outputCD.AddPos(x, y)
						}
					}
				}
			}

			if ngs == gs {
				nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

				return nc, ErrComponentDoNothing
			}
		}
	} else {
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

						ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
						ngs.Indexes[x]--

						nos.Arr[x][y] = refillSymbols2.Config.DefaultSymbolVal

						if outputCD != nil {
							outputCD.AddPos(x, y)
						}
					}
				}
			}

			if ngs == gs {
				nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

				return nc, ErrComponentDoNothing
			}

			refillSymbols2.AddOtherScene(gameProp, curpr, nos, bcd)
		} else {
			for x := 0; x < gs.Width; x++ {
				for y := gs.Height - 1; y >= 0; y-- {
					if ngs.Arr[x][y] == -1 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

						ngs.Arr[x][y] = refillSymbols2.getSymbol(cr, x, ngs.Indexes[x])
						ngs.Indexes[x]--

						if outputCD != nil {
							outputCD.AddPos(x, y)
						}
					}
				}
			}

			if ngs == gs {
				nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

				return nc, ErrComponentDoNothing
			}
		}
	}

	refillSymbols2.AddScene(gameProp, curpr, ngs, bcd)

	nc := refillSymbols2.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (refillSymbols2 *RefillSymbols2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after refillSymbols2", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// EachUsedResults -
func (refillSymbols2 *RefillSymbols2) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewRefillSymbols2(name string) IComponent {
	return &RefillSymbols2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "isNeedProcSymbolVals": false,
// "emptySymbolVal": -1,
// "defaultSymbolVal": 0,
// "Height": 4
type jsonRefillSymbols2 struct {
	IsNeedProcSymbolVals bool   `json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int    `json:"emptySymbolVal"`       // 空的symbolVal是什么
	DefaultSymbolVal     int    `json:"defaultSymbolVal"`     // 重新填充的symbolVal是什么
	OutputToComponent    string `json:"outputToComponent"`    // 输出到哪个组件
	Height               int    `json:"Height"`               // height, <=0 is ignore
}

func (jcfg *jsonRefillSymbols2) build() *RefillSymbols2Config {
	cfg := &RefillSymbols2Config{
		IsNeedProcSymbolVals: jcfg.IsNeedProcSymbolVals,
		EmptySymbolVal:       jcfg.EmptySymbolVal,
		DefaultSymbolVal:     jcfg.DefaultSymbolVal,
		OutputToComponent:    jcfg.OutputToComponent,
		Height:               jcfg.Height,
	}

	return cfg
}

func parseRefillSymbols2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRefillSymbols2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRefillSymbols2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRefillSymbols2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRefillSymbols2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RefillSymbols2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

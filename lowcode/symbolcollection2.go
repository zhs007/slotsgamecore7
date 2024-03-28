package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const SymbolCollection2TypeName = "symbolCollection2"

type SymbolCollection2Data struct {
	BasicComponentData
	SymbolCodes  []int
	TotalCoinWin int64
	TotalCashWin int64
}

// OnNewGame -
func (symbolCollection2Data *SymbolCollection2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	symbolCollection2Data.BasicComponentData.OnNewGame(gameProp, component)

	symbolCollection2 := component.(*SymbolCollection2)

	symbolCollection2Data.SymbolCodes = nil

	symbolCollection2Data.SymbolCodes = append(symbolCollection2Data.SymbolCodes, symbolCollection2.Config.InitSymbolCodes...)
}

// // OnNewStep -
// func (symbolCollection2Data *SymbolCollection2Data) onNewStep() {
// 	// symbolCollection2Data.BasicComponentData.OnNewStep(gameProp, component)
// 	// symbolCollection2Data
// }

// BuildPBComponentData
func (symbolCollection2Data *SymbolCollection2Data) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolCollection2Data{
		BasicComponentData: symbolCollection2Data.BuildPBBasicComponentData(),
	}

	for _, s := range symbolCollection2Data.SymbolCodes {
		pbcd.SymbolCodes = append(pbcd.SymbolCodes, int32(s))
	}

	return pbcd
}

// GetSymbols -
func (symbolCollection2Data *SymbolCollection2Data) GetSymbols() []int {
	return symbolCollection2Data.SymbolCodes
}

// AddSymbol -
func (symbolCollection2Data *SymbolCollection2Data) AddSymbol(symbolCode int) {
	symbolCollection2Data.SymbolCodes = append(symbolCollection2Data.SymbolCodes, symbolCode)
}

// SymbolCollection2Config - configuration for SymbolCollection2 feature
type SymbolCollection2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ForeachComponent     string   `yaml:"foreachComponent" json:"foreachComponent"` // foreach
	MaxSymbolNum         int      `yaml:"maxSymbolNum" json:"maxSymbolNum"`         // 0表示不限制
	InitSymbols          []string `yaml:"initSymbols" json:"initSymbols"`           // 初始化symbols
	InitSymbolCodes      []int    `yaml:"-" json:"-"`                               // 初始化symbols
	Children             []string `yaml:"-" json:"-"`                               //
}

// SetLinkComponent
func (cfg *SymbolCollection2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "foreach" {
		cfg.ForeachComponent = componentName
	}
}

// SymbolCollection2 - 也是一个非常特殊的组件，symbol集合
type SymbolCollection2 struct {
	*BasicComponent `json:"-"`
	Config          *SymbolCollection2Config `json:"config"`
}

// Init -
func (symbolCollection2 *SymbolCollection2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolCollection2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &SymbolCollection2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolCollection2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return symbolCollection2.InitEx(cfg, pool)
}

// InitEx -
func (symbolCollection2 *SymbolCollection2) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolCollection2.Config = cfg.(*SymbolCollection2Config)
	symbolCollection2.Config.ComponentType = SymbolCollection2TypeName

	for _, v := range symbolCollection2.Config.InitSymbols {
		symbolCollection2.Config.InitSymbolCodes = append(symbolCollection2.Config.InitSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	symbolCollection2.onInit(&symbolCollection2.Config.BasicComponentConfig)

	return nil
}

// // OnNewGame -
// func (symbolCollection2 *SymbolCollection2) OnNewGame(gameProp *GameProperty) error {
// 	cd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

// 	cd.OnNewGame()

// 	cd.SymbolCodes = append(cd.SymbolCodes, symbolCollection2.Config.InitSymbolCodes...)

// 	return nil
// }

// playgame
func (symbolCollection2 *SymbolCollection2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolCollection2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	nc := symbolCollection2.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (symbolCollection2 *SymbolCollection2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult,
	mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*SymbolCollection2Data)

	if len(cd.SymbolCodes) > 0 {
		fmt.Printf("Symbols is %v\n", cd.SymbolCodes)
	}

	return nil
}

// // OnStats
// func (symbolCollection2 *SymbolCollection2) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// NewComponentData -
func (symbolCollection2 *SymbolCollection2) NewComponentData() IComponentData {
	return &SymbolCollection2Data{}
}

// // GetSymbols -
// func (symbolCollection2 *SymbolCollection2) GetSymbols(gameProp *GameProperty) []int {
// 	scd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

// 	return scd.SymbolCodes
// }

// // AddSymbol -
// func (symbolCollection2 *SymbolCollection2) AddSymbol(gameProp *GameProperty, symbolCode int) {
// 	scd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

// 	if symbolCollection2.Config.MaxSymbolNum <= 0 || len(scd.SymbolCodes) < symbolCollection2.Config.MaxSymbolNum {
// 		scd.SymbolCodes = append(scd.SymbolCodes, symbolCode)
// 	}
// }

func (symbolCollection2 *SymbolCollection2) runInEach(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	ccn := symbolCollection2.Config.ForeachComponent

	for {
		isComponentDoNothing := false
		curComponent := gameProp.Components.MapComponents[ccn]
		if curComponent == nil {
			break
		}

		ccd := gameProp.GetCurComponentData(curComponent)
		nc, err := curComponent.OnPlayGame(gameProp, curpr, gp, plugin, "", "", ps, stake, prs, ccd)
		if err != nil {
			if err != ErrComponentDoNothing {
				goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
					goutils.Err(err))

				return err
			}

			isComponentDoNothing = true
		}

		if !isComponentDoNothing {
			// if gAllowStats2 {
			// 	if !gameProp.stats2Cache.HasFeature(curComponent.GetName()) {
			// 		gameProp.stats2Cache.AddFeature(curComponent.GetName(), curComponent.NewStats2(gameProp.Components.statsNodeData.GetParent(curComponent.GetName())))
			// 	}

			// 	curComponent.OnStats2(ccd, gameProp.stats2Cache)
			// }

			gameProp.OnCallEnd(curComponent, ccd, gp)
		}

		ccn = nc

		if ccn == "" {
			break
		}
	}

	return nil
}

// EachSymbols - foreach symbols
func (symbolCollection2 *SymbolCollection2) EachSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState, stake *sgc7game.Stake,
	prs []*sgc7game.PlayResult, cd IComponentData) error {

	if len(symbolCollection2.Config.Children) > 0 {
		scd := cd.(*SymbolCollection2Data)

		scd.CashWin = 0
		scd.CoinWin = 0

		for i, curs := range scd.SymbolCodes {
			err := gameProp.callStack.StartEachSymbols(gameProp, symbolCollection2, symbolCollection2.Config.Children, curs, i)
			if err != nil {
				goutils.Error("SymbolCollection2.EachSymbols:StartEachSymbols",
					goutils.Err(err))

				return err
			}

			err = symbolCollection2.runInEach(gameProp, curpr, gp, plugin, ps, stake, prs)
			if err != nil {
				goutils.Error("SymbolCollection2.EachSymbols:runInEach",
					goutils.Err(err))

				return err
			}

			err = gameProp.callStack.onEachSymbolsEnd(gameProp, symbolCollection2, curs, i)
			if err != nil {
				goutils.Error("SymbolCollection2.EachSymbols:onEachSymbolsEnd",
					goutils.Err(err))

				return err
			}
		}
	}
	// 	curComponentName := symbolCollection2.Config.ForeachComponent
	// 	scd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

	// 	for i, curs := range scd.SymbolCodes {
	// 		for _, cc := range symbolCollection2.Config.Children {

	// 		}

	// 		componentNum := 0
	// 		for {
	// 			next, err := gameProp.ProcEachSymbol(curComponentName, curpr, gp, plugin, ps, stake, prs, i, curs)
	// 			if err != nil {
	// 				if err == ErrComponentDoNothing {
	// 					if next == "" {
	// 						break
	// 					}
	// 				} else {
	// 					goutils.Error("SymbolCollection2.EachSymbols:ProcEachSymbol",
	// 						goutils.Err(err))

	// 					return err
	// 				}
	// 			}

	// 			curComponentName = next

	// 			componentNum++

	// 			if componentNum > MaxComponentNumInStep {
	// 				break
	// 			}
	// 		}
	// 	}
	// }

	return nil
}

// OnGameInited - on game inited
func (symbolCollection2 *SymbolCollection2) OnGameInited(components *ComponentList) error {
	if symbolCollection2.Config.ForeachComponent != "" {
		symbolCollection2.Config.Children = components.GetAllLinkComponents(symbolCollection2.Config.ForeachComponent)
	}

	return nil
}

// IsForeach -
func (symbolCollection2 *SymbolCollection2) IsForeach() bool {
	return true
}

// GetAllLinkComponents - get all link components
func (symbolCollection2 *SymbolCollection2) GetAllLinkComponents() []string {
	return []string{symbolCollection2.Config.DefaultNextComponent, symbolCollection2.Config.ForeachComponent}
}

// GetChildLinkComponents - get next link components
func (symbolCollection2 *SymbolCollection2) GetChildLinkComponents() []string {
	return []string{symbolCollection2.Config.ForeachComponent}
}

// NewStats2 -
func (symbolCollection2 *SymbolCollection2) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptRootTrigger})
}

// OnStats2
func (symbolCollection2 *SymbolCollection2) OnStats2(icd IComponentData, s2 *stats2.Cache) {
	cd := icd.(*SymbolCollection2Data)

	s2.ProcStatsForeachTrigger(symbolCollection2.Name, len(cd.SymbolCodes), 0)
}

func NewSymbolCollection2(name string) IComponent {
	return &SymbolCollection2{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "configuration": {
// },
type jsonSymbolCollection2 struct {
}

func (jr *jsonSymbolCollection2) build() *SymbolCollection2Config {
	cfg := &SymbolCollection2Config{}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseSymbolCollection2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSymbolCollection2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolCollection2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonSymbolCollection2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolCollection2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SymbolCollection2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

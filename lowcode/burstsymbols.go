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

const BurstSymbolsTypeName = "burstSymbols"

type BurstSymbolsType int

const (
	BSTypeDiffusion BurstSymbolsType = 0
	BSTypeSurround4 BurstSymbolsType = 1
	BSTypeSurround8 BurstSymbolsType = 2
)

func parseBurstSymbolsType(str string) BurstSymbolsType {
	if str == "surround4" {
		return BSTypeSurround4
	} else if str == "surround8" {
		return BSTypeSurround8
	}

	return BSTypeDiffusion
}

type BurstSymbolsSourceType int

const (
	BSSTypeSymbols            BurstSymbolsSourceType = 0
	BSSTypePositionCollection BurstSymbolsSourceType = 1
)

func parseBurstSymbolsSourceType(str string) BurstSymbolsSourceType {
	if str == "positionCollection" {
		return BSSTypePositionCollection
	}

	return BSSTypeSymbols
}

// BurstSymbolsConfig - configuration for BurstSymbols
type BurstSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrBurstType         string                 `yaml:"burstType" json:"burstType"`
	BurstType            BurstSymbolsType       `yaml:"-" json:"-"`
	BurstNumber          int                    `yaml:"burstNumber" json:"burstNumber"`
	StrSourceType        string                 `yaml:"burstSymbolsSourceType" json:"burstSymbolsSourceType"`
	SourceType           BurstSymbolsSourceType `yaml:"-" json:"-"`
	SourceSymbols        []string               `yaml:"sourceSymbols" json:"sourceSymbols"`
	SourceSymbolCodes    []int                  `yaml:"-" json:"-"`
	ExcludeSymbols       []string               `yaml:"excludeSymbols" json:"excludeSymbols"`
	excludeSymbolCodes   []int                  `yaml:"-" json:"-"`
	TargetSymbol         string                 `yaml:"targetSymbol" json:"targetSymbol"`
	TargetSymbolCode     int                    `yaml:"-" json:"-"`
	PositionCollection   string                 `yaml:"positionCollection" json:"positionCollection"`
	Controllers          []*Award               `yaml:"controllers" json:"controllers"`         // 新的奖励系统
	JumpToComponent      string                 `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
}

// SetLinkComponent
func (cfg *BurstSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type BurstSymbols struct {
	*BasicComponent `json:"-"`
	Config          *BurstSymbolsConfig `json:"config"`
}

// Init -
func (burstSymbols *BurstSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BurstSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &BurstSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BurstSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return burstSymbols.InitEx(cfg, pool)
}

// InitEx -
func (burstSymbols *BurstSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	burstSymbols.Config = cfg.(*BurstSymbolsConfig)
	burstSymbols.Config.ComponentType = BurstSymbolsTypeName

	burstSymbols.Config.BurstType = parseBurstSymbolsType(burstSymbols.Config.StrBurstType)
	burstSymbols.Config.SourceType = parseBurstSymbolsSourceType(burstSymbols.Config.StrSourceType)

	// for _, v := range moveSymbol.Config.MoveData {
	// 	if v.Src.Type != SelectWithXY {
	// 		sc, isok := pool.DefaultPaytables.MapSymbols[v.Src.Symbol]
	// 		if !isok {
	// 			goutils.Error("ReplaceReel.InitEx:Src.Symbol",
	// 				slog.String("symbol", v.Src.Symbol),
	// 				goutils.Err(ErrInvalidSymbol))

	// 			return ErrInvalidSymbol
	// 		}

	// 		v.Src.SymbolCode = sc
	// 	} else {
	// 		v.Src.SymbolCode = -1
	// 	}

	// 	if v.Target.Type != SelectWithXY {
	// 		sc, isok := pool.DefaultPaytables.MapSymbols[v.Target.Symbol]
	// 		if !isok {
	// 			goutils.Error("ReplaceReel.InitEx:Target.Symbol",
	// 				slog.String("symbol", v.Target.Symbol),
	// 				goutils.Err(ErrInvalidSymbol))

	// 			return ErrInvalidSymbol
	// 		}

	// 		v.Target.SymbolCode = sc
	// 	} else {
	// 		v.Target.SymbolCode = -1
	// 	}

	// 	sc, isok := pool.DefaultPaytables.MapSymbols[v.TargetSymbol]
	// 	if isok {
	// 		v.TargetSymbolCode = sc
	// 	} else {
	// 		v.TargetSymbolCode = -1
	// 	}
	// }

	for _, ctrl := range burstSymbols.Config.Controllers {
		ctrl.Init()
	}

	burstSymbols.onInit(&burstSymbols.Config.BasicComponentConfig)

	return nil
}

// playgame
func (burstSymbols *BurstSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// moveSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	gs := burstSymbols.GetTargetScene3(gameProp, curpr, prs, 0)

	sc2 := gs

	// for _, v := range moveSymbol.Config.MoveData {
	// 	srcok, srcx, srcy := v.Src.Select(sc2)
	// 	if !srcok {
	// 		continue
	// 	}

	// 	targetok, targetx, targety := v.Target.Select(sc2)
	// 	if !targetok {
	// 		continue
	// 	}

	// 	symbolCode := v.TargetSymbolCode
	// 	if symbolCode == -1 {
	// 		symbolCode = gs.Arr[srcx][srcy]
	// 	}

	// 	if srcx == targetx && srcy == targety {
	// 		if v.OverrideSrc {
	// 			gs.Arr[srcx][srcy] = symbolCode
	// 		}

	// 		if v.OverrideTarget {
	// 			gs.Arr[targetx][targety] = symbolCode
	// 		}

	// 		continue
	// 	}

	// 	if sc2 == gs {
	// 		sc2 = gs.CloneEx(gameProp.PoolScene)
	// 	}

	// 	v.Move(sc2, srcx, srcy, targetx, targety, symbolCode)
	// }

	if sc2 == gs {
		nc := burstSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	burstSymbols.AddScene(gameProp, curpr, sc2, bcd)

	if len(burstSymbols.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, burstSymbols.Config.Controllers, curpr, gp)
	}

	nc := burstSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (burstSymbols *BurstSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	asciigame.OutputScene("after burstSymbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)

	return nil
}

// // OnStats
// func (moveSymbol *MoveSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // NewStats2 -
// func (moveSymbol *MoveSymbol) NewStats2(parent string) *stats2.Feature {
// 	return stats2.NewFeature(parent, nil)
// }

// // OnStats2
// func (moveSymbol *MoveSymbol) OnStats2(icd IComponentData, s2 *stats2.Cache) {
// 	s2.ProcStatsTrigger(moveSymbol.Name)
// 	// s2.PushStepTrigger(moveSymbol.Name, true)
// }

// // OnStats2Trigger
// func (moveSymbol *MoveSymbol) OnStats2Trigger(s2 *Stats2) {
// 	s2.pushTriggerStats(moveSymbol.Name, true)
// }

func NewBurstSymbols(name string) IComponent {
	return &BurstSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"burstSymbolsSourceType": "symbols",
//		"burstType": "diffusion",
//		"burstNumber": 4,
//		"targetSymbol": [
//			"MM"
//		],
//		"excludeSymbols": [
//			"MM",
//			"RW",
//			"SC",
//			"MY"
//		]
//	},
type jsonBurstSymbols struct {
	SourceType     string   `json:"burstSymbolsSourceType"`
	BurstType      string   `json:"burstType"`
	BurstNumber    int      `json:"burstNumber"`
	SourceSymbols  []string `json:"sourceSymbols"`
	TargetSymbols  []string `json:"targetSymbol"`
	ExcludeSymbols []string `json:"excludeSymbols"`
}

func (jcfg *jsonBurstSymbols) build() *BurstSymbolsConfig {
	cfg := &BurstSymbolsConfig{
		StrBurstType:   jcfg.BurstType,
		StrSourceType:  jcfg.SourceType,
		SourceSymbols:  jcfg.SourceSymbols,
		TargetSymbol:   jcfg.TargetSymbols[0],
		ExcludeSymbols: jcfg.ExcludeSymbols,
	}

	// for _, v := range jms.MoveData {
	// 	cmd := &MoveData{
	// 		Src:            v.Src,
	// 		Target:         v.Target,
	// 		MoveType:       v.MoveType,
	// 		TargetSymbol:   v.TargetSymbol,
	// 		OverrideSrc:    v.OverrideSrc == "true",
	// 		OverrideTarget: v.OverrideTarget == "true",
	// 		OverridePath:   v.OverridePath == "true",
	// 	}

	// 	if cmd.Src.X > 0 {
	// 		cmd.Src.X--
	// 	}

	// 	if cmd.Src.Y > 0 {
	// 		cmd.Src.Y--
	// 	}

	// 	if cmd.Target.X > 0 {
	// 		cmd.Target.X--
	// 	}

	// 	if cmd.Target.Y > 0 {
	// 		cmd.Target.Y--
	// 	}

	// 	cfg.MoveData = append(cfg.MoveData, cmd)
	// }

	// cfg.UseSceneV3 = true

	return cfg
}

func parseBurstSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseBurstSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseBurstSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonBurstSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseBurstSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseBurstSymbols:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: BurstSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

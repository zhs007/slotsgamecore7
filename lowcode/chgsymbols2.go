package lowcode

import (
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const ChgSymbols2TypeName = "chgSymbols2"

type ChgSymbols2SourceType int

const (
	CS2STypeAll                ChgSymbols2SourceType = 0
	CS2STypeReels              ChgSymbols2SourceType = 1
	CS2STypeMask               ChgSymbols2SourceType = 2
	CS2STypePositionCollection ChgSymbols2SourceType = 3
)

func (t ChgSymbols2SourceType) IsReelsMode() bool {
	return t == CS2STypeReels || t == CS2STypeMask
}

func parseChgSymbols2SourceType(str string) ChgSymbols2SourceType {
	if str == "reels" {
		return CS2STypeReels
	} else if str == "mask" {
		return CS2STypeMask
	} else if str == "positioncollection" {
		return CS2STypePositionCollection
	}

	return CS2STypeAll
}

type ChgSymbols2SourceSymbolType int

const (
	CS2SSTypeNone         ChgSymbols2SourceSymbolType = 0
	CS2SSTypeSymbols      ChgSymbols2SourceSymbolType = 1
	CS2SSTypeSymbolWeight ChgSymbols2SourceSymbolType = 2
)

func parseChgSymbols2SourceSymbolType(str string) ChgSymbols2SourceSymbolType {
	if str == "symbols" {
		return CS2SSTypeSymbols
	} else if str == "symbolweight" {
		return CS2SSTypeSymbolWeight
	}

	return CS2SSTypeNone
}

type ChgSymbols2Type int

const (
	CS2TypeSymbol         ChgSymbols2Type = 0
	CS2TypeSymbolWeight   ChgSymbols2Type = 1
	CS2TypeMystery        ChgSymbols2Type = 2
	CS2TypeMysteryOnReels ChgSymbols2Type = 3
	CS2TypeUpSymbol       ChgSymbols2Type = 4
)

func parseChgSymbols2Type(str string) ChgSymbols2Type {
	if str == "symbolweight" {
		return CS2TypeSymbolWeight
	} else if str == "mystery" {
		return CS2TypeMystery
	} else if str == "mysteryonreels" {
		return CS2TypeMysteryOnReels
	} else if str == "upsymbol" {
		return CS2TypeUpSymbol
	}

	return CS2TypeSymbol
}

type ChgSymbols2ExitType int

const (
	CS2ETypeNone       ChgSymbols2ExitType = 0
	CS2ETypeMaxNumber  ChgSymbols2ExitType = 1
	CS2ETypeNoSameReel ChgSymbols2ExitType = 2
)

func parseChgSymbols2ExitType(str string) ChgSymbols2ExitType {
	if str == "maxnumber" {
		return CS2ETypeMaxNumber
	} else if str == "nosamereel" {
		return CS2ETypeNoSameReel
	}

	return CS2ETypeNone
}

type ChgSymbols2Data struct {
	BasicComponentData
	cfg *ChgSymbols2Config
}

// OnNewGame -
func (chgSymbolsData *ChgSymbols2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	chgSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone
func (chgSymbolsData *ChgSymbols2Data) Clone() IComponentData {
	target := &ChgSymbols2Data{
		BasicComponentData: chgSymbolsData.CloneBasicComponentData(),
		cfg:                chgSymbolsData.cfg,
	}

	return target
}

// // BuildPBComponentData
// func (chgSymbolsData *ChgSymbolsData) BuildPBComponentData() proto.Message {
// 	return &sgc7pb.BasicComponentData{
// 		BasicComponentData: chgSymbolsData.BuildPBBasicComponentData(),
// 	}
// }

// ChgConfigIntVal -
func (chgSymbolsData *ChgSymbols2Data) ChgConfigIntVal(key string, off int) int {
	if key == CCVHeight {
		if chgSymbolsData.cfg.Height > 0 {
			chgSymbolsData.MapConfigIntVals[key] = chgSymbolsData.cfg.Height
		}
	}

	return chgSymbolsData.BasicComponentData.ChgConfigIntVal(key, off)
}

// ChgSymbols2Config - configuration for ChgSymbols2
type ChgSymbols2Config struct {
	BasicComponentConfig  `yaml:",inline" json:",inline"`
	StrSrcType            string                      `yaml:"srcType" json:"srcType"`
	SrcType               ChgSymbols2SourceType       `yaml:"-" json:"-"`
	StrSrcSymbolType      string                      `yaml:"srcSymbolType" json:"srcSymbolType"`
	SrcSymbolType         ChgSymbols2SourceSymbolType `yaml:"-" json:"-"`
	StrType               string                      `yaml:"type" json:"type"`
	Type                  ChgSymbols2Type             `yaml:"-" json:"-"`
	StrExitType           string                      `yaml:"exitType" json:"exitType"`
	ExitType              ChgSymbols2ExitType         `yaml:"-" json:"-"`
	IsAlwaysGen           bool                        `yaml:"isAlwaysGen" json:"isAlwaysGen"`
	Height                int                         `yaml:"Height" json:"Height"`
	SrcSymbols            []string                    `yaml:"srcSymbols" json:"srcSymbols"`
	SrcSymbolCodes        []int                       `yaml:"-" json:"-"`
	Weight                string                      `yaml:"weight" json:"weight"`
	WeightVW2             *sgc7game.ValWeights2       `yaml:"-" json:"-"`
	BlankSymbol           string                      `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode       int                         `yaml:"-" json:"-"`
	SrcPositionCollection []string                    `yaml:"srcPositionCollection" json:"srcPositionCollection"`
	SrcSymbolWeight       string                      `yaml:"srcSymbolWeight" json:"srcSymbolWeight"`
	SrcSymbolWeightVW2    *sgc7game.ValWeights2       `yaml:"-" json:"-"`
	Symbol                string                      `yaml:"symbol" json:"symbol"`
	SymbolCode            int                         `yaml:"-" json:"-"`

	// Symbols              []string                              `yaml:"symbols" json:"-"`
	// SymbolCodes          []int                                 `yaml:"-" json:"symbols"`
	// SourceWeight         string                                `yaml:"sourceWeight" json:"sourceWeight"`
	// SourceWeightVW2      *sgc7game.ValWeights2                 `yaml:"-" json:"-"`
	// MaxNumber            int                                   `yaml:"maxNumber" json:"maxNumber"`
	// IsAlwaysGen     bool     `yaml:"isAlwaysGen" json:"isAlwaysGen"`
	Controllers     []*Award `yaml:"controllers" json:"controllers"`
	JumpToComponent string   `yaml:"jumpToComponent" json:"jumpToComponent"`
	// StrTriggers          []string                              `yaml:"triggers" json:"-"`
	// StrWeightOnReels     map[int]string                        `yaml:"weightOnReels" json:"weightOnReels"`
	// WeightOnReels        map[int]*sgc7game.ValWeights2         `yaml:"-" json:"-"`
	// MysteryOnReelsWeight []*ChgSymbolsMysteryOnReelsWeightData `yaml:"mysteryOnReelsWeight" json:"mysteryOnReelsWeight"`
}

// SetLinkComponent
func (cfg *ChgSymbols2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type ChgSymbols2 struct {
	*BasicComponent `json:"-"`
	Config          *ChgSymbols2Config `json:"config"`
}

// Init -
func (chgSymbols *ChgSymbols2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ChgSymbols2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ChgSymbols2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ChgSymbols2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return chgSymbols.InitEx(cfg, pool)
}

// InitEx -
func (chgSymbols *ChgSymbols2) InitEx(cfg any, pool *GamePropertyPool) error {
	chgSymbols.Config = cfg.(*ChgSymbols2Config)
	chgSymbols.Config.ComponentType = ChgSymbols2TypeName

	chgSymbols.Config.SrcType = parseChgSymbols2SourceType(chgSymbols.Config.StrSrcType)
	chgSymbols.Config.SrcSymbolType = parseChgSymbols2SourceSymbolType(chgSymbols.Config.StrSrcSymbolType)
	chgSymbols.Config.Type = parseChgSymbols2Type(chgSymbols.Config.StrType)
	chgSymbols.Config.ExitType = parseChgSymbols2ExitType(chgSymbols.Config.StrExitType)

	for _, s := range chgSymbols.Config.SrcSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ChgSymbols2.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		chgSymbols.Config.SrcSymbolCodes = append(chgSymbols.Config.SrcSymbolCodes, sc)
	}

	blankSymbolCode, isok := pool.DefaultPaytables.MapSymbols[chgSymbols.Config.BlankSymbol]
	if isok {
		chgSymbols.Config.BlankSymbolCode = blankSymbolCode
	} else {
		chgSymbols.Config.BlankSymbolCode = -1
	}

	if chgSymbols.Config.SrcSymbolWeight != "" {
		vw2, err := pool.LoadIntWeights(chgSymbols.Config.SrcSymbolWeight, true)
		if err != nil {
			goutils.Error("ChgSymbols2.InitEx:LoadIntWeights",
				slog.String("SourceWeight", chgSymbols.Config.SrcSymbolWeight),
				goutils.Err(err))

			return err
		}

		chgSymbols.Config.SrcSymbolWeightVW2 = vw2
	}

	// if chgSymbols.Config.StrWeightOnReels != nil {
	// 	chgSymbols.Config.WeightOnReels = make(map[int]*sgc7game.ValWeights2)

	// 	for k, v := range chgSymbols.Config.StrWeightOnReels {
	// 		vw2, err := pool.LoadIntWeights(v, chgSymbols.Config.UseFileMapping)
	// 		if err != nil {
	// 			goutils.Error("ChgSymbols2.InitEx:LoadIntWeights",
	// 				slog.String("Weight", v),
	// 				goutils.Err(err))

	// 			return err
	// 		}

	// 		chgSymbols.Config.WeightOnReels[k] = vw2
	// 	}
	// }

	if chgSymbols.Config.Weight != "" {
		vw2, err := pool.LoadIntWeights(chgSymbols.Config.Weight, chgSymbols.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ChgSymbols2.InitEx:LoadIntWeights",
				slog.String("Weight", chgSymbols.Config.Weight),
				goutils.Err(err))

			return err
		}

		chgSymbols.Config.WeightVW2 = vw2
	}

	// for _, v := range chgSymbols.Config.MysteryOnReelsWeight {
	// 	vw2, err := pool.LoadIntWeights(v.StrWeight, chgSymbols.Config.UseFileMapping)
	// 	if err != nil {
	// 		goutils.Error("ChgSymbols2.InitEx:MysteryOnReelsWeight:LoadIntWeights",
	// 			slog.String("Weight", v.StrWeight),
	// 			goutils.Err(err))

	// 		return err
	// 	}

	// 	v.Weight = vw2
	// }

	for _, award := range chgSymbols.Config.Controllers {
		award.Init()
	}

	chgSymbols.onInit(&chgSymbols.Config.BasicComponentConfig)

	return nil
}

func (chgSymbols *ChgSymbols2) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return chgSymbols.Config.WeightVW2
}

// func (chgSymbols *ChgSymbols2) getHeight(basicCD *BasicComponentData) int {
// 	height, isok := basicCD.GetConfigIntVal(CCVHeight)
// 	if isok {
// 		return height
// 	}

// 	return chgSymbols.Config.Height
// }

// func (chgSymbols *ChgSymbols2) GetSymbolCodes(plugin sgc7plugin.IPlugin) ([]int, error) {
// 	if chgSymbols.Config.SourceWeightVW2 != nil {
// 		iv, err := chgSymbols.Config.SourceWeightVW2.RandVal(plugin)
// 		if err != nil {
// 			goutils.Error("ChgSymbols2.GetSymbolCodes:RandVal",
// 				goutils.Err(err))

// 			return nil, err
// 		}

// 		return []int{iv.Int()}, nil
// 	}

// 	return chgSymbols.Config.SymbolCodes, nil
// }

// func (chgSymbols *ChgSymbols2) procReels(gameProp *GameProperty, cd *BasicComponentData,
// 	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, height int) (*sgc7game.GameScene, error) {

// 	syms, err := chgSymbols.GetSymbolCodes(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procReels:GetSymbolCodes",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	ngs := gs
// 	curNumber := 0
// 	isNeedBreak := false

// 	for x, arr := range gs.Arr {
// 		arry := make([]int, 0, height)

// 		for y := len(arr) - 1; y >= len(arr)-height; y-- {
// 			s := arr[y]

// 			if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {
// 				arry = append(arry, y)
// 			}
// 		}

// 		if len(arry) > 0 {
// 			cursc, err := chgSymbols.rollSymbolOnReels(gameProp, plugin, cd, x)
// 			if err != nil {
// 				goutils.Error("ChgSymbols2.procReels:rollSymbolOnReels",
// 					goutils.Err(err))

// 				return nil, err
// 			}

// 			if cursc != chgSymbols.Config.BlankSymbolCode {
// 				if ngs == gs {
// 					ngs = gs.CloneEx(gameProp.PoolScene)
// 				}

// 				if len(arry) == 1 {
// 					ngs.Arr[x][arry[0]] = cursc
// 				} else {
// 					arryi, err := plugin.Random(context.Background(), len(arry))
// 					if err != nil {
// 						goutils.Error("ChgSymbols2.procReels:Random",
// 							goutils.Err(err))

// 						return nil, err
// 					}

// 					ngs.Arr[x][arry[arryi]] = cursc
// 				}

// 				curNumber++

// 				if chgSymbols.Config.MaxNumber > 0 && curNumber >= chgSymbols.Config.MaxNumber {
// 					isNeedBreak = true

// 					break
// 				}
// 			}

// 			if isNeedBreak {
// 				break
// 			}
// 		}
// 	}

// 	return ngs, nil
// }

// func (chgSymbols *ChgSymbols2) procMystery(gameProp *GameProperty, cd *BasicComponentData,
// 	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, height int) (*sgc7game.GameScene, error) {

// 	syms, err := chgSymbols.GetSymbolCodes(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procMystery:GetSymbolCodes",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	cursc, err := chgSymbols.rollSymbol(gameProp, plugin, cd)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procMystery:rollSymbol",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	ngs := gs

// 	for x, arr := range gs.Arr {
// 		for y := len(arr) - 1; y >= len(arr)-height; y-- {
// 			s := arr[y]

// 			if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {
// 				if ngs == gs {
// 					ngs = gs.CloneEx(gameProp.PoolScene)
// 				}

// 				ngs.Arr[x][y] = cursc
// 			}
// 		}
// 	}

// 	return ngs, nil
// }

// func (chgSymbols *ChgSymbols2) procMysteryOnReels(gameProp *GameProperty, _ *BasicComponentData,
// 	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, height int) (*sgc7game.GameScene, error) {

// 	syms, err := chgSymbols.GetSymbolCodes(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procMysteryOnReels:GetSymbolCodes",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	ngs := gs

// 	for _, dat := range chgSymbols.Config.MysteryOnReelsWeight {
// 		cursc, err := dat.Weight.RandVal(plugin)
// 		if err != nil {
// 			goutils.Error("ChgSymbols2.procMysteryOnReels:RandVal",
// 				goutils.Err(err))

// 			return nil, err
// 		}

// 		for x, arr := range gs.Arr {
// 			if goutils.IndexOfIntSlice(dat.Index, x, 0) >= 0 {
// 				for y := len(arr) - 1; y >= len(arr)-height; y-- {
// 					s := arr[y]

// 					if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {
// 						if ngs == gs {
// 							ngs = gs.CloneEx(gameProp.PoolScene)
// 						}

// 						ngs.Arr[x][y] = cursc.Int()
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return ngs, nil
// }

// func (chgSymbols *ChgSymbols2) procUpgradeSymbolOfCategory(gameProp *GameProperty,
// 	cd *BasicComponentData, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene,
// 	height int) (*sgc7game.GameScene, error) {

// 	syms, err := chgSymbols.GetSymbolCodes(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procUpgradeSymbolOfCategory:GetSymbolCodes",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	cursc, err := chgSymbols.rollUpgradeSymbol(gameProp, plugin, cd, syms[0])
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procUpgradeSymbolOfCategory:rollUpgradeSymbol",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	if cursc == syms[0] {
// 		return gs, nil
// 	}

// 	ngs := gs

// 	for x, arr := range gs.Arr {
// 		for y := len(arr) - 1; y >= len(arr)-height; y-- {
// 			s := arr[y]

// 			if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {
// 				if ngs == gs {
// 					ngs = gs.CloneEx(gameProp.PoolScene)
// 				}

// 				ngs.Arr[x][y] = cursc
// 			}
// 		}
// 	}

// 	return ngs, nil
// }

// func (chgSymbols *ChgSymbols2) procNormal(gameProp *GameProperty, cd *BasicComponentData,
// 	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, height int) (*sgc7game.GameScene, error) {

// 	syms, err := chgSymbols.GetSymbolCodes(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procNormal:GetSymbolCodes",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	ngs := gs
// 	curNumber := 0
// 	isNeedBreak := false

// 	for x, arr := range gs.Arr {
// 		for y := len(arr) - 1; y >= len(arr)-height; y-- {
// 			s := arr[y]

// 			if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {
// 				cursc, err := chgSymbols.rollSymbol(gameProp, plugin, cd)
// 				if err != nil {
// 					goutils.Error("ChgSymbols2.procNormal:rollSymbol",
// 						goutils.Err(err))

// 					return nil, err
// 				}

// 				if cursc != chgSymbols.Config.BlankSymbolCode {
// 					if ngs == gs {
// 						ngs = gs.CloneEx(gameProp.PoolScene)
// 					}

// 					ngs.Arr[x][y] = cursc

// 					curNumber++

// 					if chgSymbols.Config.MaxNumber > 0 && curNumber >= chgSymbols.Config.MaxNumber {
// 						isNeedBreak = true

// 						break
// 					}
// 				}
// 			}
// 		}

// 		if isNeedBreak {
// 			break
// 		}
// 	}

// 	return ngs, nil
// }

// func (chgSymbols *ChgSymbols2) procRandomWithNoTrigger(gameProp *GameProperty, cd *BasicComponentData, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult,
// 	stake *sgc7game.Stake, gs *sgc7game.GameScene, height int) (*sgc7game.GameScene, error) {

// 	syms, err := chgSymbols.GetSymbolCodes(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procRandomWithNoTrigger:GetSymbolCodes",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	posx := []int{}
// 	posy := []int{}

// 	for x, arr := range gs.Arr {
// 		for y := len(arr) - 1; y >= len(arr)-height; y-- {
// 			if goutils.IndexOfIntSlice(syms, arr[y], 0) >= 0 {
// 				posx = append(posx, x)
// 				posy = append(posy, y)
// 			}
// 		}
// 	}

// 	if len(posx) == 0 {
// 		return gs, nil
// 	}

// 	ngs := gs

// 	curNumber := 0
// 	isNeedBreak := false

// 	srcVW2 := chgSymbols.getWeight(gameProp, cd)
// 	if srcVW2 == nil {
// 		goutils.Error("ChgSymbols2.procRandomWithNoTrigger:getWeight",
// 			goutils.Err(ErrNoWeight))

// 		return nil, ErrNoWeight
// 	}

// 	for {
// 		pi := 0

// 		if len(posx) > 1 {
// 			pi1, err := plugin.Random(context.Background(), len(posx))
// 			if err != nil {
// 				goutils.Error("ChgSymbols2.procRandomWithNoTrigger:roll pos",
// 					goutils.Err(err))

// 				return nil, err
// 			}

// 			pi = pi1
// 		}

// 		x := posx[pi]
// 		y := posy[pi]

// 		s := gs.Arr[x][y]

// 		vw2 := srcVW2.Clone()

// 		for {
// 			curscv, err := vw2.RandVal(plugin)
// 			if err != nil {
// 				goutils.Error("ChgSymbols2.procRandomWithNoTrigger:RollSymbol",
// 					goutils.Err(err))

// 				return nil, err
// 			}

// 			cursc := curscv.Int()

// 			if ngs == gs {
// 				ngs = gs.CloneEx(gameProp.PoolScene)
// 			}

// 			ngs.Arr[x][y] = cursc

// 			isTrigger := false
// 			for _, trigger := range chgSymbols.Config.StrTriggers {
// 				if gameProp.CanTrigger(trigger, ngs, curpr, stake) {
// 					isTrigger = true

// 					break
// 				}
// 			}

// 			if isTrigger {
// 				if len(vw2.Vals) == 1 {

// 					ngs.Arr[x][y] = s
// 					posx = append(posx[:pi], posx[pi+1:]...)
// 					posy = append(posy[:pi], posy[pi+1:]...)

// 					break
// 				}

// 				vw2.RemoveVal(curscv)

// 				continue
// 			}

// 			curNumber++

// 			if chgSymbols.Config.MaxNumber > 0 && curNumber >= chgSymbols.Config.MaxNumber {
// 				isNeedBreak = true

// 				break
// 			}
// 		}

// 		if isNeedBreak {
// 			break
// 		}

// 	}

// 	if curNumber == 0 {
// 		return gs, nil
// 	}

// 	return ngs, nil
// }

// OnProcControllers -
func (chgSymbols *ChgSymbols2) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(chgSymbols.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, chgSymbols.Config.Controllers, curpr, gp)
	}
}

// getSrcPos
func (chgSymbols2 *ChgSymbols2) getSrcPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult,
	prs []*sgc7game.PlayResult, gs *sgc7game.GameScene) ([]int, error) {

	pos := make([]int, 0, gameProp.GetVal(GamePropWidth)*gameProp.GetVal(GamePropHeight)*2)

	if chgSymbols2.Config.SrcType == CS2STypePositionCollection {
		for _, pc := range chgSymbols2.Config.SrcPositionCollection {
			curpos := gameProp.GetComponentPos(pc)
			if len(curpos) > 0 {
				for i := range len(curpos) / 2 {
					x := curpos[i*2]
					y := curpos[i*2+1]

					if goutils.IndexOfInt2Slice(pos, x, y, 0) < 0 {
						pos = append(pos, x, y)
					}
				}
			}
		}
	} else {
		for x := 0; x < gameProp.GetVal(GamePropWidth); x++ {
			for y := 0; y < gameProp.GetVal(GamePropHeight); y++ {
				pos = append(pos, x, y)
			}
		}
	}

	if len(pos) == 0 {
		return nil, nil
	}

	if chgSymbols2.Config.SrcSymbolType == CS2SSTypeSymbols {
		npos := make([]int, 0, gameProp.GetVal(GamePropWidth)*gameProp.GetVal(GamePropHeight)*2)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			if slices.Contains(chgSymbols2.Config.SrcSymbolCodes, gs.Arr[x][y]) {
				npos = append(npos, x, y)
			}
		}

		return npos, nil
	} else if chgSymbols2.Config.SrcSymbolType == CS2SSTypeSymbolWeight {
		curs, err := chgSymbols2.Config.SrcSymbolWeightVW2.RandVal(plugin)
		if err != nil {
			goutils.Error("ChgSymbols2.getSrcPos:RandVal",
				goutils.Err(err))

			return nil, err
		}

		npos := make([]int, 0, gameProp.GetVal(GamePropWidth)*gameProp.GetVal(GamePropHeight)*2)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			if gs.Arr[x][y] == curs.Int() {
				npos = append(npos, x, y)
			}
		}

		return npos, nil
	}

	return pos, nil
}

// playgame
func (chgSymbols *ChgSymbols2) procPos(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd *ChgSymbols2Data) (string, error) {

	gs := chgSymbols.GetTargetScene3(gameProp, curpr, prs, 0)

	pos, err := chgSymbols.getSrcPos(gameProp, plugin, curpr, prs, gs)
	if err != nil {
		goutils.Error("ChgSymbols2.procPos:getSrcPos",
			goutils.Err(err))

		return "", err
	}

	if len(pos) == 0 {
		if chgSymbols.Config.IsAlwaysGen {
			if gs != nil {
				ngs := gs.CloneEx(gameProp.PoolScene)

				chgSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
			}
		}

		nc := chgSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	if chgSymbols.Config.Type == CS2TypeSymbol {
		ngs, err := chgSymbols.procSymbolWithPos(gameProp, plugin, curpr, prs, gs, pos, chgSymbols.Config.SymbolCode)
		if err != nil {
			goutils.Error("ChgSymbols2.procPos:procSymbolWithPos",
				goutils.Err(err))

			return "", err
		}

		chgSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
	} else if chgSymbols.Config.Type == CS2TypeSymbolWeight {
		ngs, err := chgSymbols.procSymbolWeightWithPos(gameProp, plugin, curpr, prs, gs, pos, cd)
		if err != nil {
			goutils.Error("ChgSymbols2.procPos:procSymbolWeightWithPos",
				goutils.Err(err))

			return "", err
		}

		chgSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
	}

	nc := chgSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// procSymbolWithPos
func (chgSymbols2 *ChgSymbols2) procSymbolWithPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult,
	prs []*sgc7game.PlayResult, gs *sgc7game.GameScene, pos []int, symbolCode int) (*sgc7game.GameScene, error) {
	ngs := gs.CloneEx(gameProp.PoolScene)

	for i := range len(pos) / 2 {
		x := pos[i*2]
		y := pos[i*2+1]

		ngs.Arr[x][y] = symbolCode
	}

	return ngs, nil
}

// procSymbolWeightWithPos
func (chgSymbols2 *ChgSymbols2) procSymbolWeightWithPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin,
	curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult, gs *sgc7game.GameScene, pos []int, cd *ChgSymbols2Data) (*sgc7game.GameScene, error) {

	vw2 := chgSymbols2.getWeight(gameProp, &cd.BasicComponentData)
	curs, err := vw2.RandVal(plugin)
	if err != nil {
		goutils.Error("ChgSymbols2.procSymbolWeightWithPos:RandVal",
			goutils.Err(err))

		return nil, err
	}

	return chgSymbols2.procSymbolWithPos(gameProp, plugin, curpr, prs, gs, pos, curs.Int())
}

// // procMysteryWithPos
// func (chgSymbols2 *ChgSymbols2) procSymbolWithPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult,
// 	prs []*sgc7game.PlayResult, gs *sgc7game.GameScene, pos []int, symbolCode int) (*sgc7game.GameScene, error) {
// 	ngs := gs.CloneEx(gameProp.PoolScene)

// 	for i := range len(pos) / 2 {
// 		x := pos[i*2]
// 		y := pos[i*2+1]

// 		ngs.Arr[x][y] = symbolCode
// 	}

// 	return ngs, nil
// }

// playgame
func (chgSymbols *ChgSymbols2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*ChgSymbols2Data)

	// if chgSymbols.Config.SrcType.IsReelsMode() {

	// } else {
	return chgSymbols.procPos(gameProp, curpr, gp, plugin, ps, stake, prs, cd)
	// }

	// cd.UsedScenes = nil
	// cd.SrcScenes = nil

	// gs := chgSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	// if gs != nil {
	// 	height := chgSymbols.getHeight(&cd.BasicComponentData)
	// 	if height <= 0 || height > gs.Height {
	// 		height = gs.Height
	// 	}

	// 	ngs := gs

	// 	if chgSymbols.Config.Type == ChgSymTypeMystery {
	// 		ngs1, err := chgSymbols.procMystery(gameProp, &cd.BasicComponentData, plugin, gs, height)
	// 		if err != nil {
	// 			goutils.Error("ChgSymbols2.OnPlayGame:procMystery",
	// 				goutils.Err(err))

	// 			return "", err
	// 		}

	// 		ngs = ngs1
	// 	} else if chgSymbols.Config.Type == ChgSymTypeReels {
	// 		ngs1, err := chgSymbols.procReels(gameProp, &cd.BasicComponentData, plugin, gs, height)
	// 		if err != nil {
	// 			goutils.Error("ChgSymbols2.OnPlayGame:procReels",
	// 				goutils.Err(err))

	// 			return "", err
	// 		}

	// 		ngs = ngs1
	// 	} else if chgSymbols.Config.Type == ChgSymTypeRandomWithNoTrigger {
	// 		ngs2, err := chgSymbols.procRandomWithNoTrigger(gameProp, &cd.BasicComponentData, plugin, curpr, stake,
	// 			gs, height)
	// 		if err != nil {
	// 			goutils.Error("ChgSymbols2.OnPlayGame:procRandomWithNoTrigger",
	// 				goutils.Err(err))

	// 			return "", err
	// 		}

	// 		ngs = ngs2
	// 	} else if chgSymbols.Config.Type == ChgSymTypeUpgradeSymbolOfCategory {
	// 		ngs1, err := chgSymbols.procUpgradeSymbolOfCategory(gameProp, &cd.BasicComponentData, plugin, gs, height)
	// 		if err != nil {
	// 			goutils.Error("ChgSymbols2.OnPlayGame:procUpgradeSymbolOfCategory",
	// 				goutils.Err(err))

	// 			return "", err
	// 		}

	// 		ngs = ngs1
	// 	} else if chgSymbols.Config.Type == ChgSymTypeMysteryOnReels {
	// 		ngs1, err := chgSymbols.procMysteryOnReels(gameProp, &cd.BasicComponentData, plugin, gs, height)
	// 		if err != nil {
	// 			goutils.Error("ChgSymbols2.OnPlayGame:procMysteryOnReels",
	// 				goutils.Err(err))

	// 			return "", err
	// 		}

	// 		ngs = ngs1
	// 	} else {
	// 		ngs3, err := chgSymbols.procNormal(gameProp, &cd.BasicComponentData, plugin, gs, height)
	// 		if err != nil {
	// 			goutils.Error("ChgSymbols2.OnPlayGame:procNormal",
	// 				goutils.Err(err))

	// 			return "", err
	// 		}

	// 		ngs = ngs3
	// 	}

	// 	if ngs == gs {
	// 		if chgSymbols.Config.IsAlwaysGen {
	// 			ngs = gs.CloneEx(gameProp.PoolScene)

	// 			chgSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

	// 			nc := chgSymbols.onStepEnd(gameProp, curpr, gp, chgSymbols.Config.JumpToComponent)

	// 			return nc, nil
	// 		}

	// 		nc := chgSymbols.onStepEnd(gameProp, curpr, gp, "")

	// 		return nc, ErrComponentDoNothing
	// 	}

	// 	chgSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

	// 	chgSymbols.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

	// 	nc := chgSymbols.onStepEnd(gameProp, curpr, gp, chgSymbols.Config.JumpToComponent)

	// 	return nc, nil
	// }

	// nc := chgSymbols.onStepEnd(gameProp, curpr, gp, "")

	// return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (chgSymbols *ChgSymbols2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*ChgSymbols2Data)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after ChgSymbols2", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// GetAllLinkComponents - get all link components
func (chgSymbols *ChgSymbols2) GetAllLinkComponents() []string {
	return []string{chgSymbols.Config.DefaultNextComponent, chgSymbols.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (chgSymbols *ChgSymbols2) GetNextLinkComponents() []string {
	return []string{chgSymbols.Config.DefaultNextComponent, chgSymbols.Config.JumpToComponent}
}

// GetBranchNum -
func (chgSymbols *ChgSymbols2) GetBranchNum() int {
	return len(chgSymbols.Config.WeightVW2.Vals)
}

// GetBranchWeights -
func (chgSymbols *ChgSymbols2) GetBranchWeights() []int {
	return chgSymbols.Config.WeightVW2.Weights
}

// // rollSymbol -
// func (chgSymbols *ChgSymbols2) rollSymbol(gameProp *GameProperty, plugin sgc7plugin.IPlugin, bcd *BasicComponentData) (int, error) {
// 	vw2 := chgSymbols.getWeight(gameProp, bcd)
// 	if vw2 == nil {
// 		goutils.Error("ChgSymbols2.rollSymbol:getWeight",
// 			goutils.Err(ErrNoWeight))

// 		return 0, ErrNoWeight
// 	}

// 	curs, err := vw2.RandVal(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.rollSymbol:RandVal",
// 			goutils.Err(err))

// 		return 0, err
// 	}

// 	return curs.Int(), nil
// }

// // rollSymbolOnReels -
// func (chgSymbols *ChgSymbols2) rollSymbolOnReels(_ *GameProperty, plugin sgc7plugin.IPlugin, _ *BasicComponentData, x int) (int, error) {
// 	if chgSymbols.Config.WeightOnReels != nil {
// 		vw2, isok := chgSymbols.Config.WeightOnReels[x]
// 		if isok {
// 			curs, err := vw2.RandVal(plugin)
// 			if err != nil {
// 				goutils.Error("ChgSymbols2.rollSymbolOnReels:RandVal",
// 					goutils.Err(err))

// 				return 0, err
// 			}

// 			return curs.Int(), nil
// 		}
// 	}

// 	return chgSymbols.Config.BlankSymbolCode, nil
// }

// // rollUpgradeSymbol -
// func (chgSymbols *ChgSymbols2) rollUpgradeSymbol(gameProp *GameProperty, plugin sgc7plugin.IPlugin, bcd *BasicComponentData, s int) (int, error) {
// 	vw2 := chgSymbols.getWeight(gameProp, bcd)
// 	if vw2 == nil {
// 		goutils.Error("ChgSymbols2.rollUpgradeSymbol:getWeight",
// 			goutils.Err(ErrNoWeight))

// 		return 0, ErrNoWeight
// 	}

// 	vals := []sgc7game.IVal{}
// 	weights := []int{}

// 	for i, v := range vw2.Vals {
// 		if v.Int() < s {
// 			vals = append(vals, v)
// 			weights = append(weights, vw2.Weights[i])
// 		}
// 	}

// 	if len(vals) == 0 {
// 		return s, nil
// 	}

// 	curVW, err := sgc7game.NewValWeights2(vals, weights)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.rollUpgradeSymbol:NewValWeights2",
// 			goutils.Err(err))

// 		return 0, err
// 	}

// 	curs, err := curVW.RandVal(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.rollUpgradeSymbol:RandVal",
// 			goutils.Err(err))

// 		return 0, err
// 	}

// 	return curs.Int(), nil
// }

// NewComponentData -
func (chgSymbols *ChgSymbols2) NewComponentData() IComponentData {
	return &ChgSymbols2Data{
		cfg: chgSymbols.Config,
	}
}

func NewChgSymbols2(name string) IComponent {
	return &ChgSymbols2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "srcType": "all",
// "srcSymbolType": "symbols",
// "type": "symbolWeight",
// "exitType": "none",
// "isClearOutput": false,
// "Height": 0,
// "srcSymbols": [
//
//	"BN"
//
// ],
// "weight": "fgtoco",
// "blankSymbol": "BN"
type jsonChgSymbols2 struct {
	StrSrcType            string   `json:"srcType"`
	StrSrcSymbolType      string   `json:"srcSymbolType"`
	StrType               string   `json:"type"`
	StrExitType           string   `json:"exitType"`
	IsAlwaysGen           bool     `json:"isAlwaysGen"`
	Height                int      `json:"Height"`
	SrcSymbols            []string `json:"srcSymbols"`
	Weight                string   `json:"weight"`
	BlankSymbol           string   `json:"blankSymbol"`
	SrcPositionCollection []string `json:"srcPositionCollection"`
	SrcSymbolWeight       string   `json:"srcSymbolWeight"`
	Symbol                string   `json:"symbol"`
}

func (jcfg *jsonChgSymbols2) build() *ChgSymbols2Config {
	cfg := &ChgSymbols2Config{
		StrSrcType:            strings.ToLower(jcfg.StrSrcType),
		StrSrcSymbolType:      strings.ToLower(jcfg.StrSrcSymbolType),
		StrType:               strings.ToLower(jcfg.StrType),
		StrExitType:           strings.ToLower(jcfg.StrExitType),
		IsAlwaysGen:           jcfg.IsAlwaysGen,
		Height:                jcfg.Height,
		Weight:                jcfg.Weight,
		BlankSymbol:           jcfg.BlankSymbol,
		SrcSymbols:            slices.Clone(jcfg.SrcSymbols),
		SrcPositionCollection: slices.Clone(jcfg.SrcPositionCollection),
		SrcSymbolWeight:       jcfg.SrcSymbolWeight,
		Symbol:                jcfg.Symbol,
	}

	return cfg
}

func parseChgSymbols2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseChgSymbols2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseChgSymbols2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonChgSymbols2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseChgSymbols2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseClusterTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ChgSymbols2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

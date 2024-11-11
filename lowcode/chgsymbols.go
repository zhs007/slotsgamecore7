package lowcode

import (
	"log/slog"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const ChgSymbolsTypeName = "chgSymbols"

type ChgSymbolsType int

const (
	ChgSymTypeNormal                  ChgSymbolsType = 0
	ChgSymTypeMystery                 ChgSymbolsType = 1
	ChgSymTypeRandomWithNoTrigger     ChgSymbolsType = 2
	ChgSymTypeUpgradeSymbolOfCategory ChgSymbolsType = 3
)

func parseChgSymbolsType(str string) ChgSymbolsType {
	if str == "mystery" {
		return ChgSymTypeMystery
	} else if str == "randomwithnotrigger" {
		return ChgSymTypeRandomWithNoTrigger
	} else if str == "upgradesymbolofcategory" {
		return ChgSymTypeUpgradeSymbolOfCategory
	}

	return ChgSymTypeNormal
}

type ChgSymbolsData struct {
	BasicComponentData
	cfg *ChgSymbolsConfig
}

// OnNewGame -
func (chgSymbolsData *ChgSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	chgSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone
func (chgSymbolsData *ChgSymbolsData) Clone() IComponentData {
	target := &ChgSymbolsData{
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
func (chgSymbolsData *ChgSymbolsData) ChgConfigIntVal(key string, off int) int {
	if key == CCVHeight {
		if chgSymbolsData.cfg.Height > 0 {
			chgSymbolsData.MapConfigIntVals[key] = chgSymbolsData.cfg.Height
		}
	}

	return chgSymbolsData.BasicComponentData.ChgConfigIntVal(key, off)
}

// ChgSymbolsConfig - configuration for ChgSymbols
type ChgSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                `yaml:"chgSymbolsType" json:"-"`
	Type                 ChgSymbolsType        `yaml:"-" json:"chgSymbolsType"`
	Symbols              []string              `yaml:"symbols" json:"-"`
	SymbolCodes          []int                 `yaml:"-" json:"symbols"`
	BlankSymbol          string                `yaml:"blankSymbol" json:"-"`
	BlankSymbolCode      int                   `yaml:"-" json:"blankSymbol"`
	SourceWeight         string                `yaml:"sourceWeight" json:"sourceWeight"`
	SourceWeightVW2      *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	Weight               string                `yaml:"weight" json:"-"`
	WeightVW2            *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	MaxNumber            int                   `yaml:"maxNumber" json:"maxNumber"`
	IsAlwaysGen          bool                  `yaml:"isAlwaysGen" json:"isAlwaysGen"`
	Controllers          []*Award              `yaml:"controllers" json:"controllers"`
	JumpToComponent      string                `yaml:"jumpToComponent" json:"jumpToComponent"`
	StrTriggers          []string              `yaml:"triggers" json:"-"`
	Height               int                   `yaml:"height" json:"height"`
}

// SetLinkComponent
func (cfg *ChgSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type ChgSymbols struct {
	*BasicComponent `json:"-"`
	Config          *ChgSymbolsConfig `json:"config"`
}

// Init -
func (chgSymbols *ChgSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ChgSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ChgSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ChgSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return chgSymbols.InitEx(cfg, pool)
}

// InitEx -
func (chgSymbols *ChgSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	chgSymbols.Config = cfg.(*ChgSymbolsConfig)
	chgSymbols.Config.ComponentType = ChgSymbolsTypeName

	chgSymbols.Config.Type = parseChgSymbolsType(chgSymbols.Config.StrType)

	for _, s := range chgSymbols.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ChgSymbols.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		chgSymbols.Config.SymbolCodes = append(chgSymbols.Config.SymbolCodes, sc)
	}

	blankSymbolCode, isok := pool.DefaultPaytables.MapSymbols[chgSymbols.Config.BlankSymbol]
	if isok {
		chgSymbols.Config.BlankSymbolCode = blankSymbolCode
	} else {
		chgSymbols.Config.BlankSymbolCode = -1
	}

	if chgSymbols.Config.SourceWeight != "" {
		vw2, err := pool.LoadIntWeights(chgSymbols.Config.SourceWeight, chgSymbols.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ChgSymbols.InitEx:LoadIntWeights",
				slog.String("SourceWeight", chgSymbols.Config.SourceWeight),
				goutils.Err(err))

			return err
		}

		chgSymbols.Config.SourceWeightVW2 = vw2
	}

	if chgSymbols.Config.Weight != "" {
		vw2, err := pool.LoadIntWeights(chgSymbols.Config.Weight, chgSymbols.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ChgSymbols.InitEx:LoadIntWeights",
				slog.String("Weight", chgSymbols.Config.Weight),
				goutils.Err(err))

			return err
		}

		chgSymbols.Config.WeightVW2 = vw2
	} else {
		goutils.Error("ChgSymbols.InitEx",
			goutils.Err(ErrNoWeight))

		return ErrNoWeight
	}

	for _, award := range chgSymbols.Config.Controllers {
		award.Init()
	}

	chgSymbols.onInit(&chgSymbols.Config.BasicComponentConfig)

	return nil
}

func (chgSymbols *ChgSymbols) GetWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, chgSymbols.Config.UseFileMapping)

		return vw2
	}

	return chgSymbols.Config.WeightVW2
}

func (chgSymbols *ChgSymbols) getHeight(basicCD *BasicComponentData) int {
	height, isok := basicCD.GetConfigIntVal(CCVHeight)
	if isok {
		return height
	}

	return chgSymbols.Config.Height
}

func (chgSymbols *ChgSymbols) GetSymbolCodes(plugin sgc7plugin.IPlugin) ([]int, error) {
	if chgSymbols.Config.SourceWeightVW2 != nil {
		iv, err := chgSymbols.Config.SourceWeightVW2.RandVal(plugin)
		if err != nil {
			goutils.Error("ChgSymbols.GetSymbolCodes:RandVal",
				goutils.Err(err))

			return nil, err
		}

		return []int{iv.Int()}, nil
	}

	return chgSymbols.Config.SymbolCodes, nil
}

func (chgSymbols *ChgSymbols) procMystery(gameProp *GameProperty, cd *BasicComponentData,
	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, height int) (*sgc7game.GameScene, error) {

	syms, err := chgSymbols.GetSymbolCodes(plugin)
	if err != nil {
		goutils.Error("ChgSymbols.procMystery:GetSymbolCodes",
			goutils.Err(err))

		return nil, err
	}

	cursc, err := chgSymbols.RollSymbol(gameProp, plugin, cd)
	if err != nil {
		goutils.Error("ChgSymbols.procMystery:RollSymbol",
			goutils.Err(err))

		return nil, err
	}

	// isRealGen := false
	ngs := gs

	for x, arr := range gs.Arr {
		for y := len(arr) - 1; y >= len(arr)-height; y-- {
			s := arr[y]

			if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[x][y] = cursc

				// isRealGen = true
			}
		}
	}

	return ngs, nil
}

func (chgSymbols *ChgSymbols) procUpgradeSymbolOfCategory(gameProp *GameProperty,
	cd *BasicComponentData, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene,
	height int) (*sgc7game.GameScene, error) {

	syms, err := chgSymbols.GetSymbolCodes(plugin)
	if err != nil {
		goutils.Error("ChgSymbols.procUpgradeSymbolOfCategory:GetSymbolCodes",
			goutils.Err(err))

		return nil, err
	}

	cursc, err := chgSymbols.RollUpgradeSymbol(gameProp, plugin, cd, syms[0])
	if err != nil {
		goutils.Error("ChgSymbols.procUpgradeSymbolOfCategory:RollSymbol",
			goutils.Err(err))

		return nil, err
	}

	if cursc == syms[0] {
		return gs, nil
	}

	// isRealGen := false
	ngs := gs

	for x, arr := range gs.Arr {
		for y := len(arr) - 1; y >= len(arr)-height; y-- {
			s := arr[y]

			if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[x][y] = cursc

				// isRealGen = true
			}
		}
	}

	return ngs, nil
}

func (chgSymbols *ChgSymbols) procNormal(gameProp *GameProperty, cd *BasicComponentData,
	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, height int) (*sgc7game.GameScene, error) {

	syms, err := chgSymbols.GetSymbolCodes(plugin)
	if err != nil {
		goutils.Error("ChgSymbols.procNormal:GetSymbolCodes",
			goutils.Err(err))

		return nil, err
	}

	ngs := gs
	// if chgSymbols.Config.IsAlwaysGen {
	// 	ngs = gs.CloneEx(gameProp.PoolScene)
	// }

	// isRealGen := false
	curNumber := 0
	isNeedBreak := false

	for x, arr := range gs.Arr {
		for y := len(arr) - 1; y >= len(arr)-height; y-- {
			s := arr[y]

			if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {
				cursc, err := chgSymbols.RollSymbol(gameProp, plugin, cd)
				if err != nil {
					goutils.Error("ChgSymbols.procNormal:RollSymbol",
						goutils.Err(err))

					return nil, err
				}

				if cursc != chgSymbols.Config.BlankSymbolCode {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					ngs.Arr[x][y] = cursc

					curNumber++

					// isRealGen = true

					if chgSymbols.Config.MaxNumber > 0 && curNumber >= chgSymbols.Config.MaxNumber {
						isNeedBreak = true

						break
					}
				}
			}
		}

		if isNeedBreak {
			break
		}
	}

	return ngs, nil
}

func (chgSymbols *ChgSymbols) procRandomWithNoTrigger(gameProp *GameProperty, cd *BasicComponentData, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult,
	stake *sgc7game.Stake, gs *sgc7game.GameScene, height int) (*sgc7game.GameScene, error) {

	syms, err := chgSymbols.GetSymbolCodes(plugin)
	if err != nil {
		goutils.Error("ChgSymbols.procRandomWithNoTrigger:GetSymbolCodes",
			goutils.Err(err))

		return nil, err
	}

	ngs := gs

	curNumber := 0
	isNeedBreak := false

	srcVW2 := chgSymbols.GetWeight(gameProp, cd)

	for x, arr := range gs.Arr {
		for y := len(arr) - 1; y >= len(arr)-height; y-- {
			s := arr[y]

			if goutils.IndexOfIntSlice(syms, s, 0) >= 0 {

				vw2 := srcVW2.Clone()

				for {
					curscv, err := vw2.RandVal(plugin)
					if err != nil {
						goutils.Error("ChgSymbols.procRandomWithNoTrigger:RollSymbol",
							goutils.Err(err))

						return nil, err
					}

					cursc := curscv.Int()

					if cursc != chgSymbols.Config.BlankSymbolCode {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						ngs.Arr[x][y] = cursc

						isTrigger := false
						for _, trigger := range chgSymbols.Config.StrTriggers {
							if gameProp.CanTrigger(trigger, ngs, curpr, stake) {
								isTrigger = true

								break
							}
						}

						if isTrigger {
							if len(vw2.Vals) == 1 {

								ngs.Arr[x][y] = s

								break
							}

							vw2.RemoveVal(curscv)

							continue
						}

						curNumber++

						if chgSymbols.Config.MaxNumber > 0 && curNumber >= chgSymbols.Config.MaxNumber {
							isNeedBreak = true

							break
						}
					}
				}

				if isNeedBreak {
					break
				}
			}
		}

		if isNeedBreak {
			break
		}
	}

	return ngs, nil
}

// OnProcControllers -
func (chgSymbols *ChgSymbols) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(chgSymbols.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, chgSymbols.Config.Controllers, curpr, gp)
	}
}

// playgame
func (chgSymbols *ChgSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*ChgSymbolsData)

	cd.UsedScenes = nil
	cd.SrcScenes = nil

	gs := chgSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs != nil {
		height := chgSymbols.getHeight(&cd.BasicComponentData)
		if height <= 0 || height > gs.Height {
			height = gs.Height
		}

		ngs := gs

		if chgSymbols.Config.Type == ChgSymTypeMystery {
			ngs1, err := chgSymbols.procMystery(gameProp, &cd.BasicComponentData, plugin, gs, height)
			if err != nil {
				goutils.Error("ChgSymbols.OnPlayGame:procMystery",
					goutils.Err(err))

				return "", err
			}

			ngs = ngs1
		} else if chgSymbols.Config.Type == ChgSymTypeRandomWithNoTrigger {
			ngs2, err := chgSymbols.procRandomWithNoTrigger(gameProp, &cd.BasicComponentData, plugin, curpr, stake,
				gs, height)
			if err != nil {
				goutils.Error("ChgSymbols.OnPlayGame:procRandomWithNoTrigger",
					goutils.Err(err))

				return "", err
			}

			ngs = ngs2
		} else if chgSymbols.Config.Type == ChgSymTypeUpgradeSymbolOfCategory {
			ngs1, err := chgSymbols.procUpgradeSymbolOfCategory(gameProp, &cd.BasicComponentData, plugin, gs, height)
			if err != nil {
				goutils.Error("ChgSymbols.OnPlayGame:procUpgradeSymbolOfCategory",
					goutils.Err(err))

				return "", err
			}

			ngs = ngs1
		} else {
			ngs3, err := chgSymbols.procNormal(gameProp, &cd.BasicComponentData, plugin, gs, height)
			if err != nil {
				goutils.Error("ChgSymbols.OnPlayGame:procNormal",
					goutils.Err(err))

				return "", err
			}

			ngs = ngs3
		}

		if ngs == gs {
			if chgSymbols.Config.IsAlwaysGen {
				ngs = gs.CloneEx(gameProp.PoolScene)

				chgSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

				nc := chgSymbols.onStepEnd(gameProp, curpr, gp, chgSymbols.Config.JumpToComponent)

				return nc, nil
			}

			nc := chgSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		chgSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

		chgSymbols.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

		nc := chgSymbols.onStepEnd(gameProp, curpr, gp, chgSymbols.Config.JumpToComponent)

		return nc, nil
	}

	nc := chgSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (chgSymbols *ChgSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*ChgSymbolsData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after ChgSymbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// GetAllLinkComponents - get all link components
func (chgSymbols *ChgSymbols) GetAllLinkComponents() []string {
	return []string{chgSymbols.Config.DefaultNextComponent, chgSymbols.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (chgSymbols *ChgSymbols) GetNextLinkComponents() []string {
	return []string{chgSymbols.Config.DefaultNextComponent, chgSymbols.Config.JumpToComponent}
}

// GetBranchNum -
func (chgSymbols *ChgSymbols) GetBranchNum() int {
	return len(chgSymbols.Config.WeightVW2.Vals)
}

// GetBranchWeights -
func (chgSymbols *ChgSymbols) GetBranchWeights() []int {
	return chgSymbols.Config.WeightVW2.Weights
}

// RollSymbol -
func (chgSymbols *ChgSymbols) RollSymbol(gameProp *GameProperty, plugin sgc7plugin.IPlugin, bcd *BasicComponentData) (int, error) {
	// if bcd.ForceBranchIndex > 0 && !gIsReleaseMode {
	// 	vw2 := chgSymbols.GetWeight(gameProp, bcd)
	// 	return vw2.Vals[bcd.ForceBranchIndex-1].Int(), nil
	// }

	vw2 := chgSymbols.GetWeight(gameProp, bcd)
	curs, err := vw2.RandVal(plugin)
	if err != nil {
		goutils.Error("ChgSymbols.RollSymbol:RandVal",
			goutils.Err(err))

		return 0, err
	}

	return curs.Int(), nil
}

// RollUpgradeSymbol -
func (chgSymbols *ChgSymbols) RollUpgradeSymbol(gameProp *GameProperty, plugin sgc7plugin.IPlugin, bcd *BasicComponentData, s int) (int, error) {
	vw2 := chgSymbols.GetWeight(gameProp, bcd)
	vals := []sgc7game.IVal{}
	weights := []int{}

	for i, v := range vw2.Vals {
		if v.Int() < s {
			vals = append(vals, v)
			weights = append(weights, vw2.Weights[i])
		}
	}

	if len(vals) == 0 {
		return s, nil
	}

	curVW, err := sgc7game.NewValWeights2(vals, weights)
	if err != nil {
		goutils.Error("ChgSymbols.RollUpgradeSymbol:NewValWeights2",
			goutils.Err(err))

		return 0, err
	}

	curs, err := curVW.RandVal(plugin)
	if err != nil {
		goutils.Error("ChgSymbols.RollUpgradeSymbol:RandVal",
			goutils.Err(err))

		return 0, err
	}

	return curs.Int(), nil
}

// NewComponentData -
func (chgSymbols *ChgSymbols) NewComponentData() IComponentData {
	return &ChgSymbolsData{
		cfg: chgSymbols.Config,
	}
}

func NewChgSymbols(name string) IComponent {
	return &ChgSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "mystery",
// "symbols": [
//
//	"MY"
//
// ],
// "weight": "mweight",
// "maxNumber": 0,
// "isAlwaysGen": true,
type jsonChgSymbols struct {
	Symbols      []string `json:"symbols"`
	BlankSymbol  string   `yaml:"blankSymbol" json:"blankSymbol"`
	Weight       string   `yaml:"weight" json:"weight"`
	SourceWeight string   `yaml:"SourceWeight" json:"SourceWeight"`
	StrType      string   `json:"type"`
	MaxNumber    int      `json:"maxNumber"`
	IsAlwaysGen  bool     `json:"isAlwaysGen"`
	StrTriggers  []string `json:"trigger"`
	Height       int      `json:"Height"`
}

func (jcfg *jsonChgSymbols) build() *ChgSymbolsConfig {
	cfg := &ChgSymbolsConfig{
		Symbols:      jcfg.Symbols,
		BlankSymbol:  jcfg.BlankSymbol,
		Weight:       jcfg.Weight,
		StrType:      strings.ToLower(jcfg.StrType),
		MaxNumber:    jcfg.MaxNumber,
		IsAlwaysGen:  jcfg.IsAlwaysGen,
		StrTriggers:  jcfg.StrTriggers,
		SourceWeight: jcfg.SourceWeight,
		Height:       jcfg.Height,
	}

	return cfg
}

func parseChgSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseChgSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseChgSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonChgSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseChgSymbols:Unmarshal",
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
		Type: ChgSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

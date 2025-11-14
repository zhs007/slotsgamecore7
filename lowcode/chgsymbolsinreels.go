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
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const ChgSymbolsInReelsTypeName = "chgSymbolsInReels"

type ChgSymbolsInReelsSourceType int

const (
	CSIRSTypeAll   ChgSymbolsInReelsSourceType = 0
	CSIRSTypeReels ChgSymbolsInReelsSourceType = 1
	CSIRSTypeMask  ChgSymbolsInReelsSourceType = 2
)

func parseChgSymbolsInReelsSourceType(str string) ChgSymbolsInReelsSourceType {
	switch str {
	case "reels":
		return CSIRSTypeReels
	case "mask":
		return CSIRSTypeMask
	}

	return CSIRSTypeAll
}

type ChgSymbolsInReelsSourceSymbolType int

const (
	CSIRSSTypeNone         ChgSymbolsInReelsSourceSymbolType = 0
	CSIRSSTypeSymbols      ChgSymbolsInReelsSourceSymbolType = 1
	CSIRSSTypeSymbolWeight ChgSymbolsInReelsSourceSymbolType = 2
)

func parseChgSymbolsInReelsSourceSymbolType(str string) ChgSymbolsInReelsSourceSymbolType {
	switch str {
	case "symbols":
		return CSIRSSTypeSymbols
	case "symbolweight":
		return CSIRSSTypeSymbolWeight
	}

	return CSIRSSTypeNone
}

type ChgSymbolsInReelsCoreType int

const (
	CSIRCTypeSymbol        ChgSymbolsInReelsCoreType = 0
	CSIRCTypeSymbolWeight  ChgSymbolsInReelsCoreType = 1
	CSIRCTypeEachPosRandom ChgSymbolsInReelsCoreType = 2
)

func parseChgSymbolsInReelsCoreType(str string) ChgSymbolsInReelsCoreType {
	switch str {
	case "symbolweight":
		return CSIRCTypeSymbolWeight
	case "eachposrandom":
		return CSIRCTypeEachPosRandom
	}

	return CSIRCTypeSymbol
}

type ChgSymbolsInReelsNumberType int

const (
	CSIRNTypeNoLimit      ChgSymbolsInReelsNumberType = 0
	CSIRNTypeNumber       ChgSymbolsInReelsNumberType = 1
	CSIRNTypeNumberWeight ChgSymbolsInReelsNumberType = 2
)

func parseChgSymbolsInReelsNumberType(str string) ChgSymbolsInReelsNumberType {
	switch str {
	case "number":
		return CSIRNTypeNumber
	case "numberweight":
		return CSIRNTypeNumberWeight
	}

	return CSIRNTypeNoLimit
}

type ChgSymbolsInReelsData struct {
	BasicComponentData
	Pos []int
	cfg *ChgSymbolsInReelsConfig
}

// OnNewGame -
func (chgSymbolsData *ChgSymbolsInReelsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	chgSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (chgSymbolsData *ChgSymbolsInReelsData) OnNewStep() {
	chgSymbolsData.UsedScenes = nil
	chgSymbolsData.Pos = nil
}

// Clone
func (chgSymbolsData *ChgSymbolsInReelsData) Clone() IComponentData {
	target := &ChgSymbolsInReelsData{
		BasicComponentData: chgSymbolsData.CloneBasicComponentData(),
		cfg:                chgSymbolsData.cfg,
	}

	return target
}

// BuildPBComponentData
func (chgSymbolsData *ChgSymbolsInReelsData) BuildPBComponentData() proto.Message {
	return &sgc7pb.ChgSymbolsInReelsData{
		BasicComponentData: chgSymbolsData.BuildPBBasicComponentData(),
	}
}

// GetPos -
func (chgSymbolsData *ChgSymbolsInReelsData) GetPos() []int {
	return chgSymbolsData.Pos
}

// HasPos -
func (chgSymbolsData *ChgSymbolsInReelsData) HasPos(x int, y int) bool {
	return goutils.IndexOfInt2Slice(chgSymbolsData.Pos, x, y, 0) >= 0
}

// AddPos -
func (chgSymbolsData *ChgSymbolsInReelsData) AddPos(x int, y int) {
	chgSymbolsData.Pos = append(chgSymbolsData.Pos, x, y)
}

// ClearPos -
func (chgSymbolsData *ChgSymbolsInReelsData) ClearPos() {
	chgSymbolsData.Pos = nil
}

// AddPosEx -
func (chgSymbolsData *ChgSymbolsInReelsData) AddPosEx(x int, y int) {
	if !chgSymbolsData.HasPos(x, y) {
		chgSymbolsData.AddPos(x, y)
	}
}

// ChgConfigIntVal -
func (chgSymbolsData *ChgSymbolsInReelsData) ChgConfigIntVal(key string, off int) int {
	if key == CCVHeight {
		if chgSymbolsData.cfg.Height > 0 {
			chgSymbolsData.MapConfigIntVals[key] = chgSymbolsData.cfg.Height
		}
	}

	return chgSymbolsData.BasicComponentData.ChgConfigIntVal(key, off)
}

// ChgSymbolsInReelsConfig - configuration for ChgSymbolsInReels
type ChgSymbolsInReelsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrSrcType           string                            `yaml:"srcType" json:"srcType"`
	SrcType              ChgSymbolsInReelsSourceType       `yaml:"-" json:"-"`
	StrSrcSymbolType     string                            `yaml:"srcSymbolType" json:"srcSymbolType"`
	SrcSymbolType        ChgSymbolsInReelsSourceSymbolType `yaml:"-" json:"-"`
	StrType              string                            `yaml:"type" json:"type"`
	Type                 ChgSymbolsInReelsCoreType         `yaml:"-" json:"-"`
	StrNumberType        string                            `yaml:"numberType" json:"numberType"`
	NumberType           ChgSymbolsInReelsNumberType       `yaml:"-" json:"-"`
	IsAlwaysGen          bool                              `yaml:"isAlwaysGen" json:"isAlwaysGen"`
	Height               int                               `yaml:"Height" json:"Height"`
	BlankSymbol          string                            `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode      int                               `yaml:"-" json:"-"`
	MapSymbol            map[int]string                    `yaml:"mapSymbol" json:"mapSymbol"`
	MapSymbolCode        map[int]int                       `yaml:"-" json:"-"`
	MapSymbolWeight      map[int]string                    `yaml:"mapWeight" json:"mapWeight"`
	MapSymbolWeightVW    map[int]*sgc7game.ValWeights2     `yaml:"-" json:"-"`
	MapSrcSymbols        map[int][]string                  `yaml:"mapSrcSymbols" json:"mapSrcSymbols"`
	MapSrcSymbolCodes    map[int][]int                     `yaml:"-" json:"-"`
	MapNumber            map[int]int                       `yaml:"mapNumber" json:"mapNumber"`
	Controllers          []*Award                          `yaml:"controllers" json:"controllers"`
	JumpToComponent      string                            `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// SetLinkComponent
func (cfg *ChgSymbolsInReelsConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type ChgSymbolsInReels struct {
	*BasicComponent `json:"-"`
	Config          *ChgSymbolsInReelsConfig `json:"config"`
}

// Init -
func (chgSymbolsInReels *ChgSymbolsInReels) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ChgSymbolsInReels.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ChgSymbolsInReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ChgSymbolsInReels.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return chgSymbolsInReels.InitEx(cfg, pool)
}

// InitEx -
func (chgSymbolsInReels *ChgSymbolsInReels) InitEx(cfg any, pool *GamePropertyPool) error {
	chgSymbolsInReels.Config = cfg.(*ChgSymbolsInReelsConfig)
	chgSymbolsInReels.Config.ComponentType = ChgSymbolsInReelsTypeName

	chgSymbolsInReels.Config.SrcType = parseChgSymbolsInReelsSourceType(chgSymbolsInReels.Config.StrSrcType)
	chgSymbolsInReels.Config.SrcSymbolType = parseChgSymbolsInReelsSourceSymbolType(chgSymbolsInReels.Config.StrSrcSymbolType)
	chgSymbolsInReels.Config.Type = parseChgSymbolsInReelsCoreType(chgSymbolsInReels.Config.StrType)
	chgSymbolsInReels.Config.NumberType = parseChgSymbolsInReelsNumberType(chgSymbolsInReels.Config.StrNumberType)

	if len(chgSymbolsInReels.Config.MapSymbol) > 0 {
		chgSymbolsInReels.Config.MapSymbolCode = make(map[int]int, len(chgSymbolsInReels.Config.MapSymbol))

		for k, v := range chgSymbolsInReels.Config.MapSymbol {
			sc, isok := pool.DefaultPaytables.MapSymbols[v]
			if !isok {
				goutils.Error("ChgSymbolsInReels.InitEx:MapSymbol",
					slog.String("symbol", v),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}

			chgSymbolsInReels.Config.MapSymbolCode[k] = sc
		}
	}

	if len(chgSymbolsInReels.Config.MapSymbolWeight) > 0 {
		chgSymbolsInReels.Config.MapSymbolWeightVW = make(map[int]*sgc7game.ValWeights2, len(chgSymbolsInReels.Config.MapSymbolWeight))

		for k, v := range chgSymbolsInReels.Config.MapSymbolWeight {
			vw2, err := pool.LoadIntWeights(v, chgSymbolsInReels.Config.UseFileMapping)
			if err != nil {
				goutils.Error("ChgSymbolsInReels.InitEx:MapSymbolWeight:LoadIntWeights",
					slog.String("Weight", v),
					goutils.Err(err))
				return err
			}

			chgSymbolsInReels.Config.MapSymbolWeightVW[k] = vw2
		}
	}

	if len(chgSymbolsInReels.Config.MapSrcSymbols) > 0 {
		chgSymbolsInReels.Config.MapSrcSymbolCodes = make(map[int][]int, len(chgSymbolsInReels.Config.MapSrcSymbols))
		for k, arr := range chgSymbolsInReels.Config.MapSrcSymbols {
			for _, v := range arr {
				symbolCode, isok := pool.DefaultPaytables.MapSymbols[v]
				if !isok {
					goutils.Error("ChgSymbolsInReels.InitEx:MapSrcSymbols",
						slog.String("symbol", v),
						goutils.Err(ErrInvalidSymbol))
					return ErrInvalidSymbol
				}

				chgSymbolsInReels.Config.MapSrcSymbolCodes[k] = append(chgSymbolsInReels.Config.MapSrcSymbolCodes[k], symbolCode)
			}
		}
	}

	blankSymbolCode, isok := pool.DefaultPaytables.MapSymbols[chgSymbolsInReels.Config.BlankSymbol]
	if isok {
		chgSymbolsInReels.Config.BlankSymbolCode = blankSymbolCode
	} else {
		chgSymbolsInReels.Config.BlankSymbolCode = -1
	}

	for _, award := range chgSymbolsInReels.Config.Controllers {
		award.Init()
	}

	chgSymbolsInReels.onInit(&chgSymbolsInReels.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (chgSymbolsInReels *ChgSymbolsInReels) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(chgSymbolsInReels.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, chgSymbolsInReels.Config.Controllers, curpr, gp)
	}
}

func (chgSymbolsInReels *ChgSymbolsInReels) isEmpty(mapPos map[int][]int) bool {
	for _, arr := range mapPos {
		if len(arr) > 0 {
			return false
		}
	}

	return true
}

func (chgSymbolsInReels *ChgSymbolsInReels) newMapPos(gameProp *GameProperty) map[int][]int {
	mapPos := make(map[int][]int)

	for x := 0; x < gameProp.GetVal(GamePropWidth); x++ {
		mapPos[x] = make([]int, 0, gameProp.GetVal(GamePropHeight))
	}

	return mapPos
}

// getSrcPos
func (chgSymbolsInReels *ChgSymbolsInReels) getSrcPos(gameProp *GameProperty, _ sgc7plugin.IPlugin,
	gs *sgc7game.GameScene) (map[int][]int, error) {

	mapPos := chgSymbolsInReels.newMapPos(gameProp)

	if chgSymbolsInReels.Config.SrcType == CSIRSTypeAll {
		for x := 0; x < gameProp.GetVal(GamePropWidth); x++ {
			for y := 0; y < gameProp.GetVal(GamePropHeight); y++ {
				mapPos[x] = append(mapPos[x], y)
			}
		}
	} else {
		goutils.Error("ChgSymbolsInReels.getSrcPos:ErrUnsupportedSourceType",
			slog.String("srcType", chgSymbolsInReels.Config.StrSrcType),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	if chgSymbolsInReels.isEmpty(mapPos) {
		return nil, nil
	}

	if chgSymbolsInReels.Config.SrcSymbolType == CSIRSSTypeSymbols {
		nMapPos := chgSymbolsInReels.newMapPos(gameProp)

		for x, arr := range mapPos {
			for _, y := range arr {
				if slices.Contains(chgSymbolsInReels.Config.MapSrcSymbolCodes[x], gs.Arr[x][y]) {
					nMapPos[x] = append(nMapPos[x], y)
				}
			}
		}

		return nMapPos, nil
	}

	goutils.Error("ChgSymbolsInReels.getSrcPos:SrcSymbolType",
		slog.String("SrcSymbolType", chgSymbolsInReels.Config.StrSrcSymbolType),
		goutils.Err(ErrInvalidComponentConfig))

	return nil, ErrInvalidComponentConfig
}

// playgame
func (chgSymbolsInReels *ChgSymbolsInReels) procPos(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	_ sgc7game.IPlayerState, _ *sgc7game.Stake, prs []*sgc7game.PlayResult, cd *ChgSymbolsInReelsData) (string, error) {

	gs := chgSymbolsInReels.GetTargetScene3(gameProp, curpr, prs, 0)

	pos, err := chgSymbolsInReels.getSrcPos(gameProp, plugin, gs)
	if err != nil {
		goutils.Error("ChgSymbolsInReels.procPos:getSrcPos",
			goutils.Err(err))

		return "", err
	}

	if len(pos) == 0 {
		if chgSymbolsInReels.Config.IsAlwaysGen {
			if gs != nil {
				ngs := gs.CloneEx(gameProp.PoolScene)

				chgSymbolsInReels.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
			}
		}

		nc := chgSymbolsInReels.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	switch chgSymbolsInReels.Config.Type {
	case CSIRCTypeEachPosRandom:
		ngs, err := chgSymbolsInReels.procEachPosRandomWithPos(gameProp, plugin, gs, pos, cd)
		if err != nil {
			goutils.Error("ChgSymbolsInReels.procPos:procEachPosRandomWithPos",
				goutils.Err(err))

			return "", err
		}

		if ngs == gs {
			if chgSymbolsInReels.Config.IsAlwaysGen {
				if gs != nil {
					ngs1 := gs.CloneEx(gameProp.PoolScene)

					chgSymbolsInReels.AddScene(gameProp, curpr, ngs1, &cd.BasicComponentData)
				}
			}

			nc := chgSymbolsInReels.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		chgSymbolsInReels.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
	default:
		goutils.Error("ChgSymbolsInReels.procPos:Type",
			slog.String("Type", chgSymbolsInReels.Config.StrType),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	nc := chgSymbolsInReels.onStepEnd(gameProp, curpr, gp, chgSymbolsInReels.Config.JumpToComponent)

	return nc, nil
}

// procEachPosRandomWithPos
func (chgSymbolsInReels *ChgSymbolsInReels) procEachPosRandomWithPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin,
	gs *sgc7game.GameScene, mapPos map[int][]int, cd *ChgSymbolsInReelsData) (*sgc7game.GameScene, error) {

	if len(chgSymbolsInReels.Config.MapSymbolWeightVW) == 0 {
		goutils.Error("ChgSymbolsInReels.procEachPosRandomWithPos:MapSymbolWeightVW",
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	ngs := gs

	for x, yarr := range mapPos {
		maxNumber := 0

		switch chgSymbolsInReels.Config.NumberType {
		case CSIRNTypeNumberWeight:
			goutils.Error("ChgSymbolsInReels.procEachPosRandomWithPos:Type",
				slog.String("NumberType", chgSymbolsInReels.Config.StrNumberType),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		case CSIRNTypeNumber:
			maxNumber = chgSymbolsInReels.Config.MapNumber[x]
		}

		if maxNumber > 0 && len(yarr) > maxNumber {
			curnum := 0
			for _, y := range yarr {
				vw2 := chgSymbolsInReels.Config.MapSymbolWeightVW[x]

				curs, err := vw2.RandVal(plugin)
				if err != nil {
					goutils.Error("ChgSymbolsInReels.procEachPosRandomWithPos:RandVal",
						goutils.Err(err))

					return nil, err
				}

				sc := curs.Int()

				if sc != chgSymbolsInReels.Config.BlankSymbolCode {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					ngs.Arr[x][y] = sc
					cd.AddPos(x, y)

					curnum++
					if curnum >= maxNumber {
						break
					}
				}
			}
		} else {
			for _, y := range yarr {

				vw2 := chgSymbolsInReels.Config.MapSymbolWeightVW[x]

				curs, err := vw2.RandVal(plugin)
				if err != nil {
					goutils.Error("ChgSymbolsInReels.procEachPosRandomWithPos:RandVal",
						goutils.Err(err))

					return nil, err
				}

				sc := curs.Int()

				if sc != chgSymbolsInReels.Config.BlankSymbolCode {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					ngs.Arr[x][y] = sc

					cd.AddPos(x, y)
				}
			}
		}
	}

	return ngs, nil
}

// // procSymbolWeightWithPos
// func (chgSymbols *ChgSymbolsInReels) procSymbolWeightWithPos(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
// 	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, mapPos map[int][]int, cd *ChgSymbolsInReelsData) (*sgc7game.GameScene, error) {

// 	vw2 := chgSymbols2.getWeight(gameProp, &cd.BasicComponentData)
// 	curs, err := vw2.RandVal(plugin)
// 	if err != nil {
// 		goutils.Error("ChgSymbols2.procSymbolWeightWithPos:RandVal",
// 			goutils.Err(err))

// 		return nil, err
// 	}

// 	sc := curs.Int()

// 	if sc != chgSymbols2.Config.BlankSymbolCode {
// 		chgSymbols2.ProcControllers(gameProp, plugin, curpr, gp, -1,
// 			gameProp.Pool.DefaultPaytables.GetStringFromInt(sc))

// 		return chgSymbols2.procSymbolWithPos(gameProp, gs, pos, sc, cd)
// 	}

// 	return gs, nil
// }

// playgame
func (chgSymbols *ChgSymbolsInReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*ChgSymbolsInReelsData)

	return chgSymbols.procPos(gameProp, curpr, gp, plugin, ps, stake, prs, cd)
}

// OnAsciiGame - outpur to asciigame
func (chgSymbols *ChgSymbolsInReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*ChgSymbolsInReelsData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after ChgSymbolsInReels", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// GetAllLinkComponents - get all link components
func (chgSymbols *ChgSymbolsInReels) GetAllLinkComponents() []string {
	return []string{chgSymbols.Config.DefaultNextComponent, chgSymbols.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (chgSymbols *ChgSymbolsInReels) GetNextLinkComponents() []string {
	return []string{chgSymbols.Config.DefaultNextComponent, chgSymbols.Config.JumpToComponent}
}

// NewComponentData -
func (chgSymbols *ChgSymbolsInReels) NewComponentData() IComponentData {
	return &ChgSymbolsInReelsData{
		cfg: chgSymbols.Config,
	}
}

func NewChgSymbolsInReels(name string) IComponent {
	return &ChgSymbolsInReels{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "srcType": "all",
// "srcSymbolType": "symbols",
// "numberType": "number",
// "type": "symbolWeight",
// "isAlwaysGen": false,
// "Height": 0,
// "mapWeight": [
//
//	[
//	    1,
//	    "bg_chgvs"
//	],
//	[
//	    2,
//	    "bg_chgvs"
//	],
//	[
//	    3,
//	    "bg_chgvs"
//	],
//	[
//	    4,
//	    "bg_chgvs"
//	],
//	[
//	    5,
//	    "bg_chgvs"
//	]
//
// ],
// "blankSymbol": "BN",
// "mapNumber": [
//
//	[
//	    1,
//	    1
//	],
//	[
//	    2,
//	    1
//	],
//	[
//	    3,
//	    1
//	],
//	[
//	    4,
//	    1
//	],
//	[
//	    5,
//	    1
//	]
//
// ],
// "mapSrcSymbols": [
//
//	[
//	    1,
//	    [
//	        "WL",
//	        "L1",
//	        "L2",
//	        "L3",
//	        "L4",
//	        "L5",
//	        "H1",
//	        "H2",
//	        "H3",
//	        "H4",
//	        "H5"
//	    ]
//	],
//	[
//	    2,
//	    [
//	        "WL",
//	        "L1",
//	        "L2",
//	        "L3",
//	        "L4",
//	        "L5",
//	        "H1",
//	        "H2",
//	        "H3",
//	        "H4",
//	        "H5"
//	    ]
//	],
//	[
//	    3,
//	    [
//	        "WL",
//	        "L1",
//	        "L2",
//	        "L3",
//	        "L4",
//	        "L5",
//	        "H1",
//	        "H2",
//	        "H3",
//	        "H4",
//	        "H5"
//	    ]
//	],
//	[
//	    4,
//	    [
//	        "WL",
//	        "L1",
//	        "L2",
//	        "L3",
//	        "L4",
//	        "L5",
//	        "H1",
//	        "H2",
//	        "H3",
//	        "H4",
//	        "H5"
//	    ]
//	],
//	[
//	    5,
//	    [
//	        "WL",
//	        "L1",
//	        "L2",
//	        "L3",
//	        "L4",
//	        "L5",
//	        "H1",
//	        "H2",
//	        "H3",
//	        "H4",
//	        "H5"
//	    ]
//	]
//
// ]
type jsonChgSymbolsInReels struct {
	StrSrcType       string          `json:"srcType"`
	StrSrcSymbolType string          `json:"srcSymbolType"`
	StrNumberType    string          `json:"numberType"`
	StrType          string          `json:"type"`
	IsAlwaysGen      bool            `json:"isAlwaysGen"`
	Height           int             `json:"Height"`
	MapNumber        [][]interface{} `json:"mapNumber"`
	MapSymbol        [][]interface{} `json:"mapSymbol"`
	MapWeight        [][]interface{} `json:"mapWeight"`
	MapSrcSymbols    [][]interface{} `json:"mapSrcSymbols"`
	BlankSymbol      string          `json:"blankSymbol"`
}

func (jcfg *jsonChgSymbolsInReels) build() *ChgSymbolsInReelsConfig {
	cfg := &ChgSymbolsInReelsConfig{
		StrSrcType:       strings.ToLower(jcfg.StrSrcType),
		StrSrcSymbolType: strings.ToLower(jcfg.StrSrcSymbolType),
		StrType:          strings.ToLower(jcfg.StrType),
		StrNumberType:    strings.ToLower(jcfg.StrNumberType),
		IsAlwaysGen:      jcfg.IsAlwaysGen,
		Height:           jcfg.Height,
		BlankSymbol:      jcfg.BlankSymbol,
	}

	if len(jcfg.MapNumber) > 0 {
		cfg.MapNumber = make(map[int]int, len(jcfg.MapNumber))

		for _, arr := range jcfg.MapNumber {
			if len(arr) != 2 {
				goutils.Error("jsonChgSymbolsInReels.build:MapNumber:arr")

				return nil
			}

			k := arr[0].(float64)
			v := arr[1].(float64)
			cfg.MapNumber[int(k)-1] = int(v) // [1,w] => [0,w)
		}
	}

	if len(jcfg.MapSymbol) > 0 {
		cfg.MapSymbol = make(map[int]string, len(jcfg.MapSymbol))

		for _, arr := range jcfg.MapSymbol {
			if len(arr) != 2 {
				goutils.Error("jsonChgSymbolsInReels.build:MapSymbol:arr")

				return nil
			}

			k := arr[0].(float64)
			v := arr[1].(string)
			cfg.MapSymbol[int(k)-1] = v // [1,w] => [0,w)
		}
	}

	if len(jcfg.MapWeight) > 0 {
		cfg.MapSymbolWeight = make(map[int]string, len(jcfg.MapWeight))

		for _, arr := range jcfg.MapWeight {
			if len(arr) != 2 {
				goutils.Error("jsonChgSymbolsInReels.build:MapWeight:arr")

				return nil
			}

			k := arr[0].(float64)
			v := arr[1].(string)
			cfg.MapSymbolWeight[int(k)-1] = v // [1,w] => [0,w)
		}
	}

	if len(jcfg.MapSrcSymbols) > 0 {
		cfg.MapSrcSymbols = make(map[int][]string, len(jcfg.MapSrcSymbols))

		for _, arr := range jcfg.MapSrcSymbols {
			if len(arr) != 2 {
				goutils.Error("jsonChgSymbolsInReels.build:MapSrcSymbols:arr")

				return nil
			}

			k := arr[0].(float64)
			arr1 := arr[1].([]interface{})
			for _, v := range arr1 {
				cfg.MapSrcSymbols[int(k)-1] = append(cfg.MapSrcSymbols[int(k)-1], v.(string))
			}
			// cfg.MapSrcSymbols[int(k)-1] = v // [1,w] => [0,w)
		}
	}

	return cfg
}

func parseChgSymbolsInReels(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseChgSymbolsInReels:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseChgSymbolsInReels:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonChgSymbolsInReels{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseChgSymbolsInReels:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseChgSymbolsInReels:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ChgSymbolsInReelsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

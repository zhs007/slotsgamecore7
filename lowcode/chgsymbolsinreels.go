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
	if str == "reels" {
		return CSIRSTypeReels
	} else if str == "mask" {
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
	if str == "symbols" {
		return CSIRSSTypeSymbols
	} else if str == "symbolweight" {
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
	if str == "symbolweight" {
		return CSIRCTypeSymbolWeight
	} else if str == "eachposrandom" {
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
	if str == "number" {
		return CSIRNTypeNumber
	} else if str == "numberweight" {
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
	MaxNumber            int                               `yaml:"maxNumber" json:"maxNumber"`
	MapSymbol            map[int]string                    `json:"mapSymbol"`
	MapSymbolCode        map[int]int                       `json:"-"`
	MapSymbolWeight      map[int]string                    `json:"mapWeight"`
	MapSymbolWeightVW    map[int]*sgc7game.ValWeights2     `json:"-"`
	MapSrcSymbols        map[int]string                    `json:"mapSrcSymbols"`
	MapSrcSymbolCodes    map[int]int                       `json:"-"`
	Controllers          []*Award                          `yaml:"controllers" json:"controllers"`
	JumpToComponent      string                            `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// SetLinkComponent
func (cfg *ChgSymbolsInReelsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
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
					goutils.Err(ErrIvalidSymbol))

				return ErrIvalidSymbol
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
		chgSymbolsInReels.Config.MapSrcSymbolCodes = make(map[int]int, len(chgSymbolsInReels.Config.MapSrcSymbols))
		for k, v := range chgSymbolsInReels.Config.MapSrcSymbols {
			symbolCode, isok := pool.DefaultPaytables.MapSymbols[v]
			if !isok {
				goutils.Error("ChgSymbolsInReels.InitEx:MapSrcSymbols",
					slog.String("symbol", v),
					goutils.Err(ErrIvalidSymbol))
				return ErrIvalidSymbol
			}
			chgSymbolsInReels.Config.MapSrcSymbolCodes[k] = symbolCode
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

// getSrcPos
func (chgSymbolsInReels *ChgSymbolsInReels) getSrcPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin,
	gs *sgc7game.GameScene) ([]int, error) {

	pos := make([]int, 0, gameProp.GetVal(GamePropWidth)*gameProp.GetVal(GamePropHeight)*2)

	if chgSymbolsInReels.Config.SrcType == CSIRSTypeAll {
		for x := 0; x < gameProp.GetVal(GamePropWidth); x++ {
			for y := 0; y < gameProp.GetVal(GamePropHeight); y++ {
				pos = append(pos, x, y)
			}
		}
	} else {
		goutils.Error("ChgSymbolsInReels.getSrcPos:ErrUnsupportedSourceType",
			slog.String("srcType", chgSymbolsInReels.Config.StrSrcType),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	if len(pos) == 0 {
		return nil, nil
	}

	if chgSymbolsInReels.Config.SrcSymbolType == CSIRSSTypeSymbols {
		npos := make([]int, 0, gameProp.GetVal(GamePropWidth)*gameProp.GetVal(GamePropHeight)*2)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			if chgSymbolsInReels.Config.MapSrcSymbolCodes[x] == gs.Arr[x][y] {
				npos = append(npos, x, y)
			}
		}

		return npos, nil
	}

	goutils.Error("ChgSymbolsInReels.getSrcPos:SrcSymbolType",
		slog.String("SrcSymbolType", chgSymbolsInReels.Config.StrSrcSymbolType),
		goutils.Err(ErrInvalidComponentConfig))

	return nil, ErrInvalidComponentConfig
}

// playgame
func (chgSymbolsInReels *ChgSymbolsInReels) procPos(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd *ChgSymbols2Data) (string, error) {

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

	if chgSymbolsInReels.Config.Type == CSIRCTypeEachPosRandom {
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
	} else {
		goutils.Error("ChgSymbolsInReels.procPos:Type",
			slog.String("Type", chgSymbolsInReels.Config.StrType),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	nc := chgSymbolsInReels.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// procEachPosRandomWithPos
func (chgSymbolsInReels *ChgSymbolsInReels) procEachPosRandomWithPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin,
	gs *sgc7game.GameScene, pos []int, cd *ChgSymbols2Data) (*sgc7game.GameScene, error) {

	if len(chgSymbolsInReels.Config.MapSymbolWeightVW) == 0 {
		goutils.Error("ChgSymbolsInReels.procEachPosRandomWithPos:MapSymbolWeightVW",
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	ngs := gs
	maxNumber := chgSymbolsInReels.Config.MaxNumber

	if chgSymbolsInReels.Config.NumberType == CSIRNTypeNumberWeight {
		goutils.Error("ChgSymbolsInReels.procEachPosRandomWithPos:Type",
			slog.String("NumberType", chgSymbolsInReels.Config.StrNumberType),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	} else if chgSymbolsInReels.Config.NumberType == CSIRNTypeNoLimit {
		maxNumber = 0
	}

	if maxNumber > 0 && len(pos) > maxNumber*2 {
		curnum := 0
		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

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
				if curnum >= chgSymbolsInReels.Config.MaxNumber {
					break
				}
			}
		}
	} else {
		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

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

	return ngs, nil
}

// playgame
func (chgSymbols *ChgSymbolsInReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*ChgSymbols2Data)

	return chgSymbols.procPos(gameProp, curpr, gp, plugin, ps, stake, prs, cd)
}

// OnAsciiGame - outpur to asciigame
func (chgSymbols *ChgSymbolsInReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*ChgSymbols2Data)

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
// "numberType": "noLimit",
// "type": "symbolWeight",
// "isAlwaysGen": false,
// "Height": 0,
// "mapSymbol": [
//
//	[
//		1,
//		"CL"
//	]
//
// ],
// "mapWeight": [
//
//	[
//		1,
//		"bgmytocl"
//	],
//	[
//		2,
//		"bgmytoco"
//	],
//	[
//		3,
//		"bgmytoco"
//	],
//	[
//		4,
//		"bgmytoco"
//	],
//	[
//		5,
//		"bgmytoco"
//	]
//
// ],
// "mapSrcSymbols": [
//
//	[
//		1,
//		"MY"
//	],
//	[
//		2,
//		"MY"
//	],
//	[
//		3,
//		"MY"
//	],
//	[
//		4,
//		"MY"
//	],
//	[
//		5,
//		"MY"
//	]
//
// ]
type jsonKV struct {
	Key   int    `json:"0"` // 使用索引作为字段名
	Value string `json:"1"`
}

type jsonChgSymbolsInReels struct {
	StrSrcType       string   `json:"srcType"`
	StrSrcSymbolType string   `json:"srcSymbolType"`
	StrNumberType    string   `json:"numberType"`
	StrType          string   `json:"type"`
	IsAlwaysGen      bool     `json:"isAlwaysGen"`
	Height           int      `json:"Height"`
	MapSymbol        []jsonKV `json:"mapSymbol"`
	MapWeight        []jsonKV `json:"mapWeight"`
	MapSrcSymbols    []jsonKV `json:"mapSrcSymbols"`
	BlankSymbol      string   `json:"blankSymbol"`
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

	if len(jcfg.MapSymbol) > 0 {
		cfg.MapSymbol = make(map[int]string, len(jcfg.MapSymbol))
		for _, kv := range jcfg.MapSymbol {
			cfg.MapSymbol[kv.Key-1] = kv.Value // [1,w] => [0,w)
		}
	}

	if len(jcfg.MapWeight) > 0 {
		cfg.MapSymbolWeight = make(map[int]string, len(jcfg.MapWeight))
		for _, kv := range jcfg.MapWeight {
			cfg.MapSymbolWeight[kv.Key-1] = kv.Value // [1,w] => [0,w)
		}
	}

	if len(jcfg.MapSrcSymbols) > 0 {
		cfg.MapSrcSymbols = make(map[int]string, len(jcfg.MapSrcSymbols))
		for _, kv := range jcfg.MapWeight {
			cfg.MapSrcSymbols[kv.Key-1] = kv.Value // [1,w] => [0,w)
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

	data := &jsonChgSymbols2{}

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

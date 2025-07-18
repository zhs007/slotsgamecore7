package lowcode

import (
	"fmt"
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

const HoldAndWinTypeName = "holdAndWin"

type HoldAndWinType int

const (
	HAWTypeNormal                  HoldAndWinType = 0 // normal
	HAWTypeCollectorAndHeightLevel HoldAndWinType = 1 // Collector And Height Level
)

func parseHoldAndWinType(str string) HoldAndWinType {
	if str == "collectorandheightlevel" {
		return HAWTypeCollectorAndHeightLevel
	}

	return HAWTypeNormal
}

type HoldAndWinData struct {
	BasicComponentData
	Pos    []int
	Height int
}

// OnNewGame -
func (holdAndWinData *HoldAndWinData) OnNewGame(gameProp *GameProperty, component IComponent) {
	holdAndWinData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (holdAndWinData *HoldAndWinData) OnNewStep() {
	holdAndWinData.UsedScenes = nil
	holdAndWinData.UsedOtherScenes = nil
	holdAndWinData.Pos = nil
}

// Clone
func (holdAndWinData *HoldAndWinData) Clone() IComponentData {
	target := &HoldAndWinData{
		BasicComponentData: holdAndWinData.CloneBasicComponentData(),
		Pos:                slices.Clone(holdAndWinData.Pos),
	}

	return target
}

// BuildPBComponentData
func (holdAndWinData *HoldAndWinData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.HoldAndWinData{
		BasicComponentData: holdAndWinData.BuildPBBasicComponentData(),
		Pos:                make([]int32, len(holdAndWinData.Pos)),
	}

	for i, v := range holdAndWinData.Pos {
		pbcd.Pos[i] = int32(v)
	}

	return pbcd
}

// GetPos -
func (holdAndWinData *HoldAndWinData) GetPos() []int {
	return holdAndWinData.Pos
}

// HasPos -
func (holdAndWinData *HoldAndWinData) HasPos(x int, y int) bool {
	return goutils.IndexOfInt2Slice(holdAndWinData.Pos, x, y, 0) >= 0
}

// AddPos -
func (holdAndWinData *HoldAndWinData) AddPos(x int, y int) {
	holdAndWinData.Pos = append(holdAndWinData.Pos, x, y)
}

// AddPosEx -
func (holdAndWinData *HoldAndWinData) AddPosEx(x int, y int) {
	if !holdAndWinData.HasPos(x, y) {
		holdAndWinData.AddPos(x, y)
	}
}

// ClearPos -
func (holdAndWinData *HoldAndWinData) ClearPos() {
	holdAndWinData.Pos = nil
}

// GetValEx -
func (holdAndWinData *HoldAndWinData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVHeight {
		return holdAndWinData.Height, true
	}

	return 0, false
}

// HoldAndWinConfig - configuration for HoldAndWin
type HoldAndWinConfig struct {
	BasicComponentConfig  `yaml:",inline" json:",inline"`
	StrType               string                        `yaml:"type" json:"type"`
	Type                  HoldAndWinType                `yaml:"-" json:"-"`
	StrWeight             string                        `yaml:"weight" json:"weight"`
	WeightVW2             *sgc7game.ValWeights2         `yaml:"-" json:"-"`
	BlankSymbol           string                        `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode       int                           `yaml:"-" json:"-"`
	DefaultCoinSymbolCode int                           `yaml:"-" json:"-"`
	IgnoreSymbols         []string                      `yaml:"ignoreSymbols" json:"ignoreSymbols"`
	IgnoreSymbolCodes     []int                         `yaml:"-" json:"-"`
	MinHeight             int                           `yaml:"minHeight" json:"minHeight"`
	MaxHeight             int                           `yaml:"maxHeight" json:"maxHeight"`
	MapCoinWeight         map[string]string             `yaml:"mapCoinWeight" json:"mapCoinWeight"`
	MapCoinWeightVW2      map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`
	JumpToComponent       string                        `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
	MapAwards             map[string][]*Award           `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
func (cfg *HoldAndWinConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type HoldAndWin struct {
	*BasicComponent `json:"-"`
	Config          *HoldAndWinConfig `json:"config"`
}

// Init -
func (holdAndWin *HoldAndWin) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("HoldAndWin.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &HoldAndWinConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("HoldAndWin.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return holdAndWin.InitEx(cfg, pool)
}

// InitEx -
func (holdAndWin *HoldAndWin) InitEx(cfg any, pool *GamePropertyPool) error {
	holdAndWin.Config = cfg.(*HoldAndWinConfig)
	holdAndWin.Config.ComponentType = HoldAndWinTypeName

	holdAndWin.Config.Type = parseHoldAndWinType(holdAndWin.Config.StrType)

	if len(holdAndWin.Config.BlankSymbol) > 0 {
		sc, isok := pool.DefaultPaytables.MapSymbols[holdAndWin.Config.BlankSymbol]
		if !isok {
			goutils.Error("HoldAndWin.InitEx:BlankSymbol",
				slog.String("symbol", holdAndWin.Config.BlankSymbol),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		holdAndWin.Config.BlankSymbolCode = sc
	}

	for _, v := range holdAndWin.Config.IgnoreSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("HoldAndWin.InitEx:IgnoreSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		holdAndWin.Config.IgnoreSymbolCodes = append(holdAndWin.Config.IgnoreSymbolCodes, sc)
	}

	if holdAndWin.Config.StrWeight != "" {
		vw2, err := pool.LoadIntWeights(holdAndWin.Config.StrWeight, true)
		if err != nil {
			goutils.Error("HoldAndWin.InitEx:LoadIntWeights",
				slog.String("Weight", holdAndWin.Config.StrWeight),
				goutils.Err(err))

			return err
		}

		holdAndWin.Config.WeightVW2 = vw2
	}

	holdAndWin.Config.DefaultCoinSymbolCode = -1
	holdAndWin.Config.MapCoinWeightVW2 = make(map[int]*sgc7game.ValWeights2)
	for k, v := range holdAndWin.Config.MapCoinWeight {
		sc, isok := pool.DefaultPaytables.MapSymbols[k]
		if !isok {
			goutils.Error("HoldAndWin.InitEx:MapCoinWeight",
				slog.String("symbol", k),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		vw2, err := pool.LoadIntWeights(v, true)
		if err != nil {
			goutils.Error("HoldAndWin.InitEx:MapCoinWeight",
				slog.String("Weight", v),
				goutils.Err(err))

			return err
		}

		holdAndWin.Config.MapCoinWeightVW2[sc] = vw2

		if holdAndWin.Config.DefaultCoinSymbolCode == -1 {
			holdAndWin.Config.DefaultCoinSymbolCode = sc
		}
	}

	for _, awards := range holdAndWin.Config.MapAwards {
		for _, award := range awards {
			award.Init()
		}
	}

	holdAndWin.onInit(&holdAndWin.Config.BasicComponentConfig)

	return nil
}

func (holdAndWin *HoldAndWin) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return holdAndWin.Config.WeightVW2
}

func (holdAndWin *HoldAndWin) getCoinWeight(gameProp *GameProperty, basicCD *BasicComponentData, s int) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVMapCoinWeight + "." + strings.ToLower(gameProp.Pool.DefaultPaytables.GetStringFromInt(s)))
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return holdAndWin.Config.MapCoinWeightVW2[s]
}

// procNormal -
func (holdAndWin *HoldAndWin) procNormal(gameProp *GameProperty, plugin sgc7plugin.IPlugin, cd *HoldAndWinData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (*sgc7game.GameScene, *sgc7game.GameScene, error) {

	ngs := gs
	nos := os

	vw2 := holdAndWin.getWeight(gameProp, &cd.BasicComponentData)

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if goutils.IndexOfIntSlice(holdAndWin.Config.IgnoreSymbolCodes, s, 0) < 0 {
				cv, err := vw2.RandVal(plugin)
				if err != nil {
					goutils.Error("HoldAndWin.procNormal:RandVal",
						goutils.Err(err))

					return nil, nil, err
				}

				if cv.Int() == holdAndWin.Config.BlankSymbolCode {
					continue
				}

				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)

					if os == nil {
						nos = gameProp.PoolScene.New2(ngs.Width, ngs.Height, 0)
					} else {
						nos = os.CloneEx(gameProp.PoolScene)
					}
				}

				ngs.Arr[x][y] = cv.Int()
				cd.AddPosEx(x, y)

				cvw2 := holdAndWin.getCoinWeight(gameProp, &cd.BasicComponentData, ngs.Arr[x][y])
				coin, err := cvw2.RandVal(plugin)
				if err != nil {
					goutils.Error("HoldAndWin.procNormal:getCoinWeight:RandVal",
						goutils.Err(err))

					return nil, nil, err
				}

				nos.Arr[x][y] = coin.Int()
			}
		}
	}

	cd.Height = ngs.Height

	return ngs, nos, nil
}

func (holdAndWin *HoldAndWin) isFull(gs *sgc7game.GameScene) bool {
	for _, arr := range gs.Arr {
		for _, s := range arr {
			_, isok := holdAndWin.Config.MapCoinWeightVW2[s]
			if !isok {
				return false
			}
		}
	}

	return true
}

func (holdAndWin *HoldAndWin) isFullCollectorAndHeightLevel(gs *sgc7game.GameScene) bool {
	for x, arr := range gs.Arr {
		for y, s := range arr {
			if x == 0 && y == 0 {
				continue
			}

			_, isok := holdAndWin.Config.MapCoinWeightVW2[s]
			if !isok {
				return false
			}
		}
	}

	return true
}

// procCollectorAndHeightLevel - return gs1, os1, gs2, os2, err
func (holdAndWin *HoldAndWin) procCollectorAndHeightLevel(gameProp *GameProperty, plugin sgc7plugin.IPlugin, cd *HoldAndWinData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) (*sgc7game.GameScene, *sgc7game.GameScene, *sgc7game.GameScene, *sgc7game.GameScene, error) {

	ngs := gs
	nos := os

	vw2 := holdAndWin.getWeight(gameProp, &cd.BasicComponentData)

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if goutils.IndexOfIntSlice(holdAndWin.Config.IgnoreSymbolCodes, s, 0) < 0 {
				cv, err := vw2.RandVal(plugin)
				if err != nil {
					goutils.Error("HoldAndWin.procNormal:RandVal",
						goutils.Err(err))

					return nil, nil, nil, nil, err
				}

				if cv.Int() == holdAndWin.Config.BlankSymbolCode {
					continue
				}

				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)

					if os == nil {
						nos = gameProp.PoolScene.New2(ngs.Width, ngs.Height, 0)
					} else {
						nos = os.CloneEx(gameProp.PoolScene)
					}
				}

				ngs.Arr[x][y] = cv.Int()
				cd.AddPosEx(x, y)

				cvw2 := holdAndWin.getCoinWeight(gameProp, &cd.BasicComponentData, ngs.Arr[x][y])
				coin, err := cvw2.RandVal(plugin)
				if err != nil {
					goutils.Error("HoldAndWin.procNormal:getCoinWeight:RandVal",
						goutils.Err(err))

					return nil, nil, nil, nil, err
				}

				nos.Arr[x][y] = coin.Int()
			}
		}
	}

	if nos != nil && nos.Height < holdAndWin.Config.MaxHeight {
		if nos.Arr[0][nos.Height-1] != 0 && nos.Arr[nos.Width-1][0] != 0 && nos.Arr[nos.Width-1][nos.Height-1] != 0 {
			co := nos.Arr[0][0] + nos.Arr[0][nos.Height-1] + nos.Arr[nos.Width-1][0] + nos.Arr[nos.Width-1][nos.Height-1]

			nos2 := gameProp.PoolScene.New(nos.Width, nos.Height+1)
			ngs2 := gameProp.PoolScene.New(ngs.Width, ngs.Height+1)

			for x, arr := range nos.Arr {
				for y, s := range arr {
					if x == 0 && y == 0 {
						continue
					}

					if x == 0 && y == nos.Height-1 {
						continue
					}

					if x == nos.Width-1 && y == nos.Height-1 {
						continue
					}

					if x == nos.Width-1 && y == 0 {
						continue
					}

					nos2.Arr[x][y+1] = s
					ngs2.Arr[x][y+1] = ngs.Arr[x][y]
				}
			}

			nos2.Arr[0][0] = co
			nos2.Arr[0][nos2.Height-1] = 0
			nos2.Arr[nos2.Width-1][0] = 0
			nos2.Arr[nos2.Width-1][nos2.Height-1] = 0

			ngs2.Arr[0][0] = holdAndWin.Config.DefaultCoinSymbolCode
			ngs2.Arr[0][ngs2.Height-1] = holdAndWin.Config.BlankSymbolCode
			ngs2.Arr[ngs2.Width-1][0] = holdAndWin.Config.BlankSymbolCode
			ngs2.Arr[ngs2.Width-1][ngs2.Height-1] = holdAndWin.Config.BlankSymbolCode

			return ngs, nos, ngs2, nos2, nil
		}
	}

	cd.Height = ngs.Height

	return ngs, nos, nil, nil, nil
}

// OnProcControllers -
func (holdAndWin *HoldAndWin) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	awards, isok := holdAndWin.Config.MapAwards[strVal]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}
}

// playgame
func (holdAndWin *HoldAndWin) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*HoldAndWinData)

	cd.OnNewStep()

	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)
	sc2 := gs

	ogs := gameProp.OtherSceneStack.GetTopSceneEx(curpr, prs)
	ogs2 := ogs

	switch holdAndWin.Config.Type {
	case HAWTypeNormal:
		ngs, nos, err := holdAndWin.procNormal(gameProp, plugin, cd, gs, ogs)
		if err != nil {
			goutils.Error("HoldAndWin.OnPlayGame:procNormal",
				goutils.Err(err))

			return "", err
		}

		sc2 = ngs
		ogs2 = nos

		if sc2 == gs {
			holdAndWin.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)

			if ogs != nil {
				holdAndWin.AddOtherScene(gameProp, curpr, ogs, &cd.BasicComponentData)
			}

			nc := holdAndWin.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		holdAndWin.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)
		holdAndWin.AddOtherScene(gameProp, curpr, ogs2, &cd.BasicComponentData)

		if holdAndWin.isFull(sc2) {
			holdAndWin.ProcControllers(gameProp, plugin, curpr, gp, -1, "<full>")
		} else {
			holdAndWin.ProcControllers(gameProp, plugin, curpr, gp, -1, "<newsymbols>")
		}

		nc := holdAndWin.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	case HAWTypeCollectorAndHeightLevel:
		ngs, nos, ngs2, nos2, err := holdAndWin.procCollectorAndHeightLevel(gameProp, plugin, cd, gs, ogs)
		if err != nil {
			goutils.Error("HoldAndWin.OnPlayGame:procCollectorAndHeightLevel",
				goutils.Err(err))

			return "", err
		}

		sc2 = ngs
		ogs2 = nos

		if sc2 == gs {
			holdAndWin.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)

			if ogs != nil {
				holdAndWin.AddOtherScene(gameProp, curpr, ogs, &cd.BasicComponentData)
			}

			nc := holdAndWin.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		curgs := sc2
		holdAndWin.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)
		if ngs2 != nil {
			holdAndWin.AddScene(gameProp, curpr, ngs2, &cd.BasicComponentData)
			curgs = ngs2
		}

		holdAndWin.AddOtherScene(gameProp, curpr, ogs2, &cd.BasicComponentData)
		if nos2 != nil {
			holdAndWin.AddOtherScene(gameProp, curpr, nos2, &cd.BasicComponentData)
		}

		if ngs2 != nil {
			holdAndWin.ProcControllers(gameProp, plugin, curpr, gp, -1, fmt.Sprintf("<height=%d>", ngs2.Height))
		}

		if holdAndWin.isFullCollectorAndHeightLevel(curgs) {
			holdAndWin.ProcControllers(gameProp, plugin, curpr, gp, -1, "<full>")
		} else {
			holdAndWin.ProcControllers(gameProp, plugin, curpr, gp, -1, "<newsymbols>")
		}

		nc := holdAndWin.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	goutils.Error("HoldAndWin.OnPlayGame:InvalidType",
		slog.String("type", holdAndWin.Config.StrType),
		goutils.Err(ErrInvalidComponentConfig))

	return "", ErrInvalidComponentConfig
}

// NewComponentData -
func (flowDownSymbols *HoldAndWin) NewComponentData() IComponentData {
	return &HoldAndWinData{}
}

// OnAsciiGame - outpur to asciigame
func (flowDownSymbols *HoldAndWin) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	msd := icd.(*HoldAndWinData)

	asciigame.OutputScene("after HoldAndWin", pr.Scenes[msd.UsedScenes[0]], mapSymbolColor)

	return nil
}

func NewHoldAndWin(name string) IComponent {
	return &HoldAndWin{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "collectorAndHeightLevel",
// "weight": "weight_holdrespin",
// "blankSymbol": "BN",
// "minHeight": 4,
// "maxHeight": 8,
// "mapCoinWeight": [
// 	{
// 		"symbol": "COIN",
// 		"value": "weight_coin"
// 	}
// ],
// "ignoreSymbols": [
// 	"COIN"
// ]

type jsonHoldAndWinCoinWeight struct {
	Symbol string `json:"symbol"`
	Value  string `json:"value"`
}

type jsonHoldAndWin struct {
	StrType       string                      `json:"type"`
	StrWeight     string                      `json:"weight"`
	BlankSymbol   string                      `json:"blankSymbol"`
	IgnoreSymbols []string                    `json:"ignoreSymbols"`
	MinHeight     int                         `json:"minHeight"`
	MaxHeight     int                         `json:"maxHeight"`
	MapCoinWeight []*jsonHoldAndWinCoinWeight `json:"mapCoinWeight"`
}

func (jcfg *jsonHoldAndWin) build() *HoldAndWinConfig {
	cfg := &HoldAndWinConfig{
		StrType:       strings.ToLower(jcfg.StrType),
		StrWeight:     jcfg.StrWeight,
		BlankSymbol:   jcfg.BlankSymbol,
		IgnoreSymbols: slices.Clone(jcfg.IgnoreSymbols),
		MinHeight:     jcfg.MinHeight,
		MaxHeight:     jcfg.MaxHeight,
		MapCoinWeight: make(map[string]string),
	}

	for _, v := range jcfg.MapCoinWeight {
		cfg.MapCoinWeight[v.Symbol] = v.Value
	}

	return cfg
}

func parseHoldAndWin(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseHoldAndWin:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseHoldAndWin:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonHoldAndWin{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseHoldAndWin:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseHoldAndWin:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: HoldAndWinTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

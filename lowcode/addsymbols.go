package lowcode

import (
	"context"
	"log/slog"
	"os"

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

const AddSymbolsTypeName = "addSymbols"

type AddSymbolsType int

const (
	AddSymbolsTypeNormal              AddSymbolsType = 0 // 普通
	AddSymbolsTypeNoSameReel          AddSymbolsType = 1 // 不能加在同一轴上
	AddSymbolsTypeNoSameReelAndIgnore AddSymbolsType = 2 // 也不能加在和ignore symbol同一轴上
)

func parseAddSymbolsType(str string) AddSymbolsType {
	if str == "noSameReel" {
		return AddSymbolsTypeNoSameReel
	} else if str == "noSameReelAndIgnore" {
		return AddSymbolsTypeNoSameReelAndIgnore
	}

	return AddSymbolsTypeNormal
}

type AddSymbolNumType int

const (
	AddSymbolNumTypeNumber            AddSymbolNumType = 0 // 数字
	AddSymbolNumTypeWeight            AddSymbolNumType = 1 // 权重表
	AddSymbolNumTypeIncUntilTriggered AddSymbolNumType = 2 // 不停的加数量，直到触发器触发
)

func parseAddSymbolNumType(str string) AddSymbolNumType {
	if str == "weight" {
		return AddSymbolNumTypeWeight
	} else if str == "incUntilTriggered" {
		return AddSymbolNumTypeIncUntilTriggered
	}

	return AddSymbolNumTypeNumber
}

type AddSymbolsData struct {
	BasicComponentData
	SymbolNum int
}

// OnNewGame -
func (addSymbolsData *AddSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	addSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
	addSymbolsData.SymbolNum = 0
}

// onNewStep -
func (addSymbolsData *AddSymbolsData) onNewStep() {
	addSymbolsData.SymbolNum = 0
	addSymbolsData.UsedScenes = nil
}

// Clone
func (addSymbolsData *AddSymbolsData) Clone() IComponentData {
	target := &AddSymbolsData{
		BasicComponentData: addSymbolsData.CloneBasicComponentData(),
		SymbolNum:          addSymbolsData.SymbolNum,
	}

	return target
}

// BuildPBComponentData
func (addSymbolsData *AddSymbolsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.AddSymbolsData{
		BasicComponentData: addSymbolsData.BuildPBBasicComponentData(),
		SymbolNum:          int32(addSymbolsData.SymbolNum),
	}

	return pbcd
}

// GetValEx -
func (addSymbolsData *AddSymbolsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSymbolNum {
		return addSymbolsData.SymbolNum, true
	}

	return 0, false
}

// AddSymbolsConfig - configuration for AddSymbols
type AddSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                `yaml:"type" json:"type"`
	Type                 AddSymbolsType        `yaml:"-" json:"-"`
	Symbol               string                `yaml:"symbol" json:"symbol"`
	SymbolCode           int                   `yaml:"-" json:"-"`
	StrSymbolNumType     string                `yaml:"symbolNumType" json:"symbolNumType"`
	SymbolNumType        AddSymbolNumType      `yaml:"-" json:"-"`
	SymbolNum            int                   `yaml:"symbolNum" json:"symbolNum"`
	SymbolNumWeight      string                `yaml:"symbolNumWeight" json:"symbolNumWeight"`
	SymbolNumWeightVW    *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	IgnoreSymbols        []string              `yaml:"ignoreSymbols" json:"ignoreSymbols"`
	IgnoreSymbolCodes    []int                 `yaml:"-" json:"-"`
	SymbolNumTrigger     string                `json:"symbolNumTrigger" json:"symbolNumTrigger"`
}

// SetLinkComponent
func (cfg *AddSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type AddSymbols struct {
	*BasicComponent `json:"-"`
	Config          *AddSymbolsConfig `json:"config"`
}

// Init -
func (addSymbols *AddSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("AddSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &AddSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("AddSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return addSymbols.InitEx(cfg, pool)
}

// InitEx -
func (addSymbols *AddSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	addSymbols.Config = cfg.(*AddSymbolsConfig)
	addSymbols.Config.ComponentType = AddSymbolsTypeName

	addSymbols.Config.Type = parseAddSymbolsType(addSymbols.Config.StrType)
	addSymbols.Config.SymbolNumType = parseAddSymbolNumType(addSymbols.Config.StrSymbolNumType)

	sc, isok := pool.DefaultPaytables.MapSymbols[addSymbols.Config.Symbol]
	if !isok {
		goutils.Error("AddSymbols.InitEx:Symbol",
			slog.String("symbol", addSymbols.Config.Symbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	addSymbols.Config.SymbolCode = sc

	if addSymbols.Config.SymbolNumWeight != "" {
		vw2, err := pool.LoadIntWeights(addSymbols.Config.SymbolNumWeight, addSymbols.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightReels.Init:LoadIntWeights",
				slog.String("SymbolNumWeight", addSymbols.Config.SymbolNumWeight),
				goutils.Err(err))

			return err
		}

		addSymbols.Config.SymbolNumWeightVW = vw2
	}

	for _, v := range addSymbols.Config.IgnoreSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("AddSymbols.InitEx:IgnoreSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		addSymbols.Config.IgnoreSymbolCodes = append(addSymbols.Config.IgnoreSymbolCodes, sc)
	}

	addSymbols.onInit(&addSymbols.Config.BasicComponentConfig)

	return nil
}

// playgame
func (addSymbols *AddSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// replaceReelWithMask.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*AddSymbolsData)

	cd.onNewStep()

	gs := addSymbols.GetTargetScene3(gameProp, curpr, prs, 0)

	if addSymbols.Config.SymbolNumType == AddSymbolNumTypeIncUntilTriggered {
		if addSymbols.Config.Type == AddSymbolsTypeNormal {
			pos := make([]int, 0, gs.Width*gs.Height*2)

			for x, arr := range gs.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(addSymbols.Config.IgnoreSymbolCodes, s, 0) < 0 {
						pos = append(pos, x, y)
					}
				}
			}

			if len(pos) <= 0 {
				nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

				return nc, ErrComponentDoNothing
			}

			ngs := gs.CloneEx(gameProp.PoolScene)
			isTrigger := false

			for i := 0; i < len(pos)/2; i++ {
				cr, err := plugin.Random(context.Background(), len(pos)/2)
				if err != nil {
					goutils.Error("AddSymbols.OnPlayGame:Random",
						goutils.Err(err))

					return "", err
				}

				ngs.Arr[pos[cr*2]][pos[cr*2+1]] = addSymbols.Config.SymbolCode
				cd.SymbolNum++

				pos = append(pos[:cr*2], pos[(cr+1)*2:]...)

				if len(pos) <= 0 {
					break
				}

				if addSymbols.canTrigger(gameProp, ngs, curpr, stake) {
					isTrigger = true

					break
				}
			}

			if isTrigger {
				addSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

				nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

				return nc, nil
			}

			nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	}

	var num int

	if addSymbols.Config.SymbolNumType == AddSymbolNumTypeNumber {
		num = addSymbols.Config.SymbolNum
	} else if addSymbols.Config.SymbolNumType == AddSymbolNumTypeWeight {
		if addSymbols.Config.SymbolNumWeightVW != nil {
			cv, err := addSymbols.Config.SymbolNumWeightVW.RandVal(plugin)
			if err != nil {
				goutils.Error("AddSymbols.OnPlayGame:SymbolNumWeightVW",
					goutils.Err(err))

				return "", err
			}

			num = cv.Int()
		}
	}

	if num <= 0 {
		nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	if addSymbols.Config.Type == AddSymbolsTypeNormal {
		pos := make([]int, 0, gs.Width*gs.Height*2)

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(addSymbols.Config.IgnoreSymbolCodes, s, 0) < 0 {
					pos = append(pos, x, y)
				}
			}
		}

		if len(pos) <= 0 {
			nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		ngs := gs.CloneEx(gameProp.PoolScene)

		for i := 0; i < num; i++ {
			cr, err := plugin.Random(context.Background(), len(pos)/2)
			if err != nil {
				goutils.Error("AddSymbols.OnPlayGame:Random",
					goutils.Err(err))

				return "", err
			}

			ngs.Arr[pos[cr*2]][pos[cr*2+1]] = addSymbols.Config.SymbolCode
			cd.SymbolNum++

			pos = append(pos[:cr*2], pos[(cr+1)*2:]...)

			if len(pos) <= 0 {
				break
			}
		}

		addSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
	} else {
		xarr := make([]int, 0, gs.Width)

		if addSymbols.Config.Type == AddSymbolsTypeNoSameReel {
			for x := range gs.Arr {
				xarr = append(xarr, x)
			}
		} else if addSymbols.Config.Type == AddSymbolsTypeNoSameReelAndIgnore {
			for x, arr := range gs.Arr {
				hasIgnore := false
				for _, s := range arr {
					if goutils.IndexOfIntSlice(addSymbols.Config.IgnoreSymbolCodes, s, 0) >= 0 {
						hasIgnore = true

						break
					}
				}

				if !hasIgnore {
					xarr = append(xarr, x)
				}
			}
		}

		if len(xarr) <= 0 {
			nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		ngs := gs.CloneEx(gameProp.PoolScene)

		if len(xarr) <= num {
			for x := range xarr {
				cy, err := plugin.Random(context.Background(), gs.Height)
				if err != nil {
					goutils.Error("AddSymbols.OnPlayGame:Random",
						goutils.Err(err))

					return "", err
				}

				ngs.Arr[x][cy] = addSymbols.Config.SymbolCode
				cd.SymbolNum++
			}
		} else {
			for i := 0; i < num; i++ {
				cxi, err := plugin.Random(context.Background(), len(xarr))
				if err != nil {
					goutils.Error("AddSymbols.OnPlayGame:Random",
						goutils.Err(err))

					return "", err
				}

				cy, err := plugin.Random(context.Background(), gs.Height)
				if err != nil {
					goutils.Error("AddSymbols.OnPlayGame:Random",
						goutils.Err(err))

					return "", err
				}

				ngs.Arr[xarr[cxi]][cy] = addSymbols.Config.SymbolCode
				cd.SymbolNum++

				if len(xarr) <= 1 || i == num-1 {
					break
				}

				xarr = append(xarr[:cxi], xarr[cxi+1:]...)
			}
		}

		addSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
	}

	nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// canTrigger
func (addSymbols *AddSymbols) canTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) bool {
	return gameProp.CanTrigger(addSymbols.Config.SymbolNumTrigger, gs, curpr, stake)
}

// OnAsciiGame - outpur to asciigame
func (addSymbols *AddSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*AddSymbolsData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after addSymbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// NewComponentData -
func (addSymbols *AddSymbols) NewComponentData() IComponentData {
	return &AddSymbolsData{}
}

func NewAddSymbols(name string) IComponent {
	return &AddSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "noSameReel",
// "symbol": "WL",
// "symbolNumType": "number",
// "symbolNum": 1,
// "ignoreSymbols": [
//
//	"WL",
//	"SC"
//
// ]
type jsonAddSymbols struct {
	Type             string   `json:"type"`
	SymbolNumType    string   `json:"symbolNumType"`
	Symbol           string   `json:"symbol"`
	SymbolNum        int      `json:"symbolNum"`
	SymbolNumWeight  string   `json:"symbolNumWeight"`
	IgnoreSymbols    []string `json:"ignoreSymbols"`
	SymbolNumTrigger string   `json:"symbolNumTrigger"`
}

func (jcfg *jsonAddSymbols) build() *AddSymbolsConfig {
	cfg := &AddSymbolsConfig{
		StrType:          jcfg.Type,
		Symbol:           jcfg.Symbol,
		StrSymbolNumType: jcfg.SymbolNumType,
		SymbolNum:        jcfg.SymbolNum,
		SymbolNumWeight:  jcfg.SymbolNumWeight,
		IgnoreSymbols:    jcfg.IgnoreSymbols,
		SymbolNumTrigger: jcfg.SymbolNumTrigger,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseAddSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseAddSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseAddSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonAddSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseAddSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: AddSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

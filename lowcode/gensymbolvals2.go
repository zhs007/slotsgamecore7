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

const GenSymbolVals2TypeName = "genSymbolVals2"

type GenSymbolVals2SrcSymbolValsType int

const (
	GSV2SSVTypeNone  GenSymbolVals2SrcSymbolValsType = 0
	GSV2SSVTypeClone GenSymbolVals2SrcSymbolValsType = 1
)

func parseGenSymbolVals2SrcSymbolValsType(strType string) GenSymbolVals2SrcSymbolValsType {
	if strType == "clone" {
		return GSV2SSVTypeClone
	}

	return GSV2SSVTypeNone
}

type GenSymbolVals2CoreType int

const (
	GSV2CTypeNone         GenSymbolVals2CoreType = 0
	GSV2CTypeNumber       GenSymbolVals2CoreType = 1
	GSV2CTypeWeight       GenSymbolVals2CoreType = 2
	GSV2CTypeAdd          GenSymbolVals2CoreType = 3
	GSV2CTypeMask         GenSymbolVals2CoreType = 4
	GSV2CTypeSymbolWeight GenSymbolVals2CoreType = 5
)

func parseGenSymbolVals2CoreType(strType string) GenSymbolVals2CoreType {
	switch strType {
	case "weight":
		return GSV2CTypeWeight
	case "add":
		return GSV2CTypeAdd
	case "mask":
		return GSV2CTypeMask
	case "number":
		return GSV2CTypeNumber
	case "symbolweight":
		return GSV2CTypeSymbolWeight
	}

	return GSV2CTypeNone
}

type GenSymbolVals2Data struct {
	BasicComponentData
	cfg *GenSymbolVals2Config
}

// OnNewGame -
func (genSymbolVals2Data *GenSymbolVals2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	genSymbolVals2Data.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (genSymbolVals2Data *GenSymbolVals2Data) OnNewStep() {
	genSymbolVals2Data.UsedOtherScenes = nil
}

// Clone
func (genSymbolVals2Data *GenSymbolVals2Data) Clone() IComponentData {
	target := &GenSymbolVals2Data{
		BasicComponentData: genSymbolVals2Data.CloneBasicComponentData(),
		cfg:                genSymbolVals2Data.cfg,
	}

	return target
}

// BuildPBComponentData
func (genSymbolVals2Data *GenSymbolVals2Data) BuildPBComponentData() proto.Message {
	return &sgc7pb.BasicComponentData{
		BasicComponentData: genSymbolVals2Data.BuildPBBasicComponentData(),
	}
}

// ChgConfigIntVal -
func (genSymbolVals2Data *GenSymbolVals2Data) ChgConfigIntVal(key string, off int) int {
	if key == CCVNumber {
		_, isok := genSymbolVals2Data.MapConfigIntVals[key]
		if !isok {
			genSymbolVals2Data.MapConfigIntVals[key] = genSymbolVals2Data.cfg.Number + off

			return genSymbolVals2Data.MapConfigIntVals[key]
		}
	}

	return genSymbolVals2Data.BasicComponentData.ChgConfigIntVal(key, off)
}

// GenSymbolVals2Config - configuration for GenSymbolVals2
type GenSymbolVals2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrSrcSymbolValsType string                          `yaml:"srcSymbolValsType" json:"srcSymbolValsType"`
	SrcSymbolValsType    GenSymbolVals2SrcSymbolValsType `yaml:"-" json:"-"`
	SrcSymbols           []string                        `yaml:"srcSymbols" json:"srcSymbols"`
	SrcSymbolCodes       []int                           `yaml:"-" json:"-"`
	SrcComponents        []string                        `yaml:"srcComponents" json:"srcComponents"`
	StrGenType           string                          `yaml:"genType" json:"genType"`
	GenType              GenSymbolVals2CoreType          `yaml:"-" json:"-"`
	Number               int                             `yaml:"number" json:"number"`
	StrWeight            string                          `yaml:"number" json:"weight"`
	WeightVW             *sgc7game.ValWeights2           `yaml:"-" json:"-"`
	DefaultVal           int                             `yaml:"defaultVal" json:"defaultVal"`
	MaxVal               int                             `yaml:"maxVal" json:"maxVal"`
	IsAlwaysGen          bool                            `yaml:"isAlwaysGen" json:"isAlwaysGen"`
	MapSymbolWeights     map[string]string               `yaml:"symbolWeights" json:"symbolWeights"`
	MapSymbolWeightsVM   map[int]*sgc7game.ValWeights2   `yaml:"-" json:"-"`
	Awards               []*Award                        `yaml:"awards" json:"awards"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *GenSymbolVals2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenSymbolVals2 struct {
	*BasicComponent `json:"-"`
	Config          *GenSymbolVals2Config `json:"config"`
}

// Init -
func (genSymbolVals2 *GenSymbolVals2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenSymbolVals2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenSymbolVals2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenSymbolVals2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genSymbolVals2.InitEx(cfg, pool)
}

// InitEx -
func (genSymbolVals2 *GenSymbolVals2) InitEx(cfg any, pool *GamePropertyPool) error {
	genSymbolVals2.Config = cfg.(*GenSymbolVals2Config)
	genSymbolVals2.Config.ComponentType = GenSymbolVals2TypeName

	genSymbolVals2.Config.SrcSymbolValsType = parseGenSymbolVals2SrcSymbolValsType(genSymbolVals2.Config.StrSrcSymbolValsType)
	genSymbolVals2.Config.GenType = parseGenSymbolVals2CoreType(genSymbolVals2.Config.StrGenType)

	if genSymbolVals2.Config.StrWeight != "" {
		vw2, err := pool.LoadIntWeights(genSymbolVals2.Config.StrWeight, true)
		if err != nil {
			goutils.Error("GenSymbolVals2.Init:LoadStrWeights",
				slog.String("Weight", genSymbolVals2.Config.StrWeight),
				goutils.Err(err))

			return err
		}

		genSymbolVals2.Config.WeightVW = vw2
	}

	for _, s := range genSymbolVals2.Config.SrcSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("GenSymbolVals2.InitEx:SrcSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		genSymbolVals2.Config.SrcSymbolCodes = append(genSymbolVals2.Config.SrcSymbolCodes, sc)
	}

	if len(genSymbolVals2.Config.MapSymbolWeights) > 0 {
		genSymbolVals2.Config.MapSymbolWeightsVM = make(map[int]*sgc7game.ValWeights2, len(genSymbolVals2.Config.MapSymbolWeights))

		for s, v := range genSymbolVals2.Config.MapSymbolWeights {
			vw2, err := pool.LoadIntWeights(v, true)
			if err != nil {
				goutils.Error("GenSymbolVals2.InitEx:LoadMapSymbolWeights",
					slog.String("symbol", s),
					slog.String("value", v),
					goutils.Err(err))

				return err
			}

			sc, isok := pool.DefaultPaytables.MapSymbols[s]
			if !isok {
				goutils.Error("GenSymbolVals2.InitEx:MapSymbolWeights",
					slog.String("symbol", s),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}

			genSymbolVals2.Config.MapSymbolWeightsVM[sc] = vw2
		}
	}

	for _, award := range genSymbolVals2.Config.Awards {
		award.Init()
	}

	genSymbolVals2.onInit(&genSymbolVals2.Config.BasicComponentConfig)

	return nil
}

// getSrcPos
func (genSymbolVals2 *GenSymbolVals2) getSrcPos(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) ([]int, error) {
	if len(genSymbolVals2.Config.SrcComponents) > 0 {
		pos := []int{}

		for _, pc := range genSymbolVals2.Config.SrcComponents {
			curpos := gameProp.GetComponentPos(pc)
			if len(curpos) > 0 {
				pos = append(pos, curpos...)
			}
		}

		return pos, nil
	}

	if len(genSymbolVals2.Config.SrcSymbolCodes) > 0 {
		pos := []int{}

		gs := genSymbolVals2.GetTargetScene3(gameProp, curpr, prs, 0)
		if gs == nil {
			goutils.Error("GenSymbolVals2.getSrcPos:GetTargetScene3",
				goutils.Err(ErrInvalidScene))

			return nil, ErrInvalidScene
		}

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if slices.Index(genSymbolVals2.Config.SrcSymbolCodes, s) >= 0 {
					pos = append(pos, x, y)
				}
			}
		}

		return pos, nil
	}

	w := gameProp.GetVal(GamePropWidth)
	h := gameProp.GetVal(GamePropHeight)
	pos := make([]int, 0, w*h*2)

	for x := range w {
		for y := range h {
			pos = append(pos, x, y)
		}
	}

	return pos, nil
}

// getSrcOtherScene
func (genSymbolVals2 *GenSymbolVals2) getSrcOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	prs []*sgc7game.PlayResult) (*sgc7game.GameScene, error) {

	if genSymbolVals2.Config.SrcSymbolValsType == GSV2SSVTypeNone {
		return nil, nil
	}

	os := genSymbolVals2.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	return os, nil
}

func (genSymbolVals2 *GenSymbolVals2) getNumber(_ *GameProperty, basicCD *BasicComponentData) int {
	number, isok := basicCD.GetConfigIntVal(CCVNumber)
	if isok {
		return number
	}

	return genSymbolVals2.Config.Number
}

// procNumber
func (genSymbolVals2 *GenSymbolVals2) procNumber(gameProp *GameProperty, os *sgc7game.GameScene, pos []int, basicCD *BasicComponentData) (*sgc7game.GameScene, error) {
	number := genSymbolVals2.getNumber(gameProp, basicCD)

	// non-clone
	if os == nil {
		if len(pos) == 0 {
			if genSymbolVals2.Config.IsAlwaysGen {
				return gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight),
					genSymbolVals2.Config.DefaultVal), nil
			}

			return nil, nil
		}

		w := gameProp.GetVal(GamePropWidth)
		h := gameProp.GetVal(GamePropHeight)

		if len(pos) == w*h*2 {
			return gameProp.PoolScene.New2(w, h, number), nil
		}

		nos := gameProp.PoolScene.New2(w, h, genSymbolVals2.Config.DefaultVal)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			nos.Arr[x][y] = number
		}

		return nos, nil
	}

	// clone
	var nos *sgc7game.GameScene

	if len(pos) == 0 {
		if !genSymbolVals2.Config.IsAlwaysGen {
			return os, nil
		}

		nos = os.CloneEx(gameProp.PoolScene)

		return nos, nil
	}

	nos = os.CloneEx(gameProp.PoolScene)

	for i := range len(pos) / 2 {
		x := pos[i*2]
		y := pos[i*2+1]

		nos.Arr[x][y] = number
	}

	return nos, nil
}

func (genSymbolVals2 *GenSymbolVals2) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return genSymbolVals2.Config.WeightVW
}

func (genSymbolVals2 *GenSymbolVals2) getSymbolWeight(_ *GameProperty, _ *BasicComponentData, symbolCode int) *sgc7game.ValWeights2 {
	// str := basicCD.GetConfigVal(CCVWeight)
	// if str != "" {
	// 	vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

	// 	return vw2
	// }
	return genSymbolVals2.Config.MapSymbolWeightsVM[symbolCode]
}

// procWeight
func (genSymbolVals2 *GenSymbolVals2) procWeight(gameProp *GameProperty, os *sgc7game.GameScene, pos []int,
	plugin sgc7plugin.IPlugin, basicCD *BasicComponentData) (*sgc7game.GameScene, error) {

	vw := genSymbolVals2.getWeight(gameProp, basicCD)
	if vw == nil {
		goutils.Error("GenSymbolVals2.procWeight:getWeight",
			goutils.Err(ErrNoWeight))

		return nil, ErrNoWeight
	}

	// non-clone
	if os == nil {
		if len(pos) == 0 {
			if genSymbolVals2.Config.IsAlwaysGen {
				return gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight),
					genSymbolVals2.Config.DefaultVal), nil
			}

			return nil, nil
		}

		nos := gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight),
			genSymbolVals2.Config.DefaultVal)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			cr, err := vw.RandVal(plugin)
			if err != nil {
				goutils.Error("GenSymbolVals2.procWeight:RandVal",
					goutils.Err(err))

				return nil, err
			}

			nos.Arr[x][y] = cr.Int()
		}

		return nos, nil
	}

	// clone
	var nos *sgc7game.GameScene

	if len(pos) == 0 {
		if !genSymbolVals2.Config.IsAlwaysGen {
			return os, nil
		}

		nos = os.CloneEx(gameProp.PoolScene)

		return nos, nil
	}

	nos = os.CloneEx(gameProp.PoolScene)

	for i := range len(pos) / 2 {
		x := pos[i*2]
		y := pos[i*2+1]

		if nos.Arr[x][y] != genSymbolVals2.Config.DefaultVal {
			continue
		}

		cr, err := vw.RandVal(plugin)
		if err != nil {
			goutils.Error("GenSymbolVals2.procWeight:RandVal",
				goutils.Err(err))

			return nil, err
		}

		nos.Arr[x][y] = cr.Int()
	}

	return nos, nil
}

// procSymbolWeight
func (genSymbolVals2 *GenSymbolVals2) procSymbolWeight(gameProp *GameProperty, gs *sgc7game.GameScene, os *sgc7game.GameScene, pos []int,
	plugin sgc7plugin.IPlugin, basicCD *BasicComponentData) (*sgc7game.GameScene, error) {

	// non-clone
	if os == nil {
		if len(pos) == 0 {
			if genSymbolVals2.Config.IsAlwaysGen {
				return gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight),
					genSymbolVals2.Config.DefaultVal), nil
			}

			return nil, nil
		}

		nos := gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight),
			genSymbolVals2.Config.DefaultVal)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			vw := genSymbolVals2.getSymbolWeight(gameProp, basicCD, gs.Arr[x][y])
			if vw == nil {
				goutils.Error("GenSymbolVals2.procWeight:getSymbolWeight",
					slog.Int("symbolCode", gs.Arr[x][y]),
					goutils.Err(ErrInvalidComponentConfig))

				return nil, ErrInvalidComponentConfig
			}

			cr, err := vw.RandVal(plugin)
			if err != nil {
				goutils.Error("GenSymbolVals2.procWeight:RandVal",
					goutils.Err(err))

				return nil, err
			}

			nos.Arr[x][y] = cr.Int()
		}

		return nos, nil
	}

	// clone
	var nos *sgc7game.GameScene

	if len(pos) == 0 {
		if !genSymbolVals2.Config.IsAlwaysGen {
			return os, nil
		}

		nos = os.CloneEx(gameProp.PoolScene)

		return nos, nil
	}

	nos = os.CloneEx(gameProp.PoolScene)

	for i := range len(pos) / 2 {
		x := pos[i*2]
		y := pos[i*2+1]

		vw := genSymbolVals2.getSymbolWeight(gameProp, basicCD, gs.Arr[x][y])
		if vw == nil {
			goutils.Error("GenSymbolVals2.procWeight:getSymbolWeight",
				slog.Int("symbolCode", gs.Arr[x][y]),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		cr, err := vw.RandVal(plugin)
		if err != nil {
			goutils.Error("GenSymbolVals2.procWeight:RandVal",
				goutils.Err(err))

			return nil, err
		}

		nos.Arr[x][y] = cr.Int()
	}

	return nos, nil
}

// procAdd
func (genSymbolVals2 *GenSymbolVals2) procAdd(gameProp *GameProperty, os *sgc7game.GameScene, pos []int) (*sgc7game.GameScene, error) {
	// non-clone
	if os == nil {
		if len(pos) == 0 {
			if genSymbolVals2.Config.IsAlwaysGen {
				return gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight),
					genSymbolVals2.Config.DefaultVal), nil
			}

			return nil, nil
		}

		w := gameProp.GetVal(GamePropWidth)
		h := gameProp.GetVal(GamePropHeight)

		nos := gameProp.PoolScene.New2(w, h, genSymbolVals2.Config.DefaultVal)

		// maxVal
		if genSymbolVals2.Config.MaxVal > genSymbolVals2.Config.DefaultVal {
			for i := range len(pos) / 2 {
				x := pos[i*2]
				y := pos[i*2+1]

				if nos.Arr[x][y] < genSymbolVals2.Config.MaxVal {
					nos.Arr[x][y]++
				}
			}
		} else {
			for i := range len(pos) / 2 {
				x := pos[i*2]
				y := pos[i*2+1]

				nos.Arr[x][y]++
			}
		}

		return nos, nil
	}

	// clone
	var nos *sgc7game.GameScene

	if len(pos) == 0 {
		if !genSymbolVals2.Config.IsAlwaysGen {
			return os, nil
		}

		nos = os.CloneEx(gameProp.PoolScene)

		return nos, nil
	}

	nos = os.CloneEx(gameProp.PoolScene)

	// maxVal
	if genSymbolVals2.Config.MaxVal > genSymbolVals2.Config.DefaultVal {
		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			if nos.Arr[x][y] < genSymbolVals2.Config.MaxVal {
				nos.Arr[x][y]++
			}
		}
	} else {
		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			nos.Arr[x][y]++
		}
	}

	return nos, nil
}

// procMask
func (genSymbolVals2 *GenSymbolVals2) procMask(gameProp *GameProperty, os *sgc7game.GameScene, pos []int) (*sgc7game.GameScene, error) {
	// non-clone
	if os == nil {
		if len(pos) == 0 {
			if genSymbolVals2.Config.IsAlwaysGen {
				return gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight),
					0), nil
			}

			return nil, nil
		}

		w := gameProp.GetVal(GamePropWidth)
		h := gameProp.GetVal(GamePropHeight)

		nos := gameProp.PoolScene.New2(w, h, 0)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			if nos.Arr[x][y] == 0 {
				nos.Arr[x][y] = 1
			}
		}

		return nos, nil
	}

	// clone
	var nos *sgc7game.GameScene

	if len(pos) == 0 {
		if !genSymbolVals2.Config.IsAlwaysGen {
			return os, nil
		}

		nos = os.CloneEx(gameProp.PoolScene)

		return nos, nil
	}

	nos = os.CloneEx(gameProp.PoolScene)

	for i := range len(pos) / 2 {
		x := pos[i*2]
		y := pos[i*2+1]

		if nos.Arr[x][y] == 0 {
			nos.Arr[x][y] = 1
		}
	}

	return nos, nil
}

// playgame
func (genSymbolVals2 *GenSymbolVals2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*GenSymbolVals2Data)

	cd.UsedOtherScenes = nil

	pos, err := genSymbolVals2.getSrcPos(gameProp, curpr, prs)
	if err != nil {
		goutils.Error("GenSymbolVals2.OnPlayGame:getSrcPos",
			goutils.Err(err))

		return "", err
	}

	os, err := genSymbolVals2.getSrcOtherScene(gameProp, curpr, prs)
	if err != nil {
		goutils.Error("GenSymbolVals2.OnPlayGame:getSrcOtherScene",
			goutils.Err(err))

		return "", err
	}

	switch genSymbolVals2.Config.GenType {
	case GSV2CTypeNumber:
		nos, err := genSymbolVals2.procNumber(gameProp, os, pos, &cd.BasicComponentData)
		if err != nil {
			goutils.Error("GenSymbolVals2.OnPlayGame:procNumber",
				goutils.Err(err))

			return "", err
		}

		if nos != os {
			genSymbolVals2.AddOtherScene(gameProp, curpr, nos, &cd.BasicComponentData)
		}
	case GSV2CTypeWeight:
		nos, err := genSymbolVals2.procWeight(gameProp, os, pos, plugin, &cd.BasicComponentData)
		if err != nil {
			goutils.Error("GenSymbolVals2.OnPlayGame:procWeight",
				goutils.Err(err))

			return "", err
		}

		if nos != os {
			genSymbolVals2.AddOtherScene(gameProp, curpr, nos, &cd.BasicComponentData)
		}
	case GSV2CTypeAdd:
		nos, err := genSymbolVals2.procAdd(gameProp, os, pos)
		if err != nil {
			goutils.Error("GenSymbolVals2.OnPlayGame:procAdd",
				goutils.Err(err))

			return "", err
		}

		if nos != os {
			genSymbolVals2.AddOtherScene(gameProp, curpr, nos, &cd.BasicComponentData)
		}
	case GSV2CTypeMask:
		nos, err := genSymbolVals2.procMask(gameProp, os, pos)
		if err != nil {
			goutils.Error("GenSymbolVals2.OnPlayGame:procMask",
				goutils.Err(err))

			return "", err
		}

		if nos != os {
			genSymbolVals2.AddOtherScene(gameProp, curpr, nos, &cd.BasicComponentData)
		}
	case GSV2CTypeSymbolWeight:
		gs := genSymbolVals2.GetTargetScene3(gameProp, curpr, prs, 0)
		nos, err := genSymbolVals2.procSymbolWeight(gameProp, gs, os, pos, plugin, &cd.BasicComponentData)
		if err != nil {
			goutils.Error("GenSymbolVals2.OnPlayGame:procSymbolWeight",
				goutils.Err(err))

			return "", err
		}

		if nos != os {
			genSymbolVals2.AddOtherScene(gameProp, curpr, nos, &cd.BasicComponentData)
		}
	}

	if len(pos) <= 0 {
		nc := genSymbolVals2.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	nc := genSymbolVals2.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (genSymbolVals2 *GenSymbolVals2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*GenSymbolVals2Data)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("GenSymbolVals2", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// NewComponentData -
func (genSymbolVals2 *GenSymbolVals2) NewComponentData() IComponentData {
	return &GenSymbolVals2Data{
		cfg: genSymbolVals2.Config,
	}
}

func NewGenSymbolVals2(name string) IComponent {
	return &GenSymbolVals2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "srcSymbolValsType": "none",
// "defaultVal": 0,
// "genType": "add",
// "maxVal": 0,
// "isAlwaysGen": true,
// "srcComponents": [
//
//	"rs-pos-wm"
//
// ]
type jsonGSV2SymbolWeight struct {
	Symbol string `json:"symbol"`
	Value  string `json:"value"`
}
type jsonGenSymbolVals2 struct {
	StrSrcSymbolValsType string                 `json:"srcSymbolValsType"`
	SrcSymbols           []string               `json:"srcSymbols"`
	SrcComponents        []string               `json:"srcComponents"`
	DefaultVal           int                    `json:"defaultVal"`
	StrGenType           string                 `json:"genType"`
	Number               int                    `json:"number"`
	StrWeight            string                 `json:"weight"`
	MaxVal               int                    `json:"maxVal"`
	IsAlwaysGen          bool                   `json:"isAlwaysGen"`
	SymbolWeights        []jsonGSV2SymbolWeight `json:"symbolWeights"` // for srcSymbols
}

func (jcfg *jsonGenSymbolVals2) build() *GenSymbolVals2Config {
	cfg := &GenSymbolVals2Config{
		StrSrcSymbolValsType: strings.ToLower(jcfg.StrSrcSymbolValsType),
		SrcComponents:        slices.Clone(jcfg.SrcComponents),
		SrcSymbols:           slices.Clone(jcfg.SrcSymbols),
		DefaultVal:           jcfg.DefaultVal,
		StrGenType:           strings.ToLower(jcfg.StrGenType),
		Number:               jcfg.Number,
		StrWeight:            jcfg.StrWeight,
		MaxVal:               jcfg.MaxVal,
		IsAlwaysGen:          jcfg.IsAlwaysGen,
	}

	cfg.MapSymbolWeights = make(map[string]string, len(jcfg.SymbolWeights))
	for _, sw := range jcfg.SymbolWeights {
		cfg.MapSymbolWeights[sw.Symbol] = sw.Value
	}

	return cfg
}

func parseGenSymbolVals2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenSymbolVals2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenSymbolVals2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenSymbolVals2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenSymbolVals2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseGenSymbolVals2:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenSymbolVals2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

package lowcode

import (
	"context"
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

const ChgSymbols2TypeName = "chgSymbols2"

type ChgSymbols2SourceType int

const (
	CS2STypeAll                ChgSymbols2SourceType = 0
	CS2STypeReels              ChgSymbols2SourceType = 1
	CS2STypeMask               ChgSymbols2SourceType = 2
	CS2STypePositionCollection ChgSymbols2SourceType = 3
	CS2STypeRowMask            ChgSymbols2SourceType = 4
	CS2STypeAllMask            ChgSymbols2SourceType = 5
)

func (t ChgSymbols2SourceType) IsReelsMode() bool {
	return t == CS2STypeReels || t == CS2STypeMask
}

func parseChgSymbols2SourceType(str string) ChgSymbols2SourceType {
	switch str {
	case "reels":
		return CS2STypeReels
	case "mask":
		return CS2STypeMask
	case "positioncollection":
		return CS2STypePositionCollection
	case "rowmask":
		return CS2STypeRowMask
	case "allmask":
		return CS2STypeAllMask
	}

	return CS2STypeAll
}

type ChgSymbols2SourceSymbolType int

const (
	CS2SSTypeNone         ChgSymbols2SourceSymbolType = 0
	CS2SSTypeSymbols      ChgSymbols2SourceSymbolType = 1
	CS2SSTypeSymbolWeight ChgSymbols2SourceSymbolType = 2
	CS2SSTypeNoEmpty      ChgSymbols2SourceSymbolType = 3
)

func parseChgSymbols2SourceSymbolType(str string) ChgSymbols2SourceSymbolType {
	switch str {
	case "symbols":
		return CS2SSTypeSymbols
	case "symbolweight":
		return CS2SSTypeSymbolWeight
	case "noempty":
		return CS2SSTypeNoEmpty
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
	CS2TypeEachPosRandom  ChgSymbols2Type = 5
	CS2TypeSymbols        ChgSymbols2Type = 6
)

func parseChgSymbols2Type(str string) ChgSymbols2Type {
	switch str {
	case "symbolweight":
		return CS2TypeSymbolWeight
	case "mystery":
		return CS2TypeMystery
	case "mysteryonreels":
		return CS2TypeMysteryOnReels
	case "upsymbol":
		return CS2TypeUpSymbol
	case "eachposrandom":
		return CS2TypeEachPosRandom
	case "symbols":
		return CS2TypeSymbols
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
	switch str {
	case "maxnumber":
		return CS2ETypeMaxNumber
	case "nosamereel":
		return CS2ETypeNoSameReel
	}

	return CS2ETypeNone
}

type ChgSymbols2Data struct {
	BasicComponentData
	Pos []int
	cfg *ChgSymbols2Config
}

// OnNewGame -
func (chgSymbolsData *ChgSymbols2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	chgSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (chgSymbolsData *ChgSymbols2Data) OnNewStep() {
	chgSymbolsData.UsedScenes = nil
	chgSymbolsData.Pos = nil
}

// Clone
func (chgSymbolsData *ChgSymbols2Data) Clone() IComponentData {
	target := &ChgSymbols2Data{
		BasicComponentData: chgSymbolsData.CloneBasicComponentData(),
		cfg:                chgSymbolsData.cfg,
	}

	return target
}

// BuildPBComponentData
func (chgSymbolsData *ChgSymbols2Data) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.ChgSymbols2Data{
		BasicComponentData: chgSymbolsData.BuildPBBasicComponentData(),
		Pos:                make([]int32, len(chgSymbolsData.Pos)),
	}

	for i, v := range chgSymbolsData.Pos {
		pbcd.Pos[i] = int32(v)
	}

	return pbcd
}

// GetPos -
func (chgSymbolsData *ChgSymbols2Data) GetPos() []int {
	return chgSymbolsData.Pos
}

// HasPos -
func (chgSymbolsData *ChgSymbols2Data) HasPos(x int, y int) bool {
	return goutils.IndexOfInt2Slice(chgSymbolsData.Pos, x, y, 0) >= 0
}

// AddPos -
func (chgSymbolsData *ChgSymbols2Data) AddPos(x int, y int) {
	chgSymbolsData.Pos = append(chgSymbolsData.Pos, x, y)
}

// ClearPos -
func (chgSymbolsData *ChgSymbols2Data) ClearPos() {
	chgSymbolsData.Pos = nil
}

// AddPosEx -
func (chgSymbolsData *ChgSymbols2Data) AddPosEx(x int, y int) {
	if !chgSymbolsData.HasPos(x, y) {
		chgSymbolsData.AddPos(x, y)
	}
}

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
	SrcRowMask            string                      `yaml:"srcRowMask" json:"srcRowMask"`
	SrcMask               string                      `yaml:"srcMask" json:"srcMask"`
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
	SymbolCollection      string                      `json:"symbolCollection"`
	MaxNumber             int                         `yaml:"maxNumber" json:"maxNumber"`
	RowMask               string                      `yaml:"rowMask" json:"rowMask"`
	MapControllers        map[string][]*Award         `yaml:"controllers" json:"controllers"`
	JumpToComponent       string                      `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// SetLinkComponent
func (cfg *ChgSymbols2Config) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type ChgSymbols2 struct {
	*BasicComponent `json:"-"`
	Config          *ChgSymbols2Config `json:"config"`
}

// Init -
func (chgSymbols2 *ChgSymbols2) Init(fn string, pool *GamePropertyPool) error {
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

	return chgSymbols2.InitEx(cfg, pool)
}

// InitEx -
func (chgSymbols2 *ChgSymbols2) InitEx(cfg any, pool *GamePropertyPool) error {
	chgSymbols2.Config = cfg.(*ChgSymbols2Config)
	chgSymbols2.Config.ComponentType = ChgSymbols2TypeName

	chgSymbols2.Config.SrcType = parseChgSymbols2SourceType(chgSymbols2.Config.StrSrcType)
	chgSymbols2.Config.SrcSymbolType = parseChgSymbols2SourceSymbolType(chgSymbols2.Config.StrSrcSymbolType)
	chgSymbols2.Config.Type = parseChgSymbols2Type(chgSymbols2.Config.StrType)
	chgSymbols2.Config.ExitType = parseChgSymbols2ExitType(chgSymbols2.Config.StrExitType)

	if chgSymbols2.Config.ExitType == CS2ETypeMaxNumber && chgSymbols2.Config.MaxNumber <= 0 {
		goutils.Error("ChgSymbols2.InitEx:MaxNumber",
			slog.Int("maxNumber", chgSymbols2.Config.MaxNumber),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	for _, s := range chgSymbols2.Config.SrcSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ChgSymbols2.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		chgSymbols2.Config.SrcSymbolCodes = append(chgSymbols2.Config.SrcSymbolCodes, sc)
	}

	blankSymbolCode, isok := pool.DefaultPaytables.MapSymbols[chgSymbols2.Config.BlankSymbol]
	if isok {
		chgSymbols2.Config.BlankSymbolCode = blankSymbolCode
	} else {
		chgSymbols2.Config.BlankSymbolCode = -1
	}

	symbolCode, isok := pool.DefaultPaytables.MapSymbols[chgSymbols2.Config.Symbol]
	if isok {
		chgSymbols2.Config.SymbolCode = symbolCode
	} else {
		chgSymbols2.Config.SymbolCode = -1
	}

	if chgSymbols2.Config.SrcSymbolWeight != "" {
		vw2, err := pool.LoadIntWeights(chgSymbols2.Config.SrcSymbolWeight, true)
		if err != nil {
			goutils.Error("ChgSymbols2.InitEx:LoadIntWeights",
				slog.String("SourceWeight", chgSymbols2.Config.SrcSymbolWeight),
				goutils.Err(err))

			return err
		}

		chgSymbols2.Config.SrcSymbolWeightVW2 = vw2
	}

	if chgSymbols2.Config.Weight != "" {
		vw2, err := pool.LoadIntWeights(chgSymbols2.Config.Weight, chgSymbols2.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ChgSymbols2.InitEx:LoadIntWeights",
				slog.String("Weight", chgSymbols2.Config.Weight),
				goutils.Err(err))

			return err
		}

		chgSymbols2.Config.WeightVW2 = vw2
	}

	for _, ctrls := range chgSymbols2.Config.MapControllers {
		for _, ctrl := range ctrls {
			ctrl.Init()
		}
	}

	chgSymbols2.onInit(&chgSymbols2.Config.BasicComponentConfig)

	return nil
}

func (chgSymbols2 *ChgSymbols2) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, err := gameProp.Pool.LoadIntWeights(str, true)
		if err != nil {
			goutils.Error("ChgSymbols2.getWeight:LoadIntWeights",
				goutils.Err(err))

			return nil
		}

		return vw2
	}

	return chgSymbols2.Config.WeightVW2
}

// OnProcControllers -
func (chgSymbols2 *ChgSymbols2) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	ctrls, isok := chgSymbols2.Config.MapControllers[strVal]
	if isok {
		gameProp.procAwards(plugin, ctrls, curpr, gp)
	}
}

// getSrcPos
func (chgSymbols2 *ChgSymbols2) getSrcPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin,
	gs *sgc7game.GameScene) ([]int, error) {

	pos := make([]int, 0, gameProp.GetVal(GamePropWidth)*gameProp.GetVal(GamePropHeight)*2)

	if chgSymbols2.Config.RowMask != "" {
		imaskd := gameProp.GetComponentDataWithName(chgSymbols2.Config.RowMask)
		if imaskd == nil {
			goutils.Error("ChgSymbols2.getSrcPos:RowMask:imaskd==nil",
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		maskarr := imaskd.GetMask()

		switch chgSymbols2.Config.SrcType {
		case CS2STypePositionCollection:
			for _, pc := range chgSymbols2.Config.SrcPositionCollection {
				curpos := gameProp.GetComponentPos(pc)
				if len(curpos) > 0 {
					for i := range len(curpos) / 2 {
						x := curpos[i*2]
						y := curpos[i*2+1]

						if !maskarr[y] {
							continue
						}

						if goutils.IndexOfInt2Slice(pos, x, y, 0) < 0 {
							pos = append(pos, x, y)
						}
					}
				}
			}
		case CS2STypeAll:
			for x := 0; x < gameProp.GetVal(GamePropWidth); x++ {
				for y := 0; y < gameProp.GetVal(GamePropHeight); y++ {
					if !maskarr[y] {
						continue
					}

					pos = append(pos, x, y)
				}
			}
		case CS2STypeRowMask:
			imaskd := gameProp.GetComponentDataWithName(chgSymbols2.Config.SrcRowMask)
			if imaskd != nil {
				arr := imaskd.GetMask()
				if len(arr) != gs.Height {
					goutils.Error("ChgSymbols2.getSrcPos:RowMask:len(arr)!=gs.Height",
						goutils.Err(ErrInvalidComponentConfig))

					return nil, ErrInvalidComponentConfig
				}

				for y := 0; y < gs.Height; y++ {
					if !maskarr[y] {
						continue
					}

					if arr[y] {
						for x := 0; x < gs.Width; x++ {
							pos = append(pos, x, y)
						}
					}
				}
			}
		case CS2STypeMask:
			imaskd := gameProp.GetComponentDataWithName(chgSymbols2.Config.SrcMask)
			if imaskd != nil {
				arr := imaskd.GetMask()
				if len(arr) != gs.Width {
					goutils.Error("ChgSymbols2.getSrcPos:Mask:len(arr)!=gs.Width",
						goutils.Err(ErrInvalidComponentConfig))

					return nil, ErrInvalidComponentConfig
				}

				for x := 0; x < gs.Width; x++ {
					if arr[x] {
						for y := 0; y < gs.Height; y++ {
							if !maskarr[y] {
								continue
							}

							pos = append(pos, x, y)
						}
					}
				}
			}
		case CS2STypeAllMask:
			imaskd1 := gameProp.GetComponentDataWithName(chgSymbols2.Config.SrcMask)
			imaskd2 := gameProp.GetComponentDataWithName(chgSymbols2.Config.SrcRowMask)
			if imaskd1 != nil && imaskd2 != nil {
				arr1 := imaskd1.GetMask()
				if len(arr1) != gs.Width {
					goutils.Error("ChgSymbols2.getSrcPos:CS2STypeAllMask:len(arr1)!=gs.Width",
						goutils.Err(ErrInvalidComponentConfig))

					return nil, ErrInvalidComponentConfig
				}

				arr2 := imaskd2.GetMask()
				if len(arr2) != gs.Height {
					goutils.Error("ChgSymbols2.getSrcPos:CS2STypeAllMask:len(arr2)!=gs.Height",
						goutils.Err(ErrInvalidComponentConfig))

					return nil, ErrInvalidComponentConfig
				}

				for x := 0; x < gs.Width; x++ {
					if arr1[x] {
						for y := 0; y < gs.Height; y++ {
							if !maskarr[y] {
								continue
							}

							if arr2[y] {
								for x := 0; x < gs.Width; x++ {
									pos = append(pos, x, y)
								}
							}
						}
					}
				}
			}
		default:
			goutils.Error("ChgSymbols2.getSrcPos:ErrUnsupportedSourceType",
				slog.String("srcType", chgSymbols2.Config.StrSrcType),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}
	} else {
		switch chgSymbols2.Config.SrcType {
		case CS2STypePositionCollection:
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
		case CS2STypeAll:
			for x := 0; x < gameProp.GetVal(GamePropWidth); x++ {
				for y := 0; y < gameProp.GetVal(GamePropHeight); y++ {
					pos = append(pos, x, y)
				}
			}
		case CS2STypeRowMask:
			imaskd := gameProp.GetComponentDataWithName(chgSymbols2.Config.SrcRowMask)
			if imaskd != nil {
				arr := imaskd.GetMask()
				if len(arr) != gs.Height {
					goutils.Error("ChgSymbols2.getSrcPos:RowMask:len(arr)!=gs.Height",
						goutils.Err(ErrInvalidComponentConfig))

					return nil, ErrInvalidComponentConfig
				}

				for y := 0; y < gs.Height; y++ {
					if arr[y] {
						for x := 0; x < gs.Width; x++ {
							pos = append(pos, x, y)
						}
					}
				}
			}
		case CS2STypeMask:
			imaskd := gameProp.GetComponentDataWithName(chgSymbols2.Config.SrcMask)
			if imaskd != nil {
				arr := imaskd.GetMask()
				if len(arr) != gs.Width {
					goutils.Error("ChgSymbols2.getSrcPos:Mask:len(arr)!=gs.Width",
						goutils.Err(ErrInvalidComponentConfig))

					return nil, ErrInvalidComponentConfig
				}

				for x := 0; x < gs.Width; x++ {
					if arr[x] {
						for y := 0; y < gs.Height; y++ {
							pos = append(pos, x, y)
						}
					}
				}
			}
		case CS2STypeAllMask:
			imaskd1 := gameProp.GetComponentDataWithName(chgSymbols2.Config.SrcMask)
			imaskd2 := gameProp.GetComponentDataWithName(chgSymbols2.Config.SrcRowMask)
			if imaskd1 != nil && imaskd2 != nil {
				arr1 := imaskd1.GetMask()
				if len(arr1) != gs.Width {
					goutils.Error("ChgSymbols2.getSrcPos:CS2STypeAllMask:len(arr1)!=gs.Width",
						goutils.Err(ErrInvalidComponentConfig))

					return nil, ErrInvalidComponentConfig
				}

				arr2 := imaskd2.GetMask()
				if len(arr2) != gs.Height {
					goutils.Error("ChgSymbols2.getSrcPos:CS2STypeAllMask:len(arr2)!=gs.Height",
						goutils.Err(ErrInvalidComponentConfig))

					return nil, ErrInvalidComponentConfig
				}

				for x := 0; x < gs.Width; x++ {
					if arr1[x] {
						for y := 0; y < gs.Height; y++ {
							if arr2[y] {
								for x := 0; x < gs.Width; x++ {
									pos = append(pos, x, y)
								}
							}
						}
					}
				}
			}
		default:
			goutils.Error("ChgSymbols2.getSrcPos:ErrUnsupportedSourceType",
				slog.String("srcType", chgSymbols2.Config.StrSrcType),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}
	}

	if len(pos) == 0 {
		return nil, nil
	}

	switch chgSymbols2.Config.SrcSymbolType {
	case CS2SSTypeNoEmpty:
		npos := make([]int, 0, gameProp.GetVal(GamePropWidth)*gameProp.GetVal(GamePropHeight)*2)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			if gs.Arr[x][y] >= 0 {
				npos = append(npos, x, y)
			}
		}

		return npos, nil
	case CS2SSTypeSymbols:
		npos := make([]int, 0, gameProp.GetVal(GamePropWidth)*gameProp.GetVal(GamePropHeight)*2)

		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			if slices.Contains(chgSymbols2.Config.SrcSymbolCodes, gs.Arr[x][y]) {
				npos = append(npos, x, y)
			}
		}

		return npos, nil
	case CS2SSTypeSymbolWeight:
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

// procPos
func (chgSymbols2 *ChgSymbols2) procPos(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	_ sgc7game.IPlayerState, _ *sgc7game.Stake, prs []*sgc7game.PlayResult, cd *ChgSymbols2Data) (string, error) {

	gs := chgSymbols2.GetTargetScene3(gameProp, curpr, prs, 0)

	pos, err := chgSymbols2.getSrcPos(gameProp, plugin, gs)
	if err != nil {
		goutils.Error("ChgSymbols2.procPos:getSrcPos",
			goutils.Err(err))

		return "", err
	}

	if len(pos) == 0 {
		if chgSymbols2.Config.IsAlwaysGen {
			if gs != nil {
				ngs := gs.CloneEx(gameProp.PoolScene)

				chgSymbols2.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
			}
		}

		nc := chgSymbols2.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	switch chgSymbols2.Config.Type {
	case CS2TypeSymbol:
		symbolCode := chgSymbols2.Config.SymbolCode
		ngs, err := chgSymbols2.procSymbolWithPos(gameProp, gs, pos, symbolCode, cd)
		if err != nil {
			goutils.Error("ChgSymbols2.procPos:procSymbolWithPos",
				goutils.Err(err))

			return "", err
		}

		if ngs != gs {
			chgSymbols2.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

			chgSymbols2.ProcControllers(gameProp, plugin, curpr, gp, -1,
				gameProp.Pool.DefaultPaytables.GetStringFromInt(symbolCode))
		} else {
			nc := chgSymbols2.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	case CS2TypeSymbolWeight:
		ngs, err := chgSymbols2.procSymbolWeightWithPos(gameProp, curpr, gp, plugin, gs, pos, cd)
		if err != nil {
			goutils.Error("ChgSymbols2.procPos:procSymbolWeightWithPos",
				goutils.Err(err))

			return "", err
		}

		if ngs != gs {
			chgSymbols2.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
		} else {
			nc := chgSymbols2.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	case CS2TypeEachPosRandom:
		ngs, err := chgSymbols2.procEachPosRandomWithPos(gameProp, curpr, gp, plugin, gs, pos, cd)
		if err != nil {
			goutils.Error("ChgSymbols2.procPos:procEachPosRandomWithPos",
				goutils.Err(err))

			return "", err
		}

		if ngs != gs {
			chgSymbols2.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)
		} else {
			nc := chgSymbols2.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	case CS2TypeSymbols:
		if chgSymbols2.Config.SymbolCollection == "" {
			goutils.Error("ChgSymbols2.procPos:SymbolCollection is empty",
				goutils.Err(ErrInvalidComponentConfig))

			return "", ErrInvalidComponentConfig
		}

		var symbolcodes []int
		isc := gameProp.GetComponentDataWithName(chgSymbols2.Config.SymbolCollection)
		if isc != nil {
			symbolcodes = isc.GetSymbolCodes()
		}

		ngs, err := chgSymbols2.procSymbolsWithPos(gameProp, plugin, gs, pos, symbolcodes, cd)
		if err != nil {
			goutils.Error("ChgSymbols2.procPos:procSymbolsWithPos",
				goutils.Err(err))

			return "", err
		}

		if ngs != gs {
			chgSymbols2.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

			for _, symbolCode := range symbolcodes {
				chgSymbols2.ProcControllers(gameProp, plugin, curpr, gp, -1,
					gameProp.Pool.DefaultPaytables.GetStringFromInt(symbolCode))
			}
		} else {
			nc := chgSymbols2.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	}

	chgSymbols2.ProcControllers(gameProp, plugin, curpr, gp, -1, "<trigger>")

	nc := chgSymbols2.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// procSymbolWithPos
func (chgSymbols2 *ChgSymbols2) procSymbolWithPos(gameProp *GameProperty, gs *sgc7game.GameScene, pos []int, symbolCode int, cd *ChgSymbols2Data) (*sgc7game.GameScene, error) {
	ngs := gs.CloneEx(gameProp.PoolScene)

	if chgSymbols2.Config.ExitType == CS2ETypeMaxNumber {
		curnum := 0
		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			ngs.Arr[x][y] = symbolCode

			cd.AddPos(x, y)

			curnum++
			if curnum >= chgSymbols2.Config.MaxNumber {
				break
			}
		}
	} else {
		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			ngs.Arr[x][y] = symbolCode

			cd.AddPos(x, y)
		}
	}

	return ngs, nil
}

// procSymbolsWithPos
func (chgSymbols2 *ChgSymbols2) procSymbolsWithPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, pos []int, symbolCodes []int, cd *ChgSymbols2Data) (*sgc7game.GameScene, error) {
	ngs := gs.CloneEx(gameProp.PoolScene)

	for _, v := range symbolCodes {
		cr, err := plugin.Random(context.Background(), len(pos)/2)
		if err != nil {
			goutils.Error("ChgSymbols2.procSymbolsWithPos:Random",
				goutils.Err(err))

			return nil, err
		}

		ngs.Arr[pos[cr*2]][pos[cr*2+1]] = v

		cd.AddPos(pos[cr*2], pos[cr*2+1])

		if len(pos) <= 2 {
			goutils.Error("ChgSymbols2.procSymbolsWithPos",
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		pos = append(pos[:cr*2], pos[cr*2+2:]...)
	}

	return ngs, nil
}

// procSymbolWeightWithPos
func (chgSymbols2 *ChgSymbols2) procSymbolWeightWithPos(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, pos []int, cd *ChgSymbols2Data) (*sgc7game.GameScene, error) {

	vw2 := chgSymbols2.getWeight(gameProp, &cd.BasicComponentData)
	curs, err := vw2.RandVal(plugin)
	if err != nil {
		goutils.Error("ChgSymbols2.procSymbolWeightWithPos:RandVal",
			goutils.Err(err))

		return nil, err
	}

	sc := curs.Int()

	if sc != chgSymbols2.Config.BlankSymbolCode {
		chgSymbols2.ProcControllers(gameProp, plugin, curpr, gp, -1,
			gameProp.Pool.DefaultPaytables.GetStringFromInt(sc))

		return chgSymbols2.procSymbolWithPos(gameProp, gs, pos, sc, cd)
	}

	return gs, nil
}

// procEachPosRandomWithPos
func (chgSymbols2 *ChgSymbols2) procEachPosRandomWithPos(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, pos []int, cd *ChgSymbols2Data) (*sgc7game.GameScene, error) {

	vw2 := chgSymbols2.getWeight(gameProp, &cd.BasicComponentData)

	ngs := gs

	if chgSymbols2.Config.ExitType == CS2ETypeMaxNumber {
		curnum := 0
		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			curs, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("ChgSymbols2.procEachPosRandomWithPos:RandVal",
					goutils.Err(err))

				return nil, err
			}

			sc := curs.Int()

			if sc != chgSymbols2.Config.BlankSymbolCode {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[x][y] = sc

				cd.AddPos(x, y)

				curnum++

				chgSymbols2.ProcControllers(gameProp, plugin, curpr, gp, -1,
					gameProp.Pool.DefaultPaytables.GetStringFromInt(sc))

				if curnum >= chgSymbols2.Config.MaxNumber {
					break
				}
			}
		}
	} else {
		for i := range len(pos) / 2 {
			x := pos[i*2]
			y := pos[i*2+1]

			curs, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("ChgSymbols2.procEachPosRandomWithPos:RandVal",
					goutils.Err(err))

				return nil, err
			}

			sc := curs.Int()

			if sc != chgSymbols2.Config.BlankSymbolCode {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[x][y] = sc

				cd.AddPos(x, y)

				chgSymbols2.ProcControllers(gameProp, plugin, curpr, gp, -1,
					gameProp.Pool.DefaultPaytables.GetStringFromInt(sc))
			}
		}
	}

	return ngs, nil
}

// playgame
func (chgSymbols2 *ChgSymbols2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*ChgSymbols2Data)

	cd.OnNewStep()

	return chgSymbols2.procPos(gameProp, curpr, gp, plugin, ps, stake, prs, cd)
}

// OnAsciiGame - outpur to asciigame
func (chgSymbols2 *ChgSymbols2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*ChgSymbols2Data)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after ChgSymbols2", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// GetAllLinkComponents - get all link components
func (chgSymbols2 *ChgSymbols2) GetAllLinkComponents() []string {
	return []string{chgSymbols2.Config.DefaultNextComponent, chgSymbols2.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (chgSymbols2 *ChgSymbols2) GetNextLinkComponents() []string {
	return []string{chgSymbols2.Config.DefaultNextComponent, chgSymbols2.Config.JumpToComponent}
}

// NewComponentData -
func (chgSymbols2 *ChgSymbols2) NewComponentData() IComponentData {
	return &ChgSymbols2Data{
		cfg: chgSymbols2.Config,
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
// "srcRowMask": "mask-left"
// "rowMask": "mask-height4"
// "symbolCollection": "bg-allms"
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
	SymbolCollection      string   `json:"symbolCollection"`
	MaxNumber             int      `json:"maxNumber"`
	SrcRowMask            string   `json:"srcRowMask"`
	SrcMask               string   `json:"srcMask"`
	RowMask               string   `json:"rowMask"`
}

func (jcfg *jsonChgSymbols2) build() *ChgSymbols2Config {
	cfg := &ChgSymbols2Config{
		StrSrcType:            strings.ToLower(jcfg.StrSrcType),
		StrSrcSymbolType:      strings.ToLower(jcfg.StrSrcSymbolType),
		StrType:               strings.ToLower(jcfg.StrType),
		StrExitType:           strings.ToLower(jcfg.StrExitType),
		SrcMask:               jcfg.SrcMask,
		SrcRowMask:            jcfg.SrcRowMask,
		IsAlwaysGen:           jcfg.IsAlwaysGen,
		Height:                jcfg.Height,
		Weight:                jcfg.Weight,
		BlankSymbol:           jcfg.BlankSymbol,
		SrcSymbols:            slices.Clone(jcfg.SrcSymbols),
		SrcPositionCollection: slices.Clone(jcfg.SrcPositionCollection),
		SrcSymbolWeight:       jcfg.SrcSymbolWeight,
		Symbol:                jcfg.Symbol,
		MaxNumber:             jcfg.MaxNumber,
		RowMask:               jcfg.RowMask,
		SymbolCollection:      jcfg.SymbolCollection,
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
		mapControllers, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseChgSymbols2:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapControllers = mapControllers
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

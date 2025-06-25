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
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const CascadingRegulatorTypeName = "cascadingRegulator"

const crMaxWeightVal = 1000000

type CascadingRegulatorType int

const (
	CRTypeMysteryOnReels CascadingRegulatorType = 0
)

type CascadingRegulatorWinType int

const (
	CRWTypeLines   CascadingRegulatorWinType = 0
	CRWTypeWays    CascadingRegulatorWinType = 1
	CRWTypeScatter CascadingRegulatorWinType = 2
	CRWTypeCluster CascadingRegulatorWinType = 3
)

func parseCascadingRegulatorType(strType string) CascadingRegulatorType {
	if strType == "mysteryonreels" {
		return CRTypeMysteryOnReels
	}

	return CRTypeMysteryOnReels
}

func parseCascadingRegulatorWinType(strType string) CascadingRegulatorWinType {
	if strType == "ways" {
		return CRWTypeWays
	} else if strType == "scatter" {
		return CRWTypeScatter
	} else if strType == "cluster" {
		return CRWTypeCluster
	}

	return CRWTypeLines
}

type CascadingRegulatorData struct {
	BasicComponentData
}

// OnNewGame -
func (cascadingRegulatorData *CascadingRegulatorData) OnNewGame(gameProp *GameProperty, component IComponent) {
	cascadingRegulatorData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (cascadingRegulatorData *CascadingRegulatorData) onNewStep() {
	cascadingRegulatorData.UsedScenes = nil
}

// Clone
func (cascadingRegulatorData *CascadingRegulatorData) Clone() IComponentData {
	target := &CascadingRegulatorData{
		BasicComponentData: cascadingRegulatorData.CloneBasicComponentData(),
	}

	return target
}

// BuildPBComponentData
func (cascadingRegulatorData *CascadingRegulatorData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.CascadingRegulatorData{
		BasicComponentData: cascadingRegulatorData.BuildPBBasicComponentData(),
	}

	return pbcd
}

// CascadingRegulatorConfig - configuration for CascadingRegulator
type CascadingRegulatorConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                    `yaml:"type" json:"type"`
	Type                 CascadingRegulatorType    `yaml:"-" json:"-"`
	MaxRespinNum         int                       `yaml:"maxRespinNum" json:"maxRespinNum"` // maxRespinNum
	Mystery              string                    `yaml:"mystery" json:"mystery"`           // mystery
	MysteryCode          int                       `yaml:"-" json:"-"`
	CoreRespin           string                    `yaml:"coreRespin" json:"coreRespin"`       // coreRespin
	StrWinType           string                    `yaml:"winType" json:"winType"`             // winType
	WinType              CascadingRegulatorWinType `yaml:"-" json:"-"`                         //
	MysteryWeight        string                    `yaml:"mysteryWeight" json:"mysteryWeight"` // mysteryWeight
	MysteryWeightVW      *sgc7game.ValWeights2     `yaml:"-" json:"-"`
	minSymbolWinNum      int                       `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *CascadingRegulatorConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type CascadingRegulator struct {
	*BasicComponent `json:"-"`
	Config          *CascadingRegulatorConfig `json:"config"`
}

// Init -
func (removeSymbols *CascadingRegulator) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("CascadingRegulator.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CascadingRegulatorConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("CascadingRegulator.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return removeSymbols.InitEx(cfg, pool)
}

// InitEx -
func (cascadingRegulator *CascadingRegulator) InitEx(cfg any, pool *GamePropertyPool) error {
	cascadingRegulator.Config = cfg.(*CascadingRegulatorConfig)
	cascadingRegulator.Config.ComponentType = CascadingRegulatorTypeName

	cascadingRegulator.Config.MysteryCode = pool.DefaultPaytables.MapSymbols[cascadingRegulator.Config.Mystery]
	cascadingRegulator.Config.Type = parseCascadingRegulatorType(cascadingRegulator.Config.StrType)
	cascadingRegulator.Config.WinType = parseCascadingRegulatorWinType(cascadingRegulator.Config.StrWinType)

	if cascadingRegulator.Config.MysteryWeight != "" {
		vw, err := pool.LoadSymbolWeights(cascadingRegulator.Config.MysteryWeight, "val", "weight", pool.DefaultPaytables, true)
		if err != nil {
			goutils.Error("CascadingRegulator.InitEx:LoadSymbolWeights",
				slog.String("mysteryWeight", cascadingRegulator.Config.MysteryWeight),
				goutils.Err(err))

			return err
		}

		cascadingRegulator.Config.MysteryWeightVW = vw

		cascadingRegulator.Config.minSymbolWinNum = pool.DefaultPaytables.GetSymbolMinWinNum(vw.Vals[0].Int())
	} else {
		goutils.Error("CascadingRegulator.InitEx:non-MysteryWeight",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	cascadingRegulator.onInit(&cascadingRegulator.Config.BasicComponentConfig)

	return nil
}

func (cascadingRegulator *CascadingRegulator) getMysteryWeight(gameProp *GameProperty, crcd *CascadingRegulatorData) *sgc7game.ValWeights2 {
	return cascadingRegulator.Config.MysteryWeightVW
}

// procMysteryOnReels
func (cascadingRegulator *CascadingRegulator) procMysteryOnReel(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, crcd *CascadingRegulatorData,
	ngs *sgc7game.GameScene, x int, yarr []int, level float32, invalidSyms []int, vw2 *sgc7game.ValWeights2) error {

	if cascadingRegulator.Config.WinType == CRWTypeScatter {
		if sgc7utils.IsEqualFloat32(level, 0) {
			// 这种模式下，如果 level 是 0，不需要特殊处理 mystery weight
			if len(invalidSyms) == 0 {
				cv, err := vw2.RandVal(plugin)
				if err != nil {
					goutils.Error("CascadingRegulator.procMysteryOnReel:RandVal",
						slog.Int("x", x),
						goutils.Err(err))

					return err
				}

				for _, y := range yarr {
					ngs.Arr[x][y] = cv.Int()
				}

				return nil
			}

			nvw2 := vw2.Clone()
			symbols := vw2.GetIntVals()
			for _, cs := range invalidSyms {
				csi := slices.Index(symbols, cs)
				if csi >= 0 {
					nvw2.Weights[csi] = 0
				}
			}

			nvw2.ResetMaxWeight()

			if nvw2.MaxWeight <= 0 {
				cv, err := vw2.RandVal(plugin)
				if err != nil {
					goutils.Error("CascadingRegulator.procMysteryOnReel:RandVal",
						slog.Int("x", x),
						goutils.Err(err))

					return err
				}

				for _, y := range yarr {
					ngs.Arr[x][y] = cv.Int()
				}

				return nil
			}

			cv, err := nvw2.RandVal(plugin)
			if err != nil {
				goutils.Error("CascadingRegulator.procMysteryOnReel:RandVal",
					slog.Int("x", x),
					goutils.Err(err))

				return err
			}

			for _, y := range yarr {
				ngs.Arr[x][y] = cv.Int()
			}

			return nil
		}

		// scatter 时，应该先找出一定不能变的 symbol，然后降低这些 symbol 的权重，然后 roll 出 mystery 需要变成的 symbol
		symbols := vw2.GetIntVals()
		arrNum := ngs.CountSymbols(symbols)

		for i, v := range symbols {
			if arrNum[i] >= cascadingRegulator.Config.minSymbolWinNum {
				// 这种情况下，不能变成 mystery 的 symbol
				if slices.Index(invalidSyms, v) < 0 {
					invalidSyms = append(invalidSyms, v)
				}
			}
		}

		if level >= 1 || sgc7utils.IsEqualFloat32(level, 1) {
			// 这种模式下，如果 level 是 1，也只需要处理 invalidSyms 即可
			if len(invalidSyms) == 0 {
				cv, err := vw2.RandVal(plugin)
				if err != nil {
					goutils.Error("CascadingRegulator.procMysteryOnReel:RandVal",
						slog.Int("x", x),
						goutils.Err(err))

					return err
				}

				for _, y := range yarr {
					ngs.Arr[x][y] = cv.Int()
				}

				return nil
			}

			nvw2 := vw2.Clone()
			maxWeight := 0
			for ci, cs := range symbols {
				if slices.Index(invalidSyms, cs) < 0 {
					nvw2.Weights[ci] = int(float32(nvw2.Weights[ci]) / float32(nvw2.MaxWeight) * (1 - float32(arrNum[ci])/float32(cascadingRegulator.Config.minSymbolWinNum)) * crMaxWeightVal)
				} else {
					nvw2.Weights[ci] = 0
				}

				maxWeight += nvw2.Weights[ci]
			}

			nvw2.MaxWeight = maxWeight

			if nvw2.MaxWeight <= 0 {
				cv, err := vw2.RandVal(plugin)
				if err != nil {
					goutils.Error("CascadingRegulator.procMysteryOnReel:RandVal",
						slog.Int("x", x),
						goutils.Err(err))

					return err
				}

				for _, y := range yarr {
					ngs.Arr[x][y] = cv.Int()
				}

				return nil
			}

			cv, err := nvw2.RandVal(plugin)
			if err != nil {
				goutils.Error("CascadingRegulator.procMysteryOnReel:RandVal",
					slog.Int("x", x),
					goutils.Err(err))

				return err
			}

			for _, y := range yarr {
				ngs.Arr[x][y] = cv.Int()
			}

			return nil
		}

		// 否则，需要处理权重
		nvw2 := vw2.Clone()
		maxWeight := 0
		for ci, cs := range symbols {
			if slices.Index(invalidSyms, cs) < 0 {
				nvw2.Weights[ci] = int(float32(nvw2.Weights[ci]) / float32(nvw2.MaxWeight) * (1 - float32(arrNum[ci])/float32(cascadingRegulator.Config.minSymbolWinNum)) * crMaxWeightVal)
			} else {
				nvw2.Weights[ci] = int(float32(nvw2.Weights[ci]) / float32(nvw2.MaxWeight) * (1 - level) * 0.1 * crMaxWeightVal)
			}

			maxWeight += nvw2.Weights[ci]
		}

		nvw2.MaxWeight = maxWeight

		if nvw2.MaxWeight <= 0 {
			cv, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("CascadingRegulator.procMysteryOnReel:RandVal",
					slog.Int("x", x),
					goutils.Err(err))

				return err
			}

			for _, y := range yarr {
				ngs.Arr[x][y] = cv.Int()
			}

			return nil
		}

		cv, err := nvw2.RandVal(plugin)
		if err != nil {
			goutils.Error("CascadingRegulator.procMysteryOnReel:RandVal",
				slog.Int("x", x),
				goutils.Err(err))

			return err
		}

		for _, y := range yarr {
			ngs.Arr[x][y] = cv.Int()
		}

		return nil
	}

	return nil
}

// procMysteryOnReels
func (cascadingRegulator *CascadingRegulator) procMysteryOnReels(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, crcd *CascadingRegulatorData,
	gs *sgc7game.GameScene, level float32) error {
	ngs := gs
	vw2 := cascadingRegulator.getMysteryWeight(gameProp, crcd)

	for x, arr := range gs.Arr {
		invalidSyms := make([]int, 0, len(vw2.Weights))
		yarr := make([]int, 0, len(arr))
		prey := -1
		nexty := -1

		for y, v := range arr {
			if v == cascadingRegulator.Config.MysteryCode {
				if y > 0 && prey == y-1 {
					if slices.Index(invalidSyms, arr[prey]) < 0 {
						invalidSyms = append(invalidSyms, arr[prey])
					}
				}

				yarr = append(yarr, y)

				nexty = y + 1
			} else {
				if nexty == y {
					if slices.Index(invalidSyms, arr[nexty]) < 0 {
						invalidSyms = append(invalidSyms, arr[nexty])
					}
				}

				prey = y
			}
		}

		if len(yarr) > 0 {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			cascadingRegulator.procMysteryOnReel(gameProp, curpr, gp, plugin, crcd, ngs, x, yarr, level, invalidSyms, vw2)
		}
	}

	if ngs == gs {
		return ErrComponentDoNothing
	}

	cascadingRegulator.AddScene(gameProp, curpr, ngs, &crcd.BasicComponentData)

	return nil
}

// playgame
func (cascadingRegulator *CascadingRegulator) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	crcd := cd.(*CascadingRegulatorData)
	crcd.onNewStep()

	gs := cascadingRegulator.GetTargetScene3(gameProp, curpr, prs, 0)

	respincd := gameProp.GetComponentDataWithName(cascadingRegulator.Config.CoreRespin)
	if respincd == nil {
		goutils.Error("CascadingRegulator.OnPlayGame:GetComponentDataWithName",
			slog.String("name", cascadingRegulator.Config.CoreRespin),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	currespinnum := respincd.GetCurRespinNum()
	level := float32(currespinnum) / float32(cascadingRegulator.Config.MaxRespinNum)

	if cascadingRegulator.Config.Type == CRTypeMysteryOnReels {
		err := cascadingRegulator.procMysteryOnReels(gameProp, curpr, gp, plugin, crcd, gs, level)
		if err != nil {
			if err == ErrComponentDoNothing {
				nc := cascadingRegulator.onStepEnd(gameProp, curpr, gp, "")

				return nc, err
			}

			goutils.Error("CascadingRegulator.OnPlayGame:procMysteryOnReels",
				goutils.Err(err))

			return "", err
		}
	} else {
		goutils.Error("CascadingRegulator.OnPlayGame:InvalidType",
			slog.Int("type", int(cascadingRegulator.Config.Type)),
			goutils.Err(ErrIvalidComponentConfig))

		return "", ErrIvalidComponentConfig
	}

	nc := cascadingRegulator.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (cascadingRegulator *CascadingRegulator) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*CascadingRegulatorData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after cascadingRegulator", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// NewComponentData -
func (cascadingRegulator *CascadingRegulator) NewComponentData() IComponentData {
	return &CascadingRegulatorData{}
}

// EachUsedResults -
func (cascadingRegulator *CascadingRegulator) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewCascadingRegulator(name string) IComponent {
	return &CascadingRegulator{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "mysteryOnReels",
// "maxRespinNum": 6,
// "mystery": "MY",
// "coreRespin": "bg-respin",
// "winType": "scatter",
// "mysteryWeight": "mysteryweight"
type jsonCascadingRegulator struct {
	Type          string `json:"type"`          // type
	MaxRespinNum  int    `json:"maxRespinNum"`  // maxRespinNum
	Mystery       string `json:"mystery"`       // mystery
	CoreRespin    string `json:"coreRespin"`    // coreRespin
	WinType       string `json:"winType"`       // winType
	MysteryWeight string `json:"mysteryWeight"` // mysteryWeight
}

func (jcfg *jsonCascadingRegulator) build() *CascadingRegulatorConfig {
	cfg := &CascadingRegulatorConfig{
		StrType:       strings.ToLower(jcfg.Type),
		MaxRespinNum:  jcfg.MaxRespinNum,
		Mystery:       jcfg.Mystery,
		CoreRespin:    jcfg.CoreRespin,
		StrWinType:    strings.ToLower(jcfg.WinType),
		MysteryWeight: jcfg.MysteryWeight,
	}

	return cfg
}

func parseCascadingRegulator(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCascadingRegulator:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCascadingRegulator:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCascadingRegulator{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCascadingRegulator:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CascadingRegulatorTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

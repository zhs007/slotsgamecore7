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
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
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

var gDiffusionPos [][]int

func init() {
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(1, 0)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(1, 1)...)

	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(2, 0)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(2, 1)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(2, 2)...)

	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(3, 0)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(3, 1)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(3, 2)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(3, 3)...)

	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(4, 0)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(4, 1)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(4, 2)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(4, 3)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(4, 4)...)

	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(5, 0)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(5, 1)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(5, 2)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(5, 3)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(5, 4)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(5, 5)...)

	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(6, 0)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(6, 1)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(6, 2)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(6, 3)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(6, 4)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(6, 5)...)
	gDiffusionPos = append(gDiffusionPos, getDiffusionPosOffset2(6, 6)...)
}

func getDiffusionPosOffset2(instance int, off int) [][]int {
	pos := [][]int{}

	if off == 0 {
		pos = append(pos, []int{0, -instance})
		pos = append(pos, []int{instance, 0})
		pos = append(pos, []int{0, instance})
		pos = append(pos, []int{-instance, 0})

		return pos
	}

	pos = append(pos, []int{off, -instance})
	pos = append(pos, []int{instance, off})
	pos = append(pos, []int{-off, instance})
	pos = append(pos, []int{-instance, -off})

	return pos
}

func getDiffusionPosOffset(index int) (int, int) {
	return gDiffusionPos[index][0], gDiffusionPos[index][1]
}

func burstDiffusion(gs *sgc7game.GameScene, x int, y int, num int, overrideSym int, ignoreSyms []int, bsd *BurstSymbolsData) {
	bsd.newData()

	curnum := 0

	for i := 0; i < len(gDiffusionPos); i++ {
		cx, cy := getDiffusionPosOffset(i)

		tx := x + cx
		ty := y + cy

		if tx >= 0 && ty >= 0 && tx < gs.Width && ty < gs.Height {
			if goutils.IndexOfIntSlice(ignoreSyms, gs.Arr[tx][ty], 0) < 0 {
				bsd.AddPos(tx, ty)
				gs.Arr[tx][ty] = overrideSym
				curnum++

				if curnum >= num {
					return
				}
			}
		}
	}
}

type BurstSymbolsData struct {
	BasicComponentData
	Pos [][]int
}

// OnNewGame -
func (burstSymbolsData *BurstSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	burstSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (burstSymbolsData *BurstSymbolsData) OnNewStep() {
	burstSymbolsData.UsedScenes = nil
	burstSymbolsData.Pos = nil
}

// Clone
func (burstSymbolsData *BurstSymbolsData) Clone() IComponentData {
	target := &BurstSymbolsData{
		BasicComponentData: burstSymbolsData.CloneBasicComponentData(),
	}

	target.Pos = make([][]int, len(burstSymbolsData.Pos))
	for _, arr := range burstSymbolsData.Pos {
		dstarr := make([]int, len(arr))
		copy(dstarr, arr)
		target.Pos = append(target.Pos, dstarr)
	}

	return target
}

// BuildPBComponentData
func (burstSymbolsData *BurstSymbolsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.BurstSymbolsData{
		BasicComponentData: burstSymbolsData.BuildPBBasicComponentData(),
	}

	num := 0
	for _, arr := range burstSymbolsData.Pos {
		num += len(arr)
		num++
	}

	pbcd.Pos = make([]int32, 0, num)

	for _, arr := range burstSymbolsData.Pos {
		for _, s := range arr {
			pbcd.Pos = append(pbcd.Pos, int32(s))
		}

		pbcd.Pos = append(pbcd.Pos, -1)
	}

	return pbcd
}

// GetPos -
func (burstSymbolsData *BurstSymbolsData) GetPos() []int {
	num := 0
	for _, arr := range burstSymbolsData.Pos {
		num += len(arr)
	}

	newpos := make([]int, 0, num)

	for _, arr := range burstSymbolsData.Pos {
		newpos = append(newpos, arr...)
	}

	return newpos
}

// HasPos -
func (burstSymbolsData *BurstSymbolsData) HasPos(x int, y int) bool {
	for _, arr := range burstSymbolsData.Pos {
		if goutils.IndexOfInt2Slice(arr, x, y, 0) >= 0 {
			return true
		}
	}

	return false
}

// AddPos -
func (burstSymbolsData *BurstSymbolsData) AddPos(x int, y int) {
	if len(burstSymbolsData.Pos) == 0 {
		burstSymbolsData.Pos = append(burstSymbolsData.Pos, []int{})
	}

	burstSymbolsData.Pos[len(burstSymbolsData.Pos)-1] = append(burstSymbolsData.Pos[len(burstSymbolsData.Pos)-1], x, y)
}

// // AddPosEx -
// func (catchSymbolsData *CatchSymbolsData) AddPosEx(x int, y int) {
// 	if goutils.IndexOfInt2Slice(catchSymbolsData.Pos[len(catchSymbolsData.Pos)-1], x, y, 0) < 0 {
// 		catchSymbolsData.Pos[len(catchSymbolsData.Pos)-1] = append(catchSymbolsData.Pos[len(catchSymbolsData.Pos)-1], x, y)
// 	}
// }

// newData -
func (burstSymbolsData *BurstSymbolsData) newData() {
	burstSymbolsData.Pos = append(burstSymbolsData.Pos, []int{})
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
	IgnoreSymbols        []string               `yaml:"ignoreSymbols" json:"ignoreSymbols"`
	IgnoreSymbolCodes    []int                  `yaml:"-" json:"-"`
	OverrideSymbol       string                 `yaml:"overrideSymbol" json:"overrideSymbol"`
	OverrideSymbolCode   int                    `yaml:"-" json:"-"`
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

	for _, s := range burstSymbols.Config.SourceSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("BurstSymbols.InitEx:SourceSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		burstSymbols.Config.SourceSymbolCodes = append(burstSymbols.Config.SourceSymbolCodes, sc)
	}

	for _, s := range burstSymbols.Config.IgnoreSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("BurstSymbols.InitEx:IgnoreSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		burstSymbols.Config.IgnoreSymbolCodes = append(burstSymbols.Config.IgnoreSymbolCodes, sc)
	}

	sc0, isok := pool.DefaultPaytables.MapSymbols[burstSymbols.Config.OverrideSymbol]
	if !isok {
		goutils.Error("BurstSymbols.InitEx:OverrideSymbol",
			slog.String("symbol", burstSymbols.Config.OverrideSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	burstSymbols.Config.OverrideSymbolCode = sc0

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

	bsd := cd.(*BurstSymbolsData)

	bsd.OnNewStep()

	gs := burstSymbols.GetTargetScene3(gameProp, curpr, prs, 0)

	if burstSymbols.Config.SourceType == BSSTypePositionCollection {
		pos := gameProp.GetComponentPos(burstSymbols.Config.PositionCollection)
		if len(pos) >= 2 {
			sc2 := gs.CloneEx(gameProp.PoolScene)

			for i := 0; i < len(pos)/2; i++ {
				x := pos[i*2]
				y := pos[i*2+1]

				burstDiffusion(sc2, x, y, burstSymbols.Config.BurstNumber, burstSymbols.Config.OverrideSymbolCode, burstSymbols.Config.IgnoreSymbolCodes, bsd)
			}

			burstSymbols.AddScene(gameProp, curpr, sc2, &bsd.BasicComponentData)

			if len(burstSymbols.Config.Controllers) > 0 {
				gameProp.procAwards(plugin, burstSymbols.Config.Controllers, curpr, gp)
			}

			nc := burstSymbols.onStepEnd(gameProp, curpr, gp, burstSymbols.Config.JumpToComponent)

			return nc, nil
		}
	}

	nc := burstSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing

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

	// if sc2 == gs {
	// 	nc := burstSymbols.onStepEnd(gameProp, curpr, gp, "")

	// 	return nc, ErrComponentDoNothing
	// }

	// burstSymbols.AddScene(gameProp, curpr, sc2, &bsd.BasicComponentData)

	// if len(burstSymbols.Config.Controllers) > 0 {
	// 	gameProp.procAwards(plugin, burstSymbols.Config.Controllers, curpr, gp)
	// }

	// nc := burstSymbols.onStepEnd(gameProp, curpr, gp, burstSymbols.Config.JumpToComponent)

	// return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (burstSymbols *BurstSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bsd := cd.(*BurstSymbolsData)

	asciigame.OutputScene("after burstSymbols", pr.Scenes[bsd.UsedScenes[0]], mapSymbolColor)

	return nil
}

// NewComponentData -
func (burstSymbols *BurstSymbols) NewComponentData() IComponentData {
	return &BurstSymbolsData{}
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

// "burstSymbolsSourceType": "positionCollection",
// "burstType": "diffusion",
// "burstNumber": 4,
// "sourcePositionCollection": "bg-burstpos",
// "ignoreSymbols": [
//
//	"RW",
//	"MY",
//	"SC",
//	"RW2",
//	"MM"
//
// ],
// "overrideSymbol": "MY"
type jsonBurstSymbols struct {
	SourceType               string   `json:"burstSymbolsSourceType"`
	BurstType                string   `json:"burstType"`
	BurstNumber              int      `json:"burstNumber"`
	SourceSymbols            []string `json:"sourceSymbols"`
	OverrideSymbol           string   `json:"overrideSymbol"`
	IgnoreSymbols            []string `json:"ignoreSymbols"`
	SourcePositionCollection string   `json:"sourcePositionCollection"`
}

func (jcfg *jsonBurstSymbols) build() *BurstSymbolsConfig {
	cfg := &BurstSymbolsConfig{
		StrBurstType:       jcfg.BurstType,
		StrSourceType:      jcfg.SourceType,
		BurstNumber:        jcfg.BurstNumber,
		SourceSymbols:      jcfg.SourceSymbols,
		OverrideSymbol:     jcfg.OverrideSymbol,
		IgnoreSymbols:      jcfg.IgnoreSymbols,
		PositionCollection: jcfg.SourcePositionCollection,
	}

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

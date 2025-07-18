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
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const RandomMoveSymbolsTypeName = "randomMoveSymbols"

type RandomMoveSymbolsType int

const (
	RMSTypeNormal RandomMoveSymbolsType = 0 // normal
	RMSTypeReels  RandomMoveSymbolsType = 1 // reels
)

func parseRandomMoveSymbolsType(str string) RandomMoveSymbolsType {
	if str == "reels" {
		return RMSTypeReels
	}

	return RMSTypeNormal
}

type RandomMoveSymbolsData struct {
	BasicComponentData
	Pos [][]int
}

// OnNewGame -
func (randomMoveSymbolsData *RandomMoveSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	randomMoveSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (randomMoveSymbolsData *RandomMoveSymbolsData) OnNewStep() {
	randomMoveSymbolsData.UsedScenes = nil
	randomMoveSymbolsData.Pos = nil
}

// Clone
func (randomMoveSymbolsData *RandomMoveSymbolsData) Clone() IComponentData {
	target := &RandomMoveSymbolsData{
		BasicComponentData: randomMoveSymbolsData.CloneBasicComponentData(),
	}

	target.Pos = make([][]int, len(randomMoveSymbolsData.Pos))
	for _, arr := range randomMoveSymbolsData.Pos {
		dstarr := make([]int, len(arr))
		copy(dstarr, arr)
		target.Pos = append(target.Pos, dstarr)
	}

	return target
}

// BuildPBComponentData
func (randomMoveSymbolsData *RandomMoveSymbolsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RandomMoveSymbolsData{
		BasicComponentData: randomMoveSymbolsData.BuildPBBasicComponentData(),
	}

	num := 0
	for _, arr := range randomMoveSymbolsData.Pos {
		num += len(arr)
		num++
	}

	pbcd.Pos = make([]int32, 0, num)

	for _, arr := range randomMoveSymbolsData.Pos {
		for _, s := range arr {
			pbcd.Pos = append(pbcd.Pos, int32(s))
		}

		pbcd.Pos = append(pbcd.Pos, -1)
	}

	return pbcd
}

// GetPos -
func (randomMoveSymbolsData *RandomMoveSymbolsData) GetPos() []int {
	num := 0
	for _, arr := range randomMoveSymbolsData.Pos {
		num += len(arr)
	}

	newpos := make([]int, 0, num)

	for _, arr := range randomMoveSymbolsData.Pos {
		newpos = append(newpos, arr...)
	}

	return newpos
}

// HasPos -
func (randomMoveSymbolsData *RandomMoveSymbolsData) HasPos(x int, y int) bool {
	for _, arr := range randomMoveSymbolsData.Pos {
		if goutils.IndexOfInt2Slice(arr, x, y, 0) >= 0 {
			return true
		}
	}

	return false
}

// AddPos -
func (randomMoveSymbolsData *RandomMoveSymbolsData) AddPos(x int, y int) {
	if len(randomMoveSymbolsData.Pos) == 0 {
		randomMoveSymbolsData.Pos = append(randomMoveSymbolsData.Pos, []int{})
	}

	randomMoveSymbolsData.Pos[len(randomMoveSymbolsData.Pos)-1] = append(randomMoveSymbolsData.Pos[len(randomMoveSymbolsData.Pos)-1], x, y)
}

// ClearPos -
func (randomMoveSymbolsData *RandomMoveSymbolsData) ClearPos() {
	randomMoveSymbolsData.Pos = nil
}

// AddPosEx -
func (randomMoveSymbolsData *RandomMoveSymbolsData) AddPosEx(x int, y int) {
	if goutils.IndexOfInt2Slice(randomMoveSymbolsData.Pos[len(randomMoveSymbolsData.Pos)-1], x, y, 0) < 0 {
		randomMoveSymbolsData.Pos[len(randomMoveSymbolsData.Pos)-1] = append(randomMoveSymbolsData.Pos[len(randomMoveSymbolsData.Pos)-1], x, y)
	}
}

// newData -
func (randomMoveSymbolsData *RandomMoveSymbolsData) newData() {
	randomMoveSymbolsData.Pos = append(randomMoveSymbolsData.Pos, []int{})
}

// RandomMoveSymbolsConfig - configuration for RandomMoveSymbols
type RandomMoveSymbolsConfig struct {
	BasicComponentConfig     `yaml:",inline" json:",inline"`
	Type                     RandomMoveSymbolsType `yaml:"-" json:"-"`
	StrType                  string                `yaml:"type" json:"type"`
	TargetSymbols            []string              `yaml:"targetSymbols" json:"targetSymbols"`
	TargetSymbolCodes        []int                 `yaml:"-" json:"-"`
	IgnoreSymbols            []string              `yaml:"ignoreSymbols" json:"ignoreSymbols"`
	IgnoreSymbolCodes        []int                 `yaml:"-" json:"-"`
	ReelsWeight              string                `yaml:"reelsWeight" json:"reelsWeight"`
	ReelsWeightVW2           *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	TargetPositionCollection string                `yaml:"targetPositionCollection" json:"targetPositionCollection"`
	Controllers              []*Award              `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
func (cfg *RandomMoveSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RandomMoveSymbols struct {
	*BasicComponent `json:"-"`
	Config          *RandomMoveSymbolsConfig `json:"config"`
}

// Init -
func (randomMoveSymbols *RandomMoveSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RandomMoveSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RandomMoveSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RandomMoveSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return randomMoveSymbols.InitEx(cfg, pool)
}

// InitEx -
func (randomMoveSymbols *RandomMoveSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	randomMoveSymbols.Config = cfg.(*RandomMoveSymbolsConfig)
	randomMoveSymbols.Config.ComponentType = RandomMoveSymbolsTypeName

	randomMoveSymbols.Config.Type = parseRandomMoveSymbolsType(randomMoveSymbols.Config.StrType)

	for _, v := range randomMoveSymbols.Config.TargetSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("RandomMoveSymbols.InitEx:TargetSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		randomMoveSymbols.Config.TargetSymbolCodes = append(randomMoveSymbols.Config.TargetSymbolCodes, sc)
	}

	for _, v := range randomMoveSymbols.Config.IgnoreSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("RandomMoveSymbols.InitEx:IgnoreSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		randomMoveSymbols.Config.IgnoreSymbolCodes = append(randomMoveSymbols.Config.IgnoreSymbolCodes, sc)
	}

	if randomMoveSymbols.Config.ReelsWeight != "" {
		vw2, err := pool.LoadIntWeights(randomMoveSymbols.Config.ReelsWeight, true)
		if err != nil {
			goutils.Error("ChgSymbols.InitEx:LoadIntWeights",
				slog.String("ReelsWeight", randomMoveSymbols.Config.ReelsWeight),
				goutils.Err(err))

			return err
		}

		randomMoveSymbols.Config.ReelsWeightVW2 = vw2
	}

	for _, award := range randomMoveSymbols.Config.Controllers {
		award.Init()
	}

	randomMoveSymbols.onInit(&randomMoveSymbols.Config.BasicComponentConfig)

	return nil
}

func (randomMoveSymbols *RandomMoveSymbols) procNormal(gameProp *GameProperty, cd *RandomMoveSymbolsData,
	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene) (*sgc7game.GameScene, error) {

	posTarget := make([]int, 0, gs.Width*gs.Height*2)
	var posSrc []int

	if randomMoveSymbols.Config.TargetPositionCollection != "" {
		pccd := gameProp.GetComponentDataWithName(randomMoveSymbols.Config.TargetPositionCollection)
		if pccd == nil {
			goutils.Error("RandomMoveSymbols.procNormal:GetComponentDataWithName",
				slog.String("TargetPositionCollection", randomMoveSymbols.Config.TargetPositionCollection),
				goutils.Err(ErrNoComponent))

			return gs, ErrNoComponent
		}

		posSrc = pccd.GetPos()
	} else {
		posSrc = make([]int, 0, gs.Width*gs.Height*2)

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(randomMoveSymbols.Config.TargetSymbolCodes, s, 0) >= 0 {
					posSrc = append(posSrc, x, y)
				} else if goutils.IndexOfIntSlice(randomMoveSymbols.Config.IgnoreSymbolCodes, s, 0) < 0 {
					posTarget = append(posTarget, x, y)
				}
			}
		}
	}

	if len(posSrc) == 0 || len(posTarget) == 0 {
		return gs, nil
	}

	ngs := gs

	for i := 0; i < len(posSrc)/2; i++ {
		cd.newData()

		x := posSrc[i*2]
		y := posSrc[i*2+1]

		cr, err := plugin.Random(context.Background(), len(posTarget)/2)
		if err != nil {
			goutils.Error("RandomMoveSymbols.procNormal:Random",
				goutils.Err(err))

			return gs, err
		}

		tx := posTarget[cr*2]
		ty := posTarget[cr*2+1]

		cd.AddPos(x, y)
		cd.AddPos(tx, ty)

		if ngs == gs {
			ngs = gs.CloneEx(gameProp.PoolScene)
		}

		ngs.Arr[tx][ty] = ngs.Arr[x][y]

		if len(posTarget) == 2 {
			break
		}

		posTarget = append(posTarget[:cr*2], posTarget[(cr+1)*2:]...)
	}

	return ngs, nil
}

func (randomMoveSymbols *RandomMoveSymbols) procReels(gameProp *GameProperty, cd *RandomMoveSymbolsData,
	plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, gsSrc *sgc7game.GameScene) (*sgc7game.GameScene, error) {

	if randomMoveSymbols.Config.ReelsWeightVW2 == nil {
		goutils.Error("RandomMoveSymbols.procReels",
			goutils.Err(ErrInvalidComponentConfig))

		return gs, ErrInvalidComponentConfig
	}

	var posSrc []int

	if randomMoveSymbols.Config.TargetPositionCollection != "" {
		pccd := gameProp.GetComponentDataWithName(randomMoveSymbols.Config.TargetPositionCollection)
		if pccd == nil {
			goutils.Error("RandomMoveSymbols.procReels:GetComponentDataWithName",
				slog.String("TargetPositionCollection", randomMoveSymbols.Config.TargetPositionCollection),
				goutils.Err(ErrNoComponent))

			return gs, ErrNoComponent
		}

		posSrc = pccd.GetPos()
	} else {
		posSrc = make([]int, 0, gs.Width*gs.Height*2)

		for x, arr := range gsSrc.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(randomMoveSymbols.Config.TargetSymbolCodes, s, 0) >= 0 {
					posSrc = append(posSrc, x, y)
				}
			}
		}
	}

	if len(posSrc) == 0 {
		return gs, nil
	}

	ngs := gs

	for i := 0; i < len(posSrc)/2; i++ {
		vw2 := randomMoveSymbols.Config.ReelsWeightVW2.Clone()

		cd.newData()

		x := posSrc[i*2]
		y := posSrc[i*2+1]

		for {
			ret, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("RandomMoveSymbols.procReels:Random",
					goutils.Err(err))

				return gs, err
			}

			tx := ret.Int() - 1
			yArr := make([]int, 0, gs.Height)

			for ty := 0; ty < gs.Height; ty++ {
				if x == tx && y == ty {
					yArr = append(yArr, ty)
				} else {
					if goutils.IndexOfInt2Slice(posSrc, tx, ty, i+1) < 0 {
						s := ngs.Arr[tx][ty]

						if goutils.IndexOfIntSlice(randomMoveSymbols.Config.TargetSymbolCodes, s, 0) < 0 && goutils.IndexOfIntSlice(randomMoveSymbols.Config.IgnoreSymbolCodes, s, 0) < 0 {
							yArr = append(yArr, ty)
						}
					}
				}
			}

			if len(yArr) <= 0 {
				if len(vw2.Vals) == 1 {
					return ngs, nil
				}

				err := vw2.RemoveVal(ret)
				if err != nil {
					goutils.Error("RandomMoveSymbols.procReels:RemoveVal",
						goutils.Err(err))

					return gs, err
				}
			} else {
				cr, err := plugin.Random(context.Background(), len(yArr))
				if err != nil {
					goutils.Error("RandomMoveSymbols.procReels:Random",
						goutils.Err(err))

					return gs, err
				}

				ty := yArr[cr]

				cd.AddPos(x, y)
				cd.AddPos(tx, ty)

				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[tx][ty] = gsSrc.Arr[x][y]

				break
			}
		}
	}

	return ngs, nil
}

// OnProcControllers -
func (randomMoveSymbols *RandomMoveSymbols) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(randomMoveSymbols.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, randomMoveSymbols.Config.Controllers, curpr, gp)
	}
}

// playgame
func (randomMoveSymbols *RandomMoveSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	msd := icd.(*RandomMoveSymbolsData)

	msd.OnNewStep()

	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)
	gs2 := gameProp.SceneStack.GetPreTopSceneEx(curpr, prs)

	sc2 := gs

	switch randomMoveSymbols.Config.Type {
	case RMSTypeNormal:
		ngs, err := randomMoveSymbols.procNormal(gameProp, msd, plugin, gs)
		if err != nil {
			goutils.Error("RandomMoveSymbols.OnPlayGame:procNormal",
				goutils.Err(err))

			return "", err
		}

		sc2 = ngs
	case RMSTypeReels:
		ngs, err := randomMoveSymbols.procReels(gameProp, msd, plugin, gs, gs2)
		if err != nil {
			goutils.Error("RandomMoveSymbols.OnPlayGame:procReels",
				goutils.Err(err))

			return "", err
		}

		sc2 = ngs
	}

	if sc2 == gs {
		nc := randomMoveSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	randomMoveSymbols.AddScene(gameProp, curpr, sc2, &msd.BasicComponentData)

	randomMoveSymbols.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

	nc := randomMoveSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// NewComponentData -
func (moveSymbol *RandomMoveSymbols) NewComponentData() IComponentData {
	return &RandomMoveSymbolsData{}
}

// OnAsciiGame - outpur to asciigame
func (moveSymbol *RandomMoveSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	msd := icd.(*RandomMoveSymbolsData)

	asciigame.OutputScene("after randomMoveSymbols", pr.Scenes[msd.UsedScenes[0]], mapSymbolColor)

	return nil
}

// OnStats2
func (moveSymbol *RandomMoveSymbols) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	moveSymbol.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	if isOnStepEnd && pr.IsFinish {
		cd := icd.(*RandomMoveSymbolsData)

		s2.ProcStatsIntVal(moveSymbol.GetName(), len(cd.Pos))
	}
}

// NewStats2 -
func (moveSymbol *RandomMoveSymbols) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, []stats2.Option{stats2.OptIntVal})
}

// IsNeedOnStepEndStats2 - 除respin外，如果也有component也需要在stepEnd调用的话，这里需要返回true
func (moveSymbol *RandomMoveSymbols) IsNeedOnStepEndStats2() bool {
	return true
}

func NewRandomMoveSymbols(name string) IComponent {
	return &RandomMoveSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "reels",
// "reelsWeight": "fg_wildreel",
// "targetSymbols": [
// 	"WL2",
// 	"WL3",
// 	"WL5"
// ],
// "ignoreSymbols": [
// 	"SC",
// 	"WL"
// ]
// "targetPositionCollection": "fg-wlpos"

type jsonRandomMoveSymbols struct {
	StrType                  string   `json:"type"`
	TargetSymbols            []string `json:"targetSymbols"`
	IgnoreSymbols            []string `json:"ignoreSymbols"`
	ReelsWeight              string   `json:"reelsWeight"`
	TargetPositionCollection string   `json:"targetPositionCollection"`
}

func (jcfg *jsonRandomMoveSymbols) build() *RandomMoveSymbolsConfig {
	cfg := &RandomMoveSymbolsConfig{
		StrType:                  jcfg.StrType,
		TargetSymbols:            jcfg.TargetSymbols,
		IgnoreSymbols:            jcfg.IgnoreSymbols,
		ReelsWeight:              jcfg.ReelsWeight,
		TargetPositionCollection: jcfg.TargetPositionCollection,
	}

	return cfg
}

func parseRandomMoveSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRandomMoveSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRandomMoveSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRandomMoveSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRandomMoveSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseRandomMoveSymbols:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RandomMoveSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

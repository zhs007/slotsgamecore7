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

const MoveSymbols2TypeName = "moveSymbols2"

type MoveSymbols2Type int

const (
	MS2TypeLeft  MoveSymbols2Type = 0 // left
	MS2TypeRight MoveSymbols2Type = 1 // right
)

func parseMoveSymbols2Type(str string) MoveSymbols2Type {
	if str == "right" {
		return MS2TypeRight
	}

	return MS2TypeLeft
}

type MoveSymbols2Data struct {
	BasicComponentData
	Pos [][]int
}

// OnNewGame -
func (moveSymbols2Data *MoveSymbols2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	moveSymbols2Data.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (moveSymbols2Data *MoveSymbols2Data) OnNewStep() {
	moveSymbols2Data.UsedScenes = nil
	moveSymbols2Data.Pos = nil
}

// Clone
func (moveSymbols2Data *MoveSymbols2Data) Clone() IComponentData {
	target := &MoveSymbols2Data{
		BasicComponentData: moveSymbols2Data.CloneBasicComponentData(),
	}

	target.Pos = make([][]int, len(moveSymbols2Data.Pos))
	for _, arr := range moveSymbols2Data.Pos {
		dstarr := make([]int, len(arr))
		copy(dstarr, arr)
		target.Pos = append(target.Pos, dstarr)
	}

	return target
}

// BuildPBComponentData
func (moveSymbols2Data *MoveSymbols2Data) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.MoveSymbols2Data{
		BasicComponentData: moveSymbols2Data.BuildPBBasicComponentData(),
	}

	num := 0
	for _, arr := range moveSymbols2Data.Pos {
		num += len(arr)
		num++
	}

	pbcd.Pos = make([]int32, 0, num)

	for _, arr := range moveSymbols2Data.Pos {
		for _, s := range arr {
			pbcd.Pos = append(pbcd.Pos, int32(s))
		}

		pbcd.Pos = append(pbcd.Pos, -1)
	}

	return pbcd
}

// GetPos -
func (moveSymbols2Data *MoveSymbols2Data) GetPos() []int {
	num := 0
	for _, arr := range moveSymbols2Data.Pos {
		num += len(arr)
	}

	newpos := make([]int, 0, num)

	for _, arr := range moveSymbols2Data.Pos {
		newpos = append(newpos, arr...)
	}

	return newpos
}

// HasPos -
func (moveSymbols2Data *MoveSymbols2Data) HasPos(x int, y int) bool {
	for _, arr := range moveSymbols2Data.Pos {
		if goutils.IndexOfInt2Slice(arr, x, y, 0) >= 0 {
			return true
		}
	}

	return false
}

// AddPos -
func (moveSymbols2Data *MoveSymbols2Data) AddPos(x int, y int) {
	if len(moveSymbols2Data.Pos) == 0 {
		moveSymbols2Data.Pos = append(moveSymbols2Data.Pos, []int{})
	}

	moveSymbols2Data.Pos[len(moveSymbols2Data.Pos)-1] = append(moveSymbols2Data.Pos[len(moveSymbols2Data.Pos)-1], x, y)
}

// AddPosEx -
func (moveSymbols2Data *MoveSymbols2Data) AddPosEx(x int, y int) {
	if goutils.IndexOfInt2Slice(moveSymbols2Data.Pos[len(moveSymbols2Data.Pos)-1], x, y, 0) < 0 {
		moveSymbols2Data.Pos[len(moveSymbols2Data.Pos)-1] = append(moveSymbols2Data.Pos[len(moveSymbols2Data.Pos)-1], x, y)
	}
}

// newData -
func (moveSymbols2Data *MoveSymbols2Data) newData() {
	moveSymbols2Data.Pos = append(moveSymbols2Data.Pos, []int{})
}

// MoveSymbols2Config - configuration for MoveSymbols2
type MoveSymbols2Config struct {
	BasicComponentConfig   `yaml:",inline" json:",inline"`
	Type                   MoveSymbols2Type `yaml:"-" json:"-"`
	StrType                string           `yaml:"type" json:"type"`
	SrcSymbols             []string         `yaml:"srcSymbols" json:"srcSymbols"`
	SrcSymbolCodes         []int            `yaml:"-" json:"-"`
	SrcPositionCollections []string         `json:"srcPositionCollections"`
	FillSymbol             string           `json:"fillSymbol"`
	FillSymbolCode         int              `json:"-"`
	RemoveSymbol           string           `json:"removeSymbol"`
	RemoveSymbolCode       int              `json:"-"`
	Controllers            []*Award         `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
func (cfg *MoveSymbols2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type MoveSymbols2 struct {
	*BasicComponent `json:"-"`
	Config          *MoveSymbols2Config `json:"config"`
}

// Init -
func (moveSymbol2 *MoveSymbols2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MoveSymbols2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RandomMoveSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MoveSymbols2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return moveSymbol2.InitEx(cfg, pool)
}

// InitEx -
func (moveSymbol2 *MoveSymbols2) InitEx(cfg any, pool *GamePropertyPool) error {
	moveSymbol2.Config = cfg.(*MoveSymbols2Config)
	moveSymbol2.Config.ComponentType = RandomMoveSymbolsTypeName

	moveSymbol2.Config.Type = parseMoveSymbols2Type(moveSymbol2.Config.StrType)

	for _, v := range moveSymbol2.Config.SrcSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("MoveSymbols2.InitEx:SrcSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		moveSymbol2.Config.SrcSymbolCodes = append(moveSymbol2.Config.SrcSymbolCodes, sc)
	}

	if len(moveSymbol2.Config.FillSymbol) > 0 {
		sc, isok := pool.DefaultPaytables.MapSymbols[moveSymbol2.Config.FillSymbol]
		if !isok {
			goutils.Error("MoveSymbols2.InitEx:FillSymbol",
				slog.String("symbol", moveSymbol2.Config.FillSymbol),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		moveSymbol2.Config.FillSymbolCode = sc
	}

	if len(moveSymbol2.Config.RemoveSymbol) > 0 {
		sc, isok := pool.DefaultPaytables.MapSymbols[moveSymbol2.Config.RemoveSymbol]
		if !isok {
			goutils.Error("MoveSymbols2.InitEx:RemoveSymbol",
				slog.String("symbol", moveSymbol2.Config.RemoveSymbol),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		moveSymbol2.Config.RemoveSymbolCode = sc
	}

	for _, award := range moveSymbol2.Config.Controllers {
		award.Init()
	}

	moveSymbol2.onInit(&moveSymbol2.Config.BasicComponentConfig)

	return nil
}

// moveSymbols -
func (moveSymbol2 *MoveSymbols2) moveSymbols(gameProp *GameProperty, cd *MoveSymbols2Data, gs *sgc7game.GameScene,
	gs2 *sgc7game.GameScene, xoff int, yoff int) (*sgc7game.GameScene, error) {

	posSrc := make([]int, 0, gs.Width*gs.Height*2)

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if goutils.IndexOfIntSlice(moveSymbol2.Config.SrcSymbolCodes, s, 0) >= 0 {
				posSrc = append(posSrc, x, y)
			}
		}
	}

	if len(posSrc) == 0 {
		return gs, nil
	}

	ngs := gs

	for i := 0; i < len(posSrc)/2; i++ {
		cd.newData()

		x := posSrc[i*2]
		y := posSrc[i*2+1]

		tx := x + xoff
		ty := y + yoff

		cd.AddPos(x, y)
		cd.AddPos(tx, ty)

		if ngs == gs {
			ngs = gs.CloneEx(gameProp.PoolScene)
		}

		if ngs.IsValidPos(tx, ty) {
			ngs.Arr[tx][ty] = ngs.Arr[x][y]
		}

		ngs.Arr[x][y] = gs2.Arr[x][y]
	}

	return ngs, nil
}

// procSymbols -
func (moveSymbol2 *MoveSymbols2) procSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult,
	cd *MoveSymbols2Data, gs *sgc7game.GameScene) (*sgc7game.GameScene, error) {

	gs2 := gameProp.SceneStack.GetPreTopSceneEx(curpr, prs)

	if moveSymbol2.Config.Type == MS2TypeLeft {
		ngs, err := moveSymbol2.moveSymbols(gameProp, cd, gs, gs2, -1, 0)
		if err != nil {
			goutils.Error("MoveSymbols2.procSymbols:moveSymbols",
				goutils.Err(err))

			return nil, err
		}

		return ngs, nil
	} else if moveSymbol2.Config.Type == MS2TypeRight {
		ngs, err := moveSymbol2.moveSymbols(gameProp, cd, gs, gs2, 1, 0)
		if err != nil {
			goutils.Error("MoveSymbols2.procSymbols:moveSymbols",
				goutils.Err(err))

			return nil, err
		}

		return ngs, nil
	}

	goutils.Error("MoveSymbols2.procSymbols:procMoveSymbols",
		goutils.Err(ErrIvalidComponentConfig))

	return nil, ErrIvalidComponentConfig
}

// movePositionCollection -
func (moveSymbol2 *MoveSymbols2) movePositionCollection(gameProp *GameProperty, cd *MoveSymbols2Data, gs *sgc7game.GameScene,
	ngs *sgc7game.GameScene, cname string, xoff int, yoff int) (*sgc7game.GameScene, error) {

	pcd := gameProp.GetComponentDataWithName(cname)
	pc, isok := gameProp.Components.MapComponents[cname]
	if isok && pcd != nil {
		curpos := slices.Clone(pcd.GetPos())
		if len(curpos) > 0 {
			for i := range len(curpos) / 2 {
				x := curpos[i*2]
				y := curpos[i*2+1]

				cd.newData()

				cd.AddPos(x, y)

				srcSymbol := -99

				if x >= 0 && y >= 0 && x < gs.Width && y < gs.Height {
					srcSymbol = ngs.Arr[x][y]

					if len(moveSymbol2.Config.RemoveSymbol) > 0 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						ngs.Arr[x][y] = moveSymbol2.Config.RemoveSymbolCode
					}
				}

				x += xoff
				y += yoff

				cd.AddPos(x, y)

				if x >= 0 && y >= 0 && x < gs.Width && y < gs.Height {
					if len(moveSymbol2.Config.FillSymbol) > 0 {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						ngs.Arr[x][y] = moveSymbol2.Config.FillSymbolCode
					} else if srcSymbol != -99 {
						ngs.Arr[x][y] = srcSymbol
					}
				}

				curpos[i*2] = x
				curpos[i*2+1] = y
			}

			pc.ClearData(pcd, true)

			for i := range curpos {
				x := curpos[i*2]
				y := curpos[i*2+1]

				if x >= 0 && y >= 0 && x < gs.Width && y < gs.Height {
					pc.AddPos(pcd, x, y)
				}
			}
		}
	}

	return ngs, nil
}

// procPositionCollections -
func (moveSymbol2 *MoveSymbols2) procPositionCollections(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	prs []*sgc7game.PlayResult, cd *MoveSymbols2Data, gs *sgc7game.GameScene) (*sgc7game.GameScene, error) {

	ngs := gs

	if moveSymbol2.Config.Type == MS2TypeLeft {
		for _, v := range moveSymbol2.Config.SrcPositionCollections {
			cngs, err := moveSymbol2.movePositionCollection(gameProp, cd, gs, ngs, v, -1, 0)
			if err != nil {
				goutils.Error("MoveSymbols2.procPositionCollections:movePositionCollection",
					goutils.Err(err))

				return nil, err
			}

			ngs = cngs
		}

		return ngs, nil
	} else if moveSymbol2.Config.Type == MS2TypeRight {
		for _, v := range moveSymbol2.Config.SrcPositionCollections {
			cngs, err := moveSymbol2.movePositionCollection(gameProp, cd, gs, ngs, v, 1, 0)
			if err != nil {
				goutils.Error("MoveSymbols2.procPositionCollections:movePositionCollection",
					goutils.Err(err))

				return nil, err
			}

			ngs = cngs
		}

		return ngs, nil
	}

	goutils.Error("MoveSymbols2.procPositionCollections",
		goutils.Err(ErrIvalidComponentConfig))

	return nil, ErrIvalidComponentConfig
}

// OnProcControllers -
func (moveSymbol2 *MoveSymbols2) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(moveSymbol2.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, moveSymbol2.Config.Controllers, curpr, gp)
	}
}

// playgame
func (moveSymbol2 *MoveSymbols2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	msd := icd.(*MoveSymbols2Data)

	msd.OnNewStep()

	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)
	sc2 := gs

	if len(moveSymbol2.Config.SrcSymbolCodes) > 0 {
		ngs, err := moveSymbol2.procSymbols(gameProp, curpr, prs, msd, gs)
		if err != nil {
			goutils.Error("MoveSymbols2.OnPlayGame:procSymbols",
				goutils.Err(err))

			return "", err
		}

		sc2 = ngs
	} else if len(moveSymbol2.Config.SrcPositionCollections) > 0 {
		ngs, err := moveSymbol2.procPositionCollections(gameProp, curpr, prs, msd, gs)
		if err != nil {
			goutils.Error("MoveSymbols2.OnPlayGame:procPositionCollections",
				goutils.Err(err))

			return "", err
		}

		sc2 = ngs
	} else {
		goutils.Error("MoveSymbols2.OnPlayGame",
			goutils.Err(ErrIvalidComponentConfig))

		return "", ErrIvalidComponentConfig
	}

	if sc2 == gs {
		nc := moveSymbol2.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	moveSymbol2.AddScene(gameProp, curpr, sc2, &msd.BasicComponentData)

	moveSymbol2.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

	nc := moveSymbol2.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// NewComponentData -
func (moveSymbol2 *MoveSymbols2) NewComponentData() IComponentData {
	return &MoveSymbols2Data{}
}

// OnAsciiGame - outpur to asciigame
func (moveSymbol2 *MoveSymbols2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	msd := icd.(*MoveSymbols2Data)

	asciigame.OutputScene("after moveSymbols2", pr.Scenes[msd.UsedScenes[0]], mapSymbolColor)

	return nil
}

func NewMoveSymbols2(name string) IComponent {
	return &MoveSymbols2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "left",
// "srcPositionCollections": [
// 	"bg-pos-mw"
// ],
// "fillSymbol": "MW"

type jsonMoveSymbols2 struct {
	StrType                string   `json:"type"`
	SrcSymbols             []string `json:"srcSymbols"`
	SrcPositionCollections []string `json:"srcPositionCollections"`
	FillSymbol             string   `json:"fillSymbol"`
	RemoveSymbol           string   `json:"removeSymbol"`
}

func (jcfg *jsonMoveSymbols2) build() *MoveSymbols2Config {
	cfg := &MoveSymbols2Config{
		StrType:                strings.ToLower(jcfg.StrType),
		SrcSymbols:             slices.Clone(jcfg.SrcSymbols),
		SrcPositionCollections: slices.Clone(jcfg.SrcPositionCollections),
		FillSymbol:             jcfg.FillSymbol,
		RemoveSymbol:           jcfg.RemoveSymbol,
	}

	return cfg
}

func parseMoveSymbols2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseMoveSymbols2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseMoveSymbols2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonMoveSymbols2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseMoveSymbols2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseMoveSymbols2:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: MoveSymbols2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

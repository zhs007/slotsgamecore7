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

const MoveSymbolTypeName = "moveSymbol"

const (
	SelectSymbolR2L = "selectSymbolR2L"
	SelectSymbolL2R = "selectSymbolL2R"
	SelectWithXY    = "selectWithXY"
)

type MoveSymbolData struct {
	BasicComponentData
	Pos [][]int
}

// OnNewGame -
func (moveSymbolData *MoveSymbolData) OnNewGame(gameProp *GameProperty, component IComponent) {
	moveSymbolData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (moveSymbolData *MoveSymbolData) OnNewStep() {
	moveSymbolData.UsedScenes = nil
	moveSymbolData.Pos = nil
}

// Clone
func (moveSymbolData *MoveSymbolData) Clone() IComponentData {
	target := &MoveSymbolData{
		BasicComponentData: moveSymbolData.CloneBasicComponentData(),
	}

	target.Pos = make([][]int, len(moveSymbolData.Pos))
	for _, arr := range moveSymbolData.Pos {
		dstarr := make([]int, len(arr))
		copy(dstarr, arr)
		target.Pos = append(target.Pos, dstarr)
	}

	return target
}

// BuildPBComponentData
func (moveSymbolData *MoveSymbolData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.MoveSymbolData{
		BasicComponentData: moveSymbolData.BuildPBBasicComponentData(),
	}

	num := 0
	for _, arr := range moveSymbolData.Pos {
		num += len(arr)
		num++
	}

	pbcd.Pos = make([]int32, 0, num)

	for _, arr := range moveSymbolData.Pos {
		for _, s := range arr {
			pbcd.Pos = append(pbcd.Pos, int32(s))
		}

		pbcd.Pos = append(pbcd.Pos, -1)
	}

	return pbcd
}

// GetPos -
func (moveSymbolData *MoveSymbolData) GetPos() []int {
	num := 0
	for _, arr := range moveSymbolData.Pos {
		num += len(arr)
	}

	newpos := make([]int, 0, num)

	for _, arr := range moveSymbolData.Pos {
		newpos = append(newpos, arr...)
	}

	return newpos
}

// HasPos -
func (moveSymbolData *MoveSymbolData) HasPos(x int, y int) bool {
	for _, arr := range moveSymbolData.Pos {
		if goutils.IndexOfInt2Slice(arr, x, y, 0) >= 0 {
			return true
		}
	}

	return false
}

// AddPos -
func (moveSymbolData *MoveSymbolData) AddPos(x int, y int) {
	if len(moveSymbolData.Pos) == 0 {
		moveSymbolData.Pos = append(moveSymbolData.Pos, []int{})
	}

	moveSymbolData.Pos[len(moveSymbolData.Pos)-1] = append(moveSymbolData.Pos[len(moveSymbolData.Pos)-1], x, y)
}

// AddPosEx -
func (moveSymbolData *MoveSymbolData) AddPosEx(x int, y int) {
	if goutils.IndexOfInt2Slice(moveSymbolData.Pos[len(moveSymbolData.Pos)-1], x, y, 0) < 0 {
		moveSymbolData.Pos[len(moveSymbolData.Pos)-1] = append(moveSymbolData.Pos[len(moveSymbolData.Pos)-1], x, y)
	}
}

// newData -
func (moveSymbolData *MoveSymbolData) newData() {
	moveSymbolData.Pos = append(moveSymbolData.Pos, []int{})
}

type SelectPosData struct {
	Type       string `yaml:"type" json:"type"`
	X          int    `yaml:"x" json:"x"`
	Y          int    `yaml:"y" json:"y"`
	Symbol     string `yaml:"symbol" json:"symbol"`
	SymbolCode int    `yaml:"-" json:"-"`
}

func (spd *SelectPosData) Select(gs *sgc7game.GameScene) (bool, int, int) {
	if spd.Type == SelectWithXY {
		return true, spd.X, spd.Y
	} else if spd.Type == SelectSymbolR2L {
		for x := gs.Width - 1; x >= 0; x-- {
			if gs.Arr[x][spd.Y] == spd.SymbolCode {
				return true, x, spd.Y
			}
		}
	} else if spd.Type == SelectSymbolL2R {
		for x := 0; x < gs.Width; x++ {
			if gs.Arr[x][spd.Y] == spd.SymbolCode {
				return true, x, spd.Y
			}
		}
	}

	return false, 0, 0
}

const (
	MoveTypeXY = "xy"
	MoveTypeYX = "yx"
)

type MoveData struct {
	Src              *SelectPosData `yaml:"src" json:"src"`
	Target           *SelectPosData `yaml:"target" json:"target"`
	MoveType         string         `yaml:"moveType" json:"moveType"`
	TargetSymbol     string         `yaml:"targetSymbol" json:"targetSymbol"`
	TargetSymbolCode int            `yaml:"-" json:"-"`
	OverrideSrc      bool           `yaml:"overrideSrc" json:"overrideSrc"`
	OverrideTarget   bool           `yaml:"overrideTarget" json:"overrideTarget"`
	OverridePath     bool           `yaml:"overridePath" json:"overridePath"`
}

func (md *MoveData) moveX(gs *sgc7game.GameScene, sx, tx int, y int, symbolCode int, msd *MoveSymbolData) {
	if tx > sx {
		for x := sx + 1; x < tx; x++ {
			gs.Arr[x][y] = symbolCode

			msd.AddPosEx(x, y)
		}
	} else if tx < sx {
		for x := sx - 1; x > tx; x-- {
			gs.Arr[x][y] = symbolCode

			msd.AddPosEx(x, y)
		}
	}
}

func (md *MoveData) moveY(gs *sgc7game.GameScene, sy, ty int, x int, symbolCode int, msd *MoveSymbolData) {
	if ty > sy {
		for y := sy + 1; y < ty; y++ {
			gs.Arr[x][y] = symbolCode

			msd.AddPosEx(x, y)
		}
	} else if ty < sy {
		for y := sy - 1; y > ty; y-- {
			gs.Arr[x][y] = symbolCode

			msd.AddPosEx(x, y)
		}
	}
}

func (md *MoveData) Move(gs *sgc7game.GameScene, sx, sy, tx, ty int, symbolCode int, msd *MoveSymbolData) {
	if md.OverrideSrc {
		gs.Arr[sx][sy] = symbolCode

		msd.AddPosEx(sx, sy)
	}

	if md.OverrideTarget {
		gs.Arr[tx][ty] = symbolCode

		msd.AddPosEx(tx, ty)
	}

	if !md.OverridePath {
		return
	}

	if md.MoveType == MoveTypeXY {
		md.moveX(gs, sx, tx, sy, symbolCode, msd) // sx,sy -> tx,sy

		if sy != ty {
			gs.Arr[tx][sy] = symbolCode

			msd.AddPosEx(tx, sy)

			md.moveY(gs, sy, ty, tx, symbolCode, msd) // tx,sy -> tx,ty
		}
	} else if md.MoveType == MoveTypeYX {
		md.moveY(gs, sy, ty, sx, symbolCode, msd) // sx,sy -> sx,ty

		if sx != tx {
			gs.Arr[sx][ty] = symbolCode

			msd.AddPosEx(sx, ty)

			md.moveX(gs, sx, tx, ty, symbolCode, msd) // sx,sy -> sx,ty
		}
	}
}

// MoveSymbolConfig - configuration for MoveSymbol
type MoveSymbolConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MoveData             []*MoveData `yaml:"moveData" json:"moveData"`
}

// SetLinkComponent
func (cfg *MoveSymbolConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type MoveSymbol struct {
	*BasicComponent `json:"-"`
	Config          *MoveSymbolConfig `json:"config"`
}

// Init -
func (moveSymbol *MoveSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MoveSymbol.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &MoveSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MoveSymbol.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return moveSymbol.InitEx(cfg, pool)
}

// InitEx -
func (moveSymbol *MoveSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	moveSymbol.Config = cfg.(*MoveSymbolConfig)
	moveSymbol.Config.ComponentType = MoveSymbolTypeName

	for _, v := range moveSymbol.Config.MoveData {
		if v.Src.Type != SelectWithXY {
			sc, isok := pool.DefaultPaytables.MapSymbols[v.Src.Symbol]
			if !isok {
				goutils.Error("ReplaceReel.InitEx:Src.Symbol",
					slog.String("symbol", v.Src.Symbol),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}

			v.Src.SymbolCode = sc
		} else {
			v.Src.SymbolCode = -1
		}

		if v.Target.Type != SelectWithXY {
			sc, isok := pool.DefaultPaytables.MapSymbols[v.Target.Symbol]
			if !isok {
				goutils.Error("ReplaceReel.InitEx:Target.Symbol",
					slog.String("symbol", v.Target.Symbol),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}

			v.Target.SymbolCode = sc
		} else {
			v.Target.SymbolCode = -1
		}

		sc, isok := pool.DefaultPaytables.MapSymbols[v.TargetSymbol]
		if isok {
			v.TargetSymbolCode = sc
		} else {
			v.TargetSymbolCode = -1
		}
	}

	moveSymbol.onInit(&moveSymbol.Config.BasicComponentConfig)

	return nil
}

// playgame
func (moveSymbol *MoveSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// moveSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	msd := icd.(*MoveSymbolData)

	msd.OnNewStep()

	gs := moveSymbol.GetTargetScene3(gameProp, curpr, prs, 0)

	sc2 := gs

	for _, v := range moveSymbol.Config.MoveData {
		srcok, srcx, srcy := v.Src.Select(sc2)
		if !srcok {
			continue
		}

		targetok, targetx, targety := v.Target.Select(sc2)
		if !targetok {
			continue
		}

		if targetx == srcx && targety == srcy {
			continue
		}

		msd.newData()

		symbolCode := v.TargetSymbolCode
		if symbolCode == -1 {
			symbolCode = gs.Arr[srcx][srcy]
		}

		if srcx == targetx && srcy == targety {
			if v.OverrideSrc {
				gs.Arr[srcx][srcy] = symbolCode

				msd.AddPosEx(srcx, srcy)
			}

			if v.OverrideTarget {
				gs.Arr[targetx][targety] = symbolCode

				msd.AddPosEx(targetx, targety)
			}

			continue
		}

		if sc2 == gs {
			sc2 = gs.CloneEx(gameProp.PoolScene)
		}

		v.Move(sc2, srcx, srcy, targetx, targety, symbolCode, msd)
	}

	if sc2 == gs {
		nc := moveSymbol.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	moveSymbol.AddScene(gameProp, curpr, sc2, &msd.BasicComponentData)

	nc := moveSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// NewComponentData -
func (moveSymbol *MoveSymbol) NewComponentData() IComponentData {
	return &MoveSymbolData{}
}

// OnAsciiGame - outpur to asciigame
func (moveSymbol *MoveSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	msd := icd.(*MoveSymbolData)

	asciigame.OutputScene("after moveSymbol", pr.Scenes[msd.UsedScenes[0]], mapSymbolColor)

	return nil
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

func NewMoveSymbol(name string) IComponent {
	return &MoveSymbol{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"isExpandReel": "false",
//		"moveData": [
//			{
//				"src": {
//					"type": "selectWithXY",
//					"Y": 1,
//					"X": 1
//				},
//				"target": {
//					"type": "selectSymbolR2L",
//					"Y": 1,
//					"Symbol": "SC"
//				},
//				"moveType": "xy",
//				"targetSymbol": "SC",
//				"overrideSrc": "false",
//				"overrideTarget": "false",
//				"overridePath": "true",
//				"name": "moveData 1"
//			},
//			{
//				"src": {
//					"type": "selectWithXY",
//					"Y": 2,
//					"X": 1
//				},
//				"target": {
//					"type": "selectSymbolR2L",
//					"Y": 2,
//					"Symbol": "SC"
//				},
//				"moveType": "xy",
//				"targetSymbol": "SC",
//				"overrideSrc": "false",
//				"overrideTarget": "false",
//				"overridePath": "true",
//				"name": "moveData 2"
//			},
//			{
//				"src": {
//					"type": "selectWithXY",
//					"Y": 3,
//					"X": 1
//				},
//				"target": {
//					"type": "selectSymbolR2L",
//					"Y": 3,
//					"Symbol": "SC"
//				},
//				"moveType": "xy",
//				"targetSymbol": "SC",
//				"overrideSrc": "false",
//				"overrideTarget": "false",
//				"overridePath": "true",
//				"name": "moveData 3"
//			}
//		]
//	},
type jsonMoveData struct {
	Src            *SelectPosData `json:"src"`
	Target         *SelectPosData `json:"target"`
	MoveType       string         `json:"moveType"`
	TargetSymbol   string         `json:"targetSymbol"`
	OverrideSrc    string         `json:"overrideSrc"`
	OverrideTarget string         `json:"overrideTarget"`
	OverridePath   string         `json:"overridePath"`
}
type jsonMoveSymbol struct {
	MoveData []*jsonMoveData `json:"moveData"`
}

func (jms *jsonMoveSymbol) build() *MoveSymbolConfig {
	cfg := &MoveSymbolConfig{}

	for _, v := range jms.MoveData {
		cmd := &MoveData{
			Src:            v.Src,
			Target:         v.Target,
			MoveType:       v.MoveType,
			TargetSymbol:   v.TargetSymbol,
			OverrideSrc:    v.OverrideSrc == "true",
			OverrideTarget: v.OverrideTarget == "true",
			OverridePath:   v.OverridePath == "true",
		}

		if cmd.Src.X > 0 {
			cmd.Src.X--
		}

		if cmd.Src.Y > 0 {
			cmd.Src.Y--
		}

		if cmd.Target.X > 0 {
			cmd.Target.X--
		}

		if cmd.Target.Y > 0 {
			cmd.Target.Y--
		}

		cfg.MoveData = append(cfg.MoveData, cmd)
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseMoveSymbol(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseMoveSymbol:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseMoveSymbol:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonMoveSymbol{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseMoveSymbol:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: MoveSymbolTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}

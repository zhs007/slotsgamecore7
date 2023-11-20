package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const MoveSymbolTypeName = "moveSymbol"

const (
	SelectSymbolR2L = "selectSymbolR2L"
	SelectSymbolL2R = "selectSymbolL2R"
	SelectWithXY    = "selectWithXY"
)

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

func (md *MoveData) moveX(gs *sgc7game.GameScene, sx, tx int, y int) {
	if tx > sx {
		for x := sx + 1; x < tx; x++ {
			gs.Arr[x][y] = md.TargetSymbolCode
		}
	} else if tx < sx {
		for x := sx - 1; x > tx; x-- {
			gs.Arr[x][y] = md.TargetSymbolCode
		}
	}
}

func (md *MoveData) moveY(gs *sgc7game.GameScene, sy, ty int, x int) {
	if ty > sy {
		for y := sy + 1; y < ty; y++ {
			gs.Arr[x][y] = md.TargetSymbolCode
		}
	} else if ty < sy {
		for y := sy - 1; y > ty; y-- {
			gs.Arr[x][y] = md.TargetSymbolCode
		}
	}
}

func (md *MoveData) Move(gs *sgc7game.GameScene, sx, sy, tx, ty int) {
	if md.OverrideSrc {
		gs.Arr[sx][sy] = md.TargetSymbolCode
	}

	if md.OverrideTarget {
		gs.Arr[tx][ty] = md.TargetSymbolCode
	}

	if md.MoveType == MoveTypeXY {
		md.moveX(gs, sx, tx, sy) // sx,sy -> tx,sy

		if sy != ty {
			gs.Arr[tx][sy] = md.TargetSymbolCode

			md.moveY(gs, sy, ty, tx) // tx,sy -> tx,ty
		}
	} else if md.MoveType == MoveTypeYX {
		md.moveY(gs, sy, ty, sx) // sx,sy -> sx,ty

		if sx != tx {
			gs.Arr[sx][ty] = md.TargetSymbolCode

			md.moveX(gs, sx, tx, ty) // sx,sy -> sx,ty
		}
	}
}

// MoveSymbolConfig - configuration for MoveSymbol
type MoveSymbolConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MoveData             []*MoveData `yaml:"moveData" json:"moveData"`
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
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MoveSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MoveSymbol.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return moveSymbol.InitEx(cfg, pool)
}

// InitEx -
func (moveSymbol *MoveSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	moveSymbol.Config = cfg.(*MoveSymbolConfig)
	moveSymbol.Config.ComponentType = MoveSymbolTypeName

	for _, v := range moveSymbol.Config.MoveData {
		sc, isok := pool.DefaultPaytables.MapSymbols[v.Src.Symbol]
		if !isok {
			goutils.Error("ReplaceReel.InitEx:Src.Symbol",
				zap.String("symbol", v.Src.Symbol),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		v.Src.SymbolCode = sc

		sc, isok = pool.DefaultPaytables.MapSymbols[v.Target.Symbol]
		if !isok {
			goutils.Error("ReplaceReel.InitEx:Target.Symbol",
				zap.String("symbol", v.Target.Symbol),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		v.Target.SymbolCode = sc

		sc, isok = pool.DefaultPaytables.MapSymbols[v.TargetSymbol]
		if !isok {
			goutils.Error("ReplaceReel.InitEx:TargetSymbol",
				zap.String("symbol", v.TargetSymbol),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		v.TargetSymbolCode = sc
	}

	moveSymbol.onInit(&moveSymbol.Config.BasicComponentConfig)

	return nil
}

// playgame
func (moveSymbol *MoveSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	moveSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[moveSymbol.Name].(*BasicComponentData)

	gs := moveSymbol.GetTargetScene(gameProp, curpr, cd, "")

	sc2 := gs.CloneEx(gameProp.PoolScene)

	for _, v := range moveSymbol.Config.MoveData {
		srcok, srcx, srcy := v.Src.Select(sc2)
		if !srcok {
			continue
		}

		targetok, targetx, targety := v.Target.Select(sc2)
		if !targetok {
			continue
		}

		if srcx == targetx && srcy == targety {
			continue
		}

		v.Move(sc2, srcx, srcy, targetx, targety)
	}

	moveSymbol.AddScene(gameProp, curpr, sc2, cd)

	moveSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (moveSymbol *MoveSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[moveSymbol.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after moveSymbol", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (moveSymbol *MoveSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewMoveSymbol(name string) IComponent {
	return &MoveSymbol{
		BasicComponent: NewBasicComponent(name),
	}
}

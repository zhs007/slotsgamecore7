package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

const (
	GamePropWidth        = 1
	GamePropHeight       = 2
	GamePropCurPaytables = 3
	GamePropCurReels     = 4
	GamePropCurLineData  = 5

	GamePropCurMystery = 1000
)

var MapProperty map[string]int

func String2Property(str string) (int, error) {
	v, isok := MapProperty[str]
	if isok {
		return v, nil
	}

	goutils.Error("String2Property",
		zap.String("str", str),
		zap.Error(ErrInvalidGamePropertyString))

	return 0, ErrInvalidGamePropertyString
}

type GameProperty struct {
	Config       *Config
	MapVal       map[int]int
	MapStrVal    map[int]string
	CurPaytables *sgc7game.PayTables
	CurLineData  *sgc7game.LineData
	CurReels     *sgc7game.ReelsData
}

func (gameProp *GameProperty) SetVal(prop int, val int) error {
	gameProp.MapVal[prop] = val

	return nil
}

func (gameProp *GameProperty) SetStrVal(prop int, val string) error {
	if prop == GamePropCurMystery {
		v, isok := gameProp.CurPaytables.MapSymbols[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurMystery",
				zap.String("val", val),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		gameProp.MapVal[prop] = v
	} else if prop == GamePropCurPaytables {
		v, isok := gameProp.Config.MapPaytables[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurPaytables",
				zap.String("val", val),
				zap.Error(ErrInvalidPaytables))

			return ErrInvalidPaytables
		}

		gameProp.CurPaytables = v
	}

	gameProp.MapStrVal[prop] = val

	return nil
}

func InitGameProperty(cfgfn string) (*GameProperty, error) {
	cfg, err := LoadConfig(cfgfn)
	if err != nil {
		goutils.Error("InitGameProperty:LoadConfig",
			zap.String("cfgfn", cfgfn),
			zap.Error(err))

		return nil, err
	}

	gameProp := &GameProperty{
		Config:    cfg,
		MapVal:    make(map[int]int),
		MapStrVal: make(map[int]string),
	}

	gameProp.SetStrVal(GamePropCurPaytables, "main")
	gameProp.SetVal(GamePropWidth, cfg.Width)
	gameProp.SetVal(GamePropHeight, cfg.Height)

	return gameProp, nil
}

func init() {
	MapProperty = make(map[string]int)

	MapProperty["width"] = GamePropWidth
	MapProperty["height"] = GamePropHeight
	MapProperty["paytables"] = GamePropCurPaytables
	MapProperty["reels"] = GamePropCurReels
	MapProperty["linedata"] = GamePropCurLineData

	MapProperty["curMystery"] = GamePropCurMystery
}

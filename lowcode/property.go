package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

const (
	GamePropWidth  = 1
	GamePropHeight = 2

	GamePropCurMystery = 1000
)

type GameProperty struct {
	Config       *Config
	MapVal       map[int]int
	CurPaytables *sgc7game.PayTables
	CurLineData  *sgc7game.LineData
	CurReels     *sgc7game.ReelsData
}

func InitGameProperty(cfgfn string) (*GameProperty, error) {
	cfg, err := LoadConfig(cfgfn)
	if err != nil {
		goutils.Error("InitGameProperty:LoadConfig",
			zap.String("cfgfn", cfgfn),
			zap.Error(err))

		return nil, err
	}

	prop := &GameProperty{
		Config:       cfg,
		MapVal:       make(map[int]int),
		CurPaytables: cfg.MapPaytables["main"],
	}

	prop.MapVal[GamePropWidth] = cfg.Width
	prop.MapVal[GamePropHeight] = cfg.Height

	return prop, nil
}

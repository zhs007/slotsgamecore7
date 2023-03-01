package lowcode

import (
	"github.com/fatih/color"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
)

const (
	GamePropWidth        = 1
	GamePropHeight       = 2
	GamePropCurPaytables = 3
	GamePropCurReels     = 4
	GamePropCurLineData  = 5

	GamePropTriggerFG = 100
	GamePropFGNum     = 101

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
	Config           *Config
	MapVals          map[int]int
	MapStrVals       map[int]string
	CurPaytables     *sgc7game.PayTables
	CurLineData      *sgc7game.LineData
	CurReels         *sgc7game.ReelsData
	MapIntValWeights map[string]*sgc7game.ValWeights2
	Plugin           sgc7plugin.IPlugin
	SymbolsViewer    *SymbolsViewer
	MapSymbolColor   *asciigame.SymbolColorMap
	MapScenes        map[string]int
}

func (gameProp *GameProperty) OnNewStep() error {
	gameProp.MapScenes = make(map[string]int)

	return nil
}

func (gameProp *GameProperty) TagScene(pr *sgc7game.PlayResult, tag string, sceneIndex int) {
	gameProp.MapScenes[tag] = sceneIndex
}

func (gameProp *GameProperty) GetScene(pr *sgc7game.PlayResult, tag string) *sgc7game.GameScene {
	si, isok := gameProp.MapScenes[tag]
	if !isok {
		return pr.Scenes[len(pr.Scenes)-1]
	}

	return pr.Scenes[si]
}

func (gameProp *GameProperty) TriggerFGWithWeights(fn string) error {
	vw2, isok := gameProp.MapIntValWeights[fn]
	if !isok {
		curvw2, err := sgc7game.LoadValWeights2FromExcel(fn, "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("GameProperty.TriggerFGWithWeights:LoadValWeights2FromExcel",
				zap.String("fn", fn),
				zap.Error(err))

			return err
		}

		gameProp.MapIntValWeights[fn] = curvw2

		vw2 = curvw2
	}

	val, err := vw2.RandVal(gameProp.Plugin)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:RandVal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	if val.Int() > 0 {
		gameProp.SetVal(GamePropTriggerFG, 1)
		gameProp.SetVal(GamePropFGNum, val.Int())
	}

	return nil
}

func (gameProp *GameProperty) SetVal(prop int, val int) error {
	if prop == GamePropCurMystery {
		str := gameProp.CurPaytables.GetStringFromInt(val)

		gameProp.MapStrVals[prop] = str
	}

	gameProp.MapVals[prop] = val

	return nil
}

func (gameProp *GameProperty) GetVal(prop int) int {
	return gameProp.MapVals[prop]
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

		gameProp.MapVals[prop] = v
	} else if prop == GamePropCurPaytables {
		v, isok := gameProp.Config.MapPaytables[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurPaytables",
				zap.String("val", val),
				zap.Error(ErrInvalidPaytables))

			return ErrInvalidPaytables
		}

		gameProp.CurPaytables = v
	} else if prop == GamePropCurLineData {
		v, isok := gameProp.Config.MapLinedate[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurLineData",
				zap.String("val", val),
				zap.Error(ErrInvalidPaytables))

			return ErrInvalidPaytables
		}

		gameProp.CurLineData = v
	}

	gameProp.MapStrVals[prop] = val

	return nil
}

func (gameProp *GameProperty) GetStrVal(prop int) string {
	return gameProp.MapStrVals[prop]
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
		Config:           cfg,
		MapVals:          make(map[int]int),
		MapStrVals:       make(map[int]string),
		MapIntValWeights: make(map[string]*sgc7game.ValWeights2),
	}

	gameProp.SetStrVal(GamePropCurPaytables, cfg.DefaultPaytables)
	gameProp.SetStrVal(GamePropCurLineData, cfg.DefaultLinedata)
	gameProp.SetVal(GamePropWidth, cfg.Width)
	gameProp.SetVal(GamePropHeight, cfg.Height)

	sv, err := LoadSymbolsViewer(cfg.SymbolsViewer)
	if err != nil {
		goutils.Error("InitGameProperty:LoadSymbolsViewer",
			zap.String("fn", cfg.SymbolsViewer),
			zap.Error(err))

		return nil, err
	}

	gameProp.SymbolsViewer = sv
	gameProp.MapSymbolColor = asciigame.NewSymbolColorMap(gameProp.CurPaytables)
	wColor := color.New(color.BgRed, color.FgHiWhite)
	hColor := color.New(color.BgBlue, color.FgHiWhite)
	mColor := color.New(color.BgGreen, color.FgHiWhite)
	sColor := color.New(color.BgMagenta, color.FgHiWhite)
	for k, v := range sv.MapSymbols {
		if v.Color == "wild" {
			gameProp.MapSymbolColor.AddSymbolColor(k, wColor)
		} else if v.Color == "high" {
			gameProp.MapSymbolColor.AddSymbolColor(k, hColor)
		} else if v.Color == "medium" {
			gameProp.MapSymbolColor.AddSymbolColor(k, mColor)
		} else if v.Color == "scatter" {
			gameProp.MapSymbolColor.AddSymbolColor(k, sColor)
		}
	}

	gameProp.MapSymbolColor.OnGetSymbolString = func(s int) string {
		return gameProp.SymbolsViewer.MapSymbols[s].Output
	}

	return gameProp, nil
}

func init() {
	MapProperty = make(map[string]int)

	MapProperty["width"] = GamePropWidth
	MapProperty["height"] = GamePropHeight
	MapProperty["paytables"] = GamePropCurPaytables
	MapProperty["reels"] = GamePropCurReels
	MapProperty["linedata"] = GamePropCurLineData

	MapProperty["triggerFG"] = GamePropTriggerFG
	MapProperty["FGNum"] = GamePropFGNum

	MapProperty["curMystery"] = GamePropCurMystery
}

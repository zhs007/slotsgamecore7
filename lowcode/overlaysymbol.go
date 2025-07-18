package lowcode

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const OverlaySymbolTypeName = "overlaySymbol"

type OverlaySymbolData struct {
	BasicComponentData
	CurLevel int
}

// OnNewGame -
func (overlaySymbolData *OverlaySymbolData) OnNewGame(gameProp *GameProperty, component IComponent) {
	overlaySymbolData.BasicComponentData.OnNewGame(gameProp, component)

	overlaySymbol := component.(*OverlaySymbol)

	overlaySymbolData.CurLevel = overlaySymbol.Config.DefaultLevel
}

// // onNewStep -
// func (overlaySymbolData *OverlaySymbolData) onNewStep(gameProp *GameProperty, component IComponent) {
// 	// overlaySymbolData.BasicComponentData.OnNewStep(gameProp, component)

// 	overlaySymbol := component.(*OverlaySymbol)

// 	if overlaySymbol.Config.Collector != "" {
// 		collectorData, isok := gameProp.GetCurComponentDataWithName(overlaySymbol.Config.Collector).(*CollectorData)
// 		if !isok {
// 			goutils.Error("OverlaySymbolData.OnNewStep:GetCurComponentDataWithName",
// 				slog.String("collector", overlaySymbol.Config.Collector))
// 		} else {
// 			overlaySymbolData.CurLevel = collectorData.Val
// 		}
// 	}
// }

// Clone
func (overlaySymbolData *OverlaySymbolData) Clone() IComponentData {
	target := &OverlaySymbolData{
		BasicComponentData: overlaySymbolData.CloneBasicComponentData(),
		CurLevel:           overlaySymbolData.CurLevel,
	}

	return target
}

// BuildPBComponentData
func (overlaySymbolData *OverlaySymbolData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.OverlaySymbolData{
		BasicComponentData: overlaySymbolData.BuildPBBasicComponentData(),
	}

	pbcd.CurLevel = int32(overlaySymbolData.CurLevel)

	return pbcd
}

// OverlaySymbolConfig - configuration for OverlaySymbol feature
type OverlaySymbolConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbol               string `yaml:"symbol" json:"symbol"`
	MapPosition          string `yaml:"mapPosition" json:"mapPosition"`
	DefaultLevel         int    `yaml:"defaultLevel" json:"defaultLevel"`
	Collector            string `yaml:"collector" json:"collector"`
}

type OverlaySymbol struct {
	*BasicComponent `json:"-"`
	Config          *OverlaySymbolConfig  `json:"config"`
	SymbolCode      int                   `json:"-"`
	MapPosition     *sgc7game.ValMapping2 `json:"-"`
}

// // OnNewGame -
// func (overlaySymbol *OverlaySymbol) OnNewGame(gameProp *GameProperty) error {
// 	osd := gameProp.MapComponentData[overlaySymbol.Name].(*OverlaySymbolData)

// 	osd.OnNewGame()

// 	osd.CurLevel = overlaySymbol.Config.DefaultLevel

// 	return nil
// }

// // OnNewStep -
// func (overlaySymbol *OverlaySymbol) OnNewStep(gameProp *GameProperty) error {
// 	overlaySymbol.BasicComponent.OnNewStep(gameProp)

// 	cd := gameProp.MapComponentData[overlaySymbol.Name].(*OverlaySymbolData)

// 	if overlaySymbol.Config.Collector != "" {
// 		collectorData, isok := gameProp.MapComponentData[overlaySymbol.Config.Collector].(*CollectorData)
// 		if isok {
// 			cd.CurLevel = collectorData.Val
// 		}
// 	}

// 	return nil
// }

// Init -
func (overlaySymbol *OverlaySymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("OverlaySymbol.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &OverlaySymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("OverlaySymbol.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return overlaySymbol.InitEx(cfg, pool)
}

// InitEx -
func (overlaySymbol *OverlaySymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	overlaySymbol.Config = cfg.(*OverlaySymbolConfig)
	overlaySymbol.Config.ComponentType = OverlaySymbolTypeName

	if overlaySymbol.Config.MapPosition != "" {
		vm2, err := sgc7game.LoadValMapping2FromExcel(pool.Config.GetPath(overlaySymbol.Config.MapPosition, overlaySymbol.Config.UseFileMapping), "index", "value", sgc7game.NewIntArrVal[int])
		if err != nil {
			goutils.Error("OverlaySymbol.Init:LoadValMapping2FromExcel",
				slog.String("valmapping", overlaySymbol.Config.MapPosition),
				goutils.Err(err))

			return err
		}

		overlaySymbol.MapPosition = vm2
	}

	overlaySymbol.SymbolCode = pool.DefaultPaytables.MapSymbols[overlaySymbol.Config.Symbol]

	overlaySymbol.onInit(&overlaySymbol.Config.BasicComponentConfig)

	return nil
}

// playgame
func (overlaySymbol *OverlaySymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// overlaySymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	osd := cd.(*OverlaySymbolData)

	_, hasVal := overlaySymbol.MapPosition.MapVals[osd.CurLevel]
	if hasVal {
		gs := overlaySymbol.GetTargetScene3(gameProp, curpr, prs, 0)

		// cgs := gs.Clone()
		cgs := gs.CloneEx(gameProp.PoolScene)

		for i := 0; i <= osd.CurLevel; i++ {
			pos, isok := overlaySymbol.MapPosition.MapVals[i]
			if isok {
				cgs.Arr[pos.GetInt(0)][pos.GetInt(1)] = overlaySymbol.SymbolCode
			}
		}

		overlaySymbol.AddScene(gameProp, curpr, cgs, &osd.BasicComponentData)
	} else {
		overlaySymbol.GetTargetScene3(gameProp, curpr, prs, 0)

		// overlaySymbol.ReTagScene(gameProp, curpr, osd.TargetSceneIndex, &osd.BasicComponentData)
	}

	nc := overlaySymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (overlaySymbol *OverlaySymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	osd := cd.(*OverlaySymbolData)

	if len(osd.UsedScenes) > 0 {
		asciigame.OutputScene("The symbols after the symbol overlay", pr.Scenes[osd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (overlaySymbol *OverlaySymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // OnStatsWithPB -
// func (overlaySymbol *OverlaySymbol) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.OverlaySymbolData)
// 	if !isok {
// 		goutils.Error("OverlaySymbol.OnStatsWithPB",
// 			goutils.Err(ErrInvalidProto))

// 		return 0, ErrInvalidProto
// 	}

// 	return overlaySymbol.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// NewComponentData -
func (overlaySymbol *OverlaySymbol) NewComponentData() IComponentData {
	return &OverlaySymbolData{}
}

// EachUsedResults -
func (overlaySymbol *OverlaySymbol) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.OverlaySymbolData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("OverlaySymbol.EachUsedResults:UnmarshalTo",
			goutils.Err(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewOverlaySymbol(name string) IComponent {
	return &OverlaySymbol{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

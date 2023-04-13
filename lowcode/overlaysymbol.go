package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

type OverlaySymbolData struct {
	BasicComponentData
	CurLevel int
}

// OnNewGame -
func (overlaySymbolData *OverlaySymbolData) OnNewGame() {
	overlaySymbolData.BasicComponentData.OnNewGame()
}

// OnNewGame -
func (overlaySymbolData *OverlaySymbolData) OnNewStep() {
	overlaySymbolData.BasicComponentData.OnNewStep()

	overlaySymbolData.CurLevel = 0
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
	BasicComponentConfig `yaml:",inline"`
	Symbol               string `yaml:"symbol"`
	MapPosition          string `yaml:"mapPosition"`
	DefaultLevel         int    `yaml:"defaultLevel"`
}

type OverlaySymbol struct {
	*BasicComponent
	Config      *OverlaySymbolConfig
	SymbolCode  int
	MapPosition *sgc7game.ValMapping2
}

// OnNewGame -
func (overlaySymbol *OverlaySymbol) OnNewGame(gameProp *GameProperty) error {
	osd := gameProp.MapComponentData[overlaySymbol.Name].(*OverlaySymbolData)

	osd.OnNewGame()

	osd.CurLevel = overlaySymbol.Config.DefaultLevel

	return nil
}

// Init -
func (overlaySymbol *OverlaySymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("OverlaySymbol.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &OverlaySymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("OverlaySymbol.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	overlaySymbol.Config = cfg

	if overlaySymbol.Config.MapPosition != "" {
		vm2, err := sgc7game.LoadValMapping2FromExcel(overlaySymbol.Config.MapPosition, "index", "value", sgc7game.NewIntArrVal[int])
		if err != nil {
			goutils.Error("OverlaySymbol.Init:LoadValMapping2FromExcel",
				zap.String("valmapping", overlaySymbol.Config.MapPosition),
				zap.Error(err))

			return err
		}

		overlaySymbol.MapPosition = vm2
	}

	overlaySymbol.SymbolCode = pool.DefaultPaytables.MapSymbols[cfg.Symbol]

	overlaySymbol.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (overlaySymbol *OverlaySymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	osd := gameProp.MapComponentData[overlaySymbol.Name].(*OverlaySymbolData)

	gs := overlaySymbol.GetTargetScene(gameProp, curpr, &osd.BasicComponentData)

	cgs := gs.Clone()

	for i, pos := range overlaySymbol.MapPosition.MapVals {
		if i < osd.CurLevel {
			cgs.Arr[pos.GetInt(0)][pos.GetInt(1)] = overlaySymbol.SymbolCode
		} else {
			break
		}
	}

	overlaySymbol.AddScene(gameProp, curpr, cgs, &osd.BasicComponentData)

	overlaySymbol.onStepEnd(gameProp, curpr, gp)

	gp.AddComponentData(overlaySymbol.Name, &osd.BasicComponentData)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (overlaySymbol *OverlaySymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	osd := gameProp.MapComponentData[overlaySymbol.Name].(*OverlaySymbolData)

	if len(osd.UsedScenes) > 0 {
		asciigame.OutputScene("The symbols after the symbol overlay", pr.OtherScenes[osd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (overlaySymbol *OverlaySymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (overlaySymbol *OverlaySymbol) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.OverlaySymbolData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("OverlaySymbol.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return overlaySymbol.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

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
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewOverlaySymbol(name string) IComponent {
	return &OverlaySymbol{
		BasicComponent: NewBasicComponent(name),
	}
}

package stats

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

var featureHeaders []string

type FeatureType int

const (
	FeatureBasic    = 1
	FeatureRespin   = 2
	FeatureFreeGame = 3
)

type FuncAnalyzeFeature func(*Feature, *sgc7game.Stake, []*sgc7game.PlayResult) (bool, int64, int64)

type Feature struct {
	Name           string
	Type           FeatureType
	PlayTimes      int64
	TotalBets      int64
	TotalWins      int64
	TriggerTimes   int64
	RetriggerTimes int64
	FreeSpinTimes  int64
	RoundTimes     int64
	Parent         *Feature
	Children       []*Feature
	OnAnalyze      FuncAnalyzeFeature
	Reels          *Reels
	Symbols        *SymbolsRTP
	Obj            interface{}
}

func (feature *Feature) GetPlayTimes() int64 {
	if feature.Parent != nil {
		return feature.Parent.GetPlayTimes()
	}

	return feature.PlayTimes
}

func (feature *Feature) GetTotalBets() int64 {
	if feature.Parent != nil {
		return feature.Parent.GetTotalBets()
	}

	return feature.TotalBets
}

func (feature *Feature) OnResults(stake *sgc7game.Stake, lst []*sgc7game.PlayResult) {
	feature.PlayTimes++

	istrigger, bet, wins := feature.OnAnalyze(feature, stake, lst)
	if istrigger {
		feature.TriggerTimes++

		feature.TotalWins += wins

		feature.onTrigger(stake, lst)
	}

	feature.TotalBets += bet
}

func (feature *Feature) onTrigger(stake *sgc7game.Stake, lst []*sgc7game.PlayResult) {
	for _, v := range feature.Children {
		v.OnResults(stake, lst)
	}
}

func (feature *Feature) saveOtherSheet(f *excelize.File) error {
	if feature.Reels != nil {
		csheet := fmt.Sprintf("symbol in window - %v", feature.Name)
		f.NewSheet(csheet)
		feature.Reels.SaveSheet(f, csheet)
	}

	for _, v := range feature.Children {
		v.saveOtherSheet(f)
	}

	return nil
}

func (feature *Feature) saveSheet(f *excelize.File, sheet string, startx, starty int) int {
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+0, starty), feature.Name)

	if feature.Parent != nil {
		f.SetCellValue(sheet, goutils.Pos2Cell(startx+1, starty), feature.Parent.Name)
	}

	f.SetCellValue(sheet, goutils.Pos2Cell(startx+2, starty), feature.GetPlayTimes())
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+3, starty), feature.GetTotalBets())
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+4, starty), feature.TotalWins)
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+5, starty), float64(feature.TotalWins)/float64(feature.GetTotalBets()))
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+6, starty), feature.TriggerTimes)
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+7, starty), feature.RetriggerTimes)
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+8, starty), feature.FreeSpinTimes)
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+9, starty), float64(feature.FreeSpinTimes)/float64(feature.TriggerTimes))
	f.SetCellValue(sheet, goutils.Pos2Cell(startx+10, starty), feature.RoundTimes)

	if feature.FreeSpinTimes > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(startx+11, starty), float64(feature.RoundTimes)/float64(feature.FreeSpinTimes))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(startx+11, starty), float64(feature.RoundTimes)/float64(feature.TriggerTimes))
	}

	f.SetCellValue(sheet, goutils.Pos2Cell(startx+12, starty), float64(feature.TriggerTimes)/float64(feature.GetPlayTimes()))

	starty++

	for _, v := range feature.Children {
		starty = v.saveSheet(f, sheet, startx, starty)
	}

	return starty
}

func (feature *Feature) SaveExcel(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	for i, v := range featureHeaders {
		f.SetCellStr(sheet, goutils.Pos2Cell(i, 0), v)
	}

	feature.saveSheet(f, sheet, 0, 1)

	feature.saveOtherSheet(f)

	return f.SaveAs(fn)
}

func NewFeature(name string, ft FeatureType, onAnalyze FuncAnalyzeFeature, parent *Feature) *Feature {
	feature := &Feature{
		Name:      name,
		Type:      ft,
		OnAnalyze: onAnalyze,
		Parent:    parent,
	}

	if parent != nil {
		parent.Children = append(parent.Children, feature)
	}

	return feature
}

func init() {
	featureHeaders = []string{
		"gamemod",
		"parent",
		"playtimes",
		"bet",
		"wins",
		"rtp",
		"triggertimes",
		"retriggertimes",
		"freespintimes",
		"avgfreespintimes",
		"roundtimes",
		"avgroundtimes",
		"hit rate",
	}
}

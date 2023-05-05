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

// isTrigger, bet, wins
type FuncAnalyzeFeature func(*Feature, *sgc7game.Stake, []*sgc7game.PlayResult) (bool, int64, int64)

type Feature struct {
	Name                 string
	Type                 FeatureType
	PlayTimes            int64
	TotalBets            int64
	TotalWins            int64
	TriggerTimes         int64
	RetriggerTimes       int64
	FreeSpinTimes        int64
	RoundTimes           int64
	Parent               *Feature
	Children             []*Feature
	OnAnalyze            FuncAnalyzeFeature
	Reels                *Reels
	Symbols              *SymbolsRTP
	CurWins              *Wins
	AllWins              *Wins
	RespinEndingStatus   *Status
	RespinEndingName     string
	RespinStartStatus    *Status
	RespinStartName      string
	RespinNumStatus      *Status
	RespinWinStatus      *Status
	RespinStartNumStatus *Status
	Obj                  any
}

func (feature *Feature) CloneIncludeChildren() *Feature {
	nf := &Feature{
		Name:           feature.Name,
		Type:           feature.Type,
		PlayTimes:      feature.PlayTimes,
		TotalBets:      feature.TotalBets,
		TotalWins:      feature.TotalWins,
		TriggerTimes:   feature.TriggerTimes,
		RetriggerTimes: feature.RetriggerTimes,
		FreeSpinTimes:  feature.FreeSpinTimes,
		RoundTimes:     feature.RoundTimes,
		Parent:         feature.Parent,
		Children:       make([]*Feature, 0, len(feature.Children)),
		OnAnalyze:      feature.OnAnalyze,
		Obj:            feature.Obj,
	}

	if feature.Reels != nil {
		nf.Reels = feature.Reels.Clone()
	}

	if feature.Symbols != nil {
		nf.Symbols = feature.Symbols.Clone()
	}

	if feature.CurWins != nil {
		nf.CurWins = feature.CurWins.Clone()
	}

	if feature.AllWins != nil {
		nf.AllWins = feature.AllWins.Clone()
	}

	if feature.RespinEndingStatus != nil {
		nf.RespinEndingStatus = feature.RespinEndingStatus.Clone()
	}

	if feature.RespinStartStatus != nil {
		nf.RespinStartStatus = feature.RespinStartStatus.Clone()
	}

	if feature.RespinNumStatus != nil {
		nf.RespinNumStatus = feature.RespinNumStatus.Clone()
	}

	if feature.RespinWinStatus != nil {
		nf.RespinWinStatus = feature.RespinWinStatus.Clone()
	}

	if feature.RespinStartNumStatus != nil {
		nf.RespinStartNumStatus = feature.RespinStartNumStatus.Clone()
	}

	for _, v := range feature.Children {
		nf.Children = append(nf.Children, v.CloneIncludeChildren())
	}

	return nf
}

func (feature *Feature) Clone() *Feature {
	nf := &Feature{
		Name:           feature.Name,
		Type:           feature.Type,
		PlayTimes:      feature.PlayTimes,
		TotalBets:      feature.TotalBets,
		TotalWins:      feature.TotalWins,
		TriggerTimes:   feature.TriggerTimes,
		RetriggerTimes: feature.RetriggerTimes,
		FreeSpinTimes:  feature.FreeSpinTimes,
		RoundTimes:     feature.RoundTimes,
		Parent:         feature.Parent,
		Children:       make([]*Feature, len(feature.Children)),
		OnAnalyze:      feature.OnAnalyze,
		Obj:            feature.Obj,
	}

	if feature.Reels != nil {
		nf.Reels = feature.Reels.Clone()
	}

	if feature.Symbols != nil {
		nf.Symbols = feature.Symbols.Clone()
	}

	if feature.CurWins != nil {
		nf.CurWins = feature.CurWins.Clone()
	}

	if feature.AllWins != nil {
		nf.AllWins = feature.AllWins.Clone()
	}

	if feature.RespinEndingStatus != nil {
		nf.RespinEndingStatus = feature.RespinEndingStatus.Clone()
	}

	if feature.RespinStartStatus != nil {
		nf.RespinStartStatus = feature.RespinStartStatus.Clone()
	}

	if feature.RespinNumStatus != nil {
		nf.RespinNumStatus = feature.RespinNumStatus.Clone()
	}

	if feature.RespinWinStatus != nil {
		nf.RespinWinStatus = feature.RespinWinStatus.Clone()
	}

	if feature.RespinStartNumStatus != nil {
		nf.RespinStartNumStatus = feature.RespinStartNumStatus.Clone()
	}

	return nf
}

func (feature *Feature) Merge(src *Feature) {
	feature.PlayTimes += src.PlayTimes
	feature.TotalBets += src.TotalBets
	feature.TotalWins += src.TotalWins
	feature.TriggerTimes += src.TriggerTimes
	feature.RetriggerTimes += src.RetriggerTimes
	feature.FreeSpinTimes += src.FreeSpinTimes
	feature.RoundTimes += src.RoundTimes

	if feature.Reels != nil && src.Reels != nil {
		feature.Reels.Merge(src.Reels)
	}

	if feature.Symbols != nil && src.Symbols != nil {
		feature.Symbols.Merge(src.Symbols)
	}

	if feature.CurWins != nil && src.CurWins != nil {
		feature.CurWins.Merge(src.CurWins)
	}

	if feature.AllWins != nil && src.AllWins != nil {
		feature.AllWins.Merge(src.AllWins)
	}

	for i, v := range feature.Children {
		v.Merge(src.Children[i])
	}

	if feature.RespinEndingStatus != nil && src.RespinEndingStatus != nil {
		feature.RespinEndingStatus.Merge(src.RespinEndingStatus)
	}

	if feature.RespinStartStatus != nil && src.RespinStartStatus != nil {
		feature.RespinStartStatus.Merge(src.RespinStartStatus)
	}

	if feature.RespinNumStatus != nil && src.RespinNumStatus != nil {
		feature.RespinNumStatus.Merge(src.RespinNumStatus)
	}

	if feature.RespinWinStatus != nil && src.RespinWinStatus != nil {
		feature.RespinWinStatus.Merge(src.RespinWinStatus)
	}

	if feature.RespinStartNumStatus != nil && src.RespinStartNumStatus != nil {
		feature.RespinStartNumStatus.Merge(src.RespinStartNumStatus)
	}
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

func (feature *Feature) Retrigger() {
	feature.RetriggerTimes++
}

func (feature *Feature) OnFreeSpin() {
	feature.FreeSpinTimes++
}

func (feature *Feature) OnRound() {
	feature.RoundTimes++
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

	if feature.Symbols != nil {
		csheet := fmt.Sprintf("symbol rtp - %v", feature.Name)
		f.NewSheet(csheet)
		feature.Symbols.SaveSheet(f, csheet, feature.GetTotalBets())
	}

	if feature.AllWins != nil {
		csheet := fmt.Sprintf("total wins - %v", feature.Name)
		f.NewSheet(csheet)
		feature.AllWins.SaveSheet(f, csheet)
	}

	if feature.CurWins != nil {
		csheet := fmt.Sprintf("wins - %v", feature.Name)
		f.NewSheet(csheet)
		feature.CurWins.SaveSheet(f, csheet)
	}

	if feature.RespinEndingStatus != nil {
		csheet := fmt.Sprintf("respinEndingStatus - %v", feature.Name)
		f.NewSheet(csheet)
		feature.RespinEndingStatus.SaveSheet(f, csheet)
	}

	if feature.RespinStartStatus != nil {
		csheet := fmt.Sprintf("respinStartStatus - %v", feature.Name)
		f.NewSheet(csheet)
		feature.RespinStartStatus.SaveSheet(f, csheet)
	}

	if feature.RespinNumStatus != nil {
		csheet := fmt.Sprintf("respinNumStatus - %v", feature.Name)
		f.NewSheet(csheet)
		feature.RespinNumStatus.SaveSheet(f, csheet)
	}

	if feature.RespinWinStatus != nil {
		csheet := fmt.Sprintf("respinWinStatus - %v", feature.Name)
		f.NewSheet(csheet)
		feature.RespinWinStatus.SaveSheet(f, csheet)
	}

	if feature.RespinStartNumStatus != nil {
		csheet := fmt.Sprintf("respinStartNumStatus - %v", feature.Name)
		f.NewSheet(csheet)
		feature.RespinStartNumStatus.SaveSheet(f, csheet)
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

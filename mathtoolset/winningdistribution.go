package mathtoolset

import (
	"math"
	"os"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type WinTimesData struct {
	Win   int
	Times int64
}

type WinPerData struct {
	Win     int
	Percent float64
}

type AvgWinData struct {
	AvgWin  float64             `yaml:"win" json:"win"`
	Percent float64             `yaml:"percent" json:"percent"`
	MapWins map[int]*WinPerData `yaml:"-" json:"-"`
}

type WinningDistribution struct {
	TimesWins   map[int]*WinTimesData `yaml:"-" json:"-"`
	PercentWins map[int]*WinPerData   `yaml:"-" json:"-"`
	AvgWins     map[int]*AvgWinData   `yaml:"avgwins" json:"avgwins"`
	TotalTimes  int64                 `yaml:"-" json:"-"`
}

func (wd *WinningDistribution) getAvgWin(si, ci int) float64 {
	vw := float64(0)
	totalw := float64(0)

	for i := si; i <= ci; i++ {
		v, isok := wd.AvgWins[i]
		if isok {
			vw += v.AvgWin * v.Percent
			totalw += v.Percent
		}
	}

	return vw / totalw
}

func (wd *WinningDistribution) getMax() int {
	max := 0

	for k := range wd.AvgWins {
		if k > max {
			max = k
		}
	}

	return max
}

func (wd *WinningDistribution) AddTimesWin(win int, times int64) {
	wd.TotalTimes += times

	wind, isok := wd.TimesWins[win]
	if isok {
		wind.Times += times
	} else {
		wd.TimesWins[win] = &WinTimesData{
			Win:   win,
			Times: times,
		}
	}
}

func (wd *WinningDistribution) AddPercentWin(win int, percent float64) {
	wind, isok := wd.PercentWins[win]
	if isok {
		wind.Percent += percent
	} else {
		wd.PercentWins[win] = &WinPerData{
			Win:     win,
			Percent: percent,
		}
	}
}

func (wd *WinningDistribution) addAvgWin(bet int, win int, percent float64) {
	wini := -1
	winf := float64(win) / float64(bet)

	if win > 0 {
		wini = int(math.Floor(winf))
	}

	wind, isok := wd.AvgWins[wini]
	if isok {
		winpd, isok2 := wind.MapWins[win]
		if isok2 {
			winpd.Percent += percent
		} else {
			wind.MapWins[win] = &WinPerData{
				Win:     win,
				Percent: percent,
			}
		}
	} else {
		wind = &AvgWinData{
			MapWins: make(map[int]*WinPerData),
		}

		wind.MapWins[win] = &WinPerData{
			Win:     win,
			Percent: percent,
		}

		wd.AvgWins[wini] = wind
	}
}

func (wd *WinningDistribution) rebuildAvgWin(bet int) {
	for k, v0 := range wd.AvgWins {
		if k == -1 {
			v0.AvgWin = 0
			v0.Percent = 0

			for _, v1 := range v0.MapWins {
				v0.Percent += v1.Percent
			}

			continue
		}

		twin := float64(0)
		tper := float64(0)

		for _, v1 := range v0.MapWins {
			twin += float64(v1.Win) / float64(bet) * v1.Percent
			tper += v1.Percent
		}

		v0.AvgWin = twin / tper
		v0.Percent = tper
	}
}

func (wd *WinningDistribution) genAvgWinDataWithTimes(bet int) {
	wd.AvgWins = make(map[int]*AvgWinData)

	for _, v := range wd.TimesWins {
		wd.addAvgWin(bet, v.Win, float64(v.Times)/float64(wd.TotalTimes))
	}

	wd.rebuildAvgWin(bet)
}

func (wd *WinningDistribution) genAvgWinDataWithPercent(bet int) {
	wd.AvgWins = make(map[int]*AvgWinData)

	for _, v := range wd.PercentWins {
		wd.addAvgWin(bet, v.Win, v.Percent)
	}

	wd.rebuildAvgWin(bet)
}

func (wd *WinningDistribution) GenAvgWinData(bet int) {
	if len(wd.TimesWins) > 0 {
		wd.genAvgWinDataWithTimes(bet)
	} else if len(wd.PercentWins) > 0 {
		wd.genAvgWinDataWithPercent(bet)
	}
}

func (wd *WinningDistribution) saveTimes(f *excelize.File) {
	if len(wd.TimesWins) <= 0 {
		return
	}

	sheet := "times"
	f.NewSheet(sheet)

	lsthead := []string{
		"Win",
		"Times",
	}

	for x, v := range lsthead {
		f.SetCellStr(sheet, goutils.Pos2Cell(x, 0), v)
	}

	y := 1

	for _, v := range wd.TimesWins {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, y), v.Win)
		f.SetCellInt(sheet, goutils.Pos2Cell(1, y), int(v.Times))

		y++
	}
}

func (wd *WinningDistribution) savePercent(f *excelize.File, scale float64) {
	if len(wd.PercentWins) <= 0 {
		return
	}

	sheet := "percent"
	f.NewSheet(sheet)

	lsthead := []string{
		"Win",
		"Percent",
	}

	for x, v := range lsthead {
		f.SetCellStr(sheet, goutils.Pos2Cell(x, 0), v)
	}

	y := 1

	for _, v := range wd.PercentWins {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, y), v.Win)
		f.SetCellFloat(sheet, goutils.Pos2Cell(1, y), v.Percent*scale, 5, 64)

		y++
	}
}

func (wd *WinningDistribution) saveAvgWin(f *excelize.File, scale float64) {
	if len(wd.AvgWins) <= 0 {
		return
	}

	sheet := "avgwin"
	f.NewSheet(sheet)

	lsthead := []string{
		"AvgWin",
		"Percent",
	}

	for x, v := range lsthead {
		f.SetCellStr(sheet, goutils.Pos2Cell(x, 0), v)
	}

	y := 1

	for _, v := range wd.AvgWins {
		f.SetCellFloat(sheet, goutils.Pos2Cell(0, y), v.AvgWin, 3, 64)
		f.SetCellFloat(sheet, goutils.Pos2Cell(1, y), v.Percent*scale, 5, 64)

		y++
	}
}

func (wd *WinningDistribution) SaveExcel(fn string, scale float64) {
	f := excelize.NewFile()

	wd.saveTimes(f)
	wd.savePercent(f, scale)
	wd.saveAvgWin(f, scale)

	lstname := f.GetSheetList()
	f.DeleteSheet(lstname[0])

	f.SaveAs(fn)
}

func (wd *WinningDistribution) Save(fn string) {
	buf, err := yaml.Marshal(wd)
	if err != nil {
		goutils.Error("Save:Marshal",
			zap.String("fn", fn),
			zap.Error(err))

		return
	}

	os.WriteFile(fn, buf, 0644)
}

func (wd *WinningDistribution) mergeAvgWins(mini, maxi int) int {
	nawd := &AvgWinData{}
	totalw := float64(0)
	totalp := float64(0)

	for i := mini; i <= maxi; i++ {
		nd, isok := wd.AvgWins[i]
		if isok {
			totalw += nd.AvgWin * nd.Percent
			totalp += nd.Percent

			for k0, v0 := range nd.MapWins {
				nawd.MapWins[k0] = v0
			}

			delete(wd.AvgWins, i)
		}
	}

	nawd.AvgWin = totalw / totalp
	nawd.Percent = totalp

	newi := int(math.Floor(nawd.AvgWin))

	wd.AvgWins[newi] = nawd

	return newi
}

func LoadWinningDistribution(fn string) (*WinningDistribution, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("LoadWinningDistribution:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	wd := NewWinningDistribution()
	err = yaml.Unmarshal(data, wd)
	if err != nil {
		goutils.Error("LoadWinningDistribution:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	return wd, nil
}

func NewWinningDistribution() *WinningDistribution {
	return &WinningDistribution{
		TimesWins:   make(map[int]*WinTimesData),
		PercentWins: make(map[int]*WinPerData),
	}
}

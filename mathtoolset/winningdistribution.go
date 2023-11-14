package mathtoolset

import (
	"math"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
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

func (wd *WinningDistribution) getAllAvgWin() float64 {
	vw := float64(0)
	totalw := float64(0)

	for _, v := range wd.AvgWins {
		vw += v.AvgWin * v.Percent
		totalw += v.Percent
	}

	return vw / totalw
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

func (wd *WinningDistribution) AddAvgWin(winf float64, percent float64) error {
	if percent == 0 {
		return nil
	}

	wini := int(math.Floor(winf))

	if winf == 0 {
		wini = -1
	}

	_, isok := wd.AvgWins[wini]
	if isok {
		goutils.Error("WinningDistribution.AddAvgWin",
			zap.Int("key", wini),
			zap.Error(ErrDuplicateAvgWin))

		return ErrDuplicateAvgWin
	} else {
		wind := &AvgWinData{
			AvgWin:  winf,
			Percent: percent,
		}

		wd.AvgWins[wini] = wind
	}

	return nil
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

func (wd *WinningDistribution) fill(maxwin int) {
	lastper := float64(0)
	for i := -1; i <= maxwin; i++ {
		nd, isok := wd.AvgWins[i]
		if isok {
			lastper = nd.Percent
		} else {
			wd.AvgWins[i] = &AvgWinData{
				AvgWin:  float64(i) + 0.5,
				Percent: lastper,
			}
		}
	}
}

func (wd *WinningDistribution) setAvgWinPercent(si int, ei int, per float64) {
	for i := si; i <= ei; i++ {
		wd.AvgWins[i].Percent = per
	}
}

func (wd *WinningDistribution) mergeSmooth(si0 int, si1 int, per float64) {
	totalper := float64(0)
	totalwin := float64(0)

	for i := si0; i <= si1; i++ {
		totalper += wd.AvgWins[i].Percent * wd.AvgWins[i].AvgWin
		totalwin += wd.AvgWins[i].AvgWin
	}

	cp := totalper / float64(totalwin)
	wd.setAvgWinPercent(si0, si1, cp)
}

func (wd *WinningDistribution) getMaxPercent(si int, ei int) (int, float64) {
	maxi := si
	maxp := wd.AvgWins[si].Percent

	for i := si + 1; i <= ei; i++ {
		if wd.AvgWins[i].Percent >= maxp {
			maxi = i
			maxp = wd.AvgWins[i].Percent
		}
	}

	return maxi, maxp
}

func (wd *WinningDistribution) getPreLessPercent(si int, per float64) int {
	prei := si

	for i := si; i >= 0; i-- {
		if wd.AvgWins[i].Percent <= per {
			prei = i
		} else {
			return prei
		}
	}

	return 0
}

func (wd *WinningDistribution) smooth(maxwin int) {
	lastper := wd.AvgWins[-1].Percent
	for i := 0; i <= maxwin; i++ {
		nd, isok := wd.AvgWins[i]
		if isok {
			if nd.Percent > lastper {
				ei, maxp := wd.getMaxPercent(i, maxwin)
				prei := wd.getPreLessPercent(i, maxp)

				wd.mergeSmooth(prei, ei, maxp)
			}

			lastper = nd.Percent
		}
	}
}

func (wd *WinningDistribution) Format(maxwin int, isneedsmooth bool) {
	wd.fill(maxwin)

	if isneedsmooth {
		wd.smooth(maxwin)
	}

	totalper := float64(0)

	for _, v := range wd.AvgWins {
		totalper += v.Percent
	}

	for _, v := range wd.AvgWins {
		v.Percent /= totalper
	}
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

func LoadWinningDistributionFromExcel(fn string) (*WinningDistribution, error) {
	wd := NewWinningDistribution()
	curdat := 0
	curavgwin := 0.0
	curpercent := 0.0
	cury := -1
	err := sgc7game.LoadExcel(fn, "avgwin", func(x int, str string) string {
		return strings.TrimSpace(strings.ToLower(str))
	}, func(x int, y int, header string, data string) error {
		if y != cury {
			if curdat == 2 {
				wd.AddAvgWin(curavgwin, curpercent)
			}

			cury = y
			curdat = 0
		}

		if header == "avgwin" {
			winf, err := goutils.String2Float64(data)
			if err != nil {
				goutils.Error("LoadWinningDistributionFromExcel:String2Float64:avgwin",
					zap.String("avgwin", data),
					zap.Error(err))

				return err
			}

			curavgwin = winf
			curdat++
		} else if header == "percent" {
			perf, err := goutils.String2Float64(data)
			if err != nil {
				goutils.Error("LoadWinningDistributionFromExcel:String2Float64:percent",
					zap.String("percent", data),
					zap.Error(err))

				return err
			}

			curpercent = perf
			curdat++
		}

		return nil
	})
	if err != nil {
		goutils.Error("LoadWinningDistributionFromExcel:LoadExcel",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	if curdat == 2 {
		wd.AddAvgWin(curavgwin, curpercent)
	}

	return wd, nil
}

func NewWinningDistribution() *WinningDistribution {
	return &WinningDistribution{
		TimesWins:   make(map[int]*WinTimesData),
		PercentWins: make(map[int]*WinPerData),
		AvgWins:     make(map[int]*AvgWinData),
	}
}

package stats2

import (
	"log/slog"
	"os"
	"time"

	"github.com/bytedance/sonic"
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type Stats struct {
	MapStats       map[string]*Feature `json:"mapStats"`
	chanBet        chan int            `json:"-"`
	chanCache      chan *Cache         `json:"-"`
	TotalBet       int64               `json:"totalBet"`
	TotalWins      int64               `json:"totalWins"`
	BetTimes       int64               `json:"betTimes"`
	MaxWins        int64               `json:"maxWins"`
	MaxWinTimes    int64               `json:"maxWinTimes"`
	BetEndingTimes int64               `json:"-"`
	Components     []string            `json:"components"`
}

func (s2 *Stats) PushBet(bet int) {
	s2.chanBet <- bet
}

func (s2 *Stats) PushCache(cache *Cache) {
	s2.chanCache <- cache
}

func (s2 *Stats) onCache(cache *Cache) {
	for k, v := range cache.MapStats {
		s, isok := s2.MapStats[k]
		if isok {
			s.Merge(v)
		}
	}

	if cache.TotalWin > s2.MaxWins {
		s2.MaxWins = cache.TotalWin
		s2.MaxWinTimes = 1
	} else if cache.TotalWin == s2.MaxWins {
		s2.MaxWinTimes++
	}

	s2.TotalWins += cache.TotalWin
}

func (s2 *Stats) SaveExcel(fn string) error {
	buf, err := s2.ExportExcel()
	if err != nil {
		goutils.Error("Stats.SaveExcel:ExportExcel",
			goutils.Err(err))

		return err
	}

	os.WriteFile(fn, buf, 0644)

	return nil
	// f := excelize.NewFile()

	// for _, cn := range s2.Components {
	// 	f2, isok := s2.MapStats[cn]
	// 	if isok {
	// 		f2.SaveSheet(f, cn, s2)
	// 	}
	// }

	// return f.SaveAs(fn)
}

func (s2 *Stats) saveBasicSheet(f *excelize.File) {
	sheet := "basic"
	f.NewSheet(sheet)

	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "spin times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 1), "total bet")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 2), "total wins")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 3), "rtp")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 4), "max wins")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 5), "times of the max wins")

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), s2.BetEndingTimes)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 1), s2.TotalBet)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), s2.TotalWins)

	if s2.TotalBet > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 3), float64(s2.TotalWins)/float64(s2.TotalBet))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 3), 0)
	}

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 4), s2.MaxWins)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 5), s2.MaxWinTimes)
}

func (s2 *Stats) ExportExcel() ([]byte, error) {
	f := excelize.NewFile()

	s2.saveBasicSheet(f)

	f.DeleteSheet(f.GetSheetName(0))

	for _, cn := range s2.Components {
		f2, isok := s2.MapStats[cn]
		if isok {
			f2.SaveSheet(f, cn, s2)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		goutils.Error("Stats.ExportExcel:WriteToBuffer",
			goutils.Err(err))

		return nil, err
	}

	return buf.Bytes(), nil
}

func (s2 *Stats) Start() {
	go func() {
		for {
			bet := <-s2.chanBet

			s2.BetTimes++
			s2.TotalBet += int64(bet)
		}
	}()

	go func() {
		for {
			cache := <-s2.chanCache

			s2.onCache(cache)

			s2.BetEndingTimes++
		}
	}()
}

func (s2 *Stats) WaitEnding() {
	for {
		if s2.BetTimes == s2.BetEndingTimes {
			return
		}

		time.Sleep(time.Second)
	}
}

func (s2 *Stats) AddFeature(name string, feature *Feature) {
	s2.MapStats[name] = feature
}

func (s2 *Stats) GetRunTimes(name string) int64 {
	if name == "" {
		return s2.BetTimes
	}

	f2, isok := s2.MapStats[name]
	if isok {
		return f2.RootTrigger.RunTimes
	}

	return 0
}

func (s2 *Stats) Merge(src *Stats) {
	for k, v := range src.MapStats {
		cv, isok := s2.MapStats[k]
		if isok {
			cv.Merge(v)
		}
	}

	s2.TotalBet += src.TotalBet
	s2.TotalWins += src.TotalWins
	s2.BetTimes += src.BetTimes
	if src.MaxWins > s2.MaxWins {
		s2.MaxWins = src.MaxWins
		s2.MaxWinTimes = src.MaxWinTimes
	} else if src.MaxWins == s2.MaxWins {
		s2.MaxWinTimes += src.MaxWinTimes
	}
}

func (s2 *Stats) ToJson() string {
	str, _ := sonic.MarshalString(s2)

	return str
}

func NewStats(components []string) *Stats {
	s2 := &Stats{
		MapStats:   make(map[string]*Feature),
		chanBet:    make(chan int, 1024),
		chanCache:  make(chan *Cache, 1024),
		Components: components,
	}

	return s2
}

func LoadStats(str string) (*Stats, error) {
	s2 := &Stats{}

	err := sonic.UnmarshalString(str, s2)
	if err != nil {
		goutils.Error("LoadStats:UnmarshalString",
			slog.String("str", str),
			goutils.Err(err))

		return nil, err
	}

	return s2, nil
}

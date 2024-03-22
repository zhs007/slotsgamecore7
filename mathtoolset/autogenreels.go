package mathtoolset

import (
	"log/slog"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type agrData struct {
	reelIndex int
	symbol    SymbolType
}

func (d *agrData) IsSame(reelindex int, symbol SymbolType) bool {
	return d.reelIndex == reelindex && d.symbol == symbol
}

func newAGRData(reelindex int, symbol SymbolType) *agrData {
	return &agrData{
		reelIndex: reelindex,
		symbol:    symbol,
	}
}

type agrDataList struct {
	lst     []*agrData
	weights *sgc7game.ValWeights
}

func (lst *agrDataList) has(reelindex int, symbol SymbolType) bool {
	for _, v := range lst.lst {
		if v.IsSame(reelindex, symbol) {
			return true
		}
	}

	return false
}

func (lst *agrDataList) add(reelindex int, symbol SymbolType, weight int) error {
	if !lst.has(reelindex, symbol) {
		d := newAGRData(reelindex, symbol)

		i := len(lst.lst)
		lst.lst = append(lst.lst, d)
		lst.weights.Add(i, weight)

		return nil
	}

	goutils.Error("agrDataList.add",
		goutils.Err(ErrInvalidDataInAGRDataList))

	return ErrInvalidDataInAGRDataList
}

func (lst *agrDataList) rand(plugin sgc7plugin.IPlugin) (*agrData, error) {
	i, err := lst.weights.RandVal(plugin)
	if err != nil {
		goutils.Error("agrDataList.rand",
			goutils.Err(err))

		return nil, err
	}

	return lst.lst[i], nil
}

func newAGRDataList() *agrDataList {
	lst := &agrDataList{
		weights: sgc7game.NewValWeightsEx(),
	}

	return lst
}

func genAGRDataList(rss *ReelsStats, syms []SymbolType) *agrDataList {
	lst := newAGRDataList()

	for i, rs := range rss.Reels {
		cursyms := rs.GetCanAddSymbols(syms)

		for _, s := range cursyms {
			lst.add(i, s, 1)
		}
	}

	return lst
}

func AutoGenReels(w, h int, paytables *sgc7game.PayTables, syms []SymbolType, wilds []SymbolType,
	totalBet int, lineNum int) (*sgc7game.ReelsData, error) {

	rss, err := BuildBasicReelsStatsEx(w, syms)
	if err != nil {
		goutils.Error("AutoGenReels:BuildBasicReelsStatsEx",
			goutils.Err(err))

		return nil, err
	}

	ssws, err := AnalyzeReelsWithLineEx(paytables, rss, syms, wilds, totalBet, lineNum)
	if err != nil {
		goutils.Error("AutoGenReels:AnalyzeReelsWithLineEx",
			goutils.Err(err))

		return nil, err
	}

	plugin := sgc7plugin.NewBasicPlugin()

	lst := genAGRDataList(rss, syms)
	d, err := lst.rand(plugin)
	if err != nil {
		goutils.Error("AutoGenReels",
			goutils.Err(err))

		return nil, err
	}

	goutils.Info("AutoGenReels",
		slog.Int("ri", d.reelIndex),
		slog.Any("s", d.symbol),
		slog.Float64("rtp", ssws.CountRTP()))

	return nil, nil
}

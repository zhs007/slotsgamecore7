package lowcode

import (
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7rtp "github.com/zhs007/slotsgamecore7/rtp"
	"github.com/zhs007/slotsgamecore7/stats2"
)

type RTPSymbolFeature struct {
	Name       string   `yaml:"name"`
	Components []string `yaml:"components"`
}

type RTPSymbolModule struct {
	Name       string              `yaml:"name"`
	Components []string            `yaml:"components"`
	Features   []*RTPSymbolFeature `yaml:"features"`
}

type RTPHitRateFeature struct {
	Name       string   `yaml:"name"`
	Components []string `yaml:"components"`
}

type RTPConfig struct {
	Modules         []*RTPSymbolModule   `yaml:"modules"`
	HitRateFeatures []*RTPHitRateFeature `yaml:"hitRateFeatures"`
}

func StartRTP(gamecfg string, icore int, ispinnums int64, outputPath string, bet int64, coin int64, funcNewRNG FuncNewRNG, funcNewFeatureLevel FuncNewFeatureLevel, wincap int64) error {
	sgc7plugin.IsNoRNGCache = true

	if wincap > 0 {
		stats2.SetWinCap(int(wincap))
	}

	game, err := NewGame2(gamecfg, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	}, funcNewRNG, funcNewFeatureLevel)
	if err != nil {
		goutils.Error("StartRTP:NewGame3",
			slog.String("gamecfg", gamecfg),
			goutils.Err(err))

		return err
	}

	rtp := sgc7rtp.NewRTP()

	if bet <= 0 {
		bet = int64(game.Pool.Config.Bets[0])
	}

	if coin <= 0 {
		coin = 1
	}

	stake := &sgc7game.Stake{
		CoinBet:  coin,
		CashBet:  bet * coin,
		Currency: "EUR",
	}

	d := sgc7rtp.StartRTP2(game, rtp, icore, ispinnums, stake, 100000, func(totalnums int64, curnums int64, curtime time.Duration) {
		goutils.Info("processing...",
			slog.Int64("total nums", totalnums),
			slog.Int64("current nums", curnums),
			slog.String("cost time", curtime.String()))
	}, true, wincap)

	goutils.Info("finish.",
		slog.Int64("total nums", ispinnums),
		slog.Float64("rtp", float64(rtp.TotalWins)/float64(rtp.TotalBet)),
		slog.Duration("cost time", d))

	curtime := time.Now()

	rtp.Save2CSV(path.Join(outputPath, fmt.Sprintf("%v-%v.csv", game.Pool.Config.Name, curtime.Format("2006-01-02_15_04_05"))))

	if gAllowStats2 {
		components := game.Pool.mapComponents[int(bet)]
		components.Stats2.WaitEnding()

		components.Stats2.SaveExcel(path.Join(outputPath, fmt.Sprintf("%v-%v-stats-%v.xlsx", game.Pool.Config.Name, bet, curtime.Format("2006-01-02_15_04_05"))))

		goutils.Info("finish.",
			slog.Int64("total nums", components.Stats2.BetTimes))
	}

	return nil
}

func StartRTPWithData(gamecfg []byte, icore int, ispinnums int64, bet int64, ontimer sgc7rtp.FuncOnRTPTimer, funcNewRNG FuncNewRNG, funcNewFeatureLevel FuncNewFeatureLevel, wincap int64) (*stats2.Stats, error) {
	sgc7plugin.IsNoRNGCache = true

	game, err := NewGame2WithData(gamecfg, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	}, funcNewRNG, funcNewFeatureLevel)
	if err != nil {
		goutils.Error("StartRTPWithData:NewGame3",
			goutils.Err(err))

		return nil, err
	}

	rtp := sgc7rtp.NewRTP()

	if bet <= 0 {
		bet = int64(game.Pool.Config.Bets[0])
	}

	stake := &sgc7game.Stake{
		CoinBet:  1,
		CashBet:  int64(bet),
		Currency: "EUR",
	}

	d := sgc7rtp.StartRTP2(game, rtp, icore, ispinnums, stake, int(ispinnums/100), ontimer, true, wincap)

	goutils.Info("finish.",
		slog.Int64("total nums", ispinnums),
		slog.Float64("rtp", float64(rtp.TotalWins)/float64(rtp.TotalBet)),
		slog.String("cost time", d.String()))

	components := game.Pool.mapComponents[int(bet)]
	components.Stats2.WaitEnding()

	return components.Stats2, nil
}

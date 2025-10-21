package main

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
)

func main() {
	goutils.InitLogger2("lowcodegame", sgc7ver.Version,
		"info", true, "./logs")

	gamecfg := os.Getenv("GAMECFG")
	strBet := os.Getenv("BET")
	strCheat := os.Getenv("CHEAT")

	lowcode.SetAllowStatsV2()

	game, err := lowcode.NewGame2(gamecfg, func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	}, lowcode.NewBasicRNG, lowcode.NewEmptyFeatureLevel)
	if err != nil {
		goutils.Error("NewGame2",
			slog.String("gamecfg", gamecfg),
			goutils.Err(err))

		return
	}

	bet := int64(game.Pool.Config.Bets[0])
	if strBet != "" {
		i64, _ := goutils.String2Int64(strBet)

		bet = i64
	}

	stake := &sgc7game.Stake{
		CoinBet:  1,
		CashBet:  bet,
		Currency: "EUR",
	}

	plugin := game.NewPlugin()
	defer game.FreePlugin(plugin)

	cmd := "SPIN"
	cmdparam := ""
	ps := game.Initialize()

	ret, err := lowcode.Spin(game, ps, plugin, stake, cmd, cmdparam, strCheat, false)
	if err != nil {
		goutils.Error("Spin",
			goutils.Err(err))

		return
	}

	for _, v := range ret {
		buf, err := sgc7game.PlayResult2JSON(v)
		if err != nil {
			goutils.Error("PlayResult2JSON",
				goutils.Err(err))

			return
		}

		goutils.Info("PlayResult", slog.String("result", string(buf)))
	}
}

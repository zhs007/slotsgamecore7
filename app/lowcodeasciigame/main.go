package main

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
)

func main() {
	goutils.InitLogger2("lowcodegame", sgc7ver.Version,
		"info", true, "./logs")

	gamecfg := os.Getenv("GAMECFG")
	strAutoSpin := os.Getenv("AUTOSPIN")
	strSkipGetChar := os.Getenv("SKIPGETCHAR")
	strBreakAtFeature := os.Getenv("BREAKATFEATURE")
	strBet := os.Getenv("BET")
	autospin, _ := goutils.String2Int64(strAutoSpin)

	isSkipGetChar := false
	if strSkipGetChar != "" {
		i64, _ := goutils.String2Int64(strSkipGetChar)

		isSkipGetChar = i64 > 0
	}

	isBreakAtFeature := false
	if strSkipGetChar != "" {
		i64, _ := goutils.String2Int64(strBreakAtFeature)

		isBreakAtFeature = i64 > 0
	}

	// lowcode.SetJsonMode()
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

	// lowcode.IsStatsComponentMsg = true

	asciigame.StartGame(game, stake, func(pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, gameData any) {
		gameProp, isok := gameData.(*lowcode.GameProperty)
		if !isok {
			return
		}

		game.OnAsciiGame(gameProp, stake, pr, lst)
	}, int(autospin), isSkipGetChar, isBreakAtFeature)

	// if game.Pool.Stats != nil {
	// 	game.Pool.Stats.Wait()
	// 	game.Pool.Stats.Root.SaveExcel("stats.xlsx")
	// }
}

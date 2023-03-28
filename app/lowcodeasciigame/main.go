package main

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

func main() {
	goutils.InitLogger("lowcodegame", sgc7ver.Version,
		"info", true, "./logs")

	gamecfg := os.Getenv("GAMECFG")
	strAutoSpin := os.Getenv("AUTOSPIN")
	strSkipGetChar := os.Getenv("SKIPGETCHAR")
	strBreakAtFeature := os.Getenv("BREAKATFEATURE")
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

	game, err := lowcode.NewGame(gamecfg)
	if err != nil {
		goutils.Error("NewGame",
			zap.String("gamecfg", gamecfg),
			zap.Error(err))

		return
	}

	stake := &sgc7game.Stake{
		CoinBet:  1,
		CashBet:  int64(game.Pool.Config.Bets[0]),
		Currency: "EUR",
	}

	asciigame.StartGame(game, stake, func(pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, gameData interface{}) {
		gameProp, isok := gameData.(*lowcode.GameProperty)
		if !isok {
			return
		}

		game.OnAsciiGame(gameProp, stake, pr, lst)
	}, int(autospin), isSkipGetChar, isBreakAtFeature)

	game.Pool.Stats.Wait()
	game.Pool.Stats.Root.SaveExcel("stats.xlsx")
}

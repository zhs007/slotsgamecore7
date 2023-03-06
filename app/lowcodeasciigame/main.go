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

	// mapSymbolColor := asciigame.NewSymbolColorMap(game.GetConfig().PayTables)
	// wColor := color.New(color.BgRed, color.FgHiWhite)
	// hColor := color.New(color.BgBlue, color.FgHiWhite)
	// mColor := color.New(color.BgGreen, color.FgHiWhite)
	// sColor := color.New(color.BgMagenta, color.FgHiWhite)
	// mapSymbolColor.AddSymbolColor(0, wColor)
	// mapSymbolColor.AddSymbolColor(1, hColor)
	// mapSymbolColor.AddSymbolColor(2, hColor)
	// mapSymbolColor.AddSymbolColor(3, mColor)
	// mapSymbolColor.AddSymbolColor(4, mColor)
	// mapSymbolColor.AddSymbolColor(10, sColor)

	stake := &sgc7game.Stake{
		CoinBet:  1,
		CashBet:  int64(game.Prop.Config.Bets[0]),
		Currency: "EUR",
	}

	asciigame.StartGame(game, stake, func(pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult) {
		game.OnAsciiGame(stake, pr, lst)
		// gp := pr.CurGameModParams.(*dtgame.GameParams)
		// if len(gp.ExpSyms) > 0 {
		// 	fmt.Printf("ExpandingSymbols trigger! [ ")
		// 	for _, v := range gp.ExpSyms {
		// 		fmt.Printf("%v ", mapSymbolColor.GetSymbolString(int(v)))
		// 	}
		// 	fmt.Printf("]\n")
		// }

		// asciigame.OutputScene("base", pr.Scenes[0], mapSymbolColor)

		// if len(gp.ExpSyms) > 0 {
		// 	asciigame.OutputResults("Wins in Normal", pr, func(r *sgc7game.Result) bool {
		// 		return r.Type != sgc7game.RTScatter
		// 	}, mapSymbolColor)

		// 	asciigame.OutputResults("Wins in Expanding Symbols", pr, func(r *sgc7game.Result) bool {
		// 		return r.Type == sgc7game.RTScatter
		// 	}, mapSymbolColor)
		// } else {
		// 	asciigame.OutputResults("", pr, func(r *sgc7game.Result) bool {
		// 		return true
		// 	}, mapSymbolColor)
		// }

		// if pr.IsFinish {
		// 	bgStats.OnResults(stake, lst)
		// }
	}, int(autospin), isSkipGetChar, isBreakAtFeature)

	game.Prop.Stats.SaveExcel("stats.xlsx")
}

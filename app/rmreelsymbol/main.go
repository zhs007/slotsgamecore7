package main

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
)

func main() {
	goutils.InitLogger2("rmreelsymbol", sgc7ver.Version,
		"info", true, "./logs")

	fn := os.Getenv("REELS")
	symbol := os.Getenv("SYMBOL")
	output := os.Getenv("OUTPUT")

	err := sgc7game.RemoveSymbolInReels(fn, output, symbol)
	if err != nil {
		goutils.Error("RemoveSymbolInReels",
			slog.String("fn", fn),
			goutils.Err(err))

		return
	}

	goutils.Info("Done!")
}

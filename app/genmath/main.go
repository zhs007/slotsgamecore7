package main

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/mathtoolset"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
)

func main() {
	goutils.InitLogger2("lowcodegame", sgc7ver.Version,
		"info", true, "./logs")

	cfgfn := os.Getenv("CFG")

	err := mathtoolset.GenMath(cfgfn)
	if err != nil {
		goutils.Error("GenMath",
			slog.String("cfgfn", cfgfn),
			goutils.Err(err))

		return
	}

	goutils.Info("Done!")
}

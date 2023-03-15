package main

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/mathtoolset"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

func main() {
	goutils.InitLogger("lowcodegame", sgc7ver.Version,
		"info", true, "./logs")

	cfgfn := os.Getenv("CFG")

	err := mathtoolset.GenMath(cfgfn)
	if err != nil {
		goutils.Error("GenMath",
			zap.String("cfgfn", cfgfn),
			zap.Error(err))

		return
	}

	goutils.Info("Done!")
}

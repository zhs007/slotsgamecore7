package main

import (
	"context"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/gamecollection"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

func main() {
	goutils.InitLogger("gamecollection", sgc7ver.Version,
		"debug", true, "./logs")

	serv, err := gamecollection.NewServ(":5000", sgc7ver.Version, false)
	if err != nil {
		goutils.Error("NewServ",
			zap.Error(err))

		return
	}

	lowcode.SetAllowForceOutcome(10000)

	serv.Start(context.Background())
}

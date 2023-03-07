package main

import (
	"fmt"
	"os"

	goutils "github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/lowcode"
	"github.com/zhs007/slotsgamecore7/simserv"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

func main() {
	gamecfg := os.Getenv("GAMECFG")

	cfg, err := simserv.LoadConfig("../data/simserv.yaml")
	if err != nil {
		fmt.Printf("LoadConfig(../data/simserv.yaml) fail.")

		return
	}

	goutils.InitLogger(cfg.GameCode, sgc7ver.Version,
		cfg.LogLevel, true, "./logs")

	game, err := lowcode.NewGame(gamecfg)
	if err != nil {
		goutils.Error("NewGame",
			zap.String("gamecfg", gamecfg),
			zap.Error(err))

		return
	}

	bs, err := NewSimService(game)
	if err != nil {
		goutils.Error("NewSimService",
			zap.Error(err))

		return
	}

	serv := simserv.NewServ(bs, cfg)

	goutils.Info(cfg.GameCode+" starting ...",
		zap.String("gameCode", cfg.GameCode),
		zap.String("version", sgc7ver.Version),
		zap.String("core version", sgc7ver.Version),
		zap.String("servAddr", cfg.BindAddr))

	serv.Start()
}

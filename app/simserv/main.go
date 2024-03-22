package main

import (
	"fmt"
	"log/slog"
	"os"

	goutils "github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/lowcode"
	"github.com/zhs007/slotsgamecore7/simserv"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
)

func main() {
	gamecfg := os.Getenv("GAMECFG")

	cfg, err := simserv.LoadConfig("../data/simserv.yaml")
	if err != nil {
		fmt.Printf("LoadConfig(../data/simserv.yaml) fail.")

		return
	}

	goutils.InitLogger2(cfg.GameCode, sgc7ver.Version,
		cfg.LogLevel, true, "./logs")

	game, err := lowcode.NewGame(gamecfg)
	if err != nil {
		goutils.Error("NewGame",
			slog.String("gamecfg", gamecfg),
			goutils.Err(err))

		return
	}

	bs, err := NewSimService(game)
	if err != nil {
		goutils.Error("NewSimService",
			goutils.Err(err))

		return
	}

	serv := simserv.NewServ(bs, cfg)

	goutils.Info(cfg.GameCode+" starting ...",
		slog.String("gameCode", cfg.GameCode),
		slog.String("version", sgc7ver.Version),
		slog.String("core version", sgc7ver.Version),
		slog.String("servAddr", cfg.BindAddr))

	serv.Start()
}

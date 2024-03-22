package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/gamecollection"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
)

func main() {
	goutils.InitLogger2("gamecollection", sgc7ver.Version,
		"info", true, "./logs")

	client, err := gamecollection.NewClient(":5000")
	if err != nil {
		goutils.Error("NewClient",
			goutils.Err(err))

		return
	}

	data, err := os.ReadFile("../data/game002.json")
	if err != nil {
		goutils.Error("ReadFile",
			goutils.Err(err))

		return
	}

	ret0, err := client.InitGame(context.Background(), "moonwalk", string(data))
	if err != nil {
		goutils.Error("InitGame",
			goutils.Err(err))

		return
	}

	goutils.Info("InitGame",
		slog.Any("ret", ret0))

	ret1, err := client.GetGameConfig(context.Background(), "moonwalk")
	if err != nil {
		goutils.Error("GetGameConfig",
			goutils.Err(err))

		return
	}

	goutils.Info("GetGameConfig",
		slog.Any("ret", ret1))

	ret2, err := client.InitializeGamePlayer(context.Background(), "moonwalk")
	if err != nil {
		goutils.Error("InitializeGamePlayer error",
			goutils.Err(err))

		return
	}

	goutils.Info("InitializeGamePlayer",
		slog.Any("ret", ret2))

	ret3, err := client.PlayGame(context.Background(), "moonwalk", ret2.PlayerState, "", &sgc7pb.Stake{
		CoinBet:  1,
		CashBet:  10,
		Currency: "EUR",
	}, "", "")
	if err != nil {
		goutils.Error("PlayGame error",
			goutils.Err(err))

		return
	}

	goutils.Info("PlayGame",
		slog.Any("ret", ret3))
}

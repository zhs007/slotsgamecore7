package main

import (
	"context"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/gamecollection"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
)

func main() {
	goutils.InitLogger("gamecollection", sgc7ver.Version,
		"info", true, "./logs")

	client, err := gamecollection.NewClient(":5000")
	if err != nil {
		goutils.Error("NewClient",
			zap.Error(err))

		return
	}

	data, err := os.ReadFile("../data/game002.json")
	if err != nil {
		goutils.Error("ReadFile",
			zap.Error(err))

		return
	}

	ret0, err := client.InitGame(context.Background(), "moonwalk", string(data))
	if err != nil {
		goutils.Error("InitGame",
			zap.Error(err))

		return
	}

	goutils.Info("InitGame",
		goutils.JSON("ret", ret0))

	ret1, err := client.GetGameConfig(context.Background(), "moonwalk")
	if err != nil {
		goutils.Error("GetGameConfig",
			zap.Error(err))

		return
	}

	goutils.Info("GetGameConfig",
		goutils.JSON("ret", ret1))

	ret2, err := client.InitializeGamePlayer(context.Background(), "moonwalk")
	if err != nil {
		goutils.Error("InitializeGamePlayer error",
			zap.Error(err))

		return
	}

	goutils.Info("InitializeGamePlayer",
		goutils.JSON("ret", ret2))

	ret3, err := client.PlayGame(context.Background(), "moonwalk", ret2.PlayerState, "", &sgc7pb.Stake{
		CoinBet:  1,
		CashBet:  10,
		Currency: "EUR",
	}, "", "")
	if err != nil {
		goutils.Error("PlayGame error",
			zap.Error(err))

		return
	}

	goutils.Info("PlayGame",
		goutils.JSON("ret", ret3))
}

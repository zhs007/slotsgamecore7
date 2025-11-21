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

	data, err := os.ReadFile("../data/game.json")
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

	ret3, err := client.PlayGame(context.Background(), "moonwalk", ret2.PlayerState, "140,163,42,78,79,6,19,16,0,10,6,7,11,12,21,16,0,9,9,13,16,16,15,0,0,12,7,6,7,23,22,17,12,7,14,8,10,8,5,16,7,13,22,14,1,21,16,15,23,19,11,15,17,2,17,12,13,8,20,20,7,13,15,18,15,15,18,7,15,4,1,19,3,13,0,8,8,23,0,8,18,1,4,15,13,13,2,5,5,15,13,20,5,4,20,23,9,7,2,21,19,14", &sgc7pb.Stake{
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

	ret4, err := client.PlayGame(context.Background(), "moonwalk", ret3.Play.PlayerState, "", &sgc7pb.Stake{
		CoinBet:  2,
		CashBet:  20,
		Currency: "EUR",
	}, "", "")
	if err != nil {
		goutils.Error("PlayGame error",
			goutils.Err(err))

		return
	}

	goutils.Info("PlayGame",
		slog.Any("ret", ret4))

	ret5, err := client.PlayGame(context.Background(), "moonwalk", ret4.Play.PlayerState, "", &sgc7pb.Stake{
		CoinBet:  3,
		CashBet:  30,
		Currency: "EUR",
	}, "", "")
	if err != nil {
		goutils.Error("PlayGame error",
			goutils.Err(err))

		return
	}

	ps := ret5.Play.PlayerState
	for range 100000 {
		curret, err := client.PlayGame(context.Background(), "moonwalk", ps, "", &sgc7pb.Stake{
			CoinBet:  3,
			CashBet:  30,
			Currency: "EUR",
		}, "", "")
		if err != nil {
			goutils.Error("PlayGame error",
				goutils.Err(err))

			return
		}

		ps = curret.Play.PlayerState
	}

	goutils.Info("PlayGame",
		slog.Any("ret", ret5))
}

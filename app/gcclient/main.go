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

	ret3, err := client.PlayGame(context.Background(), "moonwalk", ret2.PlayerState, "0,1,0,0,7,6,1,0,3,0,5,0,7,6,4,7,5,4,7,2,53,3,33,7,7,4,6,0,50,7,29,0,43,4,30,3,39,6,135,128,0,74,215,157,6,1,1,6,1,5,0,0,42,6,5,3,3,1,2,7,1,2,3,6,4,1,4,3,5,43,1,7,5", &sgc7pb.Stake{
		CoinBet:  1,
		CashBet:  20,
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

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

	ret3, err := client.PlayGame(context.Background(), "moonwalk", ret2.PlayerState, "124,12,74,59,165,59,14,9,17,4,14,6,0,16,0,19,12,12,14,6,33,35,0,166,165,132,15,86,145,3,7,14,7,6,23,20,18,7,21,2,8,6,16,1,19,17,20,21,22,7,31,18,34,25,27,25,31,26,30,33,18,15,17,10,164,33,15,26,94,16,7,2,16,1,14,14,15,23,7,4,23,19,3,19,3,36,14,34,21,32,3,6,29,9,24,0,48,120,149,153,94,168,19,9,1,12,10,21,2,12,2,22,8,10,5,10,9,26,36,14,14,5,6,21,12,12,31,10,33,36,2,22,30,4,13,20,21,6,63,81,137,153,14,3,3,15,4,8,0,5,17,20,3,16,12,9,13,30,9,13,32,12,5,9,30,18,8,27,36,26,9,0,51,153,64,162,115,9,6,22,17,14,21,20,17,22,19,23,9,13,8,5,10,17,32,7,20,25,25,4,27,32,28,1,3,5,30,1,33,15,22,5,30,160,133,69,86,110,139,13,0,9,21,15,7,20,20,2,22,15,21,1,11,8,24,11,26,10,27,1,26,4,11,5,28,25,26,2,30,36,9,26,109,124,16,129,105,205,3,5,19,0,8,11,11,22,14,3,19,15,18,4,6,26,24,3,29,6,8,6,33,8,25,7,29,12,13,17,13,35,26,32", &sgc7pb.Stake{
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

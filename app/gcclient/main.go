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

	ret3, err := client.PlayGame(context.Background(), "moonwalk", ret2.PlayerState, "36,59,31,23,27,47,205,523,363,147,56,407,505,196,159,385,343,144,215,546,362,420,540,309,351,447,264,64,283,483,475,83,61,357,323,428,2947,2844,1200,1541,720,176,207,1423,2583,372,2924,1082,1755,1629,2052,3184,1505,324,466,2644,1587,2263,2225,2953,3282,2695,616,1280,103,4615,6631,6012,5626,1785,2164,4786,9273,1632,1908,2267,7652,5470,9231,183,6781,1251,6687,9033,3978,3004,1235,1139,51,209,2066,2851,1035,1873,1449,923,1976,640,1664,1813,1741,6678,2149,1762,1127,9923,7682,5223,7998,4533,4358,70428,19641,48162,24766,31967,23028,17890,35417,54901,76417,8839,3302,1322,238,9801,3500,9485,9677,3807,9858", &sgc7pb.Stake{
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

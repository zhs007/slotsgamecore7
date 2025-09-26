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

	ret3, err := client.PlayGame(context.Background(), "moonwalk", ret2.PlayerState, "53, 286, 573, 281, 327, 3, 524, 532, 193, 221, 548, 121, 8, 1, 647, 302, 527, 154, 114, 568, 472, 541, 170, 164, 110, 492, 52, 65, 379, 468, 285, 198, 85, 284, 426, 344, 512, 292, 270, 40, 99, 53, 144, 14, 70, 78, 81, 49, 87, 26, 25, 46, 44, 15, 5, 15, 4, 1, 3, 1, 1, 0, 0, 0, 301, 329, 248, 170, 91, 89, 114, 333, 202, 12, 30, 368, 316, 40, 409, 284, 3, 360, 373, 218, 312, 164, 462, 248, 170, 262, 275, 297, 138, 294, 104, 55, 341, 414, 15, 183, 114, 403, 54, 432, 335, 227, 399, 259, 397, 228, 53, 197, 203, 178, 327, 329, 302, 427, 190, 187, 150, 191, 383, 43, 85, 17, 49, 17, 379, 98, 145, 234, 332, 19, 155, 148, 318, 219, 143, 10, 193, 445, 23, 286, 60, 253, 279, 206, 164, 214, 385, 336, 174, 437, 282, 242, 343, 368, 387, 172, 382, 119, 363, 49, 59, 20, 431, 12, 135, 16, 226, 32, 224, 252, 255, 291, 81, 337, 227, 3, 284, 385, 148, 104, 195, 66, 324, 449, 322, 45, 31, 40, 267, 230, 26, 387, 97, 303, 443, 168, 13, 92, 370, 150, 132, 110, 354, 470, 67, 298, 144, 35, 427, 101, 243, 125, 331, 143, 281, 269, 114, 82, 36, 167", &sgc7pb.Stake{
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

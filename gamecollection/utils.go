package gamecollection

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/bytedance/sonic"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/grpcserv"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
)

func Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ProcCheat - process cheat
func ProcCheat(plugin sgc7plugin.IPlugin, cheat string) error {
	if cheat != "" {
		if sgc7game.IsRngString(cheat) {
			str := goutils.AppendString("[", cheat, "]")

			rngs := []int{}
			err := sonic.Unmarshal([]byte(str), &rngs)
			if err != nil {
				return err
			}

			plugin.SetCache(rngs)
		}
	}

	return nil
}

// addWinResult - add sgc7game.PlayResult
func addWinResult(sv grpcserv.IService, pr *sgc7pb.ReplyPlay, playResult *sgc7game.PlayResult) error {
	r := &sgc7pb.GameResult{
		CoinWin: int64(playResult.CoinWin),
		ClientData: &sgc7pb.PlayResult{
			CurGameMod:  playResult.CurGameMod,
			NextGameMod: playResult.NextGameMod,
			CurIndex:    int32(playResult.CurIndex),
			ParentIndex: int32(playResult.ParentIndex),
			ModType:     playResult.ModType,
		},
	}

	gp, err := sv.BuildPBGameModParam(playResult.CurGameModParams)
	if err != nil {
		goutils.Error("addWinResult:BuildPBGameModParam",
			goutils.Err(err))

		return err
	}

	r.ClientData.CurGameModParam = gp

	for _, v := range playResult.Scenes {
		cs := sgc7pbutils.BuildPBGameScene(v)
		r.ClientData.Scenes = append(r.ClientData.Scenes, cs)
	}

	for _, v := range playResult.OtherScenes {
		cs := sgc7pbutils.BuildPBGameScene(v)
		r.ClientData.OtherScenes = append(r.ClientData.OtherScenes, cs)
	}

	for _, v := range playResult.PrizeScenes {
		cs := sgc7pbutils.BuildPBGameScene(v)
		r.ClientData.PrizeScenes = append(r.ClientData.PrizeScenes, cs)
	}

	r.ClientData.PrizeCoinWin = int64(playResult.PrizeCoinWin)
	r.ClientData.PrizeCashWin = playResult.PrizeCashWin

	r.ClientData.JackpotCoinWin = int64(playResult.JackpotCoinWin)
	r.ClientData.JackpotCashWin = playResult.JackpotCashWin
	r.ClientData.JackpotType = int32(playResult.JackpotType)

	for _, v := range playResult.Results {
		cr := sgc7pbutils.BuildPBGameScenePlayResult(v)
		r.ClientData.Results = append(r.ClientData.Results, cr)
	}

	for _, v := range playResult.MulPos {
		r.ClientData.MulPos = append(r.ClientData.MulPos, int32(v))
	}

	r.CashWin = playResult.CashWin

	pr.Results = append(pr.Results, r)

	return nil
}

// AddPlayResult - []*sgc7game.PlayResult => *PlayResult
func AddPlayResult(sv grpcserv.IService, pr *sgc7pb.ReplyPlay, results []*sgc7game.PlayResult) {
	for _, v := range results {
		addWinResult(sv, pr, v)
	}
}

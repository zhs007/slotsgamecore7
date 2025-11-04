package grpcserv

import (
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
)

// addWinResult - add sgc7game.PlayResult
func addWinResult(sv IService, pr *sgc7pb.ReplyPlay, playResult *sgc7game.PlayResult) error {
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

	// SPGrid - map[string][]*GameScene -> map[string]*SPGridList
	if len(playResult.SPGrid) > 0 {
		if r.ClientData.SpGrid == nil {
			r.ClientData.SpGrid = make(map[string]*sgc7pb.SPGridList)
		}

		for k, spgrid := range playResult.SPGrid {
			pbl := &sgc7pb.SPGridList{}
			for _, gs := range spgrid.Grid {
				pbl.Scenes = append(pbl.Scenes, sgc7pbutils.BuildPBGameScene(gs))
			}

			r.ClientData.SpGrid[k] = pbl
		}
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
func AddPlayResult(sv IService, pr *sgc7pb.ReplyPlay, results []*sgc7game.PlayResult) {
	for _, v := range results {
		addWinResult(sv, pr, v)
	}
}

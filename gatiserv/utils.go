package gatiserv

import (
	jsoniter "github.com/json-iterator/go"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// BuildIPlayerState - PlayerState => sgc7game.IPlayerState
func BuildIPlayerState(ips sgc7game.IPlayerState, ps PlayerState) error {
	err := ips.SetPublicString(ps.Public)
	if err != nil {
		return err
	}

	err = ips.SetPrivateString(ps.Private)
	if err != nil {
		return err
	}

	return nil
}

// BuildStake - PlayerState => sgc7game.IPlayerState
func BuildStake(stake Stake) *sgc7game.Stake {
	return &sgc7game.Stake{
		CoinBet:  int64(stake.CoinBet * 100),
		CashBet:  int64(stake.CashBet * 100),
		Currency: stake.Currency,
	}
}

// AddWinResult - add sgc7game.PlayResult
func AddWinResult(pr *PlayResult, stake Stake, playResult *sgc7game.PlayResult) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	prb, err := json.Marshal(playResult)
	if err != nil {
		sgc7utils.Warn("gatiserv.AddWinResult:Marshal PlayResult",
			zap.Error(err))

		return err
	}

	r := &Result{
		CoinWin:    playResult.CoinWin,
		ClientData: string(prb),
	}

	r.CashWin = float64(playResult.CashWin) / 100.0

	pr.Results = append(pr.Results, r)

	return nil
}

// AddPlayResult - []*sgc7game.PlayResult => *PlayResult
func AddPlayResult(pr *PlayResult, stake Stake, results []*sgc7game.PlayResult) {

	for _, v := range results {
		AddWinResult(pr, stake, v)
	}
}

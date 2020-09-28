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

// BuildPlayerStateString - sgc7game.IPlayerState => string
func BuildPlayerStateString(ps sgc7game.IPlayerState) (string, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	if ps == nil {
		return "{\"playerStatePublic\":\"{}\",\"playerStatePrivate\":\"{}\"}", nil
	}

	dps, err := BuildPlayerState(ps)
	if err != nil {
		sgc7utils.Warn("gatiserv.BuildPlayerStateString:BuildPlayerState",
			zap.Error(err))

		return "", err
	}

	psfb, err := json.Marshal(dps)
	if err != nil {
		sgc7utils.Warn("gatiserv.BuildPlayerStateString:Marshal PlayerState",
			zap.Error(err))

		return "", err
	}

	return string(psfb), nil
}

// BuildPlayerState - sgc7game.IPlayerState => PlayerState
func BuildPlayerState(ips sgc7game.IPlayerState) (*PlayerState, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	if ips == nil {
		return nil, nil
	}

	psb, err := json.Marshal(ips.GetPublic())
	if err != nil {
		sgc7utils.Warn("gatiserv.BuildPlayerState:Marshal GetPublic",
			zap.Error(err))

		return nil, err
	}

	psp, err := json.Marshal(ips.GetPrivate())
	if err != nil {
		sgc7utils.Warn("gatiserv.BuildPlayerState:Marshal GetPrivate",
			zap.Error(err))

		return nil, err
	}

	return &PlayerState{
		Public:  string(psb),
		Private: string(psp),
	}, nil
}

// ParsePlayerState - json => PlayerState
func ParsePlayerState(str string) (*PlayerState, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	ps := &PlayerState{}
	err := json.Unmarshal([]byte(str), ps)
	if err != nil {
		sgc7utils.Error("gatiserv.ParsePlayerState:JSON",
			zap.String("str", str),
			zap.Error(err))

		return nil, err
	}

	return ps, nil
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
	// json := jsoniter.ConfigCompatibleWithStandardLibrary

	// prb, err := json.Marshal(playResult)
	// if err != nil {
	// 	sgc7utils.Warn("gatiserv.AddWinResult:Marshal PlayResult",
	// 		zap.Error(err))

	// 	return err
	// }

	r := &Result{
		CoinWin:    playResult.CoinWin,
		ClientData: playResult,
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

// ParsePlayParams - string => *PlayParams
func ParsePlayParams(str string) (*PlayParams, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	ps := &PlayParams{}
	err := json.Unmarshal([]byte(str), ps)
	if err != nil {
		sgc7utils.Error("gatiserv.ParsePlayParams:JSON",
			zap.String("str", str),
			zap.Error(err))

		return nil, err
	}

	return ps, nil
}

// ParsePlayResult - string => *PlayResult
func ParsePlayResult(str string) (*PlayResult, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	pr := &PlayResult{}
	err := json.Unmarshal([]byte(str), pr)
	if err != nil {
		sgc7utils.Error("gatiserv.ParsePlayResult:JSON",
			zap.String("str", str),
			zap.Error(err))

		return nil, err
	}

	return pr, nil
}

// NewGATIGameInfo -
func NewGATIGameInfo() *GATIGameInfo {
	return &GATIGameInfo{
		Components: make(map[int]*GATICriticalComponent),
	}
}

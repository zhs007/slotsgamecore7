package gatiserv

import (
	"os"

	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

// BuildIPlayerState - PlayerState => sgc7game.IPlayerState
func BuildIPlayerState(ips sgc7game.IPlayerState, ps *PlayerState) error {
	err := ips.SetPublic(ps.Public)
	if err != nil {
		return err
	}

	err = ips.SetPrivate(ps.Private)
	if err != nil {
		return err
	}

	return nil
}

// BuildPlayerStateString - sgc7game.IPlayerState => string
func BuildPlayerStateString(ps sgc7game.IPlayerState) (string, error) {
	if ps == nil {
		return "{\"playerStatePublic\":{},\"playerStatePrivate\":{}}", nil
	}

	dps, err := BuildPlayerState(ps)
	if err != nil {
		goutils.Warn("gatiserv.BuildPlayerStateString:BuildPlayerState",
			zap.Error(err))

		return "", err
	}

	psfb, err := sonic.Marshal(dps)
	if err != nil {
		goutils.Warn("gatiserv.BuildPlayerStateString:Marshal PlayerState",
			zap.Error(err))

		return "", err
	}

	return string(psfb), nil
}

// BuildPlayerState - sgc7game.IPlayerState => PlayerState
func BuildPlayerState(ips sgc7game.IPlayerState) (*PlayerState, error) {
	if ips == nil {
		return nil, nil
	}

	return &PlayerState{
		Public:  ips.GetPublic(),
		Private: ips.GetPrivate(),
	}, nil
}

// ParsePlayerState - json => PlayerState
func ParsePlayerState(str string, ps *PlayerState) error {
	err := sonic.Unmarshal([]byte(str), ps)
	if err != nil {
		goutils.Error("gatiserv.ParsePlayerState:JSON",
			zap.String("str", str),
			zap.Error(err))

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
func ParsePlayParams(str string, ps *PlayerState) (*PlayParams, error) {
	pp := &PlayParams{
		PlayerState: ps,
	}
	err := sonic.Unmarshal([]byte(str), pp)
	if err != nil {
		goutils.Error("gatiserv.ParsePlayParams:JSON",
			zap.String("str", str),
			zap.Error(err))

		return nil, err
	}

	return pp, nil
}

// ParsePlayResult - string => *PlayResult
func ParsePlayResult(str string) (*PlayResult, error) {
	pr := &PlayResult{}
	err := sonic.Unmarshal([]byte(str), pr)
	if err != nil {
		goutils.Error("gatiserv.ParsePlayResult:JSON",
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

// LoadGATIGameConfig - load
func LoadGATIGameConfig(fn string) (*GATIGameConfig, error) {
	if fn == "" {
		return &GATIGameConfig{}, nil
	}

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	ccs := &GATIGameConfig{}
	err = sonic.Unmarshal(data, ccs)
	if err != nil {
		goutils.Warn("gatiserv.LoadGATIGameConfig",
			zap.Error(err))

		return nil, err
	}

	return ccs, nil
}

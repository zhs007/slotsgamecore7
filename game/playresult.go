package sgc7game

import (
	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
)

type SPGrid struct {
	Width  int          `json:"width"`
	Height int          `json:"height"`
	Grid   []*GameScene `json:"grid"`
}

// PlayResult - result for play
type PlayResult struct {
	CurGameMod       string             `json:"curgamemod"`
	CurGameModParams any                `json:"curgamemodparams"`
	NextGameMod      string             `json:"nextgamemod"`
	Scenes           []*GameScene       `json:"scenes"`
	OtherScenes      []*GameScene       `json:"otherscenes"`
	PrizeScenes      []*GameScene       `json:"prizescenes"`
	PrizeCoinWin     int                `json:"prizecoinwin"`
	PrizeCashWin     int64              `json:"prizecashwin"`
	JackpotCoinWin   int                `json:"jackpotcoinwin"`
	JackpotCashWin   int64              `json:"jackpotcashwin"`
	JackpotType      int                `json:"jackpottype"`
	Results          []*Result          `json:"results"`
	MulPos           []int              `json:"mulpos"`
	NextCmds         []string           `json:"-"`
	NextCmdParams    []string           `json:"-"`
	CoinWin          int                `json:"-"`
	CashWin          int64              `json:"-"`
	IsFinish         bool               `json:"-"`
	IsWait           bool               `json:"-"`
	CurIndex         int                `json:"-"`
	ParentIndex      int                `json:"-"`
	ModType          string             `json:"-"`
	SPGrid           map[string]*SPGrid `json:"spgrid"`
}

// NewPlayResult - new a PlayResult
func NewPlayResult(curGameMod string, curIndex int, parentIndex int, modType string) *PlayResult {
	return &PlayResult{
		CurGameMod:  curGameMod,
		CurIndex:    curIndex,
		ParentIndex: parentIndex,
		ModType:     modType,
		SPGrid:      make(map[string]*SPGrid),
	}
}

// GetPlayResultCurIndex - get current index
func GetPlayResultCurIndex(prs []*PlayResult) int {
	if len(prs) == 0 {
		return 0
	}

	return prs[len(prs)-1].CurIndex + 1
}

// PlayResult2JSON - PlayResult => json
func PlayResult2JSON(pr *PlayResult) ([]byte, error) {
	b, err := sonic.Marshal(pr)
	if err != nil {
		goutils.Warn("sgc7game.PlayResult2JSON",
			goutils.Err(err))

		return nil, err
	}

	return b, nil
}

// JSON2PlayResult - json => PlayResult
func JSON2PlayResult(buf []byte, pr *PlayResult) (*PlayResult, error) {
	err := sonic.Unmarshal(buf, &pr)
	if err != nil {
		goutils.Warn("sgc7game.JSON2PlayResult",
			goutils.Err(err))
		return nil, err
	}

	return pr, nil
}

// CountEndingSymbols - count symbol number
func (pr *PlayResult) CountEndingSymbols(symbols []int) []int {
	if len(pr.Scenes) > 0 {
		cs := pr.Scenes[len(pr.Scenes)-1]

		return cs.CountSymbols(symbols)
	}

	return nil
}

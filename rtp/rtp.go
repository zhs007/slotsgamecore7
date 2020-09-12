package sgc7rtp

import (
	"strconv"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// FuncOnResult - onResult(*RTPNode, *sgc7game.PlayResult)
type FuncOnResult func(node *RTPNode, pr *sgc7game.PlayResult)

// OnRootResult - on root
func OnRootResult(node *RTPNode, pr *sgc7game.PlayResult) {
	if pr.CashWin > 0 {
		if pr.IsFinish {
			node.TriggerNums++
		}

		node.TotalWin += pr.CashWin
	}
}

// OnGameModResult - on gamemod
func OnGameModResult(node *RTPNode, pr *sgc7game.PlayResult) {
	if pr.CurGameMod == node.GameMod && pr.CashWin > 0 {
		node.TriggerNums++
		node.TotalWin += pr.CashWin
	}
}

// OnSymbolResult - on symbol
func OnSymbolResult(node *RTPNode, pr *sgc7game.PlayResult) {
	if pr.CurGameMod == node.GameMod && pr.CashWin > 0 {
		for _, v := range pr.Results {
			if v.Symbol == node.Symbol {
				node.TriggerNums++
				node.TotalWin += int64(v.CashWin)
			}
		}
	}
}

// OnSymbolNumsResult - on symbol nums
func OnSymbolNumsResult(node *RTPNode, pr *sgc7game.PlayResult) {
	if pr.CurGameMod == node.GameMod && pr.CashWin > 0 {
		for _, v := range pr.Results {
			if v.Symbol == node.Symbol && len(v.Pos) == node.SymbolNums*2 {
				node.TriggerNums++
				node.TotalWin += int64(v.CashWin)
			}
		}
	}
}

// RTPNode -
type RTPNode struct {
	TriggerNums  int64
	TotalWin     int64
	RTP          float64
	MapChildren  map[string]*RTPNode
	GameMod      string
	Symbol       int
	SymbolNums   int
	TagName      string
	funcOnResult FuncOnResult
}

// NewRTPRoot - new RTPNode
func NewRTPRoot() *RTPNode {
	return &RTPNode{
		MapChildren:  make(map[string]*RTPNode),
		funcOnResult: OnRootResult,
	}
}

// NewRTPGameMod - new RTPNode
func NewRTPGameMod(gamemod string) *RTPNode {
	return &RTPNode{
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		funcOnResult: OnGameModResult,
	}
}

// InitGameMod - new RTPNode
func InitGameMod(node *RTPNode, symbols []int, nums []int) {
	for _, sv := range symbols {
		csn := NewRTPSymbol(node.GameMod, sv)
		InitSymbol(csn, sv, nums)
		node.AddChild(strconv.Itoa(sv), csn)
	}
}

// InitSymbol - new RTPNode
func InitSymbol(node *RTPNode, symbol int, nums []int) {
	for _, nv := range nums {
		csnn := NewRTPSymbolNums(node.GameMod, symbol, nv)
		node.AddChild(strconv.Itoa(nv), csnn)
	}
}

// NewRTPSymbol - new RTPNode
func NewRTPSymbol(gamemod string, symbol int) *RTPNode {
	return &RTPNode{
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		Symbol:       symbol,
		funcOnResult: OnSymbolResult,
	}
}

// NewRTPSymbolNums - new RTPNode
func NewRTPSymbolNums(gamemod string, symbol int, nums int) *RTPNode {
	return &RTPNode{
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		Symbol:       symbol,
		SymbolNums:   nums,
		funcOnResult: OnSymbolNumsResult,
	}
}

// CalcRTP -
func (node *RTPNode) CalcRTP(totalbet int64) {
	node.RTP = float64(node.TotalWin) / float64(totalbet)

	for _, v := range node.MapChildren {
		v.CalcRTP(totalbet)
	}
}

// AddChild -
func (node *RTPNode) AddChild(name string, c *RTPNode) {
	node.MapChildren[name] = c
}

// OnResult -
func (node *RTPNode) OnResult(pr *sgc7game.PlayResult) {
	node.funcOnResult(node, pr)

	for _, v := range node.MapChildren {
		v.OnResult(pr)
	}
}

// RTP -
type RTP struct {
	BetNums  int64
	TotalBet int64
	Root     *RTPNode
}

// NewRTP - new RTP
func NewRTP() *RTP {
	return &RTP{
		Root: NewRTPRoot(),
	}
}

// CalcRTP -
func (rtp *RTP) CalcRTP() {
	rtp.Root.CalcRTP(rtp.TotalBet)
}

// Bet -
func (rtp *RTP) Bet(bet int64) {
	rtp.BetNums++
	rtp.TotalBet += bet
}

// OnResult -
func (rtp *RTP) OnResult(pr *sgc7game.PlayResult) {
	rtp.Root.OnResult(pr)
}

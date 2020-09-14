package sgc7rtp

import (
	"os"
	"sort"
	"strconv"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
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
		Symbol:       -1,
		SymbolNums:   -1,
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
		SymbolNums:   -1,
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

// GetSymbolNums -
func (node *RTPNode) GetSymbolNums(arr []int) []int {
	if node.SymbolNums > 0 {
		if sgc7utils.IndexOfIntSlice(arr, node.SymbolNums, 0) < 0 {
			arr = append(arr, node.SymbolNums)
		}
	}

	for _, v := range node.MapChildren {
		arr = v.GetSymbolNums(arr)
	}

	return arr
}

// GetGameMods -
func (node *RTPNode) GetGameMods(arr []string) []string {
	if node.GameMod != "" {
		if sgc7utils.IndexOfStringSlice(arr, node.GameMod, 0) < 0 {
			arr = append(arr, node.GameMod)
		}
	}

	for _, v := range node.MapChildren {
		arr = v.GetGameMods(arr)
	}

	return arr
}

// GetSymbols -
func (node *RTPNode) GetSymbols(arr []int) []int {
	if node.Symbol >= 0 {
		if sgc7utils.IndexOfIntSlice(arr, node.Symbol, 0) < 0 {
			arr = append(arr, node.Symbol)
		}
	}

	for _, v := range node.MapChildren {
		arr = v.GetSymbols(arr)
	}

	return arr
}

// GenRootString -
func (node *RTPNode) GenRootString(sn []int, totalbet int64) string {
	if node.GameMod == "" {
		str := sgc7utils.AppendString(",,", strconv.FormatInt(totalbet, 10))
		for range sn {
			str = sgc7utils.AppendString(str, ",")
		}
		str = sgc7utils.AppendString(str, ",", strconv.FormatInt(node.TotalWin, 10), "\n")

		return str
	}

	return ""
}

// GenGameModString -
func (node *RTPNode) GenGameModString(gamemod string, sn []int, totalbet int64) string {
	if node.GameMod == "" {
		for _, v := range node.MapChildren {
			str := v.GenGameModString(gamemod, sn, totalbet)
			if str != "" {
				return str
			}
		}

		return ""
	}

	if node.GameMod == gamemod {
		if node.Symbol < 0 && node.SymbolNums < 0 {
			str := sgc7utils.AppendString(node.GameMod, ",,", strconv.FormatInt(totalbet, 10))
			for range sn {
				str = sgc7utils.AppendString(str, ",")
			}
			str = sgc7utils.AppendString(str, ",", strconv.FormatInt(node.TotalWin, 10), "\n")

			return str
		}
	}

	return ""
}

// GenSymbolString -
func (node *RTPNode) GenSymbolString(gamemod string, symbol int, sn []int, totalbet int64) string {
	if node.GameMod == "" || (node.GameMod == gamemod && node.Symbol < 0 && node.SymbolNums < 0) {
		for _, v := range node.MapChildren {
			str := v.GenSymbolString(gamemod, symbol, sn, totalbet)
			if str != "" {
				return str
			}
		}

		return ""
	}

	if node.GameMod == gamemod && node.Symbol == symbol && node.SymbolNums < 0 {
		str := sgc7utils.AppendString(node.GameMod, ",", strconv.Itoa(symbol), ",", strconv.FormatInt(totalbet, 10))
		for _, v := range sn {
			won := node.GetSymbolNumsWon(gamemod, symbol, v)
			if won < 0 {
				str = sgc7utils.AppendString(str, ",")
			} else {
				str = sgc7utils.AppendString(str, ",", strconv.FormatInt(won, 10))
			}
		}

		str = sgc7utils.AppendString(str, ",", strconv.FormatInt(node.TotalWin, 10), "\n")

		return str
	}

	return ""
}

// GetSymbolNumsWon -
func (node *RTPNode) GetSymbolNumsWon(gamemod string, symbol int, sn int) int64 {
	if node.GameMod == "" || (node.GameMod == gamemod && node.Symbol < 0 && node.SymbolNums < 0) || (node.GameMod == gamemod && node.Symbol == symbol && node.SymbolNums < 0) {
		for _, v := range node.MapChildren {
			won := v.GetSymbolNumsWon(gamemod, symbol, sn)
			if won >= 0 {
				return won
			}
		}

		return -1
	}

	if node.GameMod == gamemod && node.Symbol == symbol && node.SymbolNums == sn {
		return node.TotalWin
	}

	return -1
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

// Save2CSV -
func (rtp *RTP) Save2CSV(fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		sgc7utils.Error("sgc7rtp.RTP.Save2CSV",
			zap.Error(err))

		return err
	}
	defer f.Close()

	gms := rtp.Root.GetGameMods(nil)
	sn := rtp.Root.GetSymbolNums(nil)
	symbols := rtp.Root.GetSymbols(nil)

	sort.Slice(sn, func(i, j int) bool {
		return sn[i] < sn[j]
	})

	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i] < symbols[j]
	})

	sort.Slice(gms, func(i, j int) bool {
		return gms[i] < gms[j]
	})

	strhead := "gamemod,symbol,totalbet"
	for _, v := range sn {
		strhead = sgc7utils.AppendString(strhead, ",X", strconv.Itoa(v))
	}
	strhead = sgc7utils.AppendString(strhead, ",totalwin\n")

	f.WriteString(strhead)

	for _, gm := range gms {
		for _, symbol := range symbols {
			str := rtp.Root.GenSymbolString(gm, symbol, sn, rtp.TotalBet)
			f.WriteString(str)
		}

		str := rtp.Root.GenGameModString(gm, sn, rtp.TotalBet)
		f.WriteString(str)
	}

	str := rtp.Root.GenRootString(sn, rtp.TotalBet)
	f.WriteString(str)

	f.Sync()

	return nil
}

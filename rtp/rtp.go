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
type FuncOnResult func(node *RTPNode, pr *sgc7game.PlayResult) bool

// OnRootResult - on root
func OnRootResult(node *RTPNode, pr *sgc7game.PlayResult) bool {
	if pr.CashWin > 0 {
		if pr.IsFinish {
			node.TriggerNums++
		}

		node.TotalWin += pr.CashWin
	}

	return true
}

// OnGameModResult - on gamemod
func OnGameModResult(node *RTPNode, pr *sgc7game.PlayResult) bool {
	if pr.CurGameMod == node.GameMod {
		if pr.CashWin > 0 {
			node.TriggerNums++
			node.TotalWin += pr.CashWin
		}

		return true
	}

	return false
}

// OnSymbolResult - on symbol
func OnSymbolResult(node *RTPNode, pr *sgc7game.PlayResult) bool {
	// if pr.CurGameMod == node.GameMod {
	if pr.CashWin > 0 {
		for _, v := range pr.Results {
			if v.Symbol == node.Symbol {
				node.TriggerNums++
				node.TotalWin += int64(v.CashWin)
			}
		}
	}

	return true
	// }

	// return false
}

// OnSymbolNumsResult - on symbol nums
func OnSymbolNumsResult(node *RTPNode, pr *sgc7game.PlayResult) bool {
	if pr.CurGameMod == node.GameMod && pr.CashWin > 0 {
		for _, v := range pr.Results {
			if v.Symbol == node.Symbol && len(v.Pos) == node.SymbolNums*2 {
				node.TriggerNums++
				node.TotalWin += int64(v.CashWin)
			}
		}
	}

	return false
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

// NewRTPGameModTag - new RTPNode
func NewRTPGameModTag(gamemod string, tag string, funcTag FuncOnResult) *RTPNode {
	return &RTPNode{
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		funcOnResult: funcTag,
		TagName:      tag,
		Symbol:       -1,
		SymbolNums:   -1,
	}
}

// InitGameMod - new RTPNode
func InitGameMod(node *RTPNode, tags []string, funcTag []FuncOnResult, symbols []int, nums []int) {
	if len(tags) > 0 {
		for i, tag := range tags {
			ctn := NewRTPGameModTag(node.GameMod, tag, funcTag[i])
			InitGameModTag(ctn, tag, symbols, nums)
			node.AddChild(tag, ctn)
		}
	} else {
		for _, sv := range symbols {
			csn := NewRTPSymbol(node.GameMod, "", sv)
			InitSymbol(csn, "", sv, nums)
			node.AddChild(strconv.Itoa(sv), csn)
		}
	}
}

// InitGameModTag - new RTPNode
func InitGameModTag(node *RTPNode, tag string, symbols []int, nums []int) {
	for _, sv := range symbols {
		csn := NewRTPSymbol(node.GameMod, tag, sv)
		InitSymbol(csn, tag, sv, nums)
		node.AddChild(strconv.Itoa(sv), csn)
	}
}

// InitSymbol - new RTPNode
func InitSymbol(node *RTPNode, tag string, symbol int, nums []int) {
	for _, nv := range nums {
		csnn := NewRTPSymbolNums(node.GameMod, tag, symbol, nv)
		node.AddChild(strconv.Itoa(nv), csnn)
	}
}

// NewRTPSymbol - new RTPNode
func NewRTPSymbol(gamemod string, tag string, symbol int) *RTPNode {
	return &RTPNode{
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		Symbol:       symbol,
		funcOnResult: OnSymbolResult,
		TagName:      tag,
		SymbolNums:   -1,
	}
}

// NewRTPSymbolNums - new RTPNode
func NewRTPSymbolNums(gamemod string, tag string, symbol int, nums int) *RTPNode {
	return &RTPNode{
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		Symbol:       symbol,
		SymbolNums:   nums,
		funcOnResult: OnSymbolNumsResult,
		TagName:      tag,
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
	ismine := node.funcOnResult(node, pr)
	if !ismine {
		return
	}

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

// GetTags -
func (node *RTPNode) GetTags(arr []string) []string {
	if node.TagName != "" {
		if sgc7utils.IndexOfStringSlice(arr, node.TagName, 0) < 0 {
			arr = append(arr, node.TagName)
		}
	}

	for _, v := range node.MapChildren {
		arr = v.GetTags(arr)
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
			str := sgc7utils.AppendString(node.GameMod, ",,,", strconv.FormatInt(totalbet, 10))
			for range sn {
				str = sgc7utils.AppendString(str, ",")
			}
			str = sgc7utils.AppendString(str, ",", strconv.FormatInt(node.TotalWin, 10), "\n")

			return str
		}
	}

	return ""
}

// GenTagString -
func (node *RTPNode) GenTagString(gamemod string, tag string, sn []int, totalbet int64) string {
	if node.GameMod == "" ||
		(node.TagName == tag && node.GameMod == gamemod) {

		for _, v := range node.MapChildren {
			str := v.GenTagString(gamemod, tag, sn, totalbet)
			if str != "" {
				return str
			}
		}

		return ""
	}

	if node.GameMod == gamemod {
		if node.Symbol < 0 && node.SymbolNums < 0 {
			str := sgc7utils.AppendString(node.GameMod, ",", tag, ",,", strconv.FormatInt(totalbet, 10))
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
func (node *RTPNode) GenSymbolString(gamemod string, tag string, symbol int, sn []int, totalbet int64) string {
	if node.GameMod == "" ||
		(node.GameMod == gamemod && (node.TagName == tag || node.TagName == "") && node.Symbol < 0 && node.SymbolNums < 0) {

		for _, v := range node.MapChildren {
			str := v.GenSymbolString(gamemod, tag, symbol, sn, totalbet)
			if str != "" {
				return str
			}
		}

		return ""
	}

	if node.GameMod == gamemod && node.TagName == tag && node.Symbol == symbol && node.SymbolNums < 0 {
		str := sgc7utils.AppendString(node.GameMod, ",", tag, ",", strconv.Itoa(symbol), ",", strconv.FormatInt(totalbet, 10))
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
	if node.GameMod == "" ||
		(node.GameMod == gamemod && node.Symbol < 0 && node.SymbolNums < 0) ||
		(node.GameMod == gamemod && node.Symbol == symbol && node.SymbolNums < 0) {

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

// ChgSymbolNumsFunc -
func (node *RTPNode) ChgSymbolNumsFunc(funcOnResult FuncOnResult) {
	if node.Symbol >= 0 && node.SymbolNums > 0 {
		node.funcOnResult = funcOnResult
	}

	for _, v := range node.MapChildren {
		v.ChgSymbolNumsFunc(funcOnResult)
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
	tags := rtp.Root.GetTags(nil)
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

	if len(tags) == 0 {
		tags = append(tags, "")
	} else {
		sort.Slice(tags, func(i, j int) bool {
			return tags[i] < tags[j]
		})
	}

	strhead := "gamemod,tag,symbol,totalbet"
	for _, v := range sn {
		strhead = sgc7utils.AppendString(strhead, ",X", strconv.Itoa(v))
	}
	strhead = sgc7utils.AppendString(strhead, ",totalwin\n")

	f.WriteString(strhead)

	for _, gm := range gms {
		for _, tag := range tags {
			for _, symbol := range symbols {
				str := rtp.Root.GenSymbolString(gm, tag, symbol, sn, rtp.TotalBet)
				f.WriteString(str)
			}

			str := rtp.Root.GenTagString(gm, tag, sn, rtp.TotalBet)
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

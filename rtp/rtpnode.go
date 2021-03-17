package sgc7rtp

import (
	"strconv"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

const (
	// RTPNodeRoot - root
	RTPNodeRoot = 0
	// RTPNodeGameMod - gamemod
	RTPNodeGameMod = 1
	// RTPNodeTag - tag
	RTPNodeTag = 2
	// RTPNodeSymbol - symbol
	RTPNodeSymbol = 3
	// RTPNodeSymbolNums - symbol nums
	RTPNodeSymbolNums = 4
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
			if v.Symbol == node.Symbol && v.SymbolNums == node.SymbolNums {
				node.TriggerNums++
				node.TotalWin += int64(v.CashWin)
			}
		}
	}

	return false
}

// RTPNode -
type RTPNode struct {
	NodeType     int
	TriggerNums  int64
	TotalWin     int64
	RTP          float64
	MapChildren  map[string]*RTPNode
	GameMod      string
	Symbol       int
	SymbolNums   int
	TagName      string
	FuncOnResult FuncOnResult
}

// NewRTPRoot - new RTPNode
func NewRTPRoot() *RTPNode {
	return &RTPNode{
		NodeType:     RTPNodeRoot,
		MapChildren:  make(map[string]*RTPNode),
		FuncOnResult: OnRootResult,
	}
}

// NewRTPGameMod - new RTPNode
func NewRTPGameMod(gamemod string) *RTPNode {
	return &RTPNode{
		NodeType:     RTPNodeGameMod,
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		FuncOnResult: OnGameModResult,
		Symbol:       -1,
		SymbolNums:   -1,
	}
}

// NewRTPGameModTag - new RTPNode
func NewRTPGameModTag(gamemod string, tag string, funcTag FuncOnResult) *RTPNode {
	return &RTPNode{
		NodeType:     RTPNodeTag,
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		FuncOnResult: funcTag,
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

// InitGameMod2 - new RTPNode
func InitGameMod2(node *RTPNode, tags []string, funcTag []FuncOnResult, symbols []int, nums []int, onSymbolResult FuncOnResult, onSymbolNumsResult FuncOnResult) {
	if len(tags) > 0 {
		for i, tag := range tags {
			ctn := NewRTPGameModTag(node.GameMod, tag, funcTag[i])
			InitGameModTag2(ctn, tag, symbols, nums, onSymbolResult, onSymbolNumsResult)
			node.AddChild(tag, ctn)
		}
	} else {
		for _, sv := range symbols {
			csn := NewRTPSymbol2(node.GameMod, "", sv, onSymbolResult)
			InitSymbol2(csn, "", sv, nums, onSymbolNumsResult)
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

// InitGameModTag2 - new RTPNode
func InitGameModTag2(node *RTPNode, tag string, symbols []int, nums []int, onSymbolResult FuncOnResult, onSymbolNumsResult FuncOnResult) {
	for _, sv := range symbols {
		csn := NewRTPSymbol2(node.GameMod, tag, sv, onSymbolResult)
		InitSymbol2(csn, tag, sv, nums, onSymbolNumsResult)
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

// InitSymbol2 - new RTPNode
func InitSymbol2(node *RTPNode, tag string, symbol int, nums []int, onSymbolNumsResult FuncOnResult) {
	for _, nv := range nums {
		csnn := NewRTPSymbolNums2(node.GameMod, tag, symbol, nv, onSymbolNumsResult)
		node.AddChild(strconv.Itoa(nv), csnn)
	}
}

// NewRTPSymbol - new RTPNode
func NewRTPSymbol(gamemod string, tag string, symbol int) *RTPNode {
	return &RTPNode{
		NodeType:     RTPNodeSymbol,
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		Symbol:       symbol,
		FuncOnResult: OnSymbolResult,
		TagName:      tag,
		SymbolNums:   -1,
	}
}

// NewRTPSymbol2 - new RTPNode
func NewRTPSymbol2(gamemod string, tag string, symbol int, onSymbolResult FuncOnResult) *RTPNode {
	return &RTPNode{
		NodeType:     RTPNodeSymbol,
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		Symbol:       symbol,
		FuncOnResult: onSymbolResult,
		TagName:      tag,
		SymbolNums:   -1,
	}
}

// NewRTPSymbolNums - new RTPNode
func NewRTPSymbolNums(gamemod string, tag string, symbol int, nums int) *RTPNode {
	return &RTPNode{
		NodeType:     RTPNodeSymbolNums,
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		Symbol:       symbol,
		SymbolNums:   nums,
		FuncOnResult: OnSymbolNumsResult,
		TagName:      tag,
	}
}

// NewRTPSymbolNums2 - new RTPNode
func NewRTPSymbolNums2(gamemod string, tag string, symbol int, nums int, onSymbolNumsResult FuncOnResult) *RTPNode {
	return &RTPNode{
		NodeType:     RTPNodeSymbolNums,
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      gamemod,
		Symbol:       symbol,
		SymbolNums:   nums,
		FuncOnResult: onSymbolNumsResult,
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
	ismine := node.FuncOnResult(node, pr)
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
func (node *RTPNode) GetTags(arr []string, gamemod string) []string {
	if node.GameMod == gamemod && node.TagName != "" {
		if sgc7utils.IndexOfStringSlice(arr, node.TagName, 0) < 0 {
			arr = append(arr, node.TagName)
		}
	}

	for _, v := range node.MapChildren {
		arr = v.GetTags(arr, gamemod)
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
	if node.NodeType == RTPNodeRoot {
		str := sgc7utils.AppendString(",,,", strconv.FormatInt(totalbet, 10))
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
	if node.NodeType == RTPNodeRoot {
		for _, v := range node.MapChildren {
			str := v.GenGameModString(gamemod, sn, totalbet)
			if str != "" {
				return str
			}
		}

		return ""
	}

	if node.NodeType == RTPNodeGameMod && node.GameMod == gamemod {
		str := sgc7utils.AppendString(node.GameMod, ",,,", strconv.FormatInt(totalbet, 10))
		for range sn {
			str = sgc7utils.AppendString(str, ",")
		}
		str = sgc7utils.AppendString(str, ",", strconv.FormatInt(node.TotalWin, 10), "\n")

		return str
	}

	return ""
}

// GenTagString -
func (node *RTPNode) GenTagString(gamemod string, tag string, sn []int, totalbet int64) string {
	if node.NodeType < RTPNodeTag {
		for _, v := range node.MapChildren {
			str := v.GenTagString(gamemod, tag, sn, totalbet)
			if str != "" {
				return str
			}
		}

		return ""
	}

	if node.NodeType == RTPNodeTag && node.GameMod == gamemod && node.TagName == tag {
		str := sgc7utils.AppendString(node.GameMod, ",", tag, ",,", strconv.FormatInt(totalbet, 10))
		for range sn {
			str = sgc7utils.AppendString(str, ",")
		}
		str = sgc7utils.AppendString(str, ",", strconv.FormatInt(node.TotalWin, 10), "\n")

		return str
	}

	return ""
}

// GenSymbolString -
func (node *RTPNode) GenSymbolString(gamemod string, tag string, symbol int, sn []int, totalbet int64) string {
	if node.NodeType == RTPNodeRoot ||
		(node.NodeType == RTPNodeGameMod && node.GameMod == gamemod) ||
		(node.NodeType == RTPNodeTag && node.GameMod == gamemod && node.TagName == tag) {

		for _, v := range node.MapChildren {
			str := v.GenSymbolString(gamemod, tag, symbol, sn, totalbet)
			if str != "" {
				return str
			}
		}

		return ""
	}

	if node.NodeType == RTPNodeSymbol && node.GameMod == gamemod && node.TagName == tag && node.Symbol == symbol {
		str := sgc7utils.AppendString(node.GameMod, ",", tag, ",", strconv.Itoa(symbol), ",", strconv.FormatInt(totalbet, 10))
		for _, v := range sn {
			won := node.GetSymbolNumsWon(gamemod, tag, symbol, v)
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
func (node *RTPNode) GetSymbolNumsWon(gamemod string, tag string, symbol int, sn int) int64 {
	if node.NodeType == RTPNodeRoot ||
		(node.NodeType == RTPNodeGameMod && node.GameMod == gamemod) ||
		(node.NodeType == RTPNodeTag && node.GameMod == gamemod && node.TagName == tag) ||
		(node.NodeType == RTPNodeSymbol && node.GameMod == gamemod && node.TagName == tag && node.Symbol == symbol) {

		for _, v := range node.MapChildren {
			won := v.GetSymbolNumsWon(gamemod, tag, symbol, sn)
			if won >= 0 {
				return won
			}
		}

		return -1
	}

	if node.NodeType == RTPNodeSymbolNums && node.GameMod == gamemod && node.TagName == tag && node.Symbol == symbol && node.SymbolNums == sn {
		return node.TotalWin
	}

	return -1
}

// ChgSymbolNumsFunc -
func (node *RTPNode) ChgSymbolNumsFunc(funcOnResult FuncOnResult) {
	if node.Symbol >= 0 && node.SymbolNums > 0 {
		node.FuncOnResult = funcOnResult
	}

	for _, v := range node.MapChildren {
		v.ChgSymbolNumsFunc(funcOnResult)
	}
}

// Clone - clone
func (node *RTPNode) Clone() *RTPNode {
	node1 := &RTPNode{
		NodeType:     node.NodeType,
		TriggerNums:  node.TriggerNums,
		TotalWin:     node.TotalWin,
		RTP:          node.RTP,
		MapChildren:  make(map[string]*RTPNode),
		GameMod:      node.GameMod,
		Symbol:       node.Symbol,
		SymbolNums:   node.SymbolNums,
		TagName:      node.TagName,
		FuncOnResult: node.FuncOnResult,
	}

	for k, v := range node.MapChildren {
		node1.MapChildren[k] = v.Clone()
	}

	return node1
}

// Add - add
func (node *RTPNode) Add(node1 *RTPNode) {
	if node.NodeType == node1.NodeType &&
		node.GameMod == node1.GameMod &&
		node.Symbol == node1.Symbol &&
		node.SymbolNums == node1.SymbolNums &&
		node.TagName == node1.TagName {

		node.TriggerNums += node1.TriggerNums
		node.TotalWin += node1.TotalWin

		for k, v := range node.MapChildren {
			v.Add(node1.MapChildren[k])
		}
	}
}

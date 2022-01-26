package sgc7rtp

import (
	"strconv"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// FuncHROnResult - onResult(*RTP, *HitRateNode, *sgc7game.PlayResult)
type FuncHROnResult func(rtp *RTP, node *HitRateNode, pr *sgc7game.PlayResult) bool

// HitRateNode -
type HitRateNode struct {
	TagName      string
	BetNums      int64
	TriggerNums  int64
	TotalNums    int64
	FuncOnResult FuncHROnResult
}

// NewSpecialHitRate - new HitRateNode
func NewSpecialHitRate(tag string, funcOnResult FuncHROnResult) *HitRateNode {
	return &HitRateNode{
		TagName:      tag,
		FuncOnResult: funcOnResult,
	}
}

// GenString -
func (node *HitRateNode) GenString() string {
	str := goutils.AppendString(node.TagName, ",",
		strconv.FormatInt(node.BetNums, 10), ",",
		strconv.FormatInt(node.TriggerNums, 10), ",",
		strconv.FormatInt(node.TotalNums, 10), ",",
		strconv.FormatFloat(float64(node.TriggerNums)/float64(node.BetNums), 'f', -1, 64), ",",
		strconv.FormatFloat(float64(node.TotalNums)/float64(node.TriggerNums), 'f', -1, 64), "\n")

	return str
}

// Clone - clone
func (node *HitRateNode) Clone() *HitRateNode {
	node1 := &HitRateNode{
		TagName:      node.TagName,
		BetNums:      node.BetNums,
		TriggerNums:  node.TriggerNums,
		TotalNums:    node.TotalNums,
		FuncOnResult: node.FuncOnResult,
	}

	return node1
}

// Add - add
func (node *HitRateNode) Add(node1 *HitRateNode) {
	if node.TagName == node1.TagName {
		node.TriggerNums += node1.TriggerNums
		node.TotalNums += node1.TotalNums
		node.BetNums += node1.BetNums
	}
}

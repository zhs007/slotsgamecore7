package sgc7rtp

import (
	"strconv"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// FuncHROnResult - onResult(*HitRateNode, *sgc7game.PlayResult)
type FuncHROnResult func(node *HitRateNode, pr *sgc7game.PlayResult) bool

// HitRateNode -
type HitRateNode struct {
	TagName      string
	BetNums      int64
	TriggerNums  int64
	TotalNums    int64
	funcOnResult FuncHROnResult
}

// NewSpecialHitRate - new HitRateNode
func NewSpecialHitRate(tag string, funcOnResult FuncHROnResult) *HitRateNode {
	return &HitRateNode{
		TagName:      tag,
		funcOnResult: funcOnResult,
	}
}

// GenString -
func (node *HitRateNode) GenString() string {
	str := sgc7utils.AppendString(node.TagName, ",",
		strconv.FormatFloat(float64(node.TriggerNums)/float64(node.BetNums), 'f', -1, 64), ",",
		strconv.FormatFloat(float64(node.TotalNums)/float64(node.TriggerNums), 'f', -1, 64), "\n")

	return str
}

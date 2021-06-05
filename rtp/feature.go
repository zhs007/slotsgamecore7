package sgc7rtp

import (
	"strconv"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// FuncFeatureOnResults - onResult(*FeatureNode, []*sgc7game.PlayResult)
type FuncFeatureOnResults func(node *FeatureNode, lst []*sgc7game.PlayResult) bool

// FeatureNode -
type FeatureNode struct {
	TagName       string
	BetNums       int64
	TriggerNums   int64
	TotalNums     int64
	FuncOnResults FuncFeatureOnResults
}

// NewFeatureNode - new FeatureNode
func NewFeatureNode(tag string, funcOnResults FuncFeatureOnResults) *FeatureNode {
	return &FeatureNode{
		TagName:       tag,
		FuncOnResults: funcOnResults,
	}
}

// GenString -
func (node *FeatureNode) GenString() string {
	str := sgc7utils.AppendString(node.TagName, ",",
		strconv.FormatInt(node.BetNums, 10), ",",
		strconv.FormatInt(node.TriggerNums, 10), ",",
		strconv.FormatInt(node.TotalNums, 10), ",",
		strconv.FormatFloat(float64(node.TriggerNums)/float64(node.BetNums), 'f', -1, 64), ",",
		strconv.FormatFloat(float64(node.TotalNums)/float64(node.TriggerNums), 'f', -1, 64), "\n")

	return str
}

// Clone - clone
func (node *FeatureNode) Clone() *FeatureNode {
	node1 := &FeatureNode{
		TagName:       node.TagName,
		BetNums:       node.BetNums,
		TriggerNums:   node.TriggerNums,
		TotalNums:     node.TotalNums,
		FuncOnResults: node.FuncOnResults,
	}

	return node1
}

// Add - add
func (node *FeatureNode) Add(node1 *FeatureNode) {
	if node.TagName == node1.TagName {
		node.TriggerNums += node1.TriggerNums
		node.TotalNums += node1.TotalNums
		node.BetNums += node1.BetNums
	}
}

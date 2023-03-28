package sgc7rtp

import (
	"strconv"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// FuncFeatureOnResults - onResult(*FeatureNode, []*sgc7game.PlayResult, interface{})
type FuncFeatureOnResults func(node *FeatureNode, lst []*sgc7game.PlayResult, gameData interface{}) bool

// FeatureNode -
type FeatureNode struct {
	TagName       string
	BetNums       int64
	TriggerNums   int64
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
	str := goutils.AppendString(node.TagName, ",",
		strconv.FormatInt(node.BetNums, 10), ",",
		strconv.FormatInt(node.TriggerNums, 10), "\n")

	return str
}

// Clone - clone
func (node *FeatureNode) Clone() *FeatureNode {
	node1 := &FeatureNode{
		TagName:       node.TagName,
		BetNums:       node.BetNums,
		TriggerNums:   node.TriggerNums,
		FuncOnResults: node.FuncOnResults,
	}

	return node1
}

// Add - add
func (node *FeatureNode) Add(node1 *FeatureNode) {
	if node.TagName == node1.TagName {
		node.TriggerNums += node1.TriggerNums
		node.BetNums += node1.BetNums
	}
}

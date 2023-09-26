package lowcode

import (
	"context"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type AwardsNode struct {
	Weight int      `yaml:"weight" json:"weight"`
	Awards []*Award `yaml:"awards" json:"awards"`
}

type AwardsWeights struct {
	Nodes     []*AwardsNode `yaml:"nodes" json:"nodes"`
	MaxWeight int           `yaml:"-" json:"-"`
}

func (aw *AwardsWeights) Init() {
	aw.MaxWeight = 0

	for _, v := range aw.Nodes {
		for _, award := range v.Awards {
			award.Init()
		}

		aw.MaxWeight += v.Weight
	}
}

func (aw *AwardsWeights) RandVal(plugin sgc7plugin.IPlugin) (*AwardsNode, error) {
	if len(aw.Nodes) == 1 {
		return aw.Nodes[0], nil
	}

	cr, err := plugin.Random(context.Background(), aw.MaxWeight)
	if err != nil {
		return nil, err
	}

	for _, v := range aw.Nodes {
		if cr < v.Weight {
			return v, nil
		}

		cr -= v.Weight
	}

	return nil, sgc7game.ErrInvalidWeights
}

package sgc7game

import (
	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
)

// ValWeights
type ValWeights struct {
	Vals      []int
	Weights   []int
	MaxWeight int
}

func NewValWeights(vals []int, weights []int) (*ValWeights, error) {
	if len(vals) != len(weights) {
		goutils.Error("NewValWeights",
			zap.Int("vals", len(vals)),
			zap.Int("weights", len(weights)),
			zap.Error(ErrInvalidValWeights))

		return nil, ErrInvalidValWeights
	}

	vw := &ValWeights{
		Vals:      make([]int, len(vals)),
		Weights:   make([]int, len(vals)),
		MaxWeight: 0,
	}

	copy(vw.Vals, vals)
	copy(vw.Weights, weights)

	for _, v := range weights {
		vw.MaxWeight += v
	}

	return vw, nil
}

func (vw *ValWeights) RandVal(plugin sgc7plugin.IPlugin) (int, error) {
	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("ValWeights.RandVal:RandWithWeights",
			zap.Error(err))

		return 0, err
	}

	return vw.Vals[ci], nil
}

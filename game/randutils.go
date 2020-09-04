package sgc7game

import (
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// RandWithWeights - random with the weights
func RandWithWeights(plugin sgc7plugin.IPlugin, max int, arr []int) (int, error) {
	if len(arr) > 0 && max > 0 {
		cr, err := plugin.Random(max)
		if err != nil {
			return -1, err
		}

		for i, v := range arr {
			if cr < v {
				return i, nil
			}

			cr -= v
		}
	}

	return -1, ErrInvalidWeights
}

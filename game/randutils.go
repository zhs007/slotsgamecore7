package sgc7game

import (
	"context"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// RandWithWeights - random with the weights
func RandWithWeights(plugin sgc7plugin.IPlugin, max int, arr []int) (int, error) {
	if len(arr) > 0 && max > 0 {
		cr, err := plugin.Random(context.Background(), max)
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

// RandList - random list
func RandList(plugin sgc7plugin.IPlugin, arr []int, num int) ([]int, error) {
	if len(arr) > 0 && num > 0 {
		if num >= len(arr) {
			return arr, nil
		}

		narr := []int{}

		for i := 0; i < num; i++ {
			cr, err := plugin.Random(context.Background(), len(arr))
			if err != nil {
				return nil, err
			}

			narr = append(narr, arr[cr])
			arr = append(arr[:cr], arr[cr+1:]...)
		}

		return narr, nil
	}

	return nil, ErrInvalidParam
}

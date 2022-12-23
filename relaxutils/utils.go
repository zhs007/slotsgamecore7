package relaxutils

import (
	"sort"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func BuildPayouts(pt *sgc7game.PayTables) *Payouts {
	po := &Payouts{}

	for k, arr := range pt.MapPay {
		s := pt.GetStringFromInt(k)

		p := &Payout{
			Name: s,
		}

		for i := range arr {
			if arr[i] > 0 {
				pw := &PayoutWin{
					Count:  i + 1,
					Payout: arr[i],
				}

				p.Wins = append(p.Wins, pw)
			}
		}

		if len(p.Wins) > 0 {
			po.Payouts = append(po.Payouts, p)
		}
	}

	sort.Slice(po.Payouts, func(i, j int) bool {
		return po.Payouts[i].Name < po.Payouts[j].Name
	})

	return po
}

func BuildSymbols(symbols []string) *Symbols {
	ret := &Symbols{}

	for i := range symbols {
		s := &Symbol{
			Name: symbols[i],
			Val:  i,
		}

		ret.Symbols = append(ret.Symbols, s)
	}

	return ret
}

func BuildReels(reels []*sgc7game.ReelsData, pt *sgc7game.PayTables) *Reels {
	ret := &Reels{}

	for _, arr2d := range reels {
		t := &Table{}

		for _, arr := range arr2d.Reels {
			lststr := &StringList{}

			for _, s := range arr {
				ss := pt.GetStringFromInt(s)
				lststr.Vals = append(lststr.Vals, ss)
			}

			t.Reel = append(t.Reel, lststr)
		}

		ret.Tables = append(ret.Tables, t)
	}

	return ret
}

func BuildWeights(weights []int) *Weights {
	ret := &Weights{}

	for i, v := range weights {
		if v > 0 {
			lst := &IntList{
				Vals: []int{i, v},
			}

			ret.Entries = append(ret.Entries, lst)
		}
	}

	return ret
}

func BuildWeightsEx(vals []int, weights []int) *Weights {
	if len(vals) != len(weights) {
		return nil
	}

	ret := &Weights{}

	for i, v := range weights {
		if v > 0 {
			lst := &IntList{
				Vals: []int{vals[i], v},
			}

			ret.Entries = append(ret.Entries, lst)
		}
	}

	return ret
}

func BuildWeightsArr(weights [][]int) *WeightsArr {
	ret := &WeightsArr{}

	for _, arr := range weights {
		w := &Weights{}

		for i, v := range arr {
			if v > 0 {
				lst := &IntList{
					Vals: []int{i, v},
				}

				w.Entries = append(w.Entries, lst)
			}
		}

		ret.Weights = append(ret.Weights, w)
	}

	return ret
}

func BuildWeightsArrEx(vals []int, weights [][]int) *WeightsArr {
	if len(vals) != len(weights[0]) {
		return nil
	}

	ret := &WeightsArr{}

	for _, arr := range weights {
		w := &Weights{}

		for i, v := range arr {
			if v > 0 {
				lst := &IntList{
					Vals: []int{vals[i], v},
				}

				w.Entries = append(w.Entries, lst)
			}
		}

		ret.Weights = append(ret.Weights, w)
	}

	return ret
}

func BuildInt3DArray(lst []*sgc7game.ReelsPosData) *Int3DArray {
	arr := &Int3DArray{}

	for _, rpd := range lst {
		arr2d := &Int2DArray{}

		for _, r := range rpd.ReelsPos {
			arr1d := &IntList{Vals: r}

			arr2d.Rows = append(arr2d.Rows, arr1d)
		}

		arr.Tables = append(arr.Tables, arr2d)
	}

	return arr
}

func BuildLineData(ld *sgc7game.LineData) *Int2DArray {
	arr2d := &Int2DArray{}

	for _, arr := range ld.Lines {
		arr1d := &IntList{Vals: make([]int, len(arr))}
		copy(arr1d.Vals, arr)

		arr2d.Rows = append(arr2d.Rows, arr1d)
	}

	return arr2d
}

func BuildInt2DArray(arr2 [][]int) *Int2DArray {
	arr2d := &Int2DArray{}

	for _, arr := range arr2 {
		arr1d := &IntList{Vals: make([]int, len(arr))}
		copy(arr1d.Vals, arr)

		arr2d.Rows = append(arr2d.Rows, arr1d)
	}

	return arr2d
}

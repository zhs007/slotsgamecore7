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
		lst := &IntList{
			Vals: []int{i, v},
		}

		ret.Entries = append(ret.Entries, lst)
	}

	return ret
}

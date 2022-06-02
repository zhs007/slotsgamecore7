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

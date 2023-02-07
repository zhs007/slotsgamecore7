package stats

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type IStats interface {
	OnResults(stake *sgc7game.Stake, lst []*sgc7game.PlayResult)
}

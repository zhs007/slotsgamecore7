package stats

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/mathtoolset"
)

type Reels struct {
	Reels []*Reel
}

func (reels *Reels) OnScene(scene *sgc7game.GameScene) {
	for _, v := range reels.Reels {
		v.OnScene(scene)
	}
}

func NewReels(width int, lst []mathtoolset.SymbolType) *Reels {
	reels := &Reels{}

	for x := 0; x < width; x++ {
		r := NewReel(x, lst)

		reels.Reels = append(reels.Reels, r)
	}

	return reels
}

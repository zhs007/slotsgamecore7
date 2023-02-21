package mathtoolset

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func getSymbolWithStops(reels *sgc7game.ReelsData, ri int, stops int) int {
	if stops < 0 {
		stops = -(-stops % len(reels.Reels[ri]))
	} else {
		stops = (stops % len(reels.Reels[ri]))
	}

	if stops < 0 {
		return reels.Reels[ri][len(reels.Reels[ri])+stops]
	}

	return reels.Reels[ri][stops]
}

func CheckReels(reels *sgc7game.ReelsData, minoff int) (int, int, error) {
	for x, reel := range reels.Reels {
		for y, s := range reel {
			for i := -minoff; i <= minoff; i++ {
				if i != 0 {
					if s == getSymbolWithStops(reels, x, y+i) {
						return x, y, ErrInvalidReelsWithMinOff
					}
				}
			}
		}
	}

	return 0, 0, nil
}

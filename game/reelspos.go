package sgc7game

import (
	"log/slog"

	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"golang.org/x/net/context"
)

// ReelsPosData - reels position data
type ReelsPosData struct {
	ReelsPos [][]int `json:"reelspos"`
}

func NewReelsPosData(reels *ReelsData) *ReelsPosData {
	reelspos := &ReelsPosData{}

	for range reels.Reels {
		reelspos.ReelsPos = append(reelspos.ReelsPos, []int{})
	}

	return reelspos
}

func (reelspos *ReelsPosData) AddPos(x, y int) {
	if x >= 0 && x < len(reelspos.ReelsPos) {
		reelspos.ReelsPos[x] = append(reelspos.ReelsPos[x], y)
	}
}

func (reelspos *ReelsPosData) RandReel(ctx context.Context, plugin sgc7plugin.IPlugin, x int) (int, error) {
	if x < 0 || x > len(reelspos.ReelsPos) {
		goutils.Error("ReelsPosData.RandReel",
			slog.Int("x", x),
			goutils.Err(ErrInvalidSceneX))

		return -1, ErrInvalidSceneX
	}

	cr, err := plugin.Random(ctx, len(reelspos.ReelsPos[x]))
	if err != nil {
		goutils.Error("ReelsPosData.RandReel",
			goutils.Err(err))

		return -1, err
	}

	return reelspos.ReelsPos[x][cr], nil
}

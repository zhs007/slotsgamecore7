package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func Test_CheckReels(t *testing.T) {
	rd, err := sgc7game.LoadReelsFromExcel("../unittestdata/reels6.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	x, y, err := CheckReels(rd, 2)
	assert.Error(t, err)

	assert.Equal(t, x, 1)
	assert.Equal(t, y, 1)

	t.Logf("Test_CheckReels OK")
}

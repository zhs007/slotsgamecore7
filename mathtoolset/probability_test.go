package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcScatterProbability(t *testing.T) {

	rss2, err := LoadReelsStats("../unittestdata/rss.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, rss2)

	prob := CalcScatterProbability(rss2, 8, 3, 3)

	assert.Equal(t, int64(prob*10000), int64(1072))

	t.Logf("Test_CalcScatterProbability OK")
}

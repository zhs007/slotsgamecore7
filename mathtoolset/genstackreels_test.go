package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenStackReels(t *testing.T) {
	rd, err := GenStackReels("../unittestdata/reelsstats2.xlsx", []int{1}, []string{"SC"})
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	SaveReels("../unittestdata/reelsstats2-output.xlsx", rd)

	t.Logf("Test_GenStackReels OK")
}

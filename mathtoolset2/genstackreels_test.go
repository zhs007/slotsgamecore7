package mathtoolset2

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenStackReels(t *testing.T) {
	file, err := os.Open("../unittestdata/reelsstats2.xlsx")
	assert.NoError(t, err)

	rd, err := GenStackReels(file, []int{2, 3}, []string{"SC"})
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	SaveReels("../unittestdata/reelsstats2-output.xlsx", rd)

	t.Logf("Test_GenStackReels OK")
}

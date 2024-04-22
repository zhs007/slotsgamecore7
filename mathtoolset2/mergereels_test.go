package mathtoolset2

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MergeReels(t *testing.T) {
	file, err := os.Open("../unittestdata/reels01.xlsx")
	assert.NoError(t, err)

	file2, err := os.Open("../unittestdata/reels01.xlsx")
	assert.NoError(t, err)

	rd, err := MergeReels([]io.Reader{file, file2})
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	SaveReels("../unittestdata/merge-output.xlsx", rd)

	t.Logf("Test_MergeReels OK")
}

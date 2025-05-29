package mathtoolset2

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenReels(t *testing.T) {
	file, err := os.Open("../unittestdata/genreels2.xlsx")
	assert.NoError(t, err)

	rd, err := GenReels(file, "SEP_5,WL,MY,MY2;SEP_1,L1,L2,L3,L4,L5,H1,H2,H3,H4,H5;EXC_5,WL,MY,MY2;")
	assert.NoError(t, err)

	SaveReels("../unittestdata/genreels2-output.xlsx", rd)

	t.Logf("Test_GenReels OK")
}

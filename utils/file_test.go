package sgc7utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsSameFile(t *testing.T) {
	issame := IsSameFile("../unittestdata/rtptestok.csv", "../unittestdata/rtptestok0.csv")
	assert.Equal(t, issame, true)

	issame = IsSameFile("../unittestdata/rtptestok.csv", "../unittestdata/rtptestok1.csv")
	assert.Equal(t, issame, false)

	issame = IsSameFile("../unittestdata/rtptestok.csv", "../unittestdata/rtptestok2.csv")
	assert.Equal(t, issame, false)

	issame = IsSameFile("../unittestdata/rtptestok.csv", "../unittestdata/rtptestok3.csv")
	assert.Equal(t, issame, false)

	issame = IsSameFile("../unittestdata/rtptestok3.csv", "../unittestdata/rtptestok.csv")
	assert.Equal(t, issame, false)

	t.Logf("Test_IsSameFile OK")
}

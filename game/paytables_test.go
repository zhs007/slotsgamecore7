package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadPayTables5JSON(t *testing.T) {
	_, err := LoadPayTables5JSON("../unittestdata/paytables1.json")
	if err == nil {
		t.Fatalf("Test_LoadPayTables5JSON LoadPayTables5JSON non-file error %v",
			err)
	}

	_, err = LoadPayTables5JSON("../unittestdata/errjson.json")
	if err == nil {
		t.Fatalf("Test_LoadPayTables5JSON LoadPayTables5JSON errjson error %v",
			err)
	}

	_, err = LoadPayTables5JSON("../unittestdata/empty.json")
	if err == nil {
		t.Fatalf("Test_LoadPayTables5JSON LoadPayTables5JSON empty error %v",
			err)
	}

	ld, err := LoadPayTables5JSON("../unittestdata/empty.linedata.json")
	if err != nil || ld != nil {
		t.Fatalf("Test_LoadPayTables5JSON LoadPayTables5JSON empty linedata error %v",
			err)
	}

	ld, err = LoadPayTables5JSON("../unittestdata/paytables.json")
	if err != nil {
		t.Fatalf("Test_LoadPayTables5JSON LoadPayTables5JSON error %v",
			err)
	}

	if ld == nil {
		t.Fatalf("Test_LoadPayTables5JSON LoadPayTables5JSON non-data")
	}

	for i := 0; i < 11; i++ {
		if len(ld.MapPay[i]) != 5 {
			t.Fatalf("Test_LoadPayTables5JSON pay %d length error %d != %d",
				i, len(ld.MapPay[i]), 5)
		}
	}

	pay0 := []int{0, 0, 50, 500, 2000}
	pay1 := []int{0, 0, 50, 200, 1000}
	pay5 := []int{0, 0, 10, 30, 120}
	pay11 := []int{0, 2, 5, 10, 100}

	for i := 0; i < 5; i++ {
		if ld.MapPay[0][i] != pay0[i] {
			t.Fatalf("Test_LoadPayTables5JSON pay0 data %d error %d != %d",
				i, ld.MapPay[0][i], pay0[i])
		}

		if ld.MapPay[1][i] != pay1[i] {
			t.Fatalf("Test_LoadPayTables5JSON pay1 data %d error %d != %d",
				i, ld.MapPay[1][i], pay1[i])
		}

		if ld.MapPay[5][i] != pay5[i] {
			t.Fatalf("Test_LoadPayTables5JSON pay5 data %d error %d != %d",
				i, ld.MapPay[5][i], pay5[i])
		}

		if ld.MapPay[11][i] != pay11[i] {
			t.Fatalf("Test_LoadPayTables5JSON pay11 data %d error %d != %d",
				i, ld.MapPay[11][i], pay11[i])
		}
	}

	t.Logf("Test_LoadPayTables5JSON OK")
}

func Test_LoadPaytablesFromExcel(t *testing.T) {
	pt, err := LoadPaytablesFromExcel("../unittestdata/paytables.xlsx")
	assert.NoError(t, err)

	assert.Equal(t, len(pt.MapPay), 13)
	assert.Equal(t, len(pt.MapSymbols), 13)

	assert.Equal(t, pt.MapSymbols["WL"], 0)
	assert.Equal(t, pt.MapSymbols["A"], 1)
	assert.Equal(t, pt.MapSymbols["MY"], 12)

	assert.Equal(t, len(pt.MapPay[0]), 5)
	assert.Equal(t, pt.MapPay[0][0], 0)
	assert.Equal(t, pt.MapPay[0][1], 0)
	assert.Equal(t, pt.MapPay[0][2], 0)
	assert.Equal(t, pt.MapPay[0][3], 0)
	assert.Equal(t, pt.MapPay[0][4], 0)

	assert.Equal(t, len(pt.MapPay[1]), 5)
	assert.Equal(t, pt.MapPay[1][0], 0)
	assert.Equal(t, pt.MapPay[1][1], 0)
	assert.Equal(t, pt.MapPay[1][2], 40)
	assert.Equal(t, pt.MapPay[1][3], 200)
	assert.Equal(t, pt.MapPay[1][4], 1000)

	assert.Equal(t, len(pt.MapPay[12]), 5)
	assert.Equal(t, pt.MapPay[12][0], 0)
	assert.Equal(t, pt.MapPay[12][1], 0)
	assert.Equal(t, pt.MapPay[12][2], 0)
	assert.Equal(t, pt.MapPay[12][3], 0)
	assert.Equal(t, pt.MapPay[12][4], 0)

	t.Logf("Test_LoadPaytablesFromExcel OK")
}

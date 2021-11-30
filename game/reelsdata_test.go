package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadReels5JSON(t *testing.T) {
	_, err := LoadReels5JSON("../unittestdata/reels1.json")
	if err == nil {
		t.Fatalf("Test_LoadReels5JSON LoadReels5JSON non-file error %v",
			err)
	}

	_, err = LoadReels5JSON("../unittestdata/errjson.json")
	if err == nil {
		t.Fatalf("Test_LoadReels5JSON LoadReels5JSON errjson error %v",
			err)
	}

	_, err = LoadReels5JSON("../unittestdata/empty.json")
	if err == nil {
		t.Fatalf("Test_LoadReels5JSON LoadReels5JSON empty error %v",
			err)
	}

	ld, err := LoadReels5JSON("../unittestdata/paytables.json")
	if err != nil || ld != nil {
		t.Fatalf("Test_LoadReels5JSON LoadReels5JSON format error %v",
			err)
	}

	ld, err = LoadReels5JSON("../unittestdata/empty.linedata.json")
	if err != nil || ld != nil {
		t.Fatalf("Test_LoadReels5JSON LoadReels5JSON empty linedata error %v",
			err)
	}

	ld, err = LoadReels5JSON("../unittestdata/reels.json")
	if err != nil {
		t.Fatalf("Test_LoadReels5JSON LoadReels5JSON error %v",
			err)
	}

	if ld == nil {
		t.Fatalf("Test_LoadReels5JSON LoadReels5JSON non-data")
	}

	// for i := 0; i < 11; i++ {
	// 	if len(ld.MapPay[i]) != 5 {
	// 		t.Fatalf("Test_LoadReels5JSON pay %d length error %d != %d",
	// 			i, len(ld.MapPay[i]), 5)
	// 	}
	// }

	// pay0 := []int{0, 0, 50, 500, 2000}
	// pay1 := []int{0, 0, 50, 200, 1000}
	// pay5 := []int{0, 0, 10, 30, 120}
	// pay11 := []int{0, 2, 5, 10, 100}

	// for i := 0; i < 5; i++ {
	// 	if ld.MapPay[0][i] != pay0[i] {
	// 		t.Fatalf("Test_LoadReels5JSON pay0 data %d error %d != %d",
	// 			i, ld.MapPay[0][i], pay0[i])
	// 	}

	// 	if ld.MapPay[1][i] != pay1[i] {
	// 		t.Fatalf("Test_LoadReels5JSON pay1 data %d error %d != %d",
	// 			i, ld.MapPay[1][i], pay1[i])
	// 	}

	// 	if ld.MapPay[5][i] != pay5[i] {
	// 		t.Fatalf("Test_LoadReels5JSON pay5 data %d error %d != %d",
	// 			i, ld.MapPay[5][i], pay5[i])
	// 	}

	// 	if ld.MapPay[11][i] != pay11[i] {
	// 		t.Fatalf("Test_LoadReels5JSON pay11 data %d error %d != %d",
	// 			i, ld.MapPay[11][i], pay11[i])
	// 	}
	// }

	t.Logf("Test_LoadReels5JSON OK")
}

func Test_LoadReelsFromExcel(t *testing.T) {
	rd, err := LoadReelsFromExcel("../unittestdata/reels.xlsx")
	assert.NoError(t, err)

	assert.Equal(t, len(rd.Reels), 5)
	assert.Equal(t, len(rd.Reels[0]), 36)
	assert.Equal(t, len(rd.Reels[1]), 196)
	assert.Equal(t, len(rd.Reels[2]), 200)
	assert.Equal(t, len(rd.Reels[3]), 260)
	assert.Equal(t, len(rd.Reels[4]), 180)

	assert.Equal(t, rd.Reels[0][0], 1)
	assert.Equal(t, rd.Reels[1][0], 2)
	assert.Equal(t, rd.Reels[2][0], 4)
	assert.Equal(t, rd.Reels[3][0], 5)
	assert.Equal(t, rd.Reels[4][0], 4)

	assert.Equal(t, rd.Reels[0][35], 12)
	assert.Equal(t, rd.Reels[1][195], 12)
	assert.Equal(t, rd.Reels[2][199], 7)
	assert.Equal(t, rd.Reels[3][259], 12)
	assert.Equal(t, rd.Reels[4][179], 12)

	t.Logf("Test_LoadReelsFromExcel OK")
}

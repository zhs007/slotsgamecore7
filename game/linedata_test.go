package sgc7game

import (
	"testing"
)

func Test_LoadLine5JSON(t *testing.T) {
	_, err := LoadLine5JSON("../unittestdata/linedata1.json")
	if err == nil {
		t.Fatalf("Test_LoadLine5JSON LoadLine5JSON non-file error %v",
			err)
	}

	_, err = LoadLine5JSON("../unittestdata/errjson.json")
	if err == nil {
		t.Fatalf("Test_LoadLine5JSON LoadLine5JSON errjson error %v",
			err)
	}

	_, err = LoadLine5JSON("../unittestdata/empty.json")
	if err == nil {
		t.Fatalf("Test_LoadLine5JSON LoadLine5JSON empty error %v",
			err)
	}

	ld, err := LoadLine5JSON("../unittestdata/paytables.json")
	if err != nil || ld != nil {
		t.Fatalf("Test_LoadLine5JSON LoadLine5JSON format error %v",
			err)
	}

	ld, err = LoadLine5JSON("../unittestdata/empty.linedata.json")
	if err != nil || ld != nil {
		t.Fatalf("Test_LoadLine5JSON LoadLine5JSON empty linedata error %v",
			err)
	}

	ld, err = LoadLine5JSON("../unittestdata/linedata.json")
	if err != nil {
		t.Fatalf("Test_LoadLine5JSON LoadLine5JSON error %v",
			err)
	}

	if ld == nil {
		t.Fatalf("Test_LoadLine5JSON LoadLine5JSON non-data")
	}

	if len(ld.Lines) != 25 {
		t.Fatalf("Test_LoadLine5JSON lines error %d != %d",
			len(ld.Lines), 25)
	}

	for i, v := range ld.Lines {
		if len(v) != 5 {
			t.Fatalf("Test_LoadLine5JSON line %d error %d != %d",
				i, len(v), 5)
		}
	}

	line1 := []int{1, 1, 1, 1, 1}
	line10 := []int{0, 1, 1, 1, 0}
	line25 := []int{2, 0, 0, 0, 2}

	for i := 0; i < 5; i++ {
		if ld.Lines[0][i] != line1[i] {
			t.Fatalf("Test_LoadLine5JSON line1 data %d error %d != %d",
				i, ld.Lines[0][i], line1[i])
		}

		if ld.Lines[9][i] != line10[i] {
			t.Fatalf("Test_LoadLine5JSON line10 data %d error %d != %d",
				i, ld.Lines[0][i], line10[i])
		}

		if ld.Lines[24][i] != line25[i] {
			t.Fatalf("Test_LoadLine5JSON line25 data %d error %d != %d",
				i, ld.Lines[0][i], line25[i])
		}
	}

	t.Logf("Test_LoadLine5JSON OK")
}

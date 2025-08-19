package sgc7game

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadExcel(t *testing.T) {
	header := []string{}
	data := [][4]string{}

	err := LoadExcel("../unittestdata/jules_test.xlsx", "Sheet1",
		func(x int, str string) string {
			header = append(header, str)
			return str
		},
		func(x int, y int, h string, d string) error {
			if len(data) < y {
				data = append(data, [4]string{})
			}
			data[y-1][x] = d
			return nil
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{"type", "val"}, header)
	assert.Equal(t, "1", data[0][0])
	assert.Equal(t, "1.1", data[0][1])
	assert.Equal(t, "2", data[1][0])
	assert.Equal(t, "2.2", data[1][1])

	t.Logf("Test_LoadExcel OK")
}

func Test_LoadExcelWithReader(t *testing.T) {
	header := []string{}
	data := [][4]string{}

	file, err := os.Open("../unittestdata/jules_test.xlsx")
	assert.NoError(t, err)
	defer file.Close()

	err = LoadExcelWithReader(file, "Sheet1",
		func(x int, str string) string {
			header = append(header, str)
			return str
		},
		func(x int, y int, h string, d string) error {
			if len(data) < y {
				data = append(data, [4]string{})
			}
			data[y-1][x] = d
			return nil
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, []string{"type", "val"}, header)
	assert.Equal(t, "1", data[0][0])
	assert.Equal(t, "1.1", data[0][1])
	assert.Equal(t, "2", data[1][0])
	assert.Equal(t, "2.2", data[1][1])

	t.Logf("Test_LoadExcelWithReader OK")
}

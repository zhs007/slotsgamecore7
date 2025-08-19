package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func Test_LoadLineDataFromExcel_Jules(t *testing.T) {
	f := excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]interface{}{"R1", "R2", "R3", "R4", "R5"})
	f.SetSheetRow("Sheet1", "A2", &[]interface{}{1, 2, 3, 4, 5})
	f.SetSheetRow("Sheet1", "A3", &[]interface{}{6, 7, 8, 9, 10})
	err := f.SaveAs("testdata/linedata_valid.xlsx")
	assert.NoError(t, err)

	ld, err := LoadLineDataFromExcel("testdata/linedata_valid.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, ld)

	assert.Equal(t, 2, len(ld.Lines))
	assert.Equal(t, 5, len(ld.Lines[0]))
	assert.Equal(t, 1, ld.Lines[0][0])
	assert.Equal(t, 7, ld.Lines[1][1])

	ld, err = LoadLineDataFromExcel("testdata/nofile.xlsx")
	assert.Error(t, err)
	assert.Nil(t, ld)

	t.Logf("Test_LoadLineDataFromExcel_Jules OK")
}

func Test_LoadLine5JSON_Jules(t *testing.T) {
	ld, err := LoadLine5JSON("testdata/line5_valid.json")
	assert.NoError(t, err)
	assert.NotNil(t, ld)

	assert.Equal(t, 2, len(ld.Lines))
	assert.Equal(t, 5, len(ld.Lines[0]))
	assert.Equal(t, 0, ld.Lines[0][0])
	assert.Equal(t, 1, ld.Lines[1][1])

	ld, err = LoadLine5JSON("testdata/line5_invalid.json")
	assert.Error(t, err)
	assert.Nil(t, ld)

	ld, err = LoadLine5JSON("testdata/nofile.json")
	assert.Error(t, err)
	assert.Nil(t, ld)

	t.Logf("Test_LoadLine5JSON_Jules OK")
}

func Test_LoadLine3JSON_Jules(t *testing.T) {
	ld, err := LoadLine3JSON("testdata/line3_valid.json")
	assert.NoError(t, err)
	assert.NotNil(t, ld)

	assert.Equal(t, 2, len(ld.Lines))
	assert.Equal(t, 3, len(ld.Lines[0]))
	assert.Equal(t, 0, ld.Lines[0][0])
	assert.Equal(t, 1, ld.Lines[1][1])

	ld, err = LoadLine3JSON("testdata/nofile.json")
	assert.Error(t, err)
	assert.Nil(t, ld)

	t.Logf("Test_LoadLine3JSON_Jules OK")
}

func Test_LoadLine6JSON_Jules(t *testing.T) {
	ld, err := LoadLine6JSON("testdata/line6_valid.json")
	assert.NoError(t, err)
	assert.NotNil(t, ld)

	assert.Equal(t, 2, len(ld.Lines))
	assert.Equal(t, 6, len(ld.Lines[0]))
	assert.Equal(t, 0, ld.Lines[0][0])
	assert.Equal(t, 1, ld.Lines[1][1])

	ld, err = LoadLine6JSON("testdata/nofile.json")
	assert.Error(t, err)
	assert.Nil(t, ld)

	t.Logf("Test_LoadLine6JSON_Jules OK")
}

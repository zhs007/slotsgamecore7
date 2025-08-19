package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func Test_LoadPaytablesFromExcel_Jules(t *testing.T) {
	f := excelize.NewFile()
	f.SetSheetRow("Sheet1", "A1", &[]interface{}{"Code", "Symbol", "X1", "X2", "X3", "X4", "X5"})
	f.SetSheetRow("Sheet1", "A2", &[]interface{}{0, "S0", 0, 0, 10, 20, 30})
	f.SetSheetRow("Sheet1", "A3", &[]interface{}{1, "S1", 0, 5, 15, 25, 35})
	err := f.SaveAs("testdata/paytables_valid.xlsx")
	assert.NoError(t, err)

	pt, err := LoadPaytablesFromExcel("testdata/paytables_valid.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, pt)

	assert.Equal(t, 2, len(pt.MapPay))
	assert.Equal(t, 5, len(pt.MapPay[0]))
	assert.Equal(t, 10, pt.MapPay[0][2])
	assert.Equal(t, 35, pt.MapPay[1][4])

	assert.Equal(t, 2, len(pt.MapSymbols))
	assert.Equal(t, 0, pt.MapSymbols["S0"])
	assert.Equal(t, 1, pt.MapSymbols["S1"])

	pt, err = LoadPaytablesFromExcel("testdata/nofile.xlsx")
	assert.Error(t, err)
	assert.Nil(t, pt)

	t.Logf("Test_LoadPaytablesFromExcel_Jules OK")
}

func Test_LoadPayTables5JSON_Jules(t *testing.T) {
	pt, err := LoadPayTables5JSON("testdata/paytables5_valid.json")
	assert.NoError(t, err)
	assert.NotNil(t, pt)

	assert.Equal(t, 2, len(pt.MapPay))
	assert.Equal(t, 5, len(pt.MapPay[0]))
	assert.Equal(t, 10, pt.MapPay[0][2])
	assert.Equal(t, 35, pt.MapPay[1][4])

	assert.Equal(t, 2, len(pt.MapSymbols))
	assert.Equal(t, 0, pt.MapSymbols["S0"])
	assert.Equal(t, 1, pt.MapSymbols["S1"])

	str := pt.GetStringFromInt(0)
	assert.Equal(t, "S0", str)

	str = pt.GetStringFromInt(1)
	assert.Equal(t, "S1", str)

	str = pt.GetStringFromInt(100)
	assert.Equal(t, "", str)

	minwin := pt.GetSymbolMinWinNum(0)
	assert.Equal(t, 3, minwin)

	minwin = pt.GetSymbolMinWinNum(1)
	assert.Equal(t, 2, minwin)

	minwin = pt.GetSymbolMinWinNum(100)
	assert.Equal(t, 0, minwin)

	pt, err = LoadPayTables5JSON("testdata/nofile.json")
	assert.Error(t, err)
	assert.Nil(t, pt)

	t.Logf("Test_LoadPayTables5JSON_Jules OK")
}

func Test_LoadPayTables3JSON_Jules(t *testing.T) {
	pt, err := LoadPayTables3JSON("testdata/paytables3_valid.json")
	assert.NoError(t, err)
	assert.NotNil(t, pt)

	assert.Equal(t, 2, len(pt.MapPay))
	assert.Equal(t, 3, len(pt.MapPay[0]))
	assert.Equal(t, 20, pt.MapPay[0][2])
	assert.Equal(t, 25, pt.MapPay[1][2])

	t.Logf("Test_LoadPayTables3JSON_Jules OK")
}

func Test_LoadPayTables6JSON_Jules(t *testing.T) {
	pt, err := LoadPayTables6JSON("testdata/paytables6_valid.json")
	assert.NoError(t, err)
	assert.NotNil(t, pt)

	assert.Equal(t, 2, len(pt.MapPay))
	assert.Equal(t, 6, len(pt.MapPay[0]))
	assert.Equal(t, 40, pt.MapPay[0][5])
	assert.Equal(t, 45, pt.MapPay[1][5])

	t.Logf("Test_LoadPayTables6JSON_Jules OK")
}

func Test_LoadPayTables15JSON_Jules(t *testing.T) {
	pt, err := LoadPayTables15JSON("testdata/paytables15_valid.json")
	assert.NoError(t, err)
	assert.NotNil(t, pt)

	assert.Equal(t, 2, len(pt.MapPay))
	assert.Equal(t, 15, len(pt.MapPay[0]))
	assert.Equal(t, 130, pt.MapPay[0][14])
	assert.Equal(t, 135, pt.MapPay[1][14])

	t.Logf("Test_LoadPayTables15JSON_Jules OK")
}

func Test_LoadPayTables25JSON_Jules(t *testing.T) {
	pt, err := LoadPayTables25JSON("testdata/paytables25_valid.json")
	assert.NoError(t, err)
	assert.NotNil(t, pt)

	assert.Equal(t, 2, len(pt.MapPay))
	assert.Equal(t, 25, len(pt.MapPay[0]))
	assert.Equal(t, 230, pt.MapPay[0][24])
	assert.Equal(t, 235, pt.MapPay[1][24])

	t.Logf("Test_LoadPayTables25JSON_Jules OK")
}

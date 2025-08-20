package main

import (
	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	// Create a new sheet.
	index, _ := f.NewSheet("Sheet1")
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "type")
	f.SetCellValue("Sheet1", "B1", "val")
	f.SetCellValue("Sheet1", "A2", 1)
	f.SetCellValue("Sheet1", "B2", "100,101,102")
	f.SetCellValue("Sheet1", "A3", 2)
	f.SetCellValue("Sheet1", "B3", "200,201")
	f.SetCellValue("Sheet1", "A4", 3)
	f.SetCellValue("Sheet1", "B4", "300")
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("./game/testdata/valarrmapping.xlsx"); err != nil {
		println(err.Error())
	}
}

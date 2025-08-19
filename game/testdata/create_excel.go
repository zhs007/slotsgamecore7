package main

import (
	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	// Create a new sheet.
	index, _ := f.NewSheet("Sheet1")
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "val1")
	f.SetCellValue("Sheet1", "B1", "val2")
	f.SetCellValue("Sheet1", "C1", "weight")
	f.SetCellValue("Sheet1", "A2", 1)
	f.SetCellValue("Sheet1", "B2", 2)
	f.SetCellValue("Sheet1", "C2", 10)
	f.SetCellValue("Sheet1", "A3", 3)
	f.SetCellValue("Sheet1", "B3", 4)
	f.SetCellValue("Sheet1", "C3", 20)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("arrvalweights.xlsx"); err != nil {
		println(err.Error())
	}
}

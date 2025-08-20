package main

import (
	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	// Create a new sheet.
	index, _ := f.NewSheet("Sheet1")
	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "val")
	f.SetCellValue("Sheet1", "B1", "weight")
	f.SetCellValue("Sheet1", "A2", "a")
	f.SetCellValue("Sheet1", "B2", 10)
	f.SetCellValue("Sheet1", "A3", "b")
	f.SetCellValue("Sheet1", "B3", 20)
	f.SetCellValue("Sheet1", "A4", "c")
	f.SetCellValue("Sheet1", "B4", 30)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("strvalweights.xlsx"); err != nil {
		println(err.Error())
	}
}

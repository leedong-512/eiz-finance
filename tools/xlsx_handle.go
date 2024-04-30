package tools

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

func Unique(file string, u string, w string)  {
	parts := strings.Split(u, ",")
	//var sheetData map[string]string
	var sheetDatas map[string][]string
	sheetDatas = make(map[string][]string)
	for _, record := range records {
		var key string
		for _, part := range parts {
			key += record[HeaderMap[part]]
		}
		_, found := sheetDatas[key]
		if !found {
			sheetDatas[key] = record
		}
	}

	nFile, _ := OpenXlsx(file)
	for _, sheetData := range sheetDatas {
		fmt.Println(sheetData)
		for i, datum := range sheetData {
			Cell, _ := excelize.CoordinatesToCellName(i + 1, i)
			nFile.SetCellValue(sheetName, Cell,  datum)
		}
	}
	if err := nFile.SaveAs(file); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
}

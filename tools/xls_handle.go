package tools

import (
	"fmt"
	"github.com/extrame/xls"
)


func XlsMarge(file1 string, file2 string)  {
	// 打开并读取文件1
	fmt.Println("Reading", file1)
	if xlsFile, err := xls.Open(file1, "utf-8"); err == nil {
		if sheet1 := xlsFile.GetSheet(0); sheet1 != nil {
			fmt.Println("Sheet 1 from file 1:")
			for i := 0; i <= int(sheet1.MaxRow); i++ {
				row := sheet1.Row(i)
				for j := row.FirstCol(); j < row.LastCol(); j++ {
					cell := row.Col(j)
					fmt.Print(cell, "\t")
				}
				fmt.Println()
			}
		}
	} else {
		fmt.Println("Failed to open", file1, ":", err)
	}

	// 打开并读取文件2
	fmt.Println("Reading", file2)
	if xlsFile, err := xls.Open(file2, "utf-8"); err == nil {
		if sheet2 := xlsFile.GetSheet(0); sheet2 != nil {
			fmt.Println("Sheet 1 from file 2:")
			for i := 0; i <= int(sheet2.MaxRow); i++ {
				row := sheet2.Row(i)
				for j := row.FirstCol(); j < row.LastCol(); j++ {
					cell := row.Col(j)
					fmt.Print(cell, "\t")
				}
				fmt.Println()
			}
		}
	} else {
		fmt.Println("Failed to open", file2, ":", err)
	}
}
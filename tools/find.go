package tools

import (
	"encoding/csv"
	"fmt"
	"github.com/extrame/xls"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)
func Find(dir string, fv string, fl string)  {
	parts := strings.Split(fl, ",")
	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 遍历目录中的文件（和子目录）
	for _, entry := range entries {
		flag := false
		// 使用 FileInfo，判断是否为文件
		info, err := entry.Info()
		if err != nil {
			fmt.Println(err)
		}
		// 如果不是目录，则打印文件名
		if !info.IsDir() {
			if filename := entry.Name(); filename[0] == '~' {
				continue
			}
			file := dir + "/" + entry.Name()
			ext := filepath.Ext(file)
			var  rows [][]string
			if ext == ".xlsx" {
				_, rows = OpenXlsx(file)
			} else if ext == ".csv" {
				rows = openCsv(file)
			} else if ext == ".xls" {
				rows = openXls(file)
			}
			for _, row := range rows[1:] {
				for _, part := range parts {
					if row[HeaderMap[part]] == fv {
						flag = true
						break
					}
				}
				if flag {
					fmt.Println(fv)
					break
				}
			}
			if flag {
				fmt.Println(entry.Name())
			}
		}
	}
}

func openCsv(file string) [][]string {
	var recordsData [][]string
	srcFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer srcFile.Close()

	reader := csv.NewReader(srcFile)
	reader.Comma = ','
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				// 文件读取完毕
				break
			}
		}
		// 处理包含换行符的字段
		//fmt.Println(record)
		/*for i, field := range record {
			fmt.Println(field)
			// 如果字段包含换行符但不被双引号包裹，则合并后续字段直到找到包含换行符的字段结束位置
			for strings.Count(field, "\"")%2 != 0 && strings.Count(field, "\n") == 0 {
				// 读取下一个字段，直到找到包含换行符的字段结束位置
				nextField, err := reader.Read()
				if err != nil {
					fmt.Println("Error:", err)
					return nil
				}
				// 将合并的字段添加到当前字段
				field += "\n" + strings.Join(nextField, ",")
			}
			record[i] = field
		}*/
		recordsData = append(recordsData, record)
	}
	// 使用 csv.NewReader 读取字符串
	/*for scanner.Scan() {
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		reader.Comma = ','
		reader.LazyQuotes = true
		record, err := reader.Read()
		recordsData = append(recordsData, record)
		if err != nil {
			break
			//fmt.Println("读取 CSV 行失败: %v", err)
		}
	}*/
	/*if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时出错: %v", err)
		return nil
	}*/
	return recordsData
}

func openXls(file string) [][]string  {
	var recordsData [][]string
	newFile, err := xls.Open(file, "utf-8")
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return nil
	}

	// 获取工作簿中的第一个工作表
	sheet := newFile.GetSheet(0)
	if sheet == nil {
		fmt.Println("Sheet not found")
		return nil
	}

	// 遍历工作表中的每一行
	for i := 0; i <= int(sheet.MaxRow); i++ {
		index := sheet.Row(1)
		row := sheet.Row(i)

		// 读取每一列的数据
		cellValues := make([]string, 0)
		for j := 0; j < index.LastCol(); j++ {
			cell := row.Col(j)
			cellValues = append(cellValues, cell)
		}
		recordsData = append(recordsData, cellValues)
	}

	return recordsData
}

func FindToSheets(file string, ml string, fl string, accountId string)  {
	findFile, rows := OpenXlsx(file)
	for _, row := range rows {
		if accountId != row[3] {
			continue
		}
		sheetName = row[3]
		sheetData, _ := findFile.GetRows(sheetName)
		var findVal map[string]int
		var flagVal map[string]int
		findVal = make(map[string]int)
		flagVal = make(map[string]int)
		for _, datum := range sheetData {
			if row[HeaderMap[ml]] ==  datum[HeaderMap[fl]] {
				findVal[row[HeaderMap[ml]]] += 1
				if findVal[row[HeaderMap[ml]]] > 1 {
					flagVal[row[HeaderMap[ml]]] = 1
				}
			}
		}
		for k, datum := range sheetData {
			_, found := flagVal[datum[HeaderMap[fl]]]
			if found {
				style, _ := findFile.NewStyle(&excelize.Style{
					Fill: excelize.Fill{Type: "pattern", Color: []string{"#FF0000"}, Pattern: 1},
				})
				fmt.Println(datum[HeaderMap[fl]], k + 1)
				for i, _ := range datum {
					cell, _ := excelize.CoordinatesToCellName(i + 1, k + 1)
					findFile.SetCellStyle(sheetName, cell,  cell, style)
				}
			}
		}
	}
	if err := findFile.SaveAs(file); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
}


/*func AddSheet(sheetName string, file string)  {
	//newSheetName := "xreo_data"
	records = nil //清空切片

	sql := `SELECT labelNumber, type, charge, fileName, accountId, fileTime, cellIndex, rowIndex FROM xero_data where type = 1 AND accountId !=''`
	rows, _ := database.DBTask.Query(sql)
	defer rows.Close()
	for rows.Next() {
		var xeroType,accountId, cellIndex, rowIndex int
		var charge float64
		var labelNumber, fileName, fileTime string
		rows.Scan(&labelNumber, &xeroType, &charge, &fileName, &accountId, &fileTime, &cellIndex, &rowIndex)
		rowData := []string{labelNumber, strconv.FormatFloat(charge, 'f', -2, 64), fileName, strconv.Itoa(accountId), fileTime}
		records = append(records, rowData)
	}
	CreateSheet(file, sheetName)
	//newFileHandle("./files/xero_data.xlsx")
	//CreateSheet(file, newSheetName)
}*/


func Check(dir string, connoteValues map[string][]string)  {
	//fileHandel, rows := OpenXlsx(file)
	//connoteValues := ReadConnote(fileHandel, rows, file, isAddSheet) //去重运单号
	Files = nil
	fileMapData := OpenDir(dir)
	var wg sync.WaitGroup
	resultChan = make(chan SearchResult, 15)
	wg.Add(len(fileMapData) - len(NotChargeCell))
	//fmt.Println("fileMapData:", len(fileMapData))
	for _, f := range fileMapData {
		_, found := NotChargeCell[f.FileName] //说明当前文件不是利润文件
		if found {
			continue
		}
		//fmt.Println("f.FileName:", f.FileName)
		go func(searchValue map[string][]string, fileName string, fileTime string) {
			defer wg.Done()
			findMultipleLabelsInFile(searchValue, fileName, fileTime)
		}(connoteValues, f.FileName, f.FileTime)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	//taskDb, _ :=  database.GetTaskDB()
	sql := "INSERT INTO xero_data (labelNumber, type, charge, fileName, accountId, fileTime, cellIndex, rowIndex) VALUES"
	insert := false
	for {
		select {
		case v, ok := <-resultChan:
			if !ok {
				if insert {
					sql = sql[:len(sql)-1]
					_, err = taskDb.Exec(sql)
				}
				if err != nil {
					fmt.Println("执行失败", err)
					return
				}
				return
			}
			insert = true
			sql += fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d'),", v.Value, "1", v.charge, v.fileMapData.FileName,v.accountId, v.fileMapData.FileTime, v.cellIndex, v.rowIndex)
		}
	}
}


func findMultipleLabelsInFile(searchValue map[string][]string, fileName string, fileTime string)  {
	chargeCell = []string{
		"charge",
		"compensation_adjustment",
		"total_amount",
		"Actual Charge (Article)",
	}

	ext := filepath.Ext(fileName)
	var  rowsData [][]string
	if ext == ".xlsx" {
		_, rowsData = OpenXlsx(fileName)
	} else if ext == ".csv" {
		rowsData = openCsv(fileName)
	} else if ext == ".xls" {
		rowsData = openXls(fileName)
	}
	chargeIndex := -1
	for rowIndex, row := range rowsData {
		if len(row) < 1 {
			continue
		}
		for cellIndex, v := range row {
			if rowIndex == 0 {
				if isInSlice(v, chargeCell) {
					chargeIndex =  cellIndex
					break
				}
			}
			for _, labels := range searchValue {

				labelNum := len(labels)
				labelSlice := strings.Split(v, ",")
				if labelNum != len(labelSlice) {
					continue
				}
				//sort.Strings(labels)
				sort.Strings(labelSlice)

				labelsStr := strings.Join(labels, ",")
				labelSliceStr := strings.Join(labelSlice, ",")
				if labelsStr != labelSliceStr {
					continue
				}
				//fmt.Println(labels, "--", labelSlice)
				mutex.Lock() // 加锁
				//fmt.Println(v, row[chargeIndex], cellIndex, rowIndex, fileName)
				resultChan <- SearchResult{
					labelSliceStr,
					row[chargeIndex],
					cellIndex,
					rowIndex,
					"",
					&fileMap{
						FileName: fileName,
						FileTime: fileTime,
					}}
				mutex.Unlock() // 解锁
			}
		}
		if chargeIndex < 0 {
			//fmt.Println(fileName)
			break
		}
	}
}

func ReadConnote(fileHandel *excelize.File, rows [][]string, file string, isAddSheet bool) map[string][]string {
	defer fileHandel.Close()
	ConnoteValues := make(map[string][]string)
	SearchValues := make(map[string][]string)
	addSheetDataIndex := make(map[string]int)
	//addSheetData := make(map[string][]string)
	for i, row := range rows[1:] {
		if len(row) > 0 {
			addSheetDataIndex[row[0]] = i+1
			ConnoteValues[row[1]] = append(ConnoteValues[row[1]], row[0])
		}
	}

	NewSheetName := "Sheet2"
	index, _ := fileHandel.NewSheet(NewSheetName)
	fileHandel.SetActiveSheet(index)

	headers := []string{
		"ArticleID",
		"Connote",
		"Plugin",
		"客户",
		"email",
		"收费方式",
		"提交时间",
		"真实成本",
		"成本日期",
		"真实利润",
		"利润日期",
		"利润率",
		"标准费率",
		"利润差",
	}
	style, _ := fileHandel.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // 从A1开始
		fileHandel.SetCellValue(NewSheetName, cell, header)
		fileHandel.SetCellStyle(NewSheetName, cell, cell, style)
	}

	//fmt.Println(NewSheetName)
	line := 1
	for k, val := range ConnoteValues {
		if len(val) > 1 {
			rowData := make([]string, 0)
			cost :=  0.0
			SearchValues[k] = val
			for _, s := range val {
				if index, found := addSheetDataIndex[s]; found {
					rowData = rows[index]
					costVal, _ := strconv.ParseFloat(rows[index][7], 64)
					cost += costVal
					//break
				}
			}
			sort.Strings(val)
			rowData[0] = strings.Join(val, ",")
			//accountId , _ := strconv.Atoi(rowData[2])
			//rowData[2] =  fmt.Sprintf("%d", accountId)
			//cst, _ := strconv.ParseFloat(rowData[6], 64)
			//rowData[6] = fmt.Sprintf("%.2f", cst)
			//rowData[7] = fmt.Sprintf("%.2f", cost)
			for i, datum := range rowData {
				if i == 12 || i == 13 {
					continue
				}
				if i == 3 {
					accountId , _ := strconv.Atoi(rowData[3])
					cell, _ := excelize.CoordinatesToCellName(i+1, line + 1) // 从A1开始
					fileHandel.SetCellValue(NewSheetName, cell, accountId)
				} else if i == 7 {
					cell, _ := excelize.CoordinatesToCellName(i+1, line + 1) // 从A1开始
					fileHandel.SetCellValue(NewSheetName, cell, cost)
				} else {
					cell, _ := excelize.CoordinatesToCellName(i+1, line + 1) // 从A1开始
					fileHandel.SetCellValue(NewSheetName, cell, datum)
				}
			}
			line++
		}
	}
	if err := fileHandel.SaveAs(file); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
	return SearchValues
}
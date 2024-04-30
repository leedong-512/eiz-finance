package tools

import (
	"bufio"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	resultChan chan SearchResult
	mutex      sync.Mutex // 互斥锁
	chargeCell []string
	notChargeFile *os.File
	NotChargeCell map[string]bool
	startTime = 0
	endTime = 0
)

func NotChargeCellHandel()  {
	file, err := os.Open("./files/notChargeFileName.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		// 读取一行数据
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// 文件读取完毕
				break
			} else {
				fmt.Println("Error:", err)
				return
			}
		}

		// 处理每一行数据
		NotChargeCell[line] = true
	}
}

func SetTimeRange(start, end int) {
	startTime = start
	endTime = end
}

func XeroData(dir string)  {
	searchValue := readColumnLabel()
	Files = nil
	fileMapData := OpenDir(dir)
	var wg sync.WaitGroup
	resultChan = make(chan SearchResult, 15)
	//notChargeFile, _ = os.Create("./files/notChargeFileName.txt")
	//defer notChargeFile.Close()
	wg.Add(len(fileMapData) - len(NotChargeCell))
	for _, f := range fileMapData {
		_, found := NotChargeCell[f.FileName]
		if found {
			continue
		}
		go func(searchValue map[string]string, fileName string, fileTime string) {
			defer wg.Done()
			//fmt.Println("正在处理文件:", fileName)
			findValueInFile(searchValue, fileName, fileTime)
		}(searchValue, f.FileName, f.FileTime)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	sql := "INSERT INTO xero_data (labelNumber, type, charge, fileName, accountId, fileTime, cellIndex, rowIndex) VALUES"
	//taskDb, _ := database.GetTaskDB()
	for {
		select {
		case v, ok := <-resultChan:
			if !ok {
				sql = sql[:len(sql)-1]
				_, err = taskDb.Exec(sql)
				if err != nil {
					//fmt.Println("执行失败", err)
					return
				}
				//fmt.Println("执行完毕")
				return
			}
			sql += fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d'),", v.Value, "1", v.charge, v.fileMapData.FileName,v.accountId, v.fileMapData.FileTime, v.cellIndex, v.rowIndex)
		}
	}
}
type SearchResult struct {
	Value string
	charge string
	cellIndex int
	rowIndex int
	accountId string
	fileMapData *fileMap
}
func isInSlice(target string, slice []string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func findValueInFile(searchValue map[string]string, fileName string, fileTime string)  {
	chargeCell = []string{
		"charge",
		"compensation_adjustment",
		"total_amount",
		"Actual Charge (Article)",
		"total_inlcuding_gst",
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

	//var chargeIndex int
	chargeIndex := -1

	for rowIndex, r1 := range rowsData {
		if len(r1) < 1 {
			continue
		}
		for cellIndex, v := range r1 {
			if rowIndex == 0 {
				if isInSlice(v, chargeCell) {
					chargeIndex = cellIndex
					break
				}
			}
			accountId, found:= searchValue[v]
			if chargeIndex > -1 && found {
				chargeVal := r1[chargeIndex]
				if accountId == "319" {
					charge, _ := strconv.ParseFloat(chargeVal, 64)
					chargeFloat := charge + charge * 0.1
					chargeVal = strconv.FormatFloat(chargeFloat, 'f', 2, 64)
				}
				mutex.Lock() // 加锁
				//fmt.Println(fileName, cellIndex, rowIndex, r1, len(r1), chargeIndex, r1[chargeIndex])
				resultChan <- SearchResult{v, chargeVal, cellIndex, rowIndex, accountId, &fileMap{
					FileName: fileName,
					FileTime: fileTime,
				}}
				mutex.Unlock() // 解锁
				break
			}
		}
		if chargeIndex < 0 {
			//fmt.Println(fileName)
			//_, err = notChargeFile.WriteString(fileName + "\n")
			break
		}
	}
}


func readColumnLabel() map[string]string {
	LabelValues := make(map[string]string)
	for _, row := range records[1:] {
		if len(row) > 0 {
			LabelValues[row[0]] = row[3]
		}
	}
	return LabelValues
}

type fileMap struct {
	FileName string
	FileTime string
}
var Files []fileMap
func OpenDir(dir string) []fileMap {
	//files = nil
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// 遍历目录中的文件（和子目录）
	for _, entry := range entries {
		// 使用 FileInfo，判断是否为文件
		info, err := entry.Info()
		if err != nil {
			fmt.Println(err)
		}
		// 如果不是目录，则打印文件名
		if !info.IsDir() {
			lastDir := filepath.Base(dir)
			if filename := entry.Name(); filename[0] == '~' {
				continue
			}
			filename := dir + "/" + entry.Name()
			Files = append(Files,  fileMap{
				FileName: filename,
				FileTime: lastDir,
			})
		} else {
			dirNameKey := strings.ReplaceAll(entry.Name(), "-", "")
			dirName, _ := strconv.Atoi(dirNameKey)
			if startTime != 0 && dirName < startTime {
				//fmt.Println("跳过", entry.Name())
				continue
			}
			if endTime != 0 && dirName > endTime {
				//fmt.Println("跳过", entry.Name())
				continue
			}

			OpenDir(dir + "/" + entry.Name())
		}
	}

	return Files

}

func NewSheet()  {
	//newSheetName := "xreo_data"
	records = nil //清空切片
	//taskDB, _ := database.GetTaskDB()
	sql := `SELECT labelNumber, type, charge, fileName, accountId, fileTime, cellIndex, rowIndex FROM xero_data where type = 1 AND accountId !=''`
	rows, _ := taskDb.Query(sql)
	defer rows.Close()
	for rows.Next() {
		var xeroType,accountId, cellIndex, rowIndex int
		var charge float64
		var labelNumber, fileName, fileTime string
		rows.Scan(&labelNumber, &xeroType, &charge, &fileName, &accountId, &fileTime, &cellIndex, &rowIndex)
		rowData := []string{labelNumber, strconv.FormatFloat(charge, 'f', -2, 64), fileName, strconv.Itoa(accountId), fileTime}
		records = append(records, rowData)
	}

	newFileHandle("./files/xero_data.xlsx")
	//CreateSheet(file, newSheetName)
}

func NewSheet2()  {
	records = nil //清空切片

	//taskDB, _ := database.GetTaskDB()
	sql := `SELECT labelNumber, type, charge, fileName, accountId, fileTime, cellIndex, rowIndex FROM xero_data where type = 1 AND accountId =''`
	rows, _ := taskDb.Query(sql)
	defer rows.Close()
	for rows.Next() {
		var xeroType,accountId, cellIndex, rowIndex int
		var charge float64
		var labelNumber, fileName, fileTime string
		rows.Scan(&labelNumber, &xeroType, &charge, &fileName, &accountId, &fileTime, &cellIndex, &rowIndex)
		rowData := []string{labelNumber, strconv.FormatFloat(charge, 'f', -2, 64), fileName, strconv.Itoa(accountId), fileTime}
		records = append(records, rowData)
	}

	newFileHandle("./files/xero_data1.xlsx")
	//CreateSheet(file, newSheetName)
}


func newFileHandle(nFile string)  {
	if !fileExists(nFile) {
		f := excelize.NewFile()
		// 创建一个工作表
		index, _ := f.NewSheet(sheetName)
		f.SetActiveSheet(index)
		if err := f.SaveAs(nFile); err != nil {
			fmt.Println(err)
			return
		}
	}
	newFile, _ := OpenXlsx(nFile)

	defer newFile.Close()
	headers := []string{"labelNumber", "charge", "fileName", "accountId", "fileTime"}
	style, _ := newFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // 从A1开始
		newFile.SetCellValue(sheetName, cell, header)
		newFile.SetCellStyle(sheetName, cell, cell, style)
	}

	for i, record := range records {
		for i2, val := range record {
			cell, _ := excelize.CoordinatesToCellName(i2 + 1, i + 2)
			newFile.SetCellValue(sheetName, cell, val)
		}
	}
	if err := newFile.SaveAs(nFile); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
}

func GetAccountId() []string {
	//taskDB, _ := database.GetTaskDB()
	sql := `SELECT accountId FROM xero_data GROUP BY accountId`
	rows, _ := taskDb.Query(sql)
	var accountGroup []string
	for rows.Next() {
		var accountId int
		rows.Scan(&accountId)
		accountGroup = append(accountGroup, strconv.Itoa(accountId))
	}

	return accountGroup
}
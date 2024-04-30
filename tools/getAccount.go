package tools

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/extrame/xls"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)
var (
	data [][]string
	once  sync.Once
	err error
	records [][]string
	sheetName = "Sheet1"
	filePath string
)

/*func init()  {
	sheetName = "Sheet1"
}*/

func SetSheetName(name string)  {
	sheetName = name
}

func GetAccount(baseFile string, accountId int, exp string) {
	newFile, err := excelize.OpenFile(baseFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer newFile.Close()

	rows, err := newFile.GetRows(sheetName)
	if err != nil {
		fmt.Println("Error getting rows:", err)
		return
	}
	lookupValues := make([]string, 0)
	for _, row := range rows[1:] {
		lookupValues = append(lookupValues, row[0])
	}

	lookupValuesString := strings.Join(lookupValues, "','")
	var query string
	if exp == "aupost" {
		query = fmt.Sprintf(`SELECT l.labelNumber,CONVERT_TZ(c.submitted_at,'Etc/UTC','Australia/Melbourne') AS submit_at,CASE
	WHEN c.plugin_id IN (10,120) THEN 'Superfreight'
	ELSE 'Shipvia'
END AS sys, a.id, a.email, a.credit FROM Fulfillments_labels l LEFT JOIN Fulfillments_consignments c ON l.consignment_id = c.id LEFT JOIN accounts a ON c.account_id = a.id WHERE c.plugin_id IN (10,120,159) AND a.id != 513 AND l.labelNumber IN ('%s')`,  lookupValuesString)
	}
	if exp == "tnt" || exp == "hunter" || exp == "aramex" || exp == "pfl" {
		query = fmt.Sprintf(`SELECT l.labelNumber,CONVERT_TZ(c.submitted_at,'Etc/UTC','Australia/Melbourne') AS submit_at,CASE
	WHEN c.plugin_id IN (10,120) THEN 'Superfreight'
	ELSE 'Shipvia'
END AS sys, a.id, a.email, a.credit FROM Fulfillments_labels l LEFT JOIN Fulfillments_consignments c ON l.consignment_id = c.id LEFT JOIN accounts a ON c.account_id = a.id WHERE c.plugin_id IN (10,120,159) AND a.id != 513 AND l.labelNumber IN ('%s') GROUP BY l.labelNumber`,  lookupValuesString)
	}


	//defer db.Close()
	rowsSqlData, _ := db.Query(query)
	defer rowsSqlData.Close()
	i := 0
	sqlData := make(map[string][]string)
	for rowsSqlData.Next() {
		var labelNumber, id, email, credit, sys string
		var submitAt sql.NullString
		rowsSqlData.Scan(&labelNumber, &submitAt, &sys, &id, &email, &credit)

		creditStr := "成本+"
		if sys == "Superfreight" {
			creditStr = "报价"
		}
		submitAtStr := ""
		formattedTime := ""
		if submitAt.Valid {
			submitAtStr = submitAt.String
			// 解析格式，对应于 yyyy-mm-dd hh:mm:ss
			layout := "2006-01-02 15:04:05"
			t, err := time.Parse(layout, submitAtStr)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}
			// 格式化为 dd/mm/yyyy 格式
			formattedTime = t.Format("02/01/2006")
		}

		sqlData[labelNumber] = []string{labelNumber, formattedTime, sys, id, email, creditStr}
		i++
	}
	for k, v := range rows[1:] {
		data, found := sqlData[v[0]]
		if !found {
			continue
		}
		pluginCell, _ := excelize.CoordinatesToCellName(3, k+2) //plugin
		newFile.SetCellValue(sheetName, pluginCell, data[2])
		accountCell, _ := excelize.CoordinatesToCellName(4, k+2) //客户
		accountId, _ := strconv.Atoi(data[3])
		newFile.SetCellValue(sheetName, accountCell, accountId)
		emailCell, _ := excelize.CoordinatesToCellName(5, k+2) //邮箱
		newFile.SetCellValue(sheetName, emailCell, data[4])
		creditCell, _ := excelize.CoordinatesToCellName(6, k+2) //收费方式
		newFile.SetCellValue(sheetName, creditCell, data[5])
		SubmitAtCell, _ := excelize.CoordinatesToCellName(7, k + 2) //提交时间
		newFile.SetCellValue(sheetName, SubmitAtCell, data[1])
	}
	if accountId != 0 {
		query = fmt.Sprintf(`SELECT id, email, "offline" as sys FROM accounts WHERE id = %d`,  accountId)
		if accountId == 6713 {
			query = fmt.Sprintf(`SELECT id, email, "Shipvia" as sys FROM accounts WHERE id = %d`,  accountId)
		}
		rowsSqlData := db.QueryRow(query)
		var eweId, eweEmail, eweSys string
		rowsSqlData.Scan(&eweId, &eweEmail, &eweSys)

		for k, v := range rows[1:] {
			if v[HeaderMap["C"]] == "" && v[HeaderMap["D"]] == "" && v[HeaderMap["E"]] == "" &&
				v[HeaderMap["F"]] == "" && v[HeaderMap["G"]] == "" {
				//fmt.Println(eweId, eweEmail, eweSys)
				pluginCell, _ := excelize.CoordinatesToCellName(HeaderMap["C"] + 1, k + 2) //plugin
				//fmt.Println(pluginCell , k + 2)
				newFile.SetCellValue(sheetName, pluginCell,  eweSys)
				accountId, _ := strconv.Atoi(eweId)
				accountCell, _ := excelize.CoordinatesToCellName(HeaderMap["D"] + 1, k + 2) //account
				//fmt.Println(accountCell , k + 2)
				newFile.SetCellValue(sheetName, accountCell,  accountId)
				emailCell, _ := excelize.CoordinatesToCellName(HeaderMap["E"] + 1, k + 2) //account
				//fmt.Println(emailCell , k + 2)
				newFile.SetCellValue(sheetName, emailCell,  eweEmail)
				typeCell, _ := excelize.CoordinatesToCellName(HeaderMap["F"] + 1, k + 2) //type
				//fmt.Println(typeCell , k + 2)
				newFile.SetCellValue(sheetName, typeCell,  "成本+")
			}
		}
	}
	if err := newFile.SaveAs(baseFile); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
}

func WriteCsv(data [][]string)  {
	file, err := os.Create("new_data.csv")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	// 写入UTF-8 BOM
	_, err = file.WriteString("\xEF\xBB\xBF")
	if err != nil {
		log.Fatal("Failed to write UTF-8 BOM", err)
	}

	// 创建一个 csv.Writer
	writer := csv.NewWriter(file)
	defer writer.Flush() // 确保所有缓存的数据都写入到文件中

	// 写入 CSV 文件的头部
	header := []string{
		"ArticleID",
		"Connote",
		"PluginId",
		"客户",
		"email",
		"收费方式",
		"提交时间",
		"真实成本",
		"成本日期",
		"真实利润",
		"利润日期",
		"利润率",
	}
	if err := writer.Write(header); err != nil {
		log.Fatalln("Error writing header to csv:", err)
	}


	// 遍历用户数据并写入
	for _, val := range data {
		fmt.Println(val)
		if err := writer.Write(val); err != nil {
			log.Fatalln("Error writing record to csv:", err)
		}
	}

}

func ReadFileData(baseFile string, sheetIndex int) [][]string {
	records = nil
	filePath = baseFile
	suffix := checkSuffix(filePath)
	if suffix == ".csv" || suffix == ".CSV" {
		/*srcFile, err := os.Open(baseFile)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		defer srcFile.Close()
		csvReader := csv.NewReader(srcFile)
		records, err = csvReader.ReadAll()*/
		srcFile, err := os.Open(baseFile)
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
			records = append(records, record)
		}

		/*if err != nil {
			fmt.Println("Error reading CSV file:", err)
			return nil
		}*/
	} else if suffix == ".xlsx" {
		newFile, err := excelize.OpenFile(filePath)
		defer newFile.Close()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		sheetName = newFile.GetSheetName(sheetIndex)
		//fmt.Println("sheetName:", sheetName)
		records, err = newFile.GetRows(sheetName)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	} else if suffix == ".xls" {
		newFile, err := xls.Open(filePath, "utf-8")
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
		//fmt.Println(sheet.MaxRow)
		for i := 0; i <= int(sheet.MaxRow); i++ {
			index := sheet.Row(1)
			row := sheet.Row(i)

			// 读取每一列的数据
			cellValues := make([]string, 0)
			//fmt.Println(i, index.LastCol())
			for j := 0; j < index.LastCol(); j++ {
				cell := row.Col(j)
				//fmt.Println(cell)
				cellValues = append(cellValues, cell)
			}
			//fmt.Println(cellValues)
			records = append(records, cellValues)
		}
	}

	return records
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

/**
 * 创建文件
 * @param {string} nFile 文件名
 * @return {boolean} 成功返回true，否则返回false
 * @author leeDong
 * @date 2022-08-18 16:57:07
 */

func SetHeader(nFile string) {
	//_, err := os.Stat(nFile)
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
	style, _ := newFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // 从A1开始
		newFile.SetCellValue(sheetName, cell, header)
		newFile.SetCellStyle(sheetName, cell, cell, style)
	}
	if err := newFile.SaveAs(nFile); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
	fmt.Println("\033[32m", "核算文件创建..................................... ok", "\033[0m")
}
/**
 * 读取excel文件，并写入到新的excel文件中
 * @param {string} baseFile 原始文件路径
 * @param {string} nFile 新的文件路径
 * @param {string} location 目标位置
 * @returns {void}
 * @example
 * Copy("D:\\test.xlsx", "D:\\test1.xlsx", "1,2")
 */

func Copy(nFile string, location string, unique bool)  {

	parts := strings.Split(location, ",")
	// 创建一个整数数组用于存储转换后的值
	srcLocation, _ := HeaderMap[parts[0]]
	targetLocation, _ := HeaderMap[parts[1]]
	newFile, _ := OpenXlsx(nFile)
	var uniqueFlag map[string]bool
	uniqueFlag = make(map[string]bool)
	num := 1
	for i, record := range records {
		if i == 0 {
			continue
		}
		data := record[srcLocation]
		if data == "" {
			continue
		}
		if unique {
			_, found := uniqueFlag[data]
			if !found {
				uniqueFlag[data] = true
				cell, _ := excelize.CoordinatesToCellName(targetLocation + 1, num+1)
				newFile.SetCellValue(sheetName, cell, data)
				num++
			}
			continue
		}

		cell, _ := excelize.CoordinatesToCellName(targetLocation + 1, i+1)
		newFile.SetCellValue(sheetName, cell, data)
	}
	if err := newFile.SaveAs(nFile); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}

	fmt.Println("\033[32m", "Label列复制....................................... ok", "\033[0m")
}

/**
 * 查找
 * @param nFile string
 * @param sl string
 * @param ml int
 * @return
 * @author leeDong
 * @date 2023-03-27 16:52
 */

func XLookUp(nFile string, sl string, ml string)  {
	parts := strings.Split(sl, ",")
	// 创建一个整数数组用于存储转换后的值
	slc, _ := HeaderMap[parts[0]]
	slMain, _ := HeaderMap[parts[1]]
	// 假设CSV的第一行是标题，Consignment Article ID和Consignment ID分别是第一和第二列
	// 创建查找映射
	lookupMap := make(map[string]string)
	for _, record := range records[1:] {
		//fmt.Println(record[slMain], record[slc])
		lookupMap[record[slMain]] = record[slc]
	}

	lookupValues := make([]string, 0)
	newFile, rows := OpenXlsx(nFile)
	if err != nil {
		fmt.Println("Error getting rows:", err)
		return
	}

	//对应列的数据
	for _, row := range rows[1:] {
		lookupValues = append(lookupValues, row[0])
	}

	// 执行查找并收集结果
	results := make([]string, 0)
	for _, val := range lookupValues {
		if consignmentID, found := lookupMap[val]; found {
			results = append(results, consignmentID)
		} else {
			results = append(results, "")
		}
	}

	for i, result := range results {
		cell, _ := excelize.CoordinatesToCellName(HeaderMap[ml] + 1, i+2)
		newFile.SetCellValue(sheetName, cell, result)
	}

	if err := newFile.SaveAs(nFile); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
	fmt.Println("\033[32m", "运单写入............................................... ok", "\033[0m")
}


type Record struct {
	AmdAmt              float64
	ConsignmentArticleID string
	InvoiceTime			string
}

func CreateSheet(nFile string, newSheetName string)  {
	newFile, err := excelize.OpenFile(nFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 检查目标工作表是否存在，如果不存在则创建
	index, err := newFile.GetSheetIndex(newSheetName)
	if index > -1 {
		return
	}
	newFile.NewSheet(newSheetName)
	for i, record := range records {
		if i == 0 {
			for i2, val := range record {
				cell, _ := excelize.CoordinatesToCellName(i2 + 1, i + 1)
				newFile.SetCellValue(newSheetName, cell, val)
			}
			continue
			//os.Exit(0)
		}

		for i2, val := range record {
			cell, _ := excelize.CoordinatesToCellName(i2 + 1, i + 1)
			newFile.SetCellValue(newSheetName, cell, val)
		}
	}

	if err := newFile.Save(); err != nil {
		fmt.Println("保存文件失败:", err)
		return
	}
	fmt.Println("\033[32m", "新增sheet:",  newSheetName, "............................................... ok", "\033[0m")
}

func GetSheet(openFile *excelize.File, getSheetName string, local string) map[string]int {
	rows, _ := openFile.GetRows(getSheetName)
	dat := make(map[string]int, 0)
	for i, val := range rows[1:] {
		dat[val[HeaderMap[local]]] = i + 2
	}

	return dat

}
/**
 * 假设CSV的第一行是标题，Amd Amt和Consignment Article ID分别是第二和第三列
 * 创建查找映射
 * @param records 包含所有记录的CSV数据
 * @return 映射
 */

func Sumifs(nFile string, suml string, rl1 string, rl2 string, wl string, accountId string, profitMargin bool, timeLocation string)  {

	var getSheetName string
	if accountId != "" {
		fileNameWithExt := filepath.Base(filePath)
		ext := filepath.Ext(fileNameWithExt)
		fileName := fileNameWithExt[:len(fileNameWithExt)-len(ext)]
		getSheetName = fileName
		CreateSheet(nFile, getSheetName)
	}
	var recordsData []Record
	for _, record := range records {
		parts := strings.Split(record[HeaderMap[rl1]], "|")
		for _, part := range parts {
			InvoiceTime := ""
			if timeLocation != "" && len(record) > (HeaderMap[timeLocation] + 1) {
				InvoiceTime = record[HeaderMap[timeLocation]]
			}
			amdAmt, _ := strconv.ParseFloat(record[HeaderMap[suml]], 64)
			consignmentArticleID := part
			recordsData = append(recordsData,  Record{AmdAmt: amdAmt, ConsignmentArticleID: consignmentArticleID,  InvoiceTime: InvoiceTime})
		}
		//consignmentArticleID := record[HeaderMap[rl1]]
	}


	newFile, rows := OpenXlsx(nFile)
	defer newFile.Close()


	var sheetData map[string]int
	if accountId != "" {
		sheetData = GetSheet(newFile, getSheetName, rl1)
	}
	for i, row := range rows[1:] {
		if accountId != "" && row[3] != accountId {
			continue
		}
		sum := 0.0
		flag := false
		var profitTime string
		for _, record := range recordsData {
			if record.ConsignmentArticleID == row[HeaderMap[rl2]] {
				sum += record.AmdAmt
				profitTime =  record.InvoiceTime
				flag = true
			}
		}
		gl := false
		if profitMargin {
			pmCell, _ := excelize.CoordinatesToCellName(HeaderMap["L"] + 1, i + 2)
			cost, _ := strconv.ParseFloat(row[7], 64)
			profitMarginValue := (sum - cost) / cost
			newFile.SetCellValue(sheetName, pmCell, profitMarginValue)
			if profitMarginValue < 0 {
				gl = true
			}
		}
		if accountId != "" && flag == false {
			style, _ := newFile.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: []string{"#FF0000"}, Pattern: 1},
			})
			for col := 'A'; col <= 'N'; col++ {
				cell, _ := excelize.CoordinatesToCellName(HeaderMap[string(col)] + 1, i + 2)
				newFile.SetCellStyle(sheetName, cell,  cell, style)
			}
		} else if gl  {
			style, _ := newFile.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFA500"}, Pattern: 1},
			})
			for col := 'A'; col <= 'N'; col++ {
				cell, _ := excelize.CoordinatesToCellName(HeaderMap[string(col)] + 1, i + 2)
				newFile.SetCellStyle(sheetName, cell,  cell, style)
			}
		} else {
			style, _ := newFile.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "none"},
			})
			for col := 'A'; col <= 'N'; col++ {
				cell, _ := excelize.CoordinatesToCellName(HeaderMap[string(col)] + 1, i + 2)
				newFile.SetCellStyle(sheetName, cell,  cell, style)
			}
		}


		cell, _ := excelize.CoordinatesToCellName(HeaderMap[wl] + 1, i + 2)
		_, found := sheetData[row[HeaderMap[rl2]]]
		if found {
			//targetCell := rl1 + fmt.Sprintf("%d", cellIndex)
			targetCell := rl1 + "1"
			linkTarget := fmt.Sprintf("%s!%s", getSheetName, targetCell)
			newFile.SetCellHyperLink(sheetName, cell, linkTarget, "Location")
		}
		newFile.SetCellValue(sheetName, cell, sum)
		if profitTime != "" {
			profitTimeCell, _ := excelize.CoordinatesToCellName(HeaderMap["K"] + 1, i + 2)
			newFile.SetCellValue(sheetName, profitTimeCell,  profitTime)
		}
	}

	if err := newFile.SaveAs(nFile); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}

	if accountId != "" {
		fmt.Println("\033[32m", accountId, "利润率计算........................................... ok", "\033[0m")
	} else {
		fmt.Println("\033[32m", "成本列计算........................................... ok", "\033[0m")
	}


}

/**
 * 获取文件扩展名
 */
func checkSuffix(filePath string) string  {
	extension := filepath.Ext(filePath)
	return extension
}

/**
 * 打开Excel文件
 * @param filePath 文件路径
 * @return 文件对象
 * @return 文件内容
 * @return 错误信息
 */

func OpenXlsx(filePath string) (*excelize.File, [][]string) {
	var xlsxRows [][]string
	newFile, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	//
	xlsxRows, err = newFile.GetRows(sheetName)
	if err != nil {
		name := newFile.GetSheetName(0)
		xlsxRows, err = newFile.GetRows(name)
		//fmt.Println(err)
		//return nil, nil
	}
	return newFile, xlsxRows
}


// 分组整理 account 数据

func Group(filePath string, gl string, isChart bool, isProfit bool)  {
	GroupData := make(map[string][]interface{})
	AccountCost := make(map[string]float64)
	AccountProfit := make(map[string]float64)
	newFile, rows := OpenXlsx(filePath)
	for _, row := range rows[1:] {
		// 成本求和
		_, exists := AccountCost[row[HeaderMap[gl]]]
		cost, _ := strconv.ParseFloat(row[7], 64)
		if exists {
			AccountCost[row[HeaderMap[gl]]] += cost
		} else {
			AccountCost[row[HeaderMap[gl]]] = cost
		}
		// 利润求和
		if isProfit {
			_, profitExists := AccountProfit[row[HeaderMap[gl]]]
			var proVal string
			if len(row) > 9 {
				proVal = row[9]
			} else {
				proVal = "0.0"
			}
			profit, _ := strconv.ParseFloat(proVal, 64)
			if profitExists {
				AccountProfit[row[HeaderMap[gl]]] += profit
			} else {
				AccountProfit[row[HeaderMap[gl]]] = profit
			}
		}


		_, found := GroupData[row[HeaderMap[gl]]]
		if found {
			GroupData[row[HeaderMap[gl]]][2] = AccountCost[row[HeaderMap[gl]]]
			if isProfit {
				GroupData[row[HeaderMap[gl]]][3] = AccountProfit[row[HeaderMap[gl]]]
			}
		} else {
			if isProfit {
				GroupData[row[HeaderMap[gl]]] = append(
					GroupData[row[HeaderMap[gl]]],
					row[2],
					row[3],
					AccountCost[row[HeaderMap[gl]]],
					AccountProfit[row[HeaderMap[gl]]],
				)
			} else {
				GroupData[row[HeaderMap[gl]]] = append(
					GroupData[row[HeaderMap[gl]]],
					row[2],
					row[3],
					AccountCost[row[HeaderMap[gl]]],
				)
			}
		}
	}

	headerData := []string{
		"Email",
		"AccountId",
		"Plugin",
		"总成本",
	}
	if isProfit {
		headerData = append(headerData, "总利润")
	}
	style, _ := newFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	for i, header := range headerData {
		cell, _ := excelize.CoordinatesToCellName(HeaderMap["P"] + i + 1, 18)
		newFile.SetCellValue(sheetName, cell, header)
		newFile.SetCellStyle(sheetName, cell, cell, style)
	}

	num := 0
	for key, value := range GroupData {
		//fmt.Println(key)
		cell, _ := excelize.CoordinatesToCellName(HeaderMap["P"] + 1, 19 + num)
		newFile.SetCellValue(sheetName, cell, key)

		//account ID
		accountCell, _ := excelize.CoordinatesToCellName(HeaderMap["Q"] + 1, 19 + num)
		newFile.SetCellValue(sheetName, accountCell, value[1])

		//Plugin
		pluginCell, _ := excelize.CoordinatesToCellName(HeaderMap["R"] + 1, 19 + num)
		newFile.SetCellValue(sheetName, pluginCell, value[0])

		//成本总和
		costCell, _ := excelize.CoordinatesToCellName(HeaderMap["S"] + 1, 19 + num)
		newFile.SetCellValue(sheetName, costCell, value[2])

		if isProfit {
			//利润总和
			profitCell, _ := excelize.CoordinatesToCellName(HeaderMap["T"] + 1, 19 + num)
			newFile.SetCellValue(sheetName, profitCell, value[3])
		}
		num++
	}


	//柱状图
	if isChart {
		var series []excelize.ChartSeries
		if isProfit {
			series = []excelize.ChartSeries{
				{
					Name:      "Sheet1!$S$18",
					Categories: "Sheet1!$Q$19:$Q$" + fmt.Sprintf("%d", num + 19),
					Values:     "Sheet1!$S$19:$S$" + fmt.Sprintf("%d", num + 19),
				},
				{
					Name:      "Sheet1!$T$18",
					Categories: "Sheet1!$Q$19:$Q$" + fmt.Sprintf("%d", num + 19),
					Values:     "Sheet1!$T$19:$T$" + fmt.Sprintf("%d", num + 19),
				},
			}
		} else {
			series = []excelize.ChartSeries{
				{
					Name:      "Sheet1!$S$18",
					Categories: "Sheet1!$Q$19:$Q$" + fmt.Sprintf("%d", num + 19),
					Values:     "Sheet1!$S$19:$S$" + fmt.Sprintf("%d", num + 19),
				},
			}
		}
		col := &excelize.Chart{
			Type: excelize.Col3D,
			Series: series,
			Title: []excelize.RichTextRun{
				{
					Text: "AUPOST",
				},
			},
			Format: excelize.GraphicOptions{
				ScaleX: 3,
				ScaleY: 2,
				OffsetX: 30,
				OffsetY: 15,
			},
			PlotArea: excelize.ChartPlotArea{
				ShowCatName:     false,
				ShowLeaderLines: false,
				ShowPercent:     false,
				ShowSerName:     false,
				ShowVal:         true,
			},
		}
		if err := newFile.AddChart(sheetName, "V19", col); err != nil {
			fmt.Println(err)
		}
	}

	if err := newFile.SaveAs(filePath); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
}


// 设置单元格的值

func SetValue(filePath string, setKey string, setValue string, accountId string)  {
	newFile, rows := OpenXlsx(filePath)
	for i, row := range rows[1:] {
		if accountId != "" && row[3] != accountId {
			continue
		}
		cell, _ := excelize.CoordinatesToCellName(HeaderMap[setKey] + 1, i + 2)
		fmt.Println(row[0], setValue)
		newFile.SetCellValue(sheetName, cell, setValue)
	}

	if err := newFile.SaveAs(filePath); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}
}

func MergeFile(filePath1 string, filePath2 string)  {

	suffix := checkSuffix(filePath)
	if suffix == ".csv" {
		MergeCsv(filePath1, filePath2)
	} else if suffix == ".xlsx" {
		newFile, err := excelize.OpenFile(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		records, err = newFile.GetRows(sheetName)
		fmt.Println(records)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if suffix == ".xls" {
		newFile, err := xls.Open(filePath, "utf-8")
		if err != nil {
			fmt.Println("Failed to open file:", err)
			return
		}

		// 获取工作簿中的第一个工作表
		sheet := newFile.GetSheet(0)
		if sheet == nil {
			fmt.Println("Sheet not found")
			return
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
			records = append(records, cellValues)
		}
	}
}

func GetCharge(filePath string, accountId string) {
	_, rows := OpenXlsx(filePath)

	lookupValues := make([]string, 0)
	for _, row := range rows[1:] {
		if accountId != "" && row[3] != accountId {
			continue
		}
		lookupValues = append(lookupValues, row[0])
	}
	lookupValuesString := strings.Join(lookupValues, "','")
	id, _ := strconv.Atoi(accountId)
	query := fmt.Sprintf(`SELECT 
l.labelNumber,c.shippingMethod_price_amount,submitted_at 
FROM Fulfillments_labels as l LEFT JOIN Fulfillments_consignments as c 
ON l.consignment_id = c.id 
WHERE l.labelNumber IN ('%s') AND c.account_id = %d`,  lookupValuesString, id)

	rowsSqlData, _ := db.Query(query)

	f := excelize.NewFile()
	headerData := []string{
		"LabelNumber",
		"报价",
		"提交时间",
	}
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})
	for i, header := range headerData {
		cell, _ := excelize.CoordinatesToCellName(i + 1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, style)
	}
	i := 0
	for rowsSqlData.Next() {
		var labelNumber string
		var shippingMethodPriceAmount  float64
		var submittedAt string
		rowsSqlData.Scan(&labelNumber, &shippingMethodPriceAmount, &submittedAt)
		fmt.Println(labelNumber, shippingMethodPriceAmount, submittedAt)
		labelNumberCell, _ := excelize.CoordinatesToCellName(HeaderMap["A"] + 1, i + 2)
		f.SetCellValue(sheetName, labelNumberCell,  labelNumber)
		shippingMethodPriceAmountCell, _ := excelize.CoordinatesToCellName(HeaderMap["B"] + 1, i + 2)
		f.SetCellValue(sheetName, shippingMethodPriceAmountCell,  shippingMethodPriceAmount)
		submittedAtCell, _ := excelize.CoordinatesToCellName(HeaderMap["C"] + 1, i + 2)
		f.SetCellValue(sheetName, submittedAtCell,  submittedAt)
		i++
	}
	// 保存文件
	if err := f.SaveAs(accountId + ".xlsx"); err != nil {
		fmt.Println(err)
	}


}

func Additional(nFile string, nl string, nSumL string, sl string, suml string, accountId string) {
	if accountId != "" {
		CreateSheet(nFile, fmt.Sprintf("recharge_%s", accountId))
	} else {
		filename := filepath.Base(filePath)
		ext := filepath.Ext(filename)
		basename := strings.TrimSuffix(filename, ext)
		CreateSheet(nFile, basename)
	}
	var recordsData []Record
	for _, record := range records {
		amdAmt, _ := strconv.ParseFloat(record[HeaderMap[suml]], 64)
		consignmentArticleID := record[HeaderMap[sl]]
		recordsData = append(recordsData,  Record{AmdAmt: amdAmt, ConsignmentArticleID: consignmentArticleID})
	}

	newFile, rows := OpenXlsx(nFile)
	for i, row := range rows[1:] {
		sum := 0.0
		for _, record := range recordsData {
			if record.ConsignmentArticleID == row[HeaderMap[nl]] {
				sum += record.AmdAmt
			}
		}

		/*style, _ := newFile.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{"#00FF00"}, Pattern: 1},
		})

		if sum > 0 {
			for col := 'A'; col <= 'L'; col++ {
				cell, _ := excelize.CoordinatesToCellName(HeaderMap[string(col)] + 1, i + 2)
				newFile.SetCellStyle(sheetName, cell,  cell, style)
			}
		}*/

		sumNl, _ := strconv.ParseFloat(row[HeaderMap[nSumL]], 64)
		sumNl = sumNl + sum
		//fmt.Println(row[0], sumNl, sum)
		cell, _ := excelize.CoordinatesToCellName(HeaderMap[nSumL]+1, i+2)
		if sum > 0 {
			targetCell := sl + "1"
			filename := filepath.Base(filePath)
			ext := filepath.Ext(filename)
			basename := strings.TrimSuffix(filename, ext)
			linkTarget := fmt.Sprintf("%s!%s", basename, targetCell)
			if accountId != "" {
				linkTarget = fmt.Sprintf("recharge_%s!%s", accountId, targetCell)
			}
			newFile.SetCellHyperLink(sheetName, cell, linkTarget, "Location")
		}
		newFile.SetCellValue(sheetName, cell, sumNl)
	}
	// 保存文件
	if err := newFile.SaveAs(nFile); err != nil {
		fmt.Println(err)
	}
}
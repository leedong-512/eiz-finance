package tools

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func CompleteData(accountIds []string) {
	var labelNumber []string
	for _, row := range records[1:] {
		if row[3] == "" {
			continue
		}
		// 检查账户ID是否已经存在
		if isInSlice(row[3], accountIds) {
			continue
		}
		labelNumber = append(labelNumber, row[0])
	}
	labelNumberString := strings.Join(labelNumber, "','")
	if labelNumberString == "" {
		return
	}
	query := fmt.Sprintf("SELECT b.labelNumber AS labels,a.account_id,a.transaction_num as id,a.transaction_amount AS charge,DATE( a.transaction_date ) AS created FROM t_transactions AS a LEFT JOIN Fulfillments_labels AS b ON a.ref_id = b.consignment_id WHERE b.labelNumber IN ('%s')",  labelNumberString)
	//query := fmt.Sprintf("SELECT fl.labelNumber AS labels, fc.account_id, fc.id, fc.shippingMethod_price_amount AS charge, fc.submitted_at AS created FROM Fulfillments_labels fl LEFT JOIN Fulfillments_consignments fc ON fl.consignment_id = fc.id WHERE fc.account_id!=513 and fl.labelNumber IN ('%s')", labelNumberString)
	//fmt.Println(len(labelNumber), query)
	rowsSqlData, _ := db.Query(query)
	defer rowsSqlData.Close()
	execSql := "INSERT INTO xero_data (labelNumber, type, charge, fileName, accountId, fileTime, cellIndex, rowIndex) VALUES"
	for rowsSqlData.Next() {
		var labels, csId string
		var accountId int
		var charge float64
		var submitAt sql.NullString
		rowsSqlData.Scan(&labels, &accountId, &csId, &charge, &submitAt)
		charge = -charge
		//fmt.Println(labels, accountId, csId, charge, submitAt)
		accountIdStr := strconv.Itoa(accountId)
		csIdStr := csId
		chargeStr := strconv.FormatFloat(charge, 'f', 2, 64)
		submitAtStr := ""
		formattedTime := ""
		if submitAt.Valid {
			submitAtStr = submitAt.String
			// 解析格式，对应于 yyyy-mm-dd hh:mm:ss
			layout := "2006-01-02"
			t, err := time.Parse(layout, submitAtStr)
			if err != nil {
				fmt.Println("Error parsing time:", err)
				return
			}
			// 格式化为 yyyy-mm-dd 格式
			formattedTime = t.Format("2006-01-02")
		}
		execSql += fmt.Sprintf("('%s', '%s', '%s', '%s', '%s', '%s', '%d', '%d'),", labels, "1", chargeStr, csIdStr, accountIdStr, formattedTime, 0, 0)

	}
	execSql = execSql[:len(execSql)-1]
	_, err = taskDb.Exec(execSql)
	if err != nil {
		fmt.Println("执行失败", err)
		return
	}
	//fmt.Println(execSql)
}

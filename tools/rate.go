package tools

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"sync"
)

var (
	rateChanSf = make(chan RateRecord, 20)
	rateChanSv = make(chan RateRecord, 20)
	rateChanOf = make(chan RateRecord, 20)
	rateChanMoreSf = make(chan RateRecord, 20)
	rateChanMoreSv = make(chan RateRecord, 20)
	rateChanMoreOf = make(chan RateRecord, 20)
	wgRateSf sync.WaitGroup
	wgRateSv sync.WaitGroup
	wgRateOf sync.WaitGroup
	wgRateMoreSf sync.WaitGroup
	wgRateMoreSv sync.WaitGroup
	wgRateMoreOf sync.WaitGroup
)



type shippingMethodData struct {
	AccountId  int
	ShippingMethodId int
}
func Rate(fileHandle *excelize.File, rows [][]string, pt string)  {
	var sfData map[string]shippingMethodData
	sfData = make(map[string]shippingMethodData)
	labels := make([]string, 0)
	labels1 := make([]string, 0)
	//labels2 := make([]string, 0)
	num := 1
	for _, row := range rows {
		if len(row) < 10  {
			continue
		}
		if (pt == "sf" || pt == "smf") && row[2] == "Superfreight" {
			num++
			label := strings.Split(row[0], ",")
			labels = append(labels, label[0])
		}
		if (pt == "sv" || pt == "smv") && row[2] == "Shipvia" {
			num++
			//label := strings.Split(row[0], ",")
			labels1 = append(labels1, row[0])
		}
		if  (pt == "of" || pt == "omf") && row[2] == "offline" {
			num++
			//label := strings.Split(row[0], ",")
			labels1 = append(labels1, row[0])
		}
	}
	if len(labels) > 0 {
		lookupValuesString := strings.Join(labels, "','")
		query := fmt.Sprintf(`SELECT l.labelNumber,c.account_id, c.shippingMethod_id FROM Fulfillments_labels l LEFT JOIN Fulfillments_consignments c ON l.consignment_id = c.id LEFT JOIN accounts a ON c.account_id = a.id WHERE c.plugin_id IN (10,120,159) AND a.id != 513 AND l.labelNumber IN ('%s')`,  lookupValuesString)
		rowsSqlData, _ := db.Query(query)
		defer rowsSqlData.Close()
		for rowsSqlData.Next() {
			var labelNumber string
			var accountId, shippingMethodId int
			rowsSqlData.Scan(&labelNumber, &accountId, &shippingMethodId)
			sfData[labelNumber] = shippingMethodData{AccountId: accountId, ShippingMethodId: shippingMethodId}
		}
	}


	if pt == "sf" {
		wgRateSf.Add(num-1)
	}
	if pt == "smf" {
		wgRateMoreSf.Add(num-1)
	}
	if pt == "sv" {
		wgRateSv.Add(num-1)
	}
	if pt == "smv" {
		wgRateMoreSv.Add(num-1)
	}
	if pt == "of" {
		wgRateOf.Add(num-1)
	}
	if pt == "omf" {
		wgRateMoreOf.Add(num-1)
	}

	for i, row := range rows {
		if  i == 0 || len(row) < 10 {
			continue
		}
		if pt == "sf" && row[2] == "Superfreight" {
			go func(rowData []string, sfData map[string]shippingMethodData, index int, pt string) {
				defer wgRateSf.Done()
				rateSf(rowData, sfData, index, pt)
			}(row, sfData, i, pt)
			continue
		}
		if pt == "sv" && row[2] == "Shipvia" {
			go func(rowData []string, index int, pt string) {
				defer wgRateSv.Done()
				rateShipVia(rowData, index, pt)
			}(row, i, pt)
			continue
		}
		if  pt == "smf" && row[2] == "Superfreight" {
			go func(rowData []string, sfData map[string]shippingMethodData, index int, pt string) {
				defer wgRateMoreSf.Done()
				rateSf(rowData, sfData, index, pt)
			}(row, sfData, i, pt)
			continue
		}
		if pt == "smv" && row[2] == "Shipvia" {
			go func(rowData []string, index int, pt string) {
				defer wgRateMoreSv.Done()
				rateShipVia(rowData, index, pt)
			}(row, i, pt)
			continue
		}

		if pt == "of" && row[2] == "offline" {
			/*go func(rowData []string, sfData map[string]shippingMethodData, index int, pt string) {
				defer wgRateOf.Done()
				rateSf(rowData, sfData, index, pt)
			}(row, sfData, i, pt)*/
			go func(rowData []string, index int, pt string) {
				defer wgRateOf.Done()
				rateShipVia(rowData, index, pt)
			}(row, i, pt)
			continue
		}
		if  pt == "omf" && row[2] == "offline" {
			/*go func(rowData []string, sfData map[string]shippingMethodData, index int, pt string) {
				defer wgRateMoreOf.Done()
				rateSf(rowData, sfData, index, pt)
			}(row, sfData, i, pt)*/
			go func(rowData []string, index int, pt string) {
				defer wgRateMoreOf.Done()
				rateShipVia(rowData, index, pt)
			}(row, i, pt)
			continue
		}
	}

	go func(pt string) {
		switch pt {
		case "sf":
			wgRateSf.Wait()
			close(rateChanSf)
		case "sv":
			wgRateSv.Wait()
			close(rateChanSv)
		case "smf":
			wgRateMoreSf.Wait()
			close(rateChanMoreSf)
		case "smv":
			wgRateMoreSv.Wait()
			close(rateChanMoreSv)
		case "of":
			wgRateOf.Wait()
			close(rateChanOf)
		case "omf":
			wgRateMoreOf.Wait()
			close(rateChanMoreOf)
		}
	}(pt)

	switch pt {
	case "sf":
		RateChanHandel(rateChanSf, fileHandle, num)
	case "sv":
		RateChanHandel(rateChanSv, fileHandle, num)
	case "smf":
		RateChanHandel(rateChanMoreSf, fileHandle, num)
	case "smv":
		RateChanHandel(rateChanMoreSv, fileHandle, num)
	case "of":
		RateChanHandel(rateChanOf, fileHandle, num)
	case "omf":
		RateChanHandel(rateChanMoreOf, fileHandle, num)
	}
	if err := fileHandle.Save(); err != nil {
		fmt.Println("Error saving Excel file:", err)
	}

}

func RateChanHandel(chanData chan RateRecord, fileHandle *excelize.File, num int)  {
	bar := pb.StartNew(num-1)
Loop:
	for {
		select {
		case v, ok := <-chanData:
			if !ok {
				if err != nil {
					return
				}
				// 完成进度条
				bar.Finish()
				break Loop
			}
			bar.Increment()
			i := v.Index
			cell, _ := excelize.CoordinatesToCellName(HeaderMap["M"] + 1, i + 1)
			fileHandle.SetCellValue(sheetName, cell, v.StandardRate)
			if !v.Pass {
				if v.Profit > 0 {
					cell, _ = excelize.CoordinatesToCellName(HeaderMap["N"] + 1, i + 1)
					fileHandle.SetCellValue(sheetName, cell, v.ProfitMargin)
					if v.ProfitMargin > 0.5 {
						style, _ := fileHandle.NewStyle(&excelize.Style{
							Fill: excelize.Fill{Type: "pattern", Color: []string{"#0000FF"}, Pattern: 1},
						})
						for col := 'A'; col <= 'N'; col++ {
							cell, _ = excelize.CoordinatesToCellName(HeaderMap[string(col)] + 1, i + 1)
							fileHandle.SetCellStyle(sheetName, cell,  cell, style)
						}
					}
				}
			}
		}
	}
}

type RateRecord struct {
	StandardRate float64
	ProfitMargin float64
	Profit 		 float64
	Pass 		 bool
	Index		 int
}
func rateSf(record []string, sfData map[string]shippingMethodData, index int, pt string) {
	if record[3] == "" {
		return
	}
	label := strings.Split(record[0], ",")
	rateLevelData, _ := sfData[label[0]]
	accountId := rateLevelData.AccountId
	shippingMethodId := rateLevelData.ShippingMethodId
	sql := fmt.Sprintf(`SELECT rateLeve FROM xero_superferight_rate_level where eizAccountId = '%d' AND shippingMethodId = '%d'`, accountId, shippingMethodId)
	row := taskDb.QueryRow(sql)
	var rateLevel float64
	row.Scan(&rateLevel)
	if record[3] == "4832" {
		rateLevel = 10
	}
	standardRate, _ := strconv.ParseFloat(record[11], 64)
	profitMargin := (rateLevel / 100) - standardRate
	pass := false
	if profitMargin < 0 {
		pass = true
	}
	chanData := RateRecord{
		StandardRate: rateLevel / 100,
		ProfitMargin: profitMargin,
		Profit: standardRate,
		Pass: pass,
		Index: index,
	}
	if pt == "sf" {
		rateChanSf <- chanData
	}
	if pt == "smf" {
		rateChanMoreSf <- chanData
	}

}

func rateShipVia(record []string, index int, pt string) {
	if record[3] == "" {
		return
	}
	label := strings.Split(record[0], ",")
	sql := fmt.Sprintf(`SELECT account_id, tracking, chargeCode FROM courier_ISF_consignments WHERE tracking LIKE '%s'`, "%"+label[0]+"%")
	row := Db3Pl.QueryRow(sql)
	var labelNumber string
	var accountId, shippingMethodId int
	row.Scan(&accountId, &labelNumber, &shippingMethodId)

	sql = fmt.Sprintf(`SELECT rateLevel FROM xero_shipvia_rate_level where accId = '%d' AND chargeCode = '%d'`, accountId, shippingMethodId)
	row = taskDb.QueryRow(sql)
	var rateLevel float64
	row.Scan(&rateLevel)
	if rateLevel == 0 {
		rateLevel = 1
	}
	standardRate, _ := strconv.ParseFloat(record[11], 64)
	profitMargin := (rateLevel-1) - standardRate
	pass := false
	if profitMargin < 0 {
		pass = true
	}
	chanData := RateRecord{
		StandardRate: rateLevel-1,
		ProfitMargin: profitMargin,
		Profit: standardRate,
		Pass: pass,
		Index: index,
	}
	if pt == "sv" {
		rateChanSv <- chanData
	}
	if pt == "smv" {
		rateChanMoreSv <- chanData
	}
	if pt == "of" {
		rateChanOf <- chanData
	}
	if pt == "omf" {
		rateChanMoreOf <- chanData
	}
}
package cmd

import (
	config "financeSys/configs"
	"financeSys/tools"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var(
	exp 			string
	Dir 			string
	SumRl1 			string
	SumL 			string
	Sl 			string
	wg1 sync.WaitGroup
	copyIndex string
	startTime int
	endTime int
	sy bool
	far bool
	AutoCmd = &cobra.Command{
		Use: "auto [flags] [value]",
		Long: `
 ______     __  __     ______   ______    
/\  __ \   /\ \/\ \   /\__  _\ /\  __ \   
\ \  __ \  \ \ \_\ \  \/_/\ \/ \ \ \/\ \  
 \ \_\ \_\  \ \_____\    \ \_\  \ \_____\ 
  \/_/\/_/   \/_____/     \/_/   \/_____/`,
		Example: `main auto --nf ./files/cost/hs/kk.xlsx --dir ./files/aupost --costFile ./files/cost/cost2005.csv --exp aupost --cpl N,A --r1 N --sml CM --sl M,N`,
		Run: autoCmd,
	}
)
func init() {
	configPath := "./configs/config.yml"
	config.Initialize(configPath)
	cfg, _ = config.GetConfig()
	defDir := cfg.Parameters.Dir
	defExp := cfg.Parameters.Exp
	defCpl := cfg.Parameters.Cpl
	defSml := cfg.Parameters.Sml
	defSl := cfg.Parameters.Sl
	defR1 := cfg.Parameters.R1
	AutoCmd.Flags().StringVar(&newFile, "nf", "", "核算文件地址")
	AutoCmd.Flags().StringVar(&Dir, "dir", defDir, "目标文件地址")
	AutoCmd.Flags().StringVar(&baseFile, "costFile", "", "成本文件")
	AutoCmd.Flags().StringVar(&exp, "exp", defExp, "快递类型")
	AutoCmd.Flags().StringVar(&copyIndex, "cpl", defCpl, "复制列的定位")
	AutoCmd.Flags().StringVar(&SumRl1, "r1", defR1, "成本列的定位")
	AutoCmd.Flags().StringVar(&SumL, "sml", defSml, "计算列")
	AutoCmd.Flags().StringVar(&Sl, "sl", defSl, "运单列")
	AutoCmd.Flags().IntVar(&startTime, "start", 0, "开始时间")
	AutoCmd.Flags().IntVar(&endTime, "end", 0, "结束时间")
	AutoCmd.Flags().BoolVar(&ewe, "ewe", false, "ewe客户")
	AutoCmd.Flags().BoolVar(&sy, "sunyee", false, "sunyee客户")
	AutoCmd.Flags().BoolVar(&far, "far", false, "far客户")
	RootCmd.AddCommand(AutoCmd)
}

type Params struct {
	starttime string
	endtime   string
	cf 		  string
	nf 		  string
}
func cmdData() {
	//var cmdParams []Params
	dirFile := "./files/cost/aramex"
	/*fileMap := "./configs/fileMap.txt"
	files, _ := os.Open(fileMap)
	scanner := bufio.NewScanner(files)
	defer files.Close()
	for scanner.Scan() {
		line := scanner.Text()
		lineArr := strings.Split(line, "    ")
		baseName := filepath.Base(lineArr[0])
		baseNameArr := strings.Split(baseName, "_")
		lower := strings.ToLower(baseNameArr[0])
		csvData := tools.ReadFileData(lineArr[0], 0)
		t, _ := time.Parse("2/01/2006", csvData[1][1])

		//减少3个月
		start_t := t.AddDate(0, -3, 0)
		startime := start_t.Format("20060102")
		//增加3个月
		end_t := t.AddDate(0, 3, 0)
		endtime := end_t.Format("20060102")

		formattedDate := t.Format("20060102")
		filename := formattedDate + "_" + lower + ".xlsx"
		cmdParams = append(cmdParams,  Params{startime, endtime, lineArr[0], "./files/cost/hs/hunter/" + filename})
	}*/
	entries, err := os.ReadDir(dirFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 遍历目录中的文件（和子目录）
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Println(err)
		}
		if info.IsDir() {
			continue
		}

		filename := entry.Name()
		ext := filepath.Ext(filename)
		if ext != ".csv" && ext != ".xlsx" {
			continue
		}
		//fmt.Println(dirFile + "/" + filename)
		csvData := tools.ReadFileData(dirFile + "/" + filename, 0)
		if csvData[1][1] == "" {
			//fmt.Println(filename)
			continue
		}
		//fmt.Println(csvData[1][1])
		t, _ := time.Parse("02-Jan-2006", csvData[1][1])
		if t.Unix() < 0 {
			t, _ = time.Parse("02-Jan-06", csvData[1][1])
			if t.Unix() < 0 {
				t, _ = time.Parse("2-Jan-06", csvData[1][1])
				if t.Unix() < 0 {
					t, _ = time.Parse("2/01/2006", csvData[1][1])
				}
			}
		}
		if t.Unix() < 0 || t.String() == ""{
			//fmt.Println(filename)
			continue
		}
		//fmt.Println(csvData[1][1], t.String())
		//减少3个月
		start_t := t.AddDate(0, -3, 0)
		startime := start_t.Format("20060102")
		//增加3个月
		end_t := t.AddDate(0, 3, 0)
		endtime := end_t.Format("20060102")

		formattedDate := t.Format("20060102")
		//filepath.Base(filename)
		nameWithoutExtension := strings.TrimSuffix(filename, filepath.Ext(filename))
		fileName := nameWithoutExtension + "_" + formattedDate + ".xlsx"

		startTime, _ = strconv.Atoi(startime)
		endTime, _ = strconv.Atoi(endtime)

		fmt.Println(fmt.Sprintf("go run main.go auto --nf %s --costFile %s --start %d --end %d",  "./files/cost/hs/aramex/" + fileName, dirFile + "/" + filename, startTime, endTime))
	}

	//return cmdParams
}
func autoCmd(cmd *cobra.Command, args []string)  {

	//cmdData()
	handle()
	//cmdParams := cmdData()
	/*for _, param := range cmdParams {
		startTime, _ = strconv.Atoi(param.starttime)
		endTime, _ = strconv.Atoi(param.endtime)
		baseFile = param.cf
		newFile = param.nf
		//fmt.Println("\033[33m", "开始执行命令..............................................", "\033[0m")
		//fmt.Println("\033[33m", fmt.Sprintf("main auto --nf %s --costFile %s --start %d --end %d", newFile, baseFile, startTime, endTime), "\033[0m")
		//handle()
		//fmt.Println("\033[32m", "执行完毕.............................................. END", "\033[0m")
		//fmt.Println(startTime, endTime, baseFile, newFile)
		fmt.Println(fmt.Sprintf("go run main.go auto --nf %s --costFile %s --start %d --end %d", newFile, baseFile, startTime, endTime))

	}*/
	//


}

func handle() {
	if exp == "aupost" || exp == "tnt" || exp == "hunter" || exp == "aramex" || exp == "pfl" {
		//删除xero_data文件
		os.Remove("./files/xero_data.xlsx")
		os.Remove("./files/xero_data1.xlsx")
		fmt.Println("\033[32m", "删除xero_data文件.............................................. ok", "\033[0m")
		//清除原始数据
		tools.Del(1, 1)
		fmt.Println("\033[32m", "数据清除..................................................... ok", "\033[0m")

		tools.ReadFileData(baseFile, 0)
		tools.SetSheetName("Sheet1")
		//设置header & 新建核算文件
		tools.SetHeader(newFile)
		//去重
		if exp == "hunter" || exp == "aramex" {
			tools.Copy(newFile, copyIndex, true)
			fmt.Println("\033[32m", "面单去重.............................................. ok", "\033[0m")
		} else {
			//复制label列
			tools.Copy(newFile, copyIndex, false)
		}

		//运单列
		tools.XLookUp(newFile, Sl, "B")
		//sumifs成本列
		if  exp == "aramex" {
			tools.Sumifs(newFile, SumL, SumRl1, "B", "H", "", false, "")
		} else {
			tools.Sumifs(newFile, SumL, SumRl1, "A", "H", "", false, "")
		}

		//客户信息获取
		if exp == "aupost" {
			tools.GetAccount(newFile, 0, "aupost")
			fmt.Println("\033[32m", "用户信息写入.............................................. ok", "\033[0m")
			if ewe {
				tools.GetAccount(newFile, 2913, "aupost")
				fmt.Println("\033[32m", "EWE特殊用户信息写入.............................................. ok", "\033[0m")
			}
			if sy {
				tools.GetAccount(newFile, 319, "aupost")
				fmt.Println("\033[32m", "SUNYEE特殊用户信息写入.............................................. ok", "\033[0m")
			}
			if far {
				tools.GetAccount(newFile, 6713, "aupost")
				fmt.Println("\033[32m", "FAR特殊用户信息写入.............................................. ok", "\033[0m")
			}
		}
		if exp == "tnt" {
			tools.GetAccount(newFile, 0, "tnt")
			fmt.Println("\033[32m", "用户信息写入.............................................. ok", "\033[0m")
		}
		if exp == "hunter" {
			tools.GetAccount(newFile, 0, "hunter")
			fmt.Println("\033[32m", "用户信息写入.............................................. ok", "\033[0m")
		}
		if exp == "aramex" {
			tools.GetAccount(newFile, 0, "aramex")
			fmt.Println("\033[32m", "用户信息写入.............................................. ok", "\033[0m")
		}
		if exp == "pfl" {
			tools.GetAccount(newFile, 0, "pfl")
			fmt.Println("\033[32m", "用户信息写入.............................................. ok", "\033[0m")
		}


		tools.ReadFileData(newFile, 0)
		//设置时间范围
		tools.SetTimeRange(startTime, endTime)
		fmt.Println("\033[33m", ".......................................数据检索中,预计需要6分钟,请等待...", "\033[0m")
		tools.XeroData(Dir)
		fmt.Println("\033[32m", "检索数据已写入到数据库........................................... ok", "\033[0m")

		var accountIds []string
		accountIds = tools.GetAccountId()
		var newAccountId string
		for i, id := range accountIds {
			if id == "0" {
				fmt.Println("\033[33m", "请处理数据库中accountId为0的数据,并指定新的accountId", "\033[0m")
				fmt.Print("\033[33m", " 请输入新的accountId:", "\033[0m")
				fmt.Scanln(&newAccountId)
				if newAccountId == "" {
					continue
				}
				accountIds[i] = newAccountId
			}
		}
		//数据补全
		tools.ReadFileData(newFile, 0)
		tools.CompleteData(accountIds)
		fmt.Println("\033[32m", "数据补全........................................... ok", "\033[0m")
		tools.SetSheetName("Sheet1")
		tools.NewSheet() //新增 xero_data 文件
		fmt.Println("\033[32m", "创建利润文件：xero_data.xlsx............................. ok", "\033[0m")


		accountIds = tools.GetAccountId()
		fmt.Println("\033[32m", "获取用户分组：", accountIds, "\033[0m")

		fmt.Println("\033[33m", "开始查询成本价格...................................... start", "\033[0m")
		if ewe {
			for _, id := range accountIds {
				if id == "0" {
					continue
				}
				//查找是否需要修改成本价格
				for _, f := range tools.RechargeFile {
					tools.ReadFileData(f, 0)
				}
				tools.Additional(newFile, "A", "H", "C", "V", id)
			}
			fmt.Println("\033[32m", "成本叠加........................................... ok", "\033[0m")
		}

		tools.ReadFileData("./files/xero_data.xlsx", 0)
		fmt.Println("\033[33m", "开始计算利润.......................................... start", "\033[0m")
		//var newAccountId string
		for _, id := range accountIds {
			tools.Sumifs(newFile, "B", "A", "A", "J", id, true, "E")
		}

		var chanNum = 2
		if ewe {
			chanNum = 3
		}
		closeChan := make(chan int, chanNum)
		wg.Add(chanNum)
		fileHandle , rows := tools.OpenXlsx(newFile)
		fmt.Println("\033[33m", "开始计算利润差............................................ start", "\033[0m")
		go func() {
			defer wg.Done()
			fmt.Println("\033[33m", "sf利润差开始计算中............................................ start", "\033[0m")
			tools.Rate(fileHandle, rows,"sf")
			closeChan <- 1
		}()
		go func() {
			defer wg.Done()
			fmt.Println("\033[33m", "sv利润差开始计算中............................................ start", "\033[0m")
			tools.Rate(fileHandle, rows, "sv")
			closeChan <- 2
		}()
		if ewe {
			go func() {
				defer wg.Done()
				fmt.Println("\033[33m", "of利润差开始计算中............................................ start", "\033[0m")
				tools.Rate(fileHandle, rows, "of")
				closeChan <- 3
			}()
		}
		wg.Wait()
		close(closeChan)

		fmt.Println("\033[32m", "单件利润差计算完毕............................................... ok", "\033[0m")


		//摘出一单多件的数据
		tools.SetSheetName("Sheet1")
		fileHandel, rows := tools.OpenXlsx(newFile)
		connoteValues := tools.ReadConnote(fileHandel, rows, newFile, false)
		fmt.Println("\033[32m", "摘出一单多件的数据............................................... ok", "\033[0m")
		if len(connoteValues) > 0 {
			fmt.Println("\033[33m", "数据检索中,预计需要5分钟,请等待...", "\033[0m")
			tools.Check(Dir, connoteValues)
			fmt.Println("\033[33m", "检测数据中多单的运单数据,并写入到数据库中...",  "\033[0m")

			tools.NewSheet2() //创建 xero_data1文件
			fmt.Println("\033[32m", "创建利润文件：xero_data1.xlsx............................. ok", "\033[0m")
			tools.ReadFileData("./files/xero_data1.xlsx", 0)
			fmt.Println("\033[33m", "开始计算多单利润............................................ start", "\033[0m")
			tools.SetSheetName("Sheet2")
			for _, id := range accountIds {
				if id == "0" {
					continue
				}
				tools.Sumifs(newFile, "B", "A", "A", "J", id, true, "E")
			}


			closeChan1 := make(chan int, chanNum)
			wg1.Add(chanNum)
			tools.SetSheetName("Sheet2")
			fileHandel , rows = tools.OpenXlsx(newFile)
			fmt.Println("\033[33m", "开始计算多单利润差............................................ start", "\033[0m")
			go func() {
				defer wg1.Done()
				fmt.Println("\033[33m", "sf利润差开始计算中............................................ start", "\033[0m")
				tools.Rate(fileHandel, rows,"smf")
				closeChan1 <- 1
			}()
			go func() {
				defer wg1.Done()
				fmt.Println("\033[33m", "sv利润差开始计算中............................................ start", "\033[0m")
				tools.Rate(fileHandel, rows, "smv")
				closeChan1 <- 2
			}()
			if ewe {
				go func() {
					defer wg1.Done()
					fmt.Println("\033[33m", "of利润差开始计算中............................................ start", "\033[0m")
					tools.Rate(fileHandel, rows, "omf")
					closeChan1 <- 3
				}()
			}

			wg1.Wait()
			close(closeChan1)
			fmt.Println("\033[32m", "一单多件利润差计算完毕............................................... ok", "\033[0m")
		}
		fmt.Println("\033[32m", "核算表........................................................... ok", "\033[0m")
	}
	tools.CloseDb()
}


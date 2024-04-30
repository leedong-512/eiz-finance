package cmd

import (
	"financeSys/tools"
	"fmt"
	"github.com/spf13/cobra"
	"sync"
)

var(

	wg sync.WaitGroup
	dir string
	newSheet bool
	ImportCmd = &cobra.Command{
		Use: "import [flags] [value]",
		Long: `
 __     __    __     ______   ______     ______     ______  
/\ \   /\ "-./  \   /\  == \ /\  __ \   /\  == \   /\__  _\ 
\ \ \  \ \ \-./\ \  \ \  _-/ \ \ \/\ \  \ \  __<   \/_/\ \/ 
 \ \_\  \ \_\ \ \_\  \ \_\    \ \_____\  \ \_\ \_\    \ \_\ 
  \/_/   \/_/  \/_/   \/_/     \/_____/   \/_/ /_/     \/_/`,
		Example: `main import --baseFile /file.xlsx --dir ./file`,
		Run: importCmd,
	}
)
func init() {
	ImportCmd.Flags().StringVar(&newFile, "nf", "", "文件地址")
	ImportCmd.Flags().StringVar(&dir, "dir", "", "读取文件的地址")
	ImportCmd.Flags().BoolVar(&newSheet, "ns", false, "是否新增sheet")
	RootCmd.AddCommand(ImportCmd)
}

func importCmd(cmd *cobra.Command, args []string)  {
	tools.ReadFileData(newFile, 0)
	tools.XeroData(dir)
	if newSheet {
		tools.NewSheet()
		accountIds := tools.GetAccountId()
		tools.ReadFileData("./files/xero_data.xlsx", 0)
		for _, id := range accountIds {
			tools.Sumifs(newFile, "B", "A", "A", "J", id, true, "E")
		}
		closeChan := make(chan int, 2)
		wg.Add(2)
		fileHandle , rows := tools.OpenXlsx(newFile)
		go func() {
			defer wg.Done()
			fmt.Println("sf利润差开始计算中...")
			tools.Rate(fileHandle, rows,"sf")
			closeChan <- 1
		}()
		go func() {
			defer wg.Done()
			fmt.Println("sv利润差开始计算中...")
			tools.Rate(fileHandle, rows, "sv")
			closeChan <- 2
		}()
		wg.Wait()
		close(closeChan)
	}

	fmt.Println("over")

}

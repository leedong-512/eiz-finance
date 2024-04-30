package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	nSumL string
	ExtraCmd = &cobra.Command{
		Use: "extra [flags] [value]",
		Long: `
 ______     __  __     ______   ______     ______    
/\  ___\   /\_\_\_\   /\__  _\ /\  == \   /\  __ \   
\ \  __\   \/_/\_\/_  \/_/\ \/ \ \  __<   \ \  __ \  
 \ \_____\   /\_\/\_\    \ \_\  \ \_\ \_\  \ \_\ \_\ 
  \/_____/   \/_/\/_/     \/_/   \/_/ /_/   \/_/\/_/`,
		Example: `main extra --nl A --nsumL H --sl C --suml N --baseFile /filebase.xlsx --nf /nf.xlsx`,
		Run: extraCmd,
	}
)
func init() {
	ExtraCmd.Flags().StringVar(&baseFile, "basePath", "", "源文件地址")
	ExtraCmd.Flags().StringVar(&newFile, "nf", "", "目标文件地址")
	ExtraCmd.Flags().StringVar(&l, "nl", "", "目标列的定位")
	ExtraCmd.Flags().StringVar(&nSumL, "nsuml", "", "目标成本计算列")
	ExtraCmd.Flags().StringVar(&sl, "sl", "", "查找的成本计算列")
	ExtraCmd.Flags().StringVar(&sumL, "suml", "", "查找的总成本计算列")
	ExtraCmd.Flags().StringVar(&accountId, "accountId", "", "客户id")
	RootCmd.AddCommand(ExtraCmd)
}

func extraCmd(cmd *cobra.Command, args []string)  {
	tools.ReadFileData(baseFile, 0)
	tools.Additional(newFile, l, nSumL, sl, sumL, accountId)
}
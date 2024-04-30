package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	profitMargin bool
	timeLocation string
	SumifsCmd = &cobra.Command{
		Use: "sumifs [flags] [value]",
		Long: `
 ______     __  __     __    __     __     ______   ______      
/\  ___\   /\ \/\ \   /\ "-./  \   /\ \   /\  ___\ /\  ___\     
\ \___  \  \ \ \_\ \  \ \ \-./\ \  \ \ \  \ \  __\ \ \___  \    
 \/\_____\  \ \_____\  \ \_\ \ \_\  \ \_\  \ \_\    \/\_____\   
  \/_____/   \/_____/   \/_/  \/_/   \/_/   \/_/     \/_____/`,
  Example: `main sumifs --rl1 N --rl2 A --suml CM --wl H --baseFile /basefile.csv --nf /nf.xlsx`,
		Run: sumifsCmd,
	}
)
func init() {
	SumifsCmd.Flags().StringVar(&baseFile, "basePath", "", "源文件地址")
	SumifsCmd.Flags().StringVar(&newFile, "nf", "", "目标文件地址")
	SumifsCmd.Flags().StringVar(&sumL, "suml", "", "求和列")
	SumifsCmd.Flags().StringVar(&sumRl1, "rl1", "", "条件范围1")
	SumifsCmd.Flags().StringVar(&sumRl2, "rl2", "", "条件范围2")
	SumifsCmd.Flags().StringVar(&sumWl, "wl", "", "填入的列")
	SumifsCmd.Flags().StringVar(&accountId, "accountId", "", "客户")
	SumifsCmd.Flags().BoolVar(&profitMargin, "pm", false, "利润率")
	SumifsCmd.Flags().StringVar(&timeLocation, "tl", "", "利润时间列")
	RootCmd.AddCommand(SumifsCmd)
}

func sumifsCmd(cmd *cobra.Command, args []string)  {
	tools.ReadFileData(baseFile, 0)
	tools.Sumifs(newFile, sumL, sumRl1, sumRl2, sumWl, accountId, profitMargin, timeLocation)
}

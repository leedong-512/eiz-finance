package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	group    bool
	groupL	string
	groupFile string
	chart bool
	profit bool
	ewe bool
	auxPosCmd = &cobra.Command{
		Use: "aupost [flags] [value]",
		Long: `
 ______     __  __     ______   ______     ______     ______  
/\  __ \   /\ \/\ \   /\  == \ /\  __ \   /\  ___\   /\__  _\ 
\ \  __ \  \ \ \_\ \  \ \  _-/ \ \ \/\ \  \ \___  \  \/_/\ \/ 
 \ \_\ \_\  \ \_____\  \ \_\    \ \_____\  \/\_____\    \ \_\ 
  \/_/\/_/   \/_____/   \/_/     \/_____/   \/_____/     \/_/ `,
  Example: `main auxpos --basePath /basefile.xlsx #获取Account数据
or
main auxpos --gl E --g true --gf /gf.xlsx #对一些列的数据进行分组并填入到指定位置`,
		Run: auxPosRun,
	}
)
func init() {
	auxPosCmd.Flags().StringVar(&baseFile, "basePath", "", "要解析的csv文件")
	auxPosCmd.Flags().BoolVar(&group,  "g", false, "是否分组")
	auxPosCmd.Flags().StringVar(&groupFile,  "gf", "", "分组文件")
	auxPosCmd.Flags().StringVar(&groupL,  "gl", "", "写入的列")
	auxPosCmd.Flags().BoolVar(&chart,  "chart", false, "是否增加图表")
	auxPosCmd.Flags().BoolVar(&profit,  "profit", false, "是否增加利润列")
	auxPosCmd.Flags().BoolVar(&ewe,  "ewe", false, "快递类型")
	RootCmd.AddCommand(auxPosCmd)
}

func auxPosRun(cmd *cobra.Command, args []string)  {
	if baseFile != "" {
		tools.GetAccount(baseFile, 0, "aupost")
	}
	if group {
		tools.Group(groupFile, groupL, chart, profit)
	}
}

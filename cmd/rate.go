package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	pt string
	RateCmd = &cobra.Command{
		Use: "rate [flags] [value]",
		Long: `
 ______     ______     ______   ______    
/\  == \   /\  __ \   /\__  _\ /\  ___\   
\ \  __<   \ \  __ \  \/_/\ \/ \ \  __\   
 \ \_\ \_\  \ \_\ \_\    \ \_\  \ \_____\ 
  \/_/ /_/   \/_/\/_/     \/_/   \/_____/`,
		Example: `main rate --nf ./files/kk/01/hesuan0101.xlsx --pt sf`,
		Run: rateCmd,
	}
)
func init() {
	RateCmd.Flags().StringVar(&pt, "pt", "", "plugin type")
	RateCmd.Flags().StringVar(&newFile, "nf", "", "文件地址")
	RootCmd.AddCommand(RateCmd)
}

func rateCmd(cmd *cobra.Command, args []string)  {
	//tools.ReadFileData(newFile)
	fileHandle , rows := tools.OpenXlsx(newFile)
	tools.Rate(fileHandle, rows, pt)
}


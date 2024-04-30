package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	fl string
	fv string
	file string
	FindCmd = &cobra.Command{
		Use: "find [flags] [value]",
		Long: `
 ______   __     __   __     _____    
/\  ___\ /\ \   /\ "-.\ \   /\  __-.  
\ \  __\ \ \ \  \ \ \-.  \  \ \ \/\ \ 
 \ \_\    \ \_\  \ \_\\"\_\  \ \____- 
  \/_/     \/_/   \/_/ \/_/   \/____/`,
		Example: `main find --dir ./file --fl A`,
		Run: findCmd,
	}
)
func init() {
	FindCmd.Flags().StringVar(&baseFile, "dir", "", "查找的路径")
	FindCmd.Flags().StringVar(&fl, "fl", "", "查找列")
	FindCmd.Flags().StringVar(&fv, "fv", "", "查找值")

	FindCmd.Flags().StringVar(&file, "file", "", "查找的文件")
	FindCmd.Flags().StringVar(&ml, "ml", "", "主列")
	FindCmd.Flags().StringVar(&accountId, "accountId", "", "客户id")
	RootCmd.AddCommand(FindCmd)
}

func findCmd(cmd *cobra.Command, args []string)  {
	if file != "" {
		tools.FindToSheets(file, ml, fl, accountId)
	} else {
		tools.Find(baseFile, fv, fl)
	}

}
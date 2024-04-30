package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	GetCmd = &cobra.Command{
		Use: "get [flags] [value]",
		Long: `
 ______     ______     ______  
/\  ___\   /\  ___\   /\__  _\ 
\ \ \__ \  \ \  __\   \/_/\ \/ 
 \ \_____\  \ \_____\    \ \_\ 
  \/_____/   \/_____/     \/_/ `,
		Run: getCmd,
	}
)
func init() {
	GetCmd.Flags().StringVar(&baseFile, "f", "", "查询文件")
	GetCmd.Flags().StringVar(&accountId, "accountId", "", "账户ID")
	RootCmd.AddCommand(GetCmd)
}

func getCmd(cmd *cobra.Command, args []string)  {
	tools.GetCharge(baseFile, accountId)
}

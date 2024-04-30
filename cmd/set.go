package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	setFile string
	setKey	string
	setValue string
	SetCmd = &cobra.Command{
		Use: "set [flags] [value]",
		Long: `
 ______     ______     ______  
/\  ___\   /\  ___\   /\__  _\ 
\ \___  \  \ \  __\   \/_/\ \/ 
 \/\_____\  \ \_____\    \ \_\ 
  \/_____/   \/_____/     \/_/ `,
  Example: `main set --f /file.xlsx --sk G --sv 19/02/2024 --accountId 4144`,
		Run: setCmd,
	}
)
func init() {
	SetCmd.Flags().StringVar(&setFile, "f", "", "目标文件")
	SetCmd.Flags().StringVar(&setKey, "sk", "", "赋值列")
	SetCmd.Flags().StringVar(&setValue, "sv", "", "值")
	SetCmd.Flags().StringVar(&accountId, "accountId", "", "账户ID")
	RootCmd.AddCommand(SetCmd)
}

func setCmd(cmd *cobra.Command, args []string)  {
	tools.SetValue(setFile, setKey, setValue, accountId)
}

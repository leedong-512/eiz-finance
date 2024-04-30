package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	SetHeaderCmd = &cobra.Command{
		Use: "head [flags] [value]",
		Long: `
 __  __     ______     ______     _____                         
/\ \_\ \   /\  __ \   /\  ___\   /\  __-.                       
\ \  __ \  \ \  __ \  \ \  __\   \ \ \/\ \                      
 \ \_\ \_\  \ \_\ \_\  \ \_____\  \ \____-                      
  \/_/\/_/   \/_/\/_/   \/_____/   \/____/ `,
  Example: `main head --nf /file.xlsx`,
		Run: setHeaderCmd,
	}
)
func init() {
	SetHeaderCmd.Flags().StringVar(&newFile, "nf", "", "新文件地址")
	RootCmd.AddCommand(SetHeaderCmd)
}

func setHeaderCmd(cmd *cobra.Command, args []string)  {
	tools.SetHeader(newFile)
}
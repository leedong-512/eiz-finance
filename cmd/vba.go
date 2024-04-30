package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	VbaCmd = &cobra.Command{
		Use: "vba [flags] [value]",
		Long: `
 __   __   ______     ______    
/\ \ / /  /\  == \   /\  __ \   
\ \ \'/   \ \  __<   \ \  __ \  
 \ \__|    \ \_____\  \ \_\ \_\ 
  \/_/      \/_____/   \/_/\/_/`,
		Example: `main vba --baseFile /filebase.xlsx`,
		Run: vbaCmd,
	}
)
func init() {
	VbaCmd.Flags().StringVar(&baseFile, "basePath", "", "源文件地址")
	RootCmd.AddCommand(VbaCmd)
}

func vbaCmd(cmd *cobra.Command, args []string)  {
	tools.Vba(baseFile)
}

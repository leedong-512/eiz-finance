package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	CopyCmd = &cobra.Command{
		Use: "copy [flags] [value]",
		Long: `
 ______     ______     ______   __  __    
/\  ___\   /\  __ \   /\  == \ /\ \_\ \   
\ \ \____  \ \ \/\ \  \ \  _-/ \ \____ \  
 \ \_____\  \ \_____\  \ \_\    \/\_____\ 
  \/_____/   \/_____/   \/_/     \/_____/`,
  Example: `main copy --l N,A --baseFile /filebase.csv --nf /nf.xlsx`,
		Run: copyCmd,
	}
)
func init() {
	CopyCmd.Flags().StringVar(&baseFile, "basePath", "", "源文件地址")
	CopyCmd.Flags().StringVar(&newFile, "nf", "", "目标文件地址")
	CopyCmd.Flags().StringVar(&l, "l", "", "复制列的定位")
	RootCmd.AddCommand(CopyCmd)
}

func copyCmd(cmd *cobra.Command, args []string)  {
	tools.ReadFileData(baseFile, 0)
	tools.Copy(newFile, l, false)
}

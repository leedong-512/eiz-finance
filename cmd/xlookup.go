package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	LookUpCmd = &cobra.Command{
		Use: "lookup [flags] [value]",
		Long: `
 __         ______     ______     __  __     __  __     ______  
/\ \       /\  __ \   /\  __ \   /\ \/ /    /\ \/\ \   /\  == \ 
\ \ \____  \ \ \/\ \  \ \ \/\ \  \ \  _"-.  \ \ \_\ \  \ \  _-/ 
 \ \_____\  \ \_____\  \ \_____\  \ \_\ \_\  \ \_____\  \ \_\   
  \/_____/   \/_____/   \/_____/   \/_/\/_/   \/_____/   \/_/  `,
  Example: `main lookup --ml B --sl M,N --baseFile /bashfile.csv --nf /nf.xlsx`,
		Run: xLookUpCmd,
	}
)
func init() {
	LookUpCmd.Flags().StringVar(&baseFile, "basePath", "", "要解析的excel文件")
	LookUpCmd.Flags().StringVar(&newFile, "nf", "", "新文件地址")
	LookUpCmd.Flags().StringVar(&sl, "sl", "M,N", "源文件的查找列")
	LookUpCmd.Flags().StringVar(&ml, "ml", "B", "源文件的查找列")
	RootCmd.AddCommand(LookUpCmd)
}

func xLookUpCmd(cmd *cobra.Command, args []string)  {
	tools.ReadFileData(baseFile, 0)
	tools.XLookUp(newFile, sl, ml)
}

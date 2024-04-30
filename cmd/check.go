package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	isAddSheet bool
	CheckCmd = &cobra.Command{
		Use: "check [flags] [value]",
		Long: ``,
		Example: `main check --l N,A --baseFile /filebase.csv --nf /nf.xlsx`,
		Run: checkCmd,
	}
)
func init() {
	CheckCmd.Flags().StringVar(&baseFile, "nf", "", "源文件地址")
	CheckCmd.Flags().StringVar(&dir, "dir", "", "目标文件地址")
	//CheckCmd.Flags().BoolVar(&isAddSheet, "add", false, "检出以逗号隔开的单号数据")
	//CheckCmd.Flags().StringVar(&l, "l", "", "复制列的定位")
	RootCmd.AddCommand(CheckCmd)
}

func checkCmd(cmd *cobra.Command, args []string)  {
	handle, rows := tools.OpenXlsx(baseFile)
	connoteValues := tools.ReadConnote(handle, rows, baseFile, false)
	tools.Check(dir, connoteValues)
}

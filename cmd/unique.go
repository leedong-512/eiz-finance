package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var(
	w string
	//aFile bool
	UniqueCmd = &cobra.Command{
		Use: "unique [flags] [value]",
		Long: ``,
		Example: `main unique --u D,E --w A --basePath ./files/20240103.xlsx --nf ./files/ewe103.xlsx`,
		Run: uniqueCmd,
	}
)
func init() {
	UniqueCmd.Flags().StringVar(&baseFile, "basePath", "", "源文件地址")
	UniqueCmd.Flags().StringVar(&newFile, "nf", "", "目标文件")
	UniqueCmd.Flags().StringVar(&l, "u", "", "去重列")
	UniqueCmd.Flags().StringVar(&w, "w", "", "写入列")
	//UniqueCmd.Flags().BoolVar(&aFile, "a", false, "是否为单文件")
	RootCmd.AddCommand(UniqueCmd)
}

func uniqueCmd(cmd *cobra.Command, args []string)  {
	tools.ReadFileData(baseFile, 0)
	tools.Unique(newFile, l, w)
}

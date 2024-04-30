package cmd

import (
	config "financeSys/configs"
	"fmt"
	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"
)

var (
	baseFile		string
	newFile 		string
	l 				string //用于copy
	sl 				string //用于检索映射
	ml 				string //用于目标文件 写入列
	sumWl			string //sumifs 用于文件 写入列
	sumRl1			string //sumifs 用于目标文件 条件列1
	sumRl2			string //sumifs 用于目标文件 条件列2
	sumL			string //sumifs 求和列
	accountId			string //客户的id
	cfg 			*config.Config
	RootCmd  = &cobra.Command{
		Use:   "tools",
		Short: "LeeHub",
		Long: `
 ________   _____   ________  
|_   __  | |_   _| |  __   _| 
  | |_ \_|   | |   |_/  / /   
  |  _| _    | |      .'.' _  
 _| |__/ |  _| |_   _/ /__/ | 
|________| |_____| |________| `,
		Version: "v0.0.1",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("cobra demo program")
		},

	}
)
func Execute() error {
	return RootCmd .Execute()
}

func init() {
	RootCmd.PersistentFlags().StringVar(&baseFile, "baseFile", "", "csv文件路径")
}

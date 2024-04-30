package cmd

import (
	"financeSys/tools"
	"github.com/spf13/cobra"
)

var (
	pf string
	PdfCmd = &cobra.Command{
		Use: "pdf [flags] [value]",
		Long: ``,
		Example: ``,
		Run: pdfCmd,
	}
)

func init() {
	PdfCmd.Flags().StringVar(&pf, "pf", "", "pdf 文件地址")
	RootCmd.AddCommand(PdfCmd)
}

func pdfCmd(cmd *cobra.Command, args []string)  {
	tools.ParsingPDF(pf)
}
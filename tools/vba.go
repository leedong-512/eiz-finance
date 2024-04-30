package tools

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func Vba(file string)  {
	// 初始化COM库
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	// 创建Excel应用程序实例
	excel, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		panic(err)
	}
	excelApp, err := excel.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		panic(err)
	}
	defer excelApp.Release()

	// 使Excel可见（可选）
	oleutil.PutProperty(excelApp, "Visible", true)

	// 添加新工作簿
	workbooks := oleutil.MustGetProperty(excelApp, "Workbooks").ToIDispatch()
	workbook := oleutil.MustCallMethod(workbooks, "Add").ToIDispatch()

	// 获取VBA项目
	vbProject := oleutil.MustGetProperty(workbook, "VBProject").ToIDispatch()

	// 添加VBA宏
	// 这里的宏代码仅作为示例，实际情况可能需要不同的宏代码
	vbaCode := `
Sub HelloWorld()
    MsgBox "Hello, world!"
End Sub
`
	// 添加一个新的模块，并设置其代码
	modules := oleutil.MustGetProperty(vbProject, "VBComponents").ToIDispatch()
	newModule := oleutil.MustCallMethod(modules, "Add", 1 /* vbext_ct_StdModule */).ToIDispatch()
	oleutil.MustPutProperty(newModule, "CodeModule", vbaCode)

	// 保存工作簿
	oleutil.MustCallMethod(workbook, "SaveAs", "./workbook.xlsm") // 请使用xlsm格式保存带有宏的工作簿

	// 关闭Excel应用程序
	oleutil.MustCallMethod(workbook, "Close", false)
	oleutil.MustCallMethod(excelApp, "Quit")
}

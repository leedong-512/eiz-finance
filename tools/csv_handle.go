package tools

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func MergeCsv(file1 string, file2 string)  {
	// 以追加模式打开文件1
	outFile, err := os.OpenFile(file1, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// 打开文件2
	inFile, err := os.Open(file2)
	if err != nil {
		fmt.Println("Error opening the file:", file2, err)
		return
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)

	// 跳过文件2的列头
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error reading the header from file2:", err)
		return
	}

	// 读取并追加数据
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // 文件末尾
		}
		if err != nil {
			fmt.Println("Error reading a record:", err)
			return
		}
		err = writer.Write(record)
		if err != nil {
			fmt.Println("Error writing a record:", err)
			return
		}
	}
}

package tools

/*func init()  {
	err := license.SetMeteredKey(`64e4954e707c26a57235c5d54c1b2833c5e49d75ee518d6caa3bec3ecbd8fd7d`)
	if err != nil {
		fmt.Printf("ERROR: Failed to set metered key: %v\n", err)
		fmt.Printf("Make sure to get a valid key from https://cloud.unidoc.io\n")
		panic(err)
	}
}*/
/*func ParsingPDF(pdfFile string) {
	// 设置Unipdf的许可证密钥。
	//license.SetLicenseKey("64e4954e707c26a57235c5d54c1b2833c5e49d75ee518d6caa3bec3ecbd8fd7d", "leedong")
	lk := license.GetLicenseKey()
	if lk == nil {
		fmt.Printf("Failed retrieving license key")
		return
	}
	state, err := license.GetMeteredState()
	if err != nil {
		fmt.Printf("ERROR getting metered state: %+v\n", err)
		panic(err)
	}
	fmt.Printf("Metered state: %+v\n", state)
	file, _ := os.Open(pdfFile) // 打开PDF文件)
	defer file.Close()
	// 创建PDF解析器。
	pdfReader, err := pdf.NewPdfReaderLazy(file)
	if err != nil {
		fmt.Println("无法创建PDF解析器:", err)
		return
	}
	numPages , _ := pdfReader.GetNumPages()
	// 创建CSV文件。
	csvFile, err := os.Create("./files/pdf/output.csv")
	if err != nil {
		fmt.Println("无法创建CSV文件:", err)
		return
	}
	defer csvFile.Close()
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		//pdfReader.Ge
		if pageNum != 3 {
			continue
		}
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			fmt.Printf("无法提取第 %d 页: %v\n", pageNum, err)
			continue
		}
		mbox, err := page.GetMediaBox()
		if err != nil {
			return
		}
		if page.Rotate != nil && *page.Rotate == 90 {
			// TODO: This is a "hack" to change the perspective of the extractor to account for the rotation.
			contents, err := page.GetContentStreams()
			if err != nil {
				fmt.Println("无法提取内容流:", err)
				return
			}

			cc := contentstream.NewContentCreator()
			cc.Translate(mbox.Width()/2, mbox.Height()/2)
			cc.RotateDeg(-90)
			cc.Translate(-mbox.Width()/2, -mbox.Height()/2)
			rotateOps := cc.Operations().String()
			contents = append([]string{rotateOps}, contents...)

			page.Duplicate()
			err = page.SetContentStreams(contents, core.NewRawEncoder())
			if err != nil {
				fmt.Println("无法设置内容流:", err)
				return

			}
			page.Rotate = nil
		}
		// 创建文本提取器。
		textExtractor, err := extractor.New(page)
		if err != nil {
			fmt.Printf("无法创建文本提取器: %v\n", err)
			continue
		}
		// 提取页面文本。
		pageText, _, _, err := textExtractor.ExtractPageText()
		if err != nil {
			fmt.Printf("无法提取第 %d 页文本: %v\n", pageNum, err)
			continue
		}
		text := pageText.Text()
		textMarks := pageText.Marks()
		// 将文本写入CSV文件。
		_, err = csvFile.WriteString(strings.TrimSpace(text) + "\n")
		if err != nil {
			fmt.Printf("无法写入CSV文件: %v\n", err)
			continue
		}
	}
}*/




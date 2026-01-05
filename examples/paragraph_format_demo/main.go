package main

import (
	"log"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func main() {
	// 创建新文档
	doc := document.New()

	// 示例1：使用单独的方法设置段落属性
	title := doc.AddParagraph("第一章：段落格式自定义功能演示")
	title.SetAlignment(document.AlignCenter)
	title.SetStyle("Heading1")
	title.SetKeepWithNext(true)     // 与下一段保持在同一页
	title.SetPageBreakBefore(true)  // 从新页开始
	title.SetOutlineLevel(0)        // 设置为一级大纲

	// 示例2：使用SetParagraphFormat一次性设置多个属性
	intro := doc.AddParagraph("本章将介绍WordZero库的段落格式自定义功能，包括分页控制、行控制、大纲级别等高级特性。")
	intro.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:    document.AlignJustify,
		LineSpacing:  1.5,
		BeforePara:   12,
		AfterPara:    12,
		FirstLineCm:  0.5,
		WidowControl: true, // 启用孤行控制
	})

	// 示例3：段落保持在一起的演示
	subtitle1 := doc.AddParagraph("1.1 段落分页控制")
	subtitle1.SetParagraphFormat(&document.ParagraphFormatConfig{
		Style:        "Heading2",
		Alignment:    document.AlignLeft,
		BeforePara:   18,
		AfterPara:    6,
		KeepWithNext: true, // 确保标题和下一段不被分页
		OutlineLevel: 1,
	})

	content1 := doc.AddParagraph("段落分页控制功能允许您精确控制段落在文档中的分页行为。" +
		"通过SetKeepWithNext方法，可以确保标题和其后的内容段落始终保持在同一页面上，" +
		"避免标题单独出现在页面底部的情况。")
	content1.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:    document.AlignJustify,
		LineSpacing:  1.5,
		FirstLineCm:  0.5,
		WidowControl: true,
	})

	// 示例4：保持段落行在一起
	subtitle2 := doc.AddParagraph("1.2 段落行保持功能")
	subtitle2.SetParagraphFormat(&document.ParagraphFormatConfig{
		Style:        "Heading2",
		Alignment:    document.AlignLeft,
		BeforePara:   18,
		AfterPara:    6,
		KeepWithNext: true,
		OutlineLevel: 1,
	})

	content2 := doc.AddParagraph("使用SetKeepLines方法可以确保段落的所有行都保持在同一页面上，" +
		"防止段落被分割到不同页面。这对于需要保持完整性的段落（如重要说明、引用等）非常有用。")
	content2.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:   document.AlignJustify,
		LineSpacing: 1.5,
		FirstLineCm: 0.5,
		KeepLines:   true, // 段落所有行保持在一起
	})

	// 示例5：章节开始新页
	chapter2 := doc.AddParagraph("第二章：高级格式特性")
	chapter2.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:       document.AlignCenter,
		Style:           "Heading1",
		BeforePara:      24,
		AfterPara:       12,
		PageBreakBefore: true, // 新章节从新页开始
		KeepWithNext:    true,
		OutlineLevel:    0,
	})

	intro2 := doc.AddParagraph("本章介绍更多高级段落格式特性的使用方法。")
	intro2.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:    document.AlignJustify,
		LineSpacing:  1.5,
		FirstLineCm:  0.5,
		BeforePara:   6,
		AfterPara:    6,
		WidowControl: true,
	})

	// 示例6：带缩进的特殊段落
	subtitle3 := doc.AddParagraph("2.1 复杂缩进示例")
	subtitle3.SetParagraphFormat(&document.ParagraphFormatConfig{
		Style:        "Heading2",
		Alignment:    document.AlignLeft,
		BeforePara:   18,
		AfterPara:    6,
		KeepWithNext: true,
		OutlineLevel: 1,
	})

	// 悬挂缩进示例
	hangingPara := doc.AddParagraph("● 这是一个使用悬挂缩进的列表项。" +
		"悬挂缩进是指首行向左突出，而后续行向右缩进的格式。" +
		"这种格式常用于编号列表、项目符号列表等场景。")
	hangingPara.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:   document.AlignLeft,
		FirstLineCm: -0.5, // 悬挂缩进
		LeftCm:      1.0,  // 左缩进
		LineSpacing: 1.2,
	})

	// 示例7：引用块样式
	quote := doc.AddParagraph("\"优秀的排版不仅仅是让文字看起来漂亮，" +
		"更重要的是让读者能够轻松理解内容的结构和层次。\"")
	quote.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:   document.AlignCenter,
		LeftCm:      1.5,
		RightCm:     1.5,
		BeforePara:  12,
		AfterPara:   12,
		LineSpacing: 1.5,
		KeepLines:   true, // 引用保持完整
	})

	// 示例8：大纲级别演示
	section := doc.AddParagraph("2.2 大纲级别的使用")
	section.SetParagraphFormat(&document.ParagraphFormatConfig{
		Style:        "Heading2",
		Alignment:    document.AlignLeft,
		BeforePara:   18,
		AfterPara:    6,
		KeepWithNext: true,
		OutlineLevel: 1, // 二级标题的大纲级别
	})

	subsection := doc.AddParagraph("2.2.1 三级标题示例")
	subsection.SetParagraphFormat(&document.ParagraphFormatConfig{
		Style:        "Heading3",
		Alignment:    document.AlignLeft,
		BeforePara:   12,
		AfterPara:    6,
		KeepWithNext: true,
		OutlineLevel: 2, // 三级标题的大纲级别
	})

	finalContent := doc.AddParagraph("大纲级别设置可以让Word在导航窗格中显示文档结构，" +
		"方便读者快速浏览和定位内容。级别范围从0到8，分别对应标题1到标题9，" +
		"而正文段落通常不设置大纲级别。")
	finalContent.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:    document.AlignJustify,
		LineSpacing:  1.5,
		FirstLineCm:  0.5,
		BeforePara:   6,
		AfterPara:    6,
		WidowControl: true,
	})

	// 示例9：总结段落
	summary := doc.AddParagraph("总结")
	summary.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:       document.AlignCenter,
		Style:           "Heading1",
		BeforePara:      24,
		PageBreakBefore: true,
		OutlineLevel:    0,
	})

	summaryContent := doc.AddParagraph("通过本文档的演示，您已经了解了WordZero库提供的全面段落格式自定义功能，" +
		"包括对齐方式、间距设置、缩进控制、分页控制、行控制、孤行控制以及大纲级别设置等。" +
		"这些功能可以帮助您创建专业、美观且结构清晰的Word文档。")
	summaryContent.SetParagraphFormat(&document.ParagraphFormatConfig{
		Alignment:    document.AlignJustify,
		LineSpacing:  1.5,
		FirstLineCm:  0.5,
		BeforePara:   12,
		AfterPara:    12,
		WidowControl: true,
		KeepLines:    true,
	})

	// 保存文档
	filename := "examples/output/paragraph_format_demo.docx"
	err := doc.Save(filename)
	if err != nil {
		log.Fatalf("保存文档失败: %v", err)
	}

	log.Printf("段落格式演示文档创建成功！文件保存在: %s", filename)
	log.Println("该文档演示了以下功能：")
	log.Println("  - SetKeepWithNext: 标题与下一段保持在同一页")
	log.Println("  - SetKeepLines: 段落行保持在一起不被分页")
	log.Println("  - SetPageBreakBefore: 章节从新页开始")
	log.Println("  - SetWidowControl: 孤行控制提升排版质量")
	log.Println("  - SetOutlineLevel: 大纲级别设置便于文档导航")
	log.Println("  - SetParagraphFormat: 一次性设置多个格式属性")
}

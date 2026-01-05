// Package main 页面设置功能示例
package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func main() {
	fmt.Println("=== WordZero 页面设置功能演示 ===")

	// 输出目录
	outputDir := "examples/output"

	// 演示1：A4纵向文档
	fmt.Println("\n1. 创建A4纵向文档")
	createA4PortraitDoc(outputDir)

	// 演示2：A4横向文档
	fmt.Println("\n2. 创建A4横向文档")
	createA4LandscapeDoc(outputDir)

	// 演示3：Letter纵向文档
	fmt.Println("\n3. 创建Letter纵向文档")
	createLetterPortraitDoc(outputDir)

	// 演示4：Legal纵向文档
	fmt.Println("\n4. 创建Legal纵向文档")
	createLegalPortraitDoc(outputDir)

	// 演示5：A3纵向文档
	fmt.Println("\n5. 创建A3纵向文档")
	createA3PortraitDoc(outputDir)

	// 演示6：A5纵向文档
	fmt.Println("\n6. 创建A5纵向文档")
	createA5PortraitDoc(outputDir)

	// 演示7：自定义尺寸文档（正方形）
	fmt.Println("\n7. 创建自定义尺寸文档（正方形）")
	createCustomSquareDoc(outputDir)

	// 演示8：自定义尺寸文档（名片尺寸）
	fmt.Println("\n8. 创建自定义尺寸文档（名片尺寸）")
	createCustomBusinessCardDoc(outputDir)

	fmt.Println("\n页面设置演示完成！所有文档已保存到 examples/output/ 目录下")
}

// createA4PortraitDoc 创建A4纵向文档
func createA4PortraitDoc(outputDir string) {
	doc := document.New()

	// 设置A4纵向页面
	if err := doc.SetPageSize(document.PageSizeA4); err != nil {
		log.Printf("设置A4页面尺寸失败: %v", err)
		return
	}
	if err := doc.SetPageOrientation(document.OrientationPortrait); err != nil {
		log.Printf("设置纵向页面失败: %v", err)
		return
	}

	// 添加内容
	title := doc.AddParagraph("A4纵向页面文档")
	title.SetAlignment(document.AlignCenter)

	doc.AddParagraph("")
	settings := doc.GetPageSettings()
	doc.AddParagraph(fmt.Sprintf("页面尺寸: %s", settings.Size))
	doc.AddParagraph(fmt.Sprintf("页面方向: %s", settings.Orientation))
	doc.AddParagraph("这是标准的A4纵向页面，常用于办公文档、报告等。")
	doc.AddParagraph("尺寸：210mm x 297mm")

	// 添加更多内容以展示页面效果
	addSampleContent(doc, "A4纵向")

	// 保存文档
	filename := filepath.Join(outputDir, "page_A4_portrait.docx")
	if err := doc.Save(filename); err != nil {
		log.Printf("保存A4纵向文档失败: %v", err)
		return
	}

	printPageSettings(settings)
	fmt.Printf("文档已保存到: %s\n", filename)
}

// createA4LandscapeDoc 创建A4横向文档
func createA4LandscapeDoc(outputDir string) {
	doc := document.New()

	// 设置A4横向页面
	if err := doc.SetPageSize(document.PageSizeA4); err != nil {
		log.Printf("设置A4页面尺寸失败: %v", err)
		return
	}
	if err := doc.SetPageOrientation(document.OrientationLandscape); err != nil {
		log.Printf("设置横向页面失败: %v", err)
		return
	}

	// 添加内容
	title := doc.AddParagraph("A4横向页面文档")
	title.SetAlignment(document.AlignCenter)

	doc.AddParagraph("")
	settings := doc.GetPageSettings()
	doc.AddParagraph(fmt.Sprintf("页面尺寸: %s", settings.Size))
	doc.AddParagraph(fmt.Sprintf("页面方向: %s", settings.Orientation))
	doc.AddParagraph("这是A4横向页面，适合展示宽表格、图表等内容。")
	doc.AddParagraph("尺寸：297mm x 210mm")

	addSampleContent(doc, "A4横向")

	// 保存文档
	filename := filepath.Join(outputDir, "page_A4_landscape.docx")
	if err := doc.Save(filename); err != nil {
		log.Printf("保存A4横向文档失败: %v", err)
		return
	}

	printPageSettings(settings)
	fmt.Printf("文档已保存到: %s\n", filename)
}

// createLetterPortraitDoc 创建Letter纵向文档
func createLetterPortraitDoc(outputDir string) {
	doc := document.New()

	// 设置Letter纵向页面
	if err := doc.SetPageSize(document.PageSizeLetter); err != nil {
		log.Printf("设置Letter页面尺寸失败: %v", err)
		return
	}
	if err := doc.SetPageOrientation(document.OrientationPortrait); err != nil {
		log.Printf("设置纵向页面失败: %v", err)
		return
	}

	// 添加内容
	title := doc.AddParagraph("Letter纵向页面文档")
	title.SetAlignment(document.AlignCenter)

	doc.AddParagraph("")
	settings := doc.GetPageSettings()
	doc.AddParagraph(fmt.Sprintf("页面尺寸: %s", settings.Size))
	doc.AddParagraph(fmt.Sprintf("页面方向: %s", settings.Orientation))
	doc.AddParagraph("这是美国标准的Letter纵向页面，在北美地区广泛使用。")
	doc.AddParagraph("尺寸：8.5\" x 11\"（215.9mm x 279.4mm）")

	addSampleContent(doc, "Letter纵向")

	// 保存文档
	filename := filepath.Join(outputDir, "page_Letter_portrait.docx")
	if err := doc.Save(filename); err != nil {
		log.Printf("保存Letter纵向文档失败: %v", err)
		return
	}

	printPageSettings(settings)
	fmt.Printf("文档已保存到: %s\n", filename)
}

// createLegalPortraitDoc 创建Legal纵向文档
func createLegalPortraitDoc(outputDir string) {
	doc := document.New()

	// 设置Legal纵向页面
	if err := doc.SetPageSize(document.PageSizeLegal); err != nil {
		log.Printf("设置Legal页面尺寸失败: %v", err)
		return
	}
	if err := doc.SetPageOrientation(document.OrientationPortrait); err != nil {
		log.Printf("设置纵向页面失败: %v", err)
		return
	}

	// 添加内容
	title := doc.AddParagraph("Legal纵向页面文档")
	title.SetAlignment(document.AlignCenter)

	doc.AddParagraph("")
	settings := doc.GetPageSettings()
	doc.AddParagraph(fmt.Sprintf("页面尺寸: %s", settings.Size))
	doc.AddParagraph(fmt.Sprintf("页面方向: %s", settings.Orientation))
	doc.AddParagraph("这是Legal纵向页面，常用于法律文档、合同等。")
	doc.AddParagraph("尺寸：8.5\" x 14\"（215.9mm x 355.6mm）")

	addSampleContent(doc, "Legal纵向")

	// 保存文档
	filename := filepath.Join(outputDir, "page_Legal_portrait.docx")
	if err := doc.Save(filename); err != nil {
		log.Printf("保存Legal纵向文档失败: %v", err)
		return
	}

	printPageSettings(settings)
	fmt.Printf("文档已保存到: %s\n", filename)
}

// createA3PortraitDoc 创建A3纵向文档
func createA3PortraitDoc(outputDir string) {
	doc := document.New()

	// 设置A3纵向页面
	if err := doc.SetPageSize(document.PageSizeA3); err != nil {
		log.Printf("设置A3页面尺寸失败: %v", err)
		return
	}
	if err := doc.SetPageOrientation(document.OrientationPortrait); err != nil {
		log.Printf("设置纵向页面失败: %v", err)
		return
	}

	// 添加内容
	title := doc.AddParagraph("A3纵向页面文档")
	title.SetAlignment(document.AlignCenter)

	doc.AddParagraph("")
	settings := doc.GetPageSettings()
	doc.AddParagraph(fmt.Sprintf("页面尺寸: %s", settings.Size))
	doc.AddParagraph(fmt.Sprintf("页面方向: %s", settings.Orientation))
	doc.AddParagraph("这是A3纵向页面，适合打印大尺寸图表、海报等。")
	doc.AddParagraph("尺寸：297mm x 420mm")

	addSampleContent(doc, "A3纵向")

	// 保存文档
	filename := filepath.Join(outputDir, "page_A3_portrait.docx")
	if err := doc.Save(filename); err != nil {
		log.Printf("保存A3纵向文档失败: %v", err)
		return
	}

	printPageSettings(settings)
	fmt.Printf("文档已保存到: %s\n", filename)
}

// createA5PortraitDoc 创建A5纵向文档
func createA5PortraitDoc(outputDir string) {
	doc := document.New()

	// 设置A5纵向页面
	if err := doc.SetPageSize(document.PageSizeA5); err != nil {
		log.Printf("设置A5页面尺寸失败: %v", err)
		return
	}
	if err := doc.SetPageOrientation(document.OrientationPortrait); err != nil {
		log.Printf("设置纵向页面失败: %v", err)
		return
	}

	// 添加内容
	title := doc.AddParagraph("A5纵向页面文档")
	title.SetAlignment(document.AlignCenter)

	doc.AddParagraph("")
	settings := doc.GetPageSettings()
	doc.AddParagraph(fmt.Sprintf("页面尺寸: %s", settings.Size))
	doc.AddParagraph(fmt.Sprintf("页面方向: %s", settings.Orientation))
	doc.AddParagraph("这是A5纵向页面，适合小册子、笔记本等。")
	doc.AddParagraph("尺寸：148mm x 210mm")

	addSampleContent(doc, "A5纵向")

	// 保存文档
	filename := filepath.Join(outputDir, "page_A5_portrait.docx")
	if err := doc.Save(filename); err != nil {
		log.Printf("保存A5纵向文档失败: %v", err)
		return
	}

	printPageSettings(settings)
	fmt.Printf("文档已保存到: %s\n", filename)
}

// createCustomSquareDoc 创建自定义正方形文档
func createCustomSquareDoc(outputDir string) {
	doc := document.New()

	// 设置自定义正方形页面
	if err := doc.SetCustomPageSize(200, 200); err != nil {
		log.Printf("设置自定义页面尺寸失败: %v", err)
		return
	}

	// 添加内容
	title := doc.AddParagraph("自定义正方形页面文档")
	title.SetAlignment(document.AlignCenter)

	doc.AddParagraph("")
	settings := doc.GetPageSettings()
	doc.AddParagraph(fmt.Sprintf("页面尺寸: %s", settings.Size))
	doc.AddParagraph(fmt.Sprintf("自定义尺寸: %.1fmm x %.1fmm", settings.CustomWidth, settings.CustomHeight))
	doc.AddParagraph("这是自定义的正方形页面，适合特殊设计需求。")

	addSampleContent(doc, "自定义正方形")

	// 保存文档
	filename := filepath.Join(outputDir, "page_Custom_square.docx")
	if err := doc.Save(filename); err != nil {
		log.Printf("保存自定义正方形文档失败: %v", err)
		return
	}

	printPageSettings(settings)
	fmt.Printf("文档已保存到: %s\n", filename)
}

// createCustomBusinessCardDoc 创建自定义名片尺寸文档
func createCustomBusinessCardDoc(outputDir string) {
	doc := document.New()

	// 设置名片尺寸（90mm x 54mm）
	if err := doc.SetCustomPageSize(90, 54); err != nil {
		log.Printf("设置自定义页面尺寸失败: %v", err)
		return
	}

	// 设置较小的边距
	if err := doc.SetPageMargins(5, 5, 5, 5); err != nil {
		log.Printf("设置页面边距失败: %v", err)
		return
	}

	// 添加内容
	title := doc.AddParagraph("名片尺寸文档")
	title.SetAlignment(document.AlignCenter)

	doc.AddParagraph("")
	settings := doc.GetPageSettings()
	doc.AddParagraph(fmt.Sprintf("页面尺寸: %s", settings.Size))
	doc.AddParagraph(fmt.Sprintf("自定义尺寸: %.1fmm x %.1fmm", settings.CustomWidth, settings.CustomHeight))
	doc.AddParagraph("标准名片尺寸，适合设计名片、标签等小型印刷品。")

	// 保存文档
	filename := filepath.Join(outputDir, "page_Custom_businesscard.docx")
	if err := doc.Save(filename); err != nil {
		log.Printf("保存自定义名片文档失败: %v", err)
		return
	}

	printPageSettings(settings)
	fmt.Printf("文档已保存到: %s\n", filename)
}

// addSampleContent 添加示例内容
func addSampleContent(doc *document.Document, pageType string) {
	doc.AddParagraph("")
	doc.AddParagraph("示例内容：")
	doc.AddParagraph("这个" + pageType + "页面演示了WordZero库的页面设置功能。")
	doc.AddParagraph("您可以使用此功能创建各种不同尺寸和方向的文档。")
	doc.AddParagraph("")
	doc.AddParagraph("支持的页面尺寸包括：")
	doc.AddParagraph("• A4 (210mm x 297mm)")
	doc.AddParagraph("• Letter (8.5\" x 11\")")
	doc.AddParagraph("• Legal (8.5\" x 14\")")
	doc.AddParagraph("• A3 (297mm x 420mm)")
	doc.AddParagraph("• A5 (148mm x 210mm)")
	doc.AddParagraph("• 自定义尺寸")
	doc.AddParagraph("")
	doc.AddParagraph("每种页面尺寸都可以设置为纵向或横向方向。")
}

// printPageSettings 打印页面设置信息
func printPageSettings(settings *document.PageSettings) {
	fmt.Printf("  页面尺寸: %s\n", settings.Size)
	if settings.Size == document.PageSizeCustom {
		fmt.Printf("  自定义尺寸: %.1fmm x %.1fmm\n", settings.CustomWidth, settings.CustomHeight)
	}
	fmt.Printf("  页面方向: %s\n", settings.Orientation)
	fmt.Printf("  页面边距: 上%.1fmm 右%.1fmm 下%.1fmm 左%.1fmm\n",
		settings.MarginTop, settings.MarginRight, settings.MarginBottom, settings.MarginLeft)
	fmt.Println()
}

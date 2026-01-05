package test

import (
	"fmt"
	"math"
	"testing"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func TestDebugPageSettings(t *testing.T) {
	// 创建文档
	doc := document.New()

	// 设置页面配置
	settings := &document.PageSettings{
		Size:           document.PageSizeLetter,
		Orientation:    document.OrientationLandscape,
		MarginTop:      25,
		MarginRight:    20,
		MarginBottom:   30,
		MarginLeft:     25,
		HeaderDistance: 12,
		FooterDistance: 15,
		GutterWidth:    5,
	}

	err := doc.SetPageSettings(settings)
	if err != nil {
		t.Fatalf("设置页面配置失败: %v", err)
	}

	// 验证设置
	currentSettings := doc.GetPageSettings()
	fmt.Printf("设置后的页面配置:\n")
	fmt.Printf("  尺寸: %s\n", currentSettings.Size)
	fmt.Printf("  方向: %s\n", currentSettings.Orientation)

	// 添加测试内容
	doc.AddParagraph("测试页面设置保存和加载")

	// 保存文档
	testFile := "debug_page_settings.docx"
	err = doc.Save(testFile)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}

	fmt.Printf("文档已保存到: %s\n", testFile)

	// 重新打开文档
	loadedDoc, err := document.Open(testFile)
	if err != nil {
		t.Fatalf("重新打开文档失败: %v", err)
	}

	// 检查加载后文档的Body.Elements
	fmt.Printf("加载后文档的Body.Elements数量: %d\n", len(loadedDoc.Body.Elements))
	for i, element := range loadedDoc.Body.Elements {
		switch elem := element.(type) {
		case *document.SectionProperties:
			fmt.Printf("  元素%d: SectionProperties found!\n", i)
			if elem.PageSize != nil {
				fmt.Printf("    PageSize: w=%s, h=%s, orient=%s\n", elem.PageSize.W, elem.PageSize.H, elem.PageSize.Orient)
			} else {
				fmt.Printf("    PageSize: nil\n")
			}
		case *document.Paragraph:
			fmt.Printf("  元素%d: Paragraph\n", i)
		default:
			fmt.Printf("  元素%d: 其他类型 (%T)\n", i, element)
		}
	}

	// 验证加载后的设置
	loadedSettings := loadedDoc.GetPageSettings()
	fmt.Printf("加载后的页面配置:\n")
	fmt.Printf("  尺寸: %s\n", loadedSettings.Size)
	fmt.Printf("  方向: %s\n", loadedSettings.Orientation)

	// 验证设置是否正确
	if loadedSettings.Size != settings.Size {
		t.Errorf("加载后页面尺寸不匹配，期望: %s, 实际: %s", settings.Size, loadedSettings.Size)
	}

	if loadedSettings.Orientation != settings.Orientation {
		t.Errorf("加载后页面方向不匹配，期望: %s, 实际: %s", settings.Orientation, loadedSettings.Orientation)
	}

	// 详细调试页面尺寸解析过程
	parts := loadedDoc.GetParts()
	if docXML, exists := parts["word/document.xml"]; exists {
		fmt.Printf("document.xml内容前500字符:\n%s\n", string(docXML)[:min(500, len(docXML))])

		// 手动验证twips转换
		fmt.Printf("调试页面尺寸转换:\n")

		// Letter尺寸：215.9mm x 279.4mm
		// 横向后应该是：279.4mm x 215.9mm
		// 转换为twips：279.4 * 56.69 ≈ 15840，215.9 * 56.69 ≈ 12240

		width_twips := 15840.0
		height_twips := 12240.0
		width_mm := width_twips / 56.692913385827
		height_mm := height_twips / 56.692913385827

		fmt.Printf("  从XML读取: 宽度=%d twips, 高度=%d twips\n", int(width_twips), int(height_twips))
		fmt.Printf("  转换为毫米: 宽度=%.1fmm, 高度=%.1fmm\n", width_mm, height_mm)

		// 测试页面尺寸识别
		fmt.Printf("  Letter纵向尺寸: 215.9mm x 279.4mm\n")
		fmt.Printf("  Letter横向尺寸: 279.4mm x 215.9mm\n")
		fmt.Printf("  实际解析尺寸: %.1fmm x %.1fmm\n", width_mm, height_mm)

		// 检查容差
		tolerance := 1.0
		letter_width := 215.9
		letter_height := 279.4

		// 检查横向匹配
		landscape_match := (math.Abs(width_mm-letter_height) < tolerance && math.Abs(height_mm-letter_width) < tolerance)
		fmt.Printf("  横向Letter匹配: %t (容差=%.1fmm)\n", landscape_match, tolerance)

		// 检查纵向匹配
		portrait_match := (math.Abs(width_mm-letter_width) < tolerance && math.Abs(height_mm-letter_height) < tolerance)
		fmt.Printf("  纵向Letter匹配: %t (容差=%.1fmm)\n", portrait_match, tolerance)
	} else {
		fmt.Printf("未找到document.xml\n")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

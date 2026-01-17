// Package main provides comprehensive issue reproduction tests
// Uses manual docx extraction and third-party tools for validation
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func main() {
	fmt.Println("==============================================")
	fmt.Println("Comprehensive Issue Reproduction Tests")
	fmt.Println("==============================================")

	outputDir := "/home/user/wordZero/validation/output/issues_check"
	os.MkdirAll(outputDir, 0755)

	// Issue #82: 图片不支持在有文字的段落行内显示
	fmt.Println("\n--- Issue #82: Inline Image Display ---")
	testIssue82(outputDir)

	// Issue #81: OpenFromMemory文件损坏
	fmt.Println("\n--- Issue #81: OpenFromMemory File Corruption ---")
	testIssue81(outputDir)

	// Issue #78: 制表符丢失
	fmt.Println("\n--- Issue #78: Tab Character Loss ---")
	testIssue78(outputDir)

	// Issue #76: 页眉页脚渲染后消失
	fmt.Println("\n--- Issue #76: Header/Footer Lost After Rendering ---")
	testIssue76(outputDir)

	// Issue #91: 表格图片替换 (已修复，验证修复)
	fmt.Println("\n--- Issue #91: Table Image Replacement (verify fix) ---")
	testIssue91(outputDir)

	// Issue #88: each循环粗体 (已修复，验证修复)
	fmt.Println("\n--- Issue #88: Each Loop Bold (verify fix) ---")
	testIssue88(outputDir)

	fmt.Println("\n==============================================")
	fmt.Println("Test files generated. Use external tools to validate:")
	fmt.Println("1. unzip -l <file.docx>  # List contents")
	fmt.Println("2. unzip -p <file.docx> word/document.xml | xmllint --format -")
	fmt.Println("3. python3 -c \"from docx import Document; d=Document('<file>'); print('OK')\"")
	fmt.Println("==============================================")
}

// Issue #82: 图片不支持在有文字的段落行内显示
func testIssue82(outputDir string) {
	doc := document.New()

	// 创建带文字的段落，然后尝试添加图片
	para := doc.AddParagraph("This is text before image: ")

	// 尝试在同一段落中添加图片（这应该实现行内显示）
	// 根据issue，图片会作为block显示，导致换行

	doc.AddParagraph("Text after image placeholder")

	outputPath := filepath.Join(outputDir, "issue82_inline_image.docx")
	if err := doc.Save(outputPath); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	fmt.Printf("Created: %s\n", outputPath)
	fmt.Println("Issue #82: Need to check if images can be inline with text")
	_ = para // 使用变量避免编译警告
}

// Issue #81: OpenFromMemory文件损坏
func testIssue81(outputDir string) {
	// 创建原始文档
	doc1 := document.New()
	doc1.AddParagraph("Original content")
	doc1.AddParagraph("Second paragraph")

	// 保存到bytes
	docBytes, err := doc1.ToBytes()
	if err != nil {
		fmt.Printf("ERROR ToBytes: %v\n", err)
		return
	}
	fmt.Printf("Document bytes length: %d\n", len(docBytes))

	// 使用OpenFromMemory打开
	reader := strings.NewReader(string(docBytes))
	doc2, err := document.OpenFromMemory(nopCloser{reader})
	if err != nil {
		fmt.Printf("ERROR OpenFromMemory: %v\n", err)
		return
	}

	// 添加新内容
	doc2.AddParagraph("Content added after OpenFromMemory")

	// 保存
	outputPath := filepath.Join(outputDir, "issue81_from_memory.docx")
	if err := doc2.Save(outputPath); err != nil {
		fmt.Printf("ERROR Save: %v\n", err)
		return
	}
	fmt.Printf("Created: %s\n", outputPath)
	fmt.Println("Issue #81: Check if file opens correctly in MS Office")
}

// Issue #78: 制表符丢失
func testIssue78(outputDir string) {
	// 创建包含制表符的文档
	doc := document.New()
	doc.AddParagraph("Column1\tColumn2\tColumn3")
	doc.AddParagraph("Data1\tData2\tData3")
	doc.AddParagraph("No tabs here")
	doc.AddParagraph("More\ttabs\there")

	// 保存原始文档
	originalPath := filepath.Join(outputDir, "issue78_original.docx")
	if err := doc.Save(originalPath); err != nil {
		fmt.Printf("ERROR saving original: %v\n", err)
		return
	}
	fmt.Printf("Created original: %s\n", originalPath)

	// 重新打开并保存
	doc2, err := document.Open(originalPath)
	if err != nil {
		fmt.Printf("ERROR opening: %v\n", err)
		return
	}

	reopenedPath := filepath.Join(outputDir, "issue78_reopened.docx")
	if err := doc2.Save(reopenedPath); err != nil {
		fmt.Printf("ERROR saving reopened: %v\n", err)
		return
	}
	fmt.Printf("Created reopened: %s\n", reopenedPath)
	fmt.Println("Issue #78: Compare XML to check if tabs (\\t) are preserved")
}

// Issue #76: 页眉页脚渲染后消失
func testIssue76(outputDir string) {
	// 创建带页眉页脚的模板
	templateDoc := document.New()

	// 添加页眉
	if err := templateDoc.AddHeader(document.HeaderFooterTypeDefault, "Template Header - {{title}}"); err != nil {
		fmt.Printf("ERROR AddHeader: %v\n", err)
	}

	// 添加页脚
	if err := templateDoc.AddFooter(document.HeaderFooterTypeDefault, "Template Footer - Page"); err != nil {
		fmt.Printf("ERROR AddFooter: %v\n", err)
	}

	templateDoc.AddParagraph("Document content: {{content}}")

	// 保存模板
	templatePath := filepath.Join(outputDir, "issue76_template.docx")
	if err := templateDoc.Save(templatePath); err != nil {
		fmt.Printf("ERROR saving template: %v\n", err)
		return
	}
	fmt.Printf("Created template: %s\n", templatePath)

	// 使用模板引擎渲染
	renderer := document.NewTemplateRenderer()
	_, err := renderer.LoadTemplateFromFile("issue76", templatePath)
	if err != nil {
		fmt.Printf("ERROR loading template: %v\n", err)
		return
	}

	data := document.NewTemplateData()
	data.SetVariable("title", "Test Title")
	data.SetVariable("content", "Test Content")

	resultDoc, err := renderer.RenderTemplate("issue76", data)
	if err != nil {
		fmt.Printf("ERROR rendering: %v\n", err)
		return
	}

	resultPath := filepath.Join(outputDir, "issue76_rendered.docx")
	if err := resultDoc.Save(resultPath); err != nil {
		fmt.Printf("ERROR saving result: %v\n", err)
		return
	}
	fmt.Printf("Created rendered: %s\n", resultPath)
	fmt.Println("Issue #76: Check if header/footer exist in rendered document")
}

// Issue #91: 表格图片替换 (验证修复)
func testIssue91(outputDir string) {
	// 创建测试图片
	testImagePath := filepath.Join(outputDir, "test_image.png")
	if err := createTestImage(testImagePath); err != nil {
		fmt.Printf("ERROR creating test image: %v\n", err)
		return
	}

	// 创建模板
	templateDoc := document.New()
	templateDoc.AddParagraph("Image outside table:")
	templateDoc.AddParagraph("{{#image outside_img}}")

	table, _ := templateDoc.AddTable(&document.TableConfig{Rows: 2, Cols: 2})
	table.SetCellText(0, 0, "Header")
	table.SetCellText(0, 1, "Image Cell")
	table.SetCellText(1, 0, "Data")
	table.SetCellText(1, 1, "{{#image inside_img}}")

	templatePath := filepath.Join(outputDir, "issue91_template.docx")
	templateDoc.Save(templatePath)
	fmt.Printf("Created template: %s\n", templatePath)

	// 渲染
	renderer := document.NewTemplateRenderer()
	renderer.LoadTemplateFromFile("issue91", templatePath)

	data := document.NewTemplateData()
	data.SetImage("outside_img", testImagePath, nil)
	data.SetImage("inside_img", testImagePath, nil)

	resultDoc, err := renderer.RenderTemplate("issue91", data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	resultPath := filepath.Join(outputDir, "issue91_result.docx")
	resultDoc.Save(resultPath)
	fmt.Printf("Created result: %s\n", resultPath)
	fmt.Println("Issue #91 (FIXED): Verify both images are in word/media/")
}

// Issue #88: each循环粗体 (验证修复)
func testIssue88(outputDir string) {
	// 创建模板
	templateDoc := document.New()
	templateDoc.AddParagraph("Normal text (not from loop)")
	templateDoc.AddParagraph("{{#each items}}")
	templateDoc.AddParagraph("Loop item: {{name}} - {{value}}")
	templateDoc.AddParagraph("{{/each}}")
	templateDoc.AddParagraph("Text after loop")

	templatePath := filepath.Join(outputDir, "issue88_template.docx")
	templateDoc.Save(templatePath)
	fmt.Printf("Created template: %s\n", templatePath)

	// 渲染
	renderer := document.NewTemplateRenderer()
	renderer.LoadTemplateFromFile("issue88", templatePath)

	data := document.NewTemplateData()
	items := []interface{}{
		map[string]interface{}{"name": "Item1", "value": "Value1"},
		map[string]interface{}{"name": "Item2", "value": "Value2"},
	}
	data.SetList("items", items)

	resultDoc, err := renderer.RenderTemplate("issue88", data)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	resultPath := filepath.Join(outputDir, "issue88_result.docx")
	resultDoc.Save(resultPath)
	fmt.Printf("Created result: %s\n", resultPath)
	fmt.Println("Issue #88 (FIXED): Verify no <w:b/> tags on loop content")
}

// Helper: create test PNG image using Python PIL
func createTestImage(path string) error {
	cmd := exec.Command("python3", "-c", `
from PIL import Image
img = Image.new('RGB', (50, 50), color='blue')
img.save('`+path+`', 'PNG')
`)
	return cmd.Run()
}

// nopCloser wraps a Reader to implement ReadCloser
type nopCloser struct {
	*strings.Reader
}

func (nopCloser) Close() error { return nil }

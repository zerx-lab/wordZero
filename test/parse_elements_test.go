package test

import (
	"testing"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
)

// TestParseElementsIntegration 测试解析不同类型元素的集成测试
func TestParseElementsIntegration(t *testing.T) {
	// 创建新文档并添加不同类型的元素
	doc := document.New()

	// 添加段落
	para1 := doc.AddParagraph("这是第一个段落")
	para1.SetAlignment(document.AlignCenter)

	// 添加格式化段落
	titleFormat := &document.TextFormat{
		Bold:     true,
		FontSize: 16,
	}
	title := doc.AddFormattedParagraph("文档标题", titleFormat)
	title.SetAlignment(document.AlignCenter)

	// 添加表格
	tableConfig := &document.TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 5000,
		Data: [][]string{
			{"标题1", "标题2", "标题3"},
			{"数据1", "数据2", "数据3"},
			{"数据4", "数据5", "数据6"},
		},
	}
	_, err := doc.AddTable(tableConfig)
	if err != nil {
		t.Fatalf("添加表格失败: %v", err)
	}

	// 添加普通段落
	doc.AddParagraph("这是表格后的段落")

	// 保存文档
	testFile := "test_parse_elements.docx"
	err = doc.Save(testFile)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}
	defer func() {
		// 清理测试文件
		// os.Remove(testFile)
	}()

	t.Logf("文档保存成功: %s", testFile)

	// 重新打开文档测试解析
	t.Log("重新打开文档测试解析...")
	reopenedDoc, err := document.Open(testFile)
	if err != nil {
		t.Fatalf("打开文档失败: %v", err)
	}

	// 检查解析结果
	elementCount := len(reopenedDoc.Body.Elements)
	t.Logf("解析到的元素数量: %d", elementCount)

	// 验证元素数量（至少应该有段落和表格）
	if elementCount == 0 {
		t.Fatal("解析结果为空，没有解析到任何元素")
	}

	// 统计不同类型的元素
	var paragraphCount, tableCount, sectPrCount, unknownCount int

	for i, element := range reopenedDoc.Body.Elements {
		switch e := element.(type) {
		case *document.Paragraph:
			paragraphCount++
			t.Logf("元素 %d: 段落 - ", i+1)
			for _, run := range e.Runs {
				t.Logf("  文本: %s", run.Text.Content)
			}
		case *document.Table:
			tableCount++
			t.Logf("元素 %d: 表格 - %d行 %d列", i+1, len(e.Rows), e.GetColumnCount())
		case *document.SectionProperties:
			sectPrCount++
			t.Logf("元素 %d: 节属性", i+1)
		default:
			unknownCount++
			t.Logf("元素 %d: 未知类型 %T", i+1, element)
		}
	}

	// 验证解析结果
	t.Logf("解析结果统计: 段落=%d, 表格=%d, 节属性=%d, 未知=%d",
		paragraphCount, tableCount, sectPrCount, unknownCount)

	// 验证基本要求
	if paragraphCount == 0 {
		t.Error("没有解析到任何段落")
	}

	if tableCount == 0 {
		t.Error("没有解析到任何表格")
	}

	// 如果有节属性，说明解析逻辑能正确识别不同类型的元素
	if sectPrCount > 0 {
		t.Log("成功解析到节属性，证明动态解析逻辑工作正常")
	}

	t.Log("解析测试完成！")
}

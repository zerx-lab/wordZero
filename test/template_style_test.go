package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
	"github.com/ZeroHawkeye/wordZero/pkg/style"
)

// TestTemplateStylePreservation 测试模板渲染时样式保持
func TestTemplateStylePreservation(t *testing.T) {
	fmt.Println("=== 模板样式保持测试 ===")

	// 1. 创建带有丰富样式的模板文档
	templateDoc := document.New()

	// 添加标题段落，使用自定义样式
	titlePara := templateDoc.AddParagraph("项目报告：{{projectName}}")
	titlePara.SetStyle(style.StyleHeading1)

	// 添加副标题段落
	subtitlePara := templateDoc.AddParagraph("报告人：{{author}}")
	subtitlePara.SetStyle(style.StyleHeading2)

	// 添加正文段落，包含格式化的文本
	bodyPara := templateDoc.AddParagraph("")

	// 创建包含不同格式的Run
	normalRun := &document.Run{
		Text: document.Text{Content: "项目状态："},
		Properties: &document.RunProperties{
			FontFamily: &document.FontFamily{
				ASCII:    "Arial",
				HAnsi:    "Arial",
				EastAsia: "微软雅黑",
			},
			FontSize: &document.FontSize{Val: "24"}, // 12pt
		},
	}

	boldRun := &document.Run{
		Properties: &document.RunProperties{
			Bold: &document.Bold{},
			FontFamily: &document.FontFamily{
				ASCII:    "Arial",
				HAnsi:    "Arial",
				EastAsia: "微软雅黑",
			},
			FontSize: &document.FontSize{Val: "24"},
			Color:    &document.Color{Val: "FF0000"}, // 红色
		},
		Text: document.Text{Content: "{{status}}"},
	}

	endRun := &document.Run{
		Text: document.Text{Content: "，进度："},
		Properties: &document.RunProperties{
			FontFamily: &document.FontFamily{
				ASCII:    "Arial",
				HAnsi:    "Arial",
				EastAsia: "微软雅黑",
			},
			FontSize: &document.FontSize{Val: "24"},
		},
	}

	progressRun := &document.Run{
		Properties: &document.RunProperties{
			Bold:   &document.Bold{},
			Italic: &document.Italic{},
			FontFamily: &document.FontFamily{
				ASCII:    "Arial",
				HAnsi:    "Arial",
				EastAsia: "微软雅黑",
			},
			FontSize: &document.FontSize{Val: "28"},  // 14pt
			Color:    &document.Color{Val: "008000"}, // 绿色
		},
		Text: document.Text{Content: "{{progress}}%"},
	}

	bodyPara.Runs = []document.Run{*normalRun, *boldRun, *endRun, *progressRun}

	// 创建带样式的表格
	tableConfig := &document.TableConfig{
		Rows: 2,
		Cols: 3,
	}
	table, err := templateDoc.AddTable(tableConfig)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 设置表格头部样式
	headerRow := table.Rows[0]
	for i, cell := range headerRow.Cells {
		headers := []string{"任务", "负责人", "状态"}
		para := &document.Paragraph{
			Properties: &document.ParagraphProperties{
				ParagraphStyle: &document.ParagraphStyle{Val: style.StyleHeading3},
			},
			Runs: []document.Run{{
				Text: document.Text{Content: headers[i]},
				Properties: &document.RunProperties{
					Bold: &document.Bold{},
					FontFamily: &document.FontFamily{
						ASCII:    "Calibri",
						HAnsi:    "Calibri",
						EastAsia: "宋体",
					},
					FontSize: &document.FontSize{Val: "22"},  // 11pt
					Color:    &document.Color{Val: "FFFFFF"}, // 白色
				},
			}},
		}
		cell.Paragraphs = []document.Paragraph{*para}
	}

	// 设置数据行模板
	dataRow := table.Rows[1]
	templates := []string{"{{#each tasks}}", "{{name}}", "{{assignee}}", "{{status}}", "{{/each}}"}
	for i, cell := range dataRow.Cells {
		if i < len(templates) {
			para := &document.Paragraph{
				Properties: &document.ParagraphProperties{
					ParagraphStyle: &document.ParagraphStyle{Val: style.StyleNormal},
				},
				Runs: []document.Run{{
					Text: document.Text{Content: templates[i]},
					Properties: &document.RunProperties{
						FontFamily: &document.FontFamily{
							ASCII:    "Times New Roman",
							HAnsi:    "Times New Roman",
							EastAsia: "宋体",
						},
						FontSize: &document.FontSize{Val: "20"},  // 10pt
						Color:    &document.Color{Val: "000080"}, // 深蓝色
					},
				}},
			}
			cell.Paragraphs = []document.Paragraph{*para}
		}
	}

	// 保存模板文档
	templateFile := "test/output/style_template.docx"
	err = templateDoc.Save(templateFile)
	if err != nil {
		t.Fatalf("保存模板文档失败: %v", err)
	}
	fmt.Printf("✓ 创建样式模板文档: %s\n", templateFile)

	// 2. 从模板文档加载模板
	engine := document.NewTemplateEngine()
	_, err = engine.LoadTemplateFromDocument("style_template", templateDoc)
	if err != nil {
		t.Fatalf("加载模板失败: %v", err)
	}

	// 3. 准备测试数据
	data := document.NewTemplateData()
	data.SetVariable("projectName", "WordZero 开发项目")
	data.SetVariable("author", "张开发")
	data.SetVariable("status", "进行中")
	data.SetVariable("progress", "75")

	// 设置表格数据
	tasks := []interface{}{
		map[string]interface{}{
			"name":     "文档解析",
			"assignee": "张三",
			"status":   "已完成",
		},
		map[string]interface{}{
			"name":     "样式系统",
			"assignee": "李四",
			"status":   "进行中",
		},
		map[string]interface{}{
			"name":     "测试用例",
			"assignee": "王五",
			"status":   "待开始",
		},
	}
	data.SetList("tasks", tasks)

	// 4. 渲染模板
	resultDoc, err := engine.RenderTemplateToDocument("style_template", data)
	if err != nil {
		t.Fatalf("渲染模板失败: %v", err)
	}

	// 5. 保存结果文档
	outputFile := "test/output/style_result_" + time.Now().Format("20060102_150405") + ".docx"
	err = resultDoc.Save(outputFile)
	if err != nil {
		t.Fatalf("保存结果文档失败: %v", err)
	}

	fmt.Printf("✓ 生成结果文档: %s\n", outputFile)

	// 6. 验证样式是否保持
	verifyDocumentStyles(t, resultDoc)

	fmt.Println("✓ 模板样式保持测试完成")
}

// verifyDocumentStyles 验证文档样式是否正确保持
func verifyDocumentStyles(t *testing.T, doc *document.Document) {
	fmt.Println("\n=== 验证样式保持情况 ===")

	// 检查文档元素
	if len(doc.Body.Elements) == 0 {
		t.Error("文档没有元素")
		return
	}

	elementCount := 0
	styledElements := 0

	for i, element := range doc.Body.Elements {
		elementCount++

		switch elem := element.(type) {
		case *document.Paragraph:
			fmt.Printf("段落 %d: ", i+1)

			// 检查段落样式
			if elem.Properties != nil && elem.Properties.ParagraphStyle != nil {
				fmt.Printf("段落样式=%s, ", elem.Properties.ParagraphStyle.Val)
				styledElements++
			} else {
				fmt.Printf("段落样式=无, ")
			}

			// 检查Run样式
			runStyleCount := 0
			for j, run := range elem.Runs {
				if run.Properties != nil {
					runStyleCount++

					// 检查关键样式属性
					hasFont := run.Properties.FontFamily != nil
					hasBold := run.Properties.Bold != nil
					hasColor := run.Properties.Color != nil
					hasSize := run.Properties.FontSize != nil

					fmt.Printf("Run%d(字体:%t,粗体:%t,颜色:%t,大小:%t) ",
						j+1, hasFont, hasBold, hasColor, hasSize)
				}
			}

			fmt.Printf("(共%d个带样式Run)\n", runStyleCount)

		case *document.Table:
			fmt.Printf("表格 %d: %d行%d列\n", i+1, len(elem.Rows), len(elem.Rows[0].Cells))

			// 检查表格样式
			tableStyledCells := 0
			for rowIdx, row := range elem.Rows {
				for cellIdx, cell := range row.Cells {
					for paraIdx, para := range cell.Paragraphs {
						if para.Properties != nil && para.Properties.ParagraphStyle != nil {
							tableStyledCells++
							fmt.Printf("  行%d列%d段落%d: 样式=%s\n",
								rowIdx+1, cellIdx+1, paraIdx+1, para.Properties.ParagraphStyle.Val)
						}

						for runIdx, run := range para.Runs {
							if run.Properties != nil {
								hasFont := run.Properties.FontFamily != nil
								hasBold := run.Properties.Bold != nil
								hasColor := run.Properties.Color != nil
								hasSize := run.Properties.FontSize != nil

								if hasFont || hasBold || hasColor || hasSize {
									fmt.Printf("    Run%d: 字体:%t,粗体:%t,颜色:%t,大小:%t\n",
										runIdx+1, hasFont, hasBold, hasColor, hasSize)
								}
							}
						}
					}
				}
			}

			if tableStyledCells > 0 {
				styledElements++
				fmt.Printf("  表格中有 %d 个带样式的单元格\n", tableStyledCells)
			}
		}
	}

	fmt.Printf("\n总结: 共 %d 个元素，其中 %d 个有样式信息\n", elementCount, styledElements)

	// 基本验证
	if elementCount == 0 {
		t.Error("文档中没有元素")
	}

	if styledElements == 0 {
		t.Error("❌ 严重问题：所有样式都丢失了！")
	} else {
		fmt.Printf("✓ 检测到 %d 个元素保持了样式\n", styledElements)
	}
}

// TestTemplateStyleIssues 专门测试样式问题的原因
func TestTemplateStyleIssues(t *testing.T) {
	fmt.Println("\n=== 样式问题诊断测试 ===")

	// 创建简单的模板文档
	templateDoc := document.New()

	// 添加一个带样式的段落
	para := templateDoc.AddParagraph("测试变量：{{testVar}}")
	para.SetStyle(style.StyleHeading1)

	// 检查模板文档的样式
	fmt.Println("模板文档检查:")
	if para.Properties != nil && para.Properties.ParagraphStyle != nil {
		fmt.Printf("✓ 段落样式: %s\n", para.Properties.ParagraphStyle.Val)
	} else {
		fmt.Println("❌ 段落没有样式")
	}

	// 检查样式管理器
	styleManager := templateDoc.GetStyleManager()
	heading1Style := styleManager.GetStyle(style.StyleHeading1)
	if heading1Style != nil {
		fmt.Printf("✓ StyleManager中存在Heading1样式: %s\n", heading1Style.StyleID)
	} else {
		fmt.Println("❌ StyleManager中缺少Heading1样式")
	}

	// 加载为模板
	engine := document.NewTemplateEngine()
	template, err := engine.LoadTemplateFromDocument("test_template", templateDoc)
	if err != nil {
		t.Fatalf("加载模板失败: %v", err)
	}

	// 检查模板的BaseDoc
	if template.BaseDoc != nil {
		fmt.Println("✓ 模板有BaseDoc")

		// 检查BaseDoc的样式管理器
		if template.BaseDoc.GetStyleManager() != nil {
			fmt.Println("✓ BaseDoc有样式管理器")

			baseHeading1 := template.BaseDoc.GetStyleManager().GetStyle(style.StyleHeading1)
			if baseHeading1 != nil {
				fmt.Printf("✓ BaseDoc样式管理器中有Heading1: %s\n", baseHeading1.StyleID)
			} else {
				fmt.Println("❌ BaseDoc样式管理器中缺少Heading1")
			}
		} else {
			fmt.Println("❌ BaseDoc没有样式管理器")
		}
	} else {
		fmt.Println("❌ 模板没有BaseDoc")
	}

	// 渲染模板
	data := document.NewTemplateData()
	data.SetVariable("testVar", "测试值")

	resultDoc, err := engine.RenderTemplateToDocument("test_template", data)
	if err != nil {
		t.Fatalf("渲染失败: %v", err)
	}

	// 检查结果文档
	fmt.Println("\n结果文档检查:")
	if len(resultDoc.Body.Elements) > 0 {
		if para, ok := resultDoc.Body.Elements[0].(*document.Paragraph); ok {
			if para.Properties != nil && para.Properties.ParagraphStyle != nil {
				fmt.Printf("✓ 结果段落样式: %s\n", para.Properties.ParagraphStyle.Val)
			} else {
				fmt.Println("❌ 结果段落没有样式")
			}

			// 检查文本内容
			fullText := ""
			for _, run := range para.Runs {
				fullText += run.Text.Content
			}
			fmt.Printf("文本内容: %s\n", fullText)
		}
	}

	// 检查结果文档的样式管理器
	resultStyleManager := resultDoc.GetStyleManager()
	if resultStyleManager != nil {
		fmt.Println("✓ 结果文档有样式管理器")

		resultHeading1 := resultStyleManager.GetStyle(style.StyleHeading1)
		if resultHeading1 != nil {
			fmt.Printf("✓ 结果样式管理器中有Heading1: %s\n", resultHeading1.StyleID)
		} else {
			fmt.Println("❌ 结果样式管理器中缺少Heading1")
		}
	} else {
		fmt.Println("❌ 结果文档没有样式管理器")
	}

	// 保存结果用于手动检查
	outputFile := "test/output/style_diagnosis_" + time.Now().Format("20060102_150405") + ".docx"
	err = resultDoc.Save(outputFile)
	if err != nil {
		t.Fatalf("保存失败: %v", err)
	}
	fmt.Printf("✓ 诊断结果保存到: %s\n", outputFile)
}

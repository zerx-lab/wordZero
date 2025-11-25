// Package test 嵌套表格模板测试
package test

import (
	"os"
	"testing"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
)

// TestNestedTableTemplate 测试嵌套表格模板功能
func TestNestedTableTemplate(t *testing.T) {
	// 确保输出目录存在
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			t.Fatalf("创建输出目录失败: %v", err)
		}
	}

	// 创建带有嵌套表格的文档
	doc := document.New()
	doc.AddParagraph("嵌套表格测试文档")
	doc.AddParagraph("")

	// 创建外层表格（2行2列）
	outerTable, err := doc.CreateTable(&document.TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 9000,
	})
	if err != nil {
		t.Fatalf("创建外层表格失败: %v", err)
	}
	doc.Body.Elements = append(doc.Body.Elements, outerTable)

	// 设置外层表格第一行内容
	outerTable.Rows[0].Cells[0].Paragraphs[0].Runs[0].Text.Content = "简历"
	outerTable.Rows[0].Cells[1].Paragraphs[0].Runs[0].Text.Content = "{{resume}}"

	// 在外层表格第二行第一列添加文本
	outerTable.Rows[1].Cells[0].Paragraphs[0].Runs[0].Text.Content = "家庭主要成为及重要社会关系"

	// 创建内层表格（嵌套在外层表格的第二行第二列）
	innerTable := &document.Table{
		Properties: &document.TableProperties{
			TableW: &document.TableWidth{
				W:    "4000",
				Type: "dxa",
			},
		},
		Rows: []document.TableRow{
			{
				Cells: []document.TableCell{
					{
						Paragraphs: []document.Paragraph{
							{
								Runs: []document.Run{
									{
										Text: document.Text{Content: "姓名"},
									},
								},
							},
						},
					},
					{
						Paragraphs: []document.Paragraph{
							{
								Runs: []document.Run{
									{
										Text: document.Text{Content: "年龄"},
									},
								},
							},
						},
					},
					{
						Paragraphs: []document.Paragraph{
							{
								Runs: []document.Run{
									{
										Text: document.Text{Content: "性别"},
									},
								},
							},
						},
					},
					{
						Paragraphs: []document.Paragraph{
							{
								Runs: []document.Run{
									{
										Text: document.Text{Content: "关系"},
									},
								},
							},
						},
					},
				},
			},
			{
				Cells: []document.TableCell{
					{
						Paragraphs: []document.Paragraph{
							{
								Runs: []document.Run{
									{
										Text: document.Text{Content: "{{#each family_members}}{{name}}"},
									},
								},
							},
						},
					},
					{
						Paragraphs: []document.Paragraph{
							{
								Runs: []document.Run{
									{
										Text: document.Text{Content: "{{age}}{{/each}}"},
									},
								},
							},
						},
					},
					{
						Paragraphs: []document.Paragraph{
							{
								Runs: []document.Run{
									{
										Text: document.Text{Content: "{{gender}}"},
									},
								},
							},
						},
					},
					{
						Paragraphs: []document.Paragraph{
							{
								Runs: []document.Run{
									{
										Text: document.Text{Content: "{{relationship}}"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// 将内层表格嵌套到外层表格的单元格中
	outerTable.Rows[1].Cells[1].Tables = []document.Table{*innerTable}
	// 确保外层表格的单元格有至少一个段落（OOXML要求）
	if len(outerTable.Rows[1].Cells[1].Paragraphs) == 0 {
		outerTable.Rows[1].Cells[1].Paragraphs = []document.Paragraph{
			{
				Runs: []document.Run{
					{
						Text: document.Text{Content: ""},
					},
				},
			},
		}
	}

	// 保存模板文档
	templatePath := "output/nested_table_template.docx"
	err = doc.Save(templatePath)
	if err != nil {
		t.Fatalf("保存模板文档失败: %v", err)
	}

	// 创建模板引擎
	engine := document.NewTemplateEngine()
	_, err = engine.LoadTemplateFromDocument("nested_template", doc)
	if err != nil {
		t.Fatalf("加载模板失败: %v", err)
	}

	// 准备测试数据
	data := document.NewTemplateData()
	data.SetVariable("resume", "个人简历")
	data.SetList("family_members", []interface{}{
		map[string]interface{}{
			"name":         "张三",
			"age":          "45",
			"gender":       "男",
			"relationship": "父亲",
		},
		map[string]interface{}{
			"name":         "李四",
			"age":          "43",
			"gender":       "女",
			"relationship": "母亲",
		},
		map[string]interface{}{
			"name":         "张小明",
			"age":          "20",
			"gender":       "男",
			"relationship": "本人",
		},
	})

	// 渲染模板
	renderedDoc, err := engine.RenderTemplateToDocument("nested_template", data)
	if err != nil {
		t.Fatalf("渲染模板失败: %v", err)
	}

	// 保存渲染后的文档
	outputPath := "output/nested_table_rendered.docx"
	err = renderedDoc.Save(outputPath)
	if err != nil {
		t.Fatalf("保存渲染文档失败: %v", err)
	}

	// 验证嵌套表格是否存在
	tables := renderedDoc.Body.GetTables()
	if len(tables) < 1 {
		t.Fatalf("渲染后的文档应该至少包含1个表格，实际包含 %d 个", len(tables))
	}

	// 检查外层表格是否有内容
	outerTableRendered := tables[0]
	if len(outerTableRendered.Rows) < 2 {
		t.Fatalf("外层表格应该有至少2行，实际有 %d 行", len(outerTableRendered.Rows))
	}

	// 验证嵌套表格是否被保留
	nestedTables := outerTableRendered.Rows[1].Cells[1].Tables
	if len(nestedTables) == 0 {
		t.Fatalf("嵌套表格在渲染后消失了！应该存在1个嵌套表格")
	}

	// 验证嵌套表格的行数（应该有1个标题行 + 3个数据行）
	nestedTable := nestedTables[0]
	expectedRows := 4 // 1 header + 3 data rows
	if len(nestedTable.Rows) != expectedRows {
		t.Errorf("嵌套表格应该有 %d 行，实际有 %d 行", expectedRows, len(nestedTable.Rows))
	}

	// 验证数据是否正确渲染
	if len(nestedTable.Rows) >= 2 {
		firstDataRow := nestedTable.Rows[1]
		if len(firstDataRow.Cells) >= 1 {
			nameCell := firstDataRow.Cells[0]
			if len(nameCell.Paragraphs) > 0 && len(nameCell.Paragraphs[0].Runs) > 0 {
				name := nameCell.Paragraphs[0].Runs[0].Text.Content
				if name != "张三" {
					t.Errorf("第一行数据的姓名应该是'张三'，实际是'%s'", name)
				}
			}
		}
	}

	t.Logf("嵌套表格模板测试通过！")
	t.Logf("模板文档: %s", templatePath)
	t.Logf("渲染文档: %s", outputPath)
}

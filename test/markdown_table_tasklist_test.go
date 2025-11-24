package test

import (
	"testing"

	"github.com/ZeroHawkeye/wordZero/pkg/markdown"
)

// TestMarkdownTableConversion 测试Markdown表格转换
func TestMarkdownTableConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantErr  bool
	}{
		{
			name: "简单表格",
			markdown: `| 姓名 | 年龄 | 城市 |
|------|------|------|
| 张三 | 25   | 北京 |
| 李四 | 30   | 上海 |`,
			wantErr: false,
		},
		{
			name: "对齐表格",
			markdown: `| 左对齐 | 居中对齐 | 右对齐 |
|:-------|:--------:|-------:|
| 内容1  |   内容2  |  内容3 |`,
			wantErr: false,
		},
		{
			name: "复杂表格",
			markdown: `| 功能 | 状态 | 描述 |
|------|------|------|
| **表格支持** | ✅ | 完整的GFM表格转换 |
| *任务列表* | ✅ | 支持复选框显示 |
| 代码块 | ✅ | 等宽字体显示 |`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := markdown.DefaultOptions()
			opts.EnableTables = true
			converter := markdown.NewConverter(opts)

			doc, err := converter.ConvertString(tt.markdown, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if doc == nil && !tt.wantErr {
				t.Error("ConvertString() returned nil document")
				return
			}

			// 验证文档包含表格
			if doc != nil {
				tables := doc.Body.GetTables()
				if len(tables) == 0 {
					t.Error("Expected document to contain at least one table")
				}
			}
		})
	}
}

// TestMarkdownTaskListConversion 测试Markdown任务列表转换
func TestMarkdownTaskListConversion(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantErr  bool
	}{
		{
			name: "简单任务列表",
			markdown: `- [x] 已完成任务
- [ ] 未完成任务`,
			wantErr: false,
		},
		{
			name: "嵌套任务列表",
			markdown: `- [x] 主要任务
  - [x] 子任务1
  - [ ] 子任务2
- [ ] 另一个主要任务`,
			wantErr: false,
		},
		{
			name: "混合格式任务列表",
			markdown: `- [x] **重要**已完成任务
- [ ] *普通*未完成任务
- [x] 包含代码的任务`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := markdown.DefaultOptions()
			opts.EnableTaskList = true
			converter := markdown.NewConverter(opts)

			doc, err := converter.ConvertString(tt.markdown, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if doc == nil && !tt.wantErr {
				t.Error("ConvertString() returned nil document")
				return
			}

			// 验证文档包含段落（任务列表项）
			if doc != nil {
				paragraphs := doc.Body.GetParagraphs()
				if len(paragraphs) == 0 {
					t.Error("Expected document to contain at least one paragraph for task items")
				}
			}
		})
	}
}

// TestMarkdownCombinedFeatures 测试表格和任务列表的组合使用
func TestMarkdownCombinedFeatures(t *testing.T) {
	markdownContent := `# 项目进度

## 功能列表

| 功能 | 状态 | 备注 |
|------|------|------|
| 表格 | ✅ | 已实现 |
| 任务列表 | ✅ | 已实现 |

## 待办事项

- [x] 实现表格功能
- [x] 实现任务列表功能
- [ ] 编写文档
- [ ] 完善测试`

	opts := markdown.HighQualityOptions()
	opts.EnableTables = true
	opts.EnableTaskList = true
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Errorf("ConvertString() error = %v", err)
		return
	}

	if doc == nil {
		t.Error("ConvertString() returned nil document")
		return
	}

	// 验证包含表格
	tables := doc.Body.GetTables()
	if len(tables) == 0 {
		t.Error("Expected document to contain at least one table")
	}

	// 验证包含段落（任务列表和其他内容）
	paragraphs := doc.Body.GetParagraphs()
	if len(paragraphs) == 0 {
		t.Error("Expected document to contain paragraphs")
	}
}

// TestTableAlignment 测试表格对齐功能
func TestTableAlignment(t *testing.T) {
	markdownContent := `| 左对齐 | 居中 | 右对齐 |
|:-------|:----:|-------:|
| Left   | Center | Right |`

	opts := markdown.DefaultOptions()
	opts.EnableTables = true
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Errorf("ConvertString() error = %v", err)
		return
	}

	if doc == nil {
		t.Error("ConvertString() returned nil document")
		return
	}

	// 验证表格存在
	tables := doc.Body.GetTables()
	if len(tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(tables))
		return
	}

	table := tables[0]
	if table.GetRowCount() != 2 {
		t.Errorf("Expected 2 rows, got %d", table.GetRowCount())
	}

	if table.GetColumnCount() != 3 {
		t.Errorf("Expected 3 columns, got %d", table.GetColumnCount())
	}
}

// TestTaskListCheckboxes 测试任务列表复选框
func TestTaskListCheckboxes(t *testing.T) {
	markdownContent := `- [x] 选中的任务
- [ ] 未选中的任务`

	opts := markdown.DefaultOptions()
	opts.EnableTaskList = true
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Errorf("ConvertString() error = %v", err)
		return
	}

	if doc == nil {
		t.Error("ConvertString() returned nil document")
		return
	}

	// 验证包含段落
	paragraphs := doc.Body.GetParagraphs()
	if len(paragraphs) < 2 {
		t.Errorf("Expected at least 2 paragraphs for task items, got %d", len(paragraphs))
	}
}

// TestMarkdownSoftLineBreaks 测试Markdown软换行（单个\n）的处理
func TestMarkdownSoftLineBreaks(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		wantErr     bool
		description string
	}{
		{
			name:        "单个软换行",
			markdown:    "第一行\n第二行",
			wantErr:     false,
			description: "两行文本用单个换行符分隔，应该在Word中用空格连接",
		},
		{
			name:        "多个软换行",
			markdown:    "第一行\n第二行\n第三行",
			wantErr:     false,
			description: "多行文本用单个换行符分隔",
		},
		{
			name:        "混合格式与软换行",
			markdown:    "**粗体文本**\n*斜体文本*\n普通文本",
			wantErr:     false,
			description: "格式化文本混合软换行",
		},
		{
			name: "问题报告中的内容",
			markdown: `
**日期：** 2024年___月___日 

---

### **附件关键条款索引** 
#### **附件一（V2）核心约束：** 
| 任务 | 法律要点 |
|------|---------|
| 任务1 | 要点1 |`,
			wantErr:     false,
			description: "实际问题报告中的Markdown内容，包含软换行和表格",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := markdown.DefaultOptions()
			opts.EnableGFM = true
			opts.EnableTables = true
			converter := markdown.NewConverter(opts)

			doc, err := converter.ConvertString(tt.markdown, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if doc == nil && !tt.wantErr {
				t.Error("ConvertString() returned nil document")
				return
			}

			// 验证文档至少有一些内容
			if doc != nil {
				if len(doc.Body.Elements) == 0 {
					t.Error("Expected document to have some content")
				}
			}
		})
	}
}

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func TestTextFormatting(t *testing.T) {
	// 设置测试日志级别
	document.SetGlobalLevel(document.LogLevelWarn)

	// 创建测试文档
	doc := document.New()

	// 测试基本格式化
	format := &document.TextFormat{
		Bold:       true,
		Italic:     true,
		FontSize:   14,
		FontColor:  "FF0000",
		FontFamily: "Arial",
	}

	p := doc.AddFormattedParagraph("测试格式化文本", format)
	if p == nil {
		t.Error("AddFormattedParagraph 返回了 nil")
	}

	// 检查段落是否被正确添加
	paragraphs := doc.Body.GetParagraphs()
	if len(paragraphs) != 1 {
		t.Errorf("预期1个段落，但得到了 %d 个", len(paragraphs))
	}

	// 检查运行属性
	if len(paragraphs[0].Runs) == 0 {
		t.Error("段落中没有运行")
	} else {
		run := paragraphs[0].Runs[0]
		if run.Properties == nil {
			t.Error("运行属性为空")
		} else {
			if run.Properties.Bold == nil {
				t.Error("粗体属性未设置")
			}
			if run.Properties.Italic == nil {
				t.Error("斜体属性未设置")
			}
			if run.Properties.FontSize == nil {
				t.Error("字体大小属性未设置")
			}
			if run.Properties.Color == nil {
				t.Error("颜色属性未设置")
			}
			if run.Properties.FontFamily == nil {
				t.Error("字体族属性未设置")
			}
		}
	}
}

func TestParagraphAlignment(t *testing.T) {
	doc := document.New()

	// 测试各种对齐方式
	alignments := []document.AlignmentType{
		document.AlignLeft,
		document.AlignCenter,
		document.AlignRight,
		document.AlignJustify,
	}

	for _, align := range alignments {
		p := doc.AddParagraph("测试对齐")
		p.SetAlignment(align)

		if p.Properties == nil {
			t.Errorf("段落属性为空，对齐方式: %s", align)
		} else if p.Properties.Justification == nil {
			t.Errorf("对齐属性为空，对齐方式: %s", align)
		} else if p.Properties.Justification.Val != string(align) {
			t.Errorf("对齐方式不匹配，预期: %s，实际: %s", align, p.Properties.Justification.Val)
		}
	}
}

func TestParagraphSpacing(t *testing.T) {
	doc := document.New()

	p := doc.AddParagraph("测试间距")

	config := &document.SpacingConfig{
		LineSpacing:     1.5,
		BeforePara:      12,
		AfterPara:       6,
		FirstLineIndent: 24,
	}

	p.SetSpacing(config)

	if p.Properties == nil {
		t.Error("段落属性为空")
		return
	}

	if p.Properties.Spacing == nil {
		t.Error("间距属性为空")
	} else {
		// 检查间距值（TWIPs单位）
		if p.Properties.Spacing.Before != "240" { // 12 * 20
			t.Errorf("段前间距不正确，预期: 240，实际: %s", p.Properties.Spacing.Before)
		}
		if p.Properties.Spacing.After != "120" { // 6 * 20
			t.Errorf("段后间距不正确，预期: 120，实际: %s", p.Properties.Spacing.After)
		}
		if p.Properties.Spacing.Line != "360" { // 1.5 * 240
			t.Errorf("行间距不正确，预期: 360，实际: %s", p.Properties.Spacing.Line)
		}
	}

	if p.Properties.Indentation == nil {
		t.Error("缩进属性为空")
	} else {
		if p.Properties.Indentation.FirstLine != "480" { // 24 * 20
			t.Errorf("首行缩进不正确，预期: 480，实际: %s", p.Properties.Indentation.FirstLine)
		}
	}
}

func TestAddFormattedText(t *testing.T) {
	doc := document.New()

	p := doc.AddParagraph("基础文本")

	// 添加格式化文本
	format := &document.TextFormat{
		Bold:      true,
		FontColor: "0000FF",
	}

	p.AddFormattedText("附加的格式化文本", format)

	if len(p.Runs) != 2 {
		t.Errorf("预期2个运行，但得到了 %d 个", len(p.Runs))
	}

	// 检查第二个运行的属性
	if len(p.Runs) >= 2 {
		run := p.Runs[1]
		if run.Properties == nil {
			t.Error("第二个运行的属性为空")
		} else {
			if run.Properties.Bold == nil {
				t.Error("第二个运行的粗体属性未设置")
			}
			if run.Properties.Color == nil {
				t.Error("第二个运行的颜色属性未设置")
			}
		}
	}
}

func TestDocumentSaveAndOpen(t *testing.T) {
	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "wordzero_test")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir)

	filename := filepath.Join(tempDir, "test_formatted.docx")

	// 创建带格式的文档
	doc := document.New()

	format := &document.TextFormat{
		Bold:     true,
		FontSize: 16,
	}

	p := doc.AddFormattedParagraph("测试保存和加载", format)
	p.SetAlignment(document.AlignCenter)

	// 保存文档
	err := doc.Save(filename)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}

	// 确认文件存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("保存的文件不存在")
	}

	// 重新打开文档
	openedDoc, err := document.Open(filename)
	if err != nil {
		t.Fatalf("打开文档失败: %v", err)
	}

	// 验证内容
	paragraphs := openedDoc.Body.GetParagraphs()
	if len(paragraphs) != 1 {
		t.Errorf("预期1个段落，但得到了 %d 个", len(paragraphs))
	}

	if len(paragraphs) > 0 {
		para := paragraphs[0]

		// 检查对齐方式
		if para.Properties == nil || para.Properties.Justification == nil {
			t.Error("段落对齐属性丢失")
		} else if para.Properties.Justification.Val != string(document.AlignCenter) {
			t.Errorf("对齐方式不匹配，预期: %s，实际: %s",
				document.AlignCenter, para.Properties.Justification.Val)
		}

		// 检查文本内容
		if len(para.Runs) > 0 {
			if para.Runs[0].Text.Content != "测试保存和加载" {
				t.Errorf("文本内容不匹配，预期: %s，实际: %s",
					"测试保存和加载", para.Runs[0].Text.Content)
			}

			// 检查格式属性
			if para.Runs[0].Properties == nil {
				t.Error("运行属性丢失")
			} else {
				if para.Runs[0].Properties.Bold == nil {
					t.Error("粗体属性丢失")
				}
				if para.Runs[0].Properties.FontSize == nil {
					t.Error("字体大小属性丢失")
				}
			}
		}
	}
}

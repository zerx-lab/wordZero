package test

import (
	"fmt"
	"testing"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func TestFootnoteConfig(t *testing.T) {
	doc := document.New()

	// 测试设置脚注配置
	config := &document.FootnoteConfig{
		NumberFormat: document.FootnoteFormatDecimal,
		StartNumber:  1,
		RestartEach:  document.FootnoteRestartContinuous,
		Position:     document.FootnotePositionPageBottom,
	}

	err := doc.SetFootnoteConfig(config)
	if err != nil {
		// 打印详细错误信息用于调试
		fmt.Printf("设置脚注配置失败的详细错误: %v\n", err)

		// 检查settings.xml内容
		parts := doc.GetParts()
		if settingsXML, exists := parts["word/settings.xml"]; exists {
			fmt.Printf("生成的settings.xml内容:\n%s\n", string(settingsXML))
		}

		t.Fatalf("设置脚注配置失败: %v", err)
	}

	// 验证settings.xml是否已创建
	_, exists := doc.GetParts()["word/settings.xml"]
	if !exists {
		t.Error("settings.xml文件未创建")
	}

	// 添加脚注测试
	err = doc.AddFootnote("这是正文文本", "这是脚注内容")
	if err != nil {
		t.Fatalf("添加脚注失败: %v", err)
	}

	// 验证脚注文件是否已创建
	_, exists = doc.GetParts()["word/footnotes.xml"]
	if !exists {
		t.Error("footnotes.xml文件未创建")
	}

	// 验证脚注数量
	count := doc.GetFootnoteCount()
	if count != 1 {
		t.Errorf("预期脚注数量为1，实际为%d", count)
	}
}

func TestEndnoteConfig(t *testing.T) {
	doc := document.New()

	// 添加尾注测试
	err := doc.AddEndnote("这是正文文本", "这是尾注内容")
	if err != nil {
		t.Fatalf("添加尾注失败: %v", err)
	}

	// 验证尾注文件是否已创建
	_, exists := doc.GetParts()["word/endnotes.xml"]
	if !exists {
		t.Error("endnotes.xml文件未创建")
	}

	// 验证尾注数量
	count := doc.GetEndnoteCount()
	if count != 1 {
		t.Errorf("预期尾注数量为1，实际为%d", count)
	}
}

func TestFootnoteNumberFormats(t *testing.T) {
	doc := document.New()

	// 测试不同的编号格式
	formats := []document.FootnoteNumberFormat{
		document.FootnoteFormatDecimal,
		document.FootnoteFormatLowerRoman,
		document.FootnoteFormatUpperRoman,
		document.FootnoteFormatLowerLetter,
		document.FootnoteFormatUpperLetter,
		document.FootnoteFormatSymbol,
	}

	for _, format := range formats {
		config := &document.FootnoteConfig{
			NumberFormat: format,
			StartNumber:  1,
			RestartEach:  document.FootnoteRestartContinuous,
			Position:     document.FootnotePositionPageBottom,
		}

		err := doc.SetFootnoteConfig(config)
		if err != nil {
			t.Fatalf("设置脚注格式%s失败: %v", format, err)
		}
	}
}

func TestFootnotePositions(t *testing.T) {
	doc := document.New()

	// 测试不同的脚注位置
	positions := []document.FootnotePosition{
		document.FootnotePositionPageBottom,
		document.FootnotePositionBeneathText,
		document.FootnotePositionSectionEnd,
		document.FootnotePositionDocumentEnd,
	}

	for _, position := range positions {
		config := &document.FootnoteConfig{
			NumberFormat: document.FootnoteFormatDecimal,
			StartNumber:  1,
			RestartEach:  document.FootnoteRestartContinuous,
			Position:     position,
		}

		err := doc.SetFootnoteConfig(config)
		if err != nil {
			t.Fatalf("设置脚注位置%s失败: %v", position, err)
		}
	}
}

func TestDefaultFootnoteConfig(t *testing.T) {
	config := document.DefaultFootnoteConfig()

	if config.NumberFormat != document.FootnoteFormatDecimal {
		t.Errorf("默认编号格式错误，预期%s，实际%s",
			document.FootnoteFormatDecimal, config.NumberFormat)
	}

	if config.StartNumber != 1 {
		t.Errorf("默认起始编号错误，预期1，实际%d", config.StartNumber)
	}

	if config.RestartEach != document.FootnoteRestartContinuous {
		t.Errorf("默认重新开始规则错误，预期%s，实际%s",
			document.FootnoteRestartContinuous, config.RestartEach)
	}

	if config.Position != document.FootnotePositionPageBottom {
		t.Errorf("默认位置错误，预期%s，实际%s",
			document.FootnotePositionPageBottom, config.Position)
	}
}

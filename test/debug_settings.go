package test

import (
	"encoding/xml"
	"fmt"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func DebugSettings() {
	doc := document.New()

	// 创建脚注配置
	config := &document.FootnoteConfig{
		NumberFormat: document.FootnoteFormatDecimal,
		StartNumber:  1,
		RestartEach:  document.FootnoteRestartContinuous,
		Position:     document.FootnotePositionPageBottom,
	}

	// 尝试设置配置
	err := doc.SetFootnoteConfig(config)
	if err != nil {
		fmt.Printf("设置脚注配置错误: %v\n", err)
		return
	}

	// 检查生成的settings.xml内容
	parts := doc.GetParts()
	if settingsXML, exists := parts["word/settings.xml"]; exists {
		fmt.Printf("Settings XML内容:\n%s\n", string(settingsXML))

		// 尝试解析生成的XML
		var settings document.Settings
		err = xml.Unmarshal(settingsXML, &settings)
		if err != nil {
			fmt.Printf("解析XML失败: %v\n", err)
		} else {
			fmt.Printf("XML解析成功!\n")
		}
	} else {
		fmt.Printf("settings.xml文件未找到\n")
	}
}

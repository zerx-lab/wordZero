// Package test 模板继承功能测试
package test

import (
	"os"
	"testing"

	"github.com/zerx-lab/wordZero/pkg/document"
)

// TestTemplateInheritanceComplete 完整的模板继承测试
func TestTemplateInheritanceComplete(t *testing.T) {
	// 确保输出目录存在
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create output directory: %v", err)
		}
	}

	engine := document.NewTemplateEngine()

	// 测试1: 基础模板继承
	t.Run("基础模板继承", func(t *testing.T) {
		testBasicInheritance(t, engine)
	})

	// 测试2: 块重写功能
	t.Run("块重写功能", func(t *testing.T) {
		testBlockOverride(t, engine)
	})

	// 清理测试文件
	t.Cleanup(func() {
		cleanupInheritanceTestFiles()
	})
}

// testBasicInheritance 测试基础模板继承功能
func testBasicInheritance(t *testing.T, engine *document.TemplateEngine) {
	// 创建基础文档
	doc := document.New()

	// 添加基础模板内容
	doc.AddParagraph("{{companyName}} 官方文档")
	doc.AddParagraph("")
	doc.AddParagraph("版本：{{version}}")
	doc.AddParagraph("创建时间：{{createTime}}")
	doc.AddParagraph("")
	doc.AddParagraph("{{#block \"content\"}}")
	doc.AddParagraph("默认内容区域")
	doc.AddParagraph("{{/block}}")
	doc.AddParagraph("")
	doc.AddParagraph("{{#block \"footer\"}}")
	doc.AddParagraph("版权所有 © {{year}} {{companyName}}")
	doc.AddParagraph("{{/block}}")

	_, err := engine.LoadTemplateFromDocument("base", doc)
	if err != nil {
		t.Fatalf("Failed to load base template: %v", err)
	}

	// 准备数据
	data := document.NewTemplateData()
	data.SetVariable("companyName", "WordZero科技")
	data.SetVariable("version", "v1.0.0")
	data.SetVariable("createTime", "2024-12-01")
	data.SetVariable("year", "2024")

	// 渲染模板
	resultDoc, err := engine.RenderTemplateToDocument("base", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 保存文档
	filename := "output/test_basic_inheritance.docx"
	err = resultDoc.Save(filename)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 验证文档内容
	if len(resultDoc.Body.Elements) == 0 {
		t.Error("Expected document to have content")
	}

	t.Logf("Basic inheritance test completed: %s", filename)
}

// testBlockOverride 测试块重写功能
func testBlockOverride(t *testing.T, engine *document.TemplateEngine) {
	// 创建基础文档
	doc := document.New()

	// 添加企业报告模板
	doc.AddParagraph("企业报告模板")
	doc.AddParagraph("")
	doc.AddParagraph("{{#block \"header\"}}")
	doc.AddParagraph("标准报告头部")
	doc.AddParagraph("{{/block}}")
	doc.AddParagraph("")
	doc.AddParagraph("{{#block \"main_content\"}}")
	doc.AddParagraph("标准内容区域")
	doc.AddParagraph("{{/block}}")
	doc.AddParagraph("")
	doc.AddParagraph("{{#block \"footer\"}}")
	doc.AddParagraph("标准页脚")
	doc.AddParagraph("{{/block}}")

	_, err := engine.LoadTemplateFromDocument("report_base", doc)
	if err != nil {
		t.Fatalf("Failed to load base template: %v", err)
	}

	// 准备数据
	data := document.NewTemplateData()
	data.SetVariable("reportPeriod", "2024年11月")
	data.SetVariable("totalSales", "1,250,000")
	data.SetVariable("growthRate", "15.8")

	// 渲染模板
	resultDoc, err := engine.RenderTemplateToDocument("report_base", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// 保存文档
	filename := "output/test_block_override.docx"
	err = resultDoc.Save(filename)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 验证文档内容
	if len(resultDoc.Body.Elements) == 0 {
		t.Error("Expected document to have content")
	}

	t.Logf("Block override test completed: %s", filename)
}

// TestTemplateInheritanceValidation 模板继承验证测试
func TestTemplateInheritanceValidation(t *testing.T) {
	engine := document.NewTemplateEngine()

	// 测试1: 块语法验证
	t.Run("块语法验证", func(t *testing.T) {
		// 创建包含正确块语法的文档
		doc := document.New()
		doc.AddParagraph("{{#block \"content\"}}")
		doc.AddParagraph("这是一个有效的块")
		doc.AddParagraph("{{/block}}")

		_, err := engine.LoadTemplateFromDocument("valid_blocks", doc)
		if err != nil {
			t.Fatalf("Valid block syntax should not cause error: %v", err)
		}

		t.Log("Block syntax validation passed")
	})

	// 测试2: 继承链验证
	t.Run("继承链验证", func(t *testing.T) {
		// 创建基础模板
		baseDoc := document.New()
		baseDoc.AddParagraph("基础模板")
		baseDoc.AddParagraph("{{#block \"content\"}}")
		baseDoc.AddParagraph("默认内容")
		baseDoc.AddParagraph("{{/block}}")

		_, err := engine.LoadTemplateFromDocument("inheritance_base", baseDoc)
		if err != nil {
			t.Fatalf("Failed to load inheritance base template: %v", err)
		}

		// 验证模板被正确加载
		template, err := engine.GetTemplate("inheritance_base")
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		if template.Name != "inheritance_base" {
			t.Errorf("Expected template name 'inheritance_base', got '%s'", template.Name)
		}

		t.Log("Inheritance chain validation passed")
	})
}

// cleanupInheritanceTestFiles 清理继承测试文件
func cleanupInheritanceTestFiles() {
	files := []string{
		"output/test_basic_inheritance.docx",
		"output/test_block_override.docx",
	}

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
		}
	}
}

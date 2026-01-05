// Package test 页面设置功能集成测试
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/zerx-lab/wordZero/pkg/document"
)

// TestPageSettingsIntegration 页面设置功能集成测试
func TestPageSettingsIntegration(t *testing.T) {
	// 创建测试文档
	doc := document.New()

	// 测试基本页面设置
	t.Run("基本页面设置", func(t *testing.T) {
		testBasicPageSettings(t, doc)
	})

	// 测试页面尺寸设置
	t.Run("页面尺寸设置", func(t *testing.T) {
		testPageSizes(t, doc)
	})

	// 测试页面方向设置
	t.Run("页面方向设置", func(t *testing.T) {
		testPageOrientation(t, doc)
	})

	// 测试页面边距设置
	t.Run("页面边距设置", func(t *testing.T) {
		testPageMargins(t, doc)
	})

	// 测试自定义页面尺寸
	t.Run("自定义页面尺寸", func(t *testing.T) {
		testCustomPageSize(t, doc)
	})

	// 测试完整文档保存和加载
	t.Run("文档保存和加载", func(t *testing.T) {
		testDocumentSaveLoad(t, doc)
	})
}

// testBasicPageSettings 测试基本页面设置
func testBasicPageSettings(t *testing.T, doc *document.Document) {
	// 获取默认设置
	settings := doc.GetPageSettings()

	if settings.Size != document.PageSizeA4 {
		t.Errorf("默认页面尺寸应为A4，实际为: %s", settings.Size)
	}

	if settings.Orientation != document.OrientationPortrait {
		t.Errorf("默认页面方向应为纵向，实际为: %s", settings.Orientation)
	}

	// 添加测试内容
	doc.AddParagraph("页面设置集成测试 - 基本设置")
}

// testPageSizes 测试页面尺寸设置
func testPageSizes(t *testing.T, doc *document.Document) {
	// 测试各种预定义尺寸
	sizes := []document.PageSize{
		document.PageSizeLetter,
		document.PageSizeLegal,
		document.PageSizeA3,
		document.PageSizeA5,
		document.PageSizeA4, // 恢复到A4
	}

	for _, size := range sizes {
		err := doc.SetPageSize(size)
		if err != nil {
			t.Errorf("设置页面尺寸 %s 失败: %v", size, err)
			continue
		}

		settings := doc.GetPageSettings()
		if settings.Size != size {
			t.Errorf("页面尺寸设置不正确，期望: %s, 实际: %s", size, settings.Size)
		}

		// 添加测试内容
		doc.AddParagraph(fmt.Sprintf("页面尺寸已设置为: %s", size))
	}
}

// testPageOrientation 测试页面方向设置
func testPageOrientation(t *testing.T, doc *document.Document) {
	// 测试横向
	err := doc.SetPageOrientation(document.OrientationLandscape)
	if err != nil {
		t.Errorf("设置横向页面失败: %v", err)
		return
	}

	settings := doc.GetPageSettings()
	if settings.Orientation != document.OrientationLandscape {
		t.Errorf("页面方向应为横向，实际为: %s", settings.Orientation)
	}

	doc.AddParagraph("页面方向已设置为横向")

	// 恢复纵向
	err = doc.SetPageOrientation(document.OrientationPortrait)
	if err != nil {
		t.Errorf("设置纵向页面失败: %v", err)
		return
	}

	settings = doc.GetPageSettings()
	if settings.Orientation != document.OrientationPortrait {
		t.Errorf("页面方向应为纵向，实际为: %s", settings.Orientation)
	}

	doc.AddParagraph("页面方向已恢复为纵向")
}

// testPageMargins 测试页面边距设置
func testPageMargins(t *testing.T, doc *document.Document) {
	// 测试自定义边距
	top, right, bottom, left := 30.0, 20.0, 25.0, 35.0
	err := doc.SetPageMargins(top, right, bottom, left)
	if err != nil {
		t.Errorf("设置页面边距失败: %v", err)
		return
	}

	settings := doc.GetPageSettings()
	if abs(settings.MarginTop-top) > 0.1 {
		t.Errorf("上边距不匹配，期望: %.1fmm, 实际: %.1fmm", top, settings.MarginTop)
	}
	if abs(settings.MarginRight-right) > 0.1 {
		t.Errorf("右边距不匹配，期望: %.1fmm, 实际: %.1fmm", right, settings.MarginRight)
	}
	if abs(settings.MarginBottom-bottom) > 0.1 {
		t.Errorf("下边距不匹配，期望: %.1fmm, 实际: %.1fmm", bottom, settings.MarginBottom)
	}
	if abs(settings.MarginLeft-left) > 0.1 {
		t.Errorf("左边距不匹配，期望: %.1fmm, 实际: %.1fmm", left, settings.MarginLeft)
	}

	doc.AddParagraph("页面边距已设置为自定义值")

	// 测试页眉页脚距离
	header, footer := 15.0, 20.0
	err = doc.SetHeaderFooterDistance(header, footer)
	if err != nil {
		t.Errorf("设置页眉页脚距离失败: %v", err)
		return
	}

	settings = doc.GetPageSettings()
	if abs(settings.HeaderDistance-header) > 0.1 {
		t.Errorf("页眉距离不匹配，期望: %.1fmm, 实际: %.1fmm", header, settings.HeaderDistance)
	}
	if abs(settings.FooterDistance-footer) > 0.1 {
		t.Errorf("页脚距离不匹配，期望: %.1fmm, 实际: %.1fmm", footer, settings.FooterDistance)
	}

	// 测试装订线
	gutter := 8.0
	err = doc.SetGutterWidth(gutter)
	if err != nil {
		t.Errorf("设置装订线宽度失败: %v", err)
		return
	}

	settings = doc.GetPageSettings()
	if abs(settings.GutterWidth-gutter) > 0.1 {
		t.Errorf("装订线宽度不匹配，期望: %.1fmm, 实际: %.1fmm", gutter, settings.GutterWidth)
	}

	doc.AddParagraph("页眉页脚距离和装订线已设置")
}

// testCustomPageSize 测试自定义页面尺寸
func testCustomPageSize(t *testing.T, doc *document.Document) {
	// 测试自定义尺寸
	width, height := 200.0, 250.0
	err := doc.SetCustomPageSize(width, height)
	if err != nil {
		t.Errorf("设置自定义页面尺寸失败: %v", err)
		return
	}

	settings := doc.GetPageSettings()
	if settings.Size != document.PageSizeCustom {
		t.Errorf("页面尺寸应为Custom，实际为: %s", settings.Size)
	}
	if abs(settings.CustomWidth-width) > 0.1 {
		t.Errorf("自定义宽度不匹配，期望: %.1fmm, 实际: %.1fmm", width, settings.CustomWidth)
	}
	if abs(settings.CustomHeight-height) > 0.1 {
		t.Errorf("自定义高度不匹配，期望: %.1fmm, 实际: %.1fmm", height, settings.CustomHeight)
	}

	doc.AddParagraph("页面已设置为自定义尺寸")

	// 恢复到A4
	err = doc.SetPageSize(document.PageSizeA4)
	if err != nil {
		t.Errorf("恢复A4页面尺寸失败: %v", err)
	}
}

// testDocumentSaveLoad 测试文档保存和加载
func testDocumentSaveLoad(t *testing.T, doc *document.Document) {
	// 设置一个完整的页面配置
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
		t.Errorf("设置完整页面配置失败: %v", err)
		return
	}

	// 添加最终测试内容
	doc.AddParagraph("页面设置集成测试完成 - 最终配置")
	doc.AddParagraph("文档将以Letter横向格式保存")

	// 保存文档
	testFile := filepath.Join("testdata", "page_settings_integration_test.docx")

	// 确保测试目录存在
	err = os.MkdirAll(filepath.Dir(testFile), 0755)
	if err != nil {
		t.Errorf("创建测试目录失败: %v", err)
		return
	}

	err = doc.Save(testFile)
	if err != nil {
		t.Errorf("保存测试文档失败: %v", err)
		return
	}

	// 验证文件存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("保存的文档文件不存在: %s", testFile)
		return
	}

	// 重新打开文档验证页面设置
	loadedDoc, err := document.Open(testFile)
	if err != nil {
		t.Errorf("重新打开文档失败: %v", err)
		return
	}

	loadedSettings := loadedDoc.GetPageSettings()

	// 验证关键设置是否保持
	if loadedSettings.Size != settings.Size {
		t.Errorf("加载后页面尺寸不匹配，期望: %s, 实际: %s", settings.Size, loadedSettings.Size)
	}

	if loadedSettings.Orientation != settings.Orientation {
		t.Errorf("加载后页面方向不匹配，期望: %s, 实际: %s", settings.Orientation, loadedSettings.Orientation)
	}

	// 清理测试文件
	os.Remove(testFile)
}

// abs 返回浮点数的绝对值
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

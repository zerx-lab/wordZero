package test

import (
	"testing"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
)

// TestTableDefaultStyle 测试表格默认样式
func TestTableDefaultStyle(t *testing.T) {
	// 创建文档
	doc := document.New()
	if doc == nil {
		t.Fatal("创建文档失败")
	}

	// 创建表格
	config := &document.TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 5000,
		Data: [][]string{
			{"姓名", "年龄", "职业"},
			{"张三", "25", "工程师"},
			{"李四", "30", "设计师"},
		},
	}

	table, err := doc.AddTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 验证表格属性是否设置了默认边框
	if table.Properties == nil {
		t.Fatal("表格属性为空")
	}

	if table.Properties.TableBorders == nil {
		t.Fatal("表格边框未设置")
	}

	// 验证边框样式
	borders := table.Properties.TableBorders
	if borders.Top == nil || borders.Top.Val != "single" {
		t.Error("顶部边框样式不正确")
	}
	if borders.Left == nil || borders.Left.Val != "single" {
		t.Error("左边框样式不正确")
	}
	if borders.Bottom == nil || borders.Bottom.Val != "single" {
		t.Error("底部边框样式不正确")
	}
	if borders.Right == nil || borders.Right.Val != "single" {
		t.Error("右边框样式不正确")
	}
	if borders.InsideH == nil || borders.InsideH.Val != "single" {
		t.Error("内部水平边框样式不正确")
	}
	if borders.InsideV == nil || borders.InsideV.Val != "single" {
		t.Error("内部垂直边框样式不正确")
	}

	// 验证边框粗细
	if borders.Top.Sz != "4" {
		t.Error("边框粗细不正确")
	}

	// 验证表格布局
	if table.Properties.TableLayout == nil || table.Properties.TableLayout.Type != "autofit" {
		t.Error("表格布局设置不正确")
	}

	// 验证单元格边距
	if table.Properties.TableCellMar == nil {
		t.Error("表格单元格边距未设置")
	}

	margins := table.Properties.TableCellMar
	if margins.Left == nil || margins.Left.W != "108" {
		t.Error("左边距设置不正确")
	}
	if margins.Right == nil || margins.Right.W != "108" {
		t.Error("右边距设置不正确")
	}

	t.Log("表格默认样式测试通过")
}

// TestDefaultStyleMatchesTmpTest 测试默认样式是否与tmp_test参考表格匹配
func TestDefaultStyleMatchesTmpTest(t *testing.T) {
	// 创建文档
	doc := document.New()
	if doc == nil {
		t.Fatal("创建文档失败")
	}

	// 创建与tmp_test相同规格的表格
	config := &document.TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 8522, // 与tmp_test中的总宽度匹配
		Data: [][]string{
			{"Cell A1", "Cell B1", "Cell C1"},
			{"Cell A2", "Cell B2", "Cell C2"},
			{"Cell A3", "Cell B3", "Cell C3"},
		},
	}

	table, err := doc.AddTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 验证表格的边框样式与tmp_test参考表格一致
	borders := table.Properties.TableBorders

	// 验证所有边框都是单线样式
	expectedBorderStyle := "single"
	expectedBorderSize := "4"

	if borders.Top.Val != expectedBorderStyle {
		t.Errorf("顶部边框样式不匹配，期望: %s, 实际: %s", expectedBorderStyle, borders.Top.Val)
	}
	if borders.Left.Val != expectedBorderStyle {
		t.Errorf("左边框样式不匹配，期望: %s, 实际: %s", expectedBorderStyle, borders.Left.Val)
	}
	if borders.Bottom.Val != expectedBorderStyle {
		t.Errorf("底部边框样式不匹配，期望: %s, 实际: %s", expectedBorderStyle, borders.Bottom.Val)
	}
	if borders.Right.Val != expectedBorderStyle {
		t.Errorf("右边框样式不匹配，期望: %s, 实际: %s", expectedBorderStyle, borders.Right.Val)
	}
	if borders.InsideH.Val != expectedBorderStyle {
		t.Errorf("内部水平边框样式不匹配，期望: %s, 实际: %s", expectedBorderStyle, borders.InsideH.Val)
	}
	if borders.InsideV.Val != expectedBorderStyle {
		t.Errorf("内部垂直边框样式不匹配，期望: %s, 实际: %s", expectedBorderStyle, borders.InsideV.Val)
	}

	// 验证边框粗细
	if borders.Top.Sz != expectedBorderSize {
		t.Errorf("边框粗细不匹配，期望: %s, 实际: %s", expectedBorderSize, borders.Top.Sz)
	}

	// 验证单元格边距与tmp_test一致
	margins := table.Properties.TableCellMar
	expectedMargin := "108"

	if margins.Left.W != expectedMargin {
		t.Errorf("左边距不匹配，期望: %s, 实际: %s", expectedMargin, margins.Left.W)
	}
	if margins.Right.W != expectedMargin {
		t.Errorf("右边距不匹配，期望: %s, 实际: %s", expectedMargin, margins.Right.W)
	}

	t.Log("默认样式与tmp_test参考表格匹配测试通过")
}

// TestDefaultStyleOverride 测试默认样式可以被覆盖
func TestDefaultStyleOverride(t *testing.T) {
	// 创建文档
	doc := document.New()
	if doc == nil {
		t.Fatal("创建文档失败")
	}

	// 创建表格
	config := &document.TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 3000,
		Data: [][]string{
			{"A1", "B1"},
			{"A2", "B2"},
		},
	}

	table, err := doc.AddTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 验证默认边框已设置
	if table.Properties.TableBorders.Top.Val != "single" {
		t.Error("默认边框样式设置失败")
	}

	// 覆盖默认样式
	borderConfig := &document.TableBorderConfig{
		Top: &document.BorderConfig{
			Style: document.BorderStyleThick,
			Width: 12,
			Color: "FF0000",
			Space: 0,
		},
		Left: &document.BorderConfig{
			Style: document.BorderStyleDouble,
			Width: 8,
			Color: "0000FF",
			Space: 0,
		},
		Bottom: &document.BorderConfig{
			Style: document.BorderStyleDashed,
			Width: 6,
			Color: "00FF00",
			Space: 0,
		},
		Right: &document.BorderConfig{
			Style: document.BorderStyleDotted,
			Width: 4,
			Color: "FF00FF",
			Space: 0,
		},
		InsideH: &document.BorderConfig{
			Style: document.BorderStyleNone,
			Width: 0,
			Color: "auto",
			Space: 0,
		},
		InsideV: &document.BorderConfig{
			Style: document.BorderStyleWave,
			Width: 6,
			Color: "FFFF00",
			Space: 0,
		},
	}

	err = table.SetTableBorders(borderConfig)
	if err != nil {
		t.Fatalf("覆盖边框样式失败: %v", err)
	}

	// 验证样式已被覆盖
	borders := table.Properties.TableBorders
	if borders.Top.Val != "thick" {
		t.Error("顶部边框样式覆盖失败")
	}
	if borders.Left.Val != "double" {
		t.Error("左边框样式覆盖失败")
	}
	if borders.Bottom.Val != "dashed" {
		t.Error("底部边框样式覆盖失败")
	}
	if borders.Right.Val != "dotted" {
		t.Error("右边框样式覆盖失败")
	}
	if borders.InsideH.Val != "none" {
		t.Error("内部水平边框样式覆盖失败")
	}
	if borders.InsideV.Val != "wave" {
		t.Error("内部垂直边框样式覆盖失败")
	}

	t.Log("默认样式覆盖测试通过")
}

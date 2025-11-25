package document

import (
	"testing"
)

// TestCreateTable 测试表格创建功能
func TestCreateTable(t *testing.T) {
	doc := New()

	// 测试基础表格创建
	config := &TableConfig{
		Rows:  3,
		Cols:  4,
		Width: 8000,
		Data: [][]string{
			{"A1", "B1", "C1", "D1"},
			{"A2", "B2", "C2", "D2"},
			{"A3", "B3", "C3", "D3"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 验证表格尺寸
	if table.GetRowCount() != 3 {
		t.Errorf("期望行数3，实际%d", table.GetRowCount())
	}
	if table.GetColumnCount() != 4 {
		t.Errorf("期望列数4，实际%d", table.GetColumnCount())
	}

	// 验证表格内容
	cellText, err := table.GetCellText(0, 0)
	if err != nil {
		t.Errorf("获取单元格内容失败: %v", err)
	}
	if cellText != "A1" {
		t.Errorf("期望单元格内容'A1'，实际'%s'", cellText)
	}

	cellText, err = table.GetCellText(2, 3)
	if err != nil {
		t.Errorf("获取单元格内容失败: %v", err)
	}
	if cellText != "D3" {
		t.Errorf("期望单元格内容'D3'，实际'%s'", cellText)
	}
}

// TestCreateTableWithInvalidConfig 测试无效配置的表格创建
func TestCreateTableWithInvalidConfig(t *testing.T) {
	doc := New()

	// 测试行数为0
	config := &TableConfig{
		Rows:  0,
		Cols:  3,
		Width: 6000,
	}
	_, err := doc.CreateTable(config)
	if err == nil {
		t.Error("期望创建失败，但成功了")
	}

	// 测试列数为0
	config = &TableConfig{
		Rows:  3,
		Cols:  0,
		Width: 6000,
	}
	_, err = doc.CreateTable(config)
	if err == nil {
		t.Error("期望创建失败，但成功了")
	}

	// 测试列宽数量不匹配
	config = &TableConfig{
		Rows:      3,
		Cols:      4,
		Width:     6000,
		ColWidths: []int{1000, 2000}, // 只有2个列宽，但有4列
	}
	_, err = doc.CreateTable(config)
	if err == nil {
		t.Error("期望创建失败，但成功了")
	}
}

// TestAddTable 测试将表格添加到文档
func TestAddTable(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  3,
		Width: 6000,
	}

	initialTableCount := len(doc.Body.GetTables())
	table, err := doc.AddTable(config)

	if err != nil {
		t.Fatalf("添加表格失败: %v", err)
	}

	if len(doc.Body.GetTables()) != initialTableCount+1 {
		t.Errorf("期望表格数量%d，实际%d", initialTableCount+1, len(doc.Body.GetTables()))
	}

	_ = table // 使用table变量避免编译警告
}

// TestTableCellOperations 测试单元格操作
func TestTableCellOperations(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试设置单元格内容
	err = table.SetCellText(0, 0, "测试内容")
	if err != nil {
		t.Errorf("设置单元格内容失败: %v", err)
	}

	// 测试获取单元格内容
	cellText, err := table.GetCellText(0, 0)
	if err != nil {
		t.Errorf("获取单元格内容失败: %v", err)
	}
	if cellText != "测试内容" {
		t.Errorf("期望单元格内容'测试内容'，实际'%s'", cellText)
	}

	// 测试无效的单元格索引
	err = table.SetCellText(5, 5, "无效")
	if err == nil {
		t.Error("期望设置无效单元格失败，但成功了")
	}

	_, err = table.GetCellText(5, 5)
	if err == nil {
		t.Error("期望获取无效单元格失败，但成功了")
	}
}

// TestInsertRow 测试插入行功能
func TestInsertRow(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  3,
		Width: 6000,
		Data: [][]string{
			{"A1", "B1", "C1"},
			{"A2", "B2", "C2"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	initialRowCount := table.GetRowCount()

	// 在中间插入行
	err = table.InsertRow(1, []string{"A1.5", "B1.5", "C1.5"})
	if err != nil {
		t.Errorf("插入行失败: %v", err)
	}

	if table.GetRowCount() != initialRowCount+1 {
		t.Errorf("期望行数%d，实际%d", initialRowCount+1, table.GetRowCount())
	}

	// 验证插入的内容
	cellText, err := table.GetCellText(1, 0)
	if err != nil {
		t.Errorf("获取插入行内容失败: %v", err)
	}
	if cellText != "A1.5" {
		t.Errorf("期望插入行内容'A1.5'，实际'%s'", cellText)
	}

	// 测试在末尾添加行
	err = table.AppendRow([]string{"A末", "B末", "C末"})
	if err != nil {
		t.Errorf("添加行失败: %v", err)
	}

	if table.GetRowCount() != initialRowCount+2 {
		t.Errorf("期望行数%d，实际%d", initialRowCount+2, table.GetRowCount())
	}
}

// TestInsertRowInvalidCases 测试插入行的无效情况
func TestInsertRowInvalidCases(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  3,
		Width: 6000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试无效位置
	err = table.InsertRow(-1, []string{"A", "B", "C"})
	if err == nil {
		t.Error("期望插入无效位置失败，但成功了")
	}

	err = table.InsertRow(10, []string{"A", "B", "C"})
	if err == nil {
		t.Error("期望插入无效位置失败，但成功了")
	}

	// 测试数据列数过多
	err = table.InsertRow(1, []string{"A", "B", "C", "D", "E"})
	if err == nil {
		t.Error("期望插入过多列数据失败，但成功了")
	}
}

// TestDeleteRow 测试删除行功能
func TestDeleteRow(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  4,
		Cols:  3,
		Width: 6000,
		Data: [][]string{
			{"A1", "B1", "C1"},
			{"A2", "B2", "C2"},
			{"A3", "B3", "C3"},
			{"A4", "B4", "C4"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	initialRowCount := table.GetRowCount()

	// 删除第2行（索引1）
	err = table.DeleteRow(1)
	if err != nil {
		t.Errorf("删除行失败: %v", err)
	}

	if table.GetRowCount() != initialRowCount-1 {
		t.Errorf("期望行数%d，实际%d", initialRowCount-1, table.GetRowCount())
	}

	// 验证删除后的内容（原第3行现在应该是第2行）
	cellText, err := table.GetCellText(1, 0)
	if err != nil {
		t.Errorf("获取删除后内容失败: %v", err)
	}
	if cellText != "A3" {
		t.Errorf("期望删除后内容'A3'，实际'%s'", cellText)
	}
}

// TestDeleteRowInvalidCases 测试删除行的无效情况
func TestDeleteRowInvalidCases(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  1,
		Cols:  3,
		Width: 6000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试删除唯一的行
	err = table.DeleteRow(0)
	if err == nil {
		t.Error("期望删除唯一行失败，但成功了")
	}

	// 添加一行以便测试无效索引
	table.AppendRow([]string{"A", "B", "C"})

	// 测试无效索引
	err = table.DeleteRow(-1)
	if err == nil {
		t.Error("期望删除无效索引失败，但成功了")
	}

	err = table.DeleteRow(10)
	if err == nil {
		t.Error("期望删除无效索引失败，但成功了")
	}
}

// TestDeleteRows 测试删除多行功能
func TestDeleteRows(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  5,
		Cols:  2,
		Width: 4000,
		Data: [][]string{
			{"A1", "B1"},
			{"A2", "B2"},
			{"A3", "B3"},
			{"A4", "B4"},
			{"A5", "B5"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	initialRowCount := table.GetRowCount()

	// 删除第2到第4行（索引1到3）
	err = table.DeleteRows(1, 3)
	if err != nil {
		t.Errorf("删除多行失败: %v", err)
	}

	expectedRowCount := initialRowCount - 3
	if table.GetRowCount() != expectedRowCount {
		t.Errorf("期望行数%d，实际%d", expectedRowCount, table.GetRowCount())
	}

	// 验证剩余内容
	cellText, err := table.GetCellText(1, 0)
	if err != nil {
		t.Errorf("获取删除后内容失败: %v", err)
	}
	if cellText != "A5" {
		t.Errorf("期望删除后内容'A5'，实际'%s'", cellText)
	}
}

// TestInsertColumn 测试插入列功能
func TestInsertColumn(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  3,
		Cols:  2,
		Width: 4000,
		Data: [][]string{
			{"A1", "B1"},
			{"A2", "B2"},
			{"A3", "B3"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	initialColCount := table.GetColumnCount()

	// 在中间插入列
	err = table.InsertColumn(1, []string{"C1", "C2", "C3"}, 1000)
	if err != nil {
		t.Errorf("插入列失败: %v", err)
	}

	if table.GetColumnCount() != initialColCount+1 {
		t.Errorf("期望列数%d，实际%d", initialColCount+1, table.GetColumnCount())
	}

	// 验证插入的内容
	cellText, err := table.GetCellText(0, 1)
	if err != nil {
		t.Errorf("获取插入列内容失败: %v", err)
	}
	if cellText != "C1" {
		t.Errorf("期望插入列内容'C1'，实际'%s'", cellText)
	}

	// 测试在末尾添加列
	err = table.AppendColumn([]string{"D1", "D2", "D3"}, 1000)
	if err != nil {
		t.Errorf("添加列失败: %v", err)
	}

	if table.GetColumnCount() != initialColCount+2 {
		t.Errorf("期望列数%d，实际%d", initialColCount+2, table.GetColumnCount())
	}
}

// TestDeleteColumn 测试删除列功能
func TestDeleteColumn(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  3,
		Cols:  4,
		Width: 8000,
		Data: [][]string{
			{"A1", "B1", "C1", "D1"},
			{"A2", "B2", "C2", "D2"},
			{"A3", "B3", "C3", "D3"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	initialColCount := table.GetColumnCount()

	// 删除第2列（索引1）
	err = table.DeleteColumn(1)
	if err != nil {
		t.Errorf("删除列失败: %v", err)
	}

	if table.GetColumnCount() != initialColCount-1 {
		t.Errorf("期望列数%d，实际%d", initialColCount-1, table.GetColumnCount())
	}

	// 验证删除后的内容（原第3列现在应该是第2列）
	cellText, err := table.GetCellText(0, 1)
	if err != nil {
		t.Errorf("获取删除后内容失败: %v", err)
	}
	if cellText != "C1" {
		t.Errorf("期望删除后内容'C1'，实际'%s'", cellText)
	}
}

// TestDeleteColumns 测试删除多列功能
func TestDeleteColumns(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  5,
		Width: 10000,
		Data: [][]string{
			{"A1", "B1", "C1", "D1", "E1"},
			{"A2", "B2", "C2", "D2", "E2"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	initialColCount := table.GetColumnCount()

	// 删除第2到第4列（索引1到3）
	err = table.DeleteColumns(1, 3)
	if err != nil {
		t.Errorf("删除多列失败: %v", err)
	}

	expectedColCount := initialColCount - 3
	if table.GetColumnCount() != expectedColCount {
		t.Errorf("期望列数%d，实际%d", expectedColCount, table.GetColumnCount())
	}

	// 验证剩余内容
	cellText, err := table.GetCellText(0, 1)
	if err != nil {
		t.Errorf("获取删除后内容失败: %v", err)
	}
	if cellText != "E1" {
		t.Errorf("期望删除后内容'E1'，实际'%s'", cellText)
	}
}

// TestClearTable 测试清空表格功能
func TestClearTable(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
		Data: [][]string{
			{"A1", "B1"},
			{"A2", "B2"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 清空表格
	table.ClearTable()

	// 验证所有单元格都为空
	for i := 0; i < table.GetRowCount(); i++ {
		for j := 0; j < table.GetColumnCount(); j++ {
			cellText, err := table.GetCellText(i, j)
			if err != nil {
				t.Errorf("获取清空后单元格内容失败: %v", err)
			}
			if cellText != "" {
				t.Errorf("期望清空后单元格为空，实际'%s'", cellText)
			}
		}
	}

	// 验证表格结构保持不变
	if table.GetRowCount() != 2 {
		t.Errorf("期望清空后行数2，实际%d", table.GetRowCount())
	}
	if table.GetColumnCount() != 2 {
		t.Errorf("期望清空后列数2，实际%d", table.GetColumnCount())
	}
}

// TestCopyTable 测试复制表格功能
func TestCopyTable(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
		Data: [][]string{
			{"原始1", "原始2"},
			{"原始3", "原始4"},
		},
	}

	originalTable, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建原始表格失败: %v", err)
	}

	// 复制表格
	copiedTable := originalTable.CopyTable()
	if copiedTable == nil {
		t.Fatal("复制表格失败")
	}

	// 验证复制的表格结构
	if copiedTable.GetRowCount() != originalTable.GetRowCount() {
		t.Errorf("复制表格行数不匹配：期望%d，实际%d",
			originalTable.GetRowCount(), copiedTable.GetRowCount())
	}
	if copiedTable.GetColumnCount() != originalTable.GetColumnCount() {
		t.Errorf("复制表格列数不匹配：期望%d，实际%d",
			originalTable.GetColumnCount(), copiedTable.GetColumnCount())
	}

	// 验证复制的表格内容
	for i := 0; i < originalTable.GetRowCount(); i++ {
		for j := 0; j < originalTable.GetColumnCount(); j++ {
			originalText, _ := originalTable.GetCellText(i, j)
			copiedText, _ := copiedTable.GetCellText(i, j)
			if originalText != copiedText {
				t.Errorf("复制表格内容不匹配：位置(%d,%d) 期望'%s'，实际'%s'",
					i, j, originalText, copiedText)
			}
		}
	}

	// 修改复制的表格，验证独立性
	err = copiedTable.SetCellText(0, 0, "修改后")
	if err != nil {
		t.Errorf("修改复制表格失败: %v", err)
	}

	originalText, _ := originalTable.GetCellText(0, 0)
	copiedText, _ := copiedTable.GetCellText(0, 0)

	if originalText == copiedText {
		t.Error("复制的表格不是独立的，修改影响了原表格")
	}
}

// TestTableWithCustomColumnWidths 测试自定义列宽的表格
func TestTableWithCustomColumnWidths(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:      2,
		Cols:      3,
		Width:     6000,
		ColWidths: []int{1000, 2000, 3000},
		Data: [][]string{
			{"窄列", "中列", "宽列"},
			{"A", "B", "C"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建自定义列宽表格失败: %v", err)
	}

	// 验证表格创建成功
	if table.GetRowCount() != 2 {
		t.Errorf("期望行数2，实际%d", table.GetRowCount())
	}
	if table.GetColumnCount() != 3 {
		t.Errorf("期望列数3，实际%d", table.GetColumnCount())
	}

	// 验证网格列数量
	if len(table.Grid.Cols) != 3 {
		t.Errorf("期望网格列数3，实际%d", len(table.Grid.Cols))
	}
}

// TestTableElementType 测试表格元素类型
func TestTableElementType(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  1,
		Cols:  1,
		Width: 2000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试表格元素类型
	if table.ElementType() != "table" {
		t.Errorf("期望表格元素类型'table'，实际'%s'", table.ElementType())
	}

	// 测试段落元素类型
	para := doc.AddParagraph("测试段落")
	if para.ElementType() != "paragraph" {
		t.Errorf("期望段落元素类型'paragraph'，实际'%s'", para.ElementType())
	}
}

// TestCellFormattedText 测试单元格富文本功能
func TestCellFormattedText(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试设置富文本内容
	format := &TextFormat{
		Bold:       true,
		Italic:     true,
		FontSize:   14,
		FontColor:  "FF0000",
		FontFamily: "Arial",
	}

	err = table.SetCellFormattedText(0, 0, "富文本测试", format)
	if err != nil {
		t.Errorf("设置富文本内容失败: %v", err)
	}

	// 验证内容
	cellText, err := table.GetCellText(0, 0)
	if err != nil {
		t.Errorf("获取单元格内容失败: %v", err)
	}
	if cellText != "富文本测试" {
		t.Errorf("期望内容'富文本测试'，实际'%s'", cellText)
	}

	// 测试添加格式化文本
	err = table.AddCellFormattedText(0, 0, " 追加文本", &TextFormat{Bold: false, FontColor: "00FF00"})
	if err != nil {
		t.Errorf("添加格式化文本失败: %v", err)
	}
}

// TestCellFormat 测试单元格格式设置
func TestCellFormat(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 设置单元格内容
	err = table.SetCellText(0, 0, "格式测试")
	if err != nil {
		t.Errorf("设置单元格内容失败: %v", err)
	}

	// 测试设置单元格格式
	format := &CellFormat{
		TextFormat: &TextFormat{
			Bold:     true,
			FontSize: 16,
		},
		HorizontalAlign: CellAlignCenter,
		VerticalAlign:   CellVAlignCenter,
	}

	err = table.SetCellFormat(0, 0, format)
	if err != nil {
		t.Errorf("设置单元格格式失败: %v", err)
	}

	// 获取并验证格式
	retrievedFormat, err := table.GetCellFormat(0, 0)
	if err != nil {
		t.Errorf("获取单元格格式失败: %v", err)
	}

	if retrievedFormat.HorizontalAlign != CellAlignCenter {
		t.Errorf("期望水平对齐'center'，实际'%s'", retrievedFormat.HorizontalAlign)
	}

	if retrievedFormat.VerticalAlign != CellVAlignCenter {
		t.Errorf("期望垂直对齐'center'，实际'%s'", retrievedFormat.VerticalAlign)
	}

	if retrievedFormat.TextFormat == nil || !retrievedFormat.TextFormat.Bold {
		t.Error("期望文字格式为粗体")
	}
}

// TestCellMergeHorizontal 测试水平合并单元格
func TestCellMergeHorizontal(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  3,
		Cols:  4,
		Width: 8000,
		Data: [][]string{
			{"A1", "B1", "C1", "D1"},
			{"A2", "B2", "C2", "D2"},
			{"A3", "B3", "C3", "D3"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	initialColCount := table.GetColumnCount()

	// 合并第一行的第2到第4列（索引1到3）
	err = table.MergeCellsHorizontal(0, 1, 3)
	if err != nil {
		t.Errorf("水平合并单元格失败: %v", err)
	}

	// 验证合并后第一行的列数减少
	if len(table.Rows[0].Cells) != initialColCount-2 {
		t.Errorf("期望第一行列数%d，实际%d", initialColCount-2, len(table.Rows[0].Cells))
	}

	// 验证合并状态
	isMerged, err := table.IsCellMerged(0, 1)
	if err != nil {
		t.Errorf("检查合并状态失败: %v", err)
	}
	if !isMerged {
		t.Error("期望单元格已合并")
	}

	// 获取合并信息
	mergeInfo, err := table.GetMergedCellInfo(0, 1)
	if err != nil {
		t.Errorf("获取合并信息失败: %v", err)
	}

	if !mergeInfo["is_merged"].(bool) {
		t.Error("期望单元格处于合并状态")
	}

	if mergeInfo["horizontal_span"].(int) != 3 {
		t.Errorf("期望水平跨度3，实际%d", mergeInfo["horizontal_span"].(int))
	}
}

// TestCellMergeVertical 测试垂直合并单元格
func TestCellMergeVertical(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  4,
		Cols:  3,
		Width: 6000,
		Data: [][]string{
			{"A1", "B1", "C1"},
			{"A2", "B2", "C2"},
			{"A3", "B3", "C3"},
			{"A4", "B4", "C4"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 合并第2列的第1到第3行（索引0到2）
	err = table.MergeCellsVertical(0, 2, 1)
	if err != nil {
		t.Errorf("垂直合并单元格失败: %v", err)
	}

	// 验证合并状态
	isMerged, err := table.IsCellMerged(0, 1)
	if err != nil {
		t.Errorf("检查合并状态失败: %v", err)
	}
	if !isMerged {
		t.Error("期望单元格已合并")
	}

	// 验证被合并的单元格也有合并标记
	isMerged, err = table.IsCellMerged(1, 1)
	if err != nil {
		t.Errorf("检查合并状态失败: %v", err)
	}
	if !isMerged {
		t.Error("期望被合并单元格也有合并标记")
	}
}

// TestCellMergeRange 测试区域合并
func TestCellMergeRange(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  4,
		Cols:  4,
		Width: 8000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 合并2x2区域（行0-1，列1-2）
	err = table.MergeCellsRange(0, 1, 1, 2)
	if err != nil {
		t.Errorf("区域合并失败: %v", err)
	}

	// 验证合并状态
	isMerged, err := table.IsCellMerged(0, 1)
	if err != nil {
		t.Errorf("检查合并状态失败: %v", err)
	}
	if !isMerged {
		t.Error("期望单元格已合并")
	}
}

// TestUnmergeCells 测试取消合并
func TestUnmergeCells(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 6000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 先进行水平合并
	err = table.MergeCellsHorizontal(0, 0, 1)
	if err != nil {
		t.Errorf("水平合并失败: %v", err)
	}

	// 验证合并状态
	isMerged, err := table.IsCellMerged(0, 0)
	if err != nil {
		t.Errorf("检查合并状态失败: %v", err)
	}
	if !isMerged {
		t.Error("期望单元格已合并")
	}

	// 取消合并
	err = table.UnmergeCells(0, 0)
	if err != nil {
		t.Errorf("取消合并失败: %v", err)
	}

	// 验证取消合并后的状态
	isMerged, err = table.IsCellMerged(0, 0)
	if err != nil {
		t.Errorf("检查合并状态失败: %v", err)
	}
	if isMerged {
		t.Error("期望单元格已取消合并")
	}
}

// TestCellContentOperations 测试单元格内容操作
func TestCellContentOperations(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 设置格式化内容
	format := &TextFormat{
		Bold:     true,
		FontSize: 12,
	}
	err = table.SetCellFormattedText(0, 0, "测试内容", format)
	if err != nil {
		t.Errorf("设置格式化内容失败: %v", err)
	}

	// 清空内容但保留格式
	err = table.ClearCellContent(0, 0)
	if err != nil {
		t.Errorf("清空单元格内容失败: %v", err)
	}

	// 验证内容已清空
	content, err := table.GetCellText(0, 0)
	if err != nil {
		t.Errorf("获取单元格内容失败: %v", err)
	}
	if content != "" {
		t.Errorf("期望内容为空，实际'%s'", content)
	}

	// 重新设置内容
	err = table.SetCellText(0, 0, "新内容")
	if err != nil {
		t.Errorf("设置新内容失败: %v", err)
	}

	// 清空格式但保留内容
	err = table.ClearCellFormat(0, 0)
	if err != nil {
		t.Errorf("清空单元格格式失败: %v", err)
	}

	// 验证内容保留
	content, err = table.GetCellText(0, 0)
	if err != nil {
		t.Errorf("获取单元格内容失败: %v", err)
	}
	if content != "新内容" {
		t.Errorf("期望内容'新内容'，实际'%s'", content)
	}
}

// TestCellMergeInvalidCases 测试合并的无效情况
func TestCellMergeInvalidCases(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试无效的水平合并
	err = table.MergeCellsHorizontal(0, 0, 0)
	if err == nil {
		t.Error("期望相同列合并失败，但成功了")
	}

	err = table.MergeCellsHorizontal(0, 1, 0)
	if err == nil {
		t.Error("期望反向合并失败，但成功了")
	}

	// 测试无效的垂直合并
	err = table.MergeCellsVertical(0, 0, 0)
	if err == nil {
		t.Error("期望相同行合并失败，但成功了")
	}

	err = table.MergeCellsVertical(1, 0, 0)
	if err == nil {
		t.Error("期望反向合并失败，但成功了")
	}

	// 测试无效的索引
	err = table.MergeCellsHorizontal(-1, 0, 1)
	if err == nil {
		t.Error("期望无效行索引失败，但成功了")
	}

	err = table.MergeCellsVertical(0, 1, -1)
	if err == nil {
		t.Error("期望无效列索引失败，但成功了")
	}
}

// TestCellPadding 测试单元格内边距
func TestCellPadding(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  1,
		Cols:  1,
		Width: 2000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试设置内边距
	err = table.SetCellPadding(0, 0, 10)
	if err != nil {
		t.Errorf("设置单元格内边距失败: %v", err)
	}

	// 测试无效索引
	err = table.SetCellPadding(5, 5, 10)
	if err == nil {
		t.Error("期望无效索引失败，但成功了")
	}
}

// TestCellTextDirection 测试单元格文字方向设置
func TestCellTextDirection(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 6000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试设置不同的文字方向
	testCases := []struct {
		name      string
		direction CellTextDirection
		row       int
		col       int
	}{
		{"从左到右", TextDirectionLR, 0, 0},
		{"从上到下", TextDirectionTB, 0, 1},
		{"从下到上", TextDirectionBT, 0, 2},
		{"从右到左", TextDirectionRL, 1, 0},
		{"从上到下垂直", TextDirectionTBV, 1, 1},
		{"从下到上垂直", TextDirectionBTV, 1, 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置文字内容
			err := table.SetCellText(tc.row, tc.col, tc.name)
			if err != nil {
				t.Errorf("设置单元格文本失败: %v", err)
			}

			// 设置文字方向
			err = table.SetCellTextDirection(tc.row, tc.col, tc.direction)
			if err != nil {
				t.Errorf("设置文字方向失败: %v", err)
			}

			// 验证文字方向
			actualDirection, err := table.GetCellTextDirection(tc.row, tc.col)
			if err != nil {
				t.Errorf("获取文字方向失败: %v", err)
			}

			if actualDirection != tc.direction {
				t.Errorf("文字方向不匹配，期望: %s，实际: %s", tc.direction, actualDirection)
			}
		})
	}
}

// TestCellFormatWithTextDirection 测试通过CellFormat设置文字方向
func TestCellFormatWithTextDirection(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 通过CellFormat设置完整格式，包括文字方向
	format := &CellFormat{
		TextFormat: &TextFormat{
			Bold:     true,
			FontSize: 14,
		},
		HorizontalAlign: CellAlignCenter,
		VerticalAlign:   CellVAlignCenter,
		TextDirection:   TextDirectionTB, // 从上到下
	}

	err = table.SetCellText(0, 0, "竖排文字")
	if err != nil {
		t.Errorf("设置单元格文本失败: %v", err)
	}

	err = table.SetCellFormat(0, 0, format)
	if err != nil {
		t.Errorf("设置单元格格式失败: %v", err)
	}

	// 验证格式是否正确设置
	retrievedFormat, err := table.GetCellFormat(0, 0)
	if err != nil {
		t.Errorf("获取单元格格式失败: %v", err)
	}

	if retrievedFormat.TextDirection != TextDirectionTB {
		t.Errorf("文字方向不匹配，期望: %s，实际: %s", TextDirectionTB, retrievedFormat.TextDirection)
	}

	if retrievedFormat.HorizontalAlign != CellAlignCenter {
		t.Errorf("水平对齐不匹配，期望: %s，实际: %s", CellAlignCenter, retrievedFormat.HorizontalAlign)
	}

	if retrievedFormat.VerticalAlign != CellVAlignCenter {
		t.Errorf("垂直对齐不匹配，期望: %s，实际: %s", CellVAlignCenter, retrievedFormat.VerticalAlign)
	}
}

// TestTextDirectionConstants 测试文字方向常量
func TestTextDirectionConstants(t *testing.T) {
	directions := []CellTextDirection{
		TextDirectionLR,
		TextDirectionTB,
		TextDirectionBT,
		TextDirectionRL,
		TextDirectionTBV,
		TextDirectionBTV,
	}

	expectedValues := []string{
		"lrTb",
		"tbRl",
		"btLr",
		"rlTb",
		"tbLrV",
		"btLrV",
	}

	for i, direction := range directions {
		if string(direction) != expectedValues[i] {
			t.Errorf("文字方向常量值不匹配，期望: %s，实际: %s", expectedValues[i], string(direction))
		}
	}
}

// TestRowHeight 测试行高设置功能
func TestRowHeight(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  3,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试设置固定行高
	heightConfig := &RowHeightConfig{
		Height: 30,
		Rule:   RowHeightExact,
	}

	err = table.SetRowHeight(0, heightConfig)
	if err != nil {
		t.Errorf("设置行高失败: %v", err)
	}

	// 测试获取行高
	retrievedConfig, err := table.GetRowHeight(0)
	if err != nil {
		t.Errorf("获取行高失败: %v", err)
	}

	if retrievedConfig.Height != 30 {
		t.Errorf("期望行高30，实际%d", retrievedConfig.Height)
	}

	if retrievedConfig.Rule != RowHeightExact {
		t.Errorf("期望行高规则%s，实际%s", RowHeightExact, retrievedConfig.Rule)
	}

	// 测试批量设置行高
	batchConfig := &RowHeightConfig{
		Height: 25,
		Rule:   RowHeightMinimum,
	}

	err = table.SetRowHeightRange(1, 2, batchConfig)
	if err != nil {
		t.Errorf("批量设置行高失败: %v", err)
	}

	// 验证批量设置结果
	for i := 1; i <= 2; i++ {
		config, err := table.GetRowHeight(i)
		if err != nil {
			t.Errorf("获取第%d行高度失败: %v", i, err)
		}
		if config.Height != 25 {
			t.Errorf("第%d行期望高度25，实际%d", i, config.Height)
		}
		if config.Rule != RowHeightMinimum {
			t.Errorf("第%d行期望规则%s，实际%s", i, RowHeightMinimum, config.Rule)
		}
	}

	// 测试无效索引
	err = table.SetRowHeight(10, heightConfig)
	if err == nil {
		t.Error("期望设置无效行索引失败，但成功了")
	}

	_, err = table.GetRowHeight(10)
	if err == nil {
		t.Error("期望获取无效行索引失败，但成功了")
	}
}

// TestTableLayout 测试表格布局和定位功能
func TestTableLayout(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试设置表格布局
	layoutConfig := &TableLayoutConfig{
		Alignment: TableAlignCenter,
		TextWrap:  TextWrapNone,
		Position:  PositionInline,
	}

	err = table.SetTableLayout(layoutConfig)
	if err != nil {
		t.Errorf("设置表格布局失败: %v", err)
	}

	// 测试获取表格布局
	retrievedLayout := table.GetTableLayout()
	if retrievedLayout.Alignment != TableAlignCenter {
		t.Errorf("期望对齐方式%s，实际%s", TableAlignCenter, retrievedLayout.Alignment)
	}

	// 测试快捷方法设置对齐
	err = table.SetTableAlignment(TableAlignRight)
	if err != nil {
		t.Errorf("设置表格对齐失败: %v", err)
	}

	retrievedLayout = table.GetTableLayout()
	if retrievedLayout.Alignment != TableAlignRight {
		t.Errorf("期望对齐方式%s，实际%s", TableAlignRight, retrievedLayout.Alignment)
	}
}

// TestTablePageBreak 测试表格分页控制功能
func TestTablePageBreak(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  4,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试设置行禁止跨页分割
	err = table.SetRowKeepTogether(0, true)
	if err != nil {
		t.Errorf("设置行禁止跨页分割失败: %v", err)
	}

	// 测试检查行是否禁止跨页分割
	keepTogether, err := table.IsRowKeepTogether(0)
	if err != nil {
		t.Errorf("检查行跨页分割设置失败: %v", err)
	}
	if !keepTogether {
		t.Error("期望行禁止跨页分割为true，实际为false")
	}

	// 测试设置标题行
	err = table.SetRowAsHeader(0, true)
	if err != nil {
		t.Errorf("设置标题行失败: %v", err)
	}

	// 测试检查是否为标题行
	isHeader, err := table.IsRowHeader(0)
	if err != nil {
		t.Errorf("检查标题行设置失败: %v", err)
	}
	if !isHeader {
		t.Error("期望第0行为标题行，实际不是")
	}

	// 测试设置标题行范围
	err = table.SetHeaderRows(0, 1)
	if err != nil {
		t.Errorf("设置标题行范围失败: %v", err)
	}

	// 验证标题行范围设置
	for i := 0; i <= 1; i++ {
		isHeader, err := table.IsRowHeader(i)
		if err != nil {
			t.Errorf("检查第%d行标题行设置失败: %v", i, err)
		}
		if !isHeader {
			t.Errorf("期望第%d行为标题行，实际不是", i)
		}
	}

	// 测试表格分页信息
	breakInfo := table.GetTableBreakInfo()
	if breakInfo["total_rows"] != 4 {
		t.Errorf("期望总行数4，实际%v", breakInfo["total_rows"])
	}
	if breakInfo["header_rows"] != 2 {
		t.Errorf("期望标题行数2，实际%v", breakInfo["header_rows"])
	}

	// 测试表格分页配置
	pageBreakConfig := &TablePageBreakConfig{
		KeepWithNext:    true,
		KeepLines:       true,
		PageBreakBefore: false,
		WidowControl:    true,
	}

	err = table.SetTablePageBreak(pageBreakConfig)
	if err != nil {
		t.Errorf("设置表格分页配置失败: %v", err)
	}

	// 测试行与下一行保持在同一页
	err = table.SetRowKeepWithNext(1, true)
	if err != nil {
		t.Errorf("设置行与下一行保持在同一页失败: %v", err)
	}

	// 测试无效索引
	err = table.SetRowKeepTogether(10, true)
	if err == nil {
		t.Error("期望设置无效行索引失败，但成功了")
	}

	err = table.SetRowAsHeader(10, true)
	if err == nil {
		t.Error("期望设置无效行索引失败，但成功了")
	}

	_, err = table.IsRowHeader(10)
	if err == nil {
		t.Error("期望检查无效行索引失败，但成功了")
	}

	_, err = table.IsRowKeepTogether(10)
	if err == nil {
		t.Error("期望检查无效行索引失败，但成功了")
	}
}

// TestRowHeightConstants 测试行高规则常量
func TestRowHeightConstants(t *testing.T) {
	// 验证行高规则常量定义正确
	if RowHeightAuto != "auto" {
		t.Errorf("期望RowHeightAuto为'auto'，实际'%s'", RowHeightAuto)
	}
	if RowHeightMinimum != "atLeast" {
		t.Errorf("期望RowHeightMinimum为'atLeast'，实际'%s'", RowHeightMinimum)
	}
	if RowHeightExact != "exact" {
		t.Errorf("期望RowHeightExact为'exact'，实际'%s'", RowHeightExact)
	}
}

// TestTableAlignmentConstants 测试表格对齐常量
func TestTableAlignmentConstants(t *testing.T) {
	// 验证表格对齐常量定义正确
	if TableAlignLeft != "left" {
		t.Errorf("期望TableAlignLeft为'left'，实际'%s'", TableAlignLeft)
	}
	if TableAlignCenter != "center" {
		t.Errorf("期望TableAlignCenter为'center'，实际'%s'", TableAlignCenter)
	}
	if TableAlignRight != "right" {
		t.Errorf("期望TableAlignRight为'right'，实际'%s'", TableAlignRight)
	}
}

// TestTableRowPropertiesExtensions 测试TableRowProperties扩展方法
func TestTableRowPropertiesExtensions(t *testing.T) {
	trp := &TableRowProperties{}

	// 测试SetCantSplit方法
	trp.SetCantSplit(true)
	if trp.CantSplit == nil || trp.CantSplit.Val != "1" {
		t.Error("设置CantSplit失败")
	}

	trp.SetCantSplit(false)
	if trp.CantSplit != nil {
		t.Error("清除CantSplit失败")
	}

	// 测试SetTblHeader方法
	trp.SetTblHeader(true)
	if trp.TblHeader == nil || trp.TblHeader.Val != "1" {
		t.Error("设置TblHeader失败")
	}

	trp.SetTblHeader(false)
	if trp.TblHeader != nil {
		t.Error("清除TblHeader失败")
	}
}

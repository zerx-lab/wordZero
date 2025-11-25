package document

import (
	"fmt"
	"testing"
)

// TestCellIterator 测试基本的单元格迭代器功能
func TestCellIterator(t *testing.T) {
	// 创建一个3x3的测试表格
	doc := New()
	config := &TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 5000,
		Data: [][]string{
			{"A1", "B1", "C1"},
			{"A2", "B2", "C2"},
			{"A3", "B3", "C3"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试迭代器创建
	iterator := table.NewCellIterator()
	if iterator == nil {
		t.Fatal("创建迭代器失败")
	}

	// 测试Total方法
	expectedTotal := 9
	if iterator.Total() != expectedTotal {
		t.Errorf("Total()期望返回%d，实际返回%d", expectedTotal, iterator.Total())
	}

	// 测试迭代器遍历
	cellCount := 0
	expectedCells := []struct {
		row  int
		col  int
		text string
	}{
		{0, 0, "A1"}, {0, 1, "B1"}, {0, 2, "C1"},
		{1, 0, "A2"}, {1, 1, "B2"}, {1, 2, "C2"},
		{2, 0, "A3"}, {2, 1, "B3"}, {2, 2, "C3"},
	}

	for iterator.HasNext() {
		cellInfo, err := iterator.Next()
		if err != nil {
			t.Fatalf("迭代器Next()失败: %v", err)
		}

		if cellCount >= len(expectedCells) {
			t.Fatalf("迭代器返回了过多的单元格")
		}

		expected := expectedCells[cellCount]
		if cellInfo.Row != expected.row || cellInfo.Col != expected.col {
			t.Errorf("单元格位置不匹配: 期望(%d,%d)，实际(%d,%d)",
				expected.row, expected.col, cellInfo.Row, cellInfo.Col)
		}

		if cellInfo.Text != expected.text {
			t.Errorf("单元格文本不匹配: 期望'%s'，实际'%s'",
				expected.text, cellInfo.Text)
		}

		if cellInfo.Cell == nil {
			t.Error("单元格引用为nil")
		}

		// 测试IsLast标记
		if cellCount == len(expectedCells)-1 && !cellInfo.IsLast {
			t.Error("最后一个单元格的IsLast标记应为true")
		}

		cellCount++
	}

	if cellCount != expectedTotal {
		t.Errorf("迭代的单元格数量不匹配: 期望%d，实际%d", expectedTotal, cellCount)
	}
}

// TestCellIteratorReset 测试迭代器重置功能
func TestCellIteratorReset(t *testing.T) {
	doc := New()
	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 3000,
		Data: [][]string{
			{"A1", "B1"},
			{"A2", "B2"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}
	iterator := table.NewCellIterator()

	// 先迭代一些单元格
	iterator.Next()
	iterator.Next()

	// 检查当前位置
	row, col := iterator.Current()
	if row != 1 || col != 0 {
		t.Errorf("迭代器位置不正确: 期望(1,0)，实际(%d,%d)", row, col)
	}

	// 重置迭代器
	iterator.Reset()

	// 检查重置后的位置
	row, col = iterator.Current()
	if row != 0 || col != 0 {
		t.Errorf("重置后位置不正确: 期望(0,0)，实际(%d,%d)", row, col)
	}

	// 确保能重新遍历
	if !iterator.HasNext() {
		t.Error("重置后应该有下一个单元格")
	}
}

// TestCellIteratorProgress 测试进度计算
func TestCellIteratorProgress(t *testing.T) {
	doc := New()
	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 3000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}
	iterator := table.NewCellIterator()

	// 初始进度应为0
	if iterator.Progress() != 0.0 {
		t.Errorf("初始进度应为0.0，实际为%f", iterator.Progress())
	}

	// 迭代一个单元格
	iterator.Next()
	expectedProgress := 0.25 // 1/4
	if iterator.Progress() != expectedProgress {
		t.Errorf("迭代一个单元格后进度应为%f，实际为%f", expectedProgress, iterator.Progress())
	}

	// 迭代到最后
	for iterator.HasNext() {
		iterator.Next()
	}

	if iterator.Progress() != 1.0 {
		t.Errorf("完成迭代后进度应为1.0，实际为%f", iterator.Progress())
	}
}

// TestTableForEach 测试ForEach方法
func TestTableForEach(t *testing.T) {
	doc := New()
	config := &TableConfig{
		Rows:  2,
		Cols:  3,
		Width: 4000,
		Data: [][]string{
			{"A1", "B1", "C1"},
			{"A2", "B2", "C2"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试ForEach遍历
	var visitedCells []string
	err = table.ForEach(func(row, col int, cell *TableCell, text string) error {
		visitedCells = append(visitedCells, fmt.Sprintf("%d-%d:%s", row, col, text))
		return nil
	})

	if err != nil {
		t.Fatalf("ForEach执行失败: %v", err)
	}

	expectedCells := []string{
		"0-0:A1", "0-1:B1", "0-2:C1",
		"1-0:A2", "1-1:B2", "1-2:C2",
	}

	if len(visitedCells) != len(expectedCells) {
		t.Errorf("访问的单元格数量不匹配: 期望%d，实际%d", len(expectedCells), len(visitedCells))
	}

	for i, expected := range expectedCells {
		if i < len(visitedCells) && visitedCells[i] != expected {
			t.Errorf("单元格访问顺序不正确: 期望'%s'，实际'%s'", expected, visitedCells[i])
		}
	}
}

// TestForEachInRow 测试按行遍历
func TestForEachInRow(t *testing.T) {
	doc := New()
	config := &TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 4000,
		Data: [][]string{
			{"A1", "B1", "C1"},
			{"A2", "B2", "C2"},
			{"A3", "B3", "C3"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试遍历第2行（索引1）
	var visitedCells []string
	err = table.ForEachInRow(1, func(col int, cell *TableCell, text string) error {
		visitedCells = append(visitedCells, fmt.Sprintf("%d:%s", col, text))
		return nil
	})

	if err != nil {
		t.Fatalf("ForEachInRow执行失败: %v", err)
	}

	expectedCells := []string{"0:A2", "1:B2", "2:C2"}
	if len(visitedCells) != len(expectedCells) {
		t.Errorf("访问的单元格数量不匹配: 期望%d，实际%d", len(expectedCells), len(visitedCells))
	}

	for i, expected := range expectedCells {
		if i < len(visitedCells) && visitedCells[i] != expected {
			t.Errorf("单元格访问顺序不正确: 期望'%s'，实际'%s'", expected, visitedCells[i])
		}
	}

	// 测试无效行索引
	err = table.ForEachInRow(5, func(col int, cell *TableCell, text string) error {
		return nil
	})
	if err == nil {
		t.Error("应该返回无效行索引错误")
	}
}

// TestForEachInColumn 测试按列遍历
func TestForEachInColumn(t *testing.T) {
	doc := New()
	config := &TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 4000,
		Data: [][]string{
			{"A1", "B1", "C1"},
			{"A2", "B2", "C2"},
			{"A3", "B3", "C3"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试遍历第2列（索引1）
	var visitedCells []string
	err = table.ForEachInColumn(1, func(row int, cell *TableCell, text string) error {
		visitedCells = append(visitedCells, fmt.Sprintf("%d:%s", row, text))
		return nil
	})

	if err != nil {
		t.Fatalf("ForEachInColumn执行失败: %v", err)
	}

	expectedCells := []string{"0:B1", "1:B2", "2:B3"}
	if len(visitedCells) != len(expectedCells) {
		t.Errorf("访问的单元格数量不匹配: 期望%d，实际%d", len(expectedCells), len(visitedCells))
	}

	for i, expected := range expectedCells {
		if i < len(visitedCells) && visitedCells[i] != expected {
			t.Errorf("单元格访问顺序不正确: 期望'%s'，实际'%s'", expected, visitedCells[i])
		}
	}

	// 测试无效列索引
	err = table.ForEachInColumn(5, func(row int, cell *TableCell, text string) error {
		return nil
	})
	if err == nil {
		t.Error("应该返回无效列索引错误")
	}
}

// TestGetCellRange 测试获取单元格范围
func TestGetCellRange(t *testing.T) {
	doc := New()
	config := &TableConfig{
		Rows:  4,
		Cols:  4,
		Width: 5000,
		Data: [][]string{
			{"A1", "B1", "C1", "D1"},
			{"A2", "B2", "C2", "D2"},
			{"A3", "B3", "C3", "D3"},
			{"A4", "B4", "C4", "D4"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试获取2x2范围 (1,1) 到 (2,2)
	cells, err := table.GetCellRange(1, 1, 2, 2)
	if err != nil {
		t.Fatalf("GetCellRange执行失败: %v", err)
	}

	expectedCells := []struct {
		row  int
		col  int
		text string
	}{
		{1, 1, "B2"}, {1, 2, "C2"},
		{2, 1, "B3"}, {2, 2, "C3"},
	}

	if len(cells) != len(expectedCells) {
		t.Errorf("返回的单元格数量不匹配: 期望%d，实际%d", len(expectedCells), len(cells))
	}

	for i, expected := range expectedCells {
		if i < len(cells) {
			cell := cells[i]
			if cell.Row != expected.row || cell.Col != expected.col {
				t.Errorf("单元格位置不匹配: 期望(%d,%d)，实际(%d,%d)",
					expected.row, expected.col, cell.Row, cell.Col)
			}
			if cell.Text != expected.text {
				t.Errorf("单元格文本不匹配: 期望'%s'，实际'%s'",
					expected.text, cell.Text)
			}
		}
	}

	// 测试无效范围
	_, err = table.GetCellRange(2, 2, 1, 1) // 开始位置大于结束位置
	if err == nil {
		t.Error("应该返回无效范围错误")
	}

	_, err = table.GetCellRange(0, 0, 10, 10) // 超出表格范围
	if err == nil {
		t.Error("应该返回超出范围错误")
	}
}

// TestFindCells 测试查找单元格功能
func TestFindCells(t *testing.T) {
	doc := New()
	config := &TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 4000,
		Data: [][]string{
			{"apple", "banana", "cherry"},
			{"dog", "apple", "fish"},
			{"grape", "horse", "apple"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 查找包含"apple"的单元格
	cells, err := table.FindCells(func(row, col int, cell *TableCell, text string) bool {
		return text == "apple"
	})

	if err != nil {
		t.Fatalf("FindCells执行失败: %v", err)
	}

	expectedPositions := [][2]int{{0, 0}, {1, 1}, {2, 2}}
	if len(cells) != len(expectedPositions) {
		t.Errorf("找到的单元格数量不匹配: 期望%d，实际%d", len(expectedPositions), len(cells))
	}

	for i, expected := range expectedPositions {
		if i < len(cells) {
			cell := cells[i]
			if cell.Row != expected[0] || cell.Col != expected[1] {
				t.Errorf("找到的单元格位置不正确: 期望(%d,%d)，实际(%d,%d)",
					expected[0], expected[1], cell.Row, cell.Col)
			}
			if cell.Text != "apple" {
				t.Errorf("找到的单元格文本不正确: 期望'apple'，实际'%s'", cell.Text)
			}
		}
	}
}

// TestFindCellsByText 测试按文本查找单元格
func TestFindCellsByText(t *testing.T) {
	doc := New()
	config := &TableConfig{
		Rows:  2,
		Cols:  3,
		Width: 4000,
		Data: [][]string{
			{"test", "testing", "other"},
			{"notest", "test123", "test"},
		},
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 精确匹配
	cells, err := table.FindCellsByText("test", true)
	if err != nil {
		t.Fatalf("FindCellsByText执行失败: %v", err)
	}

	expectedCount := 2 // (0,0) 和 (1,2)
	if len(cells) != expectedCount {
		t.Errorf("精确匹配找到的单元格数量不匹配: 期望%d，实际%d", expectedCount, len(cells))
	}

	// 模糊匹配
	cells, err = table.FindCellsByText("test", false)
	if err != nil {
		t.Fatalf("FindCellsByText执行失败: %v", err)
	}

	expectedCount = 5 // 所有包含"test"的单元格
	if len(cells) != expectedCount {
		t.Errorf("模糊匹配找到的单元格数量不匹配: 期望%d，实际%d", expectedCount, len(cells))
	}
}

// TestEmptyTable 测试空表格的迭代器
func TestEmptyTable(t *testing.T) {
	doc := New()

	// 创建一个空表格（实际上最小1x1）
	config := &TableConfig{
		Rows:  1,
		Cols:  1,
		Width: 2000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}
	iterator := table.NewCellIterator()

	if iterator.Total() != 1 {
		t.Errorf("空表格的总单元格数应为1，实际为%d", iterator.Total())
	}

	if !iterator.HasNext() {
		t.Error("1x1表格应该有一个单元格")
	}
}

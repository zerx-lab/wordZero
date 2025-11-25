package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
)

// TestCellIteratorIntegration 单元格迭代器功能集成测试
func TestCellIteratorIntegration(t *testing.T) {
	// 创建测试文档
	doc := document.New()

	// 创建测试表格
	config := &document.TableConfig{
		Rows:  4,
		Cols:  3,
		Width: 6000,
		Data: [][]string{
			{"姓名", "年龄", "城市"},
			{"张三", "25", "北京"},
			{"李四", "30", "上海"},
			{"王五", "28", "广州"},
		},
	}

	table, err := doc.AddTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试1: 基本迭代器功能
	t.Run("基本迭代器", func(t *testing.T) {
		iterator := table.NewCellIterator()

		// 验证总数
		expectedTotal := 12
		if iterator.Total() != expectedTotal {
			t.Errorf("总单元格数不正确: 期望%d，实际%d", expectedTotal, iterator.Total())
		}

		// 验证完整遍历
		count := 0
		for iterator.HasNext() {
			cellInfo, err := iterator.Next()
			if err != nil {
				t.Errorf("迭代器错误: %v", err)
				break
			}

			if cellInfo == nil {
				t.Error("单元格信息为空")
				continue
			}

			if cellInfo.Cell == nil {
				t.Error("单元格引用为空")
			}

			count++
		}

		if count != expectedTotal {
			t.Errorf("实际遍历数量不正确: 期望%d，实际%d", expectedTotal, count)
		}
	})

	// 测试2: 重置功能
	t.Run("迭代器重置", func(t *testing.T) {
		iterator := table.NewCellIterator()

		// 迭代几个单元格
		iterator.Next()
		iterator.Next()

		// 重置并验证
		iterator.Reset()
		row, col := iterator.Current()
		if row != 0 || col != 0 {
			t.Errorf("重置后位置不正确: 期望(0,0)，实际(%d,%d)", row, col)
		}
	})

	// 测试3: ForEach功能
	t.Run("ForEach遍历", func(t *testing.T) {
		visitedCount := 0
		err := table.ForEach(func(row, col int, cell *document.TableCell, text string) error {
			visitedCount++
			if cell == nil {
				t.Errorf("位置(%d,%d)的单元格为空", row, col)
			}
			return nil
		})

		if err != nil {
			t.Errorf("ForEach执行失败: %v", err)
		}

		if visitedCount != 12 {
			t.Errorf("ForEach访问数量不正确: 期望12，实际%d", visitedCount)
		}
	})

	// 测试4: 按行遍历
	t.Run("按行遍历", func(t *testing.T) {
		for row := 0; row < table.GetRowCount(); row++ {
			visitedCount := 0
			err := table.ForEachInRow(row, func(col int, cell *document.TableCell, text string) error {
				visitedCount++
				return nil
			})

			if err != nil {
				t.Errorf("第%d行遍历失败: %v", row, err)
			}

			if visitedCount != 3 {
				t.Errorf("第%d行单元格数量不正确: 期望3，实际%d", row, visitedCount)
			}
		}
	})

	// 测试5: 按列遍历
	t.Run("按列遍历", func(t *testing.T) {
		for col := 0; col < table.GetColumnCount(); col++ {
			visitedCount := 0
			err := table.ForEachInColumn(col, func(row int, cell *document.TableCell, text string) error {
				visitedCount++
				return nil
			})

			if err != nil {
				t.Errorf("第%d列遍历失败: %v", col, err)
			}

			if visitedCount != 4 {
				t.Errorf("第%d列单元格数量不正确: 期望4，实际%d", col, visitedCount)
			}
		}
	})

	// 测试6: 范围获取
	t.Run("单元格范围", func(t *testing.T) {
		// 获取数据区域 (1,0) 到 (3,2)
		cells, err := table.GetCellRange(1, 0, 3, 2)
		if err != nil {
			t.Errorf("获取范围失败: %v", err)
		}

		expectedCount := 9 // 3行x3列
		if len(cells) != expectedCount {
			t.Errorf("范围单元格数量不正确: 期望%d，实际%d", expectedCount, len(cells))
		}

		// 验证范围内容
		if cells[0].Row != 1 || cells[0].Col != 0 {
			t.Errorf("范围起始位置不正确: 期望(1,0)，实际(%d,%d)", cells[0].Row, cells[0].Col)
		}

		lastIndex := len(cells) - 1
		if cells[lastIndex].Row != 3 || cells[lastIndex].Col != 2 {
			t.Errorf("范围结束位置不正确: 期望(3,2)，实际(%d,%d)",
				cells[lastIndex].Row, cells[lastIndex].Col)
		}
	})

	// 测试7: 查找功能
	t.Run("单元格查找", func(t *testing.T) {
		// 查找包含"张"的单元格
		cells, err := table.FindCellsByText("张", false)
		if err != nil {
			t.Errorf("查找失败: %v", err)
		}

		if len(cells) != 1 {
			t.Errorf("查找结果数量不正确: 期望1，实际%d", len(cells))
		}

		if len(cells) > 0 && cells[0].Text != "张三" {
			t.Errorf("查找内容不正确: 期望'张三'，实际'%s'", cells[0].Text)
		}

		// 精确查找
		exactCells, err := table.FindCellsByText("25", true)
		if err != nil {
			t.Errorf("精确查找失败: %v", err)
		}

		if len(exactCells) != 1 {
			t.Errorf("精确查找结果数量不正确: 期望1，实际%d", len(exactCells))
		}
	})

	// 测试8: 自定义查找条件
	t.Run("自定义查找", func(t *testing.T) {
		// 查找年龄大于26的行
		ageCells, err := table.FindCells(func(row, col int, cell *document.TableCell, text string) bool {
			// 检查年龄列（第2列）
			if col == 1 && row > 0 {
				// 简单检查是否包含数字且可能大于26
				return text == "30" || text == "28"
			}
			return false
		})

		if err != nil {
			t.Errorf("自定义查找失败: %v", err)
		}

		if len(ageCells) != 2 {
			t.Errorf("自定义查找结果数量不正确: 期望2，实际%d", len(ageCells))
		}
	})

	// 保存测试文档
	outputDir := "../examples/output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Logf("创建输出目录失败: %v", err)
	}

	filename := filepath.Join(outputDir, "cell_iterator_integration_test.docx")
	if err := doc.Save(filename); err != nil {
		t.Errorf("保存测试文档失败: %v", err)
	} else {
		t.Logf("测试文档已保存到: %s", filename)
	}
}

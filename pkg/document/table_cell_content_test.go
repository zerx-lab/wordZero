package document

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// createCellTestImage 创建测试用的PNG图片数据
func createCellTestImage(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充红色背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{255, 100, 100, 255})
		}
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

// TestAddCellParagraph 测试向单元格添加段落
func TestAddCellParagraph(t *testing.T) {
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

	// 测试添加段落
	para, err := table.AddCellParagraph(0, 0, "第一段内容")
	if err != nil {
		t.Errorf("添加段落失败: %v", err)
	}
	if para == nil {
		t.Error("返回的段落不应为空")
	}

	// 添加第二个段落
	para2, err := table.AddCellParagraph(0, 0, "第二段内容")
	if err != nil {
		t.Errorf("添加第二段落失败: %v", err)
	}
	if para2 == nil {
		t.Error("返回的第二段落不应为空")
	}

	// 验证段落数量
	paragraphs, err := table.GetCellParagraphs(0, 0)
	if err != nil {
		t.Errorf("获取段落失败: %v", err)
	}

	// 初始有一个空段落，加上两个新段落
	if len(paragraphs) < 3 {
		t.Errorf("期望至少3个段落，实际%d", len(paragraphs))
	}

	// 测试无效索引
	_, err = table.AddCellParagraph(10, 10, "无效")
	if err == nil {
		t.Error("期望无效索引失败，但成功了")
	}
}

// TestAddCellFormattedParagraph 测试向单元格添加格式化段落
func TestAddCellFormattedParagraph(t *testing.T) {
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

	// 测试添加格式化段落
	format := &TextFormat{
		Bold:       true,
		Italic:     true,
		FontSize:   14,
		FontColor:  "FF0000",
		FontFamily: "Arial",
		Underline:  true,
	}

	para, err := table.AddCellFormattedParagraph(0, 0, "格式化内容", format)
	if err != nil {
		t.Errorf("添加格式化段落失败: %v", err)
	}
	if para == nil {
		t.Error("返回的段落不应为空")
	}

	// 验证格式
	if len(para.Runs) == 0 {
		t.Error("段落应包含至少一个Run")
	}

	run := para.Runs[0]
	if run.Properties == nil {
		t.Error("Run应有属性")
	} else {
		if run.Properties.Bold == nil {
			t.Error("期望粗体属性")
		}
		if run.Properties.Italic == nil {
			t.Error("期望斜体属性")
		}
		if run.Properties.Underline == nil {
			t.Error("期望下划线属性")
		}
	}
}

// TestClearCellParagraphs 测试清空单元格段落
func TestClearCellParagraphs(t *testing.T) {
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

	// 添加多个段落
	table.AddCellParagraph(0, 0, "段落1")
	table.AddCellParagraph(0, 0, "段落2")

	// 清空段落
	err = table.ClearCellParagraphs(0, 0)
	if err != nil {
		t.Errorf("清空段落失败: %v", err)
	}

	// 验证清空后只有一个空段落
	paragraphs, err := table.GetCellParagraphs(0, 0)
	if err != nil {
		t.Errorf("获取段落失败: %v", err)
	}

	if len(paragraphs) != 1 {
		t.Errorf("期望清空后只有1个段落，实际%d", len(paragraphs))
	}

	// 测试无效索引
	err = table.ClearCellParagraphs(10, 10)
	if err == nil {
		t.Error("期望无效索引失败，但成功了")
	}
}

// TestAddNestedTable 测试向单元格添加嵌套表格
func TestAddNestedTable(t *testing.T) {
	doc := New()

	// 创建主表格
	mainConfig := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 8000,
	}

	mainTable, err := doc.CreateTable(mainConfig)
	if err != nil {
		t.Fatalf("创建主表格失败: %v", err)
	}

	// 创建嵌套表格配置
	nestedConfig := &TableConfig{
		Rows:  2,
		Cols:  3,
		Width: 3000,
		Data: [][]string{
			{"嵌套1", "嵌套2", "嵌套3"},
			{"数据1", "数据2", "数据3"},
		},
	}

	// 添加嵌套表格
	nestedTable, err := mainTable.AddNestedTable(0, 0, nestedConfig)
	if err != nil {
		t.Errorf("添加嵌套表格失败: %v", err)
	}
	if nestedTable == nil {
		t.Error("返回的嵌套表格不应为空")
	}

	// 验证嵌套表格结构
	if nestedTable.GetRowCount() != 2 {
		t.Errorf("期望嵌套表格2行，实际%d", nestedTable.GetRowCount())
	}
	if nestedTable.GetColumnCount() != 3 {
		t.Errorf("期望嵌套表格3列，实际%d", nestedTable.GetColumnCount())
	}

	// 验证嵌套表格内容
	cellText, err := nestedTable.GetCellText(0, 0)
	if err != nil {
		t.Errorf("获取嵌套表格单元格内容失败: %v", err)
	}
	if cellText != "嵌套1" {
		t.Errorf("期望嵌套表格内容'嵌套1'，实际'%s'", cellText)
	}

	// 获取嵌套表格列表
	nestedTables, err := mainTable.GetNestedTables(0, 0)
	if err != nil {
		t.Errorf("获取嵌套表格列表失败: %v", err)
	}
	if len(nestedTables) != 1 {
		t.Errorf("期望1个嵌套表格，实际%d", len(nestedTables))
	}
}

// TestAddNestedTableInvalidConfig 测试嵌套表格的无效配置
func TestAddNestedTableInvalidConfig(t *testing.T) {
	doc := New()

	mainConfig := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	mainTable, err := doc.CreateTable(mainConfig)
	if err != nil {
		t.Fatalf("创建主表格失败: %v", err)
	}

	// 测试无效的行列数
	_, err = mainTable.AddNestedTable(0, 0, &TableConfig{Rows: 0, Cols: 2, Width: 2000})
	if err == nil {
		t.Error("期望行数为0时失败，但成功了")
	}

	_, err = mainTable.AddNestedTable(0, 0, &TableConfig{Rows: 2, Cols: 0, Width: 2000})
	if err == nil {
		t.Error("期望列数为0时失败，但成功了")
	}

	// 测试无效的单元格索引
	_, err = mainTable.AddNestedTable(10, 10, &TableConfig{Rows: 2, Cols: 2, Width: 2000})
	if err == nil {
		t.Error("期望无效索引失败，但成功了")
	}
}

// TestAddCellList 测试向单元格添加列表
func TestAddCellList(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  3,
		Cols:  2,
		Width: 6000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试添加无序列表
	bulletListConfig := &CellListConfig{
		Type:         ListTypeBullet,
		BulletSymbol: BulletTypeDot,
		Items:        []string{"项目一", "项目二", "项目三"},
	}

	err = table.AddCellList(0, 0, bulletListConfig)
	if err != nil {
		t.Errorf("添加无序列表失败: %v", err)
	}

	// 验证列表项数量
	paragraphs, err := table.GetCellParagraphs(0, 0)
	if err != nil {
		t.Errorf("获取段落失败: %v", err)
	}

	// 初始有一个空段落，加上3个列表项
	expectedCount := 1 + 3
	if len(paragraphs) != expectedCount {
		t.Errorf("期望%d个段落，实际%d", expectedCount, len(paragraphs))
	}

	// 测试添加有序列表
	numberListConfig := &CellListConfig{
		Type:  ListTypeNumber,
		Items: []string{"第一步", "第二步", "第三步"},
	}

	err = table.AddCellList(1, 0, numberListConfig)
	if err != nil {
		t.Errorf("添加有序列表失败: %v", err)
	}

	// 测试添加小写字母列表
	letterListConfig := &CellListConfig{
		Type:  ListTypeLowerLetter,
		Items: []string{"选项a", "选项b"},
	}

	err = table.AddCellList(2, 0, letterListConfig)
	if err != nil {
		t.Errorf("添加字母列表失败: %v", err)
	}
}

// TestAddCellListInvalidConfig 测试列表的无效配置
func TestAddCellListInvalidConfig(t *testing.T) {
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

	// 测试空配置
	err = table.AddCellList(0, 0, nil)
	if err == nil {
		t.Error("期望空配置失败，但成功了")
	}

	// 测试空列表项
	err = table.AddCellList(0, 0, &CellListConfig{Type: ListTypeBullet, Items: []string{}})
	if err == nil {
		t.Error("期望空列表项失败，但成功了")
	}

	// 测试无效索引
	err = table.AddCellList(10, 10, &CellListConfig{Type: ListTypeBullet, Items: []string{"测试"}})
	if err == nil {
		t.Error("期望无效索引失败，但成功了")
	}
}

// TestAddCellImage 测试向单元格添加图片
func TestAddCellImage(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 6000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 创建测试图片数据
	imageData := createCellTestImage(100, 100)

	// 测试从数据添加图片
	imageInfo, err := doc.AddCellImageFromData(table, 0, 0, imageData, 30)
	if err != nil {
		t.Errorf("添加图片失败: %v", err)
	}
	if imageInfo == nil {
		t.Error("返回的图片信息不应为空")
	}

	// 验证图片ID不为空
	if imageInfo.ID == "" {
		t.Error("图片ID不应为空")
	}

	// 验证关系ID不为空
	if imageInfo.RelationID == "" {
		t.Error("关系ID不应为空")
	}

	// 验证单元格段落包含图片
	paragraphs, err := table.GetCellParagraphs(0, 0)
	if err != nil {
		t.Errorf("获取段落失败: %v", err)
	}

	hasImage := false
	for _, para := range paragraphs {
		for _, run := range para.Runs {
			if run.Drawing != nil {
				hasImage = true
				break
			}
		}
	}

	if !hasImage {
		t.Error("单元格应包含图片")
	}
}

// TestAddCellImageWithConfig 测试使用配置添加图片
func TestAddCellImageWithConfig(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 6000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 创建测试图片数据
	imageData := createCellTestImage(200, 150)

	// 使用完整配置添加图片
	imageConfig := &CellImageConfig{
		Data:            imageData,
		Width:           50,
		Height:          40,
		KeepAspectRatio: false,
		AltText:         "测试图片",
		Title:           "单元格图片",
	}

	imageInfo, err := doc.AddCellImage(table, 0, 0, imageConfig)
	if err != nil {
		t.Errorf("添加图片失败: %v", err)
	}

	// 验证图片配置
	if imageInfo.Config == nil {
		t.Error("图片配置不应为空")
	} else {
		if imageInfo.Config.AltText != "测试图片" {
			t.Errorf("期望替代文字'测试图片'，实际'%s'", imageInfo.Config.AltText)
		}
		if imageInfo.Config.Title != "单元格图片" {
			t.Errorf("期望标题'单元格图片'，实际'%s'", imageInfo.Config.Title)
		}
	}
}

// TestAddCellImageInvalidCases 测试添加图片的无效情况
func TestAddCellImageInvalidCases(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 4000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 测试空表格
	_, err = doc.AddCellImage(nil, 0, 0, &CellImageConfig{Data: createCellTestImage(100, 100)})
	if err == nil {
		t.Error("期望空表格失败，但成功了")
	}

	// 测试无效索引
	_, err = doc.AddCellImage(table, 10, 10, &CellImageConfig{Data: createCellTestImage(100, 100)})
	if err == nil {
		t.Error("期望无效索引失败，但成功了")
	}

	// 测试无数据配置
	_, err = doc.AddCellImage(table, 0, 0, &CellImageConfig{})
	if err == nil {
		t.Error("期望无数据配置失败，但成功了")
	}
}

// TestComplexTableStructure 测试复杂表格结构
func TestComplexTableStructure(t *testing.T) {
	doc := New()

	// 创建主表格
	mainConfig := &TableConfig{
		Rows:  3,
		Cols:  3,
		Width: 9000,
	}

	table, err := doc.AddTable(mainConfig)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 第一个单元格：添加多个段落
	table.AddCellParagraph(0, 0, "第一段")
	table.AddCellFormattedParagraph(0, 0, "格式化段落", &TextFormat{Bold: true})

	// 第二个单元格：添加列表
	listConfig := &CellListConfig{
		Type:         ListTypeBullet,
		BulletSymbol: BulletTypeDot,
		Items:        []string{"列表项1", "列表项2"},
	}
	table.AddCellList(0, 1, listConfig)

	// 第三个单元格：添加嵌套表格
	nestedConfig := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 2500,
		Data: [][]string{
			{"A", "B"},
			{"C", "D"},
		},
	}
	table.AddNestedTable(0, 2, nestedConfig)

	// 第四个单元格：添加图片
	imageData := createCellTestImage(50, 50)
	doc.AddCellImageFromData(table, 1, 0, imageData, 20)

	// 验证复杂结构
	paragraphs00, _ := table.GetCellParagraphs(0, 0)
	if len(paragraphs00) < 3 { // 初始1个 + 添加2个
		t.Errorf("单元格(0,0)应至少有3个段落，实际%d", len(paragraphs00))
	}

	paragraphs01, _ := table.GetCellParagraphs(0, 1)
	if len(paragraphs01) < 3 { // 初始1个 + 列表2项
		t.Errorf("单元格(0,1)应至少有3个段落，实际%d", len(paragraphs01))
	}

	nestedTables, _ := table.GetNestedTables(0, 2)
	if len(nestedTables) != 1 {
		t.Errorf("单元格(0,2)应有1个嵌套表格，实际%d", len(nestedTables))
	}

	paragraphs10, _ := table.GetCellParagraphs(1, 0)
	hasImage := false
	for _, para := range paragraphs10 {
		for _, run := range para.Runs {
			if run.Drawing != nil {
				hasImage = true
				break
			}
		}
	}
	if !hasImage {
		t.Error("单元格(1,0)应包含图片")
	}
}

// TestSaveComplexTable 测试保存复杂表格
func TestSaveComplexTable(t *testing.T) {
	doc := New()

	// 创建表格
	config := &TableConfig{
		Rows:  3,
		Cols:  2,
		Width: 6000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	// 添加复杂内容
	table.AddCellParagraph(0, 0, "复杂表格测试")
	table.AddCellFormattedParagraph(0, 0, "粗体文本", &TextFormat{Bold: true})

	listConfig := &CellListConfig{
		Type:  ListTypeNumber,
		Items: []string{"第一项", "第二项"},
	}
	table.AddCellList(0, 1, listConfig)

	nestedConfig := &TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 2000,
		Data: [][]string{
			{"X", "Y"},
			{"Z", "W"},
		},
	}
	table.AddNestedTable(1, 0, nestedConfig)

	// 添加图片
	imageData := createCellTestImage(80, 60)
	doc.AddCellImageFromData(table, 1, 1, imageData, 25)

	// 保存并验证
	outputDir := "test_output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}

	outputFile := outputDir + "/complex_table_test.docx"
	err = doc.Save(outputFile)
	if err != nil {
		t.Errorf("保存文档失败: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("输出文件不存在")
	}

	// 清理
	defer os.RemoveAll(outputDir)
}

// TestRomanNumerals 测试罗马数字转换
func TestRomanNumerals(t *testing.T) {
	testCases := []struct {
		num      int
		expected string
	}{
		{1, "I"},
		{2, "II"},
		{3, "III"},
		{4, "IV"},
		{5, "V"},
		{9, "IX"},
		{10, "X"},
		{40, "XL"},
		{50, "L"},
		{90, "XC"},
		{100, "C"},
		{400, "CD"},
		{500, "D"},
		{900, "CM"},
		{1000, "M"},
		{1999, "MCMXCIX"},
		{2024, "MMXXIV"},
	}

	for _, tc := range testCases {
		result := toRomanUpper(tc.num)
		if result != tc.expected {
			t.Errorf("toRomanUpper(%d) = %s, 期望 %s", tc.num, result, tc.expected)
		}
	}

	// 测试边界情况
	if toRomanUpper(0) != "0" {
		t.Error("0应返回字符串'0'")
	}

	if toRomanUpper(4000) != "4000" {
		t.Error("4000应返回字符串'4000'")
	}
}

// TestAddCellListAllTypes 测试所有列表类型
func TestAddCellListAllTypes(t *testing.T) {
	doc := New()

	config := &TableConfig{
		Rows:  7,
		Cols:  1,
		Width: 3000,
	}

	table, err := doc.CreateTable(config)
	if err != nil {
		t.Fatalf("创建表格失败: %v", err)
	}

	testTypes := []struct {
		listType ListType
		name     string
	}{
		{ListTypeBullet, "无序列表"},
		{ListTypeNumber, "数字列表"},
		{ListTypeDecimal, "十进制列表"},
		{ListTypeLowerLetter, "小写字母列表"},
		{ListTypeUpperLetter, "大写字母列表"},
		{ListTypeLowerRoman, "小写罗马列表"},
		{ListTypeUpperRoman, "大写罗马列表"},
	}

	for i, tc := range testTypes {
		listConfig := &CellListConfig{
			Type:  tc.listType,
			Items: []string{"项目1", "项目2", "项目3"},
		}

		err := table.AddCellList(i, 0, listConfig)
		if err != nil {
			t.Errorf("添加%s失败: %v", tc.name, err)
		}
	}
}

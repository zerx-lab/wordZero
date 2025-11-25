package style

import (
	"testing"
)

func TestNewQuickStyleAPI(t *testing.T) {
	sm := NewStyleManager()
	api := NewQuickStyleAPI(sm)

	if api == nil {
		t.Fatal("NewQuickStyleAPI 返回了 nil")
	}

	if api.styleManager != sm {
		t.Error("QuickStyleAPI 的 styleManager 设置不正确")
	}
}

func TestGetStyleInfo(t *testing.T) {
	sm := NewStyleManager()
	api := NewQuickStyleAPI(sm)

	// 测试获取存在的样式信息
	info, err := api.GetStyleInfo("Heading1")
	if err != nil {
		t.Fatalf("获取样式信息失败: %v", err)
	}

	if info.ID != "Heading1" {
		t.Errorf("期望样式ID为 'Heading1'，实际为 '%s'", info.ID)
	}

	if info.Name != "heading 1" {
		t.Errorf("期望样式名称为 'heading 1'，实际为 '%s'", info.Name)
	}

	if info.Type != StyleTypeParagraph {
		t.Errorf("期望样式类型为 '%s'，实际为 '%s'", StyleTypeParagraph, info.Type)
	}

	if !info.IsBuiltIn {
		t.Error("Heading1 应该是内置样式")
	}

	// 测试获取不存在的样式信息
	_, err = api.GetStyleInfo("NonExistentStyle")
	if err == nil {
		t.Error("期望获取不存在样式时返回错误")
	}
}

func TestGetAllStylesInfo(t *testing.T) {
	sm := NewStyleManager()
	api := NewQuickStyleAPI(sm)

	allStyles := api.GetAllStylesInfo()

	if len(allStyles) == 0 {
		t.Error("期望返回样式信息列表不为空")
	}

	// 检查是否包含预期的样式
	styleFound := false
	for _, info := range allStyles {
		if info.ID == "Normal" {
			styleFound = true
			break
		}
	}

	if !styleFound {
		t.Error("期望在样式列表中找到 'Normal' 样式")
	}
}

func TestGetHeadingStylesInfo(t *testing.T) {
	sm := NewStyleManager()
	api := NewQuickStyleAPI(sm)

	headingStyles := api.GetHeadingStylesInfo()

	expectedCount := 9 // Heading1 到 Heading9
	if len(headingStyles) != expectedCount {
		t.Errorf("期望标题样式数量为 %d，实际为 %d", expectedCount, len(headingStyles))
	}

	// 检查标题样式的顺序和ID
	for i, info := range headingStyles {
		expectedID := "Heading" + string(rune('1'+i))
		if info.ID != expectedID {
			t.Errorf("期望第 %d 个标题样式ID为 '%s'，实际为 '%s'", i+1, expectedID, info.ID)
		}

		if info.Type != StyleTypeParagraph {
			t.Errorf("标题样式 '%s' 应该是段落类型", info.ID)
		}
	}
}

func TestGetParagraphStylesInfo(t *testing.T) {
	sm := NewStyleManager()
	api := NewQuickStyleAPI(sm)

	paragraphStyles := api.GetParagraphStylesInfo()

	if len(paragraphStyles) == 0 {
		t.Error("期望段落样式列表不为空")
	}

	// 检查所有返回的样式都是段落类型
	for _, info := range paragraphStyles {
		if info.Type != StyleTypeParagraph {
			t.Errorf("样式 '%s' 应该是段落类型，实际为 '%s'", info.ID, info.Type)
		}
	}
}

func TestGetCharacterStylesInfo(t *testing.T) {
	sm := NewStyleManager()
	api := NewQuickStyleAPI(sm)

	characterStyles := api.GetCharacterStylesInfo()

	if len(characterStyles) == 0 {
		t.Error("期望字符样式列表不为空")
	}

	// 检查所有返回的样式都是字符类型
	for _, info := range characterStyles {
		if info.Type != StyleTypeCharacter {
			t.Errorf("样式 '%s' 应该是字符类型，实际为 '%s'", info.ID, info.Type)
		}
	}
}

func TestCreateQuickStyle(t *testing.T) {
	sm := NewStyleManager()
	api := NewQuickStyleAPI(sm)

	// 测试创建自定义段落样式
	config := QuickStyleConfig{
		ID:      "TestCustomStyle",
		Name:    "测试自定义样式",
		Type:    StyleTypeParagraph,
		BasedOn: "Normal",
		ParagraphConfig: &QuickParagraphConfig{
			Alignment:   "center",
			LineSpacing: 1.5,
			SpaceBefore: 12,
			SpaceAfter:  6,
		},
		RunConfig: &QuickRunConfig{
			FontName:  "微软雅黑",
			FontSize:  14,
			FontColor: "FF0000",
			Bold:      true,
		},
	}

	style, err := api.CreateQuickStyle(config)
	if err != nil {
		t.Fatalf("创建自定义样式失败: %v", err)
	}

	if style.StyleID != "TestCustomStyle" {
		t.Errorf("期望样式ID为 'TestCustomStyle'，实际为 '%s'", style.StyleID)
	}

	if !style.CustomStyle {
		t.Error("创建的样式应该标记为自定义样式")
	}

	// 验证段落属性
	if style.ParagraphPr == nil {
		t.Error("自定义样式应该包含段落属性")
	} else {
		if style.ParagraphPr.Justification == nil || style.ParagraphPr.Justification.Val != "center" {
			t.Error("段落对齐方式设置不正确")
		}
	}

	// 验证字符属性
	if style.RunPr == nil {
		t.Error("自定义样式应该包含字符属性")
	} else {
		if style.RunPr.Bold == nil {
			t.Error("粗体属性设置不正确")
		}
		if style.RunPr.FontSize == nil || style.RunPr.FontSize.Val != "28" {
			t.Error("字体大小设置不正确")
		}
	}

	// 测试创建重复ID的样式
	_, err = api.CreateQuickStyle(config)
	if err == nil {
		t.Error("期望创建重复ID样式时返回错误")
	}
}

func TestCreateParagraphProperties(t *testing.T) {
	config := &QuickParagraphConfig{
		Alignment:       "center",
		LineSpacing:     1.5,
		SpaceBefore:     12,
		SpaceAfter:      6,
		FirstLineIndent: 24,
		LeftIndent:      36,
		RightIndent:     36,
	}

	props := createParagraphProperties(config)

	if props == nil {
		t.Fatal("createParagraphProperties 返回了 nil")
	}

	// 检查对齐方式
	if props.Justification == nil || props.Justification.Val != "center" {
		t.Error("对齐方式设置不正确")
	}

	// 检查间距
	if props.Spacing == nil {
		t.Error("间距属性未设置")
	} else {
		if props.Spacing.Before != "240" { // 12 * 20
			t.Errorf("段前间距设置不正确，期望 '240'，实际 '%s'", props.Spacing.Before)
		}
		if props.Spacing.After != "120" { // 6 * 20
			t.Errorf("段后间距设置不正确，期望 '120'，实际 '%s'", props.Spacing.After)
		}
	}

	// 检查缩进
	if props.Indentation == nil {
		t.Error("缩进属性未设置")
	} else {
		if props.Indentation.FirstLine != "480" { // 24 * 20
			t.Errorf("首行缩进设置不正确，期望 '480'，实际 '%s'", props.Indentation.FirstLine)
		}
	}
}

func TestCreateRunProperties(t *testing.T) {
	config := &QuickRunConfig{
		FontName:  "微软雅黑",
		FontSize:  14,
		FontColor: "FF0000",
		Bold:      true,
		Italic:    true,
		Underline: true,
		Strike:    true,
		Highlight: "yellow",
	}

	props := createRunProperties(config)

	if props == nil {
		t.Fatal("createRunProperties 返回了 nil")
	}

	// 检查字体设置
	if props.FontFamily == nil {
		t.Error("字体系列未设置")
	} else {
		if props.FontFamily.ASCII != "微软雅黑" {
			t.Errorf("ASCII字体设置不正确，期望 '微软雅黑'，实际 '%s'", props.FontFamily.ASCII)
		}
	}

	if props.FontSize == nil || props.FontSize.Val != "28" { // 14 * 2
		t.Error("字体大小设置不正确")
	}

	if props.Color == nil || props.Color.Val != "FF0000" {
		t.Error("字体颜色设置不正确")
	}

	// 检查格式设置
	if props.Bold == nil {
		t.Error("粗体设置不正确")
	}

	if props.Italic == nil {
		t.Error("斜体设置不正确")
	}

	if props.Underline == nil || props.Underline.Val != "single" {
		t.Error("下划线设置不正确")
	}

	if props.Strike == nil {
		t.Error("删除线设置不正确")
	}

	if props.Highlight == nil || props.Highlight.Val != "yellow" {
		t.Error("高亮设置不正确")
	}
}

func TestCreateParagraphPropertiesWithSnapToGrid(t *testing.T) {
	// 测试 SnapToGrid = false 时禁用网格对齐
	snapToGridFalse := false
	config := &QuickParagraphConfig{
		Alignment:   "left",
		LineSpacing: 1.5,
		SnapToGrid:  &snapToGridFalse,
	}

	props := createParagraphProperties(config)

	if props == nil {
		t.Fatal("createParagraphProperties 返回了 nil")
	}

	// 检查 SnapToGrid 设置
	if props.SnapToGrid == nil {
		t.Error("SnapToGrid 应该被设置")
	} else {
		if props.SnapToGrid.Val != "0" {
			t.Errorf("SnapToGrid.Val 设置不正确，期望 '0'，实际 '%s'", props.SnapToGrid.Val)
		}
	}

	// 检查行间距
	if props.Spacing == nil {
		t.Error("间距属性未设置")
	} else {
		if props.Spacing.Line != "360" { // 1.5 * 240
			t.Errorf("行间距设置不正确，期望 '360'，实际 '%s'", props.Spacing.Line)
		}
		if props.Spacing.LineRule != "auto" {
			t.Errorf("LineRule 设置不正确，期望 'auto'，实际 '%s'", props.Spacing.LineRule)
		}
	}

	// 测试 SnapToGrid = true 时不设置（保持默认）
	snapToGridTrue := true
	configWithGridEnabled := &QuickParagraphConfig{
		Alignment:   "left",
		LineSpacing: 1.5,
		SnapToGrid:  &snapToGridTrue,
	}

	propsWithGrid := createParagraphProperties(configWithGridEnabled)

	if propsWithGrid.SnapToGrid != nil {
		t.Error("当 SnapToGrid = true 时，不应该设置 SnapToGrid 属性（保持默认行为）")
	}

	// 测试 SnapToGrid = nil 时不设置
	configWithoutGrid := &QuickParagraphConfig{
		Alignment:   "left",
		LineSpacing: 1.5,
		SnapToGrid:  nil,
	}

	propsWithoutGrid := createParagraphProperties(configWithoutGrid)

	if propsWithoutGrid.SnapToGrid != nil {
		t.Error("当 SnapToGrid = nil 时，不应该设置 SnapToGrid 属性")
	}
}

// Package style 样式应用API
package style

import "fmt"

// StyleApplicator 样式应用器接口
type StyleApplicator interface {
	ApplyStyle(styleID string) error
	ApplyHeadingStyle(level int) error
	ApplyTitleStyle() error
	ApplySubtitleStyle() error
	ApplyQuoteStyle() error
	ApplyCodeBlockStyle() error
	ApplyListParagraphStyle() error
	ApplyNormalStyle() error
}

// QuickStyleAPI 快速样式应用API
type QuickStyleAPI struct {
	styleManager *StyleManager
}

// NewQuickStyleAPI 创建快速样式API
func NewQuickStyleAPI(styleManager *StyleManager) *QuickStyleAPI {
	return &QuickStyleAPI{
		styleManager: styleManager,
	}
}

// GetStyleInfo 获取样式信息（用于UI显示）
func (api *QuickStyleAPI) GetStyleInfo(styleID string) (*StyleInfo, error) {
	style := api.styleManager.GetStyle(styleID)
	if style == nil {
		return nil, fmt.Errorf("样式 %s 不存在", styleID)
	}

	return &StyleInfo{
		ID:          style.StyleID,
		Name:        getStyleDisplayName(style),
		Type:        StyleType(style.Type),
		Description: getStyleDescription(styleID),
		IsBuiltIn:   !style.CustomStyle,
		BasedOn:     getBasedOnStyleID(style),
	}, nil
}

// StyleInfo 样式信息结构
type StyleInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        StyleType `json:"type"`
	Description string    `json:"description"`
	IsBuiltIn   bool      `json:"isBuiltIn"`
	BasedOn     string    `json:"basedOn,omitempty"`
}

// GetAllStylesInfo 获取所有样式信息
func (api *QuickStyleAPI) GetAllStylesInfo() []*StyleInfo {
	var stylesInfo []*StyleInfo
	for _, style := range api.styleManager.GetAllStyles() {
		info := &StyleInfo{
			ID:          style.StyleID,
			Name:        getStyleDisplayName(style),
			Type:        StyleType(style.Type),
			Description: getStyleDescription(style.StyleID),
			IsBuiltIn:   !style.CustomStyle,
			BasedOn:     getBasedOnStyleID(style),
		}
		stylesInfo = append(stylesInfo, info)
	}
	return stylesInfo
}

// GetHeadingStylesInfo 获取所有标题样式信息
func (api *QuickStyleAPI) GetHeadingStylesInfo() []*StyleInfo {
	var headingStylesInfo []*StyleInfo
	for i := 1; i <= 9; i++ {
		styleID := fmt.Sprintf("Heading%d", i)
		if info, err := api.GetStyleInfo(styleID); err == nil {
			headingStylesInfo = append(headingStylesInfo, info)
		}
	}
	return headingStylesInfo
}

// GetParagraphStylesInfo 获取段落样式信息
func (api *QuickStyleAPI) GetParagraphStylesInfo() []*StyleInfo {
	var paragraphStylesInfo []*StyleInfo
	for _, style := range api.styleManager.GetStylesByType(StyleTypeParagraph) {
		info := &StyleInfo{
			ID:          style.StyleID,
			Name:        getStyleDisplayName(style),
			Type:        StyleType(style.Type),
			Description: getStyleDescription(style.StyleID),
			IsBuiltIn:   !style.CustomStyle,
			BasedOn:     getBasedOnStyleID(style),
		}
		paragraphStylesInfo = append(paragraphStylesInfo, info)
	}
	return paragraphStylesInfo
}

// GetCharacterStylesInfo 获取字符样式信息
func (api *QuickStyleAPI) GetCharacterStylesInfo() []*StyleInfo {
	var characterStylesInfo []*StyleInfo
	for _, style := range api.styleManager.GetStylesByType(StyleTypeCharacter) {
		info := &StyleInfo{
			ID:          style.StyleID,
			Name:        getStyleDisplayName(style),
			Type:        StyleType(style.Type),
			Description: getStyleDescription(style.StyleID),
			IsBuiltIn:   !style.CustomStyle,
			BasedOn:     getBasedOnStyleID(style),
		}
		characterStylesInfo = append(characterStylesInfo, info)
	}
	return characterStylesInfo
}

// CreateQuickStyle 快速创建自定义样式
func (api *QuickStyleAPI) CreateQuickStyle(config QuickStyleConfig) (*Style, error) {
	// 验证样式ID是否已存在
	if api.styleManager.StyleExists(config.ID) {
		return nil, fmt.Errorf("样式ID %s 已存在", config.ID)
	}

	// 创建基础样式
	style := api.styleManager.CreateCustomStyle(
		config.ID,
		config.Name,
		config.Type,
		config.BasedOn,
	)

	// 应用段落属性
	if config.ParagraphConfig != nil {
		style.ParagraphPr = createParagraphProperties(config.ParagraphConfig)
	}

	// 应用字符属性
	if config.RunConfig != nil {
		style.RunPr = createRunProperties(config.RunConfig)
	}

	return style, nil
}

// QuickStyleConfig 快速样式配置
type QuickStyleConfig struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Type            StyleType             `json:"type"`
	BasedOn         string                `json:"basedOn,omitempty"`
	ParagraphConfig *QuickParagraphConfig `json:"paragraphConfig,omitempty"`
	RunConfig       *QuickRunConfig       `json:"runConfig,omitempty"`
}

// QuickParagraphConfig 快速段落配置
type QuickParagraphConfig struct {
	Alignment       string  `json:"alignment,omitempty"`       // left, center, right, justify
	LineSpacing     float64 `json:"lineSpacing,omitempty"`     // 行间距倍数：1.0=单倍行距，1.5=1.5倍行距，2.0=双倍行距（内部转换为OOXML单位：值×240）
	SpaceBefore     int     `json:"spaceBefore,omitempty"`     // 段前间距（磅）
	SpaceAfter      int     `json:"spaceAfter,omitempty"`      // 段后间距（磅）
	FirstLineIndent int     `json:"firstLineIndent,omitempty"` // 首行缩进（磅）
	LeftIndent      int     `json:"leftIndent,omitempty"`      // 左缩进（磅）
	RightIndent     int     `json:"rightIndent,omitempty"`     // 右缩进（磅）
	SnapToGrid      *bool   `json:"snapToGrid,omitempty"`      // 是否对齐网格（设置为false可禁用网格对齐，使行间距精确生效）
}

// QuickRunConfig 快速字符配置
type QuickRunConfig struct {
	FontName  string `json:"fontName,omitempty"`  // 字体名称
	FontSize  int    `json:"fontSize,omitempty"`  // 字体大小（磅）
	FontColor string `json:"fontColor,omitempty"` // 字体颜色（十六进制）
	Bold      bool   `json:"bold,omitempty"`      // 粗体
	Italic    bool   `json:"italic,omitempty"`    // 斜体
	Underline bool   `json:"underline,omitempty"` // 下划线
	Strike    bool   `json:"strike,omitempty"`    // 删除线
	Highlight string `json:"highlight,omitempty"` // 高亮颜色
}

// getStyleDisplayName 获取样式显示名称
func getStyleDisplayName(style *Style) string {
	if style.Name != nil {
		return style.Name.Val
	}
	return style.StyleID
}

// getStyleDescription 获取样式描述
func getStyleDescription(styleID string) string {
	configs := GetPredefinedStyleConfigs()
	for _, config := range configs {
		if config.StyleID == styleID {
			return config.Description
		}
	}
	return ""
}

// getBasedOnStyleID 获取基础样式ID
func getBasedOnStyleID(style *Style) string {
	if style.BasedOn != nil {
		return style.BasedOn.Val
	}
	return ""
}

// createParagraphProperties 创建段落属性
func createParagraphProperties(config *QuickParagraphConfig) *ParagraphProperties {
	props := &ParagraphProperties{}

	// 对齐方式
	if config.Alignment != "" {
		props.Justification = &Justification{Val: config.Alignment}
	}

	// 网格对齐设置
	// 当设置为false时，禁用网格对齐，使自定义行间距能够精确生效
	if config.SnapToGrid != nil && !*config.SnapToGrid {
		props.SnapToGrid = &SnapToGrid{Val: "0"}
	}

	// 间距设置
	if config.LineSpacing > 0 || config.SpaceBefore > 0 || config.SpaceAfter > 0 {
		spacing := &Spacing{}
		if config.SpaceBefore > 0 {
			spacing.Before = fmt.Sprintf("%d", config.SpaceBefore*20) // 转换为twips
		}
		if config.SpaceAfter > 0 {
			spacing.After = fmt.Sprintf("%d", config.SpaceAfter*20) // 转换为twips
		}
		if config.LineSpacing > 0 {
			spacing.Line = fmt.Sprintf("%.0f", config.LineSpacing*240) // 转换为行间距单位
			spacing.LineRule = "auto"
		}
		props.Spacing = spacing
	}

	// 缩进设置
	if config.FirstLineIndent > 0 || config.LeftIndent > 0 || config.RightIndent > 0 {
		indentation := &Indentation{}
		if config.FirstLineIndent > 0 {
			indentation.FirstLine = fmt.Sprintf("%d", config.FirstLineIndent*20) // 转换为twips
		}
		if config.LeftIndent > 0 {
			indentation.Left = fmt.Sprintf("%d", config.LeftIndent*20) // 转换为twips
		}
		if config.RightIndent > 0 {
			indentation.Right = fmt.Sprintf("%d", config.RightIndent*20) // 转换为twips
		}
		props.Indentation = indentation
	}

	return props
}

// createRunProperties 创建字符属性
func createRunProperties(config *QuickRunConfig) *RunProperties {
	props := &RunProperties{}

	// 字体设置
	if config.FontName != "" {
		props.FontFamily = &FontFamily{
			ASCII:    config.FontName,
			EastAsia: config.FontName,
			HAnsi:    config.FontName,
			CS:       config.FontName,
		}
	}

	if config.FontSize > 0 {
		props.FontSize = &FontSize{Val: fmt.Sprintf("%d", config.FontSize*2)} // Word使用半磅单位
	}

	if config.FontColor != "" {
		props.Color = &Color{Val: config.FontColor}
	}

	// 格式设置
	if config.Bold {
		props.Bold = &Bold{}
	}

	if config.Italic {
		props.Italic = &Italic{}
	}

	if config.Underline {
		props.Underline = &Underline{Val: "single"}
	}

	if config.Strike {
		props.Strike = &Strike{}
	}

	if config.Highlight != "" {
		props.Highlight = &Highlight{Val: config.Highlight}
	}

	return props
}

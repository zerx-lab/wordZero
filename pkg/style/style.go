// Package style 提供Word文档样式管理功能
package style

import (
	"encoding/xml"
	"fmt"
)

// StyleType 样式类型
type StyleType string

const (
	// StyleTypeParagraph 段落样式
	StyleTypeParagraph StyleType = "paragraph"
	// StyleTypeCharacter 字符样式
	StyleTypeCharacter StyleType = "character"
	// StyleTypeTable 表格样式
	StyleTypeTable StyleType = "table"
	// StyleTypeNumbering 编号样式
	StyleTypeNumbering StyleType = "numbering"
)

// Style 样式定义
type Style struct {
	XMLName     xml.Name             `xml:"w:style"`
	Type        string               `xml:"w:type,attr"`
	StyleID     string               `xml:"w:styleId,attr"`
	Name        *StyleName           `xml:"w:name,omitempty"`
	BasedOn     *BasedOn             `xml:"w:basedOn,omitempty"`
	Next        *Next                `xml:"w:next,omitempty"`
	Default     bool                 `xml:"w:default,attr,omitempty"`
	CustomStyle bool                 `xml:"w:customStyle,attr,omitempty"`
	ParagraphPr *ParagraphProperties `xml:"w:pPr,omitempty"`
	RunPr       *RunProperties       `xml:"w:rPr,omitempty"`
	TablePr     *TableProperties     `xml:"w:tblPr,omitempty"`
	TableRowPr  *TableRowProperties  `xml:"w:trPr,omitempty"`
	TableCellPr *TableCellProperties `xml:"w:tcPr,omitempty"`
}

// StyleName 样式名称
type StyleName struct {
	XMLName xml.Name `xml:"w:name"`
	Val     string   `xml:"w:val,attr"`
}

// BasedOn 基于样式
type BasedOn struct {
	XMLName xml.Name `xml:"w:basedOn"`
	Val     string   `xml:"w:val,attr"`
}

// Next 下一个样式
type Next struct {
	XMLName xml.Name `xml:"w:next"`
	Val     string   `xml:"w:val,attr"`
}

// ParagraphProperties 段落样式属性
// 注意：字段顺序必须符合OpenXML标准
type ParagraphProperties struct {
	XMLName         xml.Name         `xml:"w:pPr"`
	KeepNext        *KeepNext        `xml:"w:keepNext,omitempty"`
	KeepLines       *KeepLines       `xml:"w:keepLines,omitempty"`
	PageBreak       *PageBreak       `xml:"w:pageBreakBefore,omitempty"`
	ParagraphBorder *ParagraphBorder `xml:"w:pBdr,omitempty"`
	Shading         *Shading         `xml:"w:shd,omitempty"`
	SnapToGrid      *SnapToGrid      `xml:"w:snapToGrid,omitempty"`
	Spacing         *Spacing         `xml:"w:spacing,omitempty"`
	Indentation     *Indentation     `xml:"w:ind,omitempty"`
	Justification   *Justification   `xml:"w:jc,omitempty"`
	OutlineLevel    *OutlineLevel    `xml:"w:outlineLvl,omitempty"`
}

// ParagraphBorder 段落边框
type ParagraphBorder struct {
	XMLName xml.Name             `xml:"w:pBdr"`
	Top     *ParagraphBorderLine `xml:"w:top,omitempty"`
	Left    *ParagraphBorderLine `xml:"w:left,omitempty"`
	Bottom  *ParagraphBorderLine `xml:"w:bottom,omitempty"`
	Right   *ParagraphBorderLine `xml:"w:right,omitempty"`
}

// ParagraphBorderLine 段落边框线
type ParagraphBorderLine struct {
	XMLName xml.Name `xml:""`
	Val     string   `xml:"w:val,attr"`
	Color   string   `xml:"w:color,attr"`
	Sz      string   `xml:"w:sz,attr"`
	Space   string   `xml:"w:space,attr"`
}

// Shading 阴影/填充色
type Shading struct {
	XMLName xml.Name `xml:"w:shd"`
	Fill    string   `xml:"w:fill,attr"`
	Val     string   `xml:"w:val,attr,omitempty"`
}

// RunProperties 字符样式属性
// 注意：字段顺序必须符合OpenXML标准，w:rFonts必须在w:color之前
type RunProperties struct {
	XMLName    xml.Name    `xml:"w:rPr"`
	FontFamily *FontFamily `xml:"w:rFonts,omitempty"`
	Bold       *Bold       `xml:"w:b,omitempty"`
	Italic     *Italic     `xml:"w:i,omitempty"`
	Underline  *Underline  `xml:"w:u,omitempty"`
	Strike     *Strike     `xml:"w:strike,omitempty"`
	Color      *Color      `xml:"w:color,omitempty"`
	FontSize   *FontSize   `xml:"w:sz,omitempty"`
	Highlight  *Highlight  `xml:"w:highlight,omitempty"`
}

// TableProperties 表格样式属性
type TableProperties struct {
	XMLName    xml.Name       `xml:"w:tblPr"`
	TblInd     *TblIndent     `xml:"w:tblInd,omitempty"`     // 表格缩进
	TblBorders *TblBorders    `xml:"w:tblBorders,omitempty"` // 表格边框
	TblCellMar *TblCellMargin `xml:"w:tblCellMar,omitempty"` // 表格单元格边距
}

// TblIndent 表格缩进
type TblIndent struct {
	XMLName xml.Name `xml:"w:tblInd"`
	W       string   `xml:"w:w,attr"`
	Type    string   `xml:"w:type,attr"`
}

// TblBorders 表格边框
type TblBorders struct {
	XMLName xml.Name   `xml:"w:tblBorders"`
	Top     *TblBorder `xml:"w:top,omitempty"`
	Left    *TblBorder `xml:"w:left,omitempty"`
	Bottom  *TblBorder `xml:"w:bottom,omitempty"`
	Right   *TblBorder `xml:"w:right,omitempty"`
	InsideH *TblBorder `xml:"w:insideH,omitempty"`
	InsideV *TblBorder `xml:"w:insideV,omitempty"`
}

// TblBorder 表格边框定义
type TblBorder struct {
	Val   string `xml:"w:val,attr"`
	Sz    string `xml:"w:sz,attr"`
	Space string `xml:"w:space,attr"`
	Color string `xml:"w:color,attr"`
}

// TblCellMargin 表格单元格边距
type TblCellMargin struct {
	XMLName xml.Name      `xml:"w:tblCellMar"`
	Top     *TblCellSpace `xml:"w:top,omitempty"`
	Left    *TblCellSpace `xml:"w:left,omitempty"`
	Bottom  *TblCellSpace `xml:"w:bottom,omitempty"`
	Right   *TblCellSpace `xml:"w:right,omitempty"`
}

// TblCellSpace 表格单元格空间
type TblCellSpace struct {
	W    string `xml:"w:w,attr"`
	Type string `xml:"w:type,attr"`
}

// TableRowProperties 表格行样式属性
type TableRowProperties struct {
	XMLName xml.Name `xml:"w:trPr"`
	// 表格行样式属性将在后续实现
}

// TableCellProperties 表格单元格样式属性
type TableCellProperties struct {
	XMLName xml.Name `xml:"w:tcPr"`
	// 表格单元格样式属性将在后续实现
}

// 基础样式元素定义
type Spacing struct {
	XMLName  xml.Name `xml:"w:spacing"`
	Before   string   `xml:"w:before,attr,omitempty"`
	After    string   `xml:"w:after,attr,omitempty"`
	Line     string   `xml:"w:line,attr,omitempty"`
	LineRule string   `xml:"w:lineRule,attr,omitempty"`
}

type Justification struct {
	XMLName xml.Name `xml:"w:jc"`
	Val     string   `xml:"w:val,attr"`
}

type Indentation struct {
	XMLName   xml.Name `xml:"w:ind"`
	FirstLine string   `xml:"w:firstLine,attr,omitempty"`
	Left      string   `xml:"w:left,attr,omitempty"`
	Right     string   `xml:"w:right,attr,omitempty"`
}

type KeepNext struct {
	XMLName xml.Name `xml:"w:keepNext"`
}

type KeepLines struct {
	XMLName xml.Name `xml:"w:keepLines"`
}

type PageBreak struct {
	XMLName xml.Name `xml:"w:pageBreakBefore"`
}

type OutlineLevel struct {
	XMLName xml.Name `xml:"w:outlineLvl"`
	Val     string   `xml:"w:val,attr"`
}

// SnapToGrid 网格对齐设置
// 设置为 "0" 时禁用网格对齐，"1" 时启用网格对齐，允许自定义行间距生效（符合 OOXML 规范，仅支持 "0" 或 "1"）
// 注意：此类型在 document 包中有相同定义，这是有意为之，因为两个包可独立使用
type SnapToGrid struct {
	XMLName xml.Name `xml:"w:snapToGrid"`
	Val     string   `xml:"w:val,attr,omitempty"`
}

type Bold struct {
	XMLName xml.Name `xml:"w:b"`
}

type Italic struct {
	XMLName xml.Name `xml:"w:i"`
}

type Underline struct {
	XMLName xml.Name `xml:"w:u"`
	Val     string   `xml:"w:val,attr,omitempty"`
}

type Strike struct {
	XMLName xml.Name `xml:"w:strike"`
}

type FontSize struct {
	XMLName xml.Name `xml:"w:sz"`
	Val     string   `xml:"w:val,attr"`
}

type Color struct {
	XMLName xml.Name `xml:"w:color"`
	Val     string   `xml:"w:val,attr"`
}

type FontFamily struct {
	XMLName  xml.Name `xml:"w:rFonts"`
	ASCII    string   `xml:"w:ascii,attr,omitempty"`
	EastAsia string   `xml:"w:eastAsia,attr,omitempty"`
	HAnsi    string   `xml:"w:hAnsi,attr,omitempty"`
	CS       string   `xml:"w:cs,attr,omitempty"`
}

type Highlight struct {
	XMLName xml.Name `xml:"w:highlight"`
	Val     string   `xml:"w:val,attr"`
}

// Styles 样式集合
type Styles struct {
	XMLName xml.Name `xml:"w:styles"`
	Xmlns   string   `xml:"xmlns:w,attr"`
	Styles  []Style  `xml:"w:style"`
}

// StyleManager 样式管理器
type StyleManager struct {
	styles map[string]*Style
}

// NewStyleManager 创建新的样式管理器
func NewStyleManager() *StyleManager {
	sm := &StyleManager{
		styles: make(map[string]*Style),
	}
	sm.initializePredefinedStyles()
	return sm
}

// GetStyle 获取指定ID的样式
func (sm *StyleManager) GetStyle(styleID string) *Style {
	return sm.styles[styleID]
}

// AddStyle 添加样式
func (sm *StyleManager) AddStyle(style *Style) {
	sm.styles[style.StyleID] = style
}

// GetAllStyles 获取所有样式
func (sm *StyleManager) GetAllStyles() []*Style {
	styles := make([]*Style, 0, len(sm.styles))
	for _, style := range sm.styles {
		styles = append(styles, style)
	}
	return styles
}

// initializePredefinedStyles 初始化预定义样式
func (sm *StyleManager) initializePredefinedStyles() {
	// 普通文本样式
	sm.addNormalStyle()

	// 标题样式
	sm.addHeadingStyles()

	// 其他预定义样式
	sm.addSpecialStyles()
}

// addNormalStyle 添加Normal样式
func (sm *StyleManager) addNormalStyle() {
	normalStyle := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Normal",
		Default: true,
		Name: &StyleName{
			Val: "Normal",
		},
		RunPr: &RunProperties{
			FontSize: &FontSize{
				Val: "21", // 五号字体（10.5磅，Word中以半磅为单位）
			},
			FontFamily: &FontFamily{
				ASCII:    "Calibri",
				EastAsia: "宋体",
				HAnsi:    "Calibri",
				CS:       "Times New Roman",
			},
		},
	}
	sm.AddStyle(normalStyle)
}

// addHeadingStyles 添加标题样式
func (sm *StyleManager) addHeadingStyles() {
	// 标题1
	heading1 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading1",
		Name: &StyleName{
			Val: "heading 1",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "240", // 12磅段前间距
				After:  "0",   // 0磅段段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "0",
			},
		},
		RunPr: &RunProperties{
			Bold: &Bold{},
			FontSize: &FontSize{
				Val: "32", // 16磅
			},
			Color: &Color{
				Val: "2F5496", // 深蓝色
			},
		},
	}
	sm.AddStyle(heading1)

	// 标题2
	heading2 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading2",
		Name: &StyleName{
			Val: "heading 2",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "0",   // 0磅段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "1",
			},
		},
		RunPr: &RunProperties{
			Bold: &Bold{},
			FontSize: &FontSize{
				Val: "26", // 13磅
			},
			Color: &Color{
				Val: "2F5496", // 深蓝色
			},
		},
	}
	sm.AddStyle(heading2)

	// 标题3
	heading3 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading3",
		Name: &StyleName{
			Val: "heading 3",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "0",   // 0磅段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "2",
			},
		},
		RunPr: &RunProperties{
			Bold: &Bold{},
			FontSize: &FontSize{
				Val: "24", // 12磅
			},
			Color: &Color{
				Val: "1F3763", // 深蓝色
			},
		},
	}
	sm.AddStyle(heading3)

	// 标题4
	heading4 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading4",
		Name: &StyleName{
			Val: "heading 4",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "0",   // 0磅段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "3",
			},
		},
		RunPr: &RunProperties{
			Bold:   &Bold{},
			Italic: &Italic{},
			FontSize: &FontSize{
				Val: "22", // 11磅
			},
			Color: &Color{
				Val: "2F5496", // 深蓝色
			},
		},
	}
	sm.AddStyle(heading4)

	// 标题5
	heading5 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading5",
		Name: &StyleName{
			Val: "heading 5",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "0",   // 0磅段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "4",
			},
		},
		RunPr: &RunProperties{
			FontSize: &FontSize{
				Val: "22", // 11磅
			},
			Color: &Color{
				Val: "2F5496", // 深蓝色
			},
		},
	}
	sm.AddStyle(heading5)

	// 标题6
	heading6 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading6",
		Name: &StyleName{
			Val: "heading 6",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "0",   // 0磅段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "5",
			},
		},
		RunPr: &RunProperties{
			Italic: &Italic{},
			FontSize: &FontSize{
				Val: "22", // 11磅
			},
			Color: &Color{
				Val: "1F3763", // 深蓝色
			},
		},
	}
	sm.AddStyle(heading6)

	// 标题7
	heading7 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading7",
		Name: &StyleName{
			Val: "heading 7",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "0",   // 0磅段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "6",
			},
		},
		RunPr: &RunProperties{
			FontSize: &FontSize{
				Val: "20", // 10磅
			},
			Color: &Color{
				Val: "1F3763", // 深蓝色
			},
		},
	}
	sm.AddStyle(heading7)

	// 标题8
	heading8 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading8",
		Name: &StyleName{
			Val: "heading 8",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "0",   // 0磅段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "7",
			},
		},
		RunPr: &RunProperties{
			Italic: &Italic{},
			FontSize: &FontSize{
				Val: "20", // 10磅
			},
			Color: &Color{
				Val: "272727", // 深灰色
			},
		},
	}
	sm.AddStyle(heading8)

	// 标题9
	heading9 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Heading9",
		Name: &StyleName{
			Val: "heading 9",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			KeepNext:  &KeepNext{},
			KeepLines: &KeepLines{},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "0",   // 0磅段后间距
			},
			OutlineLevel: &OutlineLevel{
				Val: "8",
			},
		},
		RunPr: &RunProperties{
			FontSize: &FontSize{
				Val: "18", // 9磅
			},
			Color: &Color{
				Val: "272727", // 深灰色
			},
		},
	}
	sm.AddStyle(heading9)
}

// addSpecialStyles 添加其他特殊样式
func (sm *StyleManager) addSpecialStyles() {
	// 文档标题样式
	title := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Title",
		Name: &StyleName{
			Val: "标题",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Justification: &Justification{
				Val: "center", // 居中对齐
			},
			Spacing: &Spacing{
				Before: "240", // 12磅段前间距
				After:  "60",  // 3磅段后间距
			},
		},
		RunPr: &RunProperties{
			Bold: &Bold{},
			FontSize: &FontSize{
				Val: "56", // 28磅
			},
			FontFamily: &FontFamily{
				ASCII:    "Calibri Light",
				EastAsia: "微软雅黑 Light",
				HAnsi:    "Calibri Light",
				CS:       "Calibri Light",
			},
			Color: &Color{
				Val: "2F5496", // 深蓝色
			},
		},
	}
	sm.AddStyle(title)

	// 副标题样式
	subtitle := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Subtitle",
		Name: &StyleName{
			Val: "副标题",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Justification: &Justification{
				Val: "center", // 居中对齐
			},
			Spacing: &Spacing{
				Before: "0",   // 0磅段前间距
				After:  "160", // 8磅段后间距
			},
		},
		RunPr: &RunProperties{
			Italic: &Italic{},
			FontSize: &FontSize{
				Val: "30", // 15磅
			},
			FontFamily: &FontFamily{
				ASCII:    "Calibri",
				EastAsia: "微软雅黑",
				HAnsi:    "Calibri",
				CS:       "Calibri",
			},
			Color: &Color{
				Val: "7030A0", // 紫色
			},
		},
	}
	sm.AddStyle(subtitle)

	// 添加TOC样式（目录样式）
	sm.addTOCStyles()

	// 列表段落样式
	listParagraph := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "ListParagraph",
		Name: &StyleName{
			Val: "列表段落",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Indentation: &Indentation{
				Left: "720", // 左缩进0.5英寸（36磅）
			},
			Spacing: &Spacing{
				After:    "120", // 6磅段后间距
				Line:     "276", // 1.15倍行间距
				LineRule: "auto",
			},
		},
	}
	sm.AddStyle(listParagraph)

	// 强调样式
	emphasis := &Style{
		Type:    string(StyleTypeCharacter),
		StyleID: "Emphasis",
		Name: &StyleName{
			Val: "强调",
		},
		RunPr: &RunProperties{
			Italic: &Italic{},
		},
	}
	sm.AddStyle(emphasis)

	// 加粗样式
	strong := &Style{
		Type:    string(StyleTypeCharacter),
		StyleID: "Strong",
		Name: &StyleName{
			Val: "加粗",
		},
		RunPr: &RunProperties{
			Bold: &Bold{},
		},
	}
	sm.AddStyle(strong)

	// 引用样式
	quote := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "Quote",
		Name: &StyleName{
			Val: "引用",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Indentation: &Indentation{
				Left:  "720", // 左缩进0.5英寸
				Right: "720", // 右缩进0.5英寸
			},
			Spacing: &Spacing{
				Before: "120", // 6磅段前间距
				After:  "120", // 6磅段后间距
			},
		},
		RunPr: &RunProperties{
			Italic: &Italic{},
			Color: &Color{
				Val: "666666", // 中等灰色，更通用
			},
		},
	}
	sm.AddStyle(quote)

	// 代码块样式
	codeBlock := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "CodeBlock",
		Name: &StyleName{
			Val: "代码块",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Indentation: &Indentation{
				Left: "360", // 左缩进0.25英寸
			},
			Spacing: &Spacing{
				Before: "60", // 3磅段前间距
				After:  "60", // 3磅段后间距
			},
			ParagraphBorder: &ParagraphBorder{
				Top: &ParagraphBorderLine{
					Val:   "thick",
					Color: "E9E7E7",
					Sz:    "8",
					Space: "8",
				},
				Left: &ParagraphBorderLine{
					Val:   "thick",
					Color: "E9E7E7",
					Sz:    "8",
					Space: "8",
				},
				Bottom: &ParagraphBorderLine{
					Val:   "thick",
					Color: "E9E7E7",
					Sz:    "8",
					Space: "8",
				},
				Right: &ParagraphBorderLine{
					Val:   "thick",
					Color: "E9E7E7",
					Sz:    "8",
					Space: "8",
				},
			},
			Shading: &Shading{
				Fill: "F6F5F5",
				Val:  "clear",
			},
		},
		RunPr: &RunProperties{
			FontFamily: &FontFamily{
				ASCII:    "Consolas",
				EastAsia: "Consolas",
				HAnsi:    "Consolas",
				CS:       "Consolas",
			},
			FontSize: &FontSize{
				Val: "18", // 9磅，与code_template保持一致
			},
		},
	}
	sm.AddStyle(codeBlock)

	// 代码字符样式
	codeChar := &Style{
		Type:    string(StyleTypeCharacter),
		StyleID: "CodeChar",
		Name: &StyleName{
			Val: "代码字符",
		},
		RunPr: &RunProperties{
			FontFamily: &FontFamily{
				ASCII:    "Consolas",
				EastAsia: "Consolas",
				HAnsi:    "Consolas",
				CS:       "Consolas",
			},
			FontSize: &FontSize{
				Val: "18", // 9磅
			},
		},
	}
	sm.AddStyle(codeChar)

	// 添加表格样式
	sm.addTableStyles()
}

// addTOCStyles 添加TOC（目录）样式
func (sm *StyleManager) addTOCStyles() {
	// TOC 1 - 一级目录样式（无缩进）
	toc1 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "13", // TOC1 样式ID
		Name: &StyleName{
			Val: "toc 1",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Spacing: &Spacing{
				After: "100", // 5磅段后间距
			},
			Indentation: &Indentation{
				Left: "0", // 无左缩进
			},
		},
		RunPr: &RunProperties{
			FontSize: &FontSize{
				Val: "22", // 11磅
			},
			FontFamily: &FontFamily{
				ASCII:    "Calibri",
				EastAsia: "宋体",
				HAnsi:    "Calibri",
				CS:       "Times New Roman",
			},
		},
	}
	sm.AddStyle(toc1)

	// TOC 2 - 二级目录样式（左缩进240 TWIPs = 12磅）
	toc2 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "14", // TOC2 样式ID
		Name: &StyleName{
			Val: "toc 2",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Spacing: &Spacing{
				After: "100", // 5磅段后间距
			},
			Indentation: &Indentation{
				Left: "240", // 左缩进240 TWIPs (12磅)
			},
		},
		RunPr: &RunProperties{
			FontSize: &FontSize{
				Val: "22", // 11磅
			},
			FontFamily: &FontFamily{
				ASCII:    "Calibri",
				EastAsia: "宋体",
				HAnsi:    "Calibri",
				CS:       "Times New Roman",
			},
		},
	}
	sm.AddStyle(toc2)

	// TOC 3 - 三级目录样式（左缩进480 TWIPs = 24磅）
	toc3 := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "15", // TOC3 样式ID
		Name: &StyleName{
			Val: "toc 3",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Spacing: &Spacing{
				After: "100", // 5磅段后间距
			},
			Indentation: &Indentation{
				Left: "480", // 左缩进480 TWIPs (24磅)
			},
		},
		RunPr: &RunProperties{
			FontSize: &FontSize{
				Val: "22", // 11磅
			},
			FontFamily: &FontFamily{
				ASCII:    "Calibri",
				EastAsia: "宋体",
				HAnsi:    "Calibri",
				CS:       "Times New Roman",
			},
		},
	}
	sm.AddStyle(toc3)

	// TOC 4-9 - 四到九级目录样式
	for level := 4; level <= 9; level++ {
		styleID := fmt.Sprintf("%d", 12+level) // 16, 17, 18, 19, 20, 21
		tocStyle := &Style{
			Type:    string(StyleTypeParagraph),
			StyleID: styleID,
			Name: &StyleName{
				Val: fmt.Sprintf("toc %d", level),
			},
			BasedOn: &BasedOn{
				Val: "Normal",
			},
			Next: &Next{
				Val: "Normal",
			},
			ParagraphPr: &ParagraphProperties{
				Spacing: &Spacing{
					After: "100", // 5磅段后间距
				},
				Indentation: &Indentation{
					Left: fmt.Sprintf("%d", level*240), // 每级增加240 TWIPs (12磅)
				},
			},
			RunPr: &RunProperties{
				FontSize: &FontSize{
					Val: "22", // 11磅
				},
				FontFamily: &FontFamily{
					ASCII:    "Calibri",
					EastAsia: "宋体",
					HAnsi:    "Calibri",
					CS:       "Times New Roman",
				},
			},
		}
		sm.AddStyle(tocStyle)
	}

	// 添加基础TOC样式（样式ID为"12"的目录标题样式）
	tocBase := &Style{
		Type:    string(StyleTypeParagraph),
		StyleID: "12", // 基础TOC样式ID
		Name: &StyleName{
			Val: "TOCHeading",
		},
		BasedOn: &BasedOn{
			Val: "Normal",
		},
		Next: &Next{
			Val: "Normal",
		},
		ParagraphPr: &ParagraphProperties{
			Spacing: &Spacing{
				Before: "240", // 12磅段前间距
				After:  "120", // 6磅段后间距
			},
			Justification: &Justification{
				Val: "center", // 居中对齐
			},
		},
		RunPr: &RunProperties{
			Bold: &Bold{},
			FontSize: &FontSize{
				Val: "26", // 13磅
			},
			FontFamily: &FontFamily{
				ASCII:    "Calibri",
				EastAsia: "宋体",
				HAnsi:    "Calibri",
				CS:       "Times New Roman",
			},
		},
	}
	sm.AddStyle(tocBase)
}

// GetStyleWithInheritance 获取具有继承属性的样式
// 如果样式基于其他样式，会合并父样式的属性
func (sm *StyleManager) GetStyleWithInheritance(styleID string) *Style {
	style := sm.GetStyle(styleID)
	if style == nil {
		return nil
	}

	// 如果样式没有基础样式，直接返回
	if style.BasedOn == nil {
		return style
	}

	// 递归获取基础样式
	baseStyle := sm.GetStyleWithInheritance(style.BasedOn.Val)
	if baseStyle == nil {
		return style
	}

	// 创建合并后的样式副本
	mergedStyle := &Style{
		Type:        style.Type,
		StyleID:     style.StyleID,
		Name:        style.Name,
		BasedOn:     style.BasedOn,
		Next:        style.Next,
		Default:     style.Default,
		CustomStyle: style.CustomStyle,
	}

	// 合并段落属性
	mergedStyle.ParagraphPr = mergeParagraphProperties(baseStyle.ParagraphPr, style.ParagraphPr)

	// 合并字符属性
	mergedStyle.RunPr = mergeRunProperties(baseStyle.RunPr, style.RunPr)

	// 合并表格属性（如果有）
	if style.TablePr != nil {
		mergedStyle.TablePr = style.TablePr
	} else if baseStyle.TablePr != nil {
		mergedStyle.TablePr = baseStyle.TablePr
	}

	return mergedStyle
}

// mergeParagraphProperties 合并段落属性，优先级：override > base
func mergeParagraphProperties(base, override *ParagraphProperties) *ParagraphProperties {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	merged := &ParagraphProperties{}

	// 合并间距
	if override.Spacing != nil {
		merged.Spacing = override.Spacing
	} else if base.Spacing != nil {
		merged.Spacing = base.Spacing
	}

	// 合并对齐
	if override.Justification != nil {
		merged.Justification = override.Justification
	} else if base.Justification != nil {
		merged.Justification = base.Justification
	}

	// 合并缩进
	if override.Indentation != nil {
		merged.Indentation = override.Indentation
	} else if base.Indentation != nil {
		merged.Indentation = base.Indentation
	}

	// 合并边框
	if override.ParagraphBorder != nil {
		merged.ParagraphBorder = override.ParagraphBorder
	} else if base.ParagraphBorder != nil {
		merged.ParagraphBorder = base.ParagraphBorder
	}

	// 合并阴影
	if override.Shading != nil {
		merged.Shading = override.Shading
	} else if base.Shading != nil {
		merged.Shading = base.Shading
	}

	// 合并其他属性
	if override.KeepNext != nil {
		merged.KeepNext = override.KeepNext
	} else if base.KeepNext != nil {
		merged.KeepNext = base.KeepNext
	}

	if override.KeepLines != nil {
		merged.KeepLines = override.KeepLines
	} else if base.KeepLines != nil {
		merged.KeepLines = base.KeepLines
	}

	if override.PageBreak != nil {
		merged.PageBreak = override.PageBreak
	} else if base.PageBreak != nil {
		merged.PageBreak = base.PageBreak
	}

	if override.OutlineLevel != nil {
		merged.OutlineLevel = override.OutlineLevel
	} else if base.OutlineLevel != nil {
		merged.OutlineLevel = base.OutlineLevel
	}

	return merged
}

// mergeRunProperties 合并字符属性
func mergeRunProperties(base, override *RunProperties) *RunProperties {
	if base == nil {
		return override
	}
	if override == nil {
		return base
	}

	merged := &RunProperties{}

	// 合并文字格式
	if override.Bold != nil {
		merged.Bold = override.Bold
	} else if base.Bold != nil {
		merged.Bold = base.Bold
	}

	if override.Italic != nil {
		merged.Italic = override.Italic
	} else if base.Italic != nil {
		merged.Italic = base.Italic
	}

	if override.Underline != nil {
		merged.Underline = override.Underline
	} else if base.Underline != nil {
		merged.Underline = base.Underline
	}

	if override.Strike != nil {
		merged.Strike = override.Strike
	} else if base.Strike != nil {
		merged.Strike = base.Strike
	}

	// 合并字体属性
	if override.FontSize != nil {
		merged.FontSize = override.FontSize
	} else if base.FontSize != nil {
		merged.FontSize = base.FontSize
	}

	if override.Color != nil {
		merged.Color = override.Color
	} else if base.Color != nil {
		merged.Color = base.Color
	}

	if override.FontFamily != nil {
		merged.FontFamily = override.FontFamily
	} else if base.FontFamily != nil {
		merged.FontFamily = base.FontFamily
	}

	if override.Highlight != nil {
		merged.Highlight = override.Highlight
	} else if base.Highlight != nil {
		merged.Highlight = base.Highlight
	}

	return merged
}

// CreateCustomStyle 创建自定义样式
func (sm *StyleManager) CreateCustomStyle(styleID, name string, styleType StyleType, basedOn string) *Style {
	style := &Style{
		Type:        string(styleType),
		StyleID:     styleID,
		CustomStyle: true,
		Name: &StyleName{
			Val: name,
		},
	}

	if basedOn != "" {
		style.BasedOn = &BasedOn{
			Val: basedOn,
		}
	}

	sm.AddStyle(style)
	return style
}

// RemoveStyle 移除样式
func (sm *StyleManager) RemoveStyle(styleID string) {
	delete(sm.styles, styleID)
}

// StyleExists 检查样式是否存在
func (sm *StyleManager) StyleExists(styleID string) bool {
	_, exists := sm.styles[styleID]
	return exists
}

// Clone 深拷贝样式管理器，用于模板渲染时避免样式冲突
func (sm *StyleManager) Clone() *StyleManager {
	clonedSM := &StyleManager{
		styles: make(map[string]*Style),
	}

	// 深拷贝所有样式
	for styleID, style := range sm.styles {
		clonedSM.styles[styleID] = sm.cloneStyle(style)
	}

	return clonedSM
}

// cloneStyle 深拷贝单个样式
func (sm *StyleManager) cloneStyle(source *Style) *Style {
	if source == nil {
		return nil
	}

	cloned := &Style{
		Type:        source.Type,
		StyleID:     source.StyleID,
		Default:     source.Default,
		CustomStyle: source.CustomStyle,
	}

	// 克隆样式名称
	if source.Name != nil {
		cloned.Name = &StyleName{Val: source.Name.Val}
	}

	// 克隆基于样式
	if source.BasedOn != nil {
		cloned.BasedOn = &BasedOn{Val: source.BasedOn.Val}
	}

	// 克隆下一样式
	if source.Next != nil {
		cloned.Next = &Next{Val: source.Next.Val}
	}

	// 克隆段落属性
	if source.ParagraphPr != nil {
		cloned.ParagraphPr = sm.cloneParagraphProperties(source.ParagraphPr)
	}

	// 克隆字符属性
	if source.RunPr != nil {
		cloned.RunPr = sm.cloneRunProperties(source.RunPr)
	}

	// 克隆表格属性
	if source.TablePr != nil {
		cloned.TablePr = sm.cloneTableProperties(source.TablePr)
	}

	// 克隆表格行属性
	if source.TableRowPr != nil {
		cloned.TableRowPr = sm.cloneTableRowProperties(source.TableRowPr)
	}

	// 克隆表格单元格属性
	if source.TableCellPr != nil {
		cloned.TableCellPr = sm.cloneTableCellProperties(source.TableCellPr)
	}

	return cloned
}

// cloneParagraphProperties 深度复制段落属性
func (sm *StyleManager) cloneParagraphProperties(source *ParagraphProperties) *ParagraphProperties {
	if source == nil {
		return nil
	}

	cloned := &ParagraphProperties{}

	// 复制间距
	if source.Spacing != nil {
		cloned.Spacing = &Spacing{
			Before:   source.Spacing.Before,
			After:    source.Spacing.After,
			Line:     source.Spacing.Line,
			LineRule: source.Spacing.LineRule,
		}
	}

	// 复制对齐方式
	if source.Justification != nil {
		cloned.Justification = &Justification{
			Val: source.Justification.Val,
		}
	}

	// 复制缩进
	if source.Indentation != nil {
		cloned.Indentation = &Indentation{
			FirstLine: source.Indentation.FirstLine,
			Left:      source.Indentation.Left,
			Right:     source.Indentation.Right,
		}
	}

	// 复制段落边框
	if source.ParagraphBorder != nil {
		cloned.ParagraphBorder = &ParagraphBorder{}
		if source.ParagraphBorder.Top != nil {
			cloned.ParagraphBorder.Top = &ParagraphBorderLine{
				Val:   source.ParagraphBorder.Top.Val,
				Color: source.ParagraphBorder.Top.Color,
				Sz:    source.ParagraphBorder.Top.Sz,
				Space: source.ParagraphBorder.Top.Space,
			}
		}
		if source.ParagraphBorder.Left != nil {
			cloned.ParagraphBorder.Left = &ParagraphBorderLine{
				Val:   source.ParagraphBorder.Left.Val,
				Color: source.ParagraphBorder.Left.Color,
				Sz:    source.ParagraphBorder.Left.Sz,
				Space: source.ParagraphBorder.Left.Space,
			}
		}
		if source.ParagraphBorder.Bottom != nil {
			cloned.ParagraphBorder.Bottom = &ParagraphBorderLine{
				Val:   source.ParagraphBorder.Bottom.Val,
				Color: source.ParagraphBorder.Bottom.Color,
				Sz:    source.ParagraphBorder.Bottom.Sz,
				Space: source.ParagraphBorder.Bottom.Space,
			}
		}
		if source.ParagraphBorder.Right != nil {
			cloned.ParagraphBorder.Right = &ParagraphBorderLine{
				Val:   source.ParagraphBorder.Right.Val,
				Color: source.ParagraphBorder.Right.Color,
				Sz:    source.ParagraphBorder.Right.Sz,
				Space: source.ParagraphBorder.Right.Space,
			}
		}
	}

	// 复制阴影
	if source.Shading != nil {
		cloned.Shading = &Shading{
			Fill: source.Shading.Fill,
			Val:  source.Shading.Val,
		}
	}

	// 复制其他属性
	if source.KeepNext != nil {
		cloned.KeepNext = &KeepNext{}
	}

	if source.KeepLines != nil {
		cloned.KeepLines = &KeepLines{}
	}

	if source.PageBreak != nil {
		cloned.PageBreak = &PageBreak{}
	}

	if source.OutlineLevel != nil {
		cloned.OutlineLevel = &OutlineLevel{
			Val: source.OutlineLevel.Val,
		}
	}

	// 复制网格对齐设置
	if source.SnapToGrid != nil {
		cloned.SnapToGrid = &SnapToGrid{
			Val: source.SnapToGrid.Val,
		}
	}

	return cloned
}

// cloneRunProperties 深拷贝字符属性
func (sm *StyleManager) cloneRunProperties(source *RunProperties) *RunProperties {
	if source == nil {
		return nil
	}

	cloned := &RunProperties{}

	// 克隆字体格式
	if source.Bold != nil {
		cloned.Bold = &Bold{}
	}

	if source.Italic != nil {
		cloned.Italic = &Italic{}
	}

	if source.Underline != nil {
		cloned.Underline = &Underline{Val: source.Underline.Val}
	}

	if source.Strike != nil {
		cloned.Strike = &Strike{}
	}

	// 克隆字体大小
	if source.FontSize != nil {
		cloned.FontSize = &FontSize{Val: source.FontSize.Val}
	}

	// 克隆颜色
	if source.Color != nil {
		cloned.Color = &Color{Val: source.Color.Val}
	}

	// 克隆字体族
	if source.FontFamily != nil {
		cloned.FontFamily = &FontFamily{
			ASCII:    source.FontFamily.ASCII,
			EastAsia: source.FontFamily.EastAsia,
			HAnsi:    source.FontFamily.HAnsi,
			CS:       source.FontFamily.CS,
		}
	}

	// 克隆高亮
	if source.Highlight != nil {
		cloned.Highlight = &Highlight{Val: source.Highlight.Val}
	}

	return cloned
}

// cloneTableProperties 深拷贝表格属性
func (sm *StyleManager) cloneTableProperties(source *TableProperties) *TableProperties {
	if source == nil {
		return nil
	}

	cloned := &TableProperties{}

	// 克隆表格缩进
	if source.TblInd != nil {
		cloned.TblInd = &TblIndent{
			W:    source.TblInd.W,
			Type: source.TblInd.Type,
		}
	}

	// 克隆表格边框
	if source.TblBorders != nil {
		cloned.TblBorders = &TblBorders{}

		if source.TblBorders.Top != nil {
			cloned.TblBorders.Top = &TblBorder{
				Val:   source.TblBorders.Top.Val,
				Sz:    source.TblBorders.Top.Sz,
				Space: source.TblBorders.Top.Space,
				Color: source.TblBorders.Top.Color,
			}
		}

		if source.TblBorders.Left != nil {
			cloned.TblBorders.Left = &TblBorder{
				Val:   source.TblBorders.Left.Val,
				Sz:    source.TblBorders.Left.Sz,
				Space: source.TblBorders.Left.Space,
				Color: source.TblBorders.Left.Color,
			}
		}

		if source.TblBorders.Bottom != nil {
			cloned.TblBorders.Bottom = &TblBorder{
				Val:   source.TblBorders.Bottom.Val,
				Sz:    source.TblBorders.Bottom.Sz,
				Space: source.TblBorders.Bottom.Space,
				Color: source.TblBorders.Bottom.Color,
			}
		}

		if source.TblBorders.Right != nil {
			cloned.TblBorders.Right = &TblBorder{
				Val:   source.TblBorders.Right.Val,
				Sz:    source.TblBorders.Right.Sz,
				Space: source.TblBorders.Right.Space,
				Color: source.TblBorders.Right.Color,
			}
		}

		if source.TblBorders.InsideH != nil {
			cloned.TblBorders.InsideH = &TblBorder{
				Val:   source.TblBorders.InsideH.Val,
				Sz:    source.TblBorders.InsideH.Sz,
				Space: source.TblBorders.InsideH.Space,
				Color: source.TblBorders.InsideH.Color,
			}
		}

		if source.TblBorders.InsideV != nil {
			cloned.TblBorders.InsideV = &TblBorder{
				Val:   source.TblBorders.InsideV.Val,
				Sz:    source.TblBorders.InsideV.Sz,
				Space: source.TblBorders.InsideV.Space,
				Color: source.TblBorders.InsideV.Color,
			}
		}
	}

	// 克隆表格单元格边距
	if source.TblCellMar != nil {
		cloned.TblCellMar = &TblCellMargin{}

		if source.TblCellMar.Top != nil {
			cloned.TblCellMar.Top = &TblCellSpace{
				W:    source.TblCellMar.Top.W,
				Type: source.TblCellMar.Top.Type,
			}
		}

		if source.TblCellMar.Left != nil {
			cloned.TblCellMar.Left = &TblCellSpace{
				W:    source.TblCellMar.Left.W,
				Type: source.TblCellMar.Left.Type,
			}
		}

		if source.TblCellMar.Bottom != nil {
			cloned.TblCellMar.Bottom = &TblCellSpace{
				W:    source.TblCellMar.Bottom.W,
				Type: source.TblCellMar.Bottom.Type,
			}
		}

		if source.TblCellMar.Right != nil {
			cloned.TblCellMar.Right = &TblCellSpace{
				W:    source.TblCellMar.Right.W,
				Type: source.TblCellMar.Right.Type,
			}
		}
	}

	return cloned
}

// cloneTableRowProperties 深拷贝表格行属性
func (sm *StyleManager) cloneTableRowProperties(source *TableRowProperties) *TableRowProperties {
	if source == nil {
		return nil
	}

	// 目前表格行属性是空结构体，简单返回新实例
	cloned := &TableRowProperties{}

	return cloned
}

// cloneTableCellProperties 深拷贝表格单元格属性
func (sm *StyleManager) cloneTableCellProperties(source *TableCellProperties) *TableCellProperties {
	if source == nil {
		return nil
	}

	// 目前表格单元格属性是空结构体，简单返回新实例
	cloned := &TableCellProperties{}

	return cloned
}

// GetStylesByType 按类型获取样式
func (sm *StyleManager) GetStylesByType(styleType StyleType) []*Style {
	var styles []*Style
	for _, style := range sm.styles {
		if StyleType(style.Type) == styleType {
			styles = append(styles, style)
		}
	}
	return styles
}

// GetHeadingStyles 获取所有标题样式
func (sm *StyleManager) GetHeadingStyles() []*Style {
	var headingStyles []*Style
	for i := 1; i <= 9; i++ {
		styleID := fmt.Sprintf("Heading%d", i)
		if style := sm.GetStyle(styleID); style != nil {
			headingStyles = append(headingStyles, style)
		}
	}
	return headingStyles
}

// ApplyStyleToXML 将样式应用到XML结构（为文档集成做准备）
func (sm *StyleManager) ApplyStyleToXML(styleID string) (map[string]interface{}, error) {
	style := sm.GetStyleWithInheritance(styleID)
	if style == nil {
		return nil, fmt.Errorf("style %s not found", styleID)
	}

	result := make(map[string]interface{})
	result["styleId"] = style.StyleID
	result["type"] = style.Type

	if style.ParagraphPr != nil {
		result["paragraphProperties"] = convertParagraphPropertiesToMap(style.ParagraphPr)
	}

	if style.RunPr != nil {
		result["runProperties"] = convertRunPropertiesToMap(style.RunPr)
	}

	return result, nil
}

// convertParagraphPropertiesToMap 将段落属性转换为映射
func convertParagraphPropertiesToMap(props *ParagraphProperties) map[string]interface{} {
	result := make(map[string]interface{})

	if props.Spacing != nil {
		spacing := make(map[string]string)
		if props.Spacing.Before != "" {
			spacing["before"] = props.Spacing.Before
		}
		if props.Spacing.After != "" {
			spacing["after"] = props.Spacing.After
		}
		if props.Spacing.Line != "" {
			spacing["line"] = props.Spacing.Line
		}
		if props.Spacing.LineRule != "" {
			spacing["lineRule"] = props.Spacing.LineRule
		}
		result["spacing"] = spacing
	}

	if props.Justification != nil {
		result["justification"] = props.Justification.Val
	}

	if props.Indentation != nil {
		indentation := make(map[string]string)
		if props.Indentation.FirstLine != "" {
			indentation["firstLine"] = props.Indentation.FirstLine
		}
		if props.Indentation.Left != "" {
			indentation["left"] = props.Indentation.Left
		}
		if props.Indentation.Right != "" {
			indentation["right"] = props.Indentation.Right
		}
		result["indentation"] = indentation
	}

	if props.OutlineLevel != nil {
		result["outlineLevel"] = props.OutlineLevel.Val
	}

	return result
}

// convertRunPropertiesToMap 将字符属性转换为映射
func convertRunPropertiesToMap(props *RunProperties) map[string]interface{} {
	result := make(map[string]interface{})

	if props.Bold != nil {
		result["bold"] = true
	}

	if props.Italic != nil {
		result["italic"] = true
	}

	if props.Underline != nil {
		result["underline"] = props.Underline.Val
	}

	if props.Strike != nil {
		result["strike"] = true
	}

	if props.FontSize != nil {
		result["fontSize"] = props.FontSize.Val
	}

	if props.Color != nil {
		result["color"] = props.Color.Val
	}

	if props.FontFamily != nil {
		fontFamily := make(map[string]string)
		if props.FontFamily.ASCII != "" {
			fontFamily["ascii"] = props.FontFamily.ASCII
		}
		if props.FontFamily.EastAsia != "" {
			fontFamily["eastAsia"] = props.FontFamily.EastAsia
		}
		if props.FontFamily.HAnsi != "" {
			fontFamily["hAnsi"] = props.FontFamily.HAnsi
		}
		if props.FontFamily.CS != "" {
			fontFamily["cs"] = props.FontFamily.CS
		}
		result["fontFamily"] = fontFamily
	}

	if props.Highlight != nil {
		result["highlight"] = props.Highlight.Val
	}

	return result
}

// ParseStylesFromXML 从XML数据解析样式
func (sm *StyleManager) ParseStylesFromXML(xmlData []byte) error {
	type stylesXML struct {
		XMLName xml.Name `xml:"w:styles"`
		Styles  []Style  `xml:"w:style"`
	}

	var styles stylesXML
	if err := xml.Unmarshal(xmlData, &styles); err != nil {
		return fmt.Errorf("解析样式XML失败: %v", err)
	}

	// 清空现有样式（除非我们想要合并）
	sm.styles = make(map[string]*Style)

	// 添加解析的样式
	for i := range styles.Styles {
		sm.AddStyle(&styles.Styles[i])
	}

	return nil
}

// MergeStylesFromXML 从XML数据合并样式（保留现有样式，只添加新的）
func (sm *StyleManager) MergeStylesFromXML(xmlData []byte) error {
	type stylesXML struct {
		XMLName xml.Name `xml:"w:styles"`
		Styles  []Style  `xml:"w:style"`
	}

	var styles stylesXML
	if err := xml.Unmarshal(xmlData, &styles); err != nil {
		return fmt.Errorf("解析样式XML失败: %v", err)
	}

	// 只添加不存在的样式，保留现有样式
	for i := range styles.Styles {
		if !sm.StyleExists(styles.Styles[i].StyleID) {
			sm.AddStyle(&styles.Styles[i])
		}
	}

	return nil
}

// LoadStylesFromDocument 从现有文档加载样式，优先保留原有样式设置
func (sm *StyleManager) LoadStylesFromDocument(xmlData []byte) error {
	if len(xmlData) == 0 {
		// 如果没有样式数据，使用默认样式
		sm.initializePredefinedStyles()
		return nil
	}

	// 先解析现有样式
	if err := sm.ParseStylesFromXML(xmlData); err != nil {
		// 如果解析失败，使用默认样式
		sm.initializePredefinedStyles()
		return fmt.Errorf("解析现有样式失败，使用默认样式: %v", err)
	}

	// 确保基本样式存在，如果不存在则添加
	if !sm.StyleExists("Normal") {
		sm.addNormalStyle()
	}

	// 确保基本标题样式存在
	headingStyles := []string{"Heading1", "Heading2", "Heading3", "Heading4", "Heading5", "Heading6", "Heading7", "Heading8", "Heading9"}
	for _, styleID := range headingStyles {
		if !sm.StyleExists(styleID) {
			sm.addHeadingStyles()
			break
		}
	}

	return nil
}

// addTableStyles 添加表格样式
func (sm *StyleManager) addTableStyles() {
	// 普通表格样式 (Normal Table) - 默认表格样式
	normalTable := &Style{
		Type:    string(StyleTypeTable),
		StyleID: "a1",
		Default: true,
		Name: &StyleName{
			Val: "Normal Table",
		},
		TablePr: &TableProperties{
			TblInd: &TblIndent{
				W:    "0",
				Type: "dxa",
			},
			TblCellMar: &TblCellMargin{
				Top: &TblCellSpace{
					W:    "0",
					Type: "dxa",
				},
				Left: &TblCellSpace{
					W:    "108",
					Type: "dxa",
				},
				Bottom: &TblCellSpace{
					W:    "0",
					Type: "dxa",
				},
				Right: &TblCellSpace{
					W:    "108",
					Type: "dxa",
				},
			},
		},
	}
	sm.AddStyle(normalTable)

	// 表格网格样式 (Table Grid) - 基于Normal Table，添加边框
	tableGrid := &Style{
		Type:    string(StyleTypeTable),
		StyleID: "ab",
		Name: &StyleName{
			Val: "Table Grid",
		},
		BasedOn: &BasedOn{
			Val: "a1",
		},
		TablePr: &TableProperties{
			TblBorders: &TblBorders{
				Top: &TblBorder{
					Val:   "single",
					Sz:    "4",
					Space: "0",
					Color: "auto",
				},
				Left: &TblBorder{
					Val:   "single",
					Sz:    "4",
					Space: "0",
					Color: "auto",
				},
				Bottom: &TblBorder{
					Val:   "single",
					Sz:    "4",
					Space: "0",
					Color: "auto",
				},
				Right: &TblBorder{
					Val:   "single",
					Sz:    "4",
					Space: "0",
					Color: "auto",
				},
				InsideH: &TblBorder{
					Val:   "single",
					Sz:    "4",
					Space: "0",
					Color: "auto",
				},
				InsideV: &TblBorder{
					Val:   "single",
					Sz:    "4",
					Space: "0",
					Color: "auto",
				},
			},
		},
	}
	sm.AddStyle(tableGrid)
}

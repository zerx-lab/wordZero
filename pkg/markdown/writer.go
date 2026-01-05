package markdown

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/zerx-lab/wordZero/pkg/document"
)

// ExportOptions 导出选项配置
type ExportOptions struct {
	// 基础配置
	UseGFMTables       bool // 使用GFM表格语法
	PreserveFootnotes  bool // 保留脚注
	PreserveLineBreaks bool // 保留换行符
	WrapLongLines      bool // 自动换行长行
	MaxLineLength      int  // 最大行长度

	// 图片处理
	ExtractImages     bool   // 是否导出图片文件
	ImageOutputDir    string // 图片输出目录
	ImageNamePattern  string // 图片命名模式
	ImageRelativePath bool   // 使用相对路径引用图片

	// 链接处理
	PreserveBookmarks bool // 保留书签为锚点链接
	ConvertHyperlinks bool // 转换超链接

	// 代码块处理
	PreserveCodeStyle bool   // 保留代码样式
	DefaultCodeLang   string // 默认代码语言标识

	// 样式映射
	CustomStyleMap      map[string]string // 自定义样式映射
	IgnoreUnknownStyles bool              // 忽略未知样式

	// 内容处理
	PreserveTOC     bool // 保留目录
	IncludeMetadata bool // 包含文档元数据
	StripComments   bool // 移除注释

	// 格式化选项
	UseSetext        bool   // 使用Setext样式标题
	BulletListMarker string // 项目符号标记
	EmphasisMarker   string // 强调标记

	// 错误处理
	StrictMode    bool        // 严格模式
	IgnoreErrors  bool        // 忽略转换错误
	ErrorCallback func(error) // 错误回调

	// 进度报告
	ProgressCallback func(int, int) // 进度回调
}

// MarkdownWriter Markdown格式输出器
type MarkdownWriter struct {
	opts      *ExportOptions
	doc       *document.Document
	output    strings.Builder
	imageNum  int
	footnotes []string
}

// Write 生成Markdown内容
func (w *MarkdownWriter) Write() ([]byte, error) {
	// 处理文档元数据
	if w.opts.IncludeMetadata {
		w.writeMetadata()
	}

	// 遍历文档段落
	if w.doc.Body != nil {
		for _, para := range w.doc.Body.GetParagraphs() {
			err := w.writeParagraph(para)
			if err != nil {
				if w.opts.ErrorCallback != nil {
					w.opts.ErrorCallback(err)
				}
				if !w.opts.IgnoreErrors {
					return nil, err
				}
			}
		}

		// 处理表格
		for _, table := range w.doc.Body.GetTables() {
			err := w.writeTable(table)
			if err != nil {
				if w.opts.ErrorCallback != nil {
					w.opts.ErrorCallback(err)
				}
				if !w.opts.IgnoreErrors {
					return nil, err
				}
			}
		}
	}

	// 添加脚注
	if w.opts.PreserveFootnotes && len(w.footnotes) > 0 {
		w.writeFootnotes()
	}

	return []byte(w.output.String()), nil
}

// writeMetadata 写入文档元数据
func (w *MarkdownWriter) writeMetadata() {
	w.output.WriteString("---\n")
	w.output.WriteString("title: \"Document\"\n")
	w.output.WriteString("---\n\n")
}

// writeParagraph 写入段落
func (w *MarkdownWriter) writeParagraph(para *document.Paragraph) error {
	if para == nil {
		return nil
	}

	// 检查段落样式
	style := w.getParagraphStyle(para)

	switch {
	case strings.HasPrefix(style, "Heading"):
		return w.writeHeading(para, style)
	case style == "Quote":
		return w.writeQuote(para)
	case style == "CodeBlock":
		return w.writeCodeBlock(para)
	case w.isListParagraph(para):
		return w.writeListItem(para)
	default:
		return w.writeNormalParagraph(para)
	}
}

// writeHeading 写入标题
func (w *MarkdownWriter) writeHeading(para *document.Paragraph, style string) error {
	level := w.getHeadingLevel(style)
	if level > 6 {
		level = 6
	}

	text := w.extractParagraphText(para)
	if strings.TrimSpace(text) == "" {
		return nil
	}

	if w.opts.UseSetext && level <= 2 {
		// 使用Setext样式
		w.output.WriteString(text + "\n")
		if level == 1 {
			w.output.WriteString(strings.Repeat("=", len(text)) + "\n\n")
		} else {
			w.output.WriteString(strings.Repeat("-", len(text)) + "\n\n")
		}
	} else {
		// 使用ATX样式
		w.output.WriteString(strings.Repeat("#", level) + " " + text + "\n\n")
	}

	return nil
}

// writeQuote 写入引用
func (w *MarkdownWriter) writeQuote(para *document.Paragraph) error {
	text := w.extractParagraphText(para)
	if strings.TrimSpace(text) == "" {
		return nil
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		w.output.WriteString("> " + line + "\n")
	}
	w.output.WriteString("\n")

	return nil
}

// writeCodeBlock 写入代码块
func (w *MarkdownWriter) writeCodeBlock(para *document.Paragraph) error {
	text := w.extractParagraphText(para)
	if strings.TrimSpace(text) == "" {
		return nil
	}

	lang := w.opts.DefaultCodeLang
	w.output.WriteString("```" + lang + "\n")
	w.output.WriteString(text + "\n")
	w.output.WriteString("```\n\n")

	return nil
}

// writeListItem 写入列表项
func (w *MarkdownWriter) writeListItem(para *document.Paragraph) error {
	text := w.extractParagraphText(para)
	if strings.TrimSpace(text) == "" {
		return nil
	}

	// 简单的列表项处理
	marker := w.opts.BulletListMarker
	if w.isNumberedList(para) {
		marker = "1."
	}

	w.output.WriteString(marker + " " + text + "\n")

	return nil
}

// writeNormalParagraph 写入普通段落
func (w *MarkdownWriter) writeNormalParagraph(para *document.Paragraph) error {
	text := w.extractParagraphText(para)
	if strings.TrimSpace(text) == "" {
		w.output.WriteString("\n")
		return nil
	}

	// 处理长行换行
	if w.opts.WrapLongLines && len(text) > w.opts.MaxLineLength {
		text = w.wrapText(text, w.opts.MaxLineLength)
	}

	w.output.WriteString(text + "\n\n")

	return nil
}

// writeTable 写入表格
func (w *MarkdownWriter) writeTable(table *document.Table) error {
	if table == nil || len(table.Rows) == 0 {
		return nil
	}

	if !w.opts.UseGFMTables {
		return w.writeSimpleTable(table)
	}

	rows := table.Rows

	// 写表头
	if len(rows) > 0 {
		headerRow := rows[0]
		w.output.WriteString("|")
		for _, cell := range headerRow.Cells {
			text := w.extractCellText(&cell)
			w.output.WriteString(" " + text + " |")
		}
		w.output.WriteString("\n")

		// 写分隔行
		w.output.WriteString("|")
		for i := 0; i < len(headerRow.Cells); i++ {
			w.output.WriteString("-----|")
		}
		w.output.WriteString("\n")

		// 写数据行
		for i := 1; i < len(rows); i++ {
			w.output.WriteString("|")
			for _, cell := range rows[i].Cells {
				text := w.extractCellText(&cell)
				w.output.WriteString(" " + text + " |")
			}
			w.output.WriteString("\n")
		}
	}

	w.output.WriteString("\n")

	return nil
}

// writeSimpleTable 写入简单表格格式
func (w *MarkdownWriter) writeSimpleTable(table *document.Table) error {
	for i, row := range table.Rows {
		if i == 0 {
			w.output.WriteString("**")
		}
		for j, cell := range row.Cells {
			if j > 0 {
				w.output.WriteString(" | ")
			}
			text := w.extractCellText(&cell)
			w.output.WriteString(text)
		}
		if i == 0 {
			w.output.WriteString("**")
		}
		w.output.WriteString("\n")
	}
	w.output.WriteString("\n")

	return nil
}

// writeFootnotes 写入脚注
func (w *MarkdownWriter) writeFootnotes() {
	w.output.WriteString("\n---\n\n")
	for i, footnote := range w.footnotes {
		w.output.WriteString(fmt.Sprintf("[^%d]: %s\n", i+1, footnote))
	}
}

// extractParagraphText 提取段落文本
func (w *MarkdownWriter) extractParagraphText(para *document.Paragraph) string {
	if para == nil {
		return ""
	}

	var result strings.Builder

	for _, run := range para.Runs {
		text := w.formatRunText(&run)
		result.WriteString(text)
	}

	return result.String()
}

// formatRunText 格式化文本运行
func (w *MarkdownWriter) formatRunText(run *document.Run) string {
	if run == nil {
		return ""
	}

	text := run.Text.Content
	if text == "" {
		return ""
	}

	// 检查格式属性
	if run.Properties != nil {
		// 检查粗体
		if run.Properties.Bold != nil {
			if run.Properties.Italic != nil {
				text = "***" + text + "***" // 粗斜体
			} else {
				text = "**" + text + "**" // 粗体
			}
		} else if run.Properties.Italic != nil {
			text = w.opts.EmphasisMarker + text + w.opts.EmphasisMarker // 斜体
		}

		// 检查删除线
		if run.Properties.Strike != nil {
			text = "~~" + text + "~~" // 删除线
		}

		// 处理代码样式
		if w.isCodeStyle(run.Properties) {
			text = "`" + text + "`"
		}
	}

	return text
}

// extractCellText 提取单元格文本
func (w *MarkdownWriter) extractCellText(cell *document.TableCell) string {
	if cell == nil {
		return ""
	}

	var result strings.Builder

	for _, para := range cell.Paragraphs {
		text := w.extractParagraphText(&para)
		result.WriteString(text)
	}

	// 清理表格单元格中的换行符
	text := result.String()
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.TrimSpace(text)

	return text
}

// getParagraphStyle 获取段落样式
func (w *MarkdownWriter) getParagraphStyle(para *document.Paragraph) string {
	if para.Properties != nil && para.Properties.ParagraphStyle != nil {
		return para.Properties.ParagraphStyle.Val
	}
	return "Normal"
}

// getHeadingLevel 获取标题级别
func (w *MarkdownWriter) getHeadingLevel(style string) int {
	// 提取数字
	re := regexp.MustCompile(`\d+`)
	matches := re.FindString(style)
	if matches != "" {
		if level, err := strconv.Atoi(matches); err == nil {
			return level
		}
	}
	return 1
}

// isListParagraph 判断是否为列表段落
func (w *MarkdownWriter) isListParagraph(para *document.Paragraph) bool {
	if para.Properties == nil {
		return false
	}
	return para.Properties.NumberingProperties != nil
}

// isNumberedList 判断是否为编号列表
func (w *MarkdownWriter) isNumberedList(para *document.Paragraph) bool {
	// 简单实现，实际应该检查编号格式
	return false
}

// isCodeStyle 判断是否为代码样式
func (w *MarkdownWriter) isCodeStyle(props *document.RunProperties) bool {
	if props.FontFamily != nil {
		font := props.FontFamily.ASCII
		// 检查是否为等宽字体
		codefonts := []string{"Consolas", "Courier New", "Monaco", "Menlo", "Source Code Pro"}
		for _, codefont := range codefonts {
			if strings.Contains(font, codefont) {
				return true
			}
		}
	}
	return false
}

// wrapText 文本换行
func (w *MarkdownWriter) wrapText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	var line strings.Builder

	for _, word := range words {
		if line.Len()+len(word)+1 > maxLength {
			if line.Len() > 0 {
				result.WriteString(line.String() + "\n")
				line.Reset()
			}
		}
		if line.Len() > 0 {
			line.WriteString(" ")
		}
		line.WriteString(word)
	}

	if line.Len() > 0 {
		result.WriteString(line.String())
	}

	return result.String()
}

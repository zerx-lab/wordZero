package document

import (
	"fmt"
	"regexp"
	"strings"
)

// TemplateRenderer 专门负责模板渲染的引擎
type TemplateRenderer struct {
	engine *TemplateEngine
	logger *TemplateLogger
}

// TemplateLogger 模板日志记录器
type TemplateLogger struct {
	enabled bool
}

// NewTemplateRenderer 创建新的模板渲染器
func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		engine: NewTemplateEngine(),
		logger: &TemplateLogger{enabled: true},
	}
}

// SetLogging 设置是否启用日志记录
func (tr *TemplateRenderer) SetLogging(enabled bool) {
	tr.logger.enabled = enabled
}

// logInfo 记录信息日志
func (tr *TemplateRenderer) logInfo(format string, args ...interface{}) {
	if tr.logger.enabled {
		Infof("[模板引擎] "+format, args...)
	}
}

// logError 记录错误日志
func (tr *TemplateRenderer) logError(format string, args ...interface{}) {
	if tr.logger.enabled {
		Errorf("[模板引擎] "+format, args...)
	}
}

// LoadTemplateFromFile 从文件加载模板
func (tr *TemplateRenderer) LoadTemplateFromFile(name, filePath string) (*Template, error) {
	doc, err := Open(filePath)
	if err != nil {
		tr.logError("无法打开模板文件 %s: %v", filePath, err)
		return nil, err
	}

	template, err := tr.engine.LoadTemplateFromDocument(name, doc)
	if err != nil {
		tr.logError("无法从文档加载模板 %s: %v", name, err)
		return nil, err
	}

	tr.logInfo("成功加载模板: %s (来源: %s)", name, filePath)
	return template, nil
}

// RenderTemplate 渲染模板到新文档
func (tr *TemplateRenderer) RenderTemplate(templateName string, data *TemplateData) (*Document, error) {
	tr.logInfo("开始渲染模板: %s", templateName)

	// 首先检查数据完整性
	if err := tr.validateTemplateData(data); err != nil {
		tr.logError("模板数据验证失败: %v", err)
		return nil, err
	}

	// 渲染模板
	doc, err := tr.engine.RenderTemplateToDocument(templateName, data)
	if err != nil {
		tr.logError("模板渲染失败: %v", err)
		return nil, err
	}

	tr.logInfo("模板渲染完成: %s", templateName)
	return doc, nil
}

// validateTemplateData 验证模板数据
func (tr *TemplateRenderer) validateTemplateData(data *TemplateData) error {
	if data == nil {
		return fmt.Errorf("模板数据不能为空")
	}

	// 记录数据统计
	tr.logInfo("模板数据统计: 变量=%d, 列表=%d, 条件=%d",
		len(data.Variables), len(data.Lists), len(data.Conditions))

	// 检查列表数据格式
	for listName, listData := range data.Lists {
		if len(listData) == 0 {
			tr.logInfo("警告: 列表 '%s' 为空", listName)
			continue
		}

		// 检查列表项的格式一致性
		for i, item := range listData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if i == 0 {
					tr.logInfo("列表 '%s' 包含 %d 项，第一项字段: %v",
						listName, len(listData), tr.getMapKeys(itemMap))
				}
			} else {
				tr.logInfo("列表 '%s' 第 %d 项不是 map 类型: %T", listName, i, item)
			}
		}
	}

	return nil
}

// getMapKeys 获取map的键列表
func (tr *TemplateRenderer) getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// AnalyzeTemplate 分析模板结构
func (tr *TemplateRenderer) AnalyzeTemplate(templateName string) (*TemplateAnalysis, error) {
	template, err := tr.engine.GetTemplate(templateName)
	if err != nil {
		return nil, err
	}

	analysis := &TemplateAnalysis{
		TemplateName: templateName,
		Variables:    make(map[string]bool),
		Lists:        make(map[string]bool),
		Conditions:   make(map[string]bool),
		Tables:       make([]*TableAnalysis, 0),
	}

	// 分析基础文档
	if template.BaseDoc != nil {
		tr.analyzeDocument(template.BaseDoc, analysis)
	}

	tr.logInfo("模板分析完成: %s", templateName)
	tr.logInfo("- 变量: %d", len(analysis.Variables))
	tr.logInfo("- 列表: %d", len(analysis.Lists))
	tr.logInfo("- 条件: %d", len(analysis.Conditions))
	tr.logInfo("- 表格: %d", len(analysis.Tables))

	return analysis, nil
}

// analyzeDocument 分析文档结构
func (tr *TemplateRenderer) analyzeDocument(doc *Document, analysis *TemplateAnalysis) {
	for i, element := range doc.Body.Elements {
		switch elem := element.(type) {
		case *Paragraph:
			tr.analyzeParagraph(elem, analysis)
		case *Table:
			tableAnalysis := tr.analyzeTable(elem, i)
			analysis.Tables = append(analysis.Tables, tableAnalysis)
		}
	}

	// 分析页眉页脚中的模板变量
	tr.analyzeHeadersFooters(doc, analysis)
}

// analyzeHeadersFooters 分析页眉页脚中的模板变量
func (tr *TemplateRenderer) analyzeHeadersFooters(doc *Document, analysis *TemplateAnalysis) {
	if doc == nil || doc.parts == nil {
		return
	}

	// 遍历所有部件，查找页眉页脚文件
	for partName, partData := range doc.parts {
		if strings.HasPrefix(partName, "word/header") || strings.HasPrefix(partName, "word/footer") {
			// 解析页眉/页脚XML并提取模板变量
			tr.analyzeHeaderFooterXML(partData, analysis)
		}
	}
}

// analyzeHeaderFooterXML 分析页眉/页脚XML中的模板变量
func (tr *TemplateRenderer) analyzeHeaderFooterXML(xmlData []byte, analysis *TemplateAnalysis) {
	// 使用正则表达式提取<w:t>标签中的文本内容
	textPattern := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
	matches := textPattern.FindAllSubmatch(xmlData, -1)

	var fullText strings.Builder
	for _, match := range matches {
		if len(match) >= 2 {
			fullText.Write(match[1])
		}
	}

	// 提取模板变量
	tr.extractTemplateVariables(fullText.String(), analysis)
}

// analyzeParagraph 分析段落
func (tr *TemplateRenderer) analyzeParagraph(para *Paragraph, analysis *TemplateAnalysis) {
	fullText := ""
	for _, run := range para.Runs {
		fullText += run.Text.Content
	}

	tr.extractTemplateVariables(fullText, analysis)
}

// analyzeTable 分析表格
func (tr *TemplateRenderer) analyzeTable(table *Table, index int) *TableAnalysis {
	tableAnalysis := &TableAnalysis{
		Index:         index,
		RowCount:      len(table.Rows),
		HasTemplate:   false,
		TemplateVars:  make(map[string]bool),
		LoopVariables: make([]string, 0),
	}

	if len(table.Rows) > 0 {
		tableAnalysis.ColCount = len(table.Rows[0].Cells)
	}

	// 检查是否是模板表格
	for rowIndex, row := range table.Rows {
		rowHasLoop := false
		for _, cell := range row.Cells {
			for _, para := range cell.Paragraphs {
				fullText := ""
				for _, run := range para.Runs {
					fullText += run.Text.Content
				}

				// 检查循环语法
				eachPattern := regexp.MustCompile(`\{\{#each\s+(\w+)\}\}`)
				if matches := eachPattern.FindStringSubmatch(fullText); len(matches) > 1 {
					tableAnalysis.HasTemplate = true
					tableAnalysis.TemplateRowIndex = rowIndex
					tableAnalysis.LoopVariables = append(tableAnalysis.LoopVariables, matches[1])
					rowHasLoop = true
				}

				// 提取变量
				varPattern := regexp.MustCompile(`\{\{(\w+)\}\}`)
				varMatches := varPattern.FindAllStringSubmatch(fullText, -1)
				for _, match := range varMatches {
					if len(match) >= 2 {
						tableAnalysis.TemplateVars[match[1]] = true
					}
				}
			}
		}

		if rowHasLoop {
			break
		}
	}

	return tableAnalysis
}

// extractTemplateVariables 提取模板变量
func (tr *TemplateRenderer) extractTemplateVariables(text string, analysis *TemplateAnalysis) {
	// 变量: {{变量名}}
	varPattern := regexp.MustCompile(`\{\{(\w+)\}\}`)
	varMatches := varPattern.FindAllStringSubmatch(text, -1)
	for _, match := range varMatches {
		if len(match) >= 2 {
			analysis.Variables[match[1]] = true
		}
	}

	// 条件: {{#if 条件}}
	ifPattern := regexp.MustCompile(`\{\{#if\s+(\w+)\}\}`)
	ifMatches := ifPattern.FindAllStringSubmatch(text, -1)
	for _, match := range ifMatches {
		if len(match) >= 2 {
			analysis.Conditions[match[1]] = true
		}
	}

	// 循环: {{#each 列表}}
	eachPattern := regexp.MustCompile(`\{\{#each\s+(\w+)\}\}`)
	eachMatches := eachPattern.FindAllStringSubmatch(text, -1)
	for _, match := range eachMatches {
		if len(match) >= 2 {
			analysis.Lists[match[1]] = true
		}
	}
}

// TemplateAnalysis 模板分析结果
type TemplateAnalysis struct {
	TemplateName string           // 模板名称
	Variables    map[string]bool  // 变量列表
	Lists        map[string]bool  // 列表变量
	Conditions   map[string]bool  // 条件变量
	Tables       []*TableAnalysis // 表格分析
}

// TableAnalysis 表格分析结果
type TableAnalysis struct {
	Index            int             // 表格索引
	RowCount         int             // 行数
	ColCount         int             // 列数
	HasTemplate      bool            // 是否包含模板语法
	TemplateRowIndex int             // 模板行索引
	TemplateVars     map[string]bool // 模板变量
	LoopVariables    []string        // 循环变量
}

// GetRequiredData 获取模板所需的数据结构
func (analysis *TemplateAnalysis) GetRequiredData() *TemplateData {
	data := NewTemplateData()

	// 设置变量示例值
	for varName := range analysis.Variables {
		data.SetVariable(varName, fmt.Sprintf("示例_%s", varName))
	}

	// 设置条件示例值
	for condName := range analysis.Conditions {
		data.SetCondition(condName, true)
	}

	// 设置列表示例值
	for listName := range analysis.Lists {
		sampleList := []interface{}{
			map[string]interface{}{
				"示例字段1": "示例值1",
				"示例字段2": "示例值2",
			},
			map[string]interface{}{
				"示例字段1": "示例值3",
				"示例字段2": "示例值4",
			},
		}
		data.SetList(listName, sampleList)
	}

	return data
}

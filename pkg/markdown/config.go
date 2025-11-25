// Package markdown 提供Markdown到Word文档的转换功能
package markdown

import "github.com/ZeroHawkeye/wordZero/pkg/document"

// ConvertOptions 转换选项配置
type ConvertOptions struct {
	// 基础配置
	EnableGFM       bool // 启用GitHub Flavored Markdown
	EnableFootnotes bool // 启用脚注支持
	EnableTables    bool // 启用表格支持
	EnableTaskList  bool // 启用任务列表
	EnableMath      bool // 启用数学公式支持（LaTeX语法）

	// 样式配置
	StyleMapping      map[string]string // 自定义样式映射
	DefaultFontFamily string            // 默认字体
	DefaultFontSize   float64           // 默认字号

	// 图片处理
	ImageBasePath string  // 图片基础路径
	EmbedImages   bool    // 是否嵌入图片
	MaxImageWidth float64 // 最大图片宽度（英寸）

	// 链接处理
	PreserveLinkStyle  bool // 保留链接样式
	ConvertToBookmarks bool // 内部链接转书签

	// 文档设置
	GenerateTOC  bool                   // 生成目录
	TOCMaxLevel  int                    // 目录最大级别
	PageSettings *document.PageSettings // 页面设置（使用现有结构）

	// 错误处理
	StrictMode    bool        // 严格模式
	IgnoreErrors  bool        // 忽略转换错误
	ErrorCallback func(error) // 错误回调

	// 进度报告
	ProgressCallback func(int, int) // 进度回调
}

// DefaultOptions 返回默认的转换配置
func DefaultOptions() *ConvertOptions {
	return &ConvertOptions{
		EnableGFM:         true,
		EnableFootnotes:   true,
		EnableTables:      true,
		EnableTaskList:    true,
		EnableMath:        true, // 默认启用数学公式支持
		DefaultFontFamily: "Calibri",
		DefaultFontSize:   11.0,
		EmbedImages:       false,
		MaxImageWidth:     6.0, // 英寸
		GenerateTOC:       true,
		TOCMaxLevel:       3,
		StrictMode:        false,
		IgnoreErrors:      true,
	}
}

// HighQualityOptions 返回高质量转换配置
func HighQualityOptions() *ConvertOptions {
	opts := DefaultOptions()
	opts.EmbedImages = true
	opts.PreserveLinkStyle = true
	opts.ConvertToBookmarks = true
	opts.StrictMode = true
	opts.IgnoreErrors = false
	return opts
}

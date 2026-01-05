package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zerx-lab/wordZero/pkg/document"
)

// WordToMarkdownExporter Word到Markdown导出器接口
type WordToMarkdownExporter interface {
	// ExportToFile 导出Word文档到Markdown文件
	ExportToFile(docxPath, mdPath string, options *ExportOptions) error

	// ExportToString 导出Word文档到Markdown字符串
	ExportToString(doc *document.Document, options *ExportOptions) (string, error)

	// ExportToBytes 导出Word文档到Markdown字节数组
	ExportToBytes(doc *document.Document, options *ExportOptions) ([]byte, error)

	// BatchExport 批量导出
	BatchExport(inputs []string, outputDir string, options *ExportOptions) error
}

// Exporter Word到Markdown导出器实现
type Exporter struct {
	opts *ExportOptions
}

// NewExporter 创建新的导出器实例
func NewExporter(opts *ExportOptions) *Exporter {
	if opts == nil {
		opts = DefaultExportOptions()
	}
	return &Exporter{opts: opts}
}

// ExportToFile 导出Word文档到Markdown文件
func (e *Exporter) ExportToFile(docxPath, mdPath string, options *ExportOptions) error {
	// 加载Word文档
	doc, err := document.Open(docxPath)
	if err != nil {
		return NewExportError("DocumentOpen", fmt.Sprintf("failed to open document: %v", err), err)
	}

	// 设置图片输出路径
	if options == nil {
		options = e.opts
	}
	if options.ExtractImages && options.ImageOutputDir == "" {
		options.ImageOutputDir = filepath.Dir(mdPath)
	}

	// 转换为Markdown
	markdown, err := e.ExportToString(doc, options)
	if err != nil {
		return err
	}

	// 写入文件
	err = os.WriteFile(mdPath, []byte(markdown), 0644)
	if err != nil {
		return NewExportError("FileWrite", fmt.Sprintf("failed to write markdown file: %v", err), err)
	}

	return nil
}

// ExportToString 导出Word文档到Markdown字符串
func (e *Exporter) ExportToString(doc *document.Document, options *ExportOptions) (string, error) {
	bytes, err := e.ExportToBytes(doc, options)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ExportToBytes 导出Word文档到Markdown字节数组
func (e *Exporter) ExportToBytes(doc *document.Document, options *ExportOptions) ([]byte, error) {
	if options != nil {
		e.opts = options
	}

	writer := &MarkdownWriter{
		opts:      e.opts,
		doc:       doc,
		imageNum:  0,
		footnotes: make([]string, 0),
	}

	return writer.Write()
}

// BatchExport 批量导出
func (e *Exporter) BatchExport(inputs []string, outputDir string, options *ExportOptions) error {
	// 确保输出目录存在
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return NewExportError("DirectoryCreate", fmt.Sprintf("failed to create output directory: %v", err), err)
	}

	total := len(inputs)
	for i, input := range inputs {
		// 报告进度
		if options != nil && options.ProgressCallback != nil {
			options.ProgressCallback(i+1, total)
		}

		// 生成输出文件名
		base := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
		output := filepath.Join(outputDir, base+".md")

		// 导出单个文件
		err := e.ExportToFile(input, output, options)
		if err != nil {
			if options != nil && options.ErrorCallback != nil {
				options.ErrorCallback(err)
			}
			if options == nil || !options.IgnoreErrors {
				return err
			}
		}
	}

	return nil
}

// DefaultExportOptions 返回默认的导出配置
func DefaultExportOptions() *ExportOptions {
	return &ExportOptions{
		UseGFMTables:        true,
		PreserveFootnotes:   true,
		PreserveLineBreaks:  false,
		WrapLongLines:       false,
		MaxLineLength:       80,
		ExtractImages:       true,
		ImageNamePattern:    "image_%d.png",
		ImageRelativePath:   true,
		PreserveBookmarks:   true,
		ConvertHyperlinks:   true,
		PreserveCodeStyle:   true,
		DefaultCodeLang:     "",
		IgnoreUnknownStyles: true,
		PreserveTOC:         false,
		IncludeMetadata:     false,
		StripComments:       true,
		UseSetext:           false,
		BulletListMarker:    "-",
		EmphasisMarker:      "*",
		StrictMode:          false,
		IgnoreErrors:        true,
	}
}

// HighQualityExportOptions 返回高质量导出配置
func HighQualityExportOptions() *ExportOptions {
	opts := DefaultExportOptions()
	opts.ExtractImages = true
	opts.PreserveFootnotes = true
	opts.PreserveBookmarks = true
	opts.PreserveTOC = true
	opts.IncludeMetadata = true
	opts.StrictMode = true
	opts.IgnoreErrors = false
	return opts
}

// BidirectionalConverter 双向转换器
type BidirectionalConverter struct {
	mdToWord *Converter
	wordToMd *Exporter
}

// NewBidirectionalConverter 创建双向转换器
func NewBidirectionalConverter(mdOpts *ConvertOptions, exportOpts *ExportOptions) *BidirectionalConverter {
	return &BidirectionalConverter{
		mdToWord: NewConverter(mdOpts),
		wordToMd: NewExporter(exportOpts),
	}
}

// AutoConvert 自动检测文件类型并转换
func (bc *BidirectionalConverter) AutoConvert(inputPath, outputPath string) error {
	ext := strings.ToLower(filepath.Ext(inputPath))

	switch ext {
	case ".md", ".markdown":
		return bc.mdToWord.ConvertFile(inputPath, outputPath, nil)
	case ".docx":
		return bc.wordToMd.ExportToFile(inputPath, outputPath, nil)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
}

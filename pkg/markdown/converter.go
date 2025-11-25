package markdown

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	mathjax "github.com/litao91/goldmark-mathjax"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
)

// MarkdownConverter Markdown转换器接口
type MarkdownConverter interface {
	// ConvertFile 转换单个文件
	ConvertFile(mdPath, docxPath string, options *ConvertOptions) error

	// ConvertBytes 转换字节数据
	ConvertBytes(mdContent []byte, options *ConvertOptions) (*document.Document, error)

	// ConvertString 转换字符串
	ConvertString(mdContent string, options *ConvertOptions) (*document.Document, error)

	// BatchConvert 批量转换
	BatchConvert(inputs []string, outputDir string, options *ConvertOptions) error
}

// Converter 默认转换器实现
type Converter struct {
	md   goldmark.Markdown
	opts *ConvertOptions
}

// NewConverter 创建新的转换器实例
func NewConverter(opts *ConvertOptions) *Converter {
	if opts == nil {
		opts = DefaultOptions()
	}

	extensions := []goldmark.Extender{}
	if opts.EnableGFM {
		extensions = append(extensions, extension.GFM)
	}
	if opts.EnableFootnotes {
		extensions = append(extensions, extension.Footnote)
	}
	if opts.EnableMath {
		// 使用标准的LaTeX数学公式分隔符: $...$ 用于行内公式, $$...$$ 用于块级公式
		extensions = append(extensions, mathjax.NewMathJax(
			mathjax.WithInlineDelim("$", "$"),
			mathjax.WithBlockDelim("$$", "$$"),
		))
	}

	md := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	return &Converter{md: md, opts: opts}
}

// ConvertString 转换字符串内容为Word文档
func (c *Converter) ConvertString(content string, opts *ConvertOptions) (*document.Document, error) {
	return c.ConvertBytes([]byte(content), opts)
}

// ConvertBytes 转换字节数据为Word文档
func (c *Converter) ConvertBytes(content []byte, opts *ConvertOptions) (*document.Document, error) {
	if opts != nil {
		c.opts = opts
	}

	// 创建新的Word文档
	doc := document.New()

	// 应用页面设置
	if c.opts.PageSettings != nil {
		// 这里可以后续扩展，使用现有的页面设置API
	}

	// 解析Markdown
	reader := text.NewReader(content)
	astDoc := c.md.Parser().Parse(reader)

	// 创建渲染器并转换
	renderer := &WordRenderer{
		doc:    doc,
		opts:   c.opts,
		source: content,
	}

	err := renderer.Render(astDoc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// ConvertFile 转换文件
func (c *Converter) ConvertFile(mdPath, docxPath string, options *ConvertOptions) error {
	// 读取Markdown文件
	content, err := os.ReadFile(mdPath)
	if err != nil {
		return NewConversionError("FileRead", "failed to read markdown file", 0, 0, err)
	}

	// 设置图片基础路径（如果未指定）
	if options == nil {
		options = c.opts
	}
	if options.ImageBasePath == "" {
		options.ImageBasePath = filepath.Dir(mdPath)
	}

	// 转换内容
	doc, err := c.ConvertBytes(content, options)
	if err != nil {
		return err
	}

	// 保存Word文档
	err = doc.Save(docxPath)
	if err != nil {
		return NewConversionError("FileSave", "failed to save word document", 0, 0, err)
	}

	return nil
}

// BatchConvert 批量转换文件
func (c *Converter) BatchConvert(inputs []string, outputDir string, options *ConvertOptions) error {
	// 确保输出目录存在
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return NewConversionError("DirectoryCreate", "failed to create output directory", 0, 0, err)
	}

	total := len(inputs)
	for i, input := range inputs {
		// 报告进度
		if options != nil && options.ProgressCallback != nil {
			options.ProgressCallback(i+1, total)
		}

		// 生成输出文件名
		base := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
		output := filepath.Join(outputDir, base+".docx")

		// 转换单个文件
		err := c.ConvertFile(input, output, options)
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

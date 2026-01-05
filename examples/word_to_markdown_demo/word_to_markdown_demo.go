package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/zerx-lab/wordZero/pkg/markdown"
)

func main() {
	fmt.Println("🔄 WordZero - Word转Markdown功能演示")
	fmt.Println("=====================================")

	// 1. 准备输入和输出路径
	inputPath := "examples/output/comprehensive_markdown_demo.docx"
	outputDir := "examples/output"
	outputPath := filepath.Join(outputDir, "converted_from_word.md")
	imagesDir := filepath.Join(outputDir, "images")

	// 检查输入文件是否存在
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		fmt.Printf("⚠️  输入文件不存在: %s\n", inputPath)
		fmt.Println("💡 请先运行 table_and_tasklist_demo.go 生成示例Word文档")
		return
	}

	// 确保输出目录存在
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("❌ 创建输出目录失败: %v", err)
	}

	// 2. 演示基础Word转Markdown功能
	fmt.Println("\n📝 基础转换演示...")
	demonstrateBasicConversion(inputPath, outputPath)

	// 3. 演示高质量转换配置
	fmt.Println("\n✨ 高质量转换演示...")
	demonstrateHighQualityConversion(inputPath, imagesDir)

	// 4. 演示自定义配置转换
	fmt.Println("\n🔧 自定义配置演示...")
	demonstrateCustomConversion(inputPath, outputDir)

	// 5. 演示双向转换器
	fmt.Println("\n🔄 双向转换演示...")
	demonstrateBidirectionalConversion(inputPath, outputDir)

	// 6. 演示批量转换
	fmt.Println("\n📁 批量转换演示...")
	demonstrateBatchConversion(outputDir)

	fmt.Println("\n🎉 所有Word转Markdown演示完成！")
	fmt.Println("📂 查看输出文件: " + outputDir)
}

// demonstrateBasicConversion 演示基础转换功能
func demonstrateBasicConversion(inputPath, outputPath string) {
	fmt.Printf("   输入文件: %s\n", inputPath)
	fmt.Printf("   输出文件: %s\n", outputPath)

	// 使用默认配置创建导出器
	exporter := markdown.NewExporter(markdown.DefaultExportOptions())

	// 执行转换
	err := exporter.ExportToFile(inputPath, outputPath, nil)
	if err != nil {
		fmt.Printf("   ❌ 转换失败: %v\n", err)
		return
	}

	fmt.Println("   ✅ 基础转换完成")

	// 显示转换结果摘要
	showFileInfo(outputPath)
}

// demonstrateHighQualityConversion 演示高质量转换
func demonstrateHighQualityConversion(inputPath, imagesDir string) {
	outputPath := filepath.Join(filepath.Dir(imagesDir), "high_quality_conversion.md")

	// 使用高质量配置
	options := markdown.HighQualityExportOptions()
	options.ImageOutputDir = imagesDir
	options.ExtractImages = true
	options.PreserveFootnotes = true
	options.PreserveTOC = true
	options.IncludeMetadata = true

	// 添加进度回调
	options.ProgressCallback = func(current, total int) {
		fmt.Printf("   📊 转换进度: %d/%d (%.1f%%)\n", current, total, float64(current)/float64(total)*100)
	}

	fmt.Printf("   输入文件: %s\n", inputPath)
	fmt.Printf("   输出文件: %s\n", outputPath)
	fmt.Printf("   图片目录: %s\n", imagesDir)

	exporter := markdown.NewExporter(options)
	err := exporter.ExportToFile(inputPath, outputPath, nil)
	if err != nil {
		fmt.Printf("   ❌ 高质量转换失败: %v\n", err)
		return
	}

	fmt.Println("   ✅ 高质量转换完成")
	showFileInfo(outputPath)
}

// demonstrateCustomConversion 演示自定义配置转换
func demonstrateCustomConversion(inputPath, outputDir string) {
	outputPath := filepath.Join(outputDir, "custom_conversion.md")

	// 创建自定义配置
	options := &markdown.ExportOptions{
		// 表格和格式
		UseGFMTables:       true,
		PreserveFootnotes:  true,
		PreserveLineBreaks: false,
		WrapLongLines:      true,
		MaxLineLength:      80,

		// 图片处理
		ExtractImages:     false, // 不导出图片文件
		ImageRelativePath: true,

		// 链接处理
		PreserveBookmarks: true,
		ConvertHyperlinks: true,

		// 代码处理
		PreserveCodeStyle: true,
		DefaultCodeLang:   "text",

		// 内容处理
		PreserveTOC:     false,
		IncludeMetadata: true,
		StripComments:   true,

		// 格式化选项
		UseSetext:        false, // 使用ATX样式标题
		BulletListMarker: "*",   // 使用*作为项目符号
		EmphasisMarker:   "_",   // 使用_表示斜体

		// 错误处理
		StrictMode:   true,
		IgnoreErrors: false,
		ErrorCallback: func(err error) {
			fmt.Printf("   ⚠️  转换警告: %v\n", err)
		},
	}

	fmt.Printf("   输入文件: %s\n", inputPath)
	fmt.Printf("   输出文件: %s\n", outputPath)
	fmt.Println("   配置特点:")
	fmt.Printf("     • GFM表格: %v\n", options.UseGFMTables)
	fmt.Printf("     • 最大行长: %d字符\n", options.MaxLineLength)
	fmt.Printf("     • 项目符号: %s\n", options.BulletListMarker)
	fmt.Printf("     • 强调符号: %s\n", options.EmphasisMarker)

	exporter := markdown.NewExporter(options)
	err := exporter.ExportToFile(inputPath, outputPath, nil)
	if err != nil {
		fmt.Printf("   ❌ 自定义转换失败: %v\n", err)
		return
	}

	fmt.Println("   ✅ 自定义转换完成")
	showFileInfo(outputPath)
}

// demonstrateBidirectionalConversion 演示双向转换器
func demonstrateBidirectionalConversion(inputPath, outputDir string) {
	// 创建双向转换器
	converter := markdown.NewBidirectionalConverter(
		markdown.HighQualityOptions(),       // Markdown→Word选项
		markdown.HighQualityExportOptions(), // Word→Markdown选项
	)

	// 测试Word→Markdown
	mdPath := filepath.Join(outputDir, "bidirectional_word_to_md.md")
	fmt.Printf("   Word→Markdown: %s → %s\n", inputPath, mdPath)

	err := converter.AutoConvert(inputPath, mdPath)
	if err != nil {
		fmt.Printf("   ❌ Word→Markdown失败: %v\n", err)
		return
	}
	fmt.Println("   ✅ Word→Markdown完成")

	// 测试Markdown→Word (往回转换)
	docxPath := filepath.Join(outputDir, "bidirectional_md_to_word.docx")
	fmt.Printf("   Markdown→Word: %s → %s\n", mdPath, docxPath)

	err = converter.AutoConvert(mdPath, docxPath)
	if err != nil {
		fmt.Printf("   ❌ Markdown→Word失败: %v\n", err)
		return
	}
	fmt.Println("   ✅ Markdown→Word完成")

	showFileInfo(mdPath)
	showFileInfo(docxPath)
}

// demonstrateBatchConversion 演示批量转换
func demonstrateBatchConversion(outputDir string) {
	// 准备批量转换的输入文件
	inputFiles := []string{
		"examples/output/comprehensive_markdown_demo.docx",
	}

	// 检查是否有其他可用的docx文件
	files, err := filepath.Glob(filepath.Join(outputDir, "*.docx"))
	if err == nil {
		for _, file := range files {
			if !contains(inputFiles, file) {
				inputFiles = append(inputFiles, file)
			}
		}
	}

	if len(inputFiles) == 0 {
		fmt.Println("   ⚠️  没有找到可用于批量转换的Word文档")
		return
	}

	batchOutputDir := filepath.Join(outputDir, "batch_converted")
	fmt.Printf("   输入文件数量: %d\n", len(inputFiles))
	fmt.Printf("   输出目录: %s\n", batchOutputDir)

	// 配置批量转换选项
	options := markdown.DefaultExportOptions()
	options.ProgressCallback = func(current, total int) {
		fmt.Printf("   📊 批量转换进度: %d/%d\n", current, total)
	}
	options.ErrorCallback = func(err error) {
		fmt.Printf("   ⚠️  转换错误: %v\n", err)
	}

	// 执行批量转换
	exporter := markdown.NewExporter(options)
	err = exporter.BatchExport(inputFiles, batchOutputDir, options)
	if err != nil {
		fmt.Printf("   ❌ 批量转换失败: %v\n", err)
		return
	}

	fmt.Println("   ✅ 批量转换完成")

	// 显示转换结果
	convertedFiles, _ := filepath.Glob(filepath.Join(batchOutputDir, "*.md"))
	fmt.Printf("   📄 成功转换 %d 个文件\n", len(convertedFiles))
	for _, file := range convertedFiles {
		fmt.Printf("     • %s\n", filepath.Base(file))
	}
}

// showFileInfo 显示文件信息
func showFileInfo(filePath string) {
	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("   📄 文件信息获取失败: %v\n", err)
		return
	}

	// 读取文件内容获取行数
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("   📄 %s (大小: %d 字节)\n", filepath.Base(filePath), info.Size())
		return
	}

	lines := len(strings.Split(string(content), "\n"))
	fmt.Printf("   📄 %s (大小: %d 字节, %d 行)\n", filepath.Base(filePath), info.Size(), lines)

	// 显示前几行内容预览
	preview := strings.Split(string(content), "\n")
	maxPreview := 3
	if len(preview) > maxPreview {
		fmt.Println("   📋 内容预览:")
		for i := 0; i < maxPreview && i < len(preview); i++ {
			line := preview[i]
			if len(line) > 60 {
				line = line[:57] + "..."
			}
			fmt.Printf("      %s\n", line)
		}
		if len(preview) > maxPreview {
			fmt.Printf("      ... (还有 %d 行)\n", len(preview)-maxPreview)
		}
	}
}

// contains 检查切片是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

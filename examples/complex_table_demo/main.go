// Package main 演示WordZero复杂表格结构功能
// 展示如何在表格单元格中添加段落、列表、嵌套表格和图片
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
)

func main() {
	fmt.Println("=== WordZero 复杂表格结构演示 ===")

	// 确保输出目录存在
	if _, err := os.Stat("examples/output"); os.IsNotExist(err) {
		os.MkdirAll("examples/output", 0755)
	}

	// 创建新文档
	doc := document.New()

	// 添加文档标题
	title := doc.AddParagraph("WordZero 复杂表格结构演示")
	title.SetAlignment(document.AlignCenter)
	doc.AddParagraph("") // 空行

	// 演示1：单元格中添加多个段落
	fmt.Println("1. 单元格多段落演示...")
	demonstrateMultipleParagraphs(doc)

	// 演示2：单元格中添加列表
	fmt.Println("2. 单元格列表演示...")
	demonstrateCellLists(doc)

	// 演示3：单元格中添加嵌套表格
	fmt.Println("3. 嵌套表格演示...")
	demonstrateNestedTable(doc)

	// 演示4：单元格中添加图片
	fmt.Println("4. 单元格图片演示...")
	demonstrateCellImages(doc)

	// 演示5：综合复杂表格
	fmt.Println("5. 综合复杂表格演示...")
	demonstrateComplexTable(doc)

	// 保存文档
	outputFile := "examples/output/complex_table_demo.docx"
	err := doc.Save(outputFile)
	if err != nil {
		log.Fatalf("保存文档失败: %v", err)
	}

	fmt.Printf("\n复杂表格演示文档已保存到: %s\n", outputFile)
	fmt.Println("=== 演示完成 ===")
}

// demonstrateMultipleParagraphs 演示单元格中添加多个段落
func demonstrateMultipleParagraphs(doc *document.Document) {
	doc.AddParagraph("1. 单元格多段落演示")
	doc.AddParagraph("")

	// 创建表格
	config := &document.TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 8000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		log.Printf("创建表格失败: %v", err)
		return
	}

	// 设置表头
	table.SetCellText(0, 0, "多段落单元格")
	table.SetCellText(0, 1, "格式化段落单元格")

	// 在单元格中添加多个段落
	table.SetCellText(1, 0, "这是第一段内容，介绍了表格的基本概念。")
	table.AddCellParagraph(1, 0, "这是第二段内容，说明了表格的使用方法。")
	table.AddCellParagraph(1, 0, "这是第三段内容，总结了表格的优势。")

	// 在另一个单元格中添加格式化段落
	table.SetCellText(1, 1, "普通文本介绍")
	table.AddCellFormattedParagraph(1, 1, "粗体重点内容", &document.TextFormat{
		Bold:     true,
		FontSize: 12,
	})
	table.AddCellFormattedParagraph(1, 1, "红色提示文字", &document.TextFormat{
		FontColor: "FF0000",
		Italic:    true,
	})
	table.AddCellFormattedParagraph(1, 1, "大号蓝色标题", &document.TextFormat{
		FontColor: "0000FF",
		FontSize:  14,
		Bold:      true,
	})

	doc.AddParagraph("")
	fmt.Println("   多段落单元格演示完成")
}

// demonstrateCellLists 演示单元格中添加列表
func demonstrateCellLists(doc *document.Document) {
	doc.AddParagraph("2. 单元格列表演示")
	doc.AddParagraph("")

	// 创建表格
	config := &document.TableConfig{
		Rows:  2,
		Cols:  3,
		Width: 9000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		log.Printf("创建表格失败: %v", err)
		return
	}

	// 设置表头
	table.SetCellText(0, 0, "无序列表")
	table.SetCellText(0, 1, "有序列表")
	table.SetCellText(0, 2, "罗马数字列表")

	// 添加无序列表
	bulletList := &document.CellListConfig{
		Type:         document.ListTypeBullet,
		BulletSymbol: document.BulletTypeDot,
		Items: []string{
			"第一个要点",
			"第二个要点",
			"第三个要点",
		},
	}
	table.ClearCellParagraphs(1, 0) // 清空默认段落
	table.AddCellList(1, 0, bulletList)

	// 添加有序列表
	numberList := &document.CellListConfig{
		Type: document.ListTypeNumber,
		Items: []string{
			"第一步操作",
			"第二步操作",
			"第三步操作",
		},
	}
	table.ClearCellParagraphs(1, 1)
	table.AddCellList(1, 1, numberList)

	// 添加罗马数字列表
	romanList := &document.CellListConfig{
		Type: document.ListTypeUpperRoman,
		Items: []string{
			"主要内容",
			"次要内容",
			"补充内容",
		},
	}
	table.ClearCellParagraphs(1, 2)
	table.AddCellList(1, 2, romanList)

	doc.AddParagraph("")
	fmt.Println("   单元格列表演示完成")
}

// demonstrateNestedTable 演示嵌套表格
func demonstrateNestedTable(doc *document.Document) {
	doc.AddParagraph("3. 嵌套表格演示")
	doc.AddParagraph("")

	// 创建主表格
	config := &document.TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 8000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		log.Printf("创建表格失败: %v", err)
		return
	}

	// 设置表头
	table.SetCellText(0, 0, "产品信息")
	table.SetCellText(0, 1, "销售数据")

	// 在第一个单元格添加产品信息嵌套表格
	productNestedConfig := &document.TableConfig{
		Rows:  3,
		Cols:  2,
		Width: 3500,
		Data: [][]string{
			{"属性", "值"},
			{"名称", "智能手表"},
			{"型号", "SW-2024"},
		},
	}
	table.ClearCellParagraphs(1, 0)
	table.AddCellParagraph(1, 0, "产品详细信息：")
	table.AddNestedTable(1, 0, productNestedConfig)

	// 在第二个单元格添加销售数据嵌套表格
	salesNestedConfig := &document.TableConfig{
		Rows:  4,
		Cols:  2,
		Width: 3500,
		Data: [][]string{
			{"季度", "销量"},
			{"Q1", "1000"},
			{"Q2", "1500"},
			{"Q3", "2000"},
		},
	}
	table.ClearCellParagraphs(1, 1)
	table.AddCellParagraph(1, 1, "季度销售数据：")
	table.AddNestedTable(1, 1, salesNestedConfig)

	doc.AddParagraph("")
	fmt.Println("   嵌套表格演示完成")
}

// demonstrateCellImages 演示单元格中添加图片
func demonstrateCellImages(doc *document.Document) {
	doc.AddParagraph("4. 单元格图片演示")
	doc.AddParagraph("")

	// 创建表格
	config := &document.TableConfig{
		Rows:  2,
		Cols:  2,
		Width: 8000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		log.Printf("创建表格失败: %v", err)
		return
	}

	// 设置表头
	table.SetCellText(0, 0, "产品图片")
	table.SetCellText(0, 1, "产品描述")

	// 在单元格中添加图片
	imageData := createColorImage(150, 100, color.RGBA{100, 150, 255, 255})
	table.ClearCellParagraphs(1, 0)
	_, err = doc.AddCellImageFromData(table, 1, 0, imageData, 40) // 40mm宽度
	if err != nil {
		log.Printf("添加图片失败: %v", err)
	}

	// 在另一个单元格添加描述
	table.SetCellText(1, 1, "产品名称：智能设备")
	table.AddCellParagraph(1, 1, "规格：150mm x 100mm")
	table.AddCellFormattedParagraph(1, 1, "状态：在售", &document.TextFormat{
		FontColor: "00AA00",
		Bold:      true,
	})

	doc.AddParagraph("")
	fmt.Println("   单元格图片演示完成")
}

// demonstrateComplexTable 演示综合复杂表格
func demonstrateComplexTable(doc *document.Document) {
	doc.AddParagraph("5. 综合复杂表格演示")
	doc.AddParagraph("")

	// 创建复杂表格
	config := &document.TableConfig{
		Rows:  4,
		Cols:  3,
		Width: 10000,
	}

	table, err := doc.AddTable(config)
	if err != nil {
		log.Printf("创建表格失败: %v", err)
		return
	}

	// 第一行：表头
	table.SetCellText(0, 0, "项目")
	table.SetCellText(0, 1, "详细信息")
	table.SetCellText(0, 2, "备注")

	// 第二行：多段落内容
	table.SetCellText(1, 0, "公司简介")
	table.ClearCellParagraphs(1, 1)
	table.AddCellParagraph(1, 1, "WordZero科技是一家专注于文档处理的技术公司。")
	table.AddCellFormattedParagraph(1, 1, "成立于2024年", &document.TextFormat{Bold: true})
	table.AddCellFormattedParagraph(1, 1, "总部位于北京", &document.TextFormat{Italic: true})

	// 添加备注列表
	noteList := &document.CellListConfig{
		Type:         document.ListTypeBullet,
		BulletSymbol: document.BulletTypeArrow,
		Items:        []string{"技术驱动", "用户至上", "持续创新"},
	}
	table.ClearCellParagraphs(1, 2)
	table.AddCellList(1, 2, noteList)

	// 第三行：嵌套表格
	table.SetCellText(2, 0, "产品矩阵")
	nestedConfig := &document.TableConfig{
		Rows:  3,
		Cols:  2,
		Width: 3000,
		Data: [][]string{
			{"产品", "版本"},
			{"WordZero Core", "v1.0"},
			{"WordZero Pro", "v2.0"},
		},
	}
	table.ClearCellParagraphs(2, 1)
	table.AddNestedTable(2, 1, nestedConfig)
	table.SetCellText(2, 2, "更多产品开发中...")

	// 第四行：图片和描述
	table.SetCellText(3, 0, "团队展示")

	// 添加团队图片
	teamImage := createColorImage(120, 80, color.RGBA{150, 200, 150, 255})
	table.ClearCellParagraphs(3, 1)
	doc.AddCellImageFromData(table, 3, 1, teamImage, 35)
	table.AddCellParagraph(3, 1, "专业团队")

	// 添加联系方式
	contactList := &document.CellListConfig{
		Type: document.ListTypeNumber,
		Items: []string{
			"电话：400-123-4567",
			"邮箱：contact@wordzero.com",
			"网站：www.wordzero.com",
		},
	}
	table.ClearCellParagraphs(3, 2)
	table.AddCellList(3, 2, contactList)

	doc.AddParagraph("")
	fmt.Println("   综合复杂表格演示完成")
}

// createColorImage 创建指定颜色的示例图片
func createColorImage(width, height int, bgColor color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 添加边框
	borderColor := color.RGBA{50, 50, 50, 255}
	for x := 0; x < width; x++ {
		img.Set(x, 0, borderColor)
		img.Set(x, height-1, borderColor)
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, borderColor)
		img.Set(width-1, y, borderColor)
	}

	// 添加中心十字标记
	centerX, centerY := width/2, height/2
	markColor := color.RGBA{0, 0, 0, 255}
	for x := centerX - 10; x <= centerX+10 && x < width; x++ {
		if x >= 0 {
			img.Set(x, centerY, markColor)
		}
	}
	for y := centerY - 10; y <= centerY+10 && y < height; y++ {
		if y >= 0 {
			img.Set(centerX, y, markColor)
		}
	}

	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

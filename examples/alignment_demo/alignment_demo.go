package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/zerx-lab/wordZero/pkg/document"
)

// createSampleImageWithText 创建带文字的示例图片
func createSampleImageWithText(width, height int, bgColor color.RGBA, text string) []byte {
	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 添加边框
	borderColor := color.RGBA{0, 0, 0, 255} // 黑色边框
	for x := 0; x < width; x++ {
		img.Set(x, 0, borderColor)        // 上边框
		img.Set(x, height-1, borderColor) // 下边框
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, borderColor)       // 左边框
		img.Set(width-1, y, borderColor) // 右边框
	}

	// 添加中心标记点（简单的十字）
	centerX := width / 2
	centerY := height / 2
	markColor := color.RGBA{0, 0, 0, 255} // 黑色标记

	// 画水平线
	for x := centerX - 10; x <= centerX+10; x++ {
		if x >= 0 && x < width {
			img.Set(x, centerY, markColor)
		}
	}

	// 画垂直线
	for y := centerY - 10; y <= centerY+10; y++ {
		if y >= 0 && y < height {
			img.Set(centerX, y, markColor)
		}
	}

	// 转换为PNG字节数组
	buf := new(bytes.Buffer)
	png.Encode(buf, img)
	return buf.Bytes()
}

// saveImageToFile 保存图片到文件
func saveImageToFile(imageData []byte, filePath string) error {
	return os.WriteFile(filePath, imageData, 0644)
}

func main() {
	fmt.Println("=== WordZero 图片对齐功能演示 ===")

	// 创建新文档
	doc := document.New()

	// 添加文档标题
	titlePara := doc.AddParagraph("")
	titlePara.AddFormattedText("WordZero 图片对齐功能完整演示", &document.TextFormat{
		Bold:      true,
		FontSize:  18,
		FontColor: "000080",
	})
	titlePara.SetAlignment(document.AlignCenter)
	doc.AddParagraph("")

	// 添加说明文字
	descPara := doc.AddParagraph("")
	descPara.AddFormattedText("本文档演示了WordZero中嵌入式图片的各种对齐方式，包括左对齐、居中对齐、右对齐和两端对齐。", &document.TextFormat{
		FontSize:  12,
		FontColor: "444444",
	})
	descPara.SetAlignment(document.AlignJustify)
	doc.AddParagraph("")

	// 确保输出目录存在
	outputDir := "examples/output"
	imagesDir := outputDir + "/images"
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// ==================== 第一部分：左对齐图片 ====================
	fmt.Println("1. 创建左对齐图片演示...")

	leftTitlePara := doc.AddParagraph("")
	leftTitlePara.AddFormattedText("1. 左对齐图片", &document.TextFormat{
		Bold:     true,
		FontSize: 14,
	})
	leftTitlePara.SetAlignment(document.AlignLeft)

	doc.AddParagraph("左对齐图片会显示在段落的左侧，这是默认的对齐方式。")

	// 创建左对齐图片
	leftImageData := createSampleImageWithText(200, 100, color.RGBA{255, 200, 200, 255}, "左对齐")
	leftImagePath := filepath.Join(imagesDir, "left_align_image.png")
	if err := saveImageToFile(leftImageData, leftImagePath); err != nil {
		log.Fatalf("保存左对齐图片失败: %v", err)
	}

	leftImageInfo, err := doc.AddImageFromData(
		leftImageData,
		"left_align_image.png",
		document.ImageFormatPNG,
		200,
		100,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  50.0,
				Height: 25.0,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignLeft,
			AltText:   "左对齐图片示例",
			Title:     "这是一个左对齐的图片",
		},
	)
	if err != nil {
		log.Fatalf("添加左对齐图片失败: %v", err)
	}
	fmt.Printf("   添加左对齐图片成功，ID: %s\n", leftImageInfo.ID)

	doc.AddParagraph("上面的图片使用了左对齐方式。在大多数从左到右书写的语言中，左对齐是最常见的图片显示方式。")
	doc.AddParagraph("")

	// ==================== 第二部分：居中对齐图片 ====================
	fmt.Println("2. 创建居中对齐图片演示...")

	centerTitlePara := doc.AddParagraph("")
	centerTitlePara.AddFormattedText("2. 居中对齐图片", &document.TextFormat{
		Bold:     true,
		FontSize: 14,
	})
	centerTitlePara.SetAlignment(document.AlignLeft)

	doc.AddParagraph("居中对齐图片会显示在页面的中央，适合用作重要的配图或装饰元素。")

	// 创建居中对齐图片
	centerImageData := createSampleImageWithText(180, 120, color.RGBA{200, 255, 200, 255}, "居中对齐")
	centerImagePath := filepath.Join(imagesDir, "center_align_image.png")
	if err := saveImageToFile(centerImageData, centerImagePath); err != nil {
		log.Fatalf("保存居中对齐图片失败: %v", err)
	}

	centerImageInfo, err := doc.AddImageFromData(
		centerImageData,
		"center_align_image.png",
		document.ImageFormatPNG,
		180,
		120,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  45.0,
				Height: 30.0,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignCenter,
			AltText:   "居中对齐图片示例",
			Title:     "这是一个居中对齐的图片",
		},
	)
	if err != nil {
		log.Fatalf("添加居中对齐图片失败: %v", err)
	}
	fmt.Printf("   添加居中对齐图片成功，ID: %s\n", centerImageInfo.ID)

	doc.AddParagraph("上面的图片使用了居中对齐方式。居中对齐常用于标题图片、logo、或者需要突出显示的重要图片。")
	doc.AddParagraph("")

	// ==================== 第三部分：右对齐图片 ====================
	fmt.Println("3. 创建右对齐图片演示...")

	rightTitlePara := doc.AddParagraph("")
	rightTitlePara.AddFormattedText("3. 右对齐图片", &document.TextFormat{
		Bold:     true,
		FontSize: 14,
	})
	rightTitlePara.SetAlignment(document.AlignLeft)

	doc.AddParagraph("右对齐图片会显示在页面的右侧，适合用作辅助说明或装饰性元素。")

	// 创建右对齐图片
	rightImageData := createSampleImageWithText(160, 130, color.RGBA{200, 200, 255, 255}, "右对齐")
	rightImagePath := filepath.Join(imagesDir, "right_align_image.png")
	if err := saveImageToFile(rightImageData, rightImagePath); err != nil {
		log.Fatalf("保存右对齐图片失败: %v", err)
	}

	rightImageInfo, err := doc.AddImageFromData(
		rightImageData,
		"right_align_image.png",
		document.ImageFormatPNG,
		160,
		130,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  40.0,
				Height: 32.5,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignRight,
			AltText:   "右对齐图片示例",
			Title:     "这是一个右对齐的图片",
		},
	)
	if err != nil {
		log.Fatalf("添加右对齐图片失败: %v", err)
	}
	fmt.Printf("   添加右对齐图片成功，ID: %s\n", rightImageInfo.ID)

	doc.AddParagraph("上面的图片使用了右对齐方式。右对齐在某些设计布局中可以创造视觉平衡，特别是在页面右侧有其他元素时。")
	doc.AddParagraph("")

	// ==================== 第四部分：两端对齐图片 ====================
	fmt.Println("4. 创建两端对齐图片演示...")

	justifyTitlePara := doc.AddParagraph("")
	justifyTitlePara.AddFormattedText("4. 两端对齐图片", &document.TextFormat{
		Bold:     true,
		FontSize: 14,
	})
	justifyTitlePara.SetAlignment(document.AlignLeft)

	doc.AddParagraph("两端对齐通常用于文本，对于图片，其效果类似于左对齐。")

	// 创建两端对齐图片
	justifyImageData := createSampleImageWithText(190, 110, color.RGBA{255, 255, 200, 255}, "两端对齐")
	justifyImagePath := filepath.Join(imagesDir, "justify_align_image.png")
	if err := saveImageToFile(justifyImageData, justifyImagePath); err != nil {
		log.Fatalf("保存两端对齐图片失败: %v", err)
	}

	justifyImageInfo, err := doc.AddImageFromData(
		justifyImageData,
		"justify_align_image.png",
		document.ImageFormatPNG,
		190,
		110,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  47.5,
				Height: 27.5,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignJustify,
			AltText:   "两端对齐图片示例",
			Title:     "这是一个两端对齐的图片",
		},
	)
	if err != nil {
		log.Fatalf("添加两端对齐图片失败: %v", err)
	}
	fmt.Printf("   添加两端对齐图片成功，ID: %s\n", justifyImageInfo.ID)

	doc.AddParagraph("上面的图片使用了两端对齐方式。虽然两端对齐主要用于文本，但在某些特殊的版式设计中也可能用到。")
	doc.AddParagraph("")

	// ==================== 第五部分：动态修改对齐方式 ====================
	fmt.Println("5. 演示动态修改图片对齐方式...")

	dynamicTitlePara := doc.AddParagraph("")
	dynamicTitlePara.AddFormattedText("5. 动态修改对齐方式", &document.TextFormat{
		Bold:     true,
		FontSize: 14,
	})
	dynamicTitlePara.SetAlignment(document.AlignLeft)

	doc.AddParagraph("WordZero支持在添加图片后动态修改其对齐方式。下面的图片最初设置为左对齐，然后被修改为居中对齐。")

	// 创建图片，初始为左对齐
	dynamicImageData := createSampleImageWithText(170, 140, color.RGBA{255, 200, 255, 255}, "动态修改")
	dynamicImagePath := filepath.Join(imagesDir, "dynamic_align_image.png")
	if err := saveImageToFile(dynamicImageData, dynamicImagePath); err != nil {
		log.Fatalf("保存动态修改图片失败: %v", err)
	}

	dynamicImageInfo, err := doc.AddImageFromData(
		dynamicImageData,
		"dynamic_align_image.png",
		document.ImageFormatPNG,
		170,
		140,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  42.5,
				Height: 35.0,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignLeft, // 初始左对齐
			AltText:   "动态修改对齐方式图片示例",
			Title:     "这是一个演示动态修改对齐方式的图片",
		},
	)
	if err != nil {
		log.Fatalf("添加动态修改图片失败: %v", err)
	}

	// 动态修改为居中对齐
	err = doc.SetImageAlignment(dynamicImageInfo, document.AlignCenter)
	if err != nil {
		log.Fatalf("修改图片对齐方式失败: %v", err)
	}

	fmt.Printf("   添加动态修改图片成功，ID: %s\n", dynamicImageInfo.ID)
	fmt.Printf("   图片对齐方式已从左对齐修改为居中对齐\n")

	doc.AddParagraph("上面的图片演示了动态修改对齐方式的功能。通过SetImageAlignment方法，可以在添加图片后随时调整其对齐方式。")
	doc.AddParagraph("")

	// ==================== 第六部分：对齐方式对比 ====================
	fmt.Println("6. 创建对齐方式对比演示...")

	comparisonTitlePara := doc.AddParagraph("")
	comparisonTitlePara.AddFormattedText("6. 对齐方式对比", &document.TextFormat{
		Bold:     true,
		FontSize: 14,
	})
	comparisonTitlePara.SetAlignment(document.AlignLeft)

	doc.AddParagraph("下面连续展示四种对齐方式的图片，便于比较它们的视觉效果：")

	// 创建对比图片 - 左对齐
	compLeftImageData := createSampleImageWithText(120, 80, color.RGBA{255, 180, 180, 255}, "左对齐")
	compLeftImagePath := filepath.Join(imagesDir, "comp_left_image.png")
	if err := saveImageToFile(compLeftImageData, compLeftImagePath); err != nil {
		log.Fatalf("保存对比左对齐图片失败: %v", err)
	}

	_, err = doc.AddImageFromData(
		compLeftImageData,
		"comp_left_image.png",
		document.ImageFormatPNG,
		120,
		80,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  30.0,
				Height: 20.0,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignLeft,
			AltText:   "对比用左对齐图片",
		},
	)
	if err != nil {
		log.Fatalf("添加对比左对齐图片失败: %v", err)
	}

	// 创建对比图片 - 居中对齐
	compCenterImageData := createSampleImageWithText(120, 80, color.RGBA{180, 255, 180, 255}, "居中对齐")
	compCenterImagePath := filepath.Join(imagesDir, "comp_center_image.png")
	if err := saveImageToFile(compCenterImageData, compCenterImagePath); err != nil {
		log.Fatalf("保存对比居中对齐图片失败: %v", err)
	}

	_, err = doc.AddImageFromData(
		compCenterImageData,
		"comp_center_image.png",
		document.ImageFormatPNG,
		120,
		80,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  30.0,
				Height: 20.0,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignCenter,
			AltText:   "对比用居中对齐图片",
		},
	)
	if err != nil {
		log.Fatalf("添加对比居中对齐图片失败: %v", err)
	}

	// 创建对比图片 - 右对齐
	compRightImageData := createSampleImageWithText(120, 80, color.RGBA{180, 180, 255, 255}, "右对齐")
	compRightImagePath := filepath.Join(imagesDir, "comp_right_image.png")
	if err := saveImageToFile(compRightImageData, compRightImagePath); err != nil {
		log.Fatalf("保存对比右对齐图片失败: %v", err)
	}

	_, err = doc.AddImageFromData(
		compRightImageData,
		"comp_right_image.png",
		document.ImageFormatPNG,
		120,
		80,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  30.0,
				Height: 20.0,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignRight,
			AltText:   "对比用右对齐图片",
		},
	)
	if err != nil {
		log.Fatalf("添加对比右对齐图片失败: %v", err)
	}

	// 创建对比图片 - 两端对齐
	compJustifyImageData := createSampleImageWithText(120, 80, color.RGBA{255, 255, 180, 255}, "两端对齐")
	compJustifyImagePath := filepath.Join(imagesDir, "comp_justify_image.png")
	if err := saveImageToFile(compJustifyImageData, compJustifyImagePath); err != nil {
		log.Fatalf("保存对比两端对齐图片失败: %v", err)
	}

	_, err = doc.AddImageFromData(
		compJustifyImageData,
		"comp_justify_image.png",
		document.ImageFormatPNG,
		120,
		80,
		&document.ImageConfig{
			Size: &document.ImageSize{
				Width:  30.0,
				Height: 20.0,
			},
			Position:  document.ImagePositionInline,
			Alignment: document.AlignJustify,
			AltText:   "对比用两端对齐图片",
		},
	)
	if err != nil {
		log.Fatalf("添加对比两端对齐图片失败: %v", err)
	}

	fmt.Printf("   添加对比图片成功，数量: 4张\n")

	doc.AddParagraph("上面的四张图片分别展示了左对齐、居中对齐、右对齐和两端对齐的效果。通过对比可以看出不同对齐方式的视觉差异。")
	doc.AddParagraph("")

	// ==================== 第七部分：功能总结 ====================
	summaryTitlePara := doc.AddParagraph("")
	summaryTitlePara.AddFormattedText("7. 功能总结", &document.TextFormat{
		Bold:     true,
		FontSize: 14,
	})
	summaryTitlePara.SetAlignment(document.AlignLeft)

	doc.AddParagraph("WordZero的图片对齐功能包括以下特性：")
	doc.AddParagraph("• 支持四种对齐方式：左对齐、居中对齐、右对齐、两端对齐")
	doc.AddParagraph("• 仅适用于嵌入式图片（ImagePositionInline）")
	doc.AddParagraph("• 支持在添加图片时直接设置对齐方式")
	doc.AddParagraph("• 支持在添加图片后动态修改对齐方式")
	doc.AddParagraph("• 与浮动图片的位置控制功能互补")
	doc.AddParagraph("• 符合Word文档的标准对齐规范")

	// 保存文档
	outputFile := outputDir + "/image_alignment_demo.docx"
	err = doc.Save(outputFile)
	if err != nil {
		log.Fatalf("保存文档失败: %v", err)
	}

	fmt.Printf("\n=== 图片对齐功能演示完成 ===\n")
	fmt.Printf("文档已保存到: %s\n", outputFile)
	fmt.Printf("生成的图片文件保存在: %s\n", imagesDir)
	fmt.Printf("生成的图片文件包括:\n")

	imageFiles := []string{
		"left_align_image.png (左对齐图片，红色)",
		"center_align_image.png (居中对齐图片，绿色)",
		"right_align_image.png (右对齐图片，蓝色)",
		"justify_align_image.png (两端对齐图片，黄色)",
		"dynamic_align_image.png (动态修改图片，紫色)",
		"comp_left_image.png (对比用左对齐图片，浅红色)",
		"comp_center_image.png (对比用居中对齐图片，浅绿色)",
		"comp_right_image.png (对比用右对齐图片，浅蓝色)",
		"comp_justify_image.png (对比用两端对齐图片，浅黄色)",
	}

	for i, file := range imageFiles {
		fmt.Printf("  %d. %s\n", i+1, file)
	}

	fmt.Printf("\n功能演示统计：\n")
	fmt.Printf("- 演示的图片数量: 9张\n")
	fmt.Printf("- 覆盖的对齐方式: 4种 (左对齐、居中对齐、右对齐、两端对齐)\n")
	fmt.Printf("- 动态修改演示: 1个\n")
	fmt.Printf("- 对比演示图片: 4张\n")
	fmt.Printf("- 功能测试场景: 7个部分\n")
}

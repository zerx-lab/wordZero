// Package main 演示图片持久性修复：打开包含图片的文档，修改后重新保存，图片不会丢失
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/zerx-lab/wordZero/pkg/document"
)

// createSampleImage 创建示例图片
func createSampleImage(width, height int, bgColor color.RGBA, label string) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 添加边框
	borderColor := color.RGBA{0, 0, 0, 255}
	for x := 0; x < width; x++ {
		img.Set(x, 0, borderColor)
		img.Set(x, height-1, borderColor)
	}
	for y := 0; y < height; y++ {
		img.Set(0, y, borderColor)
		img.Set(width-1, y, borderColor)
	}

	// 添加中心十字标记
	centerX := width / 2
	centerY := height / 2
	markColor := color.RGBA{255, 255, 255, 255}

	for x := centerX - 20; x <= centerX+20; x++ {
		if x >= 0 && x < width {
			img.Set(x, centerY, markColor)
		}
	}
	for y := centerY - 20; y <= centerY+20; y++ {
		if y >= 0 && y < height {
			img.Set(centerX, y, markColor)
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func main() {
	fmt.Println("图片持久性修复演示")
	fmt.Println("===================")

	// 确保输出目录存在
	if _, err := os.Stat("examples/output"); os.IsNotExist(err) {
		os.MkdirAll("examples/output", 0755)
	}

	// 步骤1: 创建包含图片的原始文档
	fmt.Println("\n步骤1: 创建包含图片的原始文档...")
	doc1 := document.New()

	// 添加标题
	title := doc1.AddParagraph("原始文档 - 包含图片")
	title.SetAlignment(document.AlignCenter)

	// 添加描述
	doc1.AddParagraph("这是一个包含图片的文档。我们将打开它，修改内容，然后保存为新文件。")

	// 添加第一张图片（红色）
	fmt.Println("   添加第一张图片（红色）...")
	imageData1 := createSampleImage(200, 150, color.RGBA{255, 100, 100, 255}, "图片1")
	_, err := doc1.AddImageFromData(
		imageData1,
		"image1.png",
		document.ImageFormatPNG,
		200, 150,
		&document.ImageConfig{
			Position:  document.ImagePositionInline,
			Alignment: document.AlignCenter,
			AltText:   "第一张图片",
			Title:     "红色图片",
		},
	)
	if err != nil {
		log.Fatalf("添加第一张图片失败: %v", err)
	}

	doc1.AddParagraph("上图：第一张图片（红色）")

	// 添加第二张图片（蓝色）
	fmt.Println("   添加第二张图片（蓝色）...")
	imageData2 := createSampleImage(200, 150, color.RGBA{100, 100, 255, 255}, "图片2")
	_, err = doc1.AddImageFromData(
		imageData2,
		"image2.png",
		document.ImageFormatPNG,
		200, 150,
		&document.ImageConfig{
			Position:  document.ImagePositionInline,
			Alignment: document.AlignCenter,
			AltText:   "第二张图片",
			Title:     "蓝色图片",
		},
	)
	if err != nil {
		log.Fatalf("添加第二张图片失败: %v", err)
	}

	doc1.AddParagraph("上图：第二张图片（蓝色）")

	// 保存原始文档
	originalFile := "examples/output/image_persistence_original.docx"
	err = doc1.Save(originalFile)
	if err != nil {
		log.Fatalf("保存原始文档失败: %v", err)
	}
	fmt.Printf("   ✓ 原始文档已保存: %s\n", originalFile)

	// 步骤2: 打开文档
	fmt.Println("\n步骤2: 打开刚才创建的文档...")
	doc2, err := document.Open(originalFile)
	if err != nil {
		log.Fatalf("打开文档失败: %v", err)
	}
	fmt.Println("   ✓ 文档打开成功")

	// 验证图片是否正确加载
	imageCount := 0
	for partName := range doc2.GetParts() {
		if len(partName) > 11 && partName[:11] == "word/media/" {
			imageCount++
			fmt.Printf("   - 找到图片: %s\n", partName)
		}
	}
	
	if imageCount < 2 {
		log.Fatalf("打开的文档应该包含2张图片，但只找到 %d 张", imageCount)
	}
	fmt.Printf("   ✓ 确认文档包含 %d 张图片\n", imageCount)

	// 步骤3: 修改文档内容
	fmt.Println("\n步骤3: 修改文档内容...")
	
	// 在文档开头添加新段落
	doc2.AddParagraph("【修改后的内容】这个段落是在打开文档后新添加的。")
	
	// 添加第三张图片（绿色）
	fmt.Println("   添加第三张图片（绿色）...")
	imageData3 := createSampleImage(200, 150, color.RGBA{100, 255, 100, 255}, "图片3")
	_, err = doc2.AddImageFromData(
		imageData3,
		"image3.png",
		document.ImageFormatPNG,
		200, 150,
		&document.ImageConfig{
			Position:  document.ImagePositionInline,
			Alignment: document.AlignCenter,
			AltText:   "第三张图片",
			Title:     "绿色图片",
		},
	)
	if err != nil {
		log.Fatalf("添加第三张图片失败: %v", err)
	}

	doc2.AddParagraph("上图：第三张图片（绿色，新添加）")
	
	doc2.AddParagraph("文档已修改完成。原有的两张图片应该仍然存在，同时新增了一张图片。")
	
	fmt.Println("   ✓ 文档内容修改完成")

	// 步骤4: 保存修改后的文档
	fmt.Println("\n步骤4: 保存修改后的文档...")
	modifiedFile := "examples/output/image_persistence_modified.docx"
	err = doc2.Save(modifiedFile)
	if err != nil {
		log.Fatalf("保存修改后的文档失败: %v", err)
	}
	fmt.Printf("   ✓ 修改后的文档已保存: %s\n", modifiedFile)

	// 步骤5: 重新打开修改后的文档，验证所有图片都存在
	fmt.Println("\n步骤5: 验证修改后的文档...")
	doc3, err := document.Open(modifiedFile)
	if err != nil {
		log.Fatalf("打开修改后的文档失败: %v", err)
	}
	fmt.Println("   ✓ 修改后的文档打开成功")

	// 验证所有三张图片都存在
	finalImageCount := 0
	expectedImages := []string{"image1.png", "image2.png", "image3.png"}
	foundImages := make(map[string]bool)

	for partName := range doc3.GetParts() {
		if len(partName) > 11 && partName[:11] == "word/media/" {
			finalImageCount++
			imageName := partName[11:] // 提取文件名
			foundImages[imageName] = true
			fmt.Printf("   - 找到图片: %s\n", partName)
		}
	}

	// 验证每张图片都存在
	allImagesPresent := true
	for _, expectedImage := range expectedImages {
		if !foundImages[expectedImage] {
			fmt.Printf("   ✗ 图片丢失: %s\n", expectedImage)
			allImagesPresent = false
		}
	}

	if !allImagesPresent {
		log.Fatalf("【错误】部分图片在保存后丢失！")
	}

	if finalImageCount != 3 {
		log.Fatalf("【错误】修改后的文档应该包含3张图片，但只找到 %d 张", finalImageCount)
	}

	fmt.Printf("   ✓ 确认所有 %d 张图片都正确保存\n", finalImageCount)

	// 成功完成
	fmt.Println("\n===================")
	fmt.Println("✓ 图片持久性测试成功！")
	fmt.Println()
	fmt.Println("演示结果：")
	fmt.Printf("  - 原始文档: %s (包含2张图片)\n", originalFile)
	fmt.Printf("  - 修改后的文档: %s (包含3张图片)\n", modifiedFile)
	fmt.Println()
	fmt.Println("说明：")
	fmt.Println("  1. 原有的2张图片（红色和蓝色）在打开和重新保存后没有丢失")
	fmt.Println("  2. 新添加的1张图片（绿色）正确保存")
	fmt.Println("  3. 所有图片在Word中可以正常显示")
	fmt.Println()
	fmt.Println("请打开以下文件验证：")
	fmt.Printf("  %s\n", originalFile)
	fmt.Printf("  %s\n", modifiedFile)
}

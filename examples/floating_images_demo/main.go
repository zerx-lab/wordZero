package main

import (
"bytes"
"fmt"
"image"
"image/color"
"image/png"
"log"

"github.com/ZeroHawkeye/wordZero/pkg/document"
)

// createTestImage 创建一个简单的测试图片
func createTestImage(width, height int, colorR, colorG, colorB uint8) []byte {
img := image.NewRGBA(image.Rect(0, 0, width, height))
for y := 0; y < height; y++ {
for x := 0; x < width; x++ {
img.SetRGBA(x, y, color.RGBA{colorR, colorG, colorB, 255})
}
}
var buf bytes.Buffer
png.Encode(&buf, img)
return buf.Bytes()
}

func main() {
fmt.Println("演示浮动图片修复...")

doc := document.New()

// 添加标题
title := doc.AddParagraph("浮动图片功能测试文档")
title.SetAlignment(document.AlignCenter)

// 添加介绍段落
doc.AddParagraph("本文档演示了左浮动和右浮动图片的各种文字环绕方式。")
doc.AddParagraph("")

// 测试1: 左浮动 + 紧密环绕
doc.AddParagraph("1. 左浮动 + 紧密环绕")
imageData1 := createTestImage(100, 100, 255, 0, 0) // 红色
config1 := &document.ImageConfig{
Position: document.ImagePositionFloatLeft,
WrapText: document.ImageWrapTight,
Size: &document.ImageSize{
Width:  30,
Height: 30,
},
AltText: "左浮动-紧密环绕",
Title:   "左浮动图片",
}
_, err := doc.AddImageFromData(imageData1, "left_tight.png", document.ImageFormatPNG, 100, 100, config1)
if err != nil {
log.Fatalf("添加左浮动图片失败: %v", err)
}
doc.AddParagraph("这是一段文字，用于演示文字如何围绕左侧的浮动图片进行环绕。紧密环绕模式会让文字紧贴图片边缘。")
doc.AddParagraph("")

// 测试2: 右浮动 + 四周环绕
doc.AddParagraph("2. 右浮动 + 四周环绕")
imageData2 := createTestImage(100, 100, 0, 255, 0) // 绿色
config2 := &document.ImageConfig{
Position: document.ImagePositionFloatRight,
WrapText: document.ImageWrapSquare,
Size: &document.ImageSize{
Width:  30,
Height: 30,
},
AltText: "右浮动-四周环绕",
Title:   "右浮动图片",
}
_, err = doc.AddImageFromData(imageData2, "right_square.png", document.ImageFormatPNG, 100, 100, config2)
if err != nil {
log.Fatalf("添加右浮动图片失败: %v", err)
}
doc.AddParagraph("这是一段文字，用于演示文字如何围绕右侧的浮动图片进行环绕。四周环绕模式会在图片四周留出一定空间。")
doc.AddParagraph("")

// 测试3: 左浮动 + 上下环绕
doc.AddParagraph("3. 左浮动 + 上下环绕")
imageData3 := createTestImage(100, 100, 0, 0, 255) // 蓝色
config3 := &document.ImageConfig{
Position: document.ImagePositionFloatLeft,
WrapText: document.ImageWrapTopAndBottom,
Size: &document.ImageSize{
Width:  40,
Height: 30,
},
AltText: "左浮动-上下环绕",
Title:   "上下环绕图片",
}
_, err = doc.AddImageFromData(imageData3, "left_topbottom.png", document.ImageFormatPNG, 100, 100, config3)
if err != nil {
log.Fatalf("添加浮动图片失败: %v", err)
}
doc.AddParagraph("这是一段文字，用于演示上下环绕模式。在这种模式下，文字只出现在图片的上方和下方，不会出现在图片的左右两侧。")
doc.AddParagraph("")

// 测试4: 右浮动 + 无环绕
doc.AddParagraph("4. 右浮动 + 无环绕")
imageData4 := createTestImage(100, 100, 255, 255, 0) // 黄色
config4 := &document.ImageConfig{
Position: document.ImagePositionFloatRight,
WrapText: document.ImageWrapNone,
Size: &document.ImageSize{
Width:  25,
Height: 25,
},
AltText: "右浮动-无环绕",
Title:   "无环绕图片",
}
_, err = doc.AddImageFromData(imageData4, "right_none.png", document.ImageFormatPNG, 100, 100, config4)
if err != nil {
log.Fatalf("添加浮动图片失败: %v", err)
}
doc.AddParagraph("这是一段文字，用于演示无环绕模式。在这种模式下，图片会遮挡文字，没有环绕效果。")
doc.AddParagraph("")

// 添加总结
doc.AddParagraph("")
summary := doc.AddParagraph("总结：所有四种浮动图片配置都已成功生成，文档可以正常打开。")
summary.SetAlignment(document.AlignCenter)

// 保存文档
filename := "floating_images_demo.docx"
err = doc.Save(filename)
if err != nil {
log.Fatalf("保存文档失败: %v", err)
}

fmt.Printf("\n✓ 文档已成功保存: %s\n", filename)
fmt.Println("✓ 文档包含以下浮动图片配置:")
fmt.Println("  - 左浮动 + 紧密环绕")
fmt.Println("  - 右浮动 + 四周环绕")
fmt.Println("  - 左浮动 + 上下环绕")
fmt.Println("  - 右浮动 + 无环绕")
fmt.Println("\n请用Microsoft Word打开文件验证修复效果！")
}

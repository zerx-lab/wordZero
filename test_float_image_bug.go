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

// loadImageFromFile 从文件加载图片数据
func loadImageFromFile(filePath string) ([]byte, error) {
	// 打开图片文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 解码图片（自动识别格式）
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// 编码为PNG格式的字节数组
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// prepareTemplateData 准备模板渲染数据
func prepareTemplateData() *document.TemplateData {
	data := document.NewTemplateData()

	// 设置基础变量
	data.SetVariable("customerName", "xxxx")

	// 创建一个简单的测试图片（红色100x100像素）
	testImage := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			testImage.SetRGBA(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, testImage)
	imageData := buf.Bytes()

	// 配置图片样式 - 左浮动
	imageConfig := &document.ImageConfig{
		Position:  document.ImagePositionFloatLeft,
		Alignment: document.AlignCenter,
		Size: &document.ImageSize{
			Width:           23.6,
			Height:          13,
			KeepAspectRatio: true,
		},
		WrapText: document.ImageWrapTight, // 紧密环绕
		AltText:  "YT",
		Title:    "YT",
	}
	
	// 设置图片数据（支持二进制数据）
	data.SetImageFromData("productImage", imageData, imageConfig)

	return data
}

func main() {
	fmt.Println("测试浮动图片问题...")

	// 创建文档
	doc := document.New()

	// 准备模板数据
	data := prepareTemplateData()

	// 获取图片数据
	if imgData, exists := data.GetImage("productImage"); exists {
		// 添加浮动图片
		fmt.Printf("添加图片，位置: %s, 环绕: %s\n", imgData.Config.Position, imgData.Config.WrapText)
		
		_, err := doc.AddImageFromData(
			imgData.Data,
			"test.png",
			document.ImageFormatPNG,
			100,
			100,
			imgData.Config,
		)
		if err != nil {
			log.Fatalf("添加图片失败: %v", err)
		}
	}

	// 保存文档
	filename := "test_float_image.docx"
	err := doc.Save(filename)
	if err != nil {
		log.Fatalf("保存文档失败: %v", err)
	}

	fmt.Printf("文档已保存: %s\n", filename)
	fmt.Println("请用Microsoft Word打开此文件检查是否有错误")
}

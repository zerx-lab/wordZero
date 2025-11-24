package document

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// TestFloatingImageLeftWithTightWrap 测试左浮动图片 + 紧密环绕
func TestFloatingImageLeftWithTightWrap(t *testing.T) {
	doc := New()

	// 创建测试图片
	imageData := createTestImageRGBA(100, 100)

	config := &ImageConfig{
		Position: ImagePositionFloatLeft,
		WrapText: ImageWrapTight,
		Size: &ImageSize{
			Width:  23.6,
			Height: 13,
		},
		AltText: "左浮动测试图片",
		Title:   "测试",
	}

	_, err := doc.AddImageFromData(imageData, "test.png", ImageFormatPNG, 100, 100, config)
	if err != nil {
		t.Fatalf("添加左浮动图片失败: %v", err)
	}

	// 保存并验证
	filename := "test_float_left_tight.docx"
	err = doc.Save(filename)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}
	defer os.Remove(filename)

	// 重新打开文档验证
	doc2, err := Open(filename)
	if err != nil {
		t.Fatalf("打开保存的文档失败: %v", err)
	}

	if len(doc2.Body.Elements) == 0 {
		t.Fatal("文档中没有元素")
	}

	t.Logf("✓ 左浮动 + 紧密环绕图片测试通过")
}

// TestFloatingImageRightWithSquareWrap 测试右浮动图片 + 四周环绕
func TestFloatingImageRightWithSquareWrap(t *testing.T) {
	doc := New()

	imageData := createTestImageRGBA(100, 100)

	config := &ImageConfig{
		Position: ImagePositionFloatRight,
		WrapText: ImageWrapSquare,
		Size: &ImageSize{
			Width:  30,
			Height: 20,
		},
		AltText: "右浮动测试图片",
		Title:   "测试",
	}

	_, err := doc.AddImageFromData(imageData, "test.png", ImageFormatPNG, 100, 100, config)
	if err != nil {
		t.Fatalf("添加右浮动图片失败: %v", err)
	}

	filename := "test_float_right_square.docx"
	err = doc.Save(filename)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}
	defer os.Remove(filename)

	// 重新打开文档验证
	doc2, err := Open(filename)
	if err != nil {
		t.Fatalf("打开保存的文档失败: %v", err)
	}

	if len(doc2.Body.Elements) == 0 {
		t.Fatal("文档中没有元素")
	}

	t.Logf("✓ 右浮动 + 四周环绕图片测试通过")
}

// TestFloatingImageWithTopAndBottomWrap 测试浮动图片 + 上下环绕
func TestFloatingImageWithTopAndBottomWrap(t *testing.T) {
	doc := New()

	imageData := createTestImageRGBA(100, 100)

	config := &ImageConfig{
		Position: ImagePositionFloatLeft,
		WrapText: ImageWrapTopAndBottom,
		Size: &ImageSize{
			Width:  40,
			Height: 30,
		},
		AltText: "上下环绕测试图片",
		Title:   "测试",
	}

	_, err := doc.AddImageFromData(imageData, "test.png", ImageFormatPNG, 100, 100, config)
	if err != nil {
		t.Fatalf("添加浮动图片失败: %v", err)
	}

	filename := "test_float_topbottom.docx"
	err = doc.Save(filename)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}
	defer os.Remove(filename)

	// 重新打开文档验证
	doc2, err := Open(filename)
	if err != nil {
		t.Fatalf("打开保存的文档失败: %v", err)
	}

	if len(doc2.Body.Elements) == 0 {
		t.Fatal("文档中没有元素")
	}

	t.Logf("✓ 浮动 + 上下环绕图片测试通过")
}

// TestFloatingImageWithNoWrap 测试浮动图片 + 无环绕
func TestFloatingImageWithNoWrap(t *testing.T) {
	doc := New()

	imageData := createTestImageRGBA(100, 100)

	config := &ImageConfig{
		Position: ImagePositionFloatRight,
		WrapText: ImageWrapNone,
		Size: &ImageSize{
			Width:  25,
			Height: 15,
		},
		AltText: "无环绕测试图片",
		Title:   "测试",
	}

	_, err := doc.AddImageFromData(imageData, "test.png", ImageFormatPNG, 100, 100, config)
	if err != nil {
		t.Fatalf("添加浮动图片失败: %v", err)
	}

	filename := "test_float_nowrap.docx"
	err = doc.Save(filename)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}
	defer os.Remove(filename)

	// 重新打开文档验证
	doc2, err := Open(filename)
	if err != nil {
		t.Fatalf("打开保存的文档失败: %v", err)
	}

	if len(doc2.Body.Elements) == 0 {
		t.Fatal("文档中没有元素")
	}

	t.Logf("✓ 浮动 + 无环绕图片测试通过")
}

// TestMultipleFloatingImages 测试多个浮动图片
func TestMultipleFloatingImages(t *testing.T) {
	doc := New()

	// 添加文本段落
	doc.AddParagraph("这是一个包含多个浮动图片的文档测试。")

	imageData := createTestImageRGBA(80, 80)

	// 添加左浮动图片
	config1 := &ImageConfig{
		Position: ImagePositionFloatLeft,
		WrapText: ImageWrapSquare,
		Size: &ImageSize{
			Width:  20,
			Height: 20,
		},
	}
	_, err := doc.AddImageFromData(imageData, "test1.png", ImageFormatPNG, 80, 80, config1)
	if err != nil {
		t.Fatalf("添加第一个浮动图片失败: %v", err)
	}

	// 添加更多文本
	doc.AddParagraph("第一个图片已添加。")

	// 添加右浮动图片
	config2 := &ImageConfig{
		Position: ImagePositionFloatRight,
		WrapText: ImageWrapTight,
		Size: &ImageSize{
			Width:  20,
			Height: 20,
		},
	}
	_, err = doc.AddImageFromData(imageData, "test2.png", ImageFormatPNG, 80, 80, config2)
	if err != nil {
		t.Fatalf("添加第二个浮动图片失败: %v", err)
	}

	doc.AddParagraph("第二个图片也已添加。")

	filename := "test_multiple_float.docx"
	err = doc.Save(filename)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}
	defer os.Remove(filename)

	// 验证文档可以正常打开
	doc2, err := Open(filename)
	if err != nil {
		t.Fatalf("打开保存的文档失败: %v", err)
	}

	// 应该有5个元素（3个文本段落 + 2个图片段落）
	const expectedElements = 5 // 3 text paragraphs + 2 image paragraphs
	if len(doc2.Body.Elements) != expectedElements {
		t.Errorf("期望%d个元素（3个文本段落+2个图片段落），实际 %d 个", expectedElements, len(doc2.Body.Elements))
	}

	t.Logf("✓ 多个浮动图片测试通过")
}

// createTestImageRGBA 创建一个彩色测试图片
func createTestImageRGBA(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 创建渐变色彩
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8(x * 255 / width)
			g := uint8(y * 255 / height)
			b := uint8((x + y) * 255 / (width + height))
			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

// TestInlineImageNotAffected 确保嵌入式图片不受修复影响
func TestInlineImageNotAffected(t *testing.T) {
	doc := New()

	imageData := createTestImageRGBA(100, 100)

	// 测试嵌入式图片（默认）
	config := &ImageConfig{
		Position:  ImagePositionInline,
		Alignment: AlignCenter,
		Size: &ImageSize{
			Width:  30,
			Height: 30,
		},
	}

	_, err := doc.AddImageFromData(imageData, "test.png", ImageFormatPNG, 100, 100, config)
	if err != nil {
		t.Fatalf("添加嵌入式图片失败: %v", err)
	}

	filename := "test_inline_image.docx"
	err = doc.Save(filename)
	if err != nil {
		t.Fatalf("保存文档失败: %v", err)
	}
	defer os.Remove(filename)

	// 重新打开验证
	doc2, err := Open(filename)
	if err != nil {
		t.Fatalf("打开保存的文档失败: %v", err)
	}

	if len(doc2.Body.Elements) == 0 {
		t.Fatal("文档中没有元素")
	}

	t.Logf("✓ 嵌入式图片未受影响测试通过")
}

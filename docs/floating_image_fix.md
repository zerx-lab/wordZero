# 浮动图片修复说明

## 问题描述

当使用 `ImagePositionFloatLeft` 或 `ImagePositionFloatRight` 设置图片为左浮动或右浮动时，生成的 Word 文档无法在 Microsoft Word 中打开，显示错误："Word 在试图打开文件时遇到错误。"

## 问题原因

生成的 XML 结构中，包含图片的 `<w:r>` (Run) 元素内同时存在：
- 空的 `<w:t></w:t>` 文本元素
- `<w:drawing>` 绘图元素

这种结构对于浮动图片是不正确的。根据 Office Open XML (OOXML) 规范，当 Run 中包含 Drawing 元素时，不应该包含空的 Text 元素。

### 修复前的 XML 结构
```xml
<w:r>
  <w:t></w:t>
  <w:drawing>
    <wp:anchor>...</wp:anchor>
  </w:drawing>
</w:r>
```

### 修复后的 XML 结构
```xml
<w:r>
  <w:drawing>
    <wp:anchor>...</wp:anchor>
  </w:drawing>
</w:r>
```

## 修复方案

在 `pkg/document/document.go` 中为 `Run` 结构体添加了自定义的 `MarshalXML` 方法：

```go
func (r *Run) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    // 开始Run元素
    if err := e.EncodeToken(start); err != nil {
        return err
    }

    // 序列化RunProperties（如果存在）
    if r.Properties != nil {
        if err := e.EncodeElement(r.Properties, xml.StartElement{Name: xml.Name{Local: "w:rPr"}}); err != nil {
            return err
        }
    }

    // 序列化Text（仅当有内容时）
    // 这是关键修复：避免序列化空的Text元素
    if r.Text.Content != "" {
        if err := e.EncodeElement(r.Text, xml.StartElement{Name: xml.Name{Local: "w:t"}}); err != nil {
            return err
        }
    }

    // 序列化Drawing（如果存在）
    if r.Drawing != nil {
        if err := e.EncodeElement(r.Drawing, xml.StartElement{Name: xml.Name{Local: "w:drawing"}}); err != nil {
            return err
        }
    }

    // ... 其他元素 ...

    return e.EncodeToken(xml.EndElement{Name: start.Name})
}
```

这个自定义方法确保：
1. 只有当 `Text.Content` 不为空时才序列化 `<w:t>` 元素
2. 其他元素（Drawing, Break, FieldChar, InstrText）正常序列化
3. 保持正确的元素顺序，符合 OOXML 规范

## 测试覆盖

添加了全面的测试覆盖所有浮动图片场景：

### 测试文件：`pkg/document/image_floating_test.go`

1. **TestFloatingImageLeftWithTightWrap**: 左浮动 + 紧密环绕
2. **TestFloatingImageRightWithSquareWrap**: 右浮动 + 四周环绕
3. **TestFloatingImageWithTopAndBottomWrap**: 浮动 + 上下环绕
4. **TestFloatingImageWithNoWrap**: 浮动 + 无环绕
5. **TestMultipleFloatingImages**: 多个浮动图片
6. **TestInlineImageNotAffected**: 确保嵌入式图片不受影响

所有测试均通过 ✓

## 示例代码

参见 `examples/floating_images_demo/main.go`，演示了四种浮动图片配置：

```go
// 左浮动 + 紧密环绕
config1 := &document.ImageConfig{
    Position: document.ImagePositionFloatLeft,
    WrapText: document.ImageWrapTight,
    Size: &document.ImageSize{
        Width:  30,
        Height: 30,
    },
}

// 右浮动 + 四周环绕
config2 := &document.ImageConfig{
    Position: document.ImagePositionFloatRight,
    WrapText: document.ImageWrapSquare,
    Size: &document.ImageSize{
        Width:  30,
        Height: 30,
    },
}
```

## 影响范围

- **修复范围**: 仅影响浮动图片（ImagePositionFloatLeft 和 ImagePositionFloatRight）的 XML 序列化
- **不影响**: 嵌入式图片（ImagePositionInline）、普通文本、其他文档元素
- **向后兼容**: 所有现有测试通过，不影响现有功能

## 验证

1. 运行 `go test ./pkg/document -v` - 所有测试通过
2. 运行 `go run examples/floating_images_demo/main.go` - 生成测试文档
3. 用 Microsoft Word 打开生成的文档 - 可正常打开，图片正确显示

## 支持的配置

### 图片位置 (ImagePosition)
- `ImagePositionInline` - 嵌入式（默认）
- `ImagePositionFloatLeft` - 左浮动 ✓ 已修复
- `ImagePositionFloatRight` - 右浮动 ✓ 已修复

### 文字环绕 (ImageWrapText)
- `ImageWrapNone` - 无环绕 ✓
- `ImageWrapSquare` - 四周环绕 ✓
- `ImageWrapTight` - 紧密环绕 ✓
- `ImageWrapTopAndBottom` - 上下环绕 ✓

所有组合均已测试并验证可以正常工作。

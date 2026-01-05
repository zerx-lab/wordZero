# WordZero Markdown转换包

`pkg/markdown` 包提供了 Markdown 和 Word 文档之间的双向转换功能。

## 功能特性

### Markdown → Word 转换
- 基于 goldmark 解析引擎
- 支持 GitHub Flavored Markdown (GFM)
- 支持标题、格式化文本、列表、表格、图片、链接等
- 可配置的转换选项

### Word → Markdown 转换 (新增)
- 支持将 Word 文档反向导出为 Markdown
- 保持文档结构和格式
- 支持图片导出
- 多种导出配置选项

## 基本使用

### Word 到 Markdown 转换

```go
package main

import (
    "fmt"
    "github.com/zerx-lab/wordZero/pkg/markdown"
)

func main() {
    // 创建导出器
    exporter := markdown.NewExporter(markdown.DefaultExportOptions())
    
    // 导出Word文档为Markdown
    err := exporter.ExportToFile("document.docx", "output.md", nil)
    if err != nil {
        fmt.Printf("导出失败: %v\n", err)
        return
    }
    
    fmt.Println("Word文档已成功转换为Markdown!")
}
```

### Markdown 到 Word 转换

```go
package main

import (
    "fmt"
    "github.com/zerx-lab/wordZero/pkg/markdown"
)

func main() {
    // 创建转换器
    converter := markdown.NewConverter(markdown.DefaultOptions())
    
    // 转换Markdown为Word文档
    err := converter.ConvertFile("input.md", "output.docx", nil)
    if err != nil {
        fmt.Printf("转换失败: %v\n", err)
        return
    }
    
    fmt.Println("Markdown已成功转换为Word文档!")
}
```

### 双向转换器

```go
package main

import (
    "fmt"
    "github.com/zerx-lab/wordZero/pkg/markdown"
)

func main() {
    // 创建双向转换器
    converter := markdown.NewBidirectionalConverter(
        markdown.DefaultOptions(),      // Markdown→Word选项
        markdown.DefaultExportOptions(), // Word→Markdown选项
    )
    
    // 自动检测文件类型并转换
    err := converter.AutoConvert("input.docx", "output.md")
    if err != nil {
        fmt.Printf("转换失败: %v\n", err)
        return
    }
    
    fmt.Println("文档转换完成!")
}
```

## 高级配置

### Word 到 Markdown 导出选项

```go
options := &markdown.ExportOptions{
    UseGFMTables:      true,  // 使用GitHub风味Markdown表格
    ExtractImages:     true,  // 导出图片文件
    ImageOutputDir:    "images/", // 图片输出目录
    PreserveFootnotes: true,  // 保留脚注
    UseSetext:         true,  // 使用Setext样式标题
    IncludeMetadata:   true,  // 包含文档元数据
    ProgressCallback: func(current, total int) {
        fmt.Printf("进度: %d/%d\n", current, total)
    },
}

exporter := markdown.NewExporter(options)
```

### Markdown 到 Word 转换选项

```go
options := &markdown.ConvertOptions{
    EnableGFM:         true,     // 启用GitHub风味Markdown
    EnableFootnotes:   true,     // 启用脚注支持
    EnableTables:      true,     // 启用表格支持
    EnableMath:        true,     // 启用数学公式支持（LaTeX语法）
    DefaultFontFamily: "Calibri", // 默认字体
    DefaultFontSize:   11.0,     // 默认字号
    GenerateTOC:       true,     // 生成目录
    TOCMaxLevel:       3,        // 目录最大级别
}

converter := markdown.NewConverter(options)
```

## 支持的转换映射

### Word → Markdown

| Word元素 | Markdown语法 | 说明 |
|----------|-------------|------|
| Heading1-6 | `# 标题` | 标题级别对应 |
| 粗体 | `**粗体**` | 文本格式 |
| 斜体 | `*斜体*` | 文本格式 |
| 删除线 | `~~删除线~~` | 文本格式 |
| 代码 | `` `代码` `` | 行内代码 |
| 代码块 | ```` 代码块 ```` | 代码块 |
| 超链接 | `[链接](url)` | 链接转换 |
| 图片 | `![图片](src)` | 图片引用 |
| 表格 | `\| 表格 \|` | GFM表格 |
| 列表 | `- 项目` | 列表项 |

### Markdown → Word

| Markdown语法 | Word元素 | 实现方式 |
|-------------|----------|----------|
| `# 标题` | Heading1样式 | `AddHeadingParagraph()` |
| `**粗体**` | 粗体格式 | `RunProperties.Bold` |
| `*斜体*` | 斜体格式 | `RunProperties.Italic` |
| `` `代码` `` | 代码样式 | 等宽字体 |
| `[链接](url)` | 超链接 | `AddHyperlink()` |
| `![图片](src)` | 图片 | `AddImageFromFile()` |
| `\| 表格 \|` | Word表格 | `AddTable()` |
| `- 列表` | 项目符号列表 | `AddBulletList()` |
| `$公式$` | 数学公式 | Cambria Math字体 |
| `$$公式$$` | 块级数学公式 | 居中显示 |

## 批量转换

```go
// 批量Markdown转Word
converter := markdown.NewConverter(markdown.DefaultOptions())
inputs := []string{"doc1.md", "doc2.md", "doc3.md"}
err := converter.BatchConvert(inputs, "output/", nil)

// 批量Word转Markdown
exporter := markdown.NewExporter(markdown.DefaultExportOptions())
inputs := []string{"doc1.docx", "doc2.docx", "doc3.docx"}
err := exporter.BatchExport(inputs, "markdown/", nil)
```

## 错误处理

```go
options := &markdown.ExportOptions{
    StrictMode: true,  // 严格模式
    IgnoreErrors: false, // 不忽略错误
    ErrorCallback: func(err error) {
        fmt.Printf("转换错误: %v\n", err)
    },
}
```

## 兼容性说明

- 该包与现有的 `pkg/document` 包完全兼容
- 不修改任何现有API
- 可以与现有代码无缝集成
- 支持所有现有的Word文档操作功能

## 注意事项

1. Word到Markdown转换会丢失某些Word特有的格式信息
2. 复杂的表格布局可能需要手动调整
3. 图片需要单独处理导出
4. 某些Word样式在Markdown中没有直接对应
5. 数学公式转换使用Unicode字符和Cambria Math字体，支持常见的LaTeX语法

## 数学公式支持

### 行内公式
使用单个美元符号包裹：`$E = mc^2$`

### 块级公式
使用双美元符号包裹：
```
$$
x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}
$$
```

### 支持的LaTeX语法
- 希腊字母：`\alpha`, `\beta`, `\gamma`, `\pi`, `\sigma` 等
- 运算符：`\times`, `\div`, `\pm`, `\leq`, `\geq`, `\neq` 等
- 上下标：`x^2`, `x_i`, `x^{n+1}`, `x_{i,j}` 等
- 分数：`\frac{a}{b}`
- 根号：`\sqrt{x}`, `\sqrt[3]{x}`
- 特殊符号：`\infty`, `\sum`, `\int`, `\partial`, `\nabla` 等
- 箭头：`\rightarrow`, `\leftarrow`, `\Rightarrow` 等

## 未来计划

- [x] 数学公式支持
- [ ] Mermaid图表转换
- [ ] 更好的列表嵌套支持
- [ ] 自定义样式映射
- [ ] 命令行工具 
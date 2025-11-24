<div align="center">
  <img src="docs/logo-banner.svg" alt="WordZero Logo" width="400"/>
  
  <h1>WordZero - Golang Word Document Library</h1>
</div>

<div align="center">
  
[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-Passing-green.svg)](#testing)
[![Benchmark](https://img.shields.io/badge/Benchmark-Go%202.62ms%20%7C%20JS%209.63ms%20%7C%20Python%2055.98ms-success.svg)](https://github.com/ZeroHawkeye/wordZero/wiki/en-Performance-Benchmarks)
[![Performance](https://img.shields.io/badge/Performance-Golang%20Winner-brightgreen.svg)](https://github.com/ZeroHawkeye/wordZero/wiki/en-Performance-Benchmarks)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/ZeroHawkeye/wordZero)

</div>

**English** | [中文](README_zh.md)

## Project Introduction

WordZero is a Golang-based Word document manipulation library that provides basic document creation and modification operations. This library follows the latest Office Open XML (OOXML) specifications and focuses on supporting modern Word document format (.docx).

### Core Features

- 🚀 **Complete Document Operations**: Create, read, and modify Word documents
- 🎨 **Rich Style System**: 18 predefined styles with custom style and inheritance support
- 📝 **Text Formatting**: Full support for fonts, sizes, colors, bold, italic, and more
- 📐 **Paragraph Format**: Alignment, spacing, indentation, and other paragraph properties
- 🏷️ **Heading Navigation**: Complete support for Heading1-9 styles, recognizable by Word navigation pane
- 📊 **Table Functionality**: Complete table creation, editing, styling, and iterator support
- 📄 **Page Settings**: Page size, margins, headers/footers, and professional layout features
- 🔧 **Advanced Features**: Table of contents generation, footnotes/endnotes, list numbering, template engine, etc.
- 🎯 **Template Inheritance**: Support for base templates and block override mechanisms for template reuse and extension
- 📝 **Header/Footer Templates**: Support for template variables in headers and footers for dynamic content replacement
- ⚡ **Excellent Performance**: Zero-dependency pure Go implementation, average 2.62ms processing speed, 3.7x faster than JavaScript, 21x faster than Python
- 🔧 **Easy to Use**: Clean API design with fluent interface support

## Related Recommended Projects

### Excel Document Operations - Excelize

If you need to work with Excel documents, we highly recommend [**Excelize**](https://github.com/qax-os/excelize) —— the most popular Go library for Excel operations:

- ⭐ **19.2k+ GitHub Stars** - The most popular Excel processing library in the Go ecosystem
- 📊 **Complete Excel Support** - Supports all modern Excel formats including XLAM/XLSM/XLSX/XLTM/XLTX
- 🎯 **Feature Rich** - Charts, pivot tables, images, streaming APIs, and more
- 🚀 **High Performance** - Streaming read/write APIs optimized for large datasets
- 🔧 **Easy Integration** - Perfect complement to WordZero for complete Office document processing solutions

**Perfect Combination**: WordZero handles Word documents, Excelize handles Excel documents, together providing comprehensive Office document manipulation capabilities for your Go projects.

```go
// WordZero + Excelize combination example
import (
    "github.com/ZeroHawkeye/wordZero/pkg/document"
    "github.com/qax-os/excelize/v2"
)

// Create Word report
doc := document.New()
doc.AddParagraph("Data Analysis Report").SetStyle(style.StyleHeading1)

// Create Excel data sheet
xlsx := excelize.NewFile()
xlsx.SetCellValue("Sheet1", "A1", "Data Item")
xlsx.SetCellValue("Sheet1", "B1", "Value")
```

## Installation

```bash
go get github.com/ZeroHawkeye/wordZero
```

### Version Notes

We recommend using versioned installation:

```bash
# Install latest version
go get github.com/ZeroHawkeye/wordZero@latest

# Install specific version
go get github.com/ZeroHawkeye/wordZero@v1.3.7
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/ZeroHawkeye/wordZero/pkg/document"
    "github.com/ZeroHawkeye/wordZero/pkg/style"
)

func main() {
    // Create new document
    doc := document.New()
    
    // Add title
    titlePara := doc.AddParagraph("WordZero Usage Example")
    titlePara.SetStyle(style.StyleHeading1)
    
    // Add body paragraph
    para := doc.AddParagraph("This is a document example created using WordZero.")
    para.SetFontFamily("Arial")
    para.SetFontSize(12)
    para.SetColor("333333")
    
    // Create table
    tableConfig := &document.TableConfig{
        Rows:    3,
        Columns: 3,
    }
    table := doc.AddTable(tableConfig)
    table.SetCellText(0, 0, "Header1")
    table.SetCellText(0, 1, "Header2")
    table.SetCellText(0, 2, "Header3")
    
    // Save document
    if err := doc.Save("example.docx"); err != nil {
        log.Fatal(err)
    }
}
```

### Template Inheritance Feature Example

```go
// Create base template
engine := document.NewTemplateEngine()
baseTemplate := `{{companyName}} Work Report

{{#block "summary"}}
Default summary content
{{/block}}

{{#block "content"}}
Default main content
{{/block}}`

engine.LoadTemplate("base_report", baseTemplate)

// Create extended template, override specific blocks
salesTemplate := `{{extends "base_report"}}

{{#block "summary"}}
Sales Performance Summary: Achieved {{achievement}}% this month
{{/block}}

{{#block "content"}}
Sales Details:
- Total Sales: {{totalSales}}
- New Customers: {{newCustomers}}
{{/block}}`

engine.LoadTemplate("sales_report", salesTemplate)

// Render template
data := document.NewTemplateData()
data.SetVariable("companyName", "WordZero Tech")
data.SetVariable("achievement", "125")
data.SetVariable("totalSales", "1,850,000")
data.SetVariable("newCustomers", "45")

doc, _ := engine.RenderTemplateToDocument("sales_report", data)
doc.Save("sales_report.docx")
```

### Image Placeholder Template Feature Example ✨ **New**

```go
package main

import (
    "log"
    "github.com/ZeroHawkeye/wordZero/pkg/document"
)

func main() {
    // Create template with image placeholders
    engine := document.NewTemplateEngine()
    template := `Company: {{companyName}}

{{#image companyLogo}}

Project Report: {{projectName}}

Status: {{#if isCompleted}}Completed{{else}}In Progress{{/if}}

{{#image statusChart}}

Team Members:
{{#each teamMembers}}
- {{name}}: {{role}}
{{/each}}`

    engine.LoadTemplate("project_report", template)

    // Prepare template data
    data := document.NewTemplateData()
    data.SetVariable("companyName", "WordZero Tech")
    data.SetVariable("projectName", "Document Processing System")
    data.SetCondition("isCompleted", true)
    
    // Set team members list
    data.SetList("teamMembers", []interface{}{
        map[string]interface{}{"name": "Alice", "role": "Lead Developer"},
        map[string]interface{}{"name": "Bob", "role": "Frontend Developer"},
    })
    
    // Configure and set images
    logoConfig := &document.ImageConfig{
        Width:     100,
        Height:    50,
        Alignment: document.AlignCenter,
    }
    data.SetImage("companyLogo", "assets/logo.png", logoConfig)
    
    chartConfig := &document.ImageConfig{
        Width:       200,
        Height:      150,
        Alignment:   document.AlignCenter,
        AltText:     "Project Status Chart",
        Title:       "Current Project Status",
    }
    data.SetImage("statusChart", "assets/chart.png", chartConfig)
    
    // Render template to document
    doc, err := engine.RenderTemplateToDocument("project_report", data)
    if err != nil {
        log.Fatal(err)
    }
    
    // Save document
    err = doc.Save("project_report.docx")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Markdown to Word Feature Example ✨ **New**

```go
package main

import (
    "log"
    "github.com/ZeroHawkeye/wordZero/pkg/markdown"
)

func main() {
    // Create Markdown converter
    converter := markdown.NewConverter(markdown.DefaultOptions())
    
    // Markdown content
    markdownText := `# WordZero Markdown Conversion Example

Welcome to WordZero's **Markdown to Word** conversion feature!

## Supported Syntax

### Text Formatting
- **Bold text**
- *Italic text*
- ` + "`Inline code`" + `

### Lists
1. Ordered list item 1
2. Ordered list item 2

- Unordered list item A
- Unordered list item B

### Quotes and Code

> This is blockquote content
> Supporting multiple lines

` + "```" + `go
// Code block example
func main() {
    fmt.Println("Hello, WordZero!")
}
` + "```" + `

---

Conversion complete!`

    // Convert to Word document
    doc, err := converter.ConvertString(markdownText, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Save Word document
    err = doc.Save("markdown_example.docx")
    if err != nil {
        log.Fatal(err)
    }
    
    // File conversion
    err = converter.ConvertFile("input.md", "output.docx", nil)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Document Pagination and Paragraph Deletion Example ✨ **New**

```go
package main

import (
    "log"
    "github.com/ZeroHawkeye/wordZero/pkg/document"
)

func main() {
    doc := document.New()
    
    // Add first page content
    doc.AddHeadingParagraph("Chapter 1: Introduction", 1)
    doc.AddParagraph("This is the content of chapter 1.")
    
    // Add page break to start a new page
    doc.AddPageBreak()
    
    // Add second page content
    doc.AddHeadingParagraph("Chapter 2: Main Content", 1)
    tempPara := doc.AddParagraph("This is a temporary paragraph.")
    doc.AddParagraph("This is the content of chapter 2.")
    
    // Delete temporary paragraph
    doc.RemoveParagraph(tempPara)
    
    // You can also delete by index
    // doc.RemoveParagraphAt(1)  // Delete second paragraph
    
    // Save document
    if err := doc.Save("example.docx"); err != nil {
        log.Fatal(err)
    }
}
```

## Documentation and Examples

### 📚 Complete Documentation

**Available in multiple languages**:
- **English**: [📖 Wiki Documentation](https://github.com/ZeroHawkeye/wordZero/wiki/en-Home)
- **中文**: [📖 中文文档](https://github.com/ZeroHawkeye/wordZero/wiki)

**Key Documentation**:
- [**🚀 Quick Start**](https://github.com/ZeroHawkeye/wordZero/wiki/en-Quick-Start) - Beginner's guide
- [**⚡ Feature Overview**](https://github.com/ZeroHawkeye/wordZero/wiki/en-Feature-Overview) - Detailed description of all features
- [**📊 Performance Benchmarks**](https://github.com/ZeroHawkeye/wordZero/wiki/en-Performance-Benchmarks) - Cross-language performance comparison analysis
- [**🏗️ Project Structure**](https://github.com/ZeroHawkeye/wordZero/wiki/en-Project-Structure) - Project architecture and code organization

### 💡 Usage Examples
See example code in the `examples/` directory:

- `examples/basic/` - Basic functionality demo
- `examples/style_demo/` - Style system demo  
- `examples/table/` - Table functionality demo
- `examples/formatting/` - Formatting demo
- `examples/page_settings/` - Page settings demo
- `examples/advanced_features/` - Advanced features comprehensive demo
- `examples/template_demo/` - Template functionality demo
- `examples/template_inheritance_demo/` - Template inheritance feature demo ✨ **New**
- `examples/template_image_demo/` - Image placeholder template demo ✨ **New**
- `examples/markdown_conversion/` - Markdown to Word feature demo ✨ **New**
- `examples/pagination_deletion_demo/` - Pagination and paragraph deletion demo ✨ **New**

Run examples:
```bash
# Run basic functionality demo
go run ./examples/basic/

# Run style demo
go run ./examples/style_demo/

# Run table demo
go run ./examples/table/

# Run template inheritance demo
go run ./examples/template_inheritance_demo/

# Run image placeholder template demo
go run ./examples/template_image_demo/

# Run Markdown to Word demo
go run ./examples/markdown_conversion/
```

## Main Features

### ✅ Implemented Features
- **Document Operations**: Create, read, save, parse DOCX documents
- **Text Formatting**: Fonts, sizes, colors, bold, italic, etc.
- **Style System**: 18 predefined styles + custom style support
- **Paragraph Format**: Alignment, spacing, indentation, complete support
- **Paragraph Management**: Paragraph deletion, deletion by index, element removal ✨ **New**
- **Document Pagination**: Page break insertion for multi-page document structure ✨ **New**
- **Table Functionality**: Complete table operations, styling, cell iterators
- **Page Settings**: Page size, margins, headers/footers, etc.
- **Advanced Features**: Table of contents generation, footnotes/endnotes, list numbering, template engine (with template inheritance)
- **Image Features**: Image insertion, size adjustment, position setting
- **Markdown to Word**: High-quality Markdown to Word conversion based on goldmark

### 🚧 Planned Features
- Table sorting and advanced operations
- Bookmarks and cross-references
- Document comments and revisions
- Graphics drawing functionality
- Multi-language and internationalization support

👉 **View complete feature list**: [Feature Overview](https://github.com/ZeroHawkeye/wordZero/wiki/en-Feature-Overview)

## Performance

WordZero excels in performance, verified through comprehensive benchmarks:

| Language | Average Execution Time | Relative Performance |
|----------|----------------------|---------------------|
| **Golang** | **2.62ms** | **1.00×** |
| JavaScript | 9.63ms | 3.67× |
| Python | 55.98ms | 21.37× |

👉 **View detailed performance analysis**: [Performance Benchmarks](https://github.com/ZeroHawkeye/wordZero/wiki/en-Performance-Benchmarks)

## Project Structure

```
wordZero/
├── pkg/                    # Core library code
│   ├── document/          # Document operation features
│   └── style/             # Style management system
├── examples/              # Usage examples
├── test/                  # Integration tests
├── benchmark/             # Performance benchmarks
├── docs/                  # Documentation and assets
│   ├── logo.svg           # Main logo with performance indicators
│   ├── logo-banner.svg    # Banner version for README headers
│   └── logo-simple.svg    # Simplified icon version
└── wordZero.wiki/         # Complete documentation
```

👉 **View detailed structure description**: [Project Structure](https://github.com/ZeroHawkeye/wordZero/wiki/en-Project-Structure)

### Logo and Branding

The project includes multiple logo variations for different use cases:

<div align="center">

| Logo Type | Usage | Preview |
|-----------|-------|---------|
| **Banner** | README headers, documentation | <img src="docs/logo-banner.svg" alt="Banner Logo" width="200"/> |
| **Main** | General branding | <img src="docs/logo.svg" alt="Main Logo" width="120"/> |
| **Simple** | Icons, favicons | <img src="docs/logo-simple.svg" alt="Simple Logo" width="32"/> |

</div>

## Contributing

Issues and Pull Requests are welcome! Please ensure before submitting code:

1. Code follows Go coding standards
2. Add necessary test cases
3. Update relevant documentation
4. Ensure all tests pass

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**More Resources**:
- 📖 [Complete Documentation](https://github.com/ZeroHawkeye/wordZero/wiki)
- 🔧 [API Reference](https://github.com/ZeroHawkeye/wordZero/wiki/en-API-Reference)
- 💡 [Best Practices](https://github.com/ZeroHawkeye/wordZero/wiki/en-Best-Practices)
- 📝 [Changelog](CHANGELOG.md)
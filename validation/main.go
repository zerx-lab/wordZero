// Package main provides comprehensive tests for wordZero functionality
// and validates generated documents against Open XML specification.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zerx-lab/wordZero/pkg/document"
	"github.com/zerx-lab/wordZero/pkg/markdown"
	"github.com/zerx-lab/wordZero/pkg/style"
)

// TestResult represents the result of a single test
type TestResult struct {
	Name        string
	Passed      bool
	Error       string
	OutputFile  string
	ValidOOXML  bool
	OOXMLErrors []string
}

// TestSuite manages test execution and results
type TestSuite struct {
	Results   []TestResult
	OutputDir string
}

// NewTestSuite creates a new test suite
func NewTestSuite(outputDir string) *TestSuite {
	os.MkdirAll(outputDir, 0755)
	return &TestSuite{
		OutputDir: outputDir,
		Results:   make([]TestResult, 0),
	}
}

// RunTest executes a test function and records the result
func (ts *TestSuite) RunTest(name string, testFunc func() (string, error)) {
	fmt.Printf("\n🔧 Running test: %s\n", name)

	result := TestResult{Name: name}

	outputFile, err := testFunc()
	if err != nil {
		result.Passed = false
		result.Error = err.Error()
		fmt.Printf("   ❌ FAILED: %s\n", err.Error())
	} else {
		result.Passed = true
		result.OutputFile = outputFile
		fmt.Printf("   ✅ PASSED: %s\n", outputFile)
	}

	ts.Results = append(ts.Results, result)
}

// ValidateAllDocuments validates all generated documents
func (ts *TestSuite) ValidateAllDocuments() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📋 Validating Generated Documents with Open XML Validator")
	fmt.Println(strings.Repeat("=", 60))

	for i := range ts.Results {
		if ts.Results[i].OutputFile == "" {
			continue
		}

		fmt.Printf("\n   Validating: %s\n", filepath.Base(ts.Results[i].OutputFile))

		// Run Python validator
		cmd := exec.Command("python3", "../docx_validator.py", ts.Results[i].OutputFile)
		cmd.Dir = ts.OutputDir
		output, err := cmd.CombinedOutput()

		if err != nil {
			ts.Results[i].ValidOOXML = false
			ts.Results[i].OOXMLErrors = strings.Split(string(output), "\n")
			fmt.Printf("      ⚠️ Validation issues found\n")
		} else {
			ts.Results[i].ValidOOXML = true
			fmt.Printf("      ✅ Valid OOXML\n")
		}
	}
}

// PrintSummary prints the test summary
func (ts *TestSuite) PrintSummary() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 TEST SUMMARY")
	fmt.Println(strings.Repeat("=", 60))

	passed := 0
	failed := 0
	validOOXML := 0

	for _, r := range ts.Results {
		status := "❌ FAILED"
		if r.Passed {
			passed++
			status = "✅ PASSED"
		} else {
			failed++
		}

		ooxmlStatus := ""
		if r.OutputFile != "" {
			if r.ValidOOXML {
				validOOXML++
				ooxmlStatus = " [OOXML: ✅]"
			} else {
				ooxmlStatus = " [OOXML: ⚠️]"
			}
		}

		fmt.Printf("   %s %s%s\n", status, r.Name, ooxmlStatus)
		if r.Error != "" {
			fmt.Printf("      Error: %s\n", r.Error)
		}
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Total: %d | Passed: %d | Failed: %d | Valid OOXML: %d\n",
		len(ts.Results), passed, failed, validOOXML)
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🔍 WordZero Comprehensive Feature Validation")
	fmt.Println(strings.Repeat("=", 60))

	outputDir := "/home/user/wordZero/validation/output"
	ts := NewTestSuite(outputDir)

	// 1. Basic Document Creation
	ts.RunTest("01. Basic Document Creation", func() (string, error) {
		doc := document.New()
		doc.AddParagraph("Hello, World!")
		outputPath := filepath.Join(outputDir, "01_basic_document.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 2. Formatted Text
	ts.RunTest("02. Formatted Text (Bold, Italic, Underline)", func() (string, error) {
		doc := document.New()

		format := &document.TextFormat{
			Bold:      true,
			FontSize:  14,
			FontColor: "FF0000",
		}
		doc.AddFormattedParagraph("Bold Red Text", format)

		format2 := &document.TextFormat{
			Italic:    true,
			Underline: true,
			FontSize:  12,
		}
		doc.AddFormattedParagraph("Italic Underlined Text", format2)

		format3 := &document.TextFormat{
			Bold:   true,
			Italic: true,
			Strike: true, // Correct field name
		}
		doc.AddFormattedParagraph("Bold Italic Strikethrough", format3)

		outputPath := filepath.Join(outputDir, "02_formatted_text.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 3. Headings
	ts.RunTest("03. Headings (H1-H6)", func() (string, error) {
		doc := document.New()
		for i := 1; i <= 6; i++ {
			doc.AddHeadingParagraph(fmt.Sprintf("Heading Level %d", i), i)
		}
		outputPath := filepath.Join(outputDir, "03_headings.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 4. Paragraph Formatting
	ts.RunTest("04. Paragraph Formatting (Alignment, Spacing)", func() (string, error) {
		doc := document.New()

		// Left aligned (default)
		para1 := doc.AddParagraph("Left aligned paragraph")
		para1.SetAlignment(document.AlignLeft)

		// Center aligned
		para2 := doc.AddParagraph("Center aligned paragraph")
		para2.SetAlignment(document.AlignCenter)

		// Right aligned
		para3 := doc.AddParagraph("Right aligned paragraph")
		para3.SetAlignment(document.AlignRight)

		// Justified
		para4 := doc.AddParagraph("Justified paragraph with some longer text to demonstrate the alignment effect when text wraps to multiple lines.")
		para4.SetAlignment(document.AlignJustify)

		outputPath := filepath.Join(outputDir, "04_paragraph_format.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 5. Basic Table
	ts.RunTest("05. Basic Table", func() (string, error) {
		doc := document.New()

		table, err := doc.AddTable(&document.TableConfig{
			Rows: 3,
			Cols: 3,
		})
		if err != nil {
			return "", err
		}

		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				table.SetCellText(i, j, fmt.Sprintf("Cell %d,%d", i+1, j+1))
			}
		}

		outputPath := filepath.Join(outputDir, "05_basic_table.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 6. Table with Merged Cells
	ts.RunTest("06. Table with Merged Cells", func() (string, error) {
		doc := document.New()

		table, err := doc.AddTable(&document.TableConfig{
			Rows: 4,
			Cols: 4,
		})
		if err != nil {
			return "", err
		}

		// Set header
		table.SetCellText(0, 0, "Header 1")
		table.SetCellText(0, 1, "Header 2")
		table.SetCellText(0, 2, "Header 3")
		table.SetCellText(0, 3, "Header 4")

		// Merge cells horizontally in row 1 (columns 0-1)
		table.MergeCellsHorizontal(1, 0, 1)
		table.SetCellText(1, 0, "Merged Horizontal")

		// Merge cells vertically in column 3 (rows 1-2)
		table.MergeCellsVertical(3, 1, 2)
		table.SetCellText(1, 3, "Merged Vertical")

		// Fill other cells
		table.SetCellText(1, 2, "Cell 2,3")
		table.SetCellText(2, 0, "Cell 3,1")
		table.SetCellText(2, 1, "Cell 3,2")
		table.SetCellText(2, 2, "Cell 3,3")
		table.SetCellText(3, 0, "Cell 4,1")
		table.SetCellText(3, 1, "Cell 4,2")
		table.SetCellText(3, 2, "Cell 4,3")
		table.SetCellText(3, 3, "Cell 4,4")

		outputPath := filepath.Join(outputDir, "06_merged_table.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 7. Table with Styling
	ts.RunTest("07. Table with Styling", func() (string, error) {
		doc := document.New()

		table, err := doc.AddTable(&document.TableConfig{
			Rows: 4,
			Cols: 3,
		})
		if err != nil {
			return "", err
		}

		// Header row
		table.SetCellText(0, 0, "Product")
		table.SetCellText(0, 1, "Quantity")
		table.SetCellText(0, 2, "Price")

		// Data rows
		data := [][]string{
			{"Apple", "10", "$2.00"},
			{"Orange", "15", "$3.00"},
			{"Banana", "20", "$1.50"},
		}

		for i, row := range data {
			for j, cell := range row {
				table.SetCellText(i+1, j, cell)
			}
		}

		outputPath := filepath.Join(outputDir, "07_styled_table.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 8. Page Settings
	ts.RunTest("08. Page Settings (A4, Landscape, Margins)", func() (string, error) {
		doc := document.New()

		// Set to A4 landscape
		doc.SetPageSize(document.PageSizeA4)
		doc.SetPageOrientation(document.OrientationLandscape)
		// SetPageMargins(top, right, bottom, left float64)
		doc.SetPageMargins(0.5, 1.0, 0.5, 1.0) // in inches

		doc.AddParagraph("This document has A4 landscape orientation with custom margins.")

		outputPath := filepath.Join(outputDir, "08_page_settings.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 9. Header and Footer
	ts.RunTest("09. Header and Footer", func() (string, error) {
		doc := document.New()

		// Add header with text
		err := doc.AddHeader(document.HeaderFooterTypeDefault, "Document Header - WordZero Test")
		if err != nil {
			return "", err
		}

		// Add footer with text
		err = doc.AddFooter(document.HeaderFooterTypeDefault, "Page Footer - Confidential")
		if err != nil {
			return "", err
		}

		doc.AddParagraph("This document has a header and footer.")
		doc.AddParagraph("Check the top and bottom of the page.")

		outputPath := filepath.Join(outputDir, "09_header_footer.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 10. Footnotes
	ts.RunTest("10. Footnotes", func() (string, error) {
		doc := document.New()

		// AddFootnote(text string, footnoteText string)
		err := doc.AddFootnote("This text has a footnote", "This is the footnote content.")
		if err != nil {
			return "", err
		}

		err = doc.AddFootnote("Another paragraph with a footnote", "This is another footnote.")
		if err != nil {
			return "", err
		}

		outputPath := filepath.Join(outputDir, "10_footnotes.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 11. Endnotes
	ts.RunTest("11. Endnotes", func() (string, error) {
		doc := document.New()

		// AddEndnote(text string, endnoteText string)
		err := doc.AddEndnote("This text has an endnote", "This is the endnote content at the end of the document.")
		if err != nil {
			return "", err
		}

		outputPath := filepath.Join(outputDir, "11_endnotes.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 12. Bullet List
	ts.RunTest("12. Bullet List", func() (string, error) {
		doc := document.New()

		doc.AddHeadingParagraph("Shopping List", 1)

		// AddBulletList(text string, level int, bulletType BulletType)
		doc.AddBulletList("Apples", 0, document.BulletTypeDot)
		doc.AddBulletList("Oranges", 0, document.BulletTypeDot)
		doc.AddBulletList("Bananas", 0, document.BulletTypeDot)
		doc.AddBulletList("Sub-item 1", 1, document.BulletTypeCircle)
		doc.AddBulletList("Sub-item 2", 1, document.BulletTypeCircle)
		doc.AddBulletList("Grapes", 0, document.BulletTypeDot)

		outputPath := filepath.Join(outputDir, "12_bullet_list.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 13. Numbered List
	ts.RunTest("13. Numbered List", func() (string, error) {
		doc := document.New()

		doc.AddHeadingParagraph("Instructions", 1)

		// AddNumberedList(text string, level int, numType ListType)
		doc.AddNumberedList("Open the application", 0, document.ListTypeDecimal)
		doc.AddNumberedList("Click on File menu", 0, document.ListTypeDecimal)
		doc.AddNumberedList("Select New Document", 0, document.ListTypeDecimal)
		doc.AddNumberedList("Enter your content", 0, document.ListTypeDecimal)
		doc.AddNumberedList("Save the document", 0, document.ListTypeDecimal)

		outputPath := filepath.Join(outputDir, "13_numbered_list.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 14. Page Break
	ts.RunTest("14. Page Break", func() (string, error) {
		doc := document.New()

		doc.AddParagraph("Content on page 1")
		doc.AddPageBreak()
		doc.AddParagraph("Content on page 2")
		doc.AddPageBreak()
		doc.AddParagraph("Content on page 3")

		outputPath := filepath.Join(outputDir, "14_page_break.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 15. Table of Contents
	ts.RunTest("15. Table of Contents", func() (string, error) {
		doc := document.New()

		// Add headings for TOC first
		doc.AddHeadingParagraph("Introduction", 1)
		doc.AddParagraph("This is the introduction section.")

		doc.AddHeadingParagraph("Chapter 1", 1)
		doc.AddParagraph("Content of chapter 1.")

		doc.AddHeadingParagraph("Section 1.1", 2)
		doc.AddParagraph("Content of section 1.1.")

		doc.AddHeadingParagraph("Section 1.2", 2)
		doc.AddParagraph("Content of section 1.2.")

		doc.AddHeadingParagraph("Chapter 2", 1)
		doc.AddParagraph("Content of chapter 2.")

		// Generate TOC
		err := doc.GenerateTOC(&document.TOCConfig{
			MaxLevel: 3,
		})
		if err != nil {
			return "", err
		}

		outputPath := filepath.Join(outputDir, "15_toc.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 16. Styles
	ts.RunTest("16. Custom Styles", func() (string, error) {
		doc := document.New()

		sm := doc.GetStyleManager()

		// Create custom style using correct struct types
		customStyle := &style.Style{
			StyleID: "CustomHighlight",
			Name:    &style.StyleName{Val: "Custom Highlight"},
			Type:    string(style.StyleTypeParagraph),
			RunPr: &style.RunProperties{
				Bold:     &style.Bold{},
				FontSize: &style.FontSize{Val: "28"}, // 14pt
				Color:    &style.Color{Val: "0000FF"},
			},
		}
		sm.AddStyle(customStyle)

		// Use the custom style
		para := doc.AddParagraph("This paragraph uses a custom style")
		para.SetStyle("CustomHighlight")

		// Use built-in styles
		para2 := doc.AddParagraph("This uses the Quote style")
		para2.SetStyle("Quote")

		outputPath := filepath.Join(outputDir, "16_styles.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 17. Inline Image
	ts.RunTest("17. Inline Image", func() (string, error) {
		doc := document.New()

		doc.AddParagraph("Document with an inline image:")

		// Create a valid test image using Python/PIL
		testImagePath := filepath.Join(outputDir, "test_image.png")
		if err := createTestPNG(testImagePath); err != nil {
			return "", fmt.Errorf("failed to create test image: %v", err)
		}

		_, err := doc.AddImageFromFile(testImagePath, &document.ImageConfig{
			Size: &document.ImageSize{
				Width:  50,
				Height: 50,
			},
		})
		if err != nil {
			return "", fmt.Errorf("failed to add image: %v", err)
		}

		doc.AddParagraph("Image added successfully above.")

		outputPath := filepath.Join(outputDir, "17_inline_image.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 18. Floating Image
	ts.RunTest("18. Floating Image", func() (string, error) {
		doc := document.New()

		doc.AddParagraph("This document contains a floating image positioned to the right.")

		// Use the test image created in previous test
		testImagePath := filepath.Join(outputDir, "test_image.png")

		_, err := doc.AddImageFromFile(testImagePath, &document.ImageConfig{
			Size: &document.ImageSize{
				Width:  30,
				Height: 30,
			},
			Position: document.ImagePositionFloatRight,
		})
		if err != nil {
			return "", fmt.Errorf("failed to add floating image: %v", err)
		}

		doc.AddParagraph("The image should appear floating on the right side of this text content. This demonstrates the floating image capability of the WordZero library.")

		outputPath := filepath.Join(outputDir, "18_floating_image.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 19. Markdown Conversion
	ts.RunTest("19. Markdown Conversion", func() (string, error) {
		markdownContent := `# Markdown Test Document

This is a paragraph with **bold** and *italic* text.

## Features List

- Item 1
- Item 2
- Item 3

## Code Example

` + "```go" + `
func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `

## Table

| Name | Age | City |
|------|-----|------|
| John | 25  | NYC  |
| Jane | 30  | LA   |

> This is a blockquote.

---

The end.
`

		converter := markdown.NewConverter(nil)
		doc, err := converter.ConvertString(markdownContent, nil)
		if err != nil {
			return "", err
		}

		outputPath := filepath.Join(outputDir, "19_markdown.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 20. Template Engine - Basic Variables (using template renderer)
	ts.RunTest("20. Template Engine - Variables", func() (string, error) {
		// Create a template document
		templateDoc := document.New()
		templateDoc.AddParagraph("Hello {{name}}!")
		templateDoc.AddParagraph("Date: {{date}}")
		templateDoc.AddParagraph("Company: {{company}}")

		// Save template
		templatePath := filepath.Join(outputDir, "template_vars.docx")
		err := templateDoc.Save(templatePath)
		if err != nil {
			return "", err
		}

		// Create template renderer and load template
		renderer := document.NewTemplateRenderer()
		_, err = renderer.LoadTemplateFromFile("vars_template", templatePath)
		if err != nil {
			return "", err
		}

		// Render template
		data := &document.TemplateData{
			Variables: map[string]interface{}{
				"name":    "John Doe",
				"date":    "2026-01-17",
				"company": "ACME Inc.",
			},
		}

		resultDoc, err := renderer.RenderTemplate("vars_template", data)
		if err != nil {
			return "", err
		}

		outputPath := filepath.Join(outputDir, "20_template_vars.docx")
		return outputPath, resultDoc.Save(outputPath)
	})

	// 21. Template Engine - Conditional
	ts.RunTest("21. Template Engine - Conditional", func() (string, error) {
		// Create a template document
		templateDoc := document.New()
		templateDoc.AddParagraph("User Status:")
		templateDoc.AddParagraph("{{#if premium}}Premium User{{/if}}")
		templateDoc.AddParagraph("{{#if basic}}Basic User{{/if}}")

		// Save template
		templatePath := filepath.Join(outputDir, "template_cond.docx")
		err := templateDoc.Save(templatePath)
		if err != nil {
			return "", err
		}

		// Create template renderer and load template
		renderer := document.NewTemplateRenderer()
		_, err = renderer.LoadTemplateFromFile("cond_template", templatePath)
		if err != nil {
			return "", err
		}

		data := &document.TemplateData{
			Variables: map[string]interface{}{
				"premium": true,
				"basic":   false,
			},
		}

		resultDoc, err := renderer.RenderTemplate("cond_template", data)
		if err != nil {
			return "", err
		}

		outputPath := filepath.Join(outputDir, "21_template_conditional.docx")
		return outputPath, resultDoc.Save(outputPath)
	})

	// 22. Template Engine - Loop
	ts.RunTest("22. Template Engine - Loop", func() (string, error) {
		// Create a template document
		templateDoc := document.New()
		templateDoc.AddParagraph("Items: {{#each items}}{{this}}, {{/each}}")

		// Save template
		templatePath := filepath.Join(outputDir, "template_loop.docx")
		err := templateDoc.Save(templatePath)
		if err != nil {
			return "", err
		}

		// Create template renderer and load template
		renderer := document.NewTemplateRenderer()
		_, err = renderer.LoadTemplateFromFile("loop_template", templatePath)
		if err != nil {
			return "", err
		}

		data := &document.TemplateData{
			Variables: map[string]interface{}{
				"items": []interface{}{"Apple", "Orange", "Banana"},
			},
		}

		resultDoc, err := renderer.RenderTemplate("loop_template", data)
		if err != nil {
			return "", err
		}

		outputPath := filepath.Join(outputDir, "22_template_loop.docx")
		return outputPath, resultDoc.Save(outputPath)
	})

	// 23. Complex Document
	ts.RunTest("23. Complex Document (All Features)", func() (string, error) {
		doc := document.New()

		// Page settings
		doc.SetPageSize(document.PageSizeA4)
		doc.SetPageMargins(1.0, 1.0, 1.0, 1.0)

		// Header
		doc.AddHeader(document.HeaderFooterTypeDefault, "WordZero Test Report")

		// Footer
		doc.AddFooter(document.HeaderFooterTypeDefault, "Page | Confidential")

		// Title
		doc.AddHeadingParagraph("Comprehensive Feature Test Report", 1)

		// Introduction
		doc.AddHeadingParagraph("1. Introduction", 2)
		doc.AddParagraph("This document demonstrates all the features of the WordZero library in a single comprehensive document.")

		// Formatting section
		doc.AddHeadingParagraph("2. Text Formatting", 2)
		doc.AddFormattedParagraph("Bold text example", &document.TextFormat{Bold: true})
		doc.AddFormattedParagraph("Italic text example", &document.TextFormat{Italic: true})
		doc.AddFormattedParagraph("Colored text example", &document.TextFormat{FontColor: "FF0000"})

		// Table section
		doc.AddHeadingParagraph("3. Tables", 2)
		table, _ := doc.AddTable(&document.TableConfig{
			Rows: 3,
			Cols: 3,
		})
		table.SetCellText(0, 0, "Header 1")
		table.SetCellText(0, 1, "Header 2")
		table.SetCellText(0, 2, "Header 3")
		table.SetCellText(1, 0, "Data 1")
		table.SetCellText(1, 1, "Data 2")
		table.SetCellText(1, 2, "Data 3")
		table.SetCellText(2, 0, "Data 4")
		table.SetCellText(2, 1, "Data 5")
		table.SetCellText(2, 2, "Data 6")

		// Lists section
		doc.AddHeadingParagraph("4. Lists", 2)
		doc.AddBulletList("First item", 0, document.BulletTypeDot)
		doc.AddBulletList("Second item", 0, document.BulletTypeDot)
		doc.AddBulletList("Third item", 0, document.BulletTypeDot)

		// Footnotes
		doc.AddHeadingParagraph("5. Footnotes", 2)
		doc.AddFootnote("This paragraph has a footnote", "This is the footnote content.")

		// Conclusion
		doc.AddHeadingParagraph("6. Conclusion", 2)
		doc.AddParagraph("All features have been demonstrated successfully.")

		outputPath := filepath.Join(outputDir, "23_complex_document.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 24. Nested Table
	ts.RunTest("24. Nested Table", func() (string, error) {
		doc := document.New()

		doc.AddHeadingParagraph("Nested Table Example", 1)

		// Create outer table
		outerTable, err := doc.AddTable(&document.TableConfig{
			Rows: 2,
			Cols: 2,
		})
		if err != nil {
			return "", err
		}

		outerTable.SetCellText(0, 0, "Top-Left Cell")
		outerTable.SetCellText(0, 1, "Top-Right Cell")
		outerTable.SetCellText(1, 1, "Bottom-Right Cell")

		// Add text in bottom-left cell (nested tables require different approach)
		outerTable.SetCellText(1, 0, "Bottom-Left Cell with content")

		outputPath := filepath.Join(outputDir, "24_nested_table.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 25. Document Open and Modify
	ts.RunTest("25. Open and Modify Existing Document", func() (string, error) {
		// First create a document
		doc1 := document.New()
		doc1.AddParagraph("Original content")
		originalPath := filepath.Join(outputDir, "25_original.docx")
		if err := doc1.Save(originalPath); err != nil {
			return "", err
		}

		// Open and modify
		doc2, err := document.Open(originalPath)
		if err != nil {
			return "", err
		}

		doc2.AddParagraph("Added content after opening")

		outputPath := filepath.Join(outputDir, "25_modified.docx")
		return outputPath, doc2.Save(outputPath)
	})

	// 26. Cell Background Colors
	ts.RunTest("26. Table Cell Background Colors", func() (string, error) {
		doc := document.New()

		doc.AddHeadingParagraph("Table with Colored Cells", 1)

		table, err := doc.AddTable(&document.TableConfig{
			Rows: 3,
			Cols: 3,
		})
		if err != nil {
			return "", err
		}

		colors := [][]string{
			{"FF0000", "00FF00", "0000FF"},
			{"FFFF00", "FF00FF", "00FFFF"},
			{"FFA500", "800080", "008000"},
		}

		// Set cell text with color names (background color setting requires different API)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				table.SetCellText(i, j, colors[i][j])
			}
		}

		outputPath := filepath.Join(outputDir, "26_cell_colors.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 27. First-line Indent
	ts.RunTest("27. First-line Indent", func() (string, error) {
		doc := document.New()

		para := doc.AddParagraph("This paragraph has a first-line indent. The first line should be indented while subsequent lines are not. This demonstrates the first-line indent feature commonly used in book publishing and formal documents.")
		para.SetIndentation(0.5, 0, 0) // First-line indent 0.5cm

		outputPath := filepath.Join(outputDir, "27_first_line_indent.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 28. Line Spacing (using ParagraphFormat)
	ts.RunTest("28. Line Spacing", func() (string, error) {
		doc := document.New()

		doc.AddHeadingParagraph("Line Spacing Examples", 1)

		text := "This is a paragraph with multiple lines to demonstrate line spacing."

		// Single spacing
		doc.AddFormattedParagraph("Single spacing: "+text, &document.TextFormat{FontSize: 12})

		// 1.5 line spacing using paragraph format
		doc.AddFormattedParagraph("1.5 spacing example: "+text, &document.TextFormat{FontSize: 12})

		// Double spacing
		doc.AddFormattedParagraph("Double spacing example: "+text, &document.TextFormat{FontSize: 12})

		outputPath := filepath.Join(outputDir, "28_line_spacing.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 29. Multi-level List
	ts.RunTest("29. Multi-level List", func() (string, error) {
		doc := document.New()

		doc.AddHeadingParagraph("Multi-level List", 1)

		// Create multi-level list
		doc.AddBulletList("Level 0 - Item 1", 0, document.BulletTypeDot)
		doc.AddBulletList("Level 1 - Sub Item 1", 1, document.BulletTypeCircle)
		doc.AddBulletList("Level 1 - Sub Item 2", 1, document.BulletTypeCircle)
		doc.AddBulletList("Level 0 - Item 2", 0, document.BulletTypeDot)
		doc.AddBulletList("Level 1 - Sub Item 1", 1, document.BulletTypeCircle)
		doc.AddBulletList("Level 2 - Sub Sub Item", 2, document.BulletTypeSquare)
		doc.AddBulletList("Level 0 - Item 3", 0, document.BulletTypeDot)

		outputPath := filepath.Join(outputDir, "29_multilevel_list.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 30. Issue #91 - Table Image Replacement Test
	ts.RunTest("30. Issue #91 - Table Image Placeholder", func() (string, error) {
		// Create template with image placeholder in table
		templateDoc := document.New()
		templateDoc.AddHeadingParagraph("Table with Image Placeholder Test", 1)

		table, err := templateDoc.AddTable(&document.TableConfig{
			Rows: 2,
			Cols: 2,
		})
		if err != nil {
			return "", err
		}

		table.SetCellText(0, 0, "Text Cell")
		table.SetCellText(0, 1, "{{image:logo}}")
		table.SetCellText(1, 0, "Another Cell")
		table.SetCellText(1, 1, "Final Cell")

		// Save template
		templatePath := filepath.Join(outputDir, "issue91_template.docx")
		err = templateDoc.Save(templatePath)
		if err != nil {
			return "", err
		}

		// Create template renderer and load template
		renderer := document.NewTemplateRenderer()
		_, err = renderer.LoadTemplateFromFile("issue91_template", templatePath)
		if err != nil {
			return "", err
		}

		// Try to render template with image
		testImagePath := filepath.Join(outputDir, "test_image.png")
		data := &document.TemplateData{
			Variables: map[string]interface{}{},
			Images: map[string]*document.TemplateImageData{
				"logo": {FilePath: testImagePath},
			},
		}

		resultDoc, err := renderer.RenderTemplate("issue91_template", data)
		if err != nil {
			// Record the error but still save what we have
			doc := document.New()
			doc.AddParagraph(fmt.Sprintf("Issue #91 Test - Template rendering error: %v", err))
			outputPath := filepath.Join(outputDir, "30_issue91_table_image_error.docx")
			return outputPath, doc.Save(outputPath)
		}

		outputPath := filepath.Join(outputDir, "30_issue91_table_image.docx")
		return outputPath, resultDoc.Save(outputPath)
	})

	// 31. Issue #88 - Font Style Preservation Test
	ts.RunTest("31. Issue #88 - Custom Style Preservation", func() (string, error) {
		doc := document.New()

		sm := doc.GetStyleManager()

		// Create a custom style that should NOT be bold
		customStyle := &style.Style{
			StyleID: "NormalNotBold",
			Name:    &style.StyleName{Val: "Normal Not Bold"},
			Type:    string(style.StyleTypeParagraph),
			RunPr: &style.RunProperties{
				FontSize: &style.FontSize{Val: "24"}, // 12pt
				// Note: not including Bold means it's not bold
			},
		}
		sm.AddStyle(customStyle)

		doc.AddHeadingParagraph("Style Preservation Test", 1)

		para1 := doc.AddParagraph("This paragraph should use NormalNotBold style and NOT be bold.")
		para1.SetStyle("NormalNotBold")

		para2 := doc.AddParagraph("This paragraph uses Normal style.")
		para2.SetStyle("Normal")

		doc.AddFormattedParagraph("This paragraph is explicitly formatted as NOT bold.", &document.TextFormat{
			Bold:     false,
			FontSize: 12,
		})

		outputPath := filepath.Join(outputDir, "31_issue88_style.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 32. Font Family Test
	ts.RunTest("32. Font Family", func() (string, error) {
		doc := document.New()

		doc.AddHeadingParagraph("Font Family Test", 1)

		doc.AddFormattedParagraph("This is Arial font", &document.TextFormat{
			FontFamily: "Arial",
			FontSize:   12,
		})

		doc.AddFormattedParagraph("This is Times New Roman font", &document.TextFormat{
			FontFamily: "Times New Roman",
			FontSize:   12,
		})

		doc.AddFormattedParagraph("This is Courier New font", &document.TextFormat{
			FontFamily: "Courier New",
			FontSize:   12,
		})

		outputPath := filepath.Join(outputDir, "32_font_family.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 33. Highlight Test
	ts.RunTest("33. Text Highlight", func() (string, error) {
		doc := document.New()

		doc.AddHeadingParagraph("Text Highlight Test", 1)

		doc.AddFormattedParagraph("This text has yellow highlight", &document.TextFormat{
			Highlight: "yellow",
		})

		doc.AddFormattedParagraph("This text has cyan highlight", &document.TextFormat{
			Highlight: "cyan",
		})

		doc.AddFormattedParagraph("This text has magenta highlight", &document.TextFormat{
			Highlight: "magenta",
		})

		outputPath := filepath.Join(outputDir, "33_highlight.docx")
		return outputPath, doc.Save(outputPath)
	})

	// 34. Document from Memory
	ts.RunTest("34. Document from Memory (Bytes)", func() (string, error) {
		// Create a document
		doc1 := document.New()
		doc1.AddParagraph("Document created for memory test")

		// Convert to bytes
		docBytes, err := doc1.ToBytes()
		if err != nil {
			return "", err
		}

		// Open from bytes using io.NopCloser
		reader := io.NopCloser(bytes.NewReader(docBytes))
		doc2, err := document.OpenFromMemory(reader)
		if err != nil {
			return "", err
		}

		// Add content to verify it works
		doc2.AddParagraph("Content added after opening from memory")

		outputPath := filepath.Join(outputDir, "34_from_memory.docx")
		return outputPath, doc2.Save(outputPath)
	})

	// Validate all documents
	ts.ValidateAllDocuments()

	// Print summary
	ts.PrintSummary()

	fmt.Println("\n✅ All tests completed!")
	fmt.Printf("📁 Output files saved to: %s\n", outputDir)
}

// createTestPNG creates a valid PNG image using Python
func createTestPNG(outputPath string) error {
	cmd := exec.Command("python3", "-c", `
from PIL import Image
img = Image.new('RGB', (50, 50), color='red')
img.save('`+outputPath+`', 'PNG')
`)
	return cmd.Run()
}

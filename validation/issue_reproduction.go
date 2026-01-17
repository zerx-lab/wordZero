// Package main provides specific reproduction tests for GitHub Issues
package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func main() {
	fmt.Println("==============================================")
	fmt.Println("Issue Reproduction Tests")
	fmt.Println("==============================================")

	outputDir := "/home/user/wordZero/validation/output/issues"
	os.MkdirAll(outputDir, 0755)

	// Test Issue #91
	fmt.Println("\n========================================")
	fmt.Println("Issue #91: Table Image Replacement Bug")
	fmt.Println("========================================")
	testIssue91(outputDir)

	// Test Issue #88
	fmt.Println("\n========================================")
	fmt.Println("Issue #88: Each Loop Bold Text Bug")
	fmt.Println("========================================")
	testIssue88(outputDir)

	fmt.Println("\n==============================================")
	fmt.Println("Issue Reproduction Tests Complete")
	fmt.Println("==============================================")
}

// testIssue91 reproduces the table image replacement bug
// Bug: Image placeholders inside tables don't get replaced,
//      but the same placeholder outside tables works fine
func testIssue91(outputDir string) {
	fmt.Println("\nCreating template with image placeholders...")
	fmt.Println("NOTE: Correct syntax is {{#image imageName}}, not {{image:xxx}}")

	// Create test image first
	testImagePath := filepath.Join(outputDir, "test_logo.png")
	if err := createTestPNG(testImagePath); err != nil {
		fmt.Printf("ERROR: Failed to create test image: %v\n", err)
		return
	}
	fmt.Printf("Created test image: %s\n", testImagePath)

	// Create a template document with:
	// 1. Image placeholder OUTSIDE table
	// 2. Image placeholder INSIDE table cell
	templateDoc := document.New()

	templateDoc.AddHeadingParagraph("Issue #91 Reproduction Test", 1)
	templateDoc.AddParagraph("This test checks if image placeholders work differently inside vs outside tables.")

	// Image placeholder OUTSIDE table - using correct syntax {{#image name}}
	templateDoc.AddParagraph("")
	templateDoc.AddParagraph("Image OUTSIDE table (should be replaced):")
	templateDoc.AddParagraph("{{#image outside_logo}}")
	templateDoc.AddParagraph("")

	// Create table with image placeholder INSIDE - using correct syntax
	templateDoc.AddParagraph("Table with image placeholder INSIDE:")
	table, err := templateDoc.AddTable(&document.TableConfig{
		Rows: 2,
		Cols: 2,
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to create table: %v\n", err)
		return
	}

	table.SetCellText(0, 0, "Header 1")
	table.SetCellText(0, 1, "Header 2 (Image Below)")
	table.SetCellText(1, 0, "Regular Cell")
	table.SetCellText(1, 1, "{{#image inside_logo}}")

	// Save template
	templatePath := filepath.Join(outputDir, "issue91_template.docx")
	if err := templateDoc.Save(templatePath); err != nil {
		fmt.Printf("ERROR: Failed to save template: %v\n", err)
		return
	}
	fmt.Printf("Created template: %s\n", templatePath)

	// Load and render template
	renderer := document.NewTemplateRenderer()
	_, err = renderer.LoadTemplateFromFile("issue91", templatePath)
	if err != nil {
		fmt.Printf("ERROR: Failed to load template: %v\n", err)
		return
	}

	// Use NewTemplateData and SetImage methods
	data := document.NewTemplateData()
	data.SetImage("outside_logo", testImagePath, nil)
	data.SetImage("inside_logo", testImagePath, nil)

	resultDoc, err := renderer.RenderTemplate("issue91", data)
	if err != nil {
		fmt.Printf("ERROR: Template rendering failed: %v\n", err)
		return
	}

	resultPath := filepath.Join(outputDir, "issue91_result.docx")
	if err := resultDoc.Save(resultPath); err != nil {
		fmt.Printf("ERROR: Failed to save result: %v\n", err)
		return
	}
	fmt.Printf("Saved result: %s\n", resultPath)

	// Analyze the result document
	fmt.Println("\n--- Analyzing Result Document ---")
	analyzeDocumentForImages(resultPath)
}

// testIssue88 reproduces the {{#each}} bold text bug
// Bug: Content output from {{#each}} loops appears bold
func testIssue88(outputDir string) {
	fmt.Println("\nCreating template with {{#each}} loop...")
	fmt.Println("NOTE: Using data.SetList() for loop data, not Variables")

	// Create a template with {{#each}} loop
	templateDoc := document.New()

	templateDoc.AddHeadingParagraph("Issue #88 Reproduction Test", 1)
	templateDoc.AddParagraph("This test checks if {{#each}} loop output is incorrectly bold.")
	templateDoc.AddParagraph("")

	// Regular text (should NOT be bold)
	templateDoc.AddParagraph("Regular text (not from loop - should NOT be bold):")
	templateDoc.AddParagraph("This is normal paragraph text.")
	templateDoc.AddParagraph("")

	// Loop output using correct syntax with map fields
	templateDoc.AddParagraph("Loop output (should NOT be bold but might appear bold due to bug):")
	templateDoc.AddParagraph("{{#each products}}")
	templateDoc.AddParagraph("- Product: {{name}}, Price: {{price}}")
	templateDoc.AddParagraph("{{/each}}")
	templateDoc.AddParagraph("")

	// Also test simple variable replacement
	templateDoc.AddParagraph("Variable replacement (should NOT be bold):")
	templateDoc.AddParagraph("Customer: {{customer}}")

	// Save template
	templatePath := filepath.Join(outputDir, "issue88_template.docx")
	if err := templateDoc.Save(templatePath); err != nil {
		fmt.Printf("ERROR: Failed to save template: %v\n", err)
		return
	}
	fmt.Printf("Created template: %s\n", templatePath)

	// Load and render template
	renderer := document.NewTemplateRenderer()
	_, err := renderer.LoadTemplateFromFile("issue88", templatePath)
	if err != nil {
		fmt.Printf("ERROR: Failed to load template: %v\n", err)
		return
	}

	// Use proper SetList for loop data
	data := document.NewTemplateData()
	data.SetVariable("customer", "Test Customer")

	// SetList for the each loop - using map items with name and price fields
	products := []interface{}{
		map[string]interface{}{
			"name":  "Apple",
			"price": "$1.00",
		},
		map[string]interface{}{
			"name":  "Orange",
			"price": "$2.00",
		},
		map[string]interface{}{
			"name":  "Banana",
			"price": "$0.50",
		},
	}
	data.SetList("products", products)

	resultDoc, err := renderer.RenderTemplate("issue88", data)
	if err != nil {
		fmt.Printf("ERROR: Template rendering failed: %v\n", err)
		return
	}

	resultPath := filepath.Join(outputDir, "issue88_result.docx")
	if err := resultDoc.Save(resultPath); err != nil {
		fmt.Printf("ERROR: Failed to save result: %v\n", err)
		return
	}
	fmt.Printf("Saved result: %s\n", resultPath)

	// Analyze the result document for bold text
	fmt.Println("\n--- Analyzing Result Document for Bold Text ---")
	analyzeDocumentForBold(resultPath)
}

// analyzeDocumentForImages checks if images were properly embedded
func analyzeDocumentForImages(docxPath string) {
	// Open the docx as a zip file
	r, err := zip.OpenReader(docxPath)
	if err != nil {
		fmt.Printf("ERROR: Cannot open docx: %v\n", err)
		return
	}
	defer r.Close()

	// Check for image files in the archive
	imageFiles := []string{}
	documentXML := ""

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "word/media/") {
			imageFiles = append(imageFiles, f.Name)
		}
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err == nil {
				data, _ := io.ReadAll(rc)
				documentXML = string(data)
				rc.Close()
			}
		}
	}

	fmt.Printf("\nImage files found in document: %d\n", len(imageFiles))
	for _, img := range imageFiles {
		fmt.Printf("  - %s\n", img)
	}

	// Check for remaining placeholders in document.xml (correct syntax: {{#image name}})
	outsideNotReplaced := strings.Contains(documentXML, "{{#image outside_logo}}") ||
		strings.Contains(documentXML, "[IMAGE:outside_logo]")
	insideNotReplaced := strings.Contains(documentXML, "{{#image inside_logo}}") ||
		strings.Contains(documentXML, "[IMAGE:inside_logo]")

	if outsideNotReplaced {
		fmt.Println("\nBUG DETECTED: {{#image outside_logo}} placeholder OUTSIDE table was NOT replaced!")
	} else if len(imageFiles) > 0 {
		fmt.Println("\nOK: {{#image outside_logo}} placeholder was replaced with actual image")
	} else {
		fmt.Println("\nINFO: outside_logo placeholder was removed but no image embedded")
	}

	if insideNotReplaced {
		fmt.Println("BUG CONFIRMED (Issue #91): {{#image inside_logo}} placeholder in TABLE was NOT replaced!")
	} else if len(imageFiles) >= 2 {
		fmt.Println("OK: {{#image inside_logo}} placeholder in table was replaced with actual image")
	} else {
		fmt.Println("INFO: inside_logo placeholder was removed but image may not be embedded in table")
	}

	// Check for <w:drawing> elements (indicates embedded images)
	drawingCount := strings.Count(documentXML, "<w:drawing")
	fmt.Printf("\nEmbedded images (<w:drawing> elements): %d\n", drawingCount)

	if len(imageFiles) == 2 && drawingCount >= 2 {
		fmt.Println("\nRESULT: Both images appear to be properly embedded")
	} else if len(imageFiles) == 1 {
		fmt.Println("\nRESULT: Only ONE image was embedded - Issue #91 may be confirmed")
	} else if len(imageFiles) == 0 {
		fmt.Println("\nRESULT: NO images were embedded")
	}
}

// analyzeDocumentForBold checks for unexpected bold formatting
func analyzeDocumentForBold(docxPath string) {
	// Open the docx as a zip file
	r, err := zip.OpenReader(docxPath)
	if err != nil {
		fmt.Printf("ERROR: Cannot open docx: %v\n", err)
		return
	}
	defer r.Close()

	// Read document.xml
	documentXML := ""
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err == nil {
				data, _ := io.ReadAll(rc)
				documentXML = string(data)
				rc.Close()
			}
		}
	}

	// Count bold tags (excluding Heading1 which is intentionally bold)
	// Bold can appear as <w:b/>, <w:b>, or <w:b></w:b>
	boldCount := strings.Count(documentXML, "<w:b/>") +
		strings.Count(documentXML, "<w:b>") +
		strings.Count(documentXML, "<w:b ")
	fmt.Printf("Bold tags (<w:b/> or <w:b>) found: %d (including heading)\n", boldCount)

	// Check each product output for bold formatting
	products := []string{"Apple", "Orange", "Banana"}
	boldProductCount := 0

	for _, product := range products {
		if strings.Contains(documentXML, product) {
			// Find the <w:r> element that contains this product
			idx := strings.Index(documentXML, product)
			if idx > 0 {
				// Look backwards for <w:r> start (within 300 chars should be enough)
				start := idx - 300
				if start < 0 {
					start = 0
				}
				segment := documentXML[start : idx+50]

				// Find the last <w:r> opening and check for bold between it and the text
				lastRunStart := strings.LastIndex(segment, "<w:r>")
				if lastRunStart >= 0 {
					runContent := segment[lastRunStart:]
					if strings.Contains(runContent, "<w:b/>") || strings.Contains(runContent, "<w:b>") || strings.Contains(runContent, "<w:b ") {
						fmt.Printf("  - '%s': BOLD (bug!)\n", product)
						boldProductCount++
					} else {
						fmt.Printf("  - '%s': not bold (ok)\n", product)
					}
				}
			}
		}
	}

	// Check variable replacement (should not be bold)
	if strings.Contains(documentXML, "Test Customer") {
		idx := strings.Index(documentXML, "Test Customer")
		if idx > 0 {
			start := idx - 300
			if start < 0 {
				start = 0
			}
			segment := documentXML[start : idx+50]
			lastRunStart := strings.LastIndex(segment, "<w:r>")
			if lastRunStart >= 0 {
				runContent := segment[lastRunStart:]
				if strings.Contains(runContent, "<w:b/>") || strings.Contains(runContent, "<w:b>") || strings.Contains(runContent, "<w:b ") {
					fmt.Printf("  - 'Test Customer' (variable): BOLD (unexpected!)\n")
				} else {
					fmt.Printf("  - 'Test Customer' (variable): not bold (ok)\n")
				}
			}
		}
	}

	// Summary
	if boldProductCount > 0 {
		fmt.Printf("\nBUG CONFIRMED (Issue #88): %d of %d loop items are incorrectly BOLD!\n",
			boldProductCount, len(products))
	} else {
		fmt.Println("\nOK: Loop output items are not bold")
	}

	// Check for {{#each}} placeholder remnants
	if strings.Contains(documentXML, "{{#each") || strings.Contains(documentXML, "{{/each}}") {
		fmt.Println("\nWARNING: Template syntax remnants found in output!")
	}

	// Show a snippet of the XML around the text content for manual inspection
	fmt.Println("\n--- Document Content Summary ---")
	fmt.Printf("Total document.xml length: %d bytes\n", len(documentXML))
}

// createTestPNG creates a valid PNG test image
func createTestPNG(outputPath string) error {
	cmd := exec.Command("python3", "-c", `
from PIL import Image
img = Image.new('RGB', (100, 100), color='blue')
img.save('`+outputPath+`', 'PNG')
`)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%v: %s", err, stderr.String())
	}
	return nil
}

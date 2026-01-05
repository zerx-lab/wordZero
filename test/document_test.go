package test

import (
	"os"
	"testing"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func TestCreateDocument(t *testing.T) {
	// 创建新文档
	doc := document.New()

	// 添加段落
	doc.AddParagraph("Hello, World!")
	doc.AddParagraph("这是一个使用 WordZero 创建的 Word 文档。")

	// 保存文档
	err := doc.Save("test_output/test_document.docx")
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat("test_output/test_document.docx"); os.IsNotExist(err) {
		t.Fatal("Document file was not created")
	}

	// 清理测试文件
	defer os.RemoveAll("test_output")

	t.Log("Document created successfully")
}

func TestOpenDocument(t *testing.T) {
	// 首先创建一个文档
	doc := document.New()
	doc.AddParagraph("Test paragraph")

	testFile := "test_output/test_open.docx"
	err := doc.Save(testFile)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 打开文档
	openedDoc, err := document.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open document: %v", err)
	}

	// 验证内容
	paragraphs := openedDoc.Body.GetParagraphs()
	if len(paragraphs) != 1 {
		t.Fatalf("Expected 1 paragraph, got %d", len(paragraphs))
	}

	if paragraphs[0].Runs[0].Text.Content != "Test paragraph" {
		t.Fatalf("Paragraph content mismatch")
	}

	// 清理测试文件
	defer os.RemoveAll("test_output")

	t.Log("Document opened successfully")
}

func TestOpenModifySaveReopen(t *testing.T) {
	// 首先创建一个文档
	doc := document.New()
	doc.AddParagraph("Original paragraph")

	testFile := "test_output/test_modify.docx"
	err := doc.Save(testFile)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 打开文档
	openedDoc, err := document.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open document: %v", err)
	}

	// 修改文档
	openedDoc.AddParagraph("Added paragraph")

	// 保存修改后的文档
	modifiedFile := "test_output/test_modify_saved.docx"
	err = openedDoc.Save(modifiedFile)
	if err != nil {
		t.Fatalf("Failed to save modified document: %v", err)
	}

	// 再次打开修改后的文档
	reopenedDoc, err := document.Open(modifiedFile)
	if err != nil {
		t.Fatalf("Failed to reopen modified document: %v", err)
	}

	// 验证内容
	paragraphs := reopenedDoc.Body.GetParagraphs()
	if len(paragraphs) != 2 {
		t.Fatalf("Expected 2 paragraphs, got %d", len(paragraphs))
	}

	if paragraphs[0].Runs[0].Text.Content != "Original paragraph" {
		t.Fatalf("First paragraph content mismatch")
	}

	if paragraphs[1].Runs[0].Text.Content != "Added paragraph" {
		t.Fatalf("Second paragraph content mismatch")
	}

	// 清理测试文件
	defer os.RemoveAll("test_output")

	t.Log("Document open-modify-save-reopen test passed")
}

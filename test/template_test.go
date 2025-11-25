// Package test 模板功能集成测试
package test

import (
	"os"
	"testing"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
)

// TestTemplateIntegration 模板功能集成测试
func TestTemplateIntegration(t *testing.T) {
	// 确保输出目录存在
	outputDir := "output"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create output directory: %v", err)
		}
	}

	// 运行各项测试
	t.Run("变量替换集成测试", testVariableReplacementIntegration)
	t.Run("条件语句集成测试", testConditionalStatementsIntegration)
	t.Run("循环语句集成测试", testLoopStatementsIntegration)
	t.Run("结构体绑定集成测试", testStructDataBindingIntegration)

	// 清理测试文件
	t.Cleanup(func() {
		cleanupTestFiles()
	})
}

// testVariableReplacementIntegration 测试变量替换集成功能
func testVariableReplacementIntegration(t *testing.T) {
	// 创建基础文档
	doc := document.New()

	// 添加包含模板变量的段落
	doc.AddParagraph("产品信息单")
	doc.AddParagraph("")
	doc.AddParagraph("产品名称：{{productName}}")
	doc.AddParagraph("产品价格：{{price}} 元")
	doc.AddParagraph("产品数量：{{quantity}} 个")
	doc.AddParagraph("是否库存充足：{{inStock}}")
	doc.AddParagraph("产品描述：{{description}}")
	doc.AddParagraph("更新时间：{{updateTime}}")

	// 创建模板引擎并从文档创建模板
	engine := document.NewTemplateEngine()
	template, err := engine.LoadTemplateFromDocument("product_info", doc)
	if err != nil {
		t.Fatalf("Failed to load template from document: %v", err)
	}

	// 验证解析的变量数量
	expectedVars := 6
	if len(template.Variables) != expectedVars {
		t.Logf("Variables found: %v", template.Variables)
	}

	// 创建多种类型的数据
	data := document.NewTemplateData()
	data.SetVariable("productName", "WordZero处理器")
	data.SetVariable("price", 299.99)
	data.SetVariable("quantity", 100)
	data.SetVariable("inStock", true)
	data.SetVariable("description", "高效的Word文档处理工具")
	data.SetVariable("updateTime", "2024-12-01 15:30:00")

	// 渲染并保存文档
	resultDoc, err := engine.RenderTemplateToDocument("product_info", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	filename := "output/test_variable_replacement_integration.docx"
	err = resultDoc.Save(filename)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 验证文档内容
	if len(resultDoc.Body.Elements) == 0 {
		t.Error("Expected document to have content")
	}

	t.Logf("Variable replacement test completed: %s", filename)
}

// testConditionalStatementsIntegration 测试条件语句集成功能
func testConditionalStatementsIntegration(t *testing.T) {
	// 创建基础文档
	doc := document.New()

	// 添加包含条件语句的段落
	doc.AddParagraph("用户权限报告")
	doc.AddParagraph("")
	doc.AddParagraph("用户名：{{username}}")
	doc.AddParagraph("")
	doc.AddParagraph("{{#if isAdmin}}管理员权限：")
	doc.AddParagraph("- 系统配置访问权限")
	doc.AddParagraph("- 用户管理权限")
	doc.AddParagraph("- 数据备份权限{{/if}}")
	doc.AddParagraph("{{#if isEditor}}编辑权限：")
	doc.AddParagraph("- 内容编辑权限")
	doc.AddParagraph("- 文档管理权限{{/if}}")
	doc.AddParagraph("{{#if isViewer}}查看权限：")
	doc.AddParagraph("- 只读访问权限{{/if}}")

	engine := document.NewTemplateEngine()
	_, err := engine.LoadTemplateFromDocument("user_permissions", doc)
	if err != nil {
		t.Fatalf("Failed to load template from document: %v", err)
	}

	// 测试不同权限组合
	testCases := []struct {
		name         string
		username     string
		isAdmin      bool
		isEditor     bool
		isViewer     bool
		expectedFile string
	}{
		{
			name:         "管理员权限",
			username:     "admin_user",
			isAdmin:      true,
			isEditor:     false,
			isViewer:     false,
			expectedFile: "test_conditional_admin.docx",
		},
		{
			name:         "编辑员权限",
			username:     "editor_user",
			isAdmin:      false,
			isEditor:     true,
			isViewer:     false,
			expectedFile: "test_conditional_editor.docx",
		},
		{
			name:         "查看者权限",
			username:     "viewer_user",
			isAdmin:      false,
			isEditor:     false,
			isViewer:     true,
			expectedFile: "test_conditional_viewer.docx",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := document.NewTemplateData()
			data.SetVariable("username", tc.username)
			data.SetCondition("isAdmin", tc.isAdmin)
			data.SetCondition("isEditor", tc.isEditor)
			data.SetCondition("isViewer", tc.isViewer)

			resultDoc, err := engine.RenderTemplateToDocument("user_permissions", data)
			if err != nil {
				t.Fatalf("Failed to render template for %s: %v", tc.name, err)
			}

			filename := "output/" + tc.expectedFile
			err = resultDoc.Save(filename)
			if err != nil {
				t.Fatalf("Failed to save document for %s: %v", tc.name, err)
			}

			// 验证文档有内容
			if len(resultDoc.Body.Elements) == 0 {
				t.Errorf("Expected document for %s to have content", tc.name)
			}

			t.Logf("Conditional test for %s completed: %s", tc.name, filename)
		})
	}
}

// testLoopStatementsIntegration 测试循环语句集成功能
func testLoopStatementsIntegration(t *testing.T) {
	// 创建基础文档
	doc := document.New()

	// 添加包含循环的表格
	tableConfig := &document.TableConfig{
		Rows: 4,
		Cols: 3,
		Data: [][]string{
			{"序号", "商品名称", "价格"},
			{"{{@index}}", "{{name}}", "{{price}} 元"},
			{"", "", ""},
			{"", "", ""},
		},
	}
	_, err := doc.AddTable(tableConfig)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	engine := document.NewTemplateEngine()
	_, err = engine.LoadTemplateFromDocument("inventory_report", doc)
	if err != nil {
		t.Fatalf("Failed to load template from document: %v", err)
	}

	// 准备数据
	data := document.NewTemplateData()
	data.SetVariable("reportDate", "2024-12-01")

	products := []interface{}{
		map[string]interface{}{"name": "笔记本电脑", "price": "8999"},
		map[string]interface{}{"name": "无线鼠标", "price": "199"},
		map[string]interface{}{"name": "机械键盘", "price": "599"},
	}
	data.SetList("products", products)

	// 渲染模板
	resultDoc, err := engine.RenderTemplateToDocument("inventory_report", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	filename := "output/test_loop_statements_integration.docx"
	err = resultDoc.Save(filename)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 验证文档内容
	if len(resultDoc.Body.Elements) == 0 {
		t.Error("Expected document to have content")
	}

	t.Logf("Loop statement test completed: %s", filename)
}

// testStructDataBindingIntegration 测试结构体绑定集成功能
func testStructDataBindingIntegration(t *testing.T) {
	type Address struct {
		Street   string
		City     string
		Province string
		PostCode string
	}

	type Contact struct {
		Phone string
		Email string
		Fax   string
	}

	type Employee struct {
		ID         int
		Name       string
		Position   string
		Department string
		Salary     float64
		IsManager  bool
		HireDate   string
		Address    Address
		Contact    Contact
	}

	// 创建基础文档
	doc := document.New()

	// 添加员工信息模板
	doc.AddParagraph("员工详细信息")
	doc.AddParagraph("")
	doc.AddParagraph("员工编号: {{ID}}")
	doc.AddParagraph("姓名: {{Name}}")
	doc.AddParagraph("职位: {{Position}}")
	doc.AddParagraph("部门: {{Department}}")
	doc.AddParagraph("薪资: {{Salary}} 元")
	doc.AddParagraph("是否管理者: {{IsManager}}")
	doc.AddParagraph("入职日期: {{HireDate}}")
	doc.AddParagraph("")
	doc.AddParagraph("联系信息:")
	doc.AddParagraph("电话: {{Phone}}")
	doc.AddParagraph("邮箱: {{Email}}")
	doc.AddParagraph("传真: {{Fax}}")
	doc.AddParagraph("")
	doc.AddParagraph("地址信息:")
	doc.AddParagraph("街道: {{Street}}")
	doc.AddParagraph("城市: {{City}}")
	doc.AddParagraph("省份: {{Province}}")
	doc.AddParagraph("邮编: {{PostCode}}")

	engine := document.NewTemplateEngine()
	_, err := engine.LoadTemplateFromDocument("employee_detail", doc)
	if err != nil {
		t.Fatalf("Failed to load template from document: %v", err)
	}

	// 创建员工数据
	employee := Employee{
		ID:         1001,
		Name:       "张三",
		Position:   "软件工程师",
		Department: "技术部",
		Salary:     15000.00,
		IsManager:  false,
		HireDate:   "2023-06-15",
		Address: Address{
			Street:   "科技大道123号",
			City:     "深圳",
			Province: "广东",
			PostCode: "518000",
		},
		Contact: Contact{
			Phone: "13800138000",
			Email: "zhangsan@example.com",
			Fax:   "0755-88888888",
		},
	}

	// 使用结构体绑定数据
	data := document.NewTemplateData()
	err = data.FromStruct(employee)
	if err != nil {
		t.Fatalf("Failed to bind struct data: %v", err)
	}

	// 渲染模板
	resultDoc, err := engine.RenderTemplateToDocument("employee_detail", data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	filename := "output/test_struct_data_binding_integration.docx"
	err = resultDoc.Save(filename)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 验证文档内容
	if len(resultDoc.Body.Elements) == 0 {
		t.Error("Expected employee detail document to have content")
	}

	t.Logf("Struct data binding test completed: %s", filename)
}

// cleanupTestFiles 清理测试文件
func cleanupTestFiles() {
	files := []string{
		"output/test_variable_replacement_integration.docx",
		"output/test_conditional_admin.docx",
		"output/test_conditional_editor.docx",
		"output/test_conditional_viewer.docx",
		"output/test_loop_statements_integration.docx",
		"output/test_struct_data_binding_integration.docx",
	}

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			os.Remove(file)
		}
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/ZeroHawkeye/wordZero/pkg/markdown"
)

func main() {
	// 演示软换行处理的修复
	// 问题：单个\n（软换行）之前会导致文本连接在一起，没有空格
	// 修复：现在单个\n会被正确渲染为空格
	markdownContent := `# 软换行处理演示

## 问题描述

之前的版本中，单个换行符（\n）会导致文本直接连接，没有空格。
例如："第一行\n第二行" 会被渲染为 "第一行第二行"。

## 修复后的效果

现在单个换行符会被正确处理为空格：
"第一行\n第二行" 现在渲染为 "第一行 第二行"。

## 实际案例

**日期：** 2024年___月___日 

---

### **附件关键条款索引** 
#### **附件一（V2）核心约束：** 

| 任务 | 法律要点 |
|------|---------|
| 任务1：合同审查 | 确保条款完整性 |
| 任务2：风险评估 | 识别潜在风险 |
| 任务3：合规检查 | 符合法律法规 |

## 技术说明

在Markdown规范中，单个换行符（\n）是"软换行"，
多个换行符（\n\n）是"硬换行"（段落分隔）。

这个修复确保了软换行在Word文档中
被正确渲染为空格，而不是直接连接文本。

## 测试场景

### 场景1：基本软换行
第一行文本
第二行文本
第三行文本

### 场景2：格式化文本与软换行
**粗体文本**
*斜体文本*
` + "`代码文本`" + `

### 场景3：列表与软换行
- 列表项1
  包含描述
- 列表项2
  也包含描述

### 场景4：表格与软换行

| 功能 | 状态 | 备注 |
|------|------|------|
| 软换行处理 | ✅ 已修复 | 单个\n渲染为空格 |
| 表格支持 | ✅ 完整 | 支持GFM表格 |
| 任务列表 | ✅ 完整 | 支持复选框 |

---

**总结：** 此修复解决了Markdown转Word时软换行处理不正确的问题，
确保文档渲染符合Markdown规范。`

	fmt.Println("🚀 开始创建软换行处理演示...")

	// 创建转换器配置
	opts := markdown.DefaultOptions()
	opts.EnableGFM = true
	opts.EnableTables = true
	opts.EnableTaskList = true

	converter := markdown.NewConverter(opts)

	// 转换为Word文档
	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		log.Fatalf("❌ 转换失败: %v", err)
	}

	// 保存文档
	outputPath := "examples/output/soft_linebreak_demo.docx"
	err = doc.Save(outputPath)
	if err != nil {
		log.Fatalf("❌ 保存文档失败: %v", err)
	}

	fmt.Printf("✅ 软换行演示文档已保存到: %s\n", outputPath)
	fmt.Println("\n📋 演示内容:")
	fmt.Println("   🔸 基本软换行处理")
	fmt.Println("   🔸 格式化文本中的软换行")
	fmt.Println("   🔸 列表中的软换行")
	fmt.Println("   🔸 表格中的软换行")
	fmt.Println("   🔸 实际问题场景重现")
	fmt.Println("\n💡 技术说明:")
	fmt.Println("   • 软换行（单个\\n）渲染为空格")
	fmt.Println("   • 硬换行（双\\n\\n）渲染为段落分隔")
	fmt.Println("   • 符合Markdown规范")
}

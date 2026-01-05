// Package main 模板继承功能演示
package main

import (
	"fmt"
	"log"

	"github.com/zerx-lab/wordZero/pkg/document"
)

func main() {
	fmt.Println("WordZero 模板继承功能演示")
	fmt.Println("=====================================")

	// 演示基础块继承
	demonstrateBasicInheritance()

	fmt.Println("\n=====================================")
	fmt.Println("模板继承功能演示完成！")
	fmt.Println("生成的文档保存在 examples/output/ 目录下")
}

// demonstrateBasicInheritance 演示基础模板继承功能
func demonstrateBasicInheritance() {
	engine := document.NewTemplateEngine()

	// 创建基础模板
	baseTemplate := `{{companyName}} 工作报告

报告日期：{{reportDate}}

{{#block "summary"}}
默认摘要内容
本报告总结了本期的工作情况。
{{/block}}

{{#block "main_content"}}
默认主要内容
这里是详细的工作内容描述。
{{/block}}

{{#block "conclusion"}}
默认结论
总体来说，工作进展顺利。
{{/block}}

{{#block "signature"}}
报告人：{{reporterName}}
部门：{{department}}
{{/block}}`

	_, err := engine.LoadTemplate("base_report", baseTemplate)
	if err != nil {
		log.Fatalf("创建基础模板失败: %v", err)
	}

	// 创建销售报告模板（重写部分块）
	salesReportTemplate := `{{extends "base_report"}}

{{#block "summary"}}
销售业绩摘要
本月销售目标已达成 {{achievementRate}}%，超额完成预定指标。
{{/block}}

{{#block "main_content"}}
销售数据分析

本月销售业绩：
- 总销售额：{{totalSales}} 元
- 新增客户：{{newCustomers}} 个
- 成交订单：{{orders}} 笔
- 平均客单价：{{avgOrderValue}} 元

销售渠道分析：
{{#each channels}}
- {{name}}：{{sales}} 元 ({{percentage}}%)
{{/each}}

重点客户维护：
{{#each keyAccounts}}
- {{name}}：{{revenue}} 元
{{/each}}
{{/block}}

{{#block "conclusion"}}
销售总结
本月销售业绩超出预期，特别是在{{topChannel}}渠道表现突出。
建议下月重点关注{{focusArea}}市场开拓。
{{/block}}`

	_, err = engine.LoadTemplate("sales_report", salesReportTemplate)
	if err != nil {
		log.Fatalf("创建销售报告模板失败: %v", err)
	}

	// 准备数据
	data := document.NewTemplateData()
	data.SetVariable("companyName", "WordZero科技有限公司")
	data.SetVariable("reportDate", "2024年12月01日")
	data.SetVariable("reporterName", "张三")
	data.SetVariable("department", "销售部")
	data.SetVariable("achievementRate", "125")
	data.SetVariable("totalSales", "1,850,000")
	data.SetVariable("newCustomers", "45")
	data.SetVariable("orders", "183")
	data.SetVariable("avgOrderValue", "10,109")
	data.SetVariable("topChannel", "线上电商")
	data.SetVariable("focusArea", "企业级")

	// 销售渠道数据
	channels := []interface{}{
		map[string]interface{}{"name": "线上电商", "sales": "742,000", "percentage": "40.1"},
		map[string]interface{}{"name": "直销团队", "sales": "555,000", "percentage": "30.0"},
		map[string]interface{}{"name": "合作伙伴", "sales": "370,000", "percentage": "20.0"},
		map[string]interface{}{"name": "其他渠道", "sales": "183,000", "percentage": "9.9"},
	}
	data.SetList("channels", channels)

	// 重点客户数据
	keyAccounts := []interface{}{
		map[string]interface{}{"name": "大型企业A", "revenue": "280,000"},
		map[string]interface{}{"name": "知名公司B", "revenue": "195,000"},
		map[string]interface{}{"name": "科技集团C", "revenue": "150,000"},
	}
	data.SetList("keyAccounts", keyAccounts)

	// 渲染并保存
	doc, err := engine.RenderToDocument("sales_report", data)
	if err != nil {
		log.Fatalf("渲染模板失败: %v", err)
	}

	err = doc.Save("examples/output/template_inheritance_demo.docx")
	if err != nil {
		log.Fatalf("保存文档失败: %v", err)
	}

	fmt.Println("✓ 销售报告已生成：template_inheritance_demo.docx")
	fmt.Println("  - 重写了摘要、主要内容和结论块")
	fmt.Println("  - 保留了签名块的默认内容")
	fmt.Println("  - 演示了完整的块重写机制")
}

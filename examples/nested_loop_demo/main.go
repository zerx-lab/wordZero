// Package main 演示嵌套循环功能
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ZeroHawkeye/wordZero/pkg/document"
)

func main() {
	fmt.Println("=== WordZero 嵌套循环功能演示 ===")
	fmt.Println("本示例演示如何使用嵌套的 {{#each}} 循环来生成会议纪要")
	fmt.Println()

	// 创建模板引擎
	engine := document.NewTemplateEngine()

	// 创建会议纪要模板（包含嵌套循环）
	meetingMinutesTemplate := `会议纪要

会议主题：{{topic}}
会议时间：{{date}}
会议地点：{{location}}

===============================================

参会人员及任务分配：

{{#each attendees}}
【{{@index}}】 {{name}} - {{role}}
   联系方式：{{contact}}
   
   本次会议任务：
   {{#each tasks}}
   ✓ 任务 {{@index}}：{{taskName}}
     优先级：{{priority}}
     截止日期：{{deadline}}
     当前状态：{{status}}
     {{#if notes}}备注：{{notes}}{{/if}}
   
   {{/each}}

{{/each}}

===============================================

会议决议：
{{#each decisions}}
{{@index}}. {{title}}
   详情：{{details}}
   执行人：{{owner}}

{{/each}}

下次会议安排：
时间：{{nextMeetingDate}}
议题：{{nextMeetingTopic}}

记录人：{{recorder}}
审核人：{{reviewer}}`

	// 加载模板
	template, err := engine.LoadTemplate("meeting_minutes", meetingMinutesTemplate)
	if err != nil {
		log.Fatalf("加载模板失败: %v", err)
	}

	fmt.Printf("✓ 模板加载成功，包含 %d 个模板块\n", len(template.Blocks))
	fmt.Println()

	// 准备会议数据
	data := document.NewTemplateData()

	// 基本信息
	data.SetVariable("topic", "Q4季度项目进度评审会议")
	data.SetVariable("date", time.Now().Format("2006年01月02日 15:04"))
	data.SetVariable("location", "会议室A")
	data.SetVariable("nextMeetingDate", "2024年12月15日 14:00")
	data.SetVariable("nextMeetingTopic", "项目交付准备会")
	data.SetVariable("recorder", "李秘书")
	data.SetVariable("reviewer", "王总监")

	// 参会人员及其任务（嵌套结构）
	attendees := []interface{}{
		map[string]interface{}{
			"name":    "张三",
			"role":    "项目经理",
			"contact": "zhangsan@example.com",
			"tasks": []interface{}{
				map[string]interface{}{
					"taskName": "完成项目进度报告",
					"priority": "高",
					"deadline": "2024-12-10",
					"status":   "进行中",
					"notes":    "需要包含风险评估部分",
				},
				map[string]interface{}{
					"taskName": "协调跨部门资源",
					"priority": "中",
					"deadline": "2024-12-08",
					"status":   "已完成",
					"notes":    "",
				},
			},
		},
		map[string]interface{}{
			"name":    "李四",
			"role":    "技术负责人",
			"contact": "lisi@example.com",
			"tasks": []interface{}{
				map[string]interface{}{
					"taskName": "完成核心模块开发",
					"priority": "高",
					"deadline": "2024-12-12",
					"status":   "进行中",
					"notes":    "已完成80%",
				},
				map[string]interface{}{
					"taskName": "代码审查",
					"priority": "中",
					"deadline": "2024-12-09",
					"status":   "待开始",
					"notes":    "",
				},
				map[string]interface{}{
					"taskName": "编写技术文档",
					"priority": "低",
					"deadline": "2024-12-15",
					"status":   "待开始",
					"notes":    "",
				},
			},
		},
		map[string]interface{}{
			"name":    "王五",
			"role":    "测试工程师",
			"contact": "wangwu@example.com",
			"tasks": []interface{}{
				map[string]interface{}{
					"taskName": "编写测试用例",
					"priority": "高",
					"deadline": "2024-12-11",
					"status":   "进行中",
					"notes":    "已完成功能测试用例",
				},
				map[string]interface{}{
					"taskName": "执行回归测试",
					"priority": "高",
					"deadline": "2024-12-13",
					"status":   "待开始",
					"notes":    "",
				},
			},
		},
	}
	data.SetList("attendees", attendees)

	// 会议决议
	decisions := []interface{}{
		map[string]interface{}{
			"title":   "加快项目进度",
			"details": "各负责人需要在本周内完成当前阶段任务",
			"owner":   "张三",
		},
		map[string]interface{}{
			"title":   "增加测试覆盖率",
			"details": "测试覆盖率需达到85%以上",
			"owner":   "王五",
		},
		map[string]interface{}{
			"title":   "完善技术文档",
			"details": "所有核心模块需要有完整的API文档",
			"owner":   "李四",
		},
	}
	data.SetList("decisions", decisions)

	fmt.Println("✓ 数据准备完成：")
	fmt.Printf("  - %d 位参会人员\n", len(attendees))
	fmt.Printf("  - %d 项会议决议\n", len(decisions))
	fmt.Println()

	// 渲染模板
	fmt.Println("开始渲染会议纪要...")
	doc, err := engine.RenderToDocument("meeting_minutes", data)
	if err != nil {
		log.Fatalf("渲染模板失败: %v", err)
	}

	// 保存文档
	outputFile := fmt.Sprintf("examples/output/meeting_minutes_nested_loop_%s.docx",
		time.Now().Format("20060102_150405"))
	err = doc.Save(outputFile)
	if err != nil {
		log.Fatalf("保存文档失败: %v", err)
	}

	fmt.Println()
	fmt.Println("✅ 会议纪要生成成功！")
	fmt.Printf("   文件已保存至：%s\n", outputFile)
	fmt.Println()
	fmt.Println("===============================================")
	fmt.Println("功能说明：")
	fmt.Println("1. 外层循环：遍历所有参会人员")
	fmt.Println("2. 内层循环：遍历每个人的任务列表")
	fmt.Println("3. 条件判断：根据任务是否有备注显示备注信息")
	fmt.Println("4. 索引变量：使用 {{@index}} 自动编号")
	fmt.Println()
	fmt.Println("这个示例展示了嵌套循环的正确使用方式，")
	fmt.Println("现在可以正确处理任意层级的嵌套循环结构！")
	fmt.Println("===============================================")
}

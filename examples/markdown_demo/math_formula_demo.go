package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ZeroHawkeye/wordZero/pkg/markdown"
)

// 数学公式转换示例
// 演示如何将包含LaTeX数学公式的Markdown转换为Word文档
func main() {
	// 包含数学公式的Markdown内容
	markdownContent := `# 数学公式示例

本文档演示了Markdown到Word的数学公式转换功能。

## 1. 行内公式

爱因斯坦的质能方程 $E = mc^2$ 是物理学中最著名的公式之一。

圆的面积公式是 $S = \pi r^2$，其中 $\pi \approx 3.14159$。

## 2. 块级公式

### 二次方程求根公式
$$x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$$

### 勾股定理
$$a^2 + b^2 = c^2$$

## 3. 希腊字母

物理学和数学中常用的希腊字母：
- $\alpha$ (alpha) - 角度、系数
- $\beta$ (beta) - 角度、系数
- $\gamma$ (gamma) - 角度、伽马函数
- $\delta$ (delta) - 变化量
- $\theta$ (theta) - 角度
- $\lambda$ (lambda) - 波长
- $\mu$ (mu) - 微、摩擦系数
- $\pi$ (pi) - 圆周率
- $\sigma$ (sigma) - 求和、标准差
- $\omega$ (omega) - 角频率

## 4. 数学运算符

- 乘法：$a \times b$
- 除法：$a \div b$
- 加减：$a \pm b$
- 小于等于：$a \leq b$
- 大于等于：$a \geq b$
- 不等于：$a \neq b$
- 约等于：$a \approx b$
- 恒等于：$a \equiv b$

## 5. 上下标

- 平方：$x^2$
- 立方：$x^3$
- n次方：$x^n$
- 下标：$x_i$
- 复合：$x_i^2$
- 求和下标：$\sum_{i=1}^n x_i$

## 6. 分数与根号

分数：$\frac{1}{2}$, $\frac{a+b}{c-d}$

平方根：$\sqrt{2}$, $\sqrt{x^2 + y^2}$

立方根：$\sqrt[3]{8} = 2$

## 7. 微积分符号

积分：$\int_0^1 x dx = \frac{1}{2}$

偏导：$\frac{\partial f}{\partial x}$

梯度：$\nabla f$

## 8. 集合符号

- 属于：$x \in A$
- 不属于：$x \notin A$
- 子集：$A \subset B$
- 交集：$A \cap B$
- 并集：$A \cup B$
- 空集：$\emptyset$

## 9. 逻辑符号

- 任意：$\forall x$
- 存在：$\exists x$
- 非：$\neg p$
- 与：$p \land q$
- 或：$p \lor q$
- 蕴含：$p \Rightarrow q$
- 等价：$p \Leftrightarrow q$

## 10. 极限与无穷

- 极限：$\lim_{x \to 0} \frac{\sin x}{x} = 1$
- 无穷大：$\infty$
- 趋近于：$x \to 0$
`

	// 创建转换器，启用数学公式支持
	opts := markdown.DefaultOptions()
	opts.EnableMath = true

	converter := markdown.NewConverter(opts)

	// 转换为Word文档
	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	// 确保输出目录存在
	outputDir := "../output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	// 保存文档
	outputPath := outputDir + "/math_formula_demo.docx"
	if err := doc.Save(outputPath); err != nil {
		log.Fatalf("保存文档失败: %v", err)
	}

	fmt.Printf("数学公式文档已生成: %s\n", outputPath)
	fmt.Println("\n支持的数学公式语法:")
	fmt.Println("- 行内公式: $E = mc^2$")
	fmt.Println("- 块级公式: $$x = \\frac{-b \\pm \\sqrt{b^2 - 4ac}}{2a}$$")
	fmt.Println("- 希腊字母: \\alpha, \\beta, \\gamma, \\pi, \\sigma 等")
	fmt.Println("- 运算符: \\times, \\div, \\pm, \\leq, \\geq 等")
	fmt.Println("- 上下标: x^2, x_i, x^{n+1}, x_{i,j}")
	fmt.Println("- 分数: \\frac{a}{b}")
	fmt.Println("- 根号: \\sqrt{x}, \\sqrt[3]{x}")
}

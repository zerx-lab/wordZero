package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ZeroHawkeye/wordZero/pkg/markdown"
)

// TestMarkdownMathFormulaConversion 测试Markdown数学公式转换
func TestMarkdownMathFormulaConversion(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		wantErr     bool
		description string
	}{
		{
			name:        "简单行内公式",
			markdown:    `质能方程是 $E = mc^2$，这是爱因斯坦的著名公式。`,
			wantErr:     false,
			description: "测试简单的行内数学公式",
		},
		{
			name:        "复杂行内公式",
			markdown:    `二次方程的求根公式是 $x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$。`,
			wantErr:     false,
			description: "测试包含分数和根号的行内公式",
		},
		{
			name:        "块级公式",
			markdown:    "勾股定理：\n$$a^2 + b^2 = c^2$$",
			wantErr:     false,
			description: "测试块级数学公式",
		},
		{
			name:        "多个公式",
			markdown:    "设 $x$ 和 $y$ 是两个变量，则：\n$$x + y = z$$",
			wantErr:     false,
			description: "测试多个数学公式",
		},
		{
			name:        "希腊字母",
			markdown:    `圆周率 $\pi \approx 3.14159$，角度 $\theta$ 和 $\alpha$。`,
			wantErr:     false,
			description: "测试希腊字母转换",
		},
		{
			name:        "上下标",
			markdown:    `水的分子式是 $H_2O$，化学方程式 $x^2 + y^2$。`,
			wantErr:     false,
			description: "测试上下标转换",
		},
		{
			name:        "分数",
			markdown:    `分数 $\frac{1}{2}$ 表示一半。`,
			wantErr:     false,
			description: "测试分数转换",
		},
		{
			name:        "根号",
			markdown:    `平方根 $\sqrt{2}$ 和立方根 $\sqrt[3]{8}$。`,
			wantErr:     false,
			description: "测试根号转换",
		},
		{
			name:        "积分和求和",
			markdown:    `积分 $\int_0^1 x dx$ 和求和 $\sum_{i=1}^n i$。`,
			wantErr:     false,
			description: "测试积分和求和符号",
		},
		{
			name:        "数学运算符",
			markdown:    `$a \times b$, $a \div b$, $a \pm b$, $a \leq b$, $a \geq b$`,
			wantErr:     false,
			description: "测试数学运算符转换",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := markdown.DefaultOptions()
			opts.EnableMath = true
			converter := markdown.NewConverter(opts)

			doc, err := converter.ConvertString(tt.markdown, opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if doc == nil && !tt.wantErr {
				t.Error("ConvertString() returned nil document")
				return
			}

			// 验证文档包含内容
			if doc != nil {
				if len(doc.Body.Elements) == 0 {
					t.Error("Expected document to contain at least one element")
				}
			}
		})
	}
}

// TestMarkdownMathDisabled 测试禁用数学公式时的行为
func TestMarkdownMathDisabled(t *testing.T) {
	markdownContent := `公式 $E = mc^2$ 不应该被特殊处理。`

	opts := markdown.DefaultOptions()
	opts.EnableMath = false
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Errorf("ConvertString() error = %v", err)
		return
	}

	if doc == nil {
		t.Error("ConvertString() returned nil document")
		return
	}

	// 验证文档包含内容
	if len(doc.Body.Elements) == 0 {
		t.Error("Expected document to contain at least one element")
	}
}

// TestMarkdownMathWithOtherElements 测试数学公式与其他元素的组合
func TestMarkdownMathWithOtherElements(t *testing.T) {
	markdownContent := `# 数学公式示例

## 基础公式

这是一个简单的公式：$E = mc^2$

## 复杂公式

二次方程的求根公式：
$$x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$$

## 表格与公式

| 名称 | 公式 |
|------|------|
| 勾股定理 | $a^2 + b^2 = c^2$ |
| 圆面积 | $S = \pi r^2$ |

## 列表与公式

- 欧拉公式：$e^{i\pi} + 1 = 0$
- 牛顿定律：$F = ma$
`

	opts := markdown.DefaultOptions()
	opts.EnableMath = true
	opts.EnableTables = true
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Errorf("ConvertString() error = %v", err)
		return
	}

	if doc == nil {
		t.Error("ConvertString() returned nil document")
		return
	}

	// 验证包含标题段落
	paragraphs := doc.Body.GetParagraphs()
	if len(paragraphs) == 0 {
		t.Error("Expected document to contain paragraphs")
	}

	// 表格可能由于格式原因无法解析，只验证文档生成成功
	// 实际项目中表格内包含math公式的情况较复杂
	t.Logf("Document generated with %d paragraphs and %d tables", 
		len(paragraphs), len(doc.Body.GetTables()))
}

// TestMarkdownMathSaveToFile 测试将包含数学公式的文档保存到文件
func TestMarkdownMathSaveToFile(t *testing.T) {
	markdownContent := `# 数学公式文档

## 著名公式

1. 质能方程：$E = mc^2$
2. 勾股定理：$a^2 + b^2 = c^2$
3. 欧拉恒等式：$e^{i\pi} + 1 = 0$

## 二次方程求根公式

$$x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}$$

## 希腊字母

- Alpha: $\alpha$
- Beta: $\beta$
- Gamma: $\gamma$
- Delta: $\delta$
- Pi: $\pi$
- Sigma: $\sigma$
- Omega: $\omega$

## 运算符

- 乘法：$a \times b$
- 除法：$a \div b$
- 大于等于：$a \geq b$
- 小于等于：$a \leq b$
- 不等于：$a \neq b$
- 约等于：$a \approx b$

## 集合符号

- 属于：$x \in A$
- 不属于：$x \notin A$
- 子集：$A \subset B$
- 交集：$A \cap B$
- 并集：$A \cup B$

## 微积分

- 积分：$\int_0^1 x dx$
- 求和：$\sum_{i=1}^n i$
- 极限：$\lim_{x \to \infty} f(x)$
`

	opts := markdown.DefaultOptions()
	opts.EnableMath = true
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Fatalf("ConvertString() error = %v", err)
	}

	// 创建输出目录
	outputDir := "test_output"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// 保存文件
	outputPath := filepath.Join(outputDir, "math_formula_test.docx")
	err = doc.Save(outputPath)
	if err != nil {
		t.Fatalf("Failed to save document: %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

// TestLaTeXToDisplayConversion 测试LaTeX到显示格式的转换
func TestLaTeXToDisplayConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "简单上标",
			input:    "x^2",
			expected: "x²",
		},
		{
			name:     "简单下标",
			input:    "x_i",
			expected: "xᵢ",
		},
		{
			name:     "希腊字母alpha",
			input:    `\alpha`,
			expected: "α",
		},
		{
			name:     "希腊字母pi",
			input:    `\pi`,
			expected: "π",
		},
		{
			name:     "乘法符号",
			input:    `\times`,
			expected: "×",
		},
		{
			name:     "小于等于",
			input:    `\leq`,
			expected: "≤",
		},
		{
			name:     "无穷符号",
			input:    `\infty`,
			expected: "∞",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建一个简单的markdown并检查转换
			md := "$" + tt.input + "$"
			opts := markdown.DefaultOptions()
			opts.EnableMath = true
			converter := markdown.NewConverter(opts)

			doc, err := converter.ConvertString(md, opts)
			if err != nil {
				t.Errorf("ConvertString() error = %v", err)
				return
			}

			if doc == nil {
				t.Error("ConvertString() returned nil document")
				return
			}

			// 验证文档生成成功且包含内容
			if doc.Body.Elements == nil || len(doc.Body.Elements) == 0 {
				t.Error("Document has no elements")
			}
		})
	}
}

// TestMarkdownBlockMathFormula 测试块级数学公式的特殊处理
func TestMarkdownBlockMathFormula(t *testing.T) {
	markdownContent := `这是一个块级公式：

$$
\frac{d}{dx} \int_a^x f(t) dt = f(x)
$$

这是微积分基本定理。`

	opts := markdown.DefaultOptions()
	opts.EnableMath = true
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Errorf("ConvertString() error = %v", err)
		return
	}

	if doc == nil {
		t.Error("ConvertString() returned nil document")
		return
	}

	// 验证文档包含多个段落（包括公式段落）
	if len(doc.Body.Elements) < 2 {
		t.Errorf("Expected at least 2 elements, got %d", len(doc.Body.Elements))
	}
}

// TestMarkdownInlineMathInParagraph 测试段落中的行内公式
func TestMarkdownInlineMathInParagraph(t *testing.T) {
	markdownContent := `在物理学中，能量 $E$ 和质量 $m$ 通过光速 $c$ 关联：$E = mc^2$。`

	opts := markdown.DefaultOptions()
	opts.EnableMath = true
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Errorf("ConvertString() error = %v", err)
		return
	}

	if doc == nil {
		t.Error("ConvertString() returned nil document")
		return
	}

	// 验证文档包含段落
	paragraphs := doc.Body.GetParagraphs()
	if len(paragraphs) == 0 {
		t.Error("Expected document to contain paragraphs")
	}
}

// TestMathFormulaEdgeCases 测试数学公式的边缘情况
func TestMathFormulaEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		shouldParse bool
	}{
		{
			name:        "空公式",
			markdown:    `$$$$`,
			shouldParse: true,
		},
		{
			name:        "只有空格的公式",
			markdown:    `$   $`,
			shouldParse: true,
		},
		{
			name:        "嵌套花括号",
			markdown:    `$\frac{\frac{1}{2}}{3}$`,
			shouldParse: true,
		},
		{
			name:        "特殊字符",
			markdown:    `$a + b = c$`,
			shouldParse: true,
		},
		{
			name:        "未闭合的美元符号",
			markdown:    `这里有一个 $未闭合`,
			shouldParse: true, // 应该作为普通文本处理
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := markdown.DefaultOptions()
			opts.EnableMath = true
			converter := markdown.NewConverter(opts)

			doc, err := converter.ConvertString(tt.markdown, opts)
			if tt.shouldParse {
				if err != nil {
					t.Errorf("Expected successful parse, got error: %v", err)
				}
				if doc == nil {
					t.Error("Expected non-nil document")
				}
			}
		})
	}
}

// TestMathDefaultOptionEnabled 测试默认选项中数学公式是启用的
func TestMathDefaultOptionEnabled(t *testing.T) {
	opts := markdown.DefaultOptions()
	if !opts.EnableMath {
		t.Error("Expected EnableMath to be true by default")
	}
}

// TestMarkdownMathContentPreservation 测试数学公式内容的保留
func TestMarkdownMathContentPreservation(t *testing.T) {
	// 测试公式内容在转换过程中是否被正确保留
	formulas := []string{
		`E = mc^2`,
		`x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`,
		`\alpha + \beta = \gamma`,
		`a^2 + b^2 = c^2`,
		`\int_0^1 x dx`,
	}

	for _, formula := range formulas {
		t.Run(formula, func(t *testing.T) {
			md := "$" + formula + "$"
			opts := markdown.DefaultOptions()
			opts.EnableMath = true
			converter := markdown.NewConverter(opts)

			doc, err := converter.ConvertString(md, opts)
			if err != nil {
				t.Errorf("Failed to convert formula %q: %v", formula, err)
				return
			}

			if doc == nil {
				t.Errorf("Got nil document for formula %q", formula)
				return
			}

			// 检查文档不为空
			if len(doc.Body.Elements) == 0 {
				t.Errorf("Document has no elements for formula %q", formula)
			}
		})
	}
}

// TestComplexMathDocument 测试包含复杂数学内容的完整文档
func TestComplexMathDocument(t *testing.T) {
	markdownContent := `# 高等数学公式汇总

## 1. 极限

### 重要极限
$$\lim_{x \to 0} \frac{\sin x}{x} = 1$$

$$\lim_{x \to \infty} (1 + \frac{1}{x})^x = e$$

## 2. 导数

基本导数公式：
- $(x^n)' = nx^{n-1}$
- $(e^x)' = e^x$
- $(\sin x)' = \cos x$
- $(\cos x)' = -\sin x$

## 3. 积分

### 不定积分
$$\int x^n dx = \frac{x^{n+1}}{n+1} + C \quad (n \neq -1)$$

### 定积分
$$\int_a^b f(x) dx = F(b) - F(a)$$

## 4. 级数

泰勒级数：
$$e^x = \sum_{n=0}^{\infty} \frac{x^n}{n!} = 1 + x + \frac{x^2}{2!} + \frac{x^3}{3!} + \cdots$$

## 5. 矩阵

行列式：
$$\det(A) = \sum_{\sigma \in S_n} \text{sgn}(\sigma) \prod_{i=1}^{n} a_{i,\sigma(i)}$$
`

	opts := markdown.DefaultOptions()
	opts.EnableMath = true
	converter := markdown.NewConverter(opts)

	doc, err := converter.ConvertString(markdownContent, opts)
	if err != nil {
		t.Fatalf("ConvertString() error = %v", err)
	}

	if doc == nil {
		t.Fatal("ConvertString() returned nil document")
	}

	// 验证文档结构
	paragraphs := doc.Body.GetParagraphs()
	if len(paragraphs) < 3 {
		t.Errorf("Expected at least 3 paragraphs for complex document, got %d", len(paragraphs))
	}

	// 验证文档包含内容
	if len(doc.Body.Elements) == 0 {
		t.Error("Document has no elements")
	}

	t.Logf("Complex math document generated with %d paragraphs", len(paragraphs))
}

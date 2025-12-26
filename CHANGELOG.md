# WordZero 更新日志

## [Unreleased]

## [v1.6.0] - 2025-12-26

### 🐛 修复

#### 模板页眉页脚引用修复 ✨ **重要修复**
- **修复问题**: 解决了DOCX渲染时模板页眉页脚引用丢失的问题
- **技术细节**: 
  - 解析 `<w:sectPr>` 块，即使 Word 将其嵌套在段落属性中
  - 通过新的 `setSectionProperties` 辅助方法保持读取 DOCX 文件时的关系引用
  - 防止序列化时 `r:id` 引用被丢弃
- **验证**: 添加了回归测试，确认重写 DOCX 后页眉页脚变量仍能正确替换

#### OpenFromMemory bug 修复
- **修复问题**: 修复了通过 `OpenFromMemory` 方法读取 Word 模板后，生成内容再保存时导致 Microsoft Office Word 无法打开文件的问题
- **问题表现**: 使用 Office 2016 等版本打开时会提示文件错误
- **感谢**: 感谢 @xNeo404 的贡献

#### Table.GetCellText 修复
- `Table.GetCellText` 现在返回单元格内所有段落与所有 Run 的完整文本，并以 `\n` 连接段落，修复此前只能获取第一段首个 Run 文本导致多行内容丢失的问题。
  - 影响：如果下游代码假设无换行符，需自行 `strings.ReplaceAll(text, "\n", "")` 或按需拆分。
  - 限制：同一段落内的 `<w:br/>` 软换行尚未单独解析（未来可扩展）。

---

## [v1.5.0] - 2025-11-25

### 🚀 新增功能

#### LaTeX 数学公式支持 ✨ **新增**
- **Markdown转Word**: 支持 LaTeX 数学公式转换
- **技术细节**: 修复 LaTeX 命令替换顺序，防止部分替换问题

#### 表格单元格复杂内容 API ✨ **新增**
- **AddCellParagraph**: 向表格单元格添加段落
- **AddCellFormattedParagraph**: 添加带格式的段落
- **ClearCellParagraphs/GetCellParagraphs**: 单元格段落管理
- **AddNestedTable**: 支持在单元格中添加嵌套表格
- **AddCellList**: 支持多种列表类型
- **AddCellImage/AddCellImageFromFile/AddCellImageFromData**: 单元格图片插入

#### SnapToGrid 网格对齐支持 ✨ **新增**
- **自定义样式**: 添加 SnapToGrid 字段到段落属性
- **段落格式**: 添加 SnapToGrid 选项用于禁用网格对齐
- **用途**: 设置为 false 可使自定义 LineSpacing 在启用网格的文档中正常工作

#### 段落级文本格式化方法 ✨ **新增**
- 添加下划线、粗体、斜体等段落级文本格式化方法

#### 段落分页符 ✨ **新增**
- **Paragraph.AddPageBreak()**: 在段落内添加分页符

#### 页眉页脚格式化支持 ✨ **新增**
- 支持带 TextFormat 和对齐选项的格式化页眉页脚

#### 段落格式属性 ✨ **新增**
- **KeepNext**: 与下段同页
- **KeepLines**: 段中不分页
- **PageBreakBefore**: 段前分页
- **WidowControl**: 孤行控制
- **OutlineLevel**: 大纲级别

### 🐛 修复

#### 模板引擎修复
- 修复模板变量替换时图片丢失问题，添加 drawing 元素解析
- 修复模板变量替换时段落编号丢失问题
- 修复模板引擎支持页眉页脚变量
- 修复模板渲染保持文档结构（页眉、页脚等）

#### 表格功能修复
- **AddTable/CreateTable**: 现在返回错误而非 nil，提供更好的错误处理
- 修复模板导出时内嵌表格消失问题

#### 其他修复
- 修复浮动图片导致 Word 文档打开失败的问题
- 修复 Markdown 转换时软换行处理问题
- 修复非 ASCII 图片文件名导致文档损坏的问题
- 修复嵌套 `{{#each}}` 循环渲染问题
- 修复 `UpdateTOC` 无法定位 `GenerateTOC` 创建的目录问题
- 修复表格单元格合并后对齐问题

---

## [v1.4.0] - 2025-11-11

### 🐛 修复

感谢以下贡献者的 PR：
- @litecn 修复 xml 规范问题
- @CoffeeSwt 混合格式段落字体下划线和删除线问题修复
- @xNeo404 实现从内存中打开文档的新增 `OpenFromMemory` 函数/修复 md 转 word 问题
- @Padane22-spec 修复通过 Open 打开的文档在 Save 后无法正常打开的问题

---

## [v1.3.9] - 2025-06-06

### 🐛 修复

#### 完全遵循 OOXML 规范的修复版本
- 修复添加图片 Word 无法打开问题
- 修复 md 转 word 代码块样式问题
- 模板渲染增加图片占位功能

---

## [v1.3.8] - 2025-06-05

### 🔧 改进
- 模板渲染部分代码更新，修复多处继承样式问题
- 添加缺失的样式解析
- 添加多语言文档支持

---

## [v1.3.7] - 2025-06-04

### 🚀 新增功能
- 增加 markdown 导入导出功能
- 增加 if else 模板语法支持

### 🐛 修复
- 修复模板渲染问题

---

## [v1.3.6] - 2025-06-04

### 🚀 新增功能

### markdown导入导出功能

#### if else 条件语句支持 ✨ **全新实现**
- **完整条件语句系统**: 实现了模板引擎中的 `{{#if}}...{{#else}}...{{/if}}` 语法支持
- **多种条件判断**: 支持布尔值、字符串、数字等多种数据类型的条件判断
- **嵌套条件支持**: 支持 if else 语句与循环语句的任意嵌套组合
- **语法特性**:
  - `{{#if condition}}` - 条件开始
  - `{{#else}}` - 否则分支
  - `{{/if}}` - 条件结束
  - 支持空值、零值的智能判断
- **使用示例**:
  ```
  {{#if isVIP}}
    尊贵的VIP客户：{{name}}
  {{#else}}
    普通客户：{{name}}
  {{/if}}
  ```

### 🐛 问题修复

#### 模板渲染引擎全面优化 ✨ **重要修复**
- **修复问题**: 解决了复杂模板渲染中的多个关键问题
- **影响范围**: 所有使用模板功能的场景，特别是包含条件语句和循环的复杂模板
- **问题表现**: 
  - if else 语句在某些情况下解析不正确
  - 嵌套模板语法处理异常
  - 条件判断逻辑不够准确
  - 模板语法识别存在边界问题
- **修复方案**:
  - **语法解析优化**: 改进模板语法的正则表达式匹配逻辑
  - **条件判断增强**: 完善布尔值、空值、零值的判断机制
  - **嵌套处理改进**: 优化嵌套结构的解析和渲染顺序
  - **错误处理增强**: 增加详细的错误信息和调试支持

#### 模板引擎性能优化
- **渲染效率提升**: 优化模板解析算法，提高渲染速度
- **内存使用优化**: 减少模板渲染过程中的内存占用
- **递归安全性**: 增强嵌套模板的递归处理安全性

### 🔧 技术架构改进

#### 条件语句引擎
- **新增方法**:
  - `renderIfElseStatements()` - 处理 if else 语句渲染
  - `evaluateCondition()` - 条件表达式求值
  - `findMatchingElse()` - 匹配对应的 else 分支
  - `findMatchingEndIf()` - 匹配对应的 endif 标记
- **类型安全**: 支持多种Go数据类型的条件判断
- **语法兼容**: 与现有循环语法完全兼容，支持任意嵌套

#### 模板解析优化
- **正则表达式改进**: 更精确的模板语法识别
- **解析顺序优化**: 确保嵌套结构按正确顺序处理
- **错误恢复机制**: 语法错误时提供有用的错误信息

### 🎯 功能完整性

#### 条件语句特性
- ✅ **基础条件判断**: 支持 `{{#if variable}}` 基础语法
- ✅ **else 分支**: 支持 `{{#else}}` 否则分支
- ✅ **嵌套条件**: 支持条件语句内嵌套其他模板语法
- ✅ **循环内条件**: 支持在循环内使用条件语句
- ✅ **多种数据类型**: 支持 `bool`, `string`, `int`, `float64` 等类型

#### 渲染稳定性
- ✅ **语法容错**: 对不规范的模板语法提供友好错误提示
- ✅ **性能稳定**: 复杂模板渲染性能大幅提升
- ✅ **内存安全**: 避免大模板渲染时的内存泄漏
- ✅ **类型安全**: 条件判断时的类型转换安全可靠

### 📚 示例程序更新

#### 条件语句演示
- 更新现有模板演示程序，增加 if else 语法示例
- 新增复杂嵌套模板的使用案例
- 提供不同数据类型的条件判断示例

### 🔍 质量改进

#### 模板引擎稳定性
- ✅ **语法完整性**: 支持完整的条件语句语法体系
- ✅ **嵌套兼容性**: 与现有循环、变量语法完美兼容
- ✅ **性能优化**: 模板渲染效率显著提升
- ✅ **错误处理**: 提供详细的模板语法错误信息

---

## [v1.3.5] - 2025-06-04

### 🐛 模板样式保持问题修复 ✨ **重要修复**

#### 深度样式复制机制完善
- **修复问题**: 解决了模板渲染中字体、字体大小、表格边框、单元格水平居中等样式丢失的问题
- **影响范围**: 使用手动创建的Word模板文件进行变量替换的所有场景
- **问题表现**: 
  - 字体信息丢失（FontFamily信息未正确复制）
  - 字体大小丢失（FontSize属性未正确保持）
  - 表格边框样式丢失（TableBorders和TableCellBorders属性未深度复制）
  - 单元格水平居中对齐丢失（VAlign和段落对齐属性未完整复制）
  - 文字颜色保持正常（Color属性复制正确）

#### 深度复制机制重构
- **完整属性复制**: 重构了整个样式复制机制，确保所有样式属性的深度复制
  - `cloneParagraph()` - 完整复制段落及其所有属性
  - `cloneParagraphProperties()` - 深度复制段落属性（对齐、间距、缩进等）
  - `cloneRun()` - 完整复制文本运行及其格式
  - `cloneRunProperties()` - 深度复制文本运行属性（粗体、斜体、字体、颜色等）
  - `cloneTable()` - 完整复制表格及其所有属性
  - `cloneTableProperties()` - 深度复制表格属性（边框、样式、布局等）
  - `cloneTableCellProperties()` - 深度复制单元格属性（边框、对齐、合并等）

#### 字体样式修复 ✨
- **字体族复制**: 完整复制 `FontFamily` 属性，包含 ASCII、HAnsi、EastAsia、CS 字段
- **字体大小保持**: 正确复制 `FontSize` 属性，保持原有字体大小设置
- **字体颜色保持**: 继续正确保持 `Color` 属性（原本已正常工作）

#### 表格样式修复 ✨
- **表格边框保持**: 深度复制 `TableBorders` 所有边框属性（上下左右、内部边框）
- **单元格边框保持**: 深度复制 `TableCellBorders` 包括对角线边框
- **边框详细属性**: 完整保持边框样式、粗细、颜色、主题颜色等所有属性

#### 单元格对齐修复 ✨
- **垂直对齐保持**: 正确复制 `VAlign` 属性，保持单元格垂直对齐设置
- **水平对齐保持**: 通过段落 `Justification` 属性保持水平对齐
- **文字方向保持**: 复制 `TextDirection` 属性，保持文字方向设置

#### 其他样式属性
- **网格跨度**: 正确复制 `GridSpan` 属性，保持单元格合并状态
- **垂直合并**: 正确复制 `VMerge` 属性，保持行合并状态
- **单元格边距**: 深度复制 `TableCellMarginsCell` 属性
- **底纹样式**: 正确复制表格和单元格的底纹/背景色设置

### 🔧 技术改进

#### 代码结构优化
- **模块化复制方法**: 将复杂的复制逻辑分解为多个专门的方法
- **类型安全**: 修复了所有Go结构体类型兼容性问题
- **深度复制**: 确保所有嵌套对象都被正确的深度复制而非浅拷贝

### 🎯 修复验证

#### 修复效果确认
- ✅ **字体保持**: 模板中设置的字体族信息完全保持
- ✅ **字体大小保持**: 所有字体大小设置正确保持  
- ✅ **表格边框保持**: 表格和单元格边框样式完整保持
- ✅ **单元格对齐保持**: 水平和垂直对齐设置正确保持
- ✅ **颜色保持**: 文字颜色继续正确保持（原本正常）

#### 模板兼容性
- ✅ **程序生成模板**: enhanced_template_demo等程序生成的模板继续正常工作
- ✅ **手动Word模板**: 从Microsoft Word或WPS手动创建的模板现在可以正确保持所有样式
- ✅ **复杂样式模板**: 包含复杂格式设置的模板现在可以完美渲染

---

## [v1.3.4] - 2025-06-03

### 🚀 重大功能修复

#### 模板功能重大重构 ✨
- **修复问题**: 彻底解决了模板功能中样式丢失、表格处理错误的问题
- **影响范围**: 从文档模板生成功能的完整重构，影响所有使用模板的场景
- **错误表现**: 
  - 普通变量样式丢失（粗体、颜色、字体等格式信息）
  - 表格被错误转换为标签形式而非保持原始结构
  - 文档格式在模板渲染过程中完全丢失
  - 复杂文档结构无法正确保持
- **重构方案**:
  - **保持文档结构**: 不再将文档转换为纯文本，直接在原始文档结构上进行变量替换
  - **新增深度复制**: 实现 `cloneDocument()` 方法，完整复制所有文档元素和属性
  - **直接结构替换**: 创建 `RenderTemplateToDocument()` 主要渲染方法
  - **段落格式保持**: 实现 `replaceVariablesInParagraph()` 保持所有文本格式
  - **表格模板支持**: 完整的表格模板循环功能，支持 `{{#each items}}` 语法

#### 表格模板功能 ✨ **全新实现**
- **表格循环渲染**: 完整支持表格行循环模板
- **模板语法支持**: 
  - `{{#each items}}` - 表格数据循环
  - `{{name}}`, `{{position}}` - 单元格变量替换
  - `{{/each}}` - 循环结束标记
- **样式保持**: 表头样式、单元格格式完全保持
- **智能检测**: 自动检测表格是否包含模板语法
- **关键功能**:
  - `isTableTemplate()` - 检测表格模板
  - `renderTableTemplate()` - 渲染表格模板
  - `cloneTableRow()` - 克隆表格行保持所有属性

#### 样式保持功能 ✨ **完全修复**
- **文本格式保持**: 粗体、颜色、字体大小、对齐方式等所有格式
- **段落属性保持**: 段落级别的格式设置完整保持
- **运行属性保持**: 文本运行级别的所有属性（颜色、字体等）
- **复合变量支持**: 单个段落内多个变量混合替换保持格式
- **示例效果**:
  ```
  原文: 作者：{{author}} | 日期：{{date}}  (蓝色粗体)
  结果: 作者：张开发 | 日期：2025年06月03日  (保持蓝色粗体)
  ```

### 🔧 技术架构改进

#### 结构体字段类型修复
- **Paragraph.Runs**: 从 `[]*Run` 修正为 `[]Run`
- **Text 结构**: 从 `*Text` 修正为 `Text`
- **TableRow.Cells**: 从 `[]*TableCell` 修正为 `[]TableCell`
- **TableCell 字段**: 从 `Content` 修正为 `Paragraphs`
- **类型兼容性**: 修复了大量Go结构体类型不兼容问题

#### 新增核心方法
- **`RenderTemplateToDocument()`**: 新的主要模板渲染方法
- **`replaceVariablesInDocument()`**: 文档级变量替换
- **`replaceVariablesInParagraph()`**: 段落级变量替换保持格式
- **`replaceVariablesInTable()`**: 表格变量替换和模板处理
- **`cloneDocument()`**: 深度复制文档所有元素
- **`cloneTableRow()`**: 表格行克隆保持所有属性

### 🎯 演示程序完善

#### 增强模板演示
- **新增文件**: `examples/enhanced_template_demo/enhanced_template_demo.go`
- **三大演示**:
  1. **样式变量模板**: 展示文本格式保持（粗体、颜色、字体、对齐）
  2. **表格模板功能**: 表头样式和数据循环渲染
  3. **复杂文档模板**: 多功能组合演示
- **实际效果验证**: 
  - 样式变量完全保持格式
  - 表格模板正确循环生成
  - 复杂文档结构完整保持

### 🔍 质量改进

#### 模板引擎稳定性
- ✅ **格式完整性**: 模板渲染过程中所有样式和格式保持
- ✅ **结构保持**: Word文档结构在整个模板过程中维护完整性  
- ✅ **表格支持**: 维护表格结构和样式的表格模板循环
- ✅ **类型安全**: 修复所有Go结构体类型兼容性问题
- ✅ **功能完整**: 单个段落内混合格式文本的正确保持

---

## [v1.3.3] - 2025-06-02

### 🐛 问题修复

#### 页面设置保存和加载问题修复 ✨ **重要修复**
- **修复问题**: 解决了页面设置在文档保存和重新加载后丢失的问题
- **影响范围**: 主要影响页面配置功能，包括页面尺寸、方向、边距等设置
- **错误表现**: 
  - 设置页面为Letter横向，保存后重新打开变成A4纵向
  - XML结构中SectionProperties位置不正确
  - 页面设置解析失败，返回默认配置
- **根本原因**:
  - `getSectionProperties()` 方法只检查Elements数组的最后一个元素
  - 文档序列化时SectionProperties被放在body开头，违反了Word XML规范
  - 解析时无法正确找到SectionProperties元素
- **修复方案**:
  - 修改 `getSectionProperties()` 方法，在整个Elements数组中查找SectionProperties
  - 优化 `Body.MarshalXML()` 方法，确保SectionProperties始终位于body末尾
  - 遵循OpenXML规范，将sectPr放在正确位置
- **修复后效果**:
  - 页面设置保存后正确加载：Letter横向 → Letter横向 ✓
  - XML结构符合Word标准：`<w:body><w:p>...</w:p><w:sectPr>...</w:sectPr></w:body>`
  - 所有页面配置（尺寸、方向、边距等）正确保持

#### 技术细节
- **修改文件**: 
  - `pkg/document/page.go` - 修复 `getSectionProperties()` 方法
  - `pkg/document/document.go` - 优化 `Body.MarshalXML()` 序列化逻辑
- **修改内容**: 
  - 在Elements数组中全局搜索SectionProperties而非只检查最后一个元素
  - 序列化时分离SectionProperties和其他元素，确保sectPr在body末尾
  - 移除位置假设，提高容错性
- **影响功能**: 
  - 所有页面设置功能（SetPageSettings, GetPageSettings等）
  - 文档保存和加载的完整性
  - XML文档结构的规范性

### 🔍 质量改进

#### XML结构规范性
- ✅ **符合OpenXML规范**: SectionProperties现在正确位于body末尾
- ✅ **文档结构完整性**: 页面设置在保存/加载过程中保持完整
- ✅ **解析稳定性**: 即使XML结构有变化也能正确解析SectionProperties
- ✅ **Word兼容性**: 生成的文档完全符合Microsoft Word和WPS的要求

#### 测试验证
- 通过 `TestPageSettingsIntegration` 验证修复效果
- 使用 `TestDebugPageSettings` 进行详细调试验证
- 确认页面设置在完整的保存/加载周期中保持正确

---

## [v1.3.2] - 2025-06-02

### 🐛 问题修复

#### 模板引擎循环内条件表达式修复 ✨ **重要修复**
- **修复问题**: 解决了模板引擎中循环内部条件表达式无法正确渲染的问题
- **影响范围**: 主要影响使用复杂模板的场景，特别是 `{{#each}}` 循环内包含 `{{#if}}` 条件语句
- **错误表现**: 
  - 循环内的条件表达式保持原始模板语法，未被正确渲染
  - 例如：`{{#if isLeader}}👑 团队负责人{{/if}}` 在循环中不生效
- **修复方案**:
  - 优化 `renderLoopConditionals()` 函数的布尔值转换逻辑
  - 调整模板渲染顺序，先处理循环语句，再处理条件语句
  - 改进条件表达式的数据类型支持（字符串、数字、布尔值等）
- **修复后效果**:
  - 循环内条件表达式正确渲染：`{{#each teamMembers}}{{#if isLeader}}👑 团队负责人{{/if}}{{/each}}`
  - 支持多种数据类型的条件判断：`bool`, `string`, `int`, `int64`, `float64`
  - 完美支持嵌套的条件和循环结构

#### 技术细节
- **修改文件**: `pkg/document/template.go`
- **修改内容**: 
  - 优化 `renderLoopConditionals()` 函数的类型判断逻辑
  - 调整 `renderTemplate()` 中的渲染顺序
  - 简化 `renderConditionals()` 函数，移除不必要的循环检测
- **影响功能**: 
  - 模板引擎的循环内条件渲染
  - 复杂模板的嵌套结构处理
  - 所有使用条件表达式的模板功能

### 🔍 质量改进

#### 模板引擎稳定性
- ✅ **条件表达式完整支持**: 循环内外的条件表达式都能正确工作
- ✅ **数据类型兼容性**: 支持多种数据类型的条件判断
- ✅ **嵌套结构支持**: 完美支持条件语句和循环语句的任意嵌套
- ✅ **渲染顺序优化**: 确保模板元素按正确顺序处理

#### 测试验证
- 通过 `test_loop_condition.go` 验证修复效果
- 使用复杂模板演示验证嵌套结构
- 确认所有模板测试用例通过

---

## [v1.3.1] - 2025-05-30

### 🐛 问题修复

#### XML命名空间错误修复 ✨ **重要修复**
- **修复问题**: 解决了生成的XML文档中 `w15:color` 元素命名空间未声明的错误
- **影响范围**: 主要影响使用目录功能（TOC）的文档
- **错误表现**: 
  - XML linter 报错：`The prefix "w15" for element "w15:color" is not bound`
  - 生成的document.xml缺少 `xmlns:w15` 命名空间声明
- **修复方案**:
  - 在 `serializeDocument()` 方法中添加 `w15` 命名空间声明
  - 添加 `xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml"` 到文档根元素
- **修复后效果**:
  ```xml
  <w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" 
              xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml">
  ```

#### 技术细节
- **修改文件**: `pkg/document/document.go`
- **修改内容**: 在 `documentXML` 结构体中添加 `XmlnsW15` 字段
- **影响功能**: 
  - 目录生成（TOC）
  - 结构化文档标签（SDT）
  - 所有使用 `w15:color` 的功能

### 🔍 质量改进

#### XML 标准合规性
- ✅ **符合OOXML规范**: 生成的XML现在完全符合Office Open XML标准
- ✅ **命名空间完整性**: 所有使用的XML命名空间都正确声明
- ✅ **XML验证通过**: 生成的文档能够通过XML linter验证
- ✅ **Word兼容性**: 与Microsoft Word和WPS完全兼容

#### 测试验证
- 通过 `advanced_features` 示例验证修复效果
- 使用 `unzip_tool.exe` 解压验证XML结构
- 确认所有 `w15:` 前缀元素都有正确的命名空间绑定

---

## [v1.3.0] - 2025-01-18

### ✨ 重大功能新增

#### 页眉页脚功能 ✨ **全新实现**
- **完整的页眉页脚支持**: 实现了页眉、页脚的创建和管理
- **多种页眉页脚类型**: 
  - `Default` - 默认页眉页脚
  - `First` - 首页页眉页脚  
  - `Even` - 偶数页页眉页脚
- **页码显示功能**: 支持在页眉页脚中显示页码
- **关键API**:
  - `AddHeader()` - 添加页眉
  - `AddFooter()` - 添加页脚
  - `AddHeaderWithPageNumber()` - 添加带页码的页眉
  - `AddFooterWithPageNumber()` - 添加带页码的页脚
  - `SetDifferentFirstPage()` - 设置首页不同

#### 目录生成功能 ✨ **全新实现**
- **自动目录生成**: 基于标题样式自动创建目录
- **多级目录支持**: 支持1-9级标题的目录条目
- **目录配置选项**:
  - 目录标题自定义
  - 页码显示控制
  - 超链接支持
  - 点状引导线
- **书签集成**: 标题自动生成书签，支持导航
- **关键API**:
  - `GenerateTOC()` - 生成目录
  - `UpdateTOC()` - 更新目录
  - `AddHeadingWithBookmark()` - 添加带书签的标题
  - `AutoGenerateTOC()` - 自动生成目录
  - `GetHeadingCount()` - 获取标题统计

#### 脚注和尾注功能 ✨ **全新实现**
- **脚注管理**: 完整的脚注添加、删除和配置功能
- **尾注支持**: 文档末尾的尾注功能
- **多种编号格式**:
  - 十进制数字 (`decimal`)
  - 小写/大写罗马数字 (`lowerRoman`, `upperRoman`)
  - 小写/大写字母 (`lowerLetter`, `upperLetter`)
  - 符号编号 (`symbol`)
- **脚注位置控制**:
  - 页面底部 (`pageBottom`)
  - 文本下方 (`beneathText`)
  - 节末尾 (`sectEnd`)
  - 文档末尾 (`docEnd`)
- **编号重启规则**:
  - 连续编号 (`continuous`)
  - 每节重启 (`eachSect`)
  - 每页重启 (`eachPage`)
- **关键API**:
  - `AddFootnote()` - 添加脚注
  - `AddEndnote()` - 添加尾注
  - `SetFootnoteConfig()` - 设置脚注配置
  - `GetFootnoteCount()`, `GetEndnoteCount()` - 获取数量统计

#### 列表和编号功能 ✨ **全新实现**
- **无序列表**: 支持多种项目符号
  - 圆点符号 (`•`)
  - 空心圆 (`○`)
  - 方块 (`■`)
  - 短横线 (`–`)
  - 箭头 (`→`)
- **有序列表**: 支持多种编号格式
  - 十进制数字 (`decimal`)
  - 小写/大写字母 (`lowerLetter`, `upperLetter`)
  - 小写/大写罗马数字 (`lowerRoman`, `upperRoman`)
- **多级列表**: 支持最多9级嵌套
- **编号控制**: 支持重新开始编号
- **关键API**:
  - `AddListItem()` - 添加列表项
  - `AddBulletList()` - 添加无序列表
  - `AddNumberedList()` - 添加有序列表
  - `CreateMultiLevelList()` - 创建多级列表
  - `RestartNumbering()` - 重启编号

#### 结构化文档标签（SDT） ✨ **全新实现**
- **目录SDT结构**: 专门用于目录功能的SDT实现
- **SDT属性管理**: 完整的SDT属性和内容控制
- **文档部件支持**: 支持SDT占位符和文档部件
- **关键API**:
  - `CreateTOCSDT()` - 创建目录SDT结构

#### 域字段功能 ✨ **全新实现**
- **超链接域**: 支持文档内部超链接
- **页码引用域**: 支持页码引用和导航
- **域字符控制**: 完整的域开始、分隔、结束标记
- **关键API**:
  - `CreateHyperlinkField()` - 创建超链接域
  - `CreatePageRefField()` - 创建页码引用域

### 🏗️ 架构改进

#### 新增核心文件
- `header_footer.go` - 页眉页脚功能实现
- `toc.go` - 目录生成功能实现
- `footnotes.go` - 脚注尾注功能实现
- `numbering.go` - 列表编号功能实现
- `sdt.go` - 结构化文档标签实现
- `field.go` - 域字段功能实现

#### 文档完善
- 新增 `pkg/document/README.md` 详细API文档更新
- 增加了所有新功能的使用示例和配置说明
- 新增配置结构体文档说明

### 📚 示例程序

#### 新增示例目录
- `examples/page_settings/` - 页面设置演示
- `examples/advanced_features/` - 高级功能综合演示
  - 页眉页脚演示
  - 目录生成演示
  - 脚注尾注演示
  - 列表编号演示

### 🔧 配置结构体

#### 新增配置类型
- `TOCConfig` - 目录配置
- `FootnoteConfig` - 脚注配置
- `ListConfig` - 列表配置
- `HeaderFooterType` - 页眉页脚类型枚举
- `FootnoteNumberFormat` - 脚注编号格式枚举
- `ListType` - 列表类型枚举
- `BulletType` - 项目符号类型枚举

### 📝 使用示例更新

```go
// 页眉页脚示例
doc.AddHeader(document.HeaderFooterTypeDefault, "这是页眉")
doc.AddFooterWithPageNumber(document.HeaderFooterTypeDefault, "第", true)
doc.SetDifferentFirstPage(true)

// 目录示例
doc.AddHeadingWithBookmark("第一章 概述", 1, "chapter1")
tocConfig := document.DefaultTOCConfig()
doc.GenerateTOC(tocConfig)

// 脚注示例
doc.AddFootnote("正文内容", "脚注内容")
doc.AddEndnote("更多说明", "尾注内容")

// 列表示例
doc.AddBulletList("列表项1", 0, document.BulletTypeDot)
doc.AddNumberedList("第一项", 0, document.ListTypeDecimal)
```

### 🎯 兼容性保证

- ✅ **API向下兼容**: 所有现有API保持不变
- ✅ **无破坏性变更**: 现有代码无需修改
- ✅ **渐进增强**: 新功能作为可选功能提供

### 🔍 技术改进

#### 功能模块化
- 每个新功能独立文件实现，降低代码耦合
- 统一的错误处理和日志记录
- 符合Word OOXML标准的实现

#### 代码质量
- 完整的单元测试覆盖
- 详细的API文档和注释
- 规范的Go代码风格

---

## [v1.2.0] - 2025-05-29

### ✨ 新增功能

#### 表格默认样式改进
- **表格默认边框样式**: 新创建的表格现在默认包含单线边框样式，无需手动设置
- **参考标准格式**: 默认样式参考了 Word 标准表格格式（tmp_test 目录中的参考实现）
- **详细规格**:
  - 边框样式：`single`（单线）
  - 边框粗细：`4`（1/8磅单位）
  - 边框颜色：`auto`（自动）
  - 边框间距：`0`
  - 表格布局：`autofit`（自动调整）
  - 单元格边距：左右各 `108 dxa`

#### 功能特性
- ✅ **向下兼容**: 现有代码无需修改，自动享受新的默认样式
- ✅ **样式覆盖**: 仍然支持通过 `SetTableBorders()` 等方法自定义样式
- ✅ **无边框选项**: 可通过 `RemoveTableBorders()` 方法回到原来的无边框效果
- ✅ **标准匹配**: 与 Microsoft Word 创建的表格样式保持一致

### 🔧 改进内容

#### 代码改进
- 修改 `CreateTable()` 函数，在表格属性中增加默认边框配置
- 添加表格布局和单元格边距的默认设置
- 保持原有 API 接口不变，确保兼容性

#### 测试完善
- 新增 `TestTableDefaultStyle` 测试，验证默认样式正确应用
- 新增 `TestDefaultStyleMatchesTmpTest` 测试，确保与参考格式匹配
- 新增 `TestDefaultStyleOverride` 测试，验证样式覆盖功能

#### 示例程序
- 新增 `examples/table_default_style/` 演示程序
- 展示新默认样式、原无边框效果对比、自定义样式覆盖等功能

### 📝 文档更新

#### README.md
- 更新表格功能说明，增加默认样式特性描述
- 标注新增功能和改进点

#### pkg/document/README.md
- 更新 `CreateTable` 方法说明，增加默认样式信息

### 🎯 影响范围

#### 用户体验改进
- **即开即用**: 新创建的表格具有专业的外观，无需额外设置
- **标准化**: 确保表格样式与 Word 标准一致
- **灵活性**: 保持完整的自定义能力

#### 开发者友好
- **API 稳定**: 无破坏性变更，现有代码继续工作
- **渐进增强**: 新功能作为默认行为提供，不影响现有逻辑

### 🔍 技术细节

#### 参考实现
基于 `tmp_test/word/document.xml` 中的表格定义：
```xml
<w:tblBorders>
  <w:top w:val="single" w:color="auto" w:sz="4" w:space="0"/>
  <w:left w:val="single" w:color="auto" w:sz="4" w:space="0"/>
  <w:bottom w:val="single" w:color="auto" w:sz="4" w:space="0"/>
  <w:right w:val="single" w:color="auto" w:sz="4" w:space="0"/>
  <w:insideH w:val="single" w:color="auto" w:sz="4" w:space="0"/>
  <w:insideV w:val="single" w:color="auto" w:sz="4" w:space="0"/>
</w:tblBorders>
```

#### 实现位置
- 文件：`pkg/document/table.go`
- 函数：`CreateTable()`
- 影响：所有通过 `AddTable()` 创建的新表格

---

## [v1.1.0] - 2025-05-28

### 🎨 表格样式系统
- 完整的表格边框设置功能
- 表格和单元格背景颜色支持
- 多种边框样式（单线、双线、虚线、点线等）
- 奇偶行颜色交替功能

### 📐 表格布局功能
- 表格尺寸控制（宽度、高度、列宽）
- 表格对齐和定位
- 单元格合并功能
- 行高设置和分页控制

### 🎯 样式管理系统
- 18种预定义样式支持
- 样式继承机制
- 自定义样式创建
- 样式查询和批量操作 API

---

## [v1.0.0] - 2025-05-27

### 🚀 初始版本
- 基础文档创建和操作功能
- 文本格式化支持
- 段落管理和样式设置
- 基础表格创建和单元格操作 
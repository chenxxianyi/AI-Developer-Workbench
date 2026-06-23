# Project Doctor 可访问性与响应式优化设计

## 背景

`UI_REVIEW_REPORT (1).md` 指向 Project Doctor 页面，核心问题集中在表单可访问性、上传区域键盘可用性、分析深度选项语义、焦点可见性、状态表达过度依赖颜色，以及移动端横向拥挤风险。

## 目标

在不改变后端接口和页面主流程的前提下，提升 Project Doctor 页面的可访问性、交互一致性和小屏稳定性。

## 范围

仅修改 `frontend/src/pages/tools/ProjectDoctorPage.vue`，并新增对应组件测试。

- 为所有表单字段增加稳定 `id`，并通过 `label for` 关联控件。
- 为必填字段增加 `required` 与 `aria-describedby`。
- 上传 ZIP 区域补充 `role="button"`、`tabindex="0"`、`aria-describedby`，支持 `Enter` / `Space` 触发文件选择。
- 上传说明文案更明确，文件名在小屏下可截断。
- 分析深度改为可访问的单选组语义：容器使用 `role="radiogroup"`，选项使用 `role="radio"` 和 `aria-checked`。
- 所有主要交互控件使用清晰 `focus-visible` ring。
- 结果严重级别展示中文 badge：高 / 中 / 低，并保留颜色作为辅助而非唯一表达。
- 小屏下分析深度和结果行允许换行，避免横向拥挤。

## 非目标

- 不重做页面视觉风格。
- 不改全局 Header/Sidebar。
- 不新增拖拽上传逻辑。
- 不改变 Project Doctor API 的字段名和提交逻辑。

## 测试策略

新增 `frontend/src/pages/tools/ProjectDoctorPage.test.ts`：

- 验证标题、项目名称、上传文件、技术栈、项目描述字段可以通过 label 查找。
- 验证上传区具备按钮语义、可聚焦，并有键盘操作说明。
- 验证分析深度是 radiogroup，三个选项有 radio 语义和 `aria-checked`。
- 验证有结果数据时严重级别显示中文“高”。


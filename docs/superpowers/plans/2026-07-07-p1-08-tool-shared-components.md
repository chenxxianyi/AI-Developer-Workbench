# P1-08 抽取五工具通用前端组件 实施方案

> 状态：已审阅，进入实施。本方案与 `项目优化开发任务拆分.md` 的 P1-08 任务卡对应。
> 依赖：M1 已结案（P0-01～P0-10 完成）。不引入新前端依赖。

## 1. 目标与范围

消除五个工具页（UIReview / ProjectDoctor / AgentConfig / APIDoc / DBSchema）之间重复的页头、表单分区、文件上传、代码输入与结果渲染代码；让工具结果页与 `ReportDetailPage` 共用同一套报告组件。

**硬性约束**：现有 `frontend/src/pages/tools/UIReviewPage.test.ts` 与 `ProjectDoctorPage.test.ts` 的全部断言必须保持通过，包括：
- `data-testid`：`screenshot-upload-zone`、`review-mode-screenshot/code/screenshot-code`、`code-source-paste/project-zip`、`project-zip-upload-zone`、`project-zip-input`、`ui-review-submit`、`project-zip-upload-zone`（ProjectDoctor）
- `label[for="..."]`：`project-doctor-title/name/zip/tech-stack/description`
- `role="button"` + `tabindex="0"` + `aria-describedby` + 文案 `按 Enter 或空格选择文件`、`Ctrl+V`
- `[role="radiogroup"][aria-label="分析深度"]` + 3 个 `[role="radio"]`，默认 `标准` 被选中
- severity 本地化（`高` 而非 `high`）
- h1 文案、步骤胶囊文案、空结果引导卡文案（评分维度 / 问题优先级 / 改进建议）

## 2. 现状盘点（已读证据）

| 页面 | 输入方式 | 文件上传 | 代码输入 | 结果区现状 |
| --- | --- | --- | --- | --- |
| UIReviewPage | FormData | 截图(pasteable+preview) + 项目ZIP | textarea | 内联 scores/issues/recommendations/generated-files/export |
| ProjectDoctorPage | FormData | 项目ZIP | — | 内联同上 + severity 本地化 |
| AgentConfigPage | JSON | — | — | 内联同上 |
| APIDocPage | FormData(paste) / JSON | 项目ZIP(paste模式无) | textarea | 内联同上 |
| DBSchemaPage | JSON | — | — | 内联同上 |

重复结构（五页一致）：
- 外壳 `max-w-6xl mx-auto` + 头部 icon 圆角块 + h1 + 描述 + 右上步骤胶囊
- `lg:grid-cols-2` 两栏（输入面板 + 结果面板），面板样式 `bg-surface border border-border rounded-xl p-6 shadow-sm`
- 错误块 `bg-danger/10 border border-danger/20 rounded-lg` + AlertCircle
- 提交按钮 Loader2 旋转 + disabled + focus-visible
- 空结果引导卡 + loading 旋转
- 结果区 summary + scores + issues + recommendations + generated_files + export 按钮，与 `ReportDetailPage` 用到的报告组件同源但**完全重写**了一份

报告组件清单（`frontend/src/components/report/`，目前仅 `ReportDetailPage` + `ReportsPage` 使用）：
`ScorePanel` `IssueList` `RecommendationList` `CodexPromptBox` `GeneratedFilesPanel` `ReportErrorState` `ActionItemsPanel` `ActionItemCard` `ReportStatusBadge` `ReportListCard`。

公共组件现状（`frontend/src/components/common/`）：`ConfirmDialog` `LanguageSwitcher` `PaginationBar`。

## 3. 新建组件（`frontend/src/components/tool/`）

### 3.1 `ToolPageShell.vue`
统一外壳 + 两栏 + 头部 + 错误块 + 提交按钮 + 结果/空态/loading 容器。

- props：
  - `icon: Component`（页头圆角块图标）
  - `iconColor: string`（默认 `'accent'`，可选 `success`/`orange`/`teal`/`purple`）
  - `title: string`
  - `description: string`
  - `stepText: string`（右上胶囊，如 `填写输入 → 上传素材 → 开始分析 → 查看结果`）
  - `loading: boolean`
  - `error: string | null`
  - `canSubmit: boolean`
  - `submitLabel: string`（默认 `开始分析`）
  - `submittingLabel: string`（默认 `分析中...`）
  - `submitColor: string`（默认 `accent`）
- slots：
  - `#form`（左栏表单内容，不含提交按钮）
  - `#actions`（提交按钮左侧的额外按钮，如重置）
  - `#result`（右栏结果内容，页面自行决定渲染条件）
  - `#empty`（无结果且非 loading 时的引导卡）
- emits：`submit`
- 渲染：`max-w-6xl mx-auto` → 头部 → `lg:grid-cols-2` 两栏 → 左栏底部固定错误块 + 提交/额外按钮 → 右栏 loading 旋转 或 `#empty` 或 `#result`

### 3.2 `ToolFormSection.vue`
统一 `<label>` + 控件容器。

- props：`label: string`、`required?: boolean`、`optional?: boolean`、`help?: string`、`idFor?: string`、`helpId?: string`（aria-describedby 关联）
- slot：默认控件
- 渲染：`mb-4` + `<label :for="idFor" class="block text-sm font-medium text-text-secondary mb-2">` + 可选 `*` / `(可选)` + 可选 help 文案 + slot

### 3.3 `FileUpload.vue`
统一拖放/点击/键盘/移除/预览/paste。

- v-model：`file: File | null`
- props：
  - `accept: string`（input accept）
  - `label?: string`
  - `hint?: string`（如 `支持 PNG, JPG, WebP (最大 20MB)`）
  - `helpId?: string`（aria-describedby）
  - `testid?: string`（透传：dropzone = `${testid}-upload-zone`，input = `${testid}-input`；为兼容现有测试，也允许只传一个 id）
  - `pasteable?: boolean`（图片模式支持 Ctrl+V）
  - `preview?: boolean`（图片预览）
  - `id?: string`（input id，用于 label[for]）
  - `emptyText?: string`（dropzone 主文案，如 `上传前端项目 ZIP`）
  - `emptySubText?: string`（次文案）
- emits：`update:file`、`paste`(file)
- 行为：
  - dropzone `role="button"` `tabindex="0"`，`@click` 触发隐藏 input.click()，`@keydown.enter/space.prevent` 触发 click
  - `@drop.prevent` + `@dragover.prevent` 支持拖放
  - `pasteable` 时 `@paste` 监听 ClipboardEvent，取首个 image item，`preventDefault` 并 emit
  - 选中后：`preview` 显示 `<img>` + 右上移除按钮；否则显示文件名 + 移除按钮（aria-label）
  - `aria-describedby` 绑定 `helpId`
  - 移除按钮 emit `update:file = null`

### 3.4 `CodeInput.vue`
统一代码 textarea。

- v-model：`modelValue: string`
- props：`placeholder?: string`、`rows?: number`（默认 8）、`id?: string`
- 渲染：`<textarea>` + mono font + 现有 focus-visible 样式

## 4. 逐页迁移顺序（增量，每步跑测试）

1. **建组件 + 单测**（不碰工具页）
   - `ToolPageShell.test.ts`：loading 旋转、error 文案、canSubmit disabled、submit emit、#empty/#result 切换
   - `ToolFormSection.test.ts`：label[for] 关联、required/optional 文案、helpId 绑定
   - `FileUpload.test.ts`：点击触发 input.click、Enter/Space 键盘、drop、移除、pasteable paste、preview 渲染
   - `CodeInput.test.ts`：v-model 双向、placeholder
   - 跑 `npm run test:unit`

2. **DBSchemaPage**（最简单，JSON 无文件）
   - 用 `ToolPageShell` + `ToolFormSection` + 结果区改用 `ScorePanel`/`IssueList`/`RecommendationList`/`GeneratedFilesPanel`/`CodexPromptBox`
   - 验证 API 可行性

3. **AgentConfigPage**（JSON 无文件，非评分型）
   - 同上，结果区用 `GeneratedFilesPanel` + `CodexPromptBox` + `ActionItemsPanel`，非评分提示用 `ReportDetailPage` 同款文案

4. **APIDocPage**（FormData + 可选 project_zip）
   - 项目 ZIP 用 `FileUpload`，代码用 `CodeInput`

5. **ProjectDoctorPage**（FormData + ZIP，有 radiogroup 与本地化测试）
   - ZIP 用 `FileUpload`，**保留** `label[for="project-doctor-*"]`、`role=radiogroup` + 3 `role=radio`（默认标准）、`aria-describedby="project-doctor-zip-help"`、`按 Enter 或空格选择文件` 文案、severity 本地化
   - 结果区改用报告组件，若 `IssueList` 的 severity 显示不满足 `'高'` 断言，则保留页面本地化映射或增强 `IssueList`（优先增强组件，使两处一致）
   - 跑 `ProjectDoctorPage.test.ts`

6. **UIReviewPage**（最复杂：截图 pasteable + preview + 项目 ZIP + 代码 + 多模式）
   - 截图用 `FileUpload pasteable preview`，**保留** `data-testid="screenshot-upload-zone"`、`Ctrl+V` 文案
   - 项目 ZIP 用 `FileUpload`，**保留** `project-zip-upload-zone`、`project-zip-input`
   - 代码用 `CodeInput`
   - **保留** `review-mode-*` / `code-source-*` 按钮组、`ui-review-submit`、h1 `UI 质量审查`、步骤胶囊、空结果引导卡文案（评分维度/问题优先级/改进建议）
   - 跑 `UIReviewPage.test.ts`

## 5. 结果区统一为报告组件

5 页结果区把内联块替换为：
- `ScorePanel`（评分型工具，有 scores 时）
- `IssueList`（issues）
- `RecommendationList`（recommendations）
- `ActionItemsPanel`（action_items，与 ReportDetailPage 同源）
- `CodexPromptBox`（codex_prompt）
- `GeneratedFilesPanel`（generated_files，含下载按钮）
- 导出 Markdown 按钮：复用 `ReportDetailPage` 的 `downloadReport` 调用（`@/api/reports`），而非 `window.open`，保持下载行为一致

消除两套展示逻辑（工具页内联 vs 报告组件）。

## 6. 验收

- `cd frontend; npm run test:unit` 通过（旧测试 + 新组件测试）
- `cd frontend; npm run build` 通过
- `cd frontend; npm run test:e2e` 通过（`report-flow.spec.ts` 工具页可访问）
- 五页视觉与交互一致，工具特有字段不丢失
- 重复的下载/复制/提交状态代码显著减少

## 7. 文档更新

完成后将 `项目优化开发任务拆分.md` 中 `## [ ] P1-08` 改为 `## [x] P1-08`，补完成记录与遗留项，然后继续 P1-09。

## 8. 风险与回滚

- 风险：迁移过程中破坏现有测试锚点。缓解：每页迁移后立即跑该页测试 + 全量 unit/e2e。
- 回滚：每页一个可回滚提交，提交信息 `feat(P1-08): migrate <page> to shared tool components`。

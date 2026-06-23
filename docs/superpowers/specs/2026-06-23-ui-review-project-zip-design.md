# UI Review 前端项目 ZIP 审查设计

## 背景

UI Review 当前支持截图、粘贴代码、截图 + 粘贴代码。用户希望新增上传“前端整个文件压缩包”进行 UI 审查。项目中已有 Project Doctor 和 API Doc 的 ZIP 上传、安全校验、静态扫描能力，因此本功能应复用现有 ZIP 链路，不执行用户项目代码。

## 设计决策

上传前端 ZIP 不作为第四个主分析模式，而作为“代码来源”的一种：

- 粘贴代码：轻量、适合单文件或代码片段。
- 上传前端 ZIP：适合完整 Vue/React/HTML/CSS 项目。

主分析模式保持：

- `screenshot`
- `code`
- `screenshot_code`

新增字段：

- `code_source`: `paste` 或 `project_zip`
- `project_zip`: ZIP 文件

## 前端交互

仅在 `code` 和 `screenshot_code` 模式下显示“代码来源”选择。

- 选择“粘贴代码”：显示现有代码输入框。
- 选择“上传前端 ZIP”：显示 ZIP 上传区域。
- `code` 模式下，标题 + ZIP 即可提交。
- `screenshot_code` 模式下，必须有截图 + 代码来源；代码来源可以是粘贴代码或 ZIP。

ZIP 上传文案：

- 支持 Vue / React / HTML / CSS 等前端项目。
- 系统只做静态读取，不执行代码。
- 建议删除 `node_modules`、`dist`、`.git`。
- 最大 20MB。

## 后端行为

`POST /api/tools/ui-review/run` 继续使用 multipart/form-data。

处理流程：

1. 接收 `code_source` 和 `project_zip`。
2. 创建 processing report。
3. 保存截图（如有）。
4. 如果 `code_source=project_zip`，保存 ZIP 为 `project_zip` 资产。
5. 使用现有 `ZipService.ExtractAndAnalyze` 静态扫描 ZIP。
6. 将扫描结果转为 UI Review 的源码上下文。
7. 调用 UI Review Prompt。
8. 保存报告和 Markdown。

## 安全边界

第一版不执行用户上传的项目代码：

- 不运行 `npm install`
- 不运行 `npm run build`
- 不启动 dev server
- 不自动截图页面

ZIP 扫描复用现有安全策略：

- 跳过 `node_modules`、`.git`、`dist`、`build` 等目录。
- 跳过 `.env`、私钥和证书文件。
- 防 Zip Slip、symlink、异常压缩比和超大解压体积。
- 读取文本文件并脱敏。

## 校验规则

| review_mode | 合法输入 |
| --- | --- |
| `screenshot` | 必须有截图 |
| `code` | 必须有粘贴代码或前端 ZIP |
| `screenshot_code` | 必须有截图，并且必须有粘贴代码或前端 ZIP |

`code_source` 为空时默认为 `paste`，保持兼容。

## 测试策略

- 后端测试：
  - `code` 模式允许 `code_source=project_zip` + ZIP。
  - `screenshot_code` 模式要求截图 + 代码来源。
  - `project_zip` 模式会保存 ZIP、扫描 ZIP，并把扫描摘要放进 AI Prompt。
  - `code` 模式无代码无 ZIP 时拒绝。
- 前端测试：
  - `code` 模式显示代码来源选择。
  - 切换到上传 ZIP 后显示 ZIP 上传区。
  - 选择 ZIP 后可提交，并向 FormData 写入 `code_source=project_zip` 和 `project_zip`。


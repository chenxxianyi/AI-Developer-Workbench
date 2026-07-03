# AI Developer Workbench 优化方案

> 版本：v1.1  
> 日期：2026-07-03  
> 适用范围：个人本地、单用户使用  
> 目标：将当前五工具工作台升级为“AI 生成项目的质量控制台”，优先补齐个人使用场景下的产品闭环、技术地基和长期维护能力。

## 1. 项目定位建议

当前项目已经具备五个核心工具雏形：

- UI 质量审查
- 项目诊断
- Agent 配置生成
- API 文档生成
- 数据库结构审查

建议不要继续横向堆叠更多工具，而是把产品定位从“AI 工具集合”升级为：

```text
面向 AI Coding 项目的质量控制台
```

核心价值应从“分析一次”转向“持续改进一个项目”：

```text
上传/输入项目材料
  -> 生成质量报告
  -> 输出可执行整改清单
  -> 生成 Codex/Copilot 修复 Prompt
  -> 修复后复查
  -> 沉淀项目质量趋势
```

这样产品会更像一个研发质量工作流，而不是多个彼此独立的 AI 表单页面。

## 2. 当前主要问题

### 2.1 产品闭环不完整

后端已经具备报告列表、详情、删除、导出、生成文件下载能力，但前端 `ReportsPage.vue` 和 `ReportDetailPage.vue` 仍是占位页面。

这会导致用户完成分析后无法系统地查看、复盘、下载和管理历史报告。

### 2.2 报告偏“结果展示”，缺少“行动转化”

当前报告主要输出评分、问题和建议，但用户下一步要做什么还不够明确。

建议把报告升级为可执行工作台，提供：

- 整改任务清单
- 优先级排序
- 一键复制 AI 修复 Prompt
- 导出 GitHub Issue 文案
- 修复后复查入口
- 前后两次分析结果对比

### 2.3 Mock Mode 与实现不一致

文档和 `docker-compose.yml` 写了 `AI_MOCK_MODE=true`，但当前后端配置仍要求 `AI_API_KEY`，导致无 API Key 时无法按文档进入 Mock Mode。

这会影响本地开发、个人演示、自测验收和自动化测试。

### 2.4 Docker 部署闭环不完整

`docker-compose.yml` 期望构建前端镜像，但当前前端目录缺少 `Dockerfile`。这会导致“一键启动”体验不完整。

### 2.5 缺少项目级持续上下文

五个工具目前主要围绕单次请求工作，没有“项目档案”概念。

后续应支持：

- 项目名称
- 技术栈
- 仓库地址
- 编码规则
- UI 风格
- 数据库约束
- 历史报告
- Agent 配置文件
- API 文档
- 质量趋势

## 3. 产品功能优化方案

### 3.1 P0：补齐报告中心

目标：让所有工具的分析结果可沉淀、可管理、可复用。

建议功能：

- 报告列表
  - 分页
  - 按工具筛选
  - 按时间排序
  - 按评分排序
  - 状态筛选：成功、fallback、失败、处理中
- 报告详情
  - 报告标题、工具类型、创建时间
  - 总分、等级、摘要
  - 评分维度
  - 问题列表
  - 改进建议
  - Codex Prompt
  - 生成文件
  - Markdown 导出
- 报告删除
  - 删除前确认
  - 明确提示会删除关联上传文件和生成文件
- 失败与 fallback 状态展示
  - AI 失败、JSON 解析失败时仍能展示可用降级报告

验收标准：

- 五个工具运行后都可以进入报告详情。
- 历史报告可筛选、排序、分页。
- 生成文件可下载。
- Markdown 报告可导出。
- 删除报告后列表和详情状态正确更新。

### 3.2 P0：把报告升级为行动计划

目标：让报告不仅指出问题，还能指导下一步修复。

建议在每份报告中增加统一的 `action_items` 数据结构：

```json
[
  {
    "title": "修复移动端上传区域键盘不可用问题",
    "priority": "high",
    "effort": "small",
    "category": "accessibility",
    "reason": "当前上传区域依赖点击 div，键盘用户难以操作",
    "suggested_prompt": "请修复 UIReviewPage 中上传区域的键盘可访问性，要求支持 Enter/Space 触发文件选择，并补充 aria-label 与 focus-visible 样式。"
  }
]
```

前端展示：

- 按优先级分组
- 支持勾选完成
- 支持复制单条修复 Prompt
- 支持复制全部任务为 Markdown
- 支持导出 GitHub Issue 模板

验收标准：

- 每份报告至少输出 3 条行动项。
- 用户可以一键复制修复 Prompt。
- 用户可以下载整改清单。

### 3.3 P1：项目档案 Project Profile

目标：把五个工具从“单次分析”连接成“围绕项目的持续服务”。

建议新增项目实体：

```text
projects
  id
  name
  description
  repo_url
  frontend_stack
  backend_stack
  database
  ui_style
  coding_rules
  created_at
  updated_at
```

报告归属到项目：

```text
reports.project_id -> projects.id
```

项目详情页展示：

- 项目基础信息
- 最近报告
- 平均质量评分
- 各工具使用次数
- 风险趋势
- 最新 Agent 配置
- 最新 API 文档
- 最新 DB Schema 建议

验收标准：

- 用户可以创建项目档案。
- 运行工具时可以选择所属项目。
- 项目页可以看到该项目的全部历史报告。

### 3.4 P1：复查与前后对比

目标：形成质量改进闭环。

建议功能：

- 报告详情页提供“基于本报告复查”入口。
- 新报告记录 `parent_report_id`。
- 对比维度：
  - 总分变化
  - 高危问题数量变化
  - 中低危问题数量变化
  - 已解决问题
  - 新增问题
  - 建议变化

验收标准：

- 用户可以从旧报告发起复查。
- 两份同工具报告可以生成对比视图。
- Dashboard 展示最近质量趋势。

### 3.5 P1：GitHub 工作流增强

目标：让报告更容易进入真实研发流程。

建议功能：

- 导出 GitHub Issue Markdown
- 为每条行动项生成 Issue 标题和正文
- 为 Project Doctor 生成技术债清单
- 为 API Doc 生成文档补全任务
- 为 UI Review 生成 UI 修复任务

后续可扩展：

- GitHub OAuth
- 仓库导入
- PR Review
- 自动评论到 PR

MVP 阶段不建议直接做完整 GitHub OAuth，先做 Markdown 导出即可。

### 3.6 P2：多模型与规则模板

目标：提升个人高级使用场景的可控性。

建议功能：

- 模型配置预设
- 工具级 Prompt 模板
- 个人规则模板
- UI 审查标准模板
- 数据库审查标准模板
- Agent 配置模板

适合在产品稳定后做，因为它会引入配置复杂度。

## 4. 五个工具的专项优化

### 4.1 UI 质量审查

建议增强：

- 支持截图区域标注
- 支持移动端/桌面端双截图对比
- 增加可访问性专项评分
- 增加颜色对比度检查建议
- 输出具体组件级修复 Prompt
- 支持按页面类型套用审查标准

推荐输出：

- 总分
- 视觉层级评分
- 一致性评分
- 可访问性评分
- 响应式评分
- 高优先级 UI 问题
- 修复 Prompt
- 复查建议

### 4.2 项目诊断

建议增强：

- 增加 Agent Readiness 评分
- 检查是否有 README、AGENTS.md、测试命令、项目结构说明
- 检查依赖风险和过期依赖
- 检查配置文件是否完整
- 检查测试覆盖线索
- 检查 Docker/部署线索
- 输出技术债优先级

推荐评分维度：

- 结构清晰度
- 可维护性
- 可测试性
- 可部署性
- 文档完整度
- Agent 可接手程度

### 4.3 Agent 配置生成

建议增强：

- 支持多种 AI Coding 工具配置
  - `AGENTS.md`
  - `TASK_PLAN.md`
  - Copilot instructions
  - Codex 项目规则
  - Cursor/Windsurf 规则
- 支持根据已有项目 ZIP 生成配置
- 支持从 Project Doctor 报告反向生成项目规则
- 支持个人编码规范和项目规则模板

推荐生成文件：

- `AGENTS.md`
- `TASK_PLAN.md`
- `CODING_RULES.md`
- `FRONTEND_STYLE_GUIDE.md`
- `BACKEND_ARCHITECTURE.md`
- `REVIEW_CHECKLIST.md`

### 4.4 API 文档生成

建议增强：

- OpenAPI JSON 校验
- DTO/请求参数/响应结构提取
- 错误码提取
- 认证方式说明
- 前端调用示例
- curl 示例
- 按模块分组
- 文档缺口提示

推荐生成文件：

- `API_DOCUMENTATION.md`
- `openapi.json`
- `FRONTEND_API_CLIENT_GUIDE.md`

### 4.5 数据库结构审查

建议增强：

- ER 图导出
- 索引风险解释
- 慢查询风险推断
- 字段类型建议
- 迁移脚本风险等级
- 数据完整性检查
- 多数据库差异建议

推荐输出：

- 结构评分
- 索引评分
- 扩展性评分
- 数据完整性评分
- 优化后 Schema
- 迁移建议
- 执行风险提示

## 5. 技术架构优化方案

### 5.1 P0：实现真正的 Mock Mode

建议新增配置：

```go
type AIConfig struct {
    Provider       string
    BaseURL        string
    APIKey         string
    Model          string
    VisionModel    string
    MockMode       bool
    TimeoutSeconds int
    MaxRetries     int
}
```

配置规则：

- `AI_MOCK_MODE=true` 时强制使用 Mock。
- `AI_API_KEY` 为空时自动使用 Mock。
- Mock Mode 不访问外部 AI Provider。
- Settings 接口返回 `mock_mode`。

验收标准：

- 无 API Key 可以启动后端。
- 五个工具在 Mock Mode 下能完整生成报告和文件。
- Dashboard 和 Settings 明确显示当前为 Mock Mode。

### 5.2 P0：补齐 Docker 部署闭环

建议新增：

- `frontend/Dockerfile`
- 前端 nginx 静态部署配置
- 根目录 `.env.example`
- Compose 健康检查修正

验收标准：

```bash
docker compose up -d
```

启动后：

- `http://localhost:5173` 可访问前端。
- `http://localhost:8080/api/health` 返回健康状态。
- 无 API Key 时仍可运行 Mock 工具。

### 5.3 P0：报告保存事务化

当前风险：

- 报告状态更新成功后，生成文件保存失败只记录 warn。
- 用户可能看到成功报告，但下载不到文件。

建议：

- `SucceedReport` 中报告更新和 generated files 保存放入同一事务。
- 文件保存失败时整体返回错误或进入 fallback。
- 保证报告状态和文件状态一致。

验收标准：

- generated files 写入失败时不会产生“成功但文件缺失”的报告。
- 单元测试覆盖成功、文件保存失败、事务回滚。

### 5.4 P1：工具执行异步化

当前同步请求适合 MVP，但大 ZIP、视觉模型和慢模型容易超过请求超时。

建议新增：

```text
jobs
  id
  tool_type
  report_id
  status
  progress
  error_message
  created_at
  updated_at
```

接口：

```text
POST /api/tools/{tool}/run -> 返回 job_id
GET /api/jobs/:id -> 查询状态
POST /api/jobs/:id/cancel -> 取消任务
```

前端：

- 提交后进入任务进度页或右侧进度面板。
- 支持刷新。
- 支持失败重试。

验收标准：

- 长任务不会占用前端单个请求等待。
- 用户可以看到任务状态。
- 任务失败后有明确错误信息。

### 5.5 P1：AI 质量与成本观测

建议记录：

- 工具类型
- Provider
- Model
- 是否 Mock
- 请求耗时
- 是否重试
- 是否 fallback
- JSON 解析是否成功
- 错误类型

可新增表：

```text
ai_runs
  id
  report_id
  tool_type
  provider
  model
  is_mock
  duration_ms
  retry_count
  parse_success
  fallback_used
  error_type
  created_at
```

价值：

- 知道哪个工具最容易失败。
- 知道哪个模型响应慢。
- 知道 JSON 解析失败率。
- 为模型选择、稳定性优化和个人调用成本分析提供依据。

### 5.6 P1：安全增强

建议继续强化：

- Gin 层限制请求体总大小。
- Markdown 渲染禁用或净化 HTML。
- 下载文件名严格清洗。
- ZIP 读取预算与压缩比测试覆盖。
- Prompt Injection 回归测试。
- 上传项目中的 `.env`、密钥、证书永不进入 Prompt。
- CORS 生产环境禁止通配符。
- 无鉴权情况下明确禁止公网裸露部署。

推荐测试用例：

- ZIP Slip
- ZIP Bomb
- 伪装扩展名
- 超大图片
- 含密钥项目
- Markdown XSS
- Prompt Injection 文本

### 5.7 P1：前端组件抽象

当前五个工具页有重复模式。建议沉淀通用组件：

- `ToolPageShell.vue`
- `ToolFormSection.vue`
- `FileUpload.vue`
- `CodeInput.vue`
- `ScorePanel.vue`
- `IssueList.vue`
- `RecommendationList.vue`
- `GeneratedFilesPanel.vue`
- `CodexPromptBox.vue`
- `ReportStatusBadge.vue`

收益：

- 降低五个工具页维护成本。
- 统一交互体验。
- 更容易做移动端适配和可访问性。

## 6. 用户体验优化方案

### 6.1 信息架构

建议主导航调整为：

```text
工作台
项目
报告
工具
设置
```

工具页可继续保留五个入口，但长期应围绕项目组织内容。

### 6.2 Dashboard 优化

建议 Dashboard 展示：

- 最近项目
- 最近报告
- 本周分析次数
- 平均质量分
- 高危问题数量
- 最常用工具
- 快捷入口
- Mock/真实 AI 状态

### 6.3 报告详情优化

建议报告详情优先展示：

```text
结论摘要
  -> 总分与等级
  -> 最严重的 3 个问题
  -> 建议先做的 3 个行动
  -> 复制修复 Prompt
  -> 下载/导出
  -> 完整评分与细节
```

避免让用户先面对大量卡片和长文本。

### 6.4 表单体验

建议：

- 字段级错误提示
- 自动聚焦首个错误字段
- 上传区域支持键盘操作
- 文件大小和类型提示
- 提交失败保留输入
- 长任务展示明确进度文案
- 所有输入控件有显式 label
- 所有图标按钮有 `aria-label`

## 7. 推荐开发路线图

### 阶段 1：产品闭环修复

目标：让当前 MVP 真正可用。

任务：

- 实现前端报告列表页
- 实现前端报告详情页
- 实现 Mock Mode
- 补齐 Docker 部署
- 修复 Settings 中 Mock 状态展示
- 修复报告成功与生成文件保存的一致性

完成标准：

- 无 API Key 可启动并体验五个工具。
- 每个工具都能生成报告。
- 报告可查看、下载、删除。
- Docker Compose 可一键启动。

### 阶段 2：行动计划工作台

目标：让报告能驱动修复。

任务：

- 后端结果 DTO 增加 action_items
- Prompt Builder 要求输出整改清单
- 前端展示行动项
- 支持复制单条 Prompt
- 支持导出 GitHub Issue Markdown
- 支持复查入口

完成标准：

- 报告详情页能直接指导用户下一步开发。
- 用户可以把行动项复制给 Codex/Copilot 执行。

### 阶段 3：项目档案

目标：从单次工具升级为项目级质量管理。

任务：

- 新增 projects 表
- 新增项目 CRUD API
- 新增项目列表和详情页
- 报告关联项目
- Dashboard 展示项目维度统计

完成标准：

- 用户可以围绕一个项目持续运行五个工具。
- 项目页能看到质量趋势和历史报告。

### 阶段 4：异步任务与观测

目标：提升稳定性和可维护性。

任务：

- 新增 jobs 表
- 工具执行异步化
- 增加任务进度查询
- 增加 AI run 观测表
- Dashboard 增加失败率和耗时统计

完成标准：

- 大文件和慢模型不会导致前端长时间阻塞。
- 使用者可以看到工具成功率和失败原因。

## 8. 优先级矩阵

| 优先级 | 项目 | 价值 | 原因 |
| --- | --- | --- | --- |
| P0 | 报告列表和详情 | 极高 | 当前产品闭环断点 |
| P0 | Mock Mode | 极高 | 影响演示、测试、部署 |
| P0 | Docker 部署闭环 | 高 | 影响本地部署和环境复现 |
| P0 | 报告与生成文件事务一致性 | 高 | 避免数据状态不一致 |
| P1 | 行动项和修复 Prompt | 极高 | 从报告展示变成研发工作流 |
| P1 | Project Profile | 高 | 建立持续项目上下文 |
| P1 | 异步任务 | 高 | 提升稳定性 |
| P1 | AI 观测 | 中高 | 为质量和成本优化打基础 |
| P2 | GitHub 集成 | 中高 | 衔接个人 GitHub 研发流程 |

## 9. 技术验收清单

- [ ] 无 `AI_API_KEY` 时后端可启动。
- [ ] `AI_MOCK_MODE=true` 时不会访问外部 AI。
- [ ] 五个工具在 Mock Mode 下可完整运行。
- [ ] 报告列表、详情、删除、导出可用。
- [ ] 生成文件可下载。
- [ ] Docker Compose 可一键启动。
- [ ] 上传文件大小、类型、文件头均被校验。
- [ ] ZIP Slip 和 ZIP Bomb 测试通过。
- [ ] Markdown 渲染不存在明显 XSS 风险。
- [ ] AI 非法 JSON 时生成 fallback 报告。
- [ ] 报告和生成文件保存具备事务一致性。
- [ ] 前端移动端无关键布局溢出。
- [ ] 关键页面具备加载、错误和空状态。

## 10. 产品验收清单

- [ ] 用户可以从 Dashboard 选择任意工具。
- [ ] 用户可以运行五个工具并得到报告。
- [ ] 用户可以查看历史报告。
- [ ] 用户可以下载 Markdown 报告。
- [ ] 用户可以下载生成文件。
- [ ] 用户可以复制 Codex Prompt。
- [ ] 用户知道每个报告下一步该做什么。
- [ ] 用户可以删除不需要的报告。
- [ ] 用户可以看到当前系统是否处于 Mock Mode。
- [ ] 用户无需阅读文档即可完成一次完整分析闭环。

## 11. 建议近期执行顺序

推荐按以下顺序推进：

1. 修复 Mock Mode。
2. 补齐前端报告列表页。
3. 补齐前端报告详情页。
4. 补齐 Docker 前端镜像。
5. 修复报告和生成文件事务一致性。
6. 在五个工具结果中加入行动项。
7. 报告详情页加入复制 Prompt 和导出 Issue 文案。
8. 设计 Project Profile 数据模型。
9. 增加项目列表和项目详情。
10. 再考虑异步任务和 GitHub 集成。

## 12. 总结

这个项目已经具备很好的基础：工具方向清晰，前后端架构明确，报告模型、文件上传、ZIP 安全、AI 调用和五个工具服务都已经有雏形。

下一步不建议继续增加新工具，而应优先补齐产品闭环：

```text
工具运行
  -> 报告沉淀
  -> 行动计划
  -> AI 修复 Prompt
  -> 复查对比
  -> 项目趋势
```

当这条链路跑通后，AI Developer Workbench 会从“开发者工具箱”升级为“AI Coding 项目的质量管理平台”。这才是它最有价值、也最容易形成差异化的位置。

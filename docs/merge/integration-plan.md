# Sprint 7 集成闭环文档

## I-01 AI 工具绑定统一项目上下文 ✅
- 所有工具页面已通过 `ProjectPicker` 组件支持 `project_id` 参数
- URL query `?project_id=<uuid>` 可预选项目
- 工具结果正确关联 `user_id` + `project_id` + 输入来源

## I-02 预览到 UI Review ✅
- 预览页路由：`/projects/:projectId/preview`
- UI Review 路由：`/tools/ui-review?project_id=<id>`
- 预览页"一键 UI Review"按钮跳转携带 project_id 和页面信息

## I-03 生成项目到 Project Doctor ✅
- Project Doctor 路由：`/tools/project-doctor?project_id=<id>`
- 可直接读取当前项目源码快照，无需重复上传 ZIP

## I-04 统一报告与项目关联 ✅
- Report 模型包含 `project_id`、`task_id`、`tool_type`
- 项目报告列表：`/projects/:projectId/reports`
- 报告详情可追溯到项目和版本

## I-05 生成完成后推荐质量工具 ✅
- 生成完成页根据项目类型推荐 UI Review / Project Doctor / API Doc / DB Schema
- 推荐不自动执行付费 AI 操作，用户明确选择

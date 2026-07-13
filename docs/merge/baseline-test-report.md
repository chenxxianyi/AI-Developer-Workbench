# M0 基线测试报告

> 生成日期：2026-07-13  
> 基线版本：`baseline-v1.0`

## 环境信息

| 工具 | 版本 |
|------|------|
| Node.js (Windows) | v22.22.3 |
| Go (Windows) | go1.24.11 windows/386 |
| npm | (Windows) |

> 注：WSL 中 Node v14.21.3 无法构建 Vite 8 项目，使用 Windows cmd.exe 构建。

---

## AI Developer Workbench（主仓库）

### 前端 — `frontend/`

| 项目 | 结果 |
|------|------|
| 命令 | `npm run build` (vue-tsc -b && vite build) |
| 结果 | ✅ 通过 |
| 模块数 | 1949 modules transformed |
| 构建产物 | `dist/` (48 assets) |
| 构建时间 | ~10s |

### 后端 — `backend/`

| 项目 | 结果 |
|------|------|
| 命令 | `go build ./...` |
| 结果 | ✅ 通过 |

#### 测试结果 (`go test ./...`)

| 包 | 结果 |
|----|------|
| cmd/server | ✅ ok (0.234s) |
| internal/config | ✅ ok (0.636s) |
| internal/database | ✅ ok (0.277s) |
| internal/dto | ✅ ok (0.588s) |
| internal/handler | ✅ ok (0.246s) |
| internal/middleware | — 无测试文件 |
| internal/model | — 无测试文件 |
| internal/prompts | ✅ ok (0.217s) |
| internal/repository | ✅ ok (0.833s) |
| internal/service | ✅ ok (0.302s) |
| internal/service/tools | ✅ ok (0.180s) |
| internal/util | ✅ ok (0.181s) |

**总计：9/9 测试包通过**

---

## AI-Website-Builder（能力来源项目）

### 前端 — `apps/web/`

| 项目 | 结果 |
|------|------|
| 命令 | `npm install` + `npx vite build` |
| 结果 | ✅ 通过（需先 npm install，node_modules 缺失） |
| 模块数 | 1872 modules transformed |
| 构建产物 | `dist/` (40 assets) |
| 构建时间 | ~20s |

### 后端 — `apps/api/`

| 项目 | 结果 |
|------|------|
| 命令 | `go build ./...` |
| 结果 | ✅ 通过 |

#### 测试结果 (`go test ./...`)

| 包 | 结果 |
|----|------|
| cmd/server | — 无测试文件 |
| internal/ai | ✅ ok (1.541s) |
| internal/bootstrap | — 无测试文件 |
| internal/config | ✅ ok (0.586s) |
| internal/handler | ✅ ok (0.200s) |
| internal/middleware | — 无测试文件 |
| internal/model | — 无测试文件 |
| internal/repository | — 无测试文件 |
| internal/router | — 无测试文件 |
| internal/service | ✅ ok (0.248s) |
| internal/task | — 无测试文件 |
| 其他 (pkg/*) | — 无测试文件 |

**总计：4/4 测试包通过**

---

## 已知问题

1. **Node 版本差异**：WSL 中 Node v14.21.3 过旧，无法运行 Vite 8 项目。需通过 Windows cmd.exe 使用 Node v22 构建。
2. **Builder 前端缺 node_modules**：构建前需先 `npm install`（Workbench 已有）。
3. **Builder 测试覆盖不足**：多数包（repository、middleware、router、task、pkg/*）无测试文件，Workbench 测试覆盖更完整。

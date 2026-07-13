# AI Developer Workbench — 前端

Vue 3 + TypeScript + Pinia + Vite，为五工具工作台提供界面。

## 开发

```bash
cd frontend
npm install
npm run dev          # 本地开发，默认 http://localhost:5173
npm run build        # 生产构建到 dist/
npm run test:unit    # Vitest 单元测试
npm run test:e2e     # Playwright E2E（mock /api，无需后端）
```

## 与后端联调

默认通过 Vite 代理把 `/api` 转发到 `http://localhost:8080`（见 `vite.config.ts`）。
后端默认端口 8080，Mock 模式下无需 API Key 即可体验全部功能。

## Docker

前端镜像使用 Node 构建、Nginx 运行的多阶段镜像（见 `Dockerfile`）。
`nginx.conf` 提供 Vue Router history fallback 与 `/api` 反向代理到 `backend:8080`。

```bash
# 在仓库根目录
docker compose up -d --build        # 启动 MySQL + backend + frontend
docker compose logs -f frontend     # 查看前端日志
docker compose down                 # 停止
docker compose down -v              # 停止并清理数据卷
```

访问 http://localhost:5173，刷新任意路由不会返回 Nginx 404。

## 切换真实 AI

编辑仓库根目录 `.env`（从 `.env.example` 复制），设置：

```
AI_MOCK_MODE=false
AI_API_KEY=sk-...
AI_PROVIDER=openai
AI_BASE_URL=https://api.openai.com/v1
AI_MODEL=gpt-4.1
```

然后 `docker compose up -d` 即可。

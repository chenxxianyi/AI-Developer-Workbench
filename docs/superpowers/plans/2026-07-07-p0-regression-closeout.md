# P0 Regression Closeout Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Finish the remaining v0.2/P0 regression gate so the project can move into P1 work with passing baseline checks.

**Architecture:** Keep the current Vue 3 + Pinia + Gin structure. Fix the test harness first, then expand the E2E smoke coverage around existing P0 flows without introducing new runtime behavior.

**Tech Stack:** Vue 3, TypeScript, Vitest, Playwright, Gin, GORM, Go test.

---

### Task 1: Unit Test Harness

**Files:**
- Modify: `frontend/vite.config.ts`
- Verify: `frontend/src/stores/languageStore.test.ts`

- [ ] **Step 1: Configure Vitest for browser-like unit tests**

Add a `test` block to `frontend/vite.config.ts`:

```ts
test: {
  environment: 'jsdom',
  exclude: ['tests/e2e/**', 'node_modules/**', 'dist/**'],
}
```

- [ ] **Step 2: Run unit tests**

Run: `cd frontend; npm.cmd run test:unit`

Expected: Vitest does not load Playwright specs and `localStorage` is available.

### Task 2: P0 Regression E2E Coverage

**Files:**
- Modify: `frontend/tests/e2e/report-flow.spec.ts`

- [ ] **Step 1: Mock stable backend responses in Playwright**

Intercept `/api/health`, `/api/system/status`, `/api/dashboard/stats`, `/api/reports`, `/api/reports/:id`, export, generated-file and delete endpoints.

- [ ] **Step 2: Cover P0 user paths**

Check Dashboard and Settings mock badges, report list filters, detail hierarchy, Markdown export, generated-file download link and delete confirmation cancel path.

- [ ] **Step 3: Run E2E**

Run: `cd frontend; npm.cmd run test:e2e`

Expected: Chromium project passes without requiring a real backend.

### Task 3: Documentation Status Update

**Files:**
- Modify: `项目优化开发任务拆分.md`

- [ ] **Step 1: Update P0-10 only after verification passes**

Change `## [ ] P0-10 建立 v0.2 发布回归门禁` to `## [x] P0-10 建立 v0.2 发布回归门禁` only after backend tests, frontend unit tests, frontend build and E2E pass or have a clearly documented environment exception.

- [ ] **Step 2: Record verification notes**

Add a short completion note under P0-10 listing the commands run and any environment limitation, such as Docker CLI being unavailable.

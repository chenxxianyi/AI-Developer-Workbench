# UI Review Focused Flow Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Improve the UI Review page with the approved focused task-flow design.

**Architecture:** Keep all behavior in `frontend/src/pages/tools/UIReviewPage.vue`. Use existing Vue state and Tailwind token classes; no new dependencies or shared components are required.

**Tech Stack:** Vue 3 `<script setup>`, TypeScript, Tailwind CSS v4 tokens, Vitest, Vue Test Utils, jsdom.

---

### Task 1: Lock expected UI behavior with tests

**Files:**
- Modify: `frontend/src/pages/tools/UIReviewPage.test.ts`

- [ ] **Step 1: Add focused-flow assertions**

Add a second test that mounts `UIReviewPage.vue` and asserts:

```ts
expect(wrapper.get('h1').text()).toBe('UI 质量审查')
expect(wrapper.text()).toContain('填写输入 → 上传素材 → 开始分析 → 查看结果')
expect(wrapper.text()).not.toContain('返回工作台')
expect(wrapper.text()).toContain('仅上传截图，快速评估视觉质量')
expect(wrapper.text()).toContain('仅粘贴前端代码，审查结构与样式')
expect(wrapper.text()).toContain('截图 + 代码，获得更完整建议')
expect(wrapper.text()).toContain('将输出哪些内容')
expect(wrapper.text()).toContain('评分维度')
expect(wrapper.text()).toContain('问题优先级')
expect(wrapper.text()).toContain('改进建议')
```

- [ ] **Step 2: Run test to verify it fails**

Run: `npm.cmd exec vitest run src/pages/tools/UIReviewPage.test.ts -- --environment jsdom`

Expected: FAIL because the current UI still has the duplicate page-level return button and lacks the new focused-flow copy.

### Task 2: Implement focused-flow layout copy and hierarchy

**Files:**
- Modify: `frontend/src/pages/tools/UIReviewPage.vue`
- Test: `frontend/src/pages/tools/UIReviewPage.test.ts`

- [ ] **Step 1: Remove duplicate page return button**

Remove the page-level `返回工作台` button and its `router.push('/dashboard')` click handler usage from the template. Remove `useRouter`/`ArrowLeft` imports if no longer used.

- [ ] **Step 2: Upgrade page header**

Keep one H1 `UI 质量审查`, add subtitle copy `上传截图或前端代码，快速获得视觉层级、一致性、可访问性与改进建议。`, and add a small flow chip `填写输入 → 上传素材 → 开始分析 → 查看结果`.

- [ ] **Step 3: Rename input panel and add mode help text**

Change card title to `1. 输入与素材`. Add helper text under the mode label: `选择最贴近当前材料的分析方式。`. Expand each mode button to show title and description.

- [ ] **Step 4: Improve upload and form contrast**

Use stronger border/text classes on upload and inputs, such as `border-border/80`, `text-text-primary`, and visible `focus-visible:ring-2 focus-visible:ring-accent`.

- [ ] **Step 5: Improve result empty state**

Replace the empty state with a card explaining `将输出哪些内容` and chips/cards for `评分维度`, `问题优先级`, `改进建议`, `报告文件`.

- [ ] **Step 6: Emphasize submit action**

Move submit controls into a clearer footer row, make the submit button `px-7 py-3 font-semibold shadow-sm`, and add disabled-state helper text `请先填写标题并提供截图或代码。`.

- [ ] **Step 7: Run focused tests**

Run: `npm.cmd exec vitest run src/pages/tools/UIReviewPage.test.ts -- --environment jsdom`

Expected: PASS.

### Task 3: Verify production build

**Files:**
- Modified files from tasks above

- [ ] **Step 1: Run frontend build**

Run: `npm.cmd run build`

Expected: TypeScript and Vite build exit with code 0.

- [ ] **Step 2: Review diff**

Run: `git diff -- frontend/src/pages/tools/UIReviewPage.vue frontend/src/pages/tools/UIReviewPage.test.ts docs/superpowers/specs/2026-06-23-ui-review-focused-flow-design.md docs/superpowers/plans/2026-06-23-ui-review-focused-flow.md`

Expected: Diff only contains the approved focused-flow UI optimization and its tests/docs.


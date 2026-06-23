# Project Doctor Accessibility Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Improve Project Doctor page accessibility, semantic controls, and mobile robustness based on the review report.

**Architecture:** Keep all changes local to `frontend/src/pages/tools/ProjectDoctorPage.vue`. Use native Vue state and existing Tailwind tokens; add component tests with Vue Test Utils.

**Tech Stack:** Vue 3 `<script setup>`, TypeScript, Tailwind CSS v4 tokens, Vitest, Vue Test Utils, jsdom.

---

### Task 1: Add failing accessibility tests

**Files:**
- Create: `frontend/src/pages/tools/ProjectDoctorPage.test.ts`

- [ ] **Step 1: Write tests**

Create tests that mount `ProjectDoctorPage.vue` and assert:

```ts
expect(wrapper.get('label[for="project-doctor-title"]').text()).toContain('标题')
expect(wrapper.get('#project-doctor-title').attributes('required')).toBeDefined()
expect(wrapper.get('label[for="project-doctor-zip"]').text()).toContain('项目 ZIP 文件')
expect(wrapper.get('[data-testid="project-zip-upload-zone"]').attributes('role')).toBe('button')
expect(wrapper.get('[data-testid="project-zip-upload-zone"]').attributes('tabindex')).toBe('0')
expect(wrapper.get('[role="radiogroup"]').attributes('aria-label')).toBe('分析深度')
expect(wrapper.findAll('[role="radio"]')).toHaveLength(3)
expect(wrapper.text()).toContain('按 Enter 或空格选择文件')
```

Add a result-rendering test that sets `wrapper.vm.result` and expects `高` to render for a `high` issue.

- [ ] **Step 2: Run tests to verify failure**

Run: `npm.cmd exec vitest run src/pages/tools/ProjectDoctorPage.test.ts -- --environment jsdom`

Expected: FAIL because the current page lacks the tested semantic attributes and Chinese severity display.

### Task 2: Implement semantic form controls and upload zone

**Files:**
- Modify: `frontend/src/pages/tools/ProjectDoctorPage.vue`
- Test: `frontend/src/pages/tools/ProjectDoctorPage.test.ts`

- [ ] **Step 1: Add field IDs and label binding**

Add IDs:

- `project-doctor-title`
- `project-doctor-name`
- `project-doctor-zip`
- `project-doctor-tech-stack`
- `project-doctor-description`

Bind labels with `for`, set `required` on title and file input, and add `aria-describedby` for helper text.

- [ ] **Step 2: Add keyboard upload activation**

Add `triggerFileInput()` and `handleUploadZoneKeydown(event: KeyboardEvent)`. Trigger file input on Enter or Space and prevent default for Space.

- [ ] **Step 3: Upgrade upload zone markup**

Add `data-testid="project-zip-upload-zone"`, `role="button"`, `tabindex="0"`, `aria-describedby="project-doctor-zip-help"`, `@keydown="handleUploadZoneKeydown"`, and clear focus-visible styles.

### Task 3: Implement analysis depth radio semantics and severity labels

**Files:**
- Modify: `frontend/src/pages/tools/ProjectDoctorPage.vue`
- Test: `frontend/src/pages/tools/ProjectDoctorPage.test.ts`

- [ ] **Step 1: Add depth option metadata**

Add `analysisDepthOptions` array with `value`, `label`, and `description` for basic/standard/deep.

- [ ] **Step 2: Render accessible radiogroup**

Replace three repeated depth buttons with `v-for` buttons using `role="radio"`, `aria-checked`, and `@click`.

- [ ] **Step 3: Add Chinese severity helper**

Add `getSeverityDisplayName(severity: string)` returning `高`, `中`, `低`, with fallback to the original severity.

- [ ] **Step 4: Improve result row wrapping**

Use `flex-col sm:flex-row` on score and issue heading rows; render severity in a rounded badge with text.

### Task 4: Verify

**Files:**
- Modified files from tasks above

- [ ] **Step 1: Run focused Project Doctor test**

Run: `npm.cmd exec vitest run src/pages/tools/ProjectDoctorPage.test.ts -- --environment jsdom`

Expected: PASS.

- [ ] **Step 2: Run frontend build**

Run: `npm.cmd run build`

Expected: TypeScript and Vite build exit with code 0.


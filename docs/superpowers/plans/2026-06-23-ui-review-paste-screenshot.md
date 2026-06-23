# UI Review Paste Screenshot Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add focused `Ctrl+V` screenshot paste support to the UI Review screenshot upload area.

**Architecture:** Keep paste handling local to `UIReviewPage.vue`. Reuse one shared file-loading function for both file input selection and pasted clipboard images so preview and selected file state stay consistent.

**Tech Stack:** Vue 3 `<script setup>`, TypeScript, Vitest, Vue Test Utils, jsdom.

---

### Task 1: Component paste behavior test

**Files:**
- Create: `frontend/src/pages/tools/UIReviewPage.test.ts`
- Modify: none

- [ ] **Step 1: Write the failing test**

Create a jsdom component test that mounts `UIReviewPage.vue`, stubs router/API dependencies, dispatches a paste event with an image clipboard item on the upload zone, and expects an image preview to render.

- [ ] **Step 2: Run test to verify it fails**

Run: `npm.cmd exec vitest run src/pages/tools/UIReviewPage.test.ts -- --environment jsdom`

Expected before implementation: FAIL because the upload zone does not handle paste and the preview image is not rendered.

### Task 2: Paste support implementation

**Files:**
- Modify: `frontend/src/pages/tools/UIReviewPage.vue`
- Test: `frontend/src/pages/tools/UIReviewPage.test.ts`

- [ ] **Step 1: Extract shared file preview loading**

Move the duplicated file assignment and `FileReader.readAsDataURL` logic into `setScreenshotFile(file: File)`.

- [ ] **Step 2: Add local paste handler**

Add `handleScreenshotPaste(event: ClipboardEvent)` that reads `event.clipboardData?.items`, finds the first `image/*` item, calls `item.getAsFile()`, prevents default behavior, and passes the file to `setScreenshotFile`.

- [ ] **Step 3: Make upload zone focusable and bind paste**

Add `tabindex="0"`, `role="button"`, `@paste="handleScreenshotPaste"`, and `focus:border-accent focus:outline-none` to the upload zone. Update the visible text to mention `Ctrl+V`.

- [ ] **Step 4: Run focused test**

Run: `npm.cmd exec vitest run src/pages/tools/UIReviewPage.test.ts -- --environment jsdom`

Expected after implementation: PASS.

### Task 3: Regression checks

**Files:**
- Modified files from Tasks 1-2

- [ ] **Step 1: Run frontend build**

Run: `npm.cmd run build`

Expected: TypeScript and Vite build exit with code 0.

- [ ] **Step 2: Review git diff**

Run: `git diff -- frontend/src/pages/tools/UIReviewPage.vue frontend/src/pages/tools/UIReviewPage.test.ts docs/superpowers/specs/2026-06-23-ui-review-paste-screenshot-design.md docs/superpowers/plans/2026-06-23-ui-review-paste-screenshot.md`

Expected: Diff only contains the paste screenshot feature, tests, and docs.


# UI Review Project ZIP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add frontend project ZIP upload as a code source for UI Review while keeping screenshot/code/screenshot+code modes.

**Architecture:** Reuse existing multipart UI Review endpoint. Add `code_source` and `project_zip` to UI Review input, inject `ZipService` into `UIReviewService`, and convert scanned project summaries into prompt context.

**Tech Stack:** Vue 3, TypeScript, Go, Gin, existing FileService/ZipService, Vitest, Go tests.

---

### Task 1: Backend ZIP input tests

**Files:**
- Modify: `backend/internal/service/tools/ui_review_service_test.go`

- [ ] **Step 1: Add service fakes**

Create small fakes for AI, report, file, and zip services inside the test file. The ZIP fake should return a `ProjectSummary` containing `src/App.vue`.

- [ ] **Step 2: Add test for project ZIP code mode**

Test that `UIReviewService.Run` accepts:

```go
UIReviewFormInput{
  Title: "完整项目审查",
  ReviewMode: "code",
  CodeSource: "project_zip",
  ProjectZip: fakeZipHeader,
}
```

Expected assertions:

- file service receives `AssetTypeProjectZip`
- zip service is called
- AI prompt contains `前端项目 ZIP 源码摘要`
- AI request `NeedVision` is false

- [ ] **Step 3: Add validation test**

Assert `code` mode with no code and no project ZIP returns an error containing `code or project_zip`.

- [ ] **Step 4: Run test to verify failure**

Run: `go test ./internal/service/tools -run UIReview`

Expected: FAIL before implementation because `CodeSource`, `ProjectZip`, and zip service integration do not exist.

### Task 2: Backend implementation

**Files:**
- Modify: `backend/internal/service/tools/ui_review_service.go`
- Modify: `backend/internal/handler/tool_run_handler.go`
- Modify: `backend/cmd/server/main.go`
- Modify: `backend/internal/dto/ui_review_dto.go`
- Modify: `backend/internal/prompts/ui_review_prompt.go`

- [ ] **Step 1: Extend input structures**

Add `CodeSource string` and `ProjectZip *multipart.FileHeader` to UI Review form input. Add `CodeSource string` to `dto.UIReviewRequest`.

- [ ] **Step 2: Inject ZipService**

Add `zipService service.ZipService` to `UIReviewService` and constructor. Pass the existing `zipService` from `cmd/server/main.go`.

- [ ] **Step 3: Read multipart ZIP field**

In `RunUIReview`, read `code_source` and optional `project_zip`.

- [ ] **Step 4: Save and scan ZIP**

If `code_source=project_zip`, save upload using `service.AllowedArchiveTypes()`, resolve the ZIP path, call `ExtractAndAnalyze`, marshal and truncate summary.

- [ ] **Step 5: Update prompt**

Change `BuildUIReviewPrompt` to accept `codeSource` and `projectSummaryText`, and include a section for frontend project ZIP context.

- [ ] **Step 6: Update validation**

`code` mode requires pasted code or project ZIP. `screenshot_code` mode requires screenshot and pasted code or project ZIP.

- [ ] **Step 7: Run backend tests**

Run: `go test ./internal/service/tools -run UIReview`

Expected: PASS.

### Task 3: Frontend tests and implementation

**Files:**
- Modify: `frontend/src/pages/tools/UIReviewPage.test.ts`
- Modify: `frontend/src/pages/tools/UIReviewPage.vue`
- Modify: `frontend/src/types/tool.ts`

- [ ] **Step 1: Add failing frontend test**

Test that switching to `code` mode shows code source controls, selecting ZIP source shows ZIP upload, selecting a ZIP enables submit after title, and submit passes `code_source=project_zip` and `project_zip`.

- [ ] **Step 2: Implement state**

Add:

```ts
const codeSource = ref<'paste' | 'project_zip'>('paste')
const projectZipFile = ref<File | null>(null)
const projectZipName = ref('')
```

- [ ] **Step 3: Implement UI**

In code-related modes, render source segmented controls and either textarea or ZIP upload area.

- [ ] **Step 4: Update submit/canSubmit**

Append `code_source`, `project_zip` when selected, and update validation messages.

- [ ] **Step 5: Run frontend test**

Run: `npm.cmd exec vitest run src/pages/tools/UIReviewPage.test.ts -- --environment jsdom`

Expected: PASS.

### Task 4: Full verification

- [ ] Run `go test ./...`
- [ ] Run `npm.cmd run build`
- [ ] Review diff for scoped changes only.


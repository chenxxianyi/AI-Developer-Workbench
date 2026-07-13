<script setup lang="ts">
/**
 * UI Review Page
 * Upload screenshot or code for UI/UX analysis
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { runUIReview } from '@/api/tools'
import { downloadReport } from '@/api/reports'
import type { Report, UIReviewResult } from '@/types/report'
import { Eye, Image, Code, Download } from '@lucide/vue'
import ToolPageShell from '@/components/tool/ToolPageShell.vue'
import ToolFormSection from '@/components/tool/ToolFormSection.vue'
import FileUpload from '@/components/tool/FileUpload.vue'
import CodeInput from '@/components/tool/CodeInput.vue'
import ProjectPicker from '@/components/tool/ProjectPicker.vue'
import ScorePanel from '@/components/report/ScorePanel.vue'
import IssueList from '@/components/report/IssueList.vue'
import RecommendationList from '@/components/report/RecommendationList.vue'
import CodexPromptBox from '@/components/report/CodexPromptBox.vue'
import GeneratedFilesPanel from '@/components/report/GeneratedFilesPanel.vue'
import ActionItemsPanel from '@/components/report/ActionItemsPanel.vue'
import ScreenshotAnnotation from '@/components/report/ScreenshotAnnotation.vue'

const router = useRouter()
const route = useRoute()

// Lineage: parent_report_id carried from a "re-run" action.
const parentReportId = ref('')
const projectId = ref('')

// Form state
const title = ref('')
const reviewMode = ref<'screenshot' | 'code' | 'screenshot_code'>('screenshot')
const codeSource = ref<'paste' | 'project_zip'>('paste')
const pageType = ref('')
const targetStyle = ref('')
const description = ref('')
const codeInput = ref('')
const screenshotFile = ref<File | null>(null)
const mobileScreenshotFile = ref<File | null>(null)
const desktopViewport = ref('1440x900')
const mobileViewport = ref('390x844')
const projectZipFile = ref<File | null>(null)

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<UIReviewResult> | null>(null)
const exportError = ref<string | null>(null)

// Prefill safe (text-only) fields from the query when re-running. File uploads
// (screenshot / project zip) are intentionally NOT restored.
onMounted(() => {
  const q = route.query
  parentReportId.value = typeof q.parent_report_id === 'string' ? q.parent_report_id : ''
  projectId.value = typeof q.project_id === 'string' ? q.project_id : ''
  if (typeof q.title === 'string') title.value = q.title
  if (typeof q.review_mode === 'string') {
    const v = q.review_mode as typeof reviewMode.value
    if (v === 'screenshot' || v === 'code' || v === 'screenshot_code') reviewMode.value = v
  }
  if (typeof q.code_source === 'string') {
    const v = q.code_source as typeof codeSource.value
    if (v === 'paste' || v === 'project_zip') codeSource.value = v
  }
  if (typeof q.page_type === 'string') pageType.value = q.page_type
  if (typeof q.target_style === 'string') targetStyle.value = q.target_style
  if (typeof q.description === 'string') description.value = q.description
  if (typeof q.code === 'string') codeInput.value = q.code
})

function hasCodeSourceInput(): boolean {
  if (codeSource.value === 'project_zip') return !!projectZipFile.value
  return !!codeInput.value.trim()
}

// Submit
async function handleSubmit() {
  if (!title.value.trim()) {
    error.value = '请输入标题'
    return
  }

  if (reviewMode.value === 'screenshot' && !screenshotFile.value && !mobileScreenshotFile.value) {
    error.value = '请上传截图'
    return
  }

  if (reviewMode.value === 'code' && !hasCodeSourceInput()) {
    error.value = '请粘贴代码或上传前端项目 ZIP'
    return
  }

  if (reviewMode.value === 'screenshot_code' && ((!screenshotFile.value && !mobileScreenshotFile.value) || !hasCodeSourceInput())) {
    error.value = '请上传截图，并粘贴代码或上传前端项目 ZIP'
    return
  }

  loading.value = true
  error.value = null
  exportError.value = null
  result.value = null

  try {
    const formData = new FormData()
    formData.append('title', title.value)
    formData.append('review_mode', reviewMode.value)
    formData.append('code_source', codeSource.value)
    if (pageType.value) formData.append('page_type', pageType.value)
    if (targetStyle.value) formData.append('target_style', targetStyle.value)
    if (description.value) formData.append('description', description.value)
    if (codeSource.value === 'paste' && codeInput.value) formData.append('code', codeInput.value)
    if (codeSource.value === 'project_zip' && projectZipFile.value) formData.append('project_zip', projectZipFile.value)
    if (screenshotFile.value) {
      formData.append('screenshot', screenshotFile.value) // legacy clients/tests
      formData.append('desktop_screenshot', screenshotFile.value)
    }
    if (mobileScreenshotFile.value) formData.append('mobile_screenshot', mobileScreenshotFile.value)
    formData.append('desktop_viewport', desktopViewport.value)
    formData.append('mobile_viewport', mobileViewport.value)
    if (parentReportId.value) formData.append('parent_report_id', parentReportId.value)
    if (projectId.value) formData.append('project_id', projectId.value)

    result.value = await runUIReview(formData) as unknown as Report<UIReviewResult>
  } catch (err: any) {
    error.value = err.message || '分析失败'
  } finally {
    loading.value = false
  }
}

async function exportMarkdown() {
  if (!result.value) return
  exportError.value = null
  try {
    await downloadReport(result.value.id)
  } catch (err: any) {
    exportError.value = err.message || '导出失败'
  }
}

function resetForm() {
  title.value = ''
  reviewMode.value = 'screenshot'
  pageType.value = ''
  targetStyle.value = ''
  description.value = ''
  codeInput.value = ''
  screenshotFile.value = null
  mobileScreenshotFile.value = null
  codeSource.value = 'paste'
  projectZipFile.value = null
  projectId.value = ''
  result.value = null
  error.value = null
  exportError.value = null
}

const canSubmit = computed(() => {
  if (!title.value.trim()) return false
  if (reviewMode.value === 'screenshot') return !!screenshotFile.value || !!mobileScreenshotFile.value
  if (reviewMode.value === 'code') return hasCodeSourceInput()
  return (!!screenshotFile.value || !!mobileScreenshotFile.value) && hasCodeSourceInput()
})

const reportData = computed(() => result.value?.report_data ?? null)

const hasScores = computed(() =>
  reportData.value?.scores?.length
  && result.value?.total_score !== null
  && result.value?.grade !== null,
)

function goBack() {
  router.push('/dashboard')
}
</script>

<template>
  <ToolPageShell
    :icon="Eye"
    title="UI 质量审查"
    description="上传截图或前端代码，快速获得视觉层级、一致性、可访问性与改进建议。"
    step-text="填写输入 → 上传素材 → 开始分析 → 查看结果"
    accent="accent"
    :loading="loading"
    :error="error"
    :can-submit="!!canSubmit"
    submit-label="开始分析"
    submitting-label="分析中..."
    loading-hint="AI 正在分析..."
    :back-label="''"
    submit-testid="ui-review-submit"
    @submit="handleSubmit"
    @back="goBack"
  >
    <template #form>
      <ToolFormSection label="标题" required id-for="ui-review-title">
        <input
          id="ui-review-title"
          v-model="title"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="输入分析标题..."
        />
      </ToolFormSection>

      <ToolFormSection label="关联项目" optional id-for="ui-review-project" help-id="ui-review-project-help" help="可选。关联后会使用项目 UI 风格和规则，并把报告归入该项目。">
        <ProjectPicker v-model="projectId" input-id="ui-review-project" help-id="ui-review-project-help" />
      </ToolFormSection>

      <!-- Review Mode -->
      <ToolFormSection label="分析模式" required as-label>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
          <button
            data-testid="review-mode-screenshot"
            type="button"
            @click="reviewMode = 'screenshot'"
            :class="[
              'flex items-start gap-2 px-4 py-3 rounded-lg transition-smooth text-left cursor-pointer focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
              reviewMode === 'screenshot' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border',
            ]"
          >
            <Image :size="18" class="mt-0.5 shrink-0" />
            <span>
              <span class="block font-semibold">截图</span>
              <span class="block text-xs opacity-85 mt-1">仅上传截图，快速评估视觉质量</span>
            </span>
          </button>
          <button
            data-testid="review-mode-code"
            type="button"
            @click="reviewMode = 'code'"
            :class="[
              'flex items-start gap-2 px-4 py-3 rounded-lg transition-smooth text-left cursor-pointer focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
              reviewMode === 'code' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border',
            ]"
          >
            <Code :size="18" class="mt-0.5 shrink-0" />
            <span>
              <span class="block font-semibold">代码</span>
              <span class="block text-xs opacity-85 mt-1">仅粘贴前端代码，审查结构与样式</span>
            </span>
          </button>
          <button
            data-testid="review-mode-screenshot-code"
            type="button"
            @click="reviewMode = 'screenshot_code'"
            :class="[
              'flex items-start gap-2 px-4 py-3 rounded-lg transition-smooth text-left cursor-pointer focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
              reviewMode === 'screenshot_code' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border',
            ]"
          >
            <Eye :size="18" class="mt-0.5 shrink-0" />
            <span>
              <span class="block font-semibold">两者</span>
              <span class="block text-xs opacity-85 mt-1">截图 + 代码，获得更完整建议</span>
            </span>
          </button>
        </div>
      </ToolFormSection>

      <!-- Screenshot Upload -->
      <ToolFormSection
        v-if="reviewMode === 'screenshot' || reviewMode === 'screenshot_code'"
        label="上传截图"
        required
        id-for="ui-review-screenshot"
      >
        <div class="grid gap-4 lg:grid-cols-2">
          <div>
            <label class="mb-2 block text-sm font-medium text-text-primary">Desktop screenshot</label>
            <FileUpload v-model="screenshotFile" input-id="ui-review-screenshot" testid="screenshot" accept="image/png,image/jpeg,image/webp" accent="accent" pasteable preview empty-text="Click, drag, or paste a desktop screenshot" hint="PNG, JPG, WebP; recommended 1440x900" remove-label="Remove desktop screenshot" />
            <input v-model="desktopViewport" aria-label="Desktop viewport" class="mt-2 w-full rounded-lg border border-border bg-surface-muted px-3 py-2 text-sm text-text-primary" />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-text-primary">Mobile screenshot (optional)</label>
            <FileUpload v-model="mobileScreenshotFile" input-id="ui-review-mobile-screenshot" testid="mobile-screenshot" accept="image/png,image/jpeg,image/webp" accent="accent" pasteable preview empty-text="Upload a mobile screenshot" hint="PNG, JPG, WebP; recommended 390x844" remove-label="Remove mobile screenshot" />
            <input v-model="mobileViewport" aria-label="Mobile viewport" class="mt-2 w-full rounded-lg border border-border bg-surface-muted px-3 py-2 text-sm text-text-primary" />
          </div>
        </div>
      </ToolFormSection>

      <!-- Code Source -->
      <ToolFormSection
        v-if="reviewMode === 'code' || reviewMode === 'screenshot_code'"
        label="代码来源"
        required
        as-label
      >
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 mb-3">
          <button
            data-testid="code-source-paste"
            type="button"
            @click="codeSource = 'paste'"
            :class="[
              'rounded-lg px-4 py-3 text-left transition-smooth focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
              codeSource === 'paste' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border',
            ]"
          >
            <span class="block font-semibold">粘贴代码</span>
            <span class="block text-xs opacity-85 mt-1">适合单文件或代码片段</span>
          </button>
          <button
            data-testid="code-source-project-zip"
            type="button"
            @click="codeSource = 'project_zip'"
            :class="[
              'rounded-lg px-4 py-3 text-left transition-smooth focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
              codeSource === 'project_zip' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border',
            ]"
          >
            <span class="block font-semibold">上传前端 ZIP</span>
            <span class="block text-xs opacity-85 mt-1">适合完整前端项目</span>
          </button>
        </div>

        <ToolFormSection v-if="codeSource === 'paste'" label="前端代码" id-for="ui-review-code">
          <CodeInput
            id="ui-review-code"
            v-model="codeInput"
            :rows="8"
            placeholder="粘贴 Vue/React/HTML/CSS 代码..."
          />
        </ToolFormSection>

        <div v-else>
          <FileUpload
            v-model="projectZipFile"
            input-id="ui-review-project-zip"
            testid="project-zip"
            accept=".zip,application/zip,application/x-zip-compressed"
            accent="accent"
            empty-text="上传前端项目 ZIP"
            hint="系统只做静态读取，不执行代码"
            remove-label="移除前端项目 ZIP"
          />
        </div>
      </ToolFormSection>

      <ToolFormSection label="页面类型" optional id-for="ui-review-page-type">
        <input
          id="ui-review-page-type"
          v-model="pageType"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: 登录页、Dashboard、表单页..."
        />
      </ToolFormSection>

      <ToolFormSection label="目标风格" optional id-for="ui-review-target-style">
        <input
          id="ui-review-target-style"
          v-model="targetStyle"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: 简洁现代、Material Design..."
        />
      </ToolFormSection>

      <ToolFormSection label="补充说明" optional id-for="ui-review-description">
        <textarea
          id="ui-review-description"
          v-model="description"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          rows="3"
          placeholder="描述设计目标、用户场景等..."
        ></textarea>
      </ToolFormSection>
    </template>

    <template #actions>
      <button
        type="button"
        class="px-4 py-3 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth cursor-pointer focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="resetForm"
      >
        重置
      </button>
    </template>

    <template #empty>
      <div class="py-8">
        <div class="rounded-xl border border-accent/20 bg-accent-soft/60 p-5">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-10 h-10 rounded-lg bg-accent text-white flex items-center justify-center">
              <Eye :size="20" />
            </div>
            <div>
              <h3 class="font-semibold text-text-primary">将输出哪些内容</h3>
              <p class="text-sm text-text-secondary">提交后会在这里生成可执行的 UI 审查结论。</p>
            </div>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div class="rounded-lg bg-surface border border-border p-3">
              <p class="font-medium text-text-primary">评分维度</p>
              <p class="text-sm text-text-secondary mt-1">视觉层级、一致性、可访问性、对比度等。</p>
            </div>
            <div class="rounded-lg bg-surface border border-border p-3">
              <p class="font-medium text-text-primary">问题优先级</p>
              <p class="text-sm text-text-secondary mt-1">按 high / medium / low 聚焦关键问题。</p>
            </div>
            <div class="rounded-lg bg-surface border border-border p-3">
              <p class="font-medium text-text-primary">改进建议</p>
              <p class="text-sm text-text-secondary mt-1">给出可落地的设计和实现建议。</p>
            </div>
            <div class="rounded-lg bg-surface border border-border p-3">
              <p class="font-medium text-text-primary">报告文件</p>
              <p class="text-sm text-text-secondary mt-1">分析完成后可导出 Markdown 报告。</p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <template v-if="result" #result>
      <!-- Summary -->
      <div class="mb-6 p-4 bg-surface-muted rounded-lg">
        <div class="flex items-center justify-between mb-2">
          <span class="text-text-secondary">总体评分</span>
          <span v-if="result.total_score" class="font-bold text-xl text-text-primary">
            {{ result.total_score }}/100
            <span v-if="result.grade" class="ml-2">{{ result.grade }}</span>
          </span>
        </div>
        <p class="text-text-primary">{{ result.summary }}</p>
      </div>

      <!-- Scores -->
      <div v-if="hasScores" class="mb-6">
        <ScorePanel
          :scores="reportData!.scores!"
          :total-score="result.total_score!"
          :grade="result.grade!"
        />
      </div>

      <div v-if="reportData?.issues?.some(issue => issue.region)" class="mb-6 grid gap-4 lg:grid-cols-2">
        <ScreenshotAnnotation :file="screenshotFile" viewport="desktop" :viewport-label="desktopViewport" :issues="reportData.issues" />
        <ScreenshotAnnotation :file="mobileScreenshotFile" viewport="mobile" :viewport-label="mobileViewport" :issues="reportData.issues" />
      </div>

      <!-- Issues -->
      <div v-if="reportData?.issues?.length" class="mb-6">
        <IssueList :issues="reportData.issues" />
      </div>

      <!-- Action Items -->
      <ActionItemsPanel
        v-if="reportData?.action_items?.length || reportData?.recommendations?.length"
        :report-id="result.id"
        :report-title="result.title"
        :action-items="reportData?.action_items ?? []"
        :recommendations="reportData?.recommendations ?? []"
        class="mb-6"
      />

      <!-- Codex Prompt -->
      <CodexPromptBox
        v-if="reportData?.codex_prompt"
        :prompt="reportData.codex_prompt"
        class="mb-6"
      />

      <!-- Recommendations -->
      <RecommendationList
        v-if="reportData?.recommendations?.length"
        :recommendations="reportData.recommendations"
        class="mb-6"
      />

      <!-- Generated Files -->
      <GeneratedFilesPanel
        v-if="result.generated_files?.length"
        :files="result.generated_files"
        :report-id="result.id"
        class="mb-6"
      />

      <!-- Export -->
      <div class="flex items-center gap-3">
        <button
          type="button"
          class="inline-flex items-center gap-2 px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
          @click="exportMarkdown"
        >
          <Download :size="18" />
          <span>导出 Markdown</span>
        </button>
        <span v-if="exportError" class="text-sm text-danger">{{ exportError }}</span>
      </div>
    </template>
  </ToolPageShell>
</template>

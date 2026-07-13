<script setup lang="ts">
/**
 * Project Doctor Page
 * Upload project ZIP for comprehensive health analysis
 */
import { ref, computed, nextTick, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { runProjectDoctor } from '@/api/tools'
import { downloadReport } from '@/api/reports'
import type { Report, ProjectDoctorResult } from '@/types/report'
import { Stethoscope, Download } from '@lucide/vue'
import ToolPageShell from '@/components/tool/ToolPageShell.vue'
import ToolFormSection from '@/components/tool/ToolFormSection.vue'
import FileUpload from '@/components/tool/FileUpload.vue'
import ProjectPicker from '@/components/tool/ProjectPicker.vue'
import ScorePanel from '@/components/report/ScorePanel.vue'
import IssueList from '@/components/report/IssueList.vue'
import RecommendationList from '@/components/report/RecommendationList.vue'
import CodexPromptBox from '@/components/report/CodexPromptBox.vue'
import GeneratedFilesPanel from '@/components/report/GeneratedFilesPanel.vue'
import ActionItemsPanel from '@/components/report/ActionItemsPanel.vue'
import { focusFirstError } from '@/utils/focusFirstError'

const router = useRouter()
const route = useRoute()

// Lineage: parent_report_id carried from a "re-run" action.
const parentReportId = ref('')
const projectId = ref('')

// Form state
const title = ref('')
const projectName = ref('')
const techStack = ref('')
const projectDescription = ref('')
const analysisDepth = ref<'basic' | 'standard' | 'deep'>('standard')
const projectFile = ref<File | null>(null)

// Field-level validation errors
const errors = ref<{ title?: string; projectFile?: string }>({})

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<ProjectDoctorResult> | null>(null)
const exportError = ref<string | null>(null)

const analysisDepthOptions: Array<{
  value: 'basic' | 'standard' | 'deep'
  label: string
  description: string
}> = [
  { value: 'basic', label: '基础', description: '快速扫描关键风险' },
  { value: 'standard', label: '标准', description: '平衡速度与覆盖面' },
  { value: 'deep', label: '深度', description: '更全面的工程审查' },
]

// Prefill safe (text-only) fields from the query when re-running. File uploads
// are intentionally NOT restored (browser security forbids programmatic selection).
onMounted(() => {
  const q = route.query
  parentReportId.value = typeof q.parent_report_id === 'string' ? q.parent_report_id : ''
  projectId.value = typeof q.project_id === 'string' ? q.project_id : ''
  if (typeof q.title === 'string') title.value = q.title
  if (typeof q.project_name === 'string') projectName.value = q.project_name
  if (typeof q.tech_stack === 'string') techStack.value = q.tech_stack
  if (typeof q.project_description === 'string') projectDescription.value = q.project_description
  if (typeof q.analysis_depth === 'string') {
    const v = q.analysis_depth as typeof analysisDepth.value
    if (v === 'basic' || v === 'standard' || v === 'deep') analysisDepth.value = v
  }
})

// Submit
async function handleSubmit() {
  errors.value = {}
  if (!title.value.trim()) {
    errors.value.title = '请输入标题'
  }
  if (!projectFile.value) {
    errors.value.projectFile = '请上传项目 ZIP 文件'
  }
  if (errors.value.title || errors.value.projectFile) {
    await nextTick()
    focusFirstError()
    return
  }

  loading.value = true
  error.value = null
  exportError.value = null
  result.value = null

  try {
    const formData = new FormData()
    formData.append('title', title.value)
    formData.append('project_name', projectName.value || '未命名项目')
    formData.append('analysis_depth', analysisDepth.value)
    if (techStack.value) formData.append('tech_stack', techStack.value)
    if (projectDescription.value) formData.append('project_description', projectDescription.value)
    formData.append('project_zip', projectFile.value as Blob)
    if (parentReportId.value) formData.append('parent_report_id', parentReportId.value)
    if (projectId.value) formData.append('project_id', projectId.value)

    result.value = await runProjectDoctor(formData) as unknown as Report<ProjectDoctorResult>
  } catch (err: any) {
    error.value = err.message || '诊断失败'
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
  projectName.value = ''
  techStack.value = ''
  projectDescription.value = ''
  analysisDepth.value = 'standard'
  projectFile.value = null
  projectId.value = ''
  result.value = null
  error.value = null
  exportError.value = null
  errors.value = {}
}

const canSubmit = computed(() => !!(title.value.trim() && projectFile.value))

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
    :icon="Stethoscope"
    title="项目诊断"
    description="全面检查项目健康度和工程质量"
    step-text="上传项目 ZIP → 选择深度 → 开始诊断 → 查看结果"
    accent="success"
    :loading="loading"
    :error="error"
    :can-submit="!!canSubmit"
    submit-label="开始诊断"
    submitting-label="诊断中..."
    loading-hint="AI 正在分析项目..."
    :back-label="''"
    @submit="handleSubmit"
    @back="goBack"
  >
    <template #form>
      <ToolFormSection label="标题" required id-for="project-doctor-title" help-id="project-doctor-title-help" help="用于区分本次诊断报告。" :error="errors.title">
        <input
          id="project-doctor-title"
          v-model="title"
          type="text"
          required
          :aria-invalid="errors.title ? 'true' : undefined"
          :aria-describedby="errors.title ? 'project-doctor-title-error' : 'project-doctor-title-help'"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-success focus-visible:border-success focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="输入诊断标题..."
        />
      </ToolFormSection>

      <ToolFormSection label="关联项目" optional id-for="project-doctor-project" help-id="project-doctor-project-help" help="可选。关联后会沉淀诊断报告与质量趋势。">
        <ProjectPicker v-model="projectId" input-id="project-doctor-project" help-id="project-doctor-project-help" />
      </ToolFormSection>

      <ToolFormSection label="项目名称" optional id-for="project-doctor-name">
        <input
          id="project-doctor-name"
          v-model="projectName"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-success focus-visible:border-success focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: AI Workbench..."
        />
      </ToolFormSection>

      <!-- File Upload -->
      <ToolFormSection
        label="项目 ZIP 文件"
        required
        id-for="project-doctor-zip"
        help-id="project-doctor-zip-help"
        help="按 Enter 或空格选择文件，支持 .zip，最大 20MB"
        :error="errors.projectFile"
      >
        <FileUpload
          v-model="projectFile"
          input-id="project-doctor-zip"
          testid="project-zip"
          accept=".zip,application/zip,application/x-zip-compressed"
          accent="success"
          empty-text="点击、拖拽或按 Enter / 空格上传 ZIP"
          hint="支持 .zip，最大 20MB"
          help-id="project-doctor-zip-help"
          remove-label="移除已上传的项目 ZIP 文件"
        />
      </ToolFormSection>

      <ToolFormSection label="技术栈" optional id-for="project-doctor-tech-stack">
        <input
          id="project-doctor-tech-stack"
          v-model="techStack"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-success focus-visible:border-success focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: Vue 3 + Go + MySQL..."
        />
      </ToolFormSection>

      <ToolFormSection label="分析深度" required as-label>
        <div role="radiogroup" aria-label="分析深度" class="grid grid-cols-1 sm:grid-cols-3 gap-2">
          <button
            v-for="option in analysisDepthOptions"
            :key="option.value"
            type="button"
            role="radio"
            :aria-checked="analysisDepth === option.value"
            @click="analysisDepth = option.value"
            :class="[
              'flex flex-col items-start gap-1 px-4 py-3 rounded-lg transition-smooth text-left cursor-pointer focus-visible:ring-2 focus-visible:ring-success focus:outline-none',
              analysisDepth === option.value
                ? 'bg-success text-white'
                : 'bg-surface-muted text-text-secondary hover:bg-border',
            ]"
          >
            <span class="font-semibold">{{ option.label }}</span>
            <span class="text-xs opacity-85">{{ option.description }}</span>
          </button>
        </div>
      </ToolFormSection>

      <ToolFormSection label="项目描述" optional id-for="project-doctor-description">
        <textarea
          id="project-doctor-description"
          v-model="projectDescription"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-success focus-visible:border-success focus:outline-none text-text-primary placeholder:text-text-muted"
          rows="3"
          placeholder="描述项目功能和目标..."
        ></textarea>
      </ToolFormSection>
    </template>

    <template #actions>
      <button
        type="button"
        class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth focus-visible:ring-2 focus-visible:ring-success focus:outline-none cursor-pointer"
        @click="resetForm"
      >
        重置
      </button>
    </template>

    <template #empty>
      <div class="py-12 text-center">
        <Stethoscope :size="48" class="text-text-muted mx-auto mb-4" />
        <p class="text-text-secondary">上传项目 ZIP 开始诊断</p>
      </div>
    </template>

    <template v-if="result" #result>
      <!-- Summary -->
      <div class="mb-6 p-4 bg-surface-muted rounded-lg">
        <div class="flex items-center justify-between mb-2">
          <span class="text-text-secondary">健康评分</span>
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

      <!-- Recommendations -->
      <RecommendationList
        v-if="reportData?.recommendations?.length"
        :recommendations="reportData.recommendations"
        class="mb-6"
      />

      <!-- Codex Prompt -->
      <CodexPromptBox
        v-if="reportData?.codex_prompt"
        :prompt="reportData.codex_prompt"
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
          class="inline-flex items-center gap-2 px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth focus-visible:ring-2 focus-visible:ring-success focus:outline-none"
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

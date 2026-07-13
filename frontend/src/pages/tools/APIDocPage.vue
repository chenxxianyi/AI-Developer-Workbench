<script setup lang="ts">
/**
 * API Doc Builder Page
 * Generate API documentation from code or project ZIP
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { runAPIDoc } from '@/api/tools'
import { downloadReport } from '@/api/reports'
import type { Report, APIDocResult } from '@/types/report'
import { FileText, Download } from '@lucide/vue'
import ToolPageShell from '@/components/tool/ToolPageShell.vue'
import ToolFormSection from '@/components/tool/ToolFormSection.vue'
import FileUpload from '@/components/tool/FileUpload.vue'
import CodeInput from '@/components/tool/CodeInput.vue'
import ProjectPicker from '@/components/tool/ProjectPicker.vue'
import RecommendationList from '@/components/report/RecommendationList.vue'
import CodexPromptBox from '@/components/report/CodexPromptBox.vue'
import GeneratedFilesPanel from '@/components/report/GeneratedFilesPanel.vue'
import ActionItemsPanel from '@/components/report/ActionItemsPanel.vue'

const router = useRouter()
const route = useRoute()

// Lineage: parent_report_id carried from a "re-run" action.
const parentReportId = ref('')
const projectId = ref('')

// Form state
const title = ref('')
const sourceType = ref<'code' | 'project_zip' | 'manual'>('code')
const backendStack = ref('')
const codeInput = ref('')
const apiDescription = ref('')
const outputFormat = ref<'markdown' | 'openapi' | 'both'>('markdown')
const projectFile = ref<File | null>(null)

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<APIDocResult> | null>(null)
const exportError = ref<string | null>(null)

// Prefill safe (text-only) fields from the query when re-running. File uploads
// are intentionally NOT restored.
onMounted(() => {
  const q = route.query
  parentReportId.value = typeof q.parent_report_id === 'string' ? q.parent_report_id : ''
  projectId.value = typeof q.project_id === 'string' ? q.project_id : ''
  if (typeof q.title === 'string') title.value = q.title
  if (typeof q.source_type === 'string') {
    const v = q.source_type as typeof sourceType.value
    if (v === 'code' || v === 'project_zip' || v === 'manual') sourceType.value = v
  }
  if (typeof q.backend_stack === 'string') backendStack.value = q.backend_stack
  if (typeof q.code === 'string') codeInput.value = q.code
  if (typeof q.api_description === 'string') apiDescription.value = q.api_description
  if (typeof q.output_format === 'string') {
    const v = q.output_format as typeof outputFormat.value
    if (v === 'markdown' || v === 'openapi' || v === 'both') outputFormat.value = v
  }
})

// Submit
async function handleSubmit() {
  if (!title.value.trim()) {
    error.value = '请输入标题'
    return
  }

  if (sourceType.value === 'code' && !codeInput.value.trim()) {
    error.value = '请输入代码'
    return
  }

  if (sourceType.value === 'project_zip' && !projectFile.value) {
    error.value = '请上传项目 ZIP'
    return
  }

  if (sourceType.value === 'manual' && !apiDescription.value.trim()) {
    error.value = '请输入 API 描述'
    return
  }

  loading.value = true
  error.value = null
  exportError.value = null
  result.value = null

  try {
    const formData = new FormData()
    formData.append('title', title.value)
    formData.append('source_type', sourceType.value)
    formData.append('output_format', outputFormat.value)
    if (backendStack.value) formData.append('backend_stack', backendStack.value)
    if (codeInput.value) formData.append('code', codeInput.value)
    if (apiDescription.value) formData.append('api_description', apiDescription.value)
    if (projectFile.value) formData.append('project_zip', projectFile.value)
    if (parentReportId.value) formData.append('parent_report_id', parentReportId.value)
    if (projectId.value) formData.append('project_id', projectId.value)

    result.value = await runAPIDoc(formData) as unknown as Report<APIDocResult>
  } catch (err: any) {
    error.value = err.message || '生成失败'
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
  sourceType.value = 'code'
  backendStack.value = ''
  codeInput.value = ''
  apiDescription.value = ''
  outputFormat.value = 'markdown'
  projectFile.value = null
  projectId.value = ''
  result.value = null
  error.value = null
  exportError.value = null
}

const canSubmit = computed(() => {
  if (!title.value.trim()) return false
  if (sourceType.value === 'code') return !!codeInput.value.trim()
  if (sourceType.value === 'project_zip') return !!projectFile.value
  if (sourceType.value === 'manual') return !!apiDescription.value.trim()
  return false
})

const reportData = computed(() => result.value?.report_data ?? null)

function goBack() {
  router.push('/dashboard')
}
</script>

<template>
  <ToolPageShell
    :icon="FileText"
    title="API 文档生成"
    description="根据代码或接口描述自动生成 API 文档"
    step-text="选择来源 → 提供材料 → 生成文档 → 下载"
    accent="orange"
    :loading="loading"
    :error="error"
    :can-submit="!!canSubmit"
    submit-label="生成文档"
    submitting-label="生成中..."
    loading-hint="AI 正在生成文档..."
    @submit="handleSubmit"
    @back="goBack"
  >
    <template #form>
      <ToolFormSection label="标题" required id-for="api-doc-title">
        <input
          id="api-doc-title"
          v-model="title"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-orange-500 focus-visible:border-orange-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="输入文档标题..."
        />
      </ToolFormSection>

      <ToolFormSection label="关联项目" optional id-for="api-doc-project" help-id="api-doc-project-help" help="可选。关联后会使用项目后端技术栈和编码规则辅助生成。">
        <ProjectPicker v-model="projectId" input-id="api-doc-project" help-id="api-doc-project-help" />
      </ToolFormSection>

      <ToolFormSection label="来源类型" required as-label>
        <div class="flex gap-2" role="group" aria-label="来源类型">
          <button
            v-for="opt in [{v:'code',l:'代码'},{v:'project_zip',l:'项目 ZIP'},{v:'manual',l:'手动描述'}]"
            :key="opt.v"
            type="button"
            :class="[
              'flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth focus-visible:ring-2 focus-visible:ring-orange-500 focus:outline-none',
              sourceType === opt.v ? 'bg-orange-500 text-white' : 'bg-surface-muted text-text-secondary hover:bg-border',
            ]"
            @click="sourceType = opt.v as typeof sourceType"
          >
            <span>{{ opt.l }}</span>
          </button>
        </div>
      </ToolFormSection>

      <!-- Code Input -->
      <ToolFormSection
        v-if="sourceType === 'code'"
        label="API 代码"
        required
        id-for="api-doc-code"
        help-id="api-doc-code-help"
        help="粘贴处理器或控制器代码，支持 Go Gin、Express、Django 等。"
      >
        <CodeInput
          id="api-doc-code"
          v-model="codeInput"
          described-by="api-doc-code-help"
          :rows="10"
          placeholder="粘贴 API 处理器或控制器代码..."
        />
      </ToolFormSection>

      <!-- ZIP Upload -->
      <ToolFormSection
        v-if="sourceType === 'project_zip'"
        label="项目 ZIP"
        required
        id-for="api-doc-zip"
        help-id="api-doc-zip-help"
        help="按 Enter 或空格选择文件，支持 .zip，最大 20MB"
      >
        <FileUpload
          v-model="projectFile"
          input-id="api-doc-zip"
          testid="api-doc-zip"
          accept=".zip,application/zip,application/x-zip-compressed"
          accent="orange"
          empty-text="点击、拖拽或按 Enter / 空格上传 ZIP"
          hint="支持 .zip，最大 20MB"
          help-id="api-doc-zip-help"
        />
      </ToolFormSection>

      <!-- Manual Description -->
      <ToolFormSection v-if="sourceType === 'manual'" label="API 描述" required id-for="api-doc-desc">
        <textarea
          id="api-doc-desc"
          v-model="apiDescription"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-orange-500 focus-visible:border-orange-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          rows="6"
          placeholder="描述 API 功能、端点、参数等..."
        ></textarea>
      </ToolFormSection>

      <ToolFormSection label="后端技术栈" optional id-for="api-doc-stack">
        <input
          id="api-doc-stack"
          v-model="backendStack"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-orange-500 focus-visible:border-orange-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: Go Gin, Express, Django..."
        />
      </ToolFormSection>

      <ToolFormSection label="输出格式" required as-label>
        <div class="flex gap-2" role="group" aria-label="输出格式">
          <button
            v-for="opt in [{v:'markdown',l:'Markdown'},{v:'openapi',l:'OpenAPI'},{v:'both',l:'两者'}]"
            :key="opt.v"
            type="button"
            :class="[
              'px-4 py-2 rounded-lg transition-smooth focus-visible:ring-2 focus-visible:ring-orange-500 focus:outline-none',
              outputFormat === opt.v ? 'bg-orange-500 text-white' : 'bg-surface-muted text-text-secondary hover:bg-border',
            ]"
            @click="outputFormat = opt.v as typeof outputFormat"
          >
            {{ opt.l }}
          </button>
        </div>
      </ToolFormSection>
    </template>

    <template #actions>
      <button
        type="button"
        class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth focus-visible:ring-2 focus-visible:ring-orange-500 focus:outline-none cursor-pointer"
        @click="resetForm"
      >
        重置
      </button>
    </template>

    <template #empty>
      <div class="py-12 text-center">
        <FileText :size="48" class="text-text-muted mx-auto mb-4" />
        <p class="text-text-secondary">提交输入后生成文档</p>
      </div>
    </template>

    <template v-if="result" #result>
      <!-- Non-scoring notice -->
      <div class="mb-6 rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-2">非评分型分析</h3>
        <p class="text-sm text-text-muted">本工具不产生量化评分，请查看下方的模块、文档与生成文件。</p>
      </div>

      <!-- Summary -->
      <div class="mb-6 p-4 bg-surface-muted rounded-lg">
        <p class="text-text-primary">{{ result.summary }}</p>
      </div>

      <!-- Modules -->
      <div v-if="reportData?.modules?.length" class="mb-6 rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-3">API 模块</h3>
        <div class="space-y-2">
          <div v-for="module in reportData.modules" :key="module.name" class="p-3 bg-surface-muted rounded">
            <span class="font-medium text-text-primary">{{ module.name }}</span>
            <span class="text-text-muted text-sm ml-2">{{ module.endpoints?.length || 0 }} 个端点</span>
          </div>
        </div>
      </div>

      <!-- Markdown Content -->
      <div v-if="reportData?.markdown_content" class="mb-6 rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-3">Markdown 文档</h3>
        <pre class="p-3 bg-surface-muted rounded text-sm text-text-secondary overflow-x-auto whitespace-pre-wrap">{{ reportData.markdown_content.substring(0, 500) }}...</pre>
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
          class="inline-flex items-center gap-2 px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth focus-visible:ring-2 focus-visible:ring-orange-500 focus:outline-none"
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

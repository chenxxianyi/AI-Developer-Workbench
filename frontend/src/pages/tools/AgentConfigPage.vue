<script setup lang="ts">
/**
 * Agent Config Studio Page
 * Generate AI agent configuration files
 */
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { runAgentConfig } from '@/api/tools'
import { downloadReport } from '@/api/reports'
import type { AgentConfigInput } from '@/types/tool'
import type { Report, AgentConfigResult } from '@/types/report'
import { Bot, Download } from '@lucide/vue'
import ToolPageShell from '@/components/tool/ToolPageShell.vue'
import ToolFormSection from '@/components/tool/ToolFormSection.vue'
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
const projectName = ref('')
const projectType = ref('')
const frontendStack = ref('')
const backendStack = ref('')
const database = ref('')
const uiStyle = ref('')
const codingPreferences = ref('')
const strictRules = ref('')

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<AgentConfigResult> | null>(null)
const exportError = ref<string | null>(null)

// Prefill safe (text-only) fields from the query when re-running.
onMounted(() => {
  const q = route.query
  parentReportId.value = typeof q.parent_report_id === 'string' ? q.parent_report_id : ''
  projectId.value = typeof q.project_id === 'string' ? q.project_id : ''
  if (typeof q.title === 'string') title.value = q.title
  if (typeof q.project_name === 'string') projectName.value = q.project_name
  if (typeof q.project_type === 'string') projectType.value = q.project_type
  if (typeof q.frontend_stack === 'string') frontendStack.value = q.frontend_stack
  if (typeof q.backend_stack === 'string') backendStack.value = q.backend_stack
  if (typeof q.database === 'string') database.value = q.database
  if (typeof q.ui_style === 'string') uiStyle.value = q.ui_style
  if (typeof q.coding_preferences === 'string') codingPreferences.value = q.coding_preferences
  if (typeof q.strict_rules === 'string') strictRules.value = q.strict_rules
})

// Submit
async function handleSubmit() {
  if (!title.value.trim()) {
    error.value = '请输入标题'
    return
  }
  if (!projectName.value.trim()) {
    error.value = '请输入项目名称'
    return
  }

  loading.value = true
  error.value = null
  exportError.value = null
  result.value = null

  try {
    const payload: AgentConfigInput = {
      title: title.value,
      project_id: projectId.value || undefined,
      project_name: projectName.value,
      project_type: projectType.value || undefined,
      frontend_stack: frontendStack.value || undefined,
      backend_stack: backendStack.value || undefined,
      database: database.value || undefined,
      ui_style: uiStyle.value || undefined,
      coding_preferences: codingPreferences.value || undefined,
      strict_rules: strictRules.value || undefined,
      parent_report_id: parentReportId.value || undefined,
    }

    result.value = await runAgentConfig(payload) as unknown as Report<AgentConfigResult>
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
  projectName.value = ''
  projectType.value = ''
  frontendStack.value = ''
  backendStack.value = ''
  database.value = ''
  uiStyle.value = ''
  codingPreferences.value = ''
  strictRules.value = ''
  projectId.value = ''
  result.value = null
  error.value = null
  exportError.value = null
}

const canSubmit = computed(() => title.value.trim() && projectName.value.trim())

const reportData = computed(() => result.value?.report_data ?? null)

function goBack() {
  router.push('/dashboard')
}
</script>

<template>
  <ToolPageShell
    :icon="Bot"
    title="Agent 配置生成"
    description="为项目生成 AI Agent 配置文件"
    step-text="填写项目信息 → 生成配置 → 查看文件 → 下载"
    accent="purple"
    :loading="loading"
    :error="error"
    :can-submit="!!canSubmit"
    submit-label="生成配置"
    submitting-label="生成中..."
    loading-hint="AI 正在生成配置..."
    form-heading="项目配置"
    result-heading="生成结果"
    @submit="handleSubmit"
    @back="goBack"
  >
    <template #form>
      <ToolFormSection label="标题" required id-for="agent-config-title">
        <input
          id="agent-config-title"
          v-model="title"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="输入配置标题..."
        />
      </ToolFormSection>

      <ToolFormSection label="关联项目" optional id-for="agent-config-project" help-id="agent-config-project-help" help="选择后会保存报告归属，并把项目技术栈和规则作为上下文。">
        <ProjectPicker v-model="projectId" input-id="agent-config-project" help-id="agent-config-project-help" />
      </ToolFormSection>

      <ToolFormSection label="项目名称" required id-for="agent-config-name">
        <input
          id="agent-config-name"
          v-model="projectName"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: AI Workbench..."
        />
      </ToolFormSection>

      <ToolFormSection label="项目类型" optional id-for="agent-config-type">
        <input
          id="agent-config-type"
          v-model="projectType"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: fullstack, frontend, backend..."
        />
      </ToolFormSection>

      <ToolFormSection label="前端技术栈" optional id-for="agent-config-frontend">
        <input
          id="agent-config-frontend"
          v-model="frontendStack"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: Vue 3 + TypeScript + Vite..."
        />
      </ToolFormSection>

      <ToolFormSection label="后端技术栈" optional id-for="agent-config-backend">
        <input
          id="agent-config-backend"
          v-model="backendStack"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: Go + Gin + GORM..."
        />
      </ToolFormSection>

      <ToolFormSection label="数据库" optional id-for="agent-config-database">
        <input
          id="agent-config-database"
          v-model="database"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: MySQL, PostgreSQL..."
        />
      </ToolFormSection>

      <ToolFormSection label="UI 风格" optional id-for="agent-config-ui">
        <input
          id="agent-config-ui"
          v-model="uiStyle"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: Tailwind CSS, 简洁现代..."
        />
      </ToolFormSection>

      <ToolFormSection label="编码偏好" optional id-for="agent-config-coding">
        <textarea
          id="agent-config-coding"
          v-model="codingPreferences"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          rows="3"
          placeholder="描述编码风格和偏好..."
        ></textarea>
      </ToolFormSection>

      <ToolFormSection label="严格规则" optional id-for="agent-config-rules">
        <textarea
          id="agent-config-rules"
          v-model="strictRules"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:border-purple-500 focus:outline-none text-text-primary placeholder:text-text-muted"
          rows="3"
          placeholder="描述必须遵守的规则..."
        ></textarea>
      </ToolFormSection>
    </template>

    <template #actions>
      <button
        type="button"
        class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth focus-visible:ring-2 focus-visible:ring-purple-500 focus:outline-none cursor-pointer"
        @click="resetForm"
      >
        重置
      </button>
    </template>

    <template #empty>
      <div class="py-12 text-center">
        <Bot :size="48" class="text-text-muted mx-auto mb-4" />
        <p class="text-text-secondary">填写项目信息后生成配置文件</p>
      </div>
    </template>

    <template v-if="result" #result>
      <!-- Non-scoring notice -->
      <div class="mb-6 rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-2">非评分型分析</h3>
        <p class="text-sm text-text-muted">本工具不产生量化评分，请查看下方的生成文件、建议和行动项。</p>
      </div>

      <!-- Summary -->
      <div class="mb-6 p-4 bg-surface-muted rounded-lg">
        <p class="text-text-primary">{{ result.summary }}</p>
      </div>

      <!-- Generated Files Preview -->
      <div v-if="reportData?.generated_files_content" class="mb-6 rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-3">生成的配置文件预览</h3>
        <div class="space-y-3">
          <div v-for="(content, filename) in reportData.generated_files_content" :key="filename" class="p-3 bg-surface-muted rounded">
            <span class="font-medium text-accent">{{ filename }}</span>
            <pre class="text-text-secondary text-sm mt-2 overflow-x-auto whitespace-pre-wrap">{{ content.substring(0, 200) }}...</pre>
          </div>
        </div>
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

      <!-- Generated Files (download) -->
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
          class="inline-flex items-center gap-2 px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth focus-visible:ring-2 focus-visible:ring-purple-500 focus:outline-none"
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

<script setup lang="ts">
/**
 * DB Schema Review Page
 * Review and optimize database schema
 */
import { ref, computed, nextTick, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { runDBSchema } from '@/api/tools'
import { downloadReport } from '@/api/reports'
import type { DBSchemaInput } from '@/types/tool'
import type { Report, DBSchemaResult } from '@/types/report'
import { Database, Download, CheckCircle2 } from '@lucide/vue'
import ToolPageShell from '@/components/tool/ToolPageShell.vue'
import ToolFormSection from '@/components/tool/ToolFormSection.vue'
import ProjectPicker from '@/components/tool/ProjectPicker.vue'
import CodeInput from '@/components/tool/CodeInput.vue'
import ScorePanel from '@/components/report/ScorePanel.vue'
import IssueList from '@/components/report/IssueList.vue'
import RecommendationList from '@/components/report/RecommendationList.vue'
import CodexPromptBox from '@/components/report/CodexPromptBox.vue'
import GeneratedFilesPanel from '@/components/report/GeneratedFilesPanel.vue'
import { focusFirstError } from '@/utils/focusFirstError'

const router = useRouter()
const route = useRoute()

// Lineage: parent_report_id carried from a "re-run" action.
const parentReportId = ref('')
const projectId = ref('')

// Form state
const title = ref('')
const schemaType = ref('')
const databaseType = ref('')
const businessContext = ref('')
const schemaContent = ref('')
const targetGoal = ref('')

// Prefill safe (text-only) fields from the query when re-running.
onMounted(() => {
  const q = route.query
  parentReportId.value = typeof q.parent_report_id === 'string' ? q.parent_report_id : ''
  projectId.value = typeof q.project_id === 'string' ? q.project_id : ''
  if (typeof q.title === 'string') title.value = q.title
  if (typeof q.schema_type === 'string') schemaType.value = q.schema_type
  if (typeof q.database_type === 'string') databaseType.value = q.database_type
  if (typeof q.business_context === 'string') businessContext.value = q.business_context
  if (typeof q.schema_content === 'string') schemaContent.value = q.schema_content
  if (typeof q.target_goal === 'string') targetGoal.value = q.target_goal
})

// Field-level validation errors
const errors = ref<{ title?: string; schemaContent?: string }>({})

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<DBSchemaResult> | null>(null)
const exportError = ref<string | null>(null)

// Submit
async function handleSubmit() {
  errors.value = {}
  if (!title.value.trim()) {
    errors.value.title = '请输入标题'
  }
  if (!schemaContent.value.trim()) {
    errors.value.schemaContent = '请输入数据库 Schema'
  }
  if (errors.value.title || errors.value.schemaContent) {
    await nextTick()
    focusFirstError()
    return
  }

  loading.value = true
  error.value = null
  exportError.value = null
  result.value = null

  try {
    const payload: DBSchemaInput = {
      title: title.value,
      project_id: projectId.value || undefined,
      schema_type: schemaType.value || 'SQL',
      schema_content: schemaContent.value,
      database_type: databaseType.value || undefined,
      business_context: businessContext.value || undefined,
      target_goal: targetGoal.value || undefined,
      parent_report_id: parentReportId.value || undefined,
    }

    result.value = await runDBSchema(payload) as unknown as Report<DBSchemaResult>
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
  schemaType.value = ''
  databaseType.value = ''
  businessContext.value = ''
  schemaContent.value = ''
  targetGoal.value = ''
  projectId.value = ''
  result.value = null
  error.value = null
  exportError.value = null
  errors.value = {}
}

const canSubmit = computed(() => !!(title.value.trim() && schemaContent.value.trim()))

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
    :icon="Database"
    title="数据库结构审查"
    description="通过 AI 分析审查并优化数据库结构"
    step-text="填写参数 → 粘贴 Schema → 开始分析 → 查看结果"
    accent="accent"
    :loading="loading"
    :error="error"
    :can-submit="!!canSubmit"
    submit-label="开始分析"
    submitting-label="分析中..."
    loading-hint="AI 正在分析 Schema..."
    result-heading="分析结果"
    :back-label="''"
    @submit="handleSubmit"
    @back="goBack"
  >
    <template #form>
      <ToolFormSection label="标题" required id-for="db-schema-title" :error="errors.title">
        <input
          id="db-schema-title"
          v-model="title"
          type="text"
          :aria-invalid="errors.title ? 'true' : undefined"
          :aria-describedby="errors.title ? 'db-schema-title-error' : undefined"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="输入分析标题..."
        />
      </ToolFormSection>

      <ToolFormSection label="关联项目" optional id-for="db-schema-project" help-id="db-schema-project-help" help="可选。关联后可在项目页查看本次 Schema 审查。">
        <ProjectPicker v-model="projectId" input-id="db-schema-project" help-id="db-schema-project-help" />
      </ToolFormSection>

      <ToolFormSection label="Schema 类型" optional id-for="db-schema-type">
        <input
          id="db-schema-type"
          v-model="schemaType"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: SQL, GORM, Prisma..."
        />
      </ToolFormSection>

      <ToolFormSection label="数据库类型" optional id-for="db-schema-dbtype">
        <input
          id="db-schema-dbtype"
          v-model="databaseType"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: MySQL, PostgreSQL..."
        />
      </ToolFormSection>

      <ToolFormSection label="数据库 Schema" required id-for="db-schema-content" help-id="db-schema-content-help" help="粘贴 CREATE TABLE、GORM model、Prisma schema 等。" :error="errors.schemaContent">
        <CodeInput
          id="db-schema-content"
          v-model="schemaContent"
          described-by="db-schema-content-help"
          :rows="12"
          placeholder="粘贴 CREATE TABLE、GORM model、Prisma schema 等..."
        />
      </ToolFormSection>

      <ToolFormSection label="业务背景" optional id-for="db-schema-context">
        <textarea
          id="db-schema-context"
          v-model="businessContext"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          rows="3"
          placeholder="描述业务场景和数据使用方式..."
        ></textarea>
      </ToolFormSection>

      <ToolFormSection label="优化目标" optional id-for="db-schema-goal">
        <input
          id="db-schema-goal"
          v-model="targetGoal"
          type="text"
          class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
          placeholder="如: 性能优化、扩展性、数据完整性..."
        />
      </ToolFormSection>
    </template>

    <template #actions>
      <button
        type="button"
        class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth focus-visible:ring-2 focus-visible:ring-accent focus:outline-none cursor-pointer"
        @click="resetForm"
      >
        重置
      </button>
    </template>

    <template #empty>
      <div class="py-12 text-center">
        <Database :size="48" class="text-text-muted mx-auto mb-4" />
        <p class="text-text-secondary">输入 Schema 后开始分析</p>
      </div>
    </template>

    <template v-if="result" #result>
      <!-- Summary -->
      <div class="mb-6 p-4 bg-surface-muted rounded-lg">
        <div class="flex items-center justify-between mb-2">
          <span class="text-text-secondary">评分</span>
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

      <!-- Optimized Schema -->
      <div v-if="reportData?.optimized_schema" class="mb-6 rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-3">优化后的 Schema</h3>
        <pre class="p-3 bg-surface-muted rounded text-sm text-text-secondary overflow-x-auto whitespace-pre-wrap">{{ reportData.optimized_schema }}</pre>
      </div>

      <!-- Migration Suggestions -->
      <div v-if="reportData?.migration_suggestions?.length" class="mb-6 rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-3">迁移建议</h3>
        <ul class="space-y-1">
          <li v-for="(sug, idx) in reportData.migration_suggestions" :key="idx" class="flex items-start gap-2 text-text-secondary">
            <CheckCircle2 :size="14" class="text-accent mt-1" />
            <span class="text-sm">{{ sug }}</span>
          </li>
        </ul>
      </div>

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

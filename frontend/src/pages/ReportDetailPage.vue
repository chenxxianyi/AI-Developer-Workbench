<script setup lang="ts">
/**
 * Report Detail Page
 * Displays report with unified information hierarchy:
 *   Summary → Score/Grade → Top Issues → Recommendations → Prompt / Downloads
 */
import { onMounted, computed, watch, ref } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { useReportStore } from '@/stores/reportStore'
import { getToolDisplayMeta } from '@/utils/toolDisplay'
import type { Report, Issue } from '@/types/report'
import { ArrowLeft, Download, FileText, Trash2 } from '@lucide/vue'
import ReportStatusBadge from '@/components/report/ReportStatusBadge.vue'
import ScorePanel from '@/components/report/ScorePanel.vue'
import IssueList from '@/components/report/IssueList.vue'
import RecommendationList from '@/components/report/RecommendationList.vue'
import CodexPromptBox from '@/components/report/CodexPromptBox.vue'
import GeneratedFilesPanel from '@/components/report/GeneratedFilesPanel.vue'
import ReportErrorState from '@/components/report/ReportErrorState.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'

const route = useRoute()
const router = useRouter()
const store = useReportStore()

const showDeleteConfirm = ref(false)

const reportId = computed(() => route.params.id as string)

onMounted(() => {
  store.fetchReport(reportId.value)
})

// Reload when route changes.
watch(reportId, (id) => {
  store.fetchReport(id)
})

const report = computed(() => store.currentReport)
const toolMeta = computed(() =>
  report.value ? getToolDisplayMeta(report.value.tool_type) : null,
)

// Parse report_data based on tool_type for a unified view.
interface ParsedReport {
  scores: Array<{ name: string; score: number; max_score: number; comment: string }>
  issues: Issue[]
  recommendations: string[]
  codexPrompt: string
}

const parsedData = computed<ParsedReport>(() => {
  if (!report.value) return { scores: [], issues: [], recommendations: [], codexPrompt: '' }
  const data = report.value.report_data as Record<string, unknown>

  return {
    scores: Array.isArray(data.scores) ? data.scores as ParsedReport['scores'] : [],
    issues: Array.isArray(data.issues) ? data.issues as Issue[] : [],
    recommendations: Array.isArray(data.recommendations) ? data.recommendations as string[] : [],
    codexPrompt: typeof data.codex_prompt === 'string' ? data.codex_prompt : '',
  }
})

// Top 3 issues by severity.
const topIssues = computed(() => {
  const issues = parsedData.value.issues
  if (!issues.length) return []
  const sorted = [...issues].sort((a, b) => {
    const order = { high: 0, medium: 1, low: 2 }
    return (order[a.severity] ?? 1) - (order[b.severity] ?? 1)
  })
  return sorted.slice(0, 3)
})

function hasScores(report: Report): boolean {
  return report.total_score !== null && report.grade !== null
}

function formatDate(dateString: string): string {
  const d = new Date(dateString)
  if (isNaN(d.getTime())) return dateString
  return d.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

async function handleDelete() {
  try {
    await store.deleteReportById(reportId.value)
    showDeleteConfirm.value = false
    router.push('/reports')
  } catch {
    showDeleteConfirm.value = false
  }
}

async function handleDownloadMarkdown() {
  const { downloadReport } = await import('@/api/reports')
  try {
    await downloadReport(reportId.value)
  } catch {
    // Error handled silently.
  }
}
</script>

<template>
  <div>
    <!-- Loading -->
    <ReportErrorState v-if="store.detailLoading" type="loading" message="正在加载报告..." />

    <!-- Error -->
    <ReportErrorState
      v-else-if="store.detailError"
      type="error"
      :message="store.detailError"
    />

    <!-- Not Found -->
    <ReportErrorState
      v-else-if="!report"
      type="not-found"
    />

    <!-- Report Content -->
    <div v-else>
      <!-- Header -->
      <div class="mb-6">
        <RouterLink
          to="/reports"
          class="inline-flex items-center gap-1.5 text-sm text-text-muted hover:text-accent transition-smooth mb-4"
        >
          <ArrowLeft :size="16" />
          返回列表
        </RouterLink>

        <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div class="min-w-0">
            <div class="flex items-center gap-2 mb-1.5">
              <ReportStatusBadge :status="report.status" />
              <span class="text-sm font-medium text-text-muted">
                {{ toolMeta?.name }} · {{ formatDate(report.created_at) }}
              </span>
            </div>
            <h1 class="text-2xl font-bold text-text-primary truncate">{{ report.title }}</h1>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-2 flex-shrink-0">
            <button
              v-if="report.status === 'succeeded' || report.status === 'fallback'"
              class="inline-flex items-center gap-1.5 rounded-md border border-border bg-surface px-4 py-2 text-sm font-medium text-text-primary hover:bg-surface-hover transition-smooth"
              @click="handleDownloadMarkdown"
              aria-label="下载 Markdown 报告"
            >
              <Download :size="16" />
              导出
            </button>
            <button
              class="inline-flex items-center gap-1.5 rounded-md border border-danger/30 bg-surface px-4 py-2 text-sm font-medium text-danger hover:bg-danger/5 transition-smooth"
              @click="showDeleteConfirm = true"
              aria-label="删除报告"
            >
              <Trash2 :size="16" />
              删除
            </button>
          </div>
        </div>
      </div>

      <!-- Failed state -->
      <div v-if="report.status === 'failed'" class="rounded-lg border border-danger/30 bg-danger/5 p-5 mb-6">
        <p class="text-sm text-danger font-medium">此报告生成失败。</p>
        <p class="text-sm text-text-muted mt-1">{{ report.summary }}</p>
      </div>

      <!-- Processing state -->
      <div v-else-if="report.status === 'processing'" class="rounded-lg border border-dashed border-border p-8 text-center mb-6">
        <FileText :size="32" class="mx-auto mb-3 text-text-muted" />
        <p class="text-text-muted">报告正在生成中，请稍后刷新页面。</p>
      </div>

      <!-- Success / Fallback content -->
      <div v-else class="space-y-6">
        <!-- Summary -->
        <div class="rounded-lg border border-border bg-surface p-5">
          <h3 class="text-lg font-semibold text-text-primary mb-3">结论摘要</h3>
          <p class="text-text-secondary leading-relaxed">{{ report.summary }}</p>
          <div
            v-if="report.status === 'fallback'"
            class="mt-3 flex items-center gap-2 text-sm text-amber-600 dark:text-amber-400"
          >
            <span class="inline-block h-2 w-2 rounded-full bg-amber-500" />
            此报告使用了降级数据，部分内容可能不完整。
          </div>
        </div>

        <!-- Score Panel (only for scoring tools) -->
        <ScorePanel
          v-if="hasScores(report) && parsedData.scores.length"
          :scores="parsedData.scores"
          :total-score="report.total_score!"
          :grade="report.grade!"
        />

        <!-- Non-scoring tools notice -->
        <div
          v-else-if="!hasScores(report)"
          class="rounded-lg border border-border bg-surface p-5"
        >
          <h3 class="text-lg font-semibold text-text-primary mb-2">非评分型分析</h3>
          <p class="text-sm text-text-muted">
            {{ toolMeta?.name }} 不产生量化评分，请查看下方的生成文件和建议。
          </p>
        </div>

        <!-- Top 3 Issues -->
        <div v-if="topIssues.length" class="rounded-lg border border-border bg-surface p-5">
          <h3 class="text-lg font-semibold text-text-primary mb-4">最严重问题</h3>
          <IssueList :issues="topIssues" />
        </div>

        <!-- Codex Prompt -->
        <CodexPromptBox v-if="parsedData.codexPrompt" :prompt="parsedData.codexPrompt" />

        <!-- Recommendations -->
        <RecommendationList :recommendations="parsedData.recommendations" />

        <!-- Generated Files -->
        <GeneratedFilesPanel
          :files="report.generated_files"
          :report-id="report.id"
        />

        <!-- Full Issues (if more than 3) -->
        <IssueList
          v-if="parsedData.issues.length > 3"
          :issues="parsedData.issues"
        />
      </div>
    </div>
  </div>

  <!-- Delete Confirmation -->
  <ConfirmDialog
    :open="showDeleteConfirm"
    title="删除报告"
    message="删除后将同时清除关联的上传文件和生成文件，此操作不可撤销。确定要继续吗？"
    confirm-label="删除"
    danger
    @confirm="handleDelete"
    @cancel="showDeleteConfirm = false"
  />
</template>

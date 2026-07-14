<script setup lang="ts">
/**
 * Project Detail Page
 * Overview, stats, recent reports, quick-run entry
 */
import { onMounted, computed, ref } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { useProjectStore } from '@/stores/projectStore'
import { getToolDisplayMeta } from '@/utils/toolDisplay'
import type { Report } from '@/types/report'
import type { ToolType } from '@/types/tool'
import ReportStatusBadge from '@/components/report/ReportStatusBadge.vue'
import QualityTrend from '@/components/dashboard/QualityTrend.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'
import { ArrowLeft, Pencil, Trash2, Folder, Wrench, Download, RotateCcw, GitCompareArrows } from '@lucide/vue'

const route = useRoute()
const router = useRouter()
const store = useProjectStore()

const projectId = computed(() => route.params.projectId as string)

onMounted(async () => {
  await Promise.all([
    store.fetchProject(projectId.value),
    store.fetchStats(projectId.value),
    store.fetchProjectReports(projectId.value),
  ])
})

const project = computed(() => store.currentProject)
const stats = computed(() => store.currentStats)
const historyTotalPages = computed(() => Math.ceil(store.projectReportsTotal / store.pageSize) || 1)
const deleteError = ref<string | null>(null)
const artifactToolTypes: ToolType[] = ['agent_config', 'api_doc', 'db_schema']
const missingArtifactToolTypes = computed(() => {
  const availableTools = new Set(stats.value?.latest_artifacts.map((artifact) => artifact.tool_type) ?? [])
  return artifactToolTypes.filter((toolType) => !availableTools.has(toolType))
})

function goEdit() {
  router.push(`/projects/${projectId.value}/edit`)
}

async function handleDelete() {
  if (!confirm('删除项目后，关联报告不会被删除，只会解除归属。确定继续吗？')) return
  deleteError.value = null
  try {
    await store.remove(projectId.value)
    router.push('/projects')
  } catch (error: any) {
    deleteError.value = error.message || '删除项目失败'
  }
}

function formatDate(s: string): string {
  return new Date(s).toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function toolPath(toolType: string): string {
  const paths: Record<string, string> = {
    ui_review: '/tools/ui-review',
    project_doctor: '/tools/project-doctor',
    agent_config: '/tools/agent-config',
    api_doc: '/tools/api-doc',
    db_schema: '/tools/db-schema',
  }
  return paths[toolType] || '/dashboard'
}

async function changeHistoryPage(page: number) {
  await store.fetchProjectReports(projectId.value, page)
}

function rerunQuery(report: Report): Record<string, string> {
  const query: Record<string, string> = {
    parent_report_id: report.id,
    project_id: project.value?.id || report.project_id || projectId.value,
  }

  for (const [key, value] of Object.entries(report.input_data ?? {})) {
    if (typeof value === 'string' && value && key !== 'parent_report_id' && key !== 'project_id') {
      query[key] = value
    }
  }

  return query
}
</script>

<template>
  <div class="max-w-5xl mx-auto">
    <RouterLink to="/projects" class="inline-flex items-center gap-1.5 text-sm text-text-muted hover:text-accent transition-smooth mb-4">
      <ArrowLeft :size="16" />
      返回项目列表
    </RouterLink>

    <!-- Loading -->
    <div v-if="store.detailLoading" class="py-16 text-center text-text-muted">加载中…</div>

    <!-- Error / Not Found -->
    <div v-else-if="store.detailError || !project" class="py-16 text-center">
      <Folder :size="48" class="mx-auto mb-3 text-text-muted" />
      <p class="text-text-muted">{{ store.detailError || '项目未找到' }}</p>
    </div>

    <!-- Content -->
    <div v-else>
      <div class="flex items-start justify-between gap-4 mb-6">
        <div>
          <h1 class="text-2xl font-bold text-text-primary">{{ project.name }}</h1>
          <p v-if="project.description" class="text-text-secondary mt-1">{{ project.description }}</p>
          <p v-if="project.repo_url" class="text-sm text-accent mt-1">{{ project.repo_url }}</p>
          <div class="flex gap-4 mt-2 text-xs text-text-muted">
            <span v-if="project.frontend_stack">前端: {{ project.frontend_stack }}</span>
            <span v-if="project.backend_stack">后端: {{ project.backend_stack }}</span>
            <span v-if="project.database">DB: {{ project.database }}</span>
          </div>
        </div>
        <div class="flex gap-2">
          <button class="inline-flex items-center gap-1.5 px-3 py-2 border border-border bg-surface rounded-md text-sm text-text-primary hover:bg-surface-hover transition-smooth" @click="goEdit">
            <Pencil :size="14" /> 编辑
          </button>
          <button class="inline-flex items-center gap-1.5 px-3 py-2 border border-danger/30 bg-surface rounded-md text-sm text-danger hover:bg-danger/5 transition-smooth" @click="handleDelete">
            <Trash2 :size="14" /> 删除
          </button>
        </div>
      </div>
      <p v-if="deleteError" role="alert" class="mb-4 text-sm text-danger">{{ deleteError }}</p>

      <!-- Stats -->
      <div v-if="stats" class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
        <div class="rounded-lg border border-border bg-surface p-4">
          <p class="text-xs text-text-muted">报告总数</p>
          <p class="text-2xl font-bold text-text-primary mt-1">{{ stats.total_reports }}</p>
        </div>
        <div class="rounded-lg border border-border bg-surface p-4">
          <p class="text-xs text-text-muted">平均分</p>
          <p class="text-2xl font-bold text-text-primary mt-1">{{ stats.average_score !== null ? stats.average_score.toFixed(0) : '—' }}</p>
        </div>
        <div class="rounded-lg border border-border bg-surface p-4">
          <p class="text-xs text-text-muted">高危问题</p>
          <p class="text-2xl font-bold mt-1" :class="stats.high_severity_count > 0 ? 'text-danger' : 'text-success'">{{ stats.high_severity_count }}</p>
        </div>
        <div class="rounded-lg border border-border bg-surface p-4">
          <p class="text-xs text-text-muted">工具使用</p>
          <p class="text-2xl font-bold text-text-primary mt-1">{{ Object.keys(stats?.tool_usage ?? {}).filter(k => ((stats?.tool_usage as any) ?? {})[k] > 0).length }}</p>
        </div>
      </div>

      <QualityTrend
        v-if="stats"
        :points="stats.quality_trend"
        class="mb-6"
      />

      <!-- Latest artifacts -->
      <section class="mb-6">
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-lg font-semibold text-text-primary">最新产物</h2>
          <span class="text-sm text-text-muted">Agent 配置、API 文档、Schema 建议</span>
        </div>
        <div v-if="!stats?.latest_artifacts?.length" class="border border-dashed border-border rounded-lg py-8 text-center text-sm text-text-muted">
          还没有可下载的项目产物。
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div
            v-for="artifact in stats.latest_artifacts"
            :key="`${artifact.report_id}-${artifact.filename}`"
            class="flex items-center justify-between gap-3 border border-border bg-surface px-4 py-3 rounded-lg"
          >
            <div class="min-w-0">
              <p class="text-sm font-medium text-text-primary truncate">{{ artifact.filename }}</p>
              <p class="text-xs text-text-muted mt-1">{{ getToolDisplayMeta(artifact.tool_type)?.name }} · {{ formatDate(artifact.created_at) }}</p>
            </div>
            <a
              :href="`/api/reports/${artifact.report_id}/files/${encodeURIComponent(artifact.filename)}`"
              :download="artifact.filename"
              class="inline-flex h-9 w-9 items-center justify-center rounded-md border border-border text-accent hover:bg-accent-soft focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
              :aria-label="`下载 ${artifact.filename}`"
              title="下载文件"
            >
              <Download :size="16" />
            </a>
          </div>
        </div>
        <div v-if="stats && missingArtifactToolTypes.length" class="mt-4 flex flex-wrap items-center gap-2">
          <span class="text-sm text-text-muted">缺少的产物：</span>
          <RouterLink
            v-for="toolType in missingArtifactToolTypes"
            :key="toolType"
            :to="{ path: toolPath(toolType), query: { project_id: project.id } }"
            class="inline-flex items-center rounded-md border border-border bg-surface-muted px-3 py-1.5 text-sm text-text-secondary hover:border-accent/35 hover:text-accent transition-smooth"
          >
            {{ getToolDisplayMeta(toolType)?.name }}
          </RouterLink>
        </div>
      </section>

      <!-- Project history -->
      <div class="rounded-lg border border-border bg-surface p-5 mb-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-text-primary">项目报告</h2>
          <span class="text-sm text-text-muted">{{ store.projectReportsTotal }} 份</span>
        </div>
        <div v-if="store.projectReportsLoading" class="py-6 text-center text-text-muted">正在加载报告…</div>
        <div v-else-if="store.projectReportsError" class="py-6 text-center text-danger">{{ store.projectReportsError }}</div>
        <div v-else-if="!store.projectReports.length" class="py-6 text-center text-text-muted">暂无报告</div>
        <div v-else class="space-y-2">
          <div
            v-for="r in store.projectReports"
            :key="r.id"
            class="flex items-center justify-between gap-3 rounded-md border border-border bg-background/50 px-4 py-3 hover:border-accent/35 transition-smooth"
          >
            <RouterLink :to="`/reports/${r.id}`" class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <ReportStatusBadge :status="r.status" />
                <span class="text-sm font-medium text-text-primary truncate">{{ r.title }}</span>
              </div>
              <p class="text-xs text-text-muted mt-1">{{ getToolDisplayMeta(r.tool_type)?.name }} · {{ formatDate(r.created_at) }}</p>
            </RouterLink>
            <div class="flex items-center gap-1.5 shrink-0">
            <span v-if="r.total_score !== null" class="text-sm font-semibold text-text-primary">{{ r.total_score }}</span>
              <RouterLink
                :to="{ path: toolPath(r.tool_type), query: rerunQuery(r) }"
                class="inline-flex h-8 w-8 items-center justify-center rounded-md border border-border text-text-secondary hover:border-accent/35 hover:text-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
                :aria-label="`基于 ${r.title} 复查`"
                :title="`基于 ${r.title} 复查`"
              >
                <RotateCcw :size="15" />
              </RouterLink>
              <RouterLink
                v-if="r.parent_report_id"
                :to="`/reports/${r.parent_report_id}/compare/${r.id}`"
                class="inline-flex h-8 w-8 items-center justify-center rounded-md border border-border text-text-secondary hover:border-accent/35 hover:text-accent focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
                :aria-label="`对比 ${r.title}`"
                :title="`对比 ${r.title}`"
              >
                <GitCompareArrows :size="15" />
              </RouterLink>
            </div>
          </div>
        </div>
        <PaginationBar
          v-if="historyTotalPages > 1"
          :current-page="store.projectReportPage"
          :total-pages="historyTotalPages"
          :has-prev="store.projectReportPage > 1"
          :has-next="store.projectReportPage < historyTotalPages"
          @prev="changeHistoryPage(store.projectReportPage - 1)"
          @next="changeHistoryPage(store.projectReportPage + 1)"
        />
      </div>

      <!-- Quick run -->
      <div class="rounded-lg border border-border bg-surface p-5">
        <h2 class="text-lg font-semibold text-text-primary mb-3">
          <Wrench :size="18" class="inline-block mr-1.5 text-accent" />
          快捷运行
        </h2>
        <p class="text-sm text-text-muted mb-3">选择工具后可关联本项目运行分析。</p>
        <div class="flex flex-wrap gap-2">
          <RouterLink
            v-for="t in ['ui_review','project_doctor','agent_config','api_doc','db_schema']"
            :key="t"
            :to="{ path: toolPath(t), query: { project_id: project.id } }"
            class="px-4 py-2 rounded-lg border border-border bg-surface-muted text-sm text-text-secondary hover:border-accent/35 hover:text-accent transition-smooth"
          >
            {{ getToolDisplayMeta(t as any)?.name }}
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

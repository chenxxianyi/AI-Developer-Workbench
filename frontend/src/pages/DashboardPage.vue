<script setup lang="ts">
/**
 * Dashboard Page
 * Stats, tool grid, and recent reports
 */

import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useSystemStore } from '@/stores/systemStore'
import { useToolStore } from '@/stores/toolStore'
import { getToolDisplayMeta } from '@/utils/toolDisplay'
import type { ToolType } from '@/types/tool'
import {
  Wrench,
  Eye,
  Stethoscope,
  Bot,
  FileText,
  Database,
  ArrowRight,
  ChevronRight,
  FileStack,
  BarChart3,
  Activity,
  CheckCircle2,
  Sparkles,
  FlaskConical,
} from '@lucide/vue'

const systemStore = useSystemStore()
const toolStore = useToolStore()

onMounted(async () => {
  await systemStore.fetchDashboardStats()
  await systemStore.fetchStatus()
  await toolStore.fetchTools()
})

const toolToneMap: Record<string, { icon: string; card: string; text: string; soft: string; border: string }> = {
  accent: {
    icon: 'bg-accent-soft text-accent border-accent/10',
    card: 'hover:border-accent/35 hover:shadow-[0_14px_34px_rgba(37,99,235,0.12)]',
    text: 'text-accent',
    soft: 'bg-accent-soft',
    border: 'border-accent/25',
  },
  success: {
    icon: 'bg-success/10 text-success border-success/15',
    card: 'hover:border-success/30 hover:shadow-[0_14px_34px_rgba(22,163,74,0.10)]',
    text: 'text-success',
    soft: 'bg-success/10',
    border: 'border-success/25',
  },
  warning: {
    icon: 'bg-warning/10 text-warning border-warning/15',
    card: 'hover:border-warning/30 hover:shadow-[0_14px_34px_rgba(217,119,6,0.10)]',
    text: 'text-warning',
    soft: 'bg-warning/10',
    border: 'border-warning/25',
  },
  danger: {
    icon: 'bg-danger/10 text-danger border-danger/15',
    card: 'hover:border-danger/30 hover:shadow-[0_14px_34px_rgba(220,38,38,0.10)]',
    text: 'text-danger',
    soft: 'bg-danger/10',
    border: 'border-danger/25',
  },
}

const stats = computed(() => [
  {
    label: '总报告数',
    value: systemStore.dashboardStats?.total_reports ?? 0,
    icon: FileStack,
    color: 'accent',
  },
  {
    label: '平均评分',
    value: systemStore.dashboardStats?.average_score ?? '暂无',
    suffix: '/100',
    icon: BarChart3,
    color: 'success',
  },
  {
    label: '工具使用',
    value: Object.keys(systemStore.dashboardStats?.tool_usage ?? {}).length,
    suffix: '/5',
    icon: Wrench,
    color: 'warning',
  },
  {
    label: '最近活动',
    value: latestActivityLabel.value,
    suffix: '',
    icon: Activity,
    color: 'danger',
  },
])

const recentReports = computed(() => systemStore.dashboardStats?.recent_reports?.slice(0, 4) ?? [])

const latestActivityLabel = computed(() => {
  const latestReport = systemStore.dashboardStats?.recent_reports?.[0]
  return latestReport ? formatRelativeTime(latestReport.created_at) : '暂无'
})

const totalToolUsage = computed(() => {
  const usage = systemStore.dashboardStats?.tool_usage ?? {}
  return (Object.values(usage) as number[]).reduce((sum, count) => sum + count, 0)
})

const mostUsedTool = computed(() => {
  const tools = [...toolStore.tools]
  if (!tools.length) return null

  return tools.sort((a, b) => b.usage_count - a.usage_count)[0]
})

const fallbackTools = computed(() => {
  const usage: Partial<Record<ToolType, number>> = systemStore.dashboardStats?.tool_usage ?? {}
  const toolTypes: ToolType[] = ['ui_review', 'project_doctor', 'agent_config', 'api_doc', 'db_schema']
  const colors = ['accent', 'success', 'warning', 'danger', 'accent']

  return toolTypes.map((toolType, index) => ({
    tool_type: toolType,
    color: colors[index],
    usage_count: usage[toolType] ?? 0,
  }))
})

const dashboardTools = computed(() => {
  if (toolStore.tools.length) return toolStore.tools
  return fallbackTools.value
})

// Tool icon mapping
const toolIconMap: Record<string, any> = {
  ui_review: Eye,
  project_doctor: Stethoscope,
  agent_config: Bot,
  api_doc: FileText,
  db_schema: Database,
}

function getToolIcon(toolType: ToolType | string) {
  return toolIconMap[toolType] || Wrench
}

function getToolDisplay(toolType: ToolType) {
  return getToolDisplayMeta(toolType)
}

function getToolTone(color: string) {
  return toolToneMap[color] || toolToneMap.accent
}

function formatRelativeTime(dateString: string): string {
  const date = new Date(dateString)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()

  const minutes = Math.floor(diffMs / 60_000)
  const hours = Math.floor(diffMs / 3_600_000)
  const days = Math.floor(diffMs / 86_400_000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`

  return dateString
}

function getScoreColorClass(score: number | null): string {
  if (score === null) return 'bg-surface-muted text-text-muted'
  if (score >= 80) return 'bg-success/10 text-success'
  if (score >= 60) return 'bg-warning/10 text-warning'
  return 'bg-danger/10 text-danger'
}
</script>

<template>
  <div class="space-y-8">
    <section class="dashboard-overview overflow-hidden rounded-lg border border-border bg-surface">
      <div class="grid gap-6 p-5 md:p-6 lg:grid-cols-[1fr_340px] lg:items-stretch">
        <div class="min-w-0">
          <div class="mb-5 inline-flex items-center gap-2 rounded-full border border-accent/15 bg-accent-soft px-3 py-1 text-sm font-semibold text-accent">
            <Sparkles :size="15" />
            <span>AI Developer Workbench</span>
          </div>
          <h2 class="text-2xl font-bold leading-tight text-text-primary md:text-3xl">
            选择工具，完成一次清晰的项目分析
          </h2>
          <p class="mt-3 max-w-3xl leading-relaxed text-text-secondary">
            聚合 UI、工程结构、Agent 配置、API 文档和数据库审查，把分析入口、报告记录和系统状态放在同一个工作台里。
          </p>

          <div class="mt-6 grid grid-cols-2 gap-3 lg:grid-cols-4">
            <div
              v-for="stat in stats"
              :key="stat.label"
              class="rounded-lg border border-border/80 bg-background/70 px-4 py-4 max-sm:px-3 max-sm:py-3"
            >
              <div class="mb-3 flex items-center justify-between gap-3">
                <span class="text-sm font-medium text-text-muted">{{ stat.label }}</span>
                <component :is="stat.icon" :size="18" :class="getToolTone(stat.color).text" />
              </div>
              <div class="truncate text-2xl font-bold text-text-primary max-sm:text-xl">
                {{ stat.value }}<span v-if="stat.suffix" class="ml-1 text-base font-semibold text-text-muted">{{ stat.suffix }}</span>
              </div>
            </div>
          </div>
        </div>

        <aside class="rounded-lg border border-border/80 bg-background/75 p-5 max-sm:p-4">
          <div class="mb-4 flex items-center justify-between gap-3">
            <div>
              <h3 class="font-semibold text-text-primary">当前状态</h3>
              <p class="mt-1 text-sm text-text-muted">服务、模型和最近使用概览</p>
            </div>
            <span
              :class="[
                'inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-sm font-semibold',
                systemStore.status?.healthy ? 'bg-success/10 text-success' : 'bg-danger/10 text-danger',
              ]"
            >
              <CheckCircle2 :size="15" />
              {{ systemStore.status?.healthy ? '正常' : '异常' }}
            </span>
          </div>

          <div class="grid gap-3 text-sm sm:grid-cols-3 lg:grid-cols-1">
            <div class="flex items-center justify-between gap-4 rounded-md bg-surface px-3 py-3 max-sm:py-2.5">
              <span class="text-text-muted">模型</span>
              <span class="truncate font-semibold text-text-primary">{{ systemStore.providerInfo }}</span>
            </div>
            <div
              :class="[
                'flex items-center justify-between gap-4 rounded-md px-3 py-3 max-sm:py-2.5',
                systemStore.isMockMode
                  ? 'bg-amber-50 border border-amber-200 dark:bg-amber-900/20 dark:border-amber-700/40'
                  : 'bg-emerald-50 border border-emerald-200 dark:bg-emerald-900/20 dark:border-emerald-700/40',
              ]"
            >
              <div class="flex items-center gap-2">
                <FlaskConical
                  :size="15"
                  :class="systemStore.isMockMode ? 'text-amber-600 dark:text-amber-400' : 'text-emerald-600 dark:text-emerald-400'"
                />
                <span class="text-text-muted">模式</span>
              </div>
              <span
                :class="[
                  'font-semibold text-sm',
                  systemStore.isMockMode
                    ? 'text-amber-800 dark:text-amber-200'
                    : 'text-emerald-800 dark:text-emerald-200',
                ]"
              >
                {{ systemStore.isMockMode ? '演示' : '真实 AI' }}
              </span>
            </div>
            <div class="flex items-center justify-between gap-4 rounded-md bg-surface px-3 py-3 max-sm:py-2.5">
              <span class="text-text-muted">总使用</span>
              <span class="font-semibold text-text-primary">{{ totalToolUsage }} 次</span>
            </div>
            <div class="flex items-center justify-between gap-4 rounded-md bg-surface px-3 py-3 max-sm:py-2.5">
              <span class="text-text-muted">常用工具</span>
            <span class="truncate font-semibold text-text-primary">
                {{ mostUsedTool ? getToolDisplay(mostUsedTool.tool_type).name : '暂无' }}
              </span>
            </div>
          </div>
        </aside>
      </div>
    </section>

    <div class="grid gap-8 xl:grid-cols-[minmax(0,1fr)_360px]">
      <section>
        <div class="mb-5 flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h2 class="text-2xl font-bold text-text-primary">分析工具</h2>
            <p class="mt-1 text-sm text-text-muted">选择一个入口开始分析，结果会自动沉淀为报告。</p>
          </div>
          <span class="text-sm font-medium text-text-muted">{{ dashboardTools.length }} 个工具可用</span>
        </div>

        <div class="grid gap-4 md:grid-cols-2">
          <RouterLink
            v-for="tool in dashboardTools"
            :key="tool.tool_type"
            :to="`/tools/${tool.tool_type.replace('_', '-')}`"
            :class="[
              'group flex min-h-[190px] flex-col rounded-lg border border-border bg-surface p-5 transition-smooth hover:-translate-y-0.5',
              getToolTone(tool.color).card,
            ]"
          >
            <div class="mb-4 flex items-start justify-between gap-4">
              <div class="flex items-start gap-4">
                <div :class="['flex h-11 w-11 flex-shrink-0 items-center justify-center rounded-lg border', getToolTone(tool.color).icon]">
                  <component :is="getToolIcon(tool.tool_type)" :size="22" />
                </div>
                <div class="min-w-0">
                  <h3 class="text-lg font-bold text-text-primary">{{ getToolDisplay(tool.tool_type).name }}</h3>
                  <p class="mt-1 text-sm font-medium text-text-muted">
                    {{ getToolDisplay(tool.tool_type).shortDescription }}
                  </p>
                </div>
              </div>
              <ArrowRight :size="20" class="mt-1 flex-shrink-0 text-text-muted transition-smooth group-hover:translate-x-1 group-hover:text-accent" />
            </div>

            <p class="flex-1 leading-relaxed text-text-secondary">
              {{ getToolDisplay(tool.tool_type).description }}
            </p>

            <div class="mt-5 flex items-center justify-between border-t border-border pt-4 text-sm">
              <span class="text-text-muted">{{ tool.usage_count }} 次使用</span>
              <span :class="['rounded-full border px-2.5 py-1 font-semibold', getToolTone(tool.color).soft, getToolTone(tool.color).border, getToolTone(tool.color).text]">
                开始分析
              </span>
            </div>
          </RouterLink>
        </div>
      </section>

      <aside class="space-y-6">
        <section class="rounded-lg border border-border bg-surface p-5">
          <div class="mb-5 flex items-center justify-between gap-3">
            <div>
              <h2 class="text-xl font-bold text-text-primary">最近报告</h2>
              <p class="mt-1 text-sm text-text-muted">最近 4 条分析结果</p>
            </div>
            <RouterLink to="/reports" class="inline-flex items-center gap-1 text-sm font-semibold text-accent hover:underline">
              查看全部
              <ArrowRight :size="15" />
            </RouterLink>
          </div>

          <div v-if="systemStore.statsLoading" class="rounded-lg border border-dashed border-border py-8 text-center text-text-muted">
            加载中...
          </div>

          <div v-else-if="!recentReports.length" class="rounded-lg border border-dashed border-border py-8 text-center text-text-muted">
            暂无报告
          </div>

          <div v-else class="space-y-3">
            <RouterLink
              v-for="report in recentReports"
              :key="report.id"
              :to="`/reports/${report.id}`"
              class="group block rounded-lg border border-border bg-background/70 px-4 py-3 transition-smooth hover:border-accent/35 hover:bg-surface"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="flex min-w-0 items-start gap-3">
                  <div class="flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent">
                    <component :is="getToolIcon(report.tool_type)" :size="17" />
                  </div>
                  <div class="min-w-0">
                    <h3 class="truncate font-semibold text-text-primary">{{ report.title }}</h3>
                    <p class="mt-1 truncate text-sm text-text-muted">
                      {{ getToolDisplay(report.tool_type).name }} · {{ formatRelativeTime(report.created_at) }}
                    </p>
                  </div>
                </div>
                <ChevronRight :size="18" class="mt-1 flex-shrink-0 text-text-muted transition-smooth group-hover:text-accent" />
              </div>
              <div class="mt-3 flex items-center justify-between gap-3">
                <p class="line-clamp-1 min-w-0 text-sm text-text-secondary">{{ report.summary }}</p>
                <span :class="['flex-shrink-0 rounded-full px-2.5 py-1 text-sm font-semibold', getScoreColorClass(report.total_score)]">
                  {{ report.total_score !== null ? report.total_score + '分' : '暂无' }}
                </span>
              </div>
            </RouterLink>
          </div>
        </section>

        <section class="rounded-lg border border-border bg-surface p-5">
          <h2 class="text-xl font-bold text-text-primary">建议流程</h2>
          <div class="mt-4 space-y-3">
            <div class="flex gap-3">
              <span class="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-md bg-accent-soft text-sm font-bold text-accent">1</span>
              <p class="text-sm leading-relaxed text-text-secondary">先用项目诊断看结构风险，再进入单项工具处理具体问题。</p>
            </div>
            <div class="flex gap-3">
              <span class="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-md bg-accent-soft text-sm font-bold text-accent">2</span>
              <p class="text-sm leading-relaxed text-text-secondary">把报告中的 Prompt 或建议复制到 AI Coding 工具继续迭代。</p>
            </div>
          </div>
        </section>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.dashboard-overview {
  background:
    linear-gradient(135deg, rgba(37, 99, 235, 0.06), transparent 38%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(247, 247, 245, 0.92));
}
</style>

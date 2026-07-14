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
  Wrench, Eye, Stethoscope, Bot, FileText, Database,
  ArrowRight, FileStack, BarChart3, Activity,
  CheckCircle2,
  FolderPlus, GitBranch, FileCode, Monitor, Package, WandSparkles, Search,
} from '@lucide/vue'
import QualityTrend from '@/components/dashboard/QualityTrend.vue'

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
    note: '累计报告',
  },
  {
    label: '平均评分',
    value: systemStore.dashboardStats?.average_score ?? '暂无',
    suffix: '/100',
    icon: BarChart3,
    color: 'success',
    note: '整体质量',
  },
  {
    label: '本周分析',
    value: systemStore.dashboardStats?.weekly_stats?.report_count_this_week ?? 0,
    suffix: '',
    icon: Activity,
    color: 'warning',
    note: '最近 7 天',
  },
  {
    label: '本周高危问题',
    value: systemStore.dashboardStats?.weekly_stats?.high_severity_count_this_week ?? 0,
    suffix: '',
    icon: CheckCircle2,
    color: 'danger',
    note: '最近 7 天',
  },
])

const weekly = computed(() => systemStore.dashboardStats?.weekly_stats ?? null)
const trendPoints = computed(() => systemStore.dashboardStats?.quality_trend ?? [])

const recentReports = computed(() => systemStore.dashboardStats?.recent_reports?.slice(0, 4) ?? [])

const totalToolUsage = computed(() => {
  const usage = systemStore.dashboardStats?.tool_usage ?? {}
  return (Object.values(usage) as number[]).reduce((sum, count) => sum + count, 0)
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

const maxToolUsage = computed(() => Math.max(1, ...dashboardTools.value.map((tool) => tool.usage_count)))

function getToolUsagePercent(count: number) {
  if (count <= 0) return 0
  return Math.max(8, Math.round((count / maxToolUsage.value) * 100))
}
const generationSteps = [
  { number: '01', title: '创建项目', description: '选择项目类型与技术栈', icon: FolderPlus },
  { number: '02', title: '填写需求', description: '明确目标用户与核心功能', icon: FileText },
  { number: '03', title: '确认蓝图', description: '评审页面结构与技术方案', icon: GitBranch },
  { number: '04', title: '生成代码', description: '实时查看进度与生成日志', icon: FileCode },
  { number: '05', title: '在线预览', description: '构建并检查运行效果', icon: Monitor },
  { number: '06', title: '交付文件', description: '浏览源码并导出完整项目', icon: Package },
]

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
  <div class="space-y-6">
    <header class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-[-0.025em] text-text-primary">工作台</h1>
        <p class="mt-1 text-sm text-text-muted">项目生成与质量数据概览</p>
      </div>
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
        <span
          :class="[
            'inline-flex min-h-10 items-center justify-center gap-2 rounded-lg border px-3 text-sm font-medium',
            systemStore.status?.healthy
              ? 'border-success/20 bg-success/10 text-success'
              : 'border-danger/20 bg-danger/10 text-danger',
          ]"
        >
          <span class="h-2 w-2 rounded-full" :class="systemStore.status?.healthy ? 'bg-success' : 'bg-danger'"></span>
          {{ systemStore.status?.healthy ? '服务正常' : '服务异常' }}
          <span class="text-current/60">·</span>
          <span class="max-w-32 truncate">{{ systemStore.providerInfo }}</span>
        </span>
        <RouterLink to="/tools/project-doctor" class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-semibold text-text-primary transition-smooth hover:border-accent/30 hover:bg-accent-soft">
          <Search :size="16" />分析项目
        </RouterLink>
        <RouterLink to="/projects/new" class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/90">
          <FolderPlus :size="16" />新建项目
        </RouterLink>
      </div>
    </header>

    <section class="grid grid-cols-2 gap-3 lg:grid-cols-4" aria-label="核心指标">
      <div v-for="stat in stats" :key="stat.label" class="metric-card rounded-xl border border-border bg-surface p-4 shadow-sm md:p-5">
        <div class="flex items-start justify-between gap-3">
          <div>
            <p class="text-sm font-medium text-text-muted">{{ stat.label }}</p>
            <p class="mt-3 text-2xl font-bold tracking-[-0.03em] text-text-primary md:text-3xl">
              {{ stat.value }}<span v-if="stat.suffix" class="ml-1 text-sm font-semibold text-text-muted">{{ stat.suffix }}</span>
            </p>
          </div>
          <span :class="['flex h-9 w-9 items-center justify-center rounded-lg border', getToolTone(stat.color).icon]">
            <component :is="stat.icon" :size="18" />
          </span>
        </div>
        <p class="mt-3 text-xs text-text-muted">{{ stat.note }}</p>
      </div>
    </section>

    <div class="grid gap-5 xl:grid-cols-[minmax(0,1fr)_340px]">
      <QualityTrend :points="trendPoints" />

      <section class="rounded-lg border border-border bg-surface p-5">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-text-primary">本周数据</h2>
            <p class="mt-1 text-xs text-text-muted">最近 7 天</p>
          </div>
          <Activity :size="18" class="text-accent" />
        </div>

        <div v-if="weekly" class="mt-5 grid grid-cols-3 gap-2 xl:grid-cols-1">
          <div class="rounded-lg bg-background px-3 py-3 xl:flex xl:items-center xl:justify-between">
            <span class="text-xs text-text-muted">分析次数</span>
            <strong class="mt-1 block text-lg text-text-primary xl:mt-0">{{ weekly.report_count_this_week }}</strong>
          </div>
          <div class="rounded-lg bg-background px-3 py-3 xl:flex xl:items-center xl:justify-between">
            <span class="text-xs text-text-muted">平均分</span>
            <strong class="mt-1 block text-lg text-text-primary xl:mt-0">{{ weekly.average_score_this_week !== null ? weekly.average_score_this_week.toFixed(1) : '—' }}</strong>
          </div>
          <div class="rounded-lg bg-background px-3 py-3 xl:flex xl:items-center xl:justify-between">
            <span class="text-xs text-text-muted">高危问题</span>
            <strong class="mt-1 block text-lg xl:mt-0" :class="weekly.high_severity_count_this_week > 0 ? 'text-danger' : 'text-success'">{{ weekly.high_severity_count_this_week }}</strong>
          </div>
        </div>
        <div v-else class="mt-5 rounded-lg border border-dashed border-border py-8 text-center text-sm text-text-muted">暂无本周数据</div>

        <div class="mt-4 flex items-center justify-between border-t border-border pt-4 text-sm">
          <span class="text-text-muted">最常用工具</span>
          <span class="font-semibold text-text-primary">{{ weekly?.most_used_tool_this_week ? getToolDisplayMeta(weekly.most_used_tool_this_week as ToolType)?.name : '暂无' }}</span>
        </div>
      </section>
    </div>

    <section class="generation-workflow relative overflow-hidden rounded-xl border border-slate-800 bg-[#111827] text-white">
      <div class="generation-workflow__grid absolute inset-0" aria-hidden="true"></div>
      <div class="generation-workflow__glow" aria-hidden="true"></div>
      <div class="relative p-5 md:p-6">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <div class="flex items-center gap-2">
              <WandSparkles :size="17" class="text-blue-300" />
              <h2 class="text-lg font-semibold">生成工作流</h2>
              <span class="rounded-full border border-white/10 bg-white/[0.05] px-2 py-0.5 text-xs text-slate-400">6 stages</span>
            </div>
            <p class="mt-1 text-sm text-slate-400">需求 → 蓝图 → 代码 → 预览 → 交付</p>
          </div>
          <RouterLink to="/projects/new" class="inline-flex min-h-10 items-center justify-center gap-2 rounded-lg bg-white px-4 text-sm font-semibold text-slate-950 transition-smooth hover:bg-blue-100">
            开始生成<ArrowRight :size="15" />
          </RouterLink>
        </div>

        <ol class="mt-5 grid grid-cols-3 gap-2 lg:grid-cols-6">
          <li v-for="step in generationSteps" :key="step.number" class="generation-step rounded-lg border border-white/10 bg-white/[0.045] px-3 py-3">
            <div class="flex items-center justify-between gap-2">
              <component :is="step.icon" :size="16" class="text-blue-200" />
              <span class="font-mono text-[10px] text-slate-600">{{ step.number }}</span>
            </div>
            <p class="mt-3 text-xs font-semibold text-slate-100">{{ step.title }}</p>
          </li>
        </ol>
      </div>
    </section>

    <div class="grid gap-5 xl:grid-cols-[minmax(0,1fr)_420px]">
      <section class="rounded-xl border border-border bg-surface p-5">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-text-primary">工具使用分布</h2>
            <p class="mt-1 text-xs text-text-muted">累计 {{ totalToolUsage }} 次</p>
          </div>
          <RouterLink to="/tools/project-doctor" class="text-sm font-semibold text-accent hover:underline">运行分析</RouterLink>
        </div>

        <div class="mt-5 space-y-4">
          <RouterLink v-for="tool in dashboardTools" :key="tool.tool_type" :to="`/tools/${tool.tool_type.replace('_', '-')}`" class="group grid grid-cols-[minmax(120px,180px)_1fr_auto] items-center gap-3">
            <div class="flex min-w-0 items-center gap-2.5">
              <span :class="['flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg border', getToolTone(tool.color).icon]"><component :is="getToolIcon(tool.tool_type)" :size="15" /></span>
              <span class="truncate text-sm font-medium text-text-primary group-hover:text-accent">{{ getToolDisplay(tool.tool_type).name }}</span>
            </div>
            <div class="h-2 overflow-hidden rounded-full bg-surface-muted">
              <div class="h-full rounded-full bg-accent transition-all duration-300" :style="{ width: `${getToolUsagePercent(tool.usage_count)}%` }"></div>
            </div>
            <span class="w-10 text-right text-sm font-semibold text-text-secondary">{{ tool.usage_count }}</span>
          </RouterLink>
        </div>
      </section>

      <section class="rounded-xl border border-border bg-surface p-5">
        <div class="flex items-center justify-between gap-3">
          <div>
            <h2 class="text-lg font-semibold text-text-primary">最近报告</h2>
            <p class="mt-1 text-xs text-text-muted">最新 4 条结果</p>
          </div>
          <RouterLink to="/reports" class="text-sm font-semibold text-accent hover:underline">全部</RouterLink>
        </div>

        <div v-if="systemStore.statsLoading" class="mt-5 rounded-lg border border-dashed border-border py-8 text-center text-sm text-text-muted">加载中...</div>
        <div v-else-if="!recentReports.length" class="mt-5 rounded-lg border border-dashed border-border py-8 text-center text-sm text-text-muted">暂无报告</div>
        <div v-else class="mt-4 divide-y divide-border">
          <RouterLink v-for="report in recentReports" :key="report.id" :to="`/reports/${report.id}`" class="group flex items-center gap-3 py-3 first:pt-0 last:pb-0">
            <span class="flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent"><component :is="getToolIcon(report.tool_type)" :size="15" /></span>
            <span class="min-w-0 flex-1">
              <span class="block truncate text-sm font-medium text-text-primary group-hover:text-accent">{{ report.title }}</span>
              <span class="mt-0.5 block text-xs text-text-muted">{{ getToolDisplay(report.tool_type).name }} · {{ formatRelativeTime(report.created_at) }}</span>
            </span>
            <span :class="['rounded-full px-2 py-1 text-xs font-semibold', getScoreColorClass(report.total_score)]">{{ report.total_score !== null ? report.total_score : '—' }}</span>
          </RouterLink>
        </div>
      </section>
    </div>
  </div>
</template>
<style scoped>
.metric-card {
  transition: border-color 180ms ease, box-shadow 180ms ease, transform 180ms ease;
}

.metric-card:hover {
  transform: translateY(-2px);
  border-color: rgba(37, 99, 235, 0.22);
  box-shadow: 0 12px 28px rgba(15, 23, 42, 0.07);
}

.generation-workflow {
  isolation: isolate;
}

.generation-workflow__grid {
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.045) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.045) 1px, transparent 1px);
  background-size: 34px 34px;
  mask-image: linear-gradient(135deg, black, transparent 86%);
}

.generation-workflow__glow {
  position: absolute;
  top: -140px;
  right: -80px;
  width: 380px;
  height: 380px;
  border-radius: 9999px;
  background: rgba(37, 99, 235, 0.24);
  filter: blur(80px);
  pointer-events: none;
}

.generation-step {
  transition: transform 200ms ease, border-color 200ms ease, background-color 200ms ease;
}

.generation-step:hover {
  transform: translateY(-3px);
  border-color: rgba(125, 211, 252, 0.28);
  background: rgba(255, 255, 255, 0.07);
}

@media (prefers-reduced-motion: reduce) {
  .generation-step { transition: none; }
  .generation-step:hover { transform: none; }
}
</style>

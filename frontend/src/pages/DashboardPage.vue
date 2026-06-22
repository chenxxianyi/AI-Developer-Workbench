<script setup lang="ts">
/**
 * Dashboard Page
 * Stats, tool grid, and recent reports
 */

import { onMounted } from 'vue'
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
} from '@lucide/vue'

const systemStore = useSystemStore()
const toolStore = useToolStore()

onMounted(async () => {
  await systemStore.fetchDashboardStats()
  await systemStore.fetchStatus()
  await toolStore.fetchTools()
})

// Stats cards data
const stats = [
  {
    label: '总报告数',
    value: systemStore.dashboardStats?.total_reports ?? 0,
    icon: 'file-stack',
    color: 'accent',
  },
  {
    label: '平均评分',
    value: systemStore.dashboardStats?.average_score ?? '暂无',
    suffix: '/100',
    icon: 'bar-chart-2',
    color: 'success',
  },
  {
    label: '工具使用',
    value: Object.keys(systemStore.dashboardStats?.tool_usage ?? {}).length,
    suffix: '/5',
    icon: 'wrench',
    color: 'warning',
  },
  {
    label: '最近活动',
    value: '2小时前',
    suffix: '',
    icon: 'activity',
    color: 'danger',
  },
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
  <div>
    <!-- Stats Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
      <div
        v-for="stat in stats"
        :key="stat.label"
        class="p-6 bg-surface border border-border rounded-lg"
      >
        <div class="flex items-center gap-3 mb-3">
          <component :is="stat.icon" :size="20" :class="`text-${stat.color}`" />
          <span class="text-sm text-text-muted">{{ stat.label }}</span>
        </div>
        <div class="text-3xl font-bold text-text-primary">
          {{ stat.value
          }}<span v-if="stat.suffix" class="text-lg text-text-muted">{{ stat.suffix }}</span>
        </div>
      </div>
    </div>

    <!-- Tools Grid -->
    <div class="mb-8">
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-text-primary">分析工具</h2>
        <span class="text-sm text-text-muted">选择工具开始分析</span>
      </div>

      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <RouterLink
          v-for="tool in toolStore.tools"
          :key="tool.tool_type"
          :to="`/tools/${tool.tool_type.replace('_', '-')}`"
          class="group p-6 bg-surface border border-border rounded-lg hover:border-accent hover:shadow-md transition-smooth"
        >
          <div class="flex items-start gap-4 mb-4">
            <div
              :class="[
                'w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0',
                `bg-${tool.color}-soft`,
              ]"
            >
              <component :is="getToolIcon(tool.tool_type)" :size="24" :class="`text-${tool.color}`" />
            </div>
            <div class="flex-1">
              <h3 class="text-xl font-semibold mb-1 text-text-primary">
                {{ getToolDisplay(tool.tool_type).name }}
              </h3>
              <p class="text-sm text-text-muted">
                {{ getToolDisplay(tool.tool_type).shortDescription }}
              </p>
            </div>
          </div>
          <p class="text-text-secondary mb-4 leading-relaxed">
            {{ getToolDisplay(tool.tool_type).description }}
          </p>
          <div class="flex items-center justify-between pt-4 border-t border-border">
            <div class="flex items-center gap-2 text-sm text-text-muted">
              <component :is="getToolIcon(tool.tool_type)" :size="16" />
              <span>{{ tool.usage_count }} 次使用</span>
            </div>
            <ArrowRight :size="20" class="text-text-muted group-hover:text-accent transition-smooth" />
          </div>
        </RouterLink>

        <!-- Empty slot for grid alignment -->
        <div class="hidden lg:block"></div>
      </div>
    </div>

    <!-- Recent Reports -->
    <div>
      <div class="flex items-center justify-between mb-6">
        <h2 class="text-2xl font-bold text-text-primary">最近报告</h2>
        <RouterLink
          to="/reports"
          class="flex items-center gap-2 text-accent hover:underline transition-smooth"
        >
          <span>查看全部</span>
          <ArrowRight :size="16" />
        </RouterLink>
      </div>

      <div v-if="systemStore.statsLoading" class="text-center py-8">
        <div class="text-text-muted">加载中...</div>
      </div>

      <div v-else-if="!systemStore.dashboardStats?.recent_reports?.length" class="text-center py-8">
        <div class="text-text-muted">暂无报告</div>
      </div>

      <div v-else class="space-y-4">
        <RouterLink
          v-for="report in systemStore.dashboardStats?.recent_reports?.slice(0, 3)"
          :key="report.id"
          :to="`/reports/${report.id}`"
          class="block p-4 bg-surface border border-border rounded-lg hover:border-accent hover:shadow-sm transition-smooth"
        >
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-3">
              <div
                :class="[
                  'w-8 h-8 rounded flex items-center justify-center',
                  `bg-${getToolIcon(report.tool_type) ? 'accent-soft' : 'surface-muted'}`,
                ]"
              >
                <component :is="getToolIcon(report.tool_type)" :size="16" class="text-accent" />
              </div>
              <div>
                <h3 class="font-semibold text-text-primary">{{ report.title }}</h3>
                <p class="text-sm text-text-muted">
                  {{ report.tool_type === 'ui_review' ? 'UI 审查' :
                     report.tool_type === 'project_doctor' ? '项目诊断' :
                     report.tool_type === 'agent_config' ? 'Agent 配置' :
                     report.tool_type === 'api_doc' ? 'API 文档' :
                     report.tool_type === 'db_schema' ? '数据库审查' : report.tool_type
                  }} · {{ formatRelativeTime(report.created_at) }}
                </p>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <div
                :class="['px-3 py-1 rounded-full text-sm font-medium', getScoreColorClass(report.total_score)]"
              >
                {{ report.total_score !== null ? report.total_score + '分' : '暂无评分'
                }}
              </div>
              <ChevronRight :size="20" class="text-text-muted" />
            </div>
          </div>
          <p class="text-sm text-text-secondary">{{ report.summary }}</p>
        </RouterLink>
      </div>
    </div>
  </div>
</template>

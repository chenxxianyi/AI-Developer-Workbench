<script setup lang="ts">
/**
 * ReportStatusBadge — shows report status with distinct colors and accessible icons.
 */
import type { ReportStatus } from '@/types/report'
import { Clock, CheckCircle2, AlertTriangle, XCircle } from '@lucide/vue'

defineProps<{
  status: ReportStatus
}>()

const statusConfig: Record<ReportStatus, { label: string; icon: any; classes: string }> = {
  processing: {
    label: '处理中',
    icon: Clock,
    classes: 'bg-blue-50 text-blue-700 border-blue-200 dark:bg-blue-900/20 dark:text-blue-300 dark:border-blue-700/40',
  },
  succeeded: {
    label: '成功',
    icon: CheckCircle2,
    classes: 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-900/20 dark:text-emerald-300 dark:border-emerald-700/40',
  },
  fallback: {
    label: '降级',
    icon: AlertTriangle,
    classes: 'bg-amber-50 text-amber-700 border-amber-200 dark:bg-amber-900/20 dark:text-amber-300 dark:border-amber-700/40',
  },
  failed: {
    label: '失败',
    icon: XCircle,
    classes: 'bg-red-50 text-red-700 border-red-200 dark:bg-red-900/20 dark:text-red-300 dark:border-red-700/40',
  },
}
</script>

<template>
  <span
    :class="[
      'inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-semibold',
      statusConfig[status].classes,
    ]"
    role="status"
  >
    <component :is="statusConfig[status].icon" :size="13" aria-hidden="true" />
    {{ statusConfig[status].label }}
  </span>
</template>

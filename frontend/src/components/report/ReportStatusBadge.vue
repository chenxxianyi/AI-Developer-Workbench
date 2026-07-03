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
    classes: 'bg-blue-600 text-white border-blue-600 dark:bg-blue-500 dark:text-white dark:border-blue-500',
  },
  succeeded: {
    label: '成功',
    icon: CheckCircle2,
    classes: 'bg-emerald-600 text-white border-emerald-600 dark:bg-emerald-500 dark:text-white dark:border-emerald-500',
  },
  fallback: {
    label: '降级',
    icon: AlertTriangle,
    classes: 'bg-amber-600 text-white border-amber-600 dark:bg-amber-500 dark:text-white dark:border-amber-500',
  },
  failed: {
    label: '失败',
    icon: XCircle,
    classes: 'bg-red-600 text-white border-red-600 dark:bg-red-500 dark:text-white dark:border-red-500',
  },
}
</script>

<template>
  <span
    :class="[
      'inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-bold shadow-sm',
      statusConfig[status].classes,
    ]"
    role="status"
  >
    <component :is="statusConfig[status].icon" :size="13" aria-hidden="true" />
    {{ statusConfig[status].label }}
  </span>
</template>

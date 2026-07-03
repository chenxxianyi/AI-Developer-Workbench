<script setup lang="ts">
/**
 * ReportListCard — single report row with status, score, tool type, and link.
 */
import type { Report } from '@/types/report'
import { getToolDisplayMeta } from '@/utils/toolDisplay'
import { ChevronRight } from '@lucide/vue'
import ReportStatusBadge from './ReportStatusBadge.vue'

defineProps<{
  report: Report
}>()

function getScoreBadgeClass(score: number): string {
  if (score >= 80) {
    return 'border-emerald-600 bg-emerald-600 text-white dark:border-emerald-500 dark:bg-emerald-500'
  }
  if (score >= 60) {
    return 'border-amber-600 bg-amber-600 text-white dark:border-amber-500 dark:bg-amber-500'
  }
  return 'border-red-600 bg-red-600 text-white dark:border-red-500 dark:bg-red-500'
}
</script>

<template>
  <RouterLink
    :to="`/reports/${report.id}`"
    class="group block rounded-lg border border-border bg-surface px-5 py-4 transition-smooth hover:border-accent/35 hover:bg-surface-hover focus-visible:ring-2 focus-visible:ring-accent focus-visible:outline-none"
  >
    <div class="flex items-start justify-between gap-4">
      <div class="min-w-0 flex-1">
        <div class="flex items-center gap-2 mb-1.5">
          <ReportStatusBadge :status="report.status" />
          <span class="text-xs font-medium text-text-muted uppercase tracking-wide">
            {{ getToolDisplayMeta(report.tool_type).name }}
          </span>
        </div>
        <h3 class="font-semibold text-text-primary truncate">{{ report.title }}</h3>
        <p class="mt-1 text-sm text-text-secondary line-clamp-1">{{ report.summary }}</p>
      </div>

      <div class="flex items-center gap-3 flex-shrink-0">
        <span
          v-if="report.total_score !== null"
          :class="[
            'inline-flex min-w-16 items-center justify-center rounded-lg border px-3 py-1.5 text-sm font-extrabold leading-none shadow-sm',
            getScoreBadgeClass(report.total_score),
          ]"
        >
          {{ report.total_score }}
          <span class="ml-1 border-l border-white/35 pl-1 text-xs font-bold text-white/95">{{ report.grade }}</span>
        </span>
        <span
          v-else
          class="inline-flex items-center rounded-lg border border-border bg-surface-muted px-3 py-1.5 text-xs font-semibold text-text-secondary"
        >
          无需评分
        </span>
        <ChevronRight :size="18" class="text-text-muted transition-smooth group-hover:text-accent group-hover:translate-x-0.5" />
      </div>
    </div>
  </RouterLink>
</template>

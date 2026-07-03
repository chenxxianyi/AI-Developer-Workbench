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
          class="inline-flex items-center justify-center rounded-md px-2.5 py-1 text-sm font-bold"
          :class="
            report.total_score >= 80
              ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300'
              : report.total_score >= 60
                ? 'bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300'
                : 'bg-red-50 text-red-700 dark:bg-red-900/20 dark:text-red-300'
          "
        >
          {{ report.total_score }}
          <span class="text-xs font-normal opacity-70 ml-0.5">{{ report.grade }}</span>
        </span>
        <span
          v-else
          class="text-xs font-medium text-text-muted px-2 py-1"
        >
          无需评分
        </span>
        <ChevronRight :size="18" class="text-text-muted transition-smooth group-hover:text-accent group-hover:translate-x-0.5" />
      </div>
    </div>
  </RouterLink>
</template>

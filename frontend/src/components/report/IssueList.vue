<script setup lang="ts">
import type { Issue } from '@/types/report'
import { AlertTriangle, AlertCircle, Info } from '@lucide/vue'

defineProps<{
  issues: Issue[]
}>()

const severityConfig = {
  high: { icon: AlertTriangle, classes: 'bg-red-50 border-red-200 text-red-800 dark:bg-red-900/20 dark:border-red-700/40 dark:text-red-200' },
  medium: { icon: AlertCircle, classes: 'bg-amber-50 border-amber-200 text-amber-800 dark:bg-amber-900/20 dark:border-amber-700/40 dark:text-amber-200' },
  low: { icon: Info, classes: 'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-900/20 dark:border-blue-700/40 dark:text-blue-200' },
}

const severityLabel = { high: '高', medium: '中', low: '低' }
</script>

<template>
  <div v-if="issues.length" class="rounded-lg border border-border bg-surface p-5">
    <h3 class="text-lg font-semibold text-text-primary mb-4">
      发现的问题
      <span class="text-sm font-normal text-text-muted ml-2">{{ issues.length }} 个</span>
    </h3>

    <div class="space-y-3">
      <div
        v-for="(issue, idx) in issues"
        :key="idx"
        :class="[
          'rounded-md border px-4 py-3',
          severityConfig[issue.severity]?.classes || severityConfig.medium.classes,
        ]"
      >
        <div class="flex items-start gap-2 mb-1">
          <component :is="severityConfig[issue.severity]?.icon || AlertCircle" :size="16" class="mt-0.5 flex-shrink-0" />
          <div>
            <div class="flex items-center gap-2">
              <span class="font-semibold text-sm">{{ issue.title }}</span>
              <span class="text-xs opacity-70 uppercase">{{ severityLabel[issue.severity] || '?' }}</span>
              <span v-if="issue.category" class="text-xs opacity-50">{{ issue.category }}</span>
            </div>
          </div>
        </div>
        <p class="mt-1 text-sm opacity-90">{{ issue.problem }}</p>
        <p v-if="issue.suggestion" class="mt-1 text-sm opacity-80">
          <span class="font-medium">建议：</span>{{ issue.suggestion }}
        </p>
        <p v-if="issue.action" class="mt-1 text-sm opacity-80">
          <span class="font-medium">行动：</span>{{ issue.action }}
        </p>
      </div>
    </div>
  </div>
</template>

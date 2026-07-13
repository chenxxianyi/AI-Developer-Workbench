<script setup lang="ts">
/**
 * ActionItemCard — one actionable remediation task from a report.
 */
import type { DisplayActionItem } from '@/utils/actionItems'
import { Check, Copy } from '@lucide/vue'

defineProps<{
  item: DisplayActionItem
  completed: boolean
  copied: boolean
}>()

const emit = defineEmits<{
  toggle: [id: string]
  copy: [item: DisplayActionItem]
}>()

const priorityMeta = {
  high: 'border-red-200 bg-red-50 text-red-700 dark:border-red-800/60 dark:bg-red-950/20 dark:text-red-300',
  medium: 'border-amber-200 bg-amber-50 text-amber-700 dark:border-amber-800/60 dark:bg-amber-950/20 dark:text-amber-300',
  low: 'border-slate-200 bg-slate-50 text-slate-700 dark:border-slate-700 dark:bg-slate-900/40 dark:text-slate-300',
}

const effortLabel = {
  small: '小',
  medium: '中',
  large: '大',
}
</script>

<template>
  <article class="rounded-lg border border-border bg-background/70 p-4">
    <div class="flex items-start gap-3">
      <label class="mt-1 inline-flex h-5 w-5 flex-shrink-0 cursor-pointer items-center justify-center">
        <input
          type="checkbox"
          class="sr-only"
          :checked="completed"
          @change="emit('toggle', item.id)"
        >
        <span
          :class="[
            'flex h-5 w-5 items-center justify-center rounded border transition-smooth',
            completed ? 'border-success bg-success text-white' : 'border-border bg-surface',
          ]"
          aria-hidden="true"
        >
          <Check v-if="completed" :size="14" />
        </span>
      </label>

      <div class="min-w-0 flex-1">
        <div class="mb-2 flex flex-wrap items-center gap-2">
          <h4
            :class="[
              'text-sm font-semibold text-text-primary',
              completed ? 'line-through opacity-70' : '',
            ]"
          >
            {{ item.title }}
          </h4>
          <span :class="['rounded-full border px-2 py-0.5 text-xs font-semibold', priorityMeta[item.priority]]">
            {{ item.priority }}
          </span>
          <span class="rounded-full border border-border bg-surface px-2 py-0.5 text-xs font-medium text-text-muted">
            {{ effortLabel[item.effort] }}工作量
          </span>
          <span v-if="item.legacy" class="rounded-full border border-border bg-surface px-2 py-0.5 text-xs font-medium text-text-muted">
            旧报告
          </span>
        </div>

        <p class="text-sm leading-relaxed text-text-secondary">{{ item.reason }}</p>
        <p class="mt-2 text-xs font-medium text-text-muted">{{ item.category }}</p>

        <div class="mt-3 rounded-md border border-border bg-surface p-3">
          <div class="mb-1 text-xs font-semibold text-text-muted">Prompt</div>
          <p class="whitespace-pre-wrap text-sm leading-relaxed text-text-secondary">{{ item.suggested_prompt }}</p>
        </div>
      </div>

      <button
        class="inline-flex flex-shrink-0 items-center gap-1.5 rounded-md border border-border bg-surface px-3 py-1.5 text-sm font-medium text-text-primary transition-smooth hover:bg-surface-hover"
        :aria-label="`复制行动项 Prompt：${item.title}`"
        @click="emit('copy', item)"
      >
        <Check v-if="copied" :size="15" class="text-success" />
        <Copy v-else :size="15" />
        <span class="hidden sm:inline">{{ copied ? '已复制' : '复制' }}</span>
      </button>
    </div>
  </article>
</template>

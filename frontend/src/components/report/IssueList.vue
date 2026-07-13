<script setup lang="ts">
import type { Issue } from '@/types/report'
import { AlertTriangle, AlertCircle, Info } from '@lucide/vue'

defineProps<{
  issues: Issue[]
}>()

const severityConfig = {
  high: {
    icon: AlertTriangle,
    cardClasses: 'border-border bg-surface hover:border-red-200 dark:hover:border-red-800/60',
    accentClasses: 'bg-red-500 dark:bg-red-400',
    iconFrameClasses: 'border-red-200 bg-surface text-red-600 dark:border-red-800/60 dark:text-red-300',
    iconClasses: 'text-red-600 dark:text-red-300',
    badgeClasses: 'border-red-600 bg-red-600 text-white dark:border-red-500 dark:bg-red-500 dark:text-white',
  },
  medium: {
    icon: AlertCircle,
    cardClasses: 'border-border bg-surface hover:border-amber-200 dark:hover:border-amber-800/60',
    accentClasses: 'bg-amber-500 dark:bg-amber-300',
    iconFrameClasses: 'border-amber-200 bg-surface text-amber-600 dark:border-amber-800/60 dark:text-amber-300',
    iconClasses: 'text-amber-600 dark:text-amber-300',
    badgeClasses: 'border-amber-600 bg-amber-600 text-white dark:border-amber-500 dark:bg-amber-500 dark:text-white',
  },
  low: {
    icon: Info,
    cardClasses: 'border-border bg-surface hover:border-blue-200 dark:hover:border-blue-800/60',
    accentClasses: 'bg-blue-500 dark:bg-blue-300',
    iconFrameClasses: 'border-blue-200 bg-surface text-blue-600 dark:border-blue-800/60 dark:text-blue-300',
    iconClasses: 'text-blue-600 dark:text-blue-300',
    badgeClasses: 'border-blue-600 bg-blue-600 text-white dark:border-blue-500 dark:bg-blue-500 dark:text-white',
  },
}

const severityLabel = { high: '高', medium: '中', low: '低' }
</script>

<template>
  <section v-if="issues.length" class="rounded-2xl border border-border bg-surface p-5 shadow-sm">
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <div>
        <h3 class="text-lg font-semibold text-text-primary">发现的问题</h3>
        <p class="mt-1 text-sm text-text-muted">按风险等级整理，优先处理高影响项</p>
      </div>
      <span class="inline-flex items-center rounded-full border border-border bg-surface-muted px-3 py-1 text-sm font-semibold text-text-secondary">
        {{ issues.length }} 个
      </span>
    </div>

    <div class="space-y-3">
      <article
        v-for="(issue, idx) in issues"
        :key="idx"
        :class="[
          'relative overflow-hidden rounded-xl border p-4 text-text-primary shadow-sm transition-smooth hover:shadow-md',
          severityConfig[issue.severity]?.cardClasses || severityConfig.medium.cardClasses,
        ]"
      >
        <span
          aria-hidden="true"
          :class="[
            'absolute inset-y-0 left-0 w-1.5',
            severityConfig[issue.severity]?.accentClasses || severityConfig.medium.accentClasses,
          ]"
        />

        <div class="flex gap-3 pl-2">
          <div
            :class="[
              'flex h-9 w-9 flex-shrink-0 items-center justify-center rounded-lg border',
              severityConfig[issue.severity]?.iconFrameClasses || severityConfig.medium.iconFrameClasses,
            ]"
          >
            <component
              :is="severityConfig[issue.severity]?.icon || AlertCircle"
              :size="17"
              :class="severityConfig[issue.severity]?.iconClasses || severityConfig.medium.iconClasses"
            />
          </div>

          <div class="min-w-0 flex-1">
            <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
              <h4 class="min-w-0 text-base font-semibold leading-snug text-text-primary">{{ issue.title }}</h4>
              <div class="flex flex-shrink-0 flex-wrap items-center gap-2">
                <span
                  :class="[
                    'inline-flex min-w-6 items-center justify-center rounded-full border px-2 py-0.5 text-xs font-bold leading-5 shadow-sm',
                    severityConfig[issue.severity]?.badgeClasses || severityConfig.medium.badgeClasses,
                  ]"
                >
                  {{ severityLabel[issue.severity] || '?' }}
                </span>
                <span
                  v-if="issue.category"
                  class="inline-flex items-center rounded-full border border-border bg-surface-muted px-2 py-0.5 text-xs font-semibold leading-5 text-text-secondary"
                >
                  {{ issue.category }}
                </span>
              </div>
            </div>

            <p class="mt-2 text-sm leading-relaxed text-text-secondary">{{ issue.problem }}</p>
            <div v-if="issue.viewport || issue.contrast_suggestion || issue.component_prompt" class="mt-3 space-y-2 rounded-lg border border-border/70 bg-surface-muted/60 px-3 py-2 text-sm text-text-secondary">
              <p v-if="issue.viewport"><strong class="text-text-primary">Viewport:</strong> {{ issue.viewport }}</p>
              <p v-if="issue.contrast_suggestion"><strong class="text-text-primary">Contrast:</strong> {{ issue.contrast_suggestion }}</p>
              <p v-if="issue.component_prompt"><strong class="text-text-primary">Component fix prompt:</strong> {{ issue.component_prompt }}</p>
            </div>

            <div
              v-if="issue.suggestion || issue.action"
              :class="[
                'mt-3 grid gap-2',
                issue.suggestion && issue.action ? 'lg:grid-cols-2' : 'lg:grid-cols-1',
              ]"
            >
              <div
                v-if="issue.suggestion"
                class="rounded-lg border border-border/70 bg-surface-muted/60 px-3 py-2"
              >
                <div class="mb-1 text-xs font-semibold text-text-primary">建议</div>
                <p class="text-sm leading-relaxed text-text-secondary">{{ issue.suggestion }}</p>
              </div>

              <div
                v-if="issue.action"
                class="rounded-lg border border-border/70 bg-surface-muted/60 px-3 py-2"
              >
                <div class="mb-1 text-xs font-semibold text-text-primary">行动</div>
                <p class="text-sm leading-relaxed text-text-secondary">{{ issue.action }}</p>
              </div>
            </div>
          </div>
        </div>
      </article>
    </div>
  </section>
</template>

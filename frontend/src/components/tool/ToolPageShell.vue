<script setup lang="ts">
/**
 * ToolPageShell — shared layout for the five tool pages.
 * Provides the page header, two-column grid (form | result), error block,
 * submit button, and loading/empty/result containers.
 */
import type { Component } from 'vue'
import { Loader2, AlertCircle } from '@lucide/vue'

const props = withDefaults(
  defineProps<{
    icon: Component
    title: string
    description: string
    stepText?: string
    loading?: boolean
    error?: string | null
    canSubmit?: boolean
    submitLabel?: string
    submittingLabel?: string
    /** Accent color key for the icon frame and submit button. */
    accent?: 'accent' | 'success' | 'orange' | 'teal' | 'purple'
    /** Back link label; omit to hide the back-to-dashboard button. */
    backLabel?: string
    /** Loading hint shown under the spinner. */
    loadingHint?: string
    /** Panel heading for the form column. */
    formHeading?: string
    /** Panel heading for the result column. */
    resultHeading?: string
    /** testid applied to the submit button. */
    submitTestid?: string
  }>(),
  {
    stepText: '',
    loading: false,
    error: null,
    canSubmit: true,
    submitLabel: '开始分析',
    submittingLabel: '分析中...',
    accent: 'accent',
    backLabel: '返回工作台',
    loadingHint: 'AI 正在分析...',
    formHeading: '输入参数',
    resultHeading: '分析结果',
  },
)

const emit = defineEmits<{
  submit: []
  back: []
}>()

const accentBg: Record<string, string> = {
  accent: 'bg-accent',
  success: 'bg-success',
  orange: 'bg-orange-500',
  teal: 'bg-teal-500',
  purple: 'bg-purple-500',
}

const accentRing: Record<string, string> = {
  accent: 'focus-visible:ring-accent',
  success: 'focus-visible:ring-success',
  orange: 'focus-visible:ring-orange-500',
  teal: 'focus-visible:ring-teal-500',
  purple: 'focus-visible:ring-purple-500',
}

const accentSubmit: Record<string, string> = {
  accent: 'bg-accent text-white hover:bg-accent/80',
  success: 'bg-success text-white hover:bg-success/80',
  orange: 'bg-orange-500 text-white hover:bg-orange-600',
  teal: 'bg-teal-500 text-white hover:bg-teal-600',
  purple: 'bg-purple-500 text-white hover:bg-purple-600',
}

const accentStep: Record<string, string> = {
  accent: 'border-accent/20 bg-accent-soft text-accent',
  success: 'border-success/20 bg-success/10 text-success',
  orange: 'border-orange-500/20 bg-orange-50 text-orange-600 dark:bg-orange-900/20 dark:text-orange-300',
  teal: 'border-teal-500/20 bg-teal-50 text-teal-600 dark:bg-teal-900/20 dark:text-teal-300',
  purple: 'border-purple-500/20 bg-purple-50 text-purple-600 dark:bg-purple-900/20 dark:text-purple-300',
}

function onSubmit() {
  if (props.loading || !props.canSubmit) return
  emit('submit')
}
</script>

<template>
  <div class="max-w-6xl mx-auto">
    <!-- Header -->
    <div class="mb-6">
      <button
        v-if="backLabel"
        type="button"
        class="flex items-center gap-2 text-text-secondary hover:text-text-primary transition-smooth mb-4 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none rounded"
        @click="emit('back')"
      >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M19 12H5M12 19l-7-7 7-7" />
        </svg>
        <span>{{ backLabel }}</span>
      </button>
      <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div class="flex items-start gap-3">
          <div
            :class="['w-12 h-12 rounded-xl flex items-center justify-center shrink-0 shadow-sm', accentBg[accent]]"
          >
            <component :is="icon" :size="24" class="text-white" />
          </div>
          <div>
            <h1 class="text-2xl font-bold tracking-tight text-text-primary sm:text-3xl">{{ title }}</h1>
            <p class="text-text-secondary mt-1 max-w-2xl">{{ description }}</p>
          </div>
        </div>

        <div
          v-if="stepText"
          :class="['inline-flex items-center rounded-full border px-4 py-2 text-sm font-medium', accentStep[accent]]"
        >
          {{ stepText }}
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Input Panel -->
      <div class="bg-surface border border-border rounded-xl p-6 shadow-sm">
        <div class="mb-5">
          <h2 class="text-lg font-semibold text-text-primary">{{ formHeading }}</h2>
        </div>

        <slot name="form" />

        <!-- Error -->
        <div v-if="error" class="mb-4 p-3 bg-danger/10 border border-danger/20 rounded-lg">
          <div class="flex items-center gap-2">
            <AlertCircle :size="18" class="text-danger" />
            <span class="text-danger">{{ error }}</span>
          </div>
        </div>

        <!-- Submit -->
        <div class="flex flex-col gap-3 border-t border-border pt-4 sm:flex-row sm:items-center sm:justify-between">
          <div class="flex gap-3 sm:justify-end sm:ml-auto">
            <button
              type="button"
              :data-testid="submitTestid || undefined"
              :disabled="loading || !canSubmit"
              :class="[
                'flex items-center gap-2 px-6 py-2.5 rounded-lg font-semibold shadow-sm transition-smooth focus-visible:ring-2 focus:outline-none',
                accentRing[accent],
                loading || !canSubmit
                  ? 'bg-surface-muted text-text-muted cursor-not-allowed shadow-none'
                  : `${accentSubmit[accent]} cursor-pointer`,
              ]"
              @click="onSubmit"
            >
              <Loader2 v-if="loading" :size="18" class="animate-spin" />
              <component :is="icon" v-else :size="18" />
              <span>{{ loading ? submittingLabel : submitLabel }}</span>
            </button>
            <slot name="actions" />
          </div>
        </div>
      </div>

      <!-- Result Panel -->
      <div class="bg-surface border border-border rounded-xl p-6 shadow-sm">
        <h2 class="text-lg font-semibold text-text-primary mb-4">{{ resultHeading }}</h2>

        <!-- Empty -->
        <div v-if="!loading && !$slots.result">
          <slot name="empty" />
        </div>
        <div v-else-if="!loading">
          <slot name="result" />
        </div>

        <!-- Loading -->
        <div v-if="loading" class="text-center py-12">
          <Loader2 :size="48" class="text-accent mx-auto mb-4 animate-spin" />
          <p class="text-text-secondary">{{ loadingHint }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

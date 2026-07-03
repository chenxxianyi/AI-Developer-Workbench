<script setup lang="ts">
/**
 * ReportErrorState — displays error or not-found states for report detail.
 */
import { RouterLink } from 'vue-router'
import { AlertTriangle, FileQuestion, ArrowLeft } from '@lucide/vue'

defineProps<{
  type: 'not-found' | 'error' | 'loading'
  message?: string
}>()
</script>

<template>
  <div class="rounded-lg border border-dashed border-border py-16 text-center">
    <FileQuestion v-if="type === 'not-found'" :size="48" class="mx-auto mb-4 text-text-muted" />
    <AlertTriangle v-else :size="48" class="mx-auto mb-4 text-warning" />

    <h2 class="text-xl font-semibold text-text-primary mb-2">
      {{ type === 'not-found' ? '报告未找到' : type === 'error' ? '加载失败' : '加载中...' }}
    </h2>
    <p class="text-sm text-text-muted mb-6">
      {{ message || (type === 'not-found' ? '该报告可能已被删除。' : '请稍后重试。') }}
    </p>

    <RouterLink
      v-if="type !== 'loading'"
      to="/reports"
      class="inline-flex items-center gap-2 rounded-md border border-border bg-surface px-5 py-2.5 text-sm font-semibold text-text-primary hover:bg-surface-hover transition-smooth"
    >
      <ArrowLeft :size="18" />
      返回报告列表
    </RouterLink>
  </div>
</template>

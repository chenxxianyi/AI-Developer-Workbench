<script setup lang="ts">
/**
 * PaginationBar — lightweight pagination with prev/next and page info.
 */
import { ChevronLeft, ChevronRight } from '@lucide/vue'

defineProps<{
  currentPage: number
  totalPages: number
  hasPrev: boolean
  hasNext: boolean
}>()

const emit = defineEmits<{
  prev: []
  next: []
}>()
</script>

<template>
  <div v-if="totalPages > 1" class="flex items-center justify-center gap-3 py-4">
    <button
      :disabled="!hasPrev"
      class="inline-flex items-center gap-1 rounded-md border border-border px-3 py-2 text-sm font-medium transition-smooth hover:bg-surface-hover disabled:opacity-40 disabled:cursor-not-allowed"
      @click="emit('prev')"
      aria-label="上一页"
    >
      <ChevronLeft :size="16" aria-hidden="true" />
      <span class="hidden sm:inline">上一页</span>
    </button>

    <span class="text-sm text-text-muted min-w-[80px] text-center">
      第 {{ currentPage }} / {{ totalPages }} 页
    </span>

    <button
      :disabled="!hasNext"
      class="inline-flex items-center gap-1 rounded-md border border-border px-3 py-2 text-sm font-medium transition-smooth hover:bg-surface-hover disabled:opacity-40 disabled:cursor-not-allowed"
      @click="emit('next')"
      aria-label="下一页"
    >
      <span class="hidden sm:inline">下一页</span>
      <ChevronRight :size="16" aria-hidden="true" />
    </button>
  </div>
</template>

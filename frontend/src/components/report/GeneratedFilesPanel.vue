<script setup lang="ts">
/**
 * GeneratedFilesPanel — lists generated files with download links.
 */
import type { GeneratedFileMeta } from '@/types/report'
import { Download, FileText, FileCode, FileJson } from '@lucide/vue'

defineProps<{
  files: GeneratedFileMeta[]
  reportId: string
}>()

function getFileIcon(filename: string) {
  const ext = filename.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'json': return FileJson
    case 'sql': return FileCode
    default: return FileText
  }
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
</script>

<template>
  <div v-if="files.length" class="rounded-lg border border-border bg-surface p-5">
    <h3 class="text-lg font-semibold text-text-primary mb-4">
      生成文件
      <span class="text-sm font-normal text-text-muted ml-2">{{ files.length }} 个</span>
    </h3>

    <div class="space-y-2">
      <div
        v-for="file in files"
        :key="file.id"
        class="flex items-center justify-between gap-3 rounded-md border border-border/80 bg-background/50 px-4 py-3"
      >
        <div class="flex items-center gap-3 min-w-0">
          <component :is="getFileIcon(file.filename)" :size="18" class="text-text-muted flex-shrink-0" />
          <div class="min-w-0">
            <span class="text-sm font-medium text-text-primary truncate block">{{ file.filename }}</span>
            <span class="text-xs text-text-muted">{{ formatSize(file.size_bytes) }}</span>
          </div>
        </div>
        <a
          :href="`/api/reports/${reportId}/files/${encodeURIComponent(file.filename)}`"
          class="inline-flex items-center gap-1.5 rounded-md border border-border bg-surface px-3 py-1.5 text-sm font-medium text-accent hover:bg-accent-soft transition-smooth flex-shrink-0"
          :download="file.filename"
          aria-label="下载文件"
        >
          <Download :size="15" />
          <span class="hidden sm:inline">下载</span>
        </a>
      </div>
    </div>
  </div>
</template>

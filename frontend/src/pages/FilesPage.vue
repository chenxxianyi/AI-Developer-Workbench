<script setup lang="ts">
import { onMounted, ref, type Component } from 'vue'
import { useRoute } from 'vue-router'
import {
  AlertCircle,
  Download,
  FileCode,
  FileJson,
  FileText,
  Files,
  Folder,
  Loader2,
  X,
} from '@lucide/vue'
import apiClient from '@/api/client'
import ProjectStageShell from '@/components/project/ProjectStageShell.vue'

interface FileEntry {
  name: string
  is_dir: boolean
}

interface FileContent {
  path: string
  content: string
  size?: number
  is_binary?: boolean
}

const route = useRoute()
const projectId = route.params.projectId as string
const files = ref<FileEntry[]>([])
const selectedFile = ref<FileContent | null>(null)
const selectedName = ref('')
const loading = ref(true)
const loadingFile = ref(false)
const error = ref('')

function getFileIcon(file: FileEntry): Component {
  if (file.is_dir) return Folder
  const extension = file.name.split('.').pop()?.toLowerCase()
  if (extension === 'json') return FileJson
  if (['vue', 'ts', 'tsx', 'js', 'jsx', 'css', 'html', 'go', 'py', 'java'].includes(extension || '')) {
    return FileCode
  }
  return FileText
}

function formatSize(size?: number) {
  if (size === undefined) return ''
  if (size < 1024) return `${size} B`
  return `${(size / 1024).toFixed(1)} KB`
}

async function loadFiles() {
  loading.value = true
  error.value = ''
  try {
    files.value = await apiClient.get(`/projects/${projectId}/files`)
  } catch (err: any) {
    error.value = err.message || '项目文件加载失败'
  } finally {
    loading.value = false
  }
}

async function viewFile(file: FileEntry) {
  if (file.is_dir || loadingFile.value) return
  selectedName.value = file.name
  loadingFile.value = true
  error.value = ''
  try {
    selectedFile.value = await apiClient.get(
      `/projects/${projectId}/files/content?path=${encodeURIComponent(file.name)}`,
    )
  } catch (err: any) {
    error.value = err.message || '文件内容加载失败'
  } finally {
    loadingFile.value = false
  }
}

function closePreview() {
  selectedFile.value = null
  selectedName.value = ''
}

function exportZip() {
  const base = import.meta.env.VITE_API_BASE_URL || '/api/v1'
  window.open(`${base}/projects/${projectId}/export`, '_blank')
}

onMounted(loadFiles)
</script>

<template>
  <ProjectStageShell
    :icon="Files"
    title="项目文件"
    description="浏览生成后的目录与源码，并导出完整项目文件。"
    step-text="交付阶段"
  >
    <template #actions>
      <button
        type="button"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="exportZip"
      >
        <Download :size="16" />
        导出 ZIP
      </button>
    </template>

    <div v-if="error" role="alert" class="mb-5 flex items-center gap-2 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
      <AlertCircle :size="18" class="shrink-0" />
      {{ error }}
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-[340px_minmax(0,1fr)]">
      <section class="overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
        <div class="flex min-h-13 items-center justify-between border-b border-border px-4 py-3">
          <div class="flex items-center gap-2">
            <Folder :size="17" class="text-accent" />
            <h2 class="text-sm font-semibold text-text-primary">文件浏览器</h2>
          </div>
          <span class="text-xs text-text-muted">{{ files.length }} 项</span>
        </div>

        <div v-if="loading" class="space-y-2 p-3">
          <div v-for="index in 6" :key="index" class="h-10 animate-pulse rounded-lg bg-surface-muted" />
        </div>

        <div
          v-else-if="files.length"
          class="max-h-[640px] min-h-96 overflow-y-auto p-2"
          role="list"
        >
          <button
            v-for="file in files"
            :key="file.name"
            type="button"
            :disabled="file.is_dir"
            :class="[
              'flex min-h-11 w-full items-center gap-3 rounded-lg px-3 text-left transition-smooth focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
              selectedName === file.name
                ? 'bg-accent-soft text-accent'
                : 'text-text-secondary hover:bg-surface-muted hover:text-text-primary',
              file.is_dir ? 'cursor-default' : 'cursor-pointer',
            ]"
            role="listitem"
            @click="viewFile(file)"
          >
            <component
              :is="getFileIcon(file)"
              :size="18"
              :class="file.is_dir ? 'text-warning' : selectedName === file.name ? 'text-accent' : 'text-text-muted'"
            />
            <span class="min-w-0 flex-1 truncate text-sm font-medium">{{ file.name }}</span>
            <span v-if="file.is_dir" class="text-xs text-text-muted">目录</span>
          </button>
        </div>

        <div v-else class="flex min-h-96 flex-col items-center justify-center px-5 text-center">
          <div class="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-surface-muted text-text-muted">
            <Folder :size="23" />
          </div>
          <p class="text-sm font-medium text-text-primary">暂无生成文件</p>
          <p class="mt-1 text-xs text-text-muted">完成代码生成后，文件会显示在这里。</p>
        </div>
      </section>

      <section class="min-w-0 overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
        <div class="flex min-h-13 items-center justify-between gap-3 border-b border-border px-4 py-3">
          <div class="flex min-w-0 items-center gap-2">
            <FileCode :size="17" class="shrink-0 text-accent" />
            <h2 class="truncate text-sm font-semibold text-text-primary">
              {{ selectedFile?.path || '文件预览' }}
            </h2>
          </div>
          <div class="flex items-center gap-3">
            <span v-if="selectedFile?.size !== undefined" class="text-xs text-text-muted">
              {{ formatSize(selectedFile.size) }}
            </span>
            <button
              v-if="selectedFile"
              type="button"
              title="关闭文件预览"
              aria-label="关闭文件预览"
              class="flex h-8 w-8 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-surface-muted hover:text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
              @click="closePreview"
            >
              <X :size="17" />
            </button>
          </div>
        </div>

        <div v-if="loadingFile" class="flex min-h-[640px] items-center justify-center">
          <Loader2 :size="28" class="animate-spin text-accent" />
        </div>

        <div
          v-else-if="selectedFile?.is_binary"
          class="flex min-h-[640px] flex-col items-center justify-center px-6 text-center"
        >
          <FileText :size="30" class="mb-3 text-text-muted" />
          <p class="text-sm font-medium text-text-primary">该文件不支持文本预览</p>
          <p class="mt-1 text-xs text-text-muted">可导出项目后使用本地工具打开。</p>
        </div>

        <pre
          v-else-if="selectedFile"
          class="min-h-[640px] max-h-[640px] overflow-auto bg-[#171717] p-5 font-mono text-xs leading-6 text-[#d4d4d4]"
        ><code>{{ selectedFile.content }}</code></pre>

        <div v-else class="flex min-h-[640px] flex-col items-center justify-center px-6 text-center">
          <div class="mb-3 flex h-12 w-12 items-center justify-center rounded-lg bg-accent-soft text-accent">
            <FileCode :size="23" />
          </div>
          <p class="text-sm font-medium text-text-primary">选择文件查看内容</p>
          <p class="mt-1 text-xs text-text-muted">源码会在此区域以只读方式打开。</p>
        </div>
      </section>
    </div>
  </ProjectStageShell>
</template>

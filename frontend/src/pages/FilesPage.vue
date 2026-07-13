<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold">项目文件</h1>
      <button @click="exportZip" class="px-4 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm">📦 导出 ZIP</button>
    </div>
    <div class="bg-[var(--color-surface)] rounded-lg border p-4">
      <div v-if="files.length === 0" class="text-center py-12 text-[var(--color-text-muted)]">
        <p class="mb-2">项目暂无生成文件</p>
        <p class="text-xs">请先生成项目代码</p>
      </div>
      <div v-else class="space-y-1">
        <div v-for="f in files" :key="f.name" class="flex items-center gap-3 px-3 py-2 rounded hover:bg-[var(--color-surface-muted)] cursor-pointer" @click="viewFile(f)">
          <span class="text-sm">{{ f.is_dir ? '📁' : '📄' }}</span>
          <span class="flex-1 text-sm">{{ f.name }}</span>
          <span v-if="f.is_dir" class="text-xs text-[var(--color-text-muted)]">{{ f.children?.length || 0 }} 项</span>
        </div>
      </div>
    </div>
    <!-- File content preview -->
    <div v-if="selectedFile" class="bg-[var(--color-surface)] rounded-lg border p-4">
      <div class="flex items-center justify-between mb-3"><h3 class="font-semibold text-sm">{{ selectedFile.path }}</h3><button @click="selectedFile = null" class="text-xs">✕ 关闭</button></div>
      <pre class="bg-[#1e1e1e] text-[#d4d4d4] rounded p-4 text-xs overflow-auto max-h-96"><code>{{ selectedFile.content }}</code></pre>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import apiClient from '@/api/client'
const route = useRoute()
const projectId = route.params.projectId as string
const files = ref<any[]>([])
const selectedFile = ref<any>(null)
onMounted(async () => { try { files.value = await apiClient.get(`/projects/${projectId}/files`) as any } catch {} })
async function viewFile(f: any) {
  if (f.is_dir) return
  try { selectedFile.value = await apiClient.get(`/projects/${projectId}/files/content?path=${encodeURIComponent(f.name)}`) as any } catch {}
}
function exportZip() { window.open(`${import.meta.env.VITE_API_BASE_URL || '/api/v1'}/projects/${projectId}/export`, '_blank') }
</script>

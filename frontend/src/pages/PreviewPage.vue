<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold">在线预览</h1>
      <div class="flex gap-2">
        <button @click="rebuild" class="px-4 py-2 rounded-lg border text-sm">🔄 重新构建</button>
        <a :href="previewUrl" target="_blank" class="px-4 py-2 rounded-lg border text-sm">🔗 新窗口</a>
      </div>
    </div>
    <div v-if="previewUrl" class="border rounded-lg overflow-hidden bg-white">
      <iframe :src="previewUrl" class="w-full h-[70vh]" sandbox="allow-scripts allow-same-origin" />
    </div>
    <div v-else class="text-center py-16 bg-[var(--color-surface)] rounded-lg border">
      <p class="text-[var(--color-text-muted)] mb-4">尚未构建预览</p>
      <button @click="rebuild" class="px-6 py-3 rounded-lg bg-[var(--color-accent)] text-white">🏗 构建并预览</button>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import apiClient from '@/api/client'
const route = useRoute()
const projectId = route.params.projectId as string
const previewUrl = ref('')
async function rebuild() {
  try {
    const res = await apiClient.post(`/projects/${projectId}/build`) as any
    previewUrl.value = res.preview_url || ''
  } catch {}
}
</script>

<template>
  <div class="max-w-4xl mx-auto py-8">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold">蓝图评审</h1>
      <div class="flex gap-2">
        <button @click="generate" class="px-4 py-2 rounded-lg border text-sm">🔄 重新生成</button>
        <button @click="confirm" class="px-6 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm">✅ 确认蓝图</button>
      </div>
    </div>
    <div v-if="blueprint" class="space-y-4">
      <div class="bg-[var(--color-surface)] rounded-[var(--radius-panel)] border p-5"><h3 class="font-semibold mb-2">📋 产品定位</h3><p class="text-sm text-[var(--color-text-secondary)]">{{ blueprint.product_positioning || 'AI 生成的企业网站' }}</p></div>
      <div class="grid grid-cols-2 gap-4">
        <div class="bg-[var(--color-surface)] rounded-[var(--radius-panel)] border p-5"><h3 class="font-semibold mb-2">🛠 技术栈</h3><p class="text-sm">{{ blueprint.tech_stack || 'Vue 3 + Tailwind CSS' }}</p></div>
        <div class="bg-[var(--color-surface)] rounded-[var(--radius-panel)] border p-5"><h3 class="font-semibold mb-2">📄 页面</h3><ul class="text-sm space-y-1"><li v-for="p in blueprint.pages || []" :key="p.name">{{ p.name }} → {{ p.route }}</li></ul></div>
      </div>
    </div>
    <div v-else class="text-center py-16"><p class="text-[var(--color-text-muted)] mb-4">尚未生成蓝图</p><button @click="generate" class="px-6 py-3 rounded-lg bg-[var(--color-accent)] text-white">🚀 生成蓝图</button></div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import apiClient from '@/api/client'
const route = useRoute(); const router = useRouter()
const projectId = route.params.projectId as string
const blueprint = ref<any>(null)
onMounted(async () => { try { blueprint.value = await apiClient.get(`/projects/${projectId}/blueprint`) } catch {} })
async function generate() { try { blueprint.value = await apiClient.post(`/projects/${projectId}/blueprint/generate`) } catch {} }
async function confirm() { await apiClient.post(`/projects/${projectId}/blueprint/confirm`); router.push(`/projects/${projectId}/generation`) }
</script>

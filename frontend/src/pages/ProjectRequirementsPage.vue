<template>
  <div class="max-w-3xl mx-auto py-8">
    <h1 class="text-2xl font-bold text-[var(--color-text-primary)] mb-6">项目需求</h1>
    <div class="bg-[var(--color-surface)] rounded-[var(--radius-panel)] border border-[var(--color-border)] p-6 space-y-4">
      <div><label for="project-goals" class="block text-sm font-medium mb-1">项目目标</label><textarea id="project-goals" v-model="goals" rows="3" class="w-full px-3 py-2 border rounded-lg bg-white" placeholder="描述项目要达成的目标..." /></div>
      <div><label for="project-audience" class="block text-sm font-medium mb-1">目标用户</label><input id="project-audience" v-model="audience" class="w-full px-3 py-2 border rounded-lg bg-white" placeholder="例如：中小企业管理者" /></div>
      <div><label for="project-features" class="block text-sm font-medium mb-1">核心功能（每行一个）</label><textarea id="project-features" v-model="features" rows="4" class="w-full px-3 py-2 border rounded-lg bg-white" placeholder="用户注册登录&#10;产品展示&#10;在线下单" /></div>
      <p v-if="error" class="text-sm text-[var(--color-danger)]">{{ error }}</p>
      <div class="flex justify-end gap-3 pt-4">
        <button @click="saveDraft" :disabled="saving" class="px-4 py-2 rounded-lg border text-sm disabled:opacity-50">{{ saving ? '保存中...' : '保存草稿' }}</button>
        <button @click="saveAndGenerate" :disabled="saving" class="px-6 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm disabled:opacity-50">{{ saving ? '保存中...' : '保存并生成蓝图' }}</button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import apiClient from '@/api/client'

const route = useRoute(); const router = useRouter()
const projectId = route.params.projectId as string
const goals = ref(''); const audience = ref(''); const features = ref('')
const saving = ref(false); const error = ref('')

function buildContent() {
  return JSON.stringify({
    goals: goals.value,
    audience: audience.value,
    features: features.value.split('\n').filter(Boolean),
  })
}

async function loadRequirements() {
  try {
    const res: any = await apiClient.get(`/projects/${projectId}/requirements`)
    if (res?.content) {
      const data = typeof res.content === 'string' ? JSON.parse(res.content) : res.content
      goals.value = data.goals || ''
      audience.value = data.audience || ''
      features.value = (data.features || []).join('\n')
    }
  } catch { /* no saved requirements yet */ }
}

async function saveRequirements() {
  saving.value = true; error.value = ''
  try {
    await apiClient.put(`/projects/${projectId}/requirements`, { content: buildContent() })
  } catch (err: any) {
    error.value = err.message || '保存失败'
  } finally { saving.value = false }
}

async function saveDraft() {
  await saveRequirements()
}

async function saveAndGenerate() {
  await saveRequirements()
  if (!error.value) router.push(`/projects/${projectId}/blueprint`)
}

onMounted(loadRequirements)
</script>


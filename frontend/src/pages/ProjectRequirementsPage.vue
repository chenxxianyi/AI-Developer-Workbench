<template>
  <div class="max-w-3xl mx-auto py-8">
    <h1 class="text-2xl font-bold text-[var(--color-text-primary)] mb-6">项目需求</h1>
    <div class="bg-[var(--color-surface)] rounded-[var(--radius-panel)] border border-[var(--color-border)] p-6 space-y-4">
      <div><label class="block text-sm font-medium mb-1">项目目标</label><textarea v-model="goals" rows="3" class="w-full px-3 py-2 border rounded-lg bg-white" placeholder="描述项目要达成的目标..." /></div>
      <div><label class="block text-sm font-medium mb-1">目标用户</label><input v-model="audience" class="w-full px-3 py-2 border rounded-lg bg-white" placeholder="例如：中小企业管理者" /></div>
      <div><label class="block text-sm font-medium mb-1">核心功能（每行一个）</label><textarea v-model="features" rows="4" class="w-full px-3 py-2 border rounded-lg bg-white" placeholder="用户注册登录&#10;产品展示&#10;在线下单" /></div>
      <div class="flex justify-end gap-3 pt-4">
        <button @click="saveDraft" class="px-4 py-2 rounded-lg border text-sm">保存草稿</button>
        <button @click="saveAndGenerate" class="px-6 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm">保存并生成蓝图</button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
const route = useRoute(); const router = useRouter()
const projectId = route.params.projectId as string
const goals = ref(''); const audience = ref(''); const features = ref('')
async function saveDraft() { /* API call */ }
async function saveAndGenerate() { router.push(`/projects/${projectId}/blueprint`) }
</script>

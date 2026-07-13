<template>
  <div class="max-w-3xl mx-auto py-8">
    <h1 class="text-2xl font-bold text-[var(--color-text-primary)] mb-2">创建新项目</h1>
    <p class="text-[var(--color-text-muted)] mb-8">选择项目类型，开始您的 AI 驱动开发之旅</p>

    <!-- Step 1: 选择类型 -->
    <div v-if="step === 1" class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <button v-for="opt in projectTypes" :key="opt.value"
        @click="selectedType = opt.value; step = 2"
        class="p-6 rounded-[var(--radius-card)] border-2 text-left transition-all duration-200"
        :class="selectedType === opt.value ? 'border-[var(--color-accent)] bg-[var(--color-accent-soft)]' : 'border-[var(--color-border)] hover:border-[var(--color-accent)]/50'"
      >
        <div class="text-2xl mb-2">{{ opt.icon }}</div>
        <h3 class="text-lg font-semibold text-[var(--color-text-primary)]">{{ opt.title }}</h3>
        <p class="text-sm text-[var(--color-text-muted)] mt-1">{{ opt.desc }}</p>
      </button>
    </div>

    <!-- Step 2: 填写信息 -->
    <div v-if="step === 2" class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-[var(--color-text-primary)] mb-1">项目名称 *</label>
        <input v-model="form.name" class="w-full px-3 py-2 border border-[var(--color-border)] rounded-[var(--radius-md)] bg-[var(--color-surface)] text-[var(--color-text-primary)] focus:ring-2 focus:ring-[var(--color-accent)]/30 focus:border-[var(--color-accent)]" placeholder="例如：我的企业官网" />
      </div>
      <div>
        <label class="block text-sm font-medium text-[var(--color-text-primary)] mb-1">项目描述</label>
        <textarea v-model="form.description" rows="3" class="w-full px-3 py-2 border border-[var(--color-border)] rounded-[var(--radius-md)] bg-[var(--color-surface)]" placeholder="简要描述项目目标和用途..." />
      </div>
      <div v-if="selectedType === 'website'">
        <label class="block text-sm font-medium text-[var(--color-text-primary)] mb-2">技术栈偏好</label>
        <div class="grid grid-cols-2 gap-2">
          <label v-for="stack in techStacks" :key="stack" class="flex items-center gap-2 text-sm">
            <input type="checkbox" v-model="form.techStacks" :value="stack" class="rounded" />
            {{ stack }}
          </label>
        </div>
      </div>
      <div class="flex gap-3 pt-4">
        <button @click="step = 1" class="px-4 py-2 rounded-lg border border-[var(--color-border)] text-sm">返回</button>
        <button @click="createProject" :disabled="!form.name" class="px-6 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm font-medium disabled:opacity-50">创建项目</button>
      </div>
    </div>

    <!-- Result -->
    <div v-if="step === 3" class="text-center py-12">
      <div class="text-4xl mb-3">🎉</div>
      <h2 class="text-xl font-bold text-[var(--color-text-primary)] mb-2">项目创建成功！</h2>
      <p class="text-[var(--color-text-muted)] mb-4">项目 <strong>{{ createdProject?.name }}</strong> 已创建</p>
      <div class="flex gap-3 justify-center">
        <RouterLink :to="`/projects/${createdProject?.id}`" class="px-4 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm">进入项目</RouterLink>
        <RouterLink :to="`/projects/${createdProject?.id}/requirements`" class="px-4 py-2 rounded-lg border border-[var(--color-border)] text-sm">填写需求</RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { RouterLink } from 'vue-router'
import apiClient from '@/api/client'

const step = ref(1)
const selectedType = ref('website')
const form = reactive({ name: '', description: '', techStacks: [] as string[] })
const createdProject = ref<any>(null)

const projectTypes = [
  { value: 'website', icon: '🌐', title: '生成新网站', desc: 'AI 驱动从需求到代码的完整网站生成' },
  { value: 'analysis', icon: '🔍', title: '分析已有项目', desc: '上传或导入现有项目进行质量分析' },
]

const techStacks = ['Vue 3', 'React', 'Tailwind CSS', 'TypeScript', 'Node.js', 'Go']

async function createProject() {
  try {
    const res = await apiClient.post('/projects', { name: form.name, description: form.description })
    createdProject.value = res
    step.value = 3
  } catch { /* error handled by interceptor */ }
}
</script>

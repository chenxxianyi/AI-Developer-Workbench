<script setup lang="ts">
/**
 * Agent Config Studio Page
 * Generate AI agent configuration files
 */

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { runAgentConfig } from '@/api/tools'
import type { AgentConfigInput } from '@/types/tool'
import type { Report, AgentConfigResult } from '@/types/report'
import {
  Bot,
  Loader2,
  AlertCircle,
  CheckCircle2,
  FileText,
  ArrowLeft,
  Download,
} from '@lucide/vue'

const router = useRouter()

// Form state
const title = ref('')
const projectName = ref('')
const projectType = ref('')
const frontendStack = ref('')
const backendStack = ref('')
const database = ref('')
const uiStyle = ref('')
const codingPreferences = ref('')
const strictRules = ref('')

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<AgentConfigResult> | null>(null)

// Submit
async function handleSubmit() {
  if (!title.value.trim()) {
    error.value = '请输入标题'
    return
  }
  if (!projectName.value.trim()) {
    error.value = '请输入项目名称'
    return
  }

  loading.value = true
  error.value = null
  result.value = null

  try {
    const payload: AgentConfigInput = {
      title: title.value,
      project_name: projectName.value,
      project_type: projectType.value || undefined,
      frontend_stack: frontendStack.value || undefined,
      backend_stack: backendStack.value || undefined,
      database: database.value || undefined,
      ui_style: uiStyle.value || undefined,
      coding_preferences: codingPreferences.value || undefined,
      strict_rules: strictRules.value || undefined,
    }

    result.value = await runAgentConfig(payload) as unknown as Report<AgentConfigResult>
  } catch (err: any) {
    error.value = err.message || '生成失败'
  } finally {
    loading.value = false
  }
}

function downloadFile(filename: string) {
  if (!result.value) return
  window.open(`/api/reports/${result.value.id}/files/${filename}`, '_blank')
}

function exportMarkdown() {
  if (!result.value) return
  window.open(`/api/reports/${result.value.id}/export?format=markdown`, '_blank')
}

function resetForm() {
  title.value = ''
  projectName.value = ''
  projectType.value = ''
  frontendStack.value = ''
  backendStack.value = ''
  database.value = ''
  uiStyle.value = ''
  codingPreferences.value = ''
  strictRules.value = ''
  result.value = null
  error.value = null
}

const canSubmit = computed(() => title.value.trim() && projectName.value.trim())
</script>

<template>
  <div class="max-w-6xl mx-auto">
    <!-- Header -->
    <div class="mb-6">
      <button @click="router.push('/dashboard')" class="flex items-center gap-2 text-text-secondary hover:text-text-primary transition-smooth mb-4">
        <ArrowLeft :size="20" />
        <span>返回工作台</span>
      </button>
      <div class="flex items-center gap-3">
        <div class="w-12 h-12 bg-warning rounded-xl flex items-center justify-center">
          <Bot :size="24" class="text-white" />
        </div>
        <div>
          <h1 class="text-2xl font-bold text-text-primary">Agent 配置生成</h1>
          <p class="text-text-secondary">为项目生成 AI Agent 配置文件</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Input Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">项目配置</h2>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">标题 *</label>
          <input v-model="title" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="输入配置标题..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">项目名称 *</label>
          <input v-model="projectName" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: AI Workbench..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">项目类型</label>
          <input v-model="projectType" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: fullstack, frontend, backend..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">前端技术栈</label>
          <input v-model="frontendStack" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: Vue 3 + TypeScript + Vite..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">后端技术栈</label>
          <input v-model="backendStack" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: Go + Gin + GORM..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">数据库</label>
          <input v-model="database" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: MySQL, PostgreSQL..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">UI 风格</label>
          <input v-model="uiStyle" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: Tailwind CSS, 简洁现代..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">编码偏好</label>
          <textarea v-model="codingPreferences" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" rows="3" placeholder="描述编码风格和偏好..."></textarea>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">严格规则</label>
          <textarea v-model="strictRules" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" rows="3" placeholder="描述必须遵守的规则..."></textarea>
        </div>

        <!-- Error -->
        <div v-if="error" class="mb-4 p-3 bg-danger/10 border border-danger/20 rounded-lg">
          <div class="flex items-center gap-2">
            <AlertCircle :size="18" class="text-danger" />
            <span class="text-danger">{{ error }}</span>
          </div>
        </div>

        <!-- Submit -->
        <div class="flex gap-3">
          <button @click="handleSubmit" :disabled="loading || !canSubmit" :class="['flex items-center gap-2 px-6 py-2 rounded-lg transition-smooth', loading || !canSubmit ? 'bg-surface-muted text-text-muted cursor-not-allowed' : 'bg-warning text-white hover:bg-warning/80']">
            <Loader2 v-if="loading" :size="18" class="animate-spin" />
            <Bot v-else :size="18" />
            <span>{{ loading ? '生成中...' : '生成配置' }}</span>
          </button>
          <button @click="resetForm" class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth">重置</button>
        </div>
      </div>

      <!-- Result Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">生成结果</h2>

        <div v-if="!result && !loading" class="text-center py-12">
          <Bot :size="48" class="text-text-muted mx-auto mb-4" />
          <p class="text-text-secondary">填写项目信息后生成配置文件</p>
        </div>

        <div v-if="loading" class="text-center py-12">
          <Loader2 :size="48" class="text-accent mx-auto mb-4 animate-spin" />
          <p class="text-text-secondary">AI 正在生成配置...</p>
        </div>

        <div v-if="result && !loading">
          <!-- Summary -->
          <div class="mb-6 p-4 bg-surface-muted rounded-lg">
            <p class="text-text-primary">{{ result.summary }}</p>
          </div>

          <!-- Generated Files -->
          <div v-if="result.report_data?.generated_files_content" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">生成的配置文件</h3>
            <div class="space-y-3">
              <div v-for="(content, filename) in result.report_data.generated_files_content" :key="filename" class="p-3 bg-surface-muted rounded">
                <div class="flex items-center justify-between mb-2">
                  <span class="font-medium text-accent">{{ filename }}</span>
                  <button @click="downloadFile(filename)" class="flex items-center gap-1 px-2 py-1 bg-accent-soft text-accent rounded hover:bg-accent hover:text-white transition-smooth text-sm">
                    <Download :size="14" />
                    <span>下载</span>
                  </button>
                </div>
                <pre class="text-text-secondary text-sm overflow-x-auto whitespace-pre-wrap">{{ content.substring(0, 200) }}...</pre>
              </div>
            </div>
          </div>

          <!-- Recommendations -->
          <div v-if="result.report_data?.recommendations?.length" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">建议</h3>
            <ul class="space-y-1">
              <li v-for="(rec, idx) in result.report_data.recommendations" :key="idx" class="flex items-start gap-2 text-text-secondary">
                <CheckCircle2 :size="14" class="text-success mt-1" />
                <span class="text-sm">{{ rec }}</span>
              </li>
            </ul>
          </div>

          <!-- Download All -->
          <div v-if="result.generated_files?.length" class="mb-4">
            <h3 class="text-md font-semibold text-text-primary mb-2">下载全部文件</h3>
            <div class="flex gap-2">
              <button v-for="file in result.generated_files" :key="file.id" @click="downloadFile(file.filename)" class="flex items-center gap-1 px-3 py-1 bg-accent-soft text-accent rounded hover:bg-accent hover:text-white transition-smooth text-sm">
                <FileText :size="14" />
                <span>{{ file.filename }}</span>
              </button>
            </div>
          </div>

          <button @click="exportMarkdown" class="flex items-center gap-2 px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth">
            <Download :size="18" />
            <span>导出 Markdown</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

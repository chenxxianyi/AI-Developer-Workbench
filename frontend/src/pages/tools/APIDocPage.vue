<script setup lang="ts">
/**
 * API Doc Builder Page
 * Generate API documentation from code or project ZIP
 */

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { runAPIDoc } from '@/api/tools'
import type { Report, APIDocResult } from '@/types/report'
import {
  FileText,
  Upload,
  Loader2,
  AlertCircle,
  CheckCircle2,
  ArrowLeft,
  Download,
} from '@lucide/vue'

const router = useRouter()

// Form state
const title = ref('')
const sourceType = ref<'code' | 'project_zip' | 'manual'>('code')
const backendStack = ref('')
const codeInput = ref('')
const apiDescription = ref('')
const outputFormat = ref<'markdown' | 'openapi' | 'both'>('markdown')
const projectFile = ref<File | null>(null)
const fileName = ref('')

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<APIDocResult> | null>(null)

// File handling
function handleFileSelect(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files?.length) {
    projectFile.value = target.files[0]
    fileName.value = target.files[0].name
  }
}

function clearFile() {
  projectFile.value = null
  fileName.value = ''
}

// Submit
async function handleSubmit() {
  if (!title.value.trim()) {
    error.value = '请输入标题'
    return
  }

  if (sourceType.value === 'code' && !codeInput.value.trim()) {
    error.value = '请输入代码'
    return
  }

  if (sourceType.value === 'project_zip' && !projectFile.value) {
    error.value = '请上传项目 ZIP'
    return
  }

  if (sourceType.value === 'manual' && !apiDescription.value.trim()) {
    error.value = '请输入 API 描述'
    return
  }

  loading.value = true
  error.value = null
  result.value = null

  try {
    const formData = new FormData()
    formData.append('title', title.value)
    formData.append('source_type', sourceType.value)
    formData.append('output_format', outputFormat.value)
    if (backendStack.value) formData.append('backend_stack', backendStack.value)
    if (codeInput.value) formData.append('code', codeInput.value)
    if (apiDescription.value) formData.append('api_description', apiDescription.value)
    if (projectFile.value) formData.append('project_zip', projectFile.value)

    result.value = await runAPIDoc(formData) as Report<APIDocResult>
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
  sourceType.value = 'code'
  backendStack.value = ''
  codeInput.value = ''
  apiDescription.value = ''
  outputFormat.value = 'markdown'
  projectFile.value = null
  fileName.value = ''
  result.value = null
  error.value = null
}

const canSubmit = computed(() => {
  if (!title.value.trim()) return false
  if (sourceType.value === 'code') return !!codeInput.value.trim()
  if (sourceType.value === 'project_zip') return !!projectFile.value
  if (sourceType.value === 'manual') return !!apiDescription.value.trim()
  return false
})
</script>

<template>
  <div class="max-w-6xl mx-auto">
    <!-- Header -->
    <div class="mb-6">
      <button @click="router.push('/dashboard')" class="flex items-center gap-2 text-text-secondary hover:text-text-primary transition-smooth mb-4">
        <ArrowLeft :size="20" />
        <span>返回 Dashboard</span>
      </button>
      <div class="flex items-center gap-3">
        <div class="w-12 h-12 bg-danger rounded-xl flex items-center justify-center">
          <FileText :size="24" class="text-white" />
        </div>
        <div>
          <h1 class="text-2xl font-bold text-text-primary">API Doc Builder</h1>
          <p class="text-text-secondary">Auto-generate API documentation from code or descriptions</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Input Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">输入参数</h2>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">标题 *</label>
          <input v-model="title" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="输入文档标题..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">来源类型 *</label>
          <div class="flex gap-2">
            <button @click="sourceType = 'code'" :class="['flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth', sourceType === 'code' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border']">
              <span>代码</span>
            </button>
            <button @click="sourceType = 'project_zip'" :class="['flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth', sourceType === 'project_zip' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border']">
              <span>项目 ZIP</span>
            </button>
            <button @click="sourceType = 'manual'" :class="['flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth', sourceType === 'manual' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border']">
              <span>手动描述</span>
            </button>
          </div>
        </div>

        <!-- Code Input -->
        <div v-if="sourceType === 'code'" class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">API 代码 *</label>
          <textarea v-model="codeInput" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none font-mono text-sm" rows="10" placeholder="粘贴 API handler/controller 代码..."></textarea>
        </div>

        <!-- ZIP Upload -->
        <div v-if="sourceType === 'project_zip'" class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">项目 ZIP *</label>
          <div v-if="!fileName" class="border-2 border-dashed border-border rounded-lg p-8 text-center hover:border-accent transition-smooth cursor-pointer" @click="($refs.fileInput as HTMLInputElement).click()">
            <Upload :size="32" class="text-text-muted mx-auto mb-2" />
            <p class="text-text-secondary">点击上传 ZIP</p>
            <input ref="fileInput" type="file" accept=".zip" class="hidden" @change="handleFileSelect" />
          </div>
          <div v-else class="flex items-center justify-between p-3 bg-surface-muted rounded-lg">
            <div class="flex items-center gap-2">
              <FileText :size="18" class="text-accent" />
              <span class="text-text-primary">{{ fileName }}</span>
            </div>
            <button @click="clearFile" class="p-1 text-danger"><AlertCircle :size="16" /></button>
          </div>
        </div>

        <!-- Manual Description -->
        <div v-if="sourceType === 'manual'" class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">API 描述 *</label>
          <textarea v-model="apiDescription" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" rows="6" placeholder="描述 API 功能、端点、参数等..."></textarea>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">后端技术栈</label>
          <input v-model="backendStack" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: Go Gin, Express, Django..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">输出格式 *</label>
          <div class="flex gap-2">
            <button @click="outputFormat = 'markdown'" :class="['px-4 py-2 rounded-lg transition-smooth', outputFormat === 'markdown' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary']">Markdown</button>
            <button @click="outputFormat = 'openapi'" :class="['px-4 py-2 rounded-lg transition-smooth', outputFormat === 'openapi' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary']">OpenAPI</button>
            <button @click="outputFormat = 'both'" :class="['px-4 py-2 rounded-lg transition-smooth', outputFormat === 'both' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary']">两者</button>
          </div>
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
          <button @click="handleSubmit" :disabled="loading || !canSubmit" :class="['flex items-center gap-2 px-6 py-2 rounded-lg transition-smooth', loading || !canSubmit ? 'bg-surface-muted text-text-muted cursor-not-allowed' : 'bg-accent text-white hover:bg-accent/80']">
            <Loader2 v-if="loading" :size="18" class="animate-spin" />
            <FileText v-else :size="18" />
            <span>{{ loading ? '生成中...' : '生成文档' }}</span>
          </button>
          <button @click="resetForm" class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth">重置</button>
        </div>
      </div>

      <!-- Result Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">生成结果</h2>

        <div v-if="!result && !loading" class="text-center py-12">
          <FileText :size="48" class="text-text-muted mx-auto mb-4" />
          <p class="text-text-secondary">提交输入后生成文档</p>
        </div>

        <div v-if="loading" class="text-center py-12">
          <Loader2 :size="48" class="text-accent mx-auto mb-4 animate-spin" />
          <p class="text-text-secondary">AI 正在生成文档...</p>
        </div>

        <div v-if="result && !loading">
          <div class="mb-6 p-4 bg-surface-muted rounded-lg">
            <p class="text-text-primary">{{ result.summary }}</p>
          </div>

          <!-- Modules -->
          <div v-if="result.report_data?.modules?.length" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">API 模块</h3>
            <div class="space-y-2">
              <div v-for="module in result.report_data.modules" :key="module.name" class="p-3 bg-surface-muted rounded">
                <span class="font-medium text-text-primary">{{ module.name }}</span>
                <span class="text-text-muted text-sm ml-2">{{ module.endpoints?.length || 0 }} endpoints</span>
              </div>
            </div>
          </div>

          <!-- Markdown Content -->
          <div v-if="result.report_data?.markdown_content" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">Markdown 文档</h3>
            <pre class="p-3 bg-surface-muted rounded text-sm text-text-secondary overflow-x-auto whitespace-pre-wrap">{{ result.report_data.markdown_content.substring(0, 500) }}...</pre>
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

          <!-- Generated Files -->
          <div v-if="result.generated_files?.length" class="mb-4">
            <h3 class="text-md font-semibold text-text-primary mb-2">生成的文件</h3>
            <div class="flex gap-2">
              <button v-for="file in result.generated_files" :key="file.id" @click="downloadFile(file.filename)" class="flex items-center gap-1 px-3 py-1 bg-accent-soft text-accent rounded hover:bg-accent hover:text-white transition-smooth text-sm">
                <Download :size="14" />
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
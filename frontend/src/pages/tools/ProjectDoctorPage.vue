<script setup lang="ts">
/**
 * Project Doctor Page
 * Upload project ZIP for comprehensive health analysis
 */

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { runProjectDoctor } from '@/api/tools'
import type { Report, ProjectDoctorResult } from '@/types/report'
import {
  Stethoscope,
  Upload,
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
const techStack = ref('')
const projectDescription = ref('')
const analysisDepth = ref<'basic' | 'standard' | 'deep'>('standard')
const projectFile = ref<File | null>(null)
const fileName = ref('')

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<ProjectDoctorResult> | null>(null)

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
  if (!projectFile.value) {
    error.value = '请上传项目 ZIP 文件'
    return
  }

  loading.value = true
  error.value = null
  result.value = null

  try {
    const formData = new FormData()
    formData.append('title', title.value)
    formData.append('project_name', projectName.value || '未命名项目')
    formData.append('analysis_depth', analysisDepth.value)
    if (techStack.value) formData.append('tech_stack', techStack.value)
    if (projectDescription.value) formData.append('project_description', projectDescription.value)
    formData.append('project_zip', projectFile.value)

    result.value = await runProjectDoctor(formData) as unknown as Report<ProjectDoctorResult>
  } catch (err: any) {
    error.value = err.message || '诊断失败'
  } finally {
    loading.value = false
  }
}

// Helpers
function getScoreColor(score: number, max: number): string {
  const percent = (score / max) * 100
  if (percent >= 80) return 'text-success'
  if (percent >= 60) return 'text-warning'
  return 'text-danger'
}

function getSeverityColor(severity: string): string {
  if (severity === 'high') return 'text-danger'
  if (severity === 'medium') return 'text-warning'
  return 'text-accent'
}

function getGradeColor(grade: string | null): string {
  if (!grade) return 'text-text-muted'
  if (grade === 'A') return 'text-success'
  if (grade === 'B') return 'text-accent'
  if (grade === 'C') return 'text-warning'
  return 'text-danger'
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
  techStack.value = ''
  projectDescription.value = ''
  analysisDepth.value = 'standard'
  projectFile.value = null
  fileName.value = ''
  result.value = null
  error.value = null
}

const canSubmit = computed(() => title.value.trim() && projectFile.value)
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
        <div class="w-12 h-12 bg-success rounded-xl flex items-center justify-center">
          <Stethoscope :size="24" class="text-white" />
        </div>
        <div>
          <h1 class="text-2xl font-bold text-text-primary">项目诊断</h1>
          <p class="text-text-secondary">全面检查项目健康度和工程质量</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Input Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">输入参数</h2>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">标题 *</label>
          <input v-model="title" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="输入诊断标题..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">项目名称</label>
          <input v-model="projectName" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: AI Workbench..." />
        </div>

        <!-- File Upload -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">项目 ZIP 文件 *</label>
          <div v-if="!fileName" class="border-2 border-dashed border-border rounded-lg p-8 text-center hover:border-accent transition-smooth cursor-pointer" @click="($refs.fileInput as HTMLInputElement).click()">
            <Upload :size="32" class="text-text-muted mx-auto mb-2" />
            <p class="text-text-secondary">点击或拖拽上传 ZIP</p>
            <p class="text-text-muted text-sm mt-1">最大 20MB</p>
            <input ref="fileInput" type="file" accept=".zip" class="hidden" @change="handleFileSelect" />
          </div>
          <div v-else class="flex items-center justify-between p-3 bg-surface-muted rounded-lg">
            <div class="flex items-center gap-2">
              <FileText :size="18" class="text-accent" />
              <span class="text-text-primary">{{ fileName }}</span>
            </div>
            <button @click="clearFile" class="p-1 text-danger hover:text-danger/80">
              <AlertCircle :size="16" />
            </button>
          </div>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">技术栈</label>
          <input v-model="techStack" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: Vue 3 + Go + MySQL..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">分析深度 *</label>
          <div class="flex gap-2">
            <button @click="analysisDepth = 'basic'" :class="['flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth', analysisDepth === 'basic' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border']">
              <span>基础</span>
            </button>
            <button @click="analysisDepth = 'standard'" :class="['flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth', analysisDepth === 'standard' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border']">
              <span>标准</span>
            </button>
            <button @click="analysisDepth = 'deep'" :class="['flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth', analysisDepth === 'deep' ? 'bg-accent text-white' : 'bg-surface-muted text-text-secondary hover:bg-border']">
              <span>深度</span>
            </button>
          </div>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">项目描述</label>
          <textarea v-model="projectDescription" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" rows="3" placeholder="描述项目功能和目标..."></textarea>
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
          <button @click="handleSubmit" :disabled="loading || !canSubmit" :class="['flex items-center gap-2 px-6 py-2 rounded-lg transition-smooth', loading || !canSubmit ? 'bg-surface-muted text-text-muted cursor-not-allowed' : 'bg-success text-white hover:bg-success/80']">
            <Loader2 v-if="loading" :size="18" class="animate-spin" />
            <Stethoscope v-else :size="18" />
            <span>{{ loading ? '诊断中...' : '开始诊断' }}</span>
          </button>
          <button @click="resetForm" class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth">重置</button>
        </div>
      </div>

      <!-- Result Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">诊断结果</h2>

        <div v-if="!result && !loading" class="text-center py-12">
          <Stethoscope :size="48" class="text-text-muted mx-auto mb-4" />
          <p class="text-text-secondary">上传项目 ZIP 开始诊断</p>
        </div>

        <div v-if="loading" class="text-center py-12">
          <Loader2 :size="48" class="text-accent mx-auto mb-4 animate-spin" />
          <p class="text-text-secondary">AI 正在分析项目...</p>
        </div>

        <div v-if="result && !loading">
          <!-- Summary -->
          <div class="mb-6 p-4 bg-surface-muted rounded-lg">
            <div class="flex items-center justify-between mb-2">
              <span class="text-text-secondary">健康评分</span>
              <span :class="getGradeColor(result.grade)">
                <span v-if="result.total_score">{{ result.total_score }}/100</span>
                <span v-if="result.grade" class="ml-2 font-bold text-xl">{{ result.grade }}</span>
              </span>
            </div>
            <p class="text-text-primary">{{ result.summary }}</p>
          </div>

          <!-- Scores -->
          <div v-if="result.report_data?.scores?.length" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">评分详情</h3>
            <div class="space-y-2">
              <div v-for="score in result.report_data.scores" :key="score.name" class="flex items-center justify-between p-2 bg-surface-muted rounded">
                <span class="text-text-secondary">{{ score.name }}</span>
                <span :class="getScoreColor(score.score, score.max_score)">{{ score.score }}/{{ score.max_score }}</span>
              </div>
            </div>
          </div>

          <!-- Issues -->
          <div v-if="result.report_data?.issues?.length" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">发现的问题</h3>
            <div class="space-y-2">
              <div v-for="(issue, idx) in result.report_data.issues" :key="idx" class="p-3 bg-surface-muted rounded">
                <div class="flex items-center justify-between mb-1">
                  <span class="font-medium text-text-primary">{{ issue.title }}</span>
                  <span :class="getSeverityColor(issue.severity)" class="text-sm">{{ issue.severity }}</span>
                </div>
                <p class="text-text-secondary text-sm">{{ issue.problem }}</p>
                <p class="text-accent text-sm mt-1">建议: {{ issue.suggestion }}</p>
              </div>
            </div>
          </div>

          <!-- Recommendations -->
          <div v-if="result.report_data?.recommendations?.length" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">改进建议</h3>
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

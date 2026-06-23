<script setup lang="ts">
/**
 * UI Review Page
 * Upload screenshot or code for UI/UX analysis
 */

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { runUIReview } from '@/api/tools'
import type { Report, UIReviewResult } from '@/types/report'
import { getScoreDisplayName, getSeverityDisplayName } from '@/utils/uiReviewDisplay'
import {
  Eye,
  Upload,
  Code,
  Image,
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
const reviewMode = ref<'screenshot' | 'code' | 'screenshot_code'>('screenshot')
const pageType = ref('')
const targetStyle = ref('')
const description = ref('')
const codeInput = ref('')
const screenshotFile = ref<File | null>(null)
const screenshotPreview = ref<string | null>(null)

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<UIReviewResult> | null>(null)

// File handling
function setScreenshotFile(file: File) {
  screenshotFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => {
    screenshotPreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

function handleFileSelect(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files?.length) {
    setScreenshotFile(target.files[0])
  }
}

function handleScreenshotPaste(event: ClipboardEvent) {
  const items = event.clipboardData?.items
  if (!items?.length) return

  for (let index = 0; index < items.length; index += 1) {
    const item = items[index]
    if (!item.type.startsWith('image/')) continue

    const file = item.getAsFile()
    if (!file) return

    event.preventDefault()
    setScreenshotFile(file)
    return
  }
}

function clearFile() {
  screenshotFile.value = null
  screenshotPreview.value = null
}

// Submit
async function handleSubmit() {
  if (!title.value.trim()) {
    error.value = '请输入标题'
    return
  }

  if (reviewMode.value === 'screenshot' && !screenshotFile.value) {
    error.value = '请上传截图'
    return
  }

  if (reviewMode.value === 'code' && !codeInput.value.trim()) {
    error.value = '请输入代码'
    return
  }

  if (reviewMode.value === 'screenshot_code' && !screenshotFile.value && !codeInput.value.trim()) {
    error.value = '请上传截图或输入代码'
    return
  }

  loading.value = true
  error.value = null
  result.value = null

  try {
    const formData = new FormData()
    formData.append('title', title.value)
    formData.append('review_mode', reviewMode.value)
    if (pageType.value) formData.append('page_type', pageType.value)
    if (targetStyle.value) formData.append('target_style', targetStyle.value)
    if (description.value) formData.append('description', description.value)
    if (codeInput.value) formData.append('code', codeInput.value)
    if (screenshotFile.value) formData.append('screenshot', screenshotFile.value)

    result.value = await runUIReview(formData) as unknown as Report<UIReviewResult>
  } catch (err: any) {
    error.value = err.message || '分析失败'
  } finally {
    loading.value = false
  }
}

// Score color helper
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

// Download file
function downloadFile(filename: string) {
  if (!result.value) return
  window.open(`/api/reports/${result.value.id}/files/${filename}`, '_blank')
}

// Export markdown
function exportMarkdown() {
  if (!result.value) return
  window.open(`/api/reports/${result.value.id}/export?format=markdown`, '_blank')
}

// Reset form
function resetForm() {
  title.value = ''
  reviewMode.value = 'screenshot'
  pageType.value = ''
  targetStyle.value = ''
  description.value = ''
  codeInput.value = ''
  screenshotFile.value = null
  screenshotPreview.value = null
  result.value = null
  error.value = null
}

const canSubmit = computed(() => {
  if (!title.value.trim()) return false
  if (reviewMode.value === 'screenshot') return !!screenshotFile.value
  if (reviewMode.value === 'code') return !!codeInput.value.trim()
  return !!screenshotFile.value || !!codeInput.value.trim()
})
</script>

<template>
  <div class="max-w-6xl mx-auto">
    <!-- Header -->
    <div class="mb-6">
      <button
        @click="router.push('/dashboard')"
        class="flex items-center gap-2 text-text-secondary hover:text-text-primary transition-smooth mb-4"
      >
        <ArrowLeft :size="20" />
        <span>返回工作台</span>
      </button>

      <div class="flex items-center gap-3">
        <div class="w-12 h-12 bg-accent rounded-xl flex items-center justify-center">
          <Eye :size="24" class="text-white" />
        </div>
        <div>
          <h1 class="text-2xl font-bold text-text-primary">UI 质量审查</h1>
          <p class="text-text-secondary">基于截图或代码的 AI UI/UX 质量分析</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Input Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">输入参数</h2>

        <!-- Title -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">标题 *</label>
          <input
            v-model="title"
            type="text"
            class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none"
            placeholder="输入分析标题..."
          />
        </div>

        <!-- Review Mode -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">分析模式 *</label>
          <div class="flex gap-2">
            <button
              @click="reviewMode = 'screenshot'"
              :class="[
                'flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth',
                reviewMode === 'screenshot'
                  ? 'bg-accent text-white'
                  : 'bg-surface-muted text-text-secondary hover:bg-border',
              ]"
            >
              <Image :size="18" />
              <span>截图</span>
            </button>
            <button
              @click="reviewMode = 'code'"
              :class="[
                'flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth',
                reviewMode === 'code'
                  ? 'bg-accent text-white'
                  : 'bg-surface-muted text-text-secondary hover:bg-border',
              ]"
            >
              <Code :size="18" />
              <span>代码</span>
            </button>
            <button
              @click="reviewMode = 'screenshot_code'"
              :class="[
                'flex items-center gap-2 px-4 py-2 rounded-lg transition-smooth',
                reviewMode === 'screenshot_code'
                  ? 'bg-accent text-white'
                  : 'bg-surface-muted text-text-secondary hover:bg-border',
              ]"
            >
              <Eye :size="18" />
              <span>两者</span>
            </button>
          </div>
        </div>

        <!-- Screenshot Upload -->
        <div v-if="reviewMode === 'screenshot' || reviewMode === 'screenshot_code'" class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">上传截图</label>
          <div
            v-if="!screenshotPreview"
            data-testid="screenshot-upload-zone"
            role="button"
            tabindex="0"
            class="border-2 border-dashed border-border rounded-lg p-8 text-center hover:border-accent focus:border-accent focus:outline-none transition-smooth cursor-pointer"
            @click="($refs.fileInput as HTMLInputElement).click()"
            @paste="handleScreenshotPaste"
          >
            <Upload :size="32" class="text-text-muted mx-auto mb-2" />
            <p class="text-text-secondary">点击、拖拽或 Ctrl+V 粘贴截图</p>
            <p class="text-text-muted text-sm mt-1">支持 PNG, JPG, WebP (最大 20MB)</p>
            <input
              ref="fileInput"
              type="file"
              accept="image/png,image/jpeg,image/webp"
              class="hidden"
              @change="handleFileSelect"
            />
          </div>
          <div v-else class="relative">
            <img
              :src="screenshotPreview"
              class="w-full rounded-lg border border-border"
              style="max-height: 300px; object-fit: contain;"
            />
            <button
              @click="clearFile"
              class="absolute top-2 right-2 p-1 bg-danger text-white rounded hover:bg-danger/80 transition-smooth"
            >
              <AlertCircle :size="16" />
            </button>
          </div>
        </div>

        <!-- Code Input -->
        <div v-if="reviewMode === 'code' || reviewMode === 'screenshot_code'" class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">前端代码</label>
          <textarea
            v-model="codeInput"
            class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none font-mono text-sm"
            rows="8"
            placeholder="粘贴 Vue/React/HTML/CSS 代码..."
          ></textarea>
        </div>

        <!-- Optional Fields -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">页面类型 (可选)</label>
          <input
            v-model="pageType"
            type="text"
            class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none"
            placeholder="如: 登录页、Dashboard、表单页..."
          />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">目标风格 (可选)</label>
          <input
            v-model="targetStyle"
            type="text"
            class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none"
            placeholder="如: 简洁现代、Material Design..."
          />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">补充说明 (可选)</label>
          <textarea
            v-model="description"
            class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none"
            rows="3"
            placeholder="描述设计目标、用户场景等..."
          ></textarea>
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
          <button
            @click="handleSubmit"
            :disabled="loading || !canSubmit"
            :class="[
              'flex items-center gap-2 px-6 py-2 rounded-lg transition-smooth',
              loading || !canSubmit
                ? 'bg-surface-muted text-text-muted cursor-not-allowed'
                : 'bg-accent text-white hover:bg-accent/80',
            ]"
          >
            <Loader2 v-if="loading" :size="18" class="animate-spin" />
            <Eye v-else :size="18" />
            <span>{{ loading ? '分析中...' : '开始分析' }}</span>
          </button>

          <button
            @click="resetForm"
            class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth"
          >
            重置
          </button>
        </div>
      </div>

      <!-- Result Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">分析结果</h2>

        <!-- No Result -->
        <div v-if="!result && !loading" class="text-center py-12">
          <Eye :size="48" class="text-text-muted mx-auto mb-4" />
          <p class="text-text-secondary">提交输入后开始分析</p>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="text-center py-12">
          <Loader2 :size="48" class="text-accent mx-auto mb-4 animate-spin" />
          <p class="text-text-secondary">AI 正在分析...</p>
        </div>

        <!-- Result -->
        <div v-if="result && !loading">
          <!-- Summary -->
          <div class="mb-6 p-4 bg-surface-muted rounded-lg">
            <div class="flex items-center justify-between mb-2">
              <span class="text-text-secondary">总体评分</span>
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
              <div
                v-for="score in result.report_data.scores"
                :key="score.name"
                class="flex items-center justify-between p-2 bg-surface-muted rounded"
              >
                <span class="text-text-secondary">{{ getScoreDisplayName(score.name) }}</span>
                <div class="flex items-center gap-2">
                  <span :class="getScoreColor(score.score, score.max_score)">{{ score.score }}/{{ score.max_score }}</span>
                </div>
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
                  <span :class="getSeverityColor(issue.severity)" class="text-sm">{{ getSeverityDisplayName(issue.severity) }}</span>
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
              <button
                v-for="file in result.generated_files"
                :key="file.id"
                @click="downloadFile(file.filename)"
                class="flex items-center gap-1 px-3 py-1 bg-accent-soft text-accent rounded hover:bg-accent hover:text-white transition-smooth text-sm"
              >
                <FileText :size="14" />
                <span>{{ file.filename }}</span>
              </button>
            </div>
          </div>

          <!-- Export -->
          <button
            @click="exportMarkdown"
            class="flex items-center gap-2 px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth"
          >
            <Download :size="18" />
            <span>导出 Markdown</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

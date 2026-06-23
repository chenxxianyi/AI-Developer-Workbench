<script setup lang="ts">
/**
 * UI Review Page
 * Upload screenshot or code for UI/UX analysis
 */

import { ref, computed } from 'vue'
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
  Download,
} from '@lucide/vue'

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
      <div class="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div class="flex items-start gap-3">
          <div class="w-12 h-12 bg-accent rounded-xl flex items-center justify-center shrink-0 shadow-sm">
          <Eye :size="24" class="text-white" />
          </div>
          <div>
            <h1 class="text-3xl font-bold tracking-tight text-text-primary">UI 质量审查</h1>
            <p class="text-text-secondary mt-1 max-w-2xl">
              上传截图或前端代码，快速获得视觉层级、一致性、可访问性与改进建议。
            </p>
          </div>
        </div>

        <div class="inline-flex items-center rounded-full border border-accent/20 bg-accent-soft px-4 py-2 text-sm font-medium text-accent">
          填写输入 → 上传素材 → 开始分析 → 查看结果
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Input Panel -->
      <div class="bg-surface border border-border rounded-xl p-6 shadow-sm">
        <div class="mb-5">
          <h2 class="text-lg font-semibold text-text-primary">1. 输入与素材</h2>
          <p class="text-sm text-text-secondary mt-1">提供最关键的上下文，AI 会按当前材料生成审查报告。</p>
        </div>

        <!-- Title -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">标题 *</label>
          <input
            v-model="title"
            type="text"
            class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
            placeholder="输入分析标题..."
          />
        </div>

        <!-- Review Mode -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">分析模式 *</label>
          <p class="text-sm text-text-muted mb-3">选择最贴近当前材料的分析方式。</p>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <button
              @click="reviewMode = 'screenshot'"
              :class="[
                'flex items-start gap-2 px-4 py-3 rounded-lg transition-smooth text-left cursor-pointer focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
                reviewMode === 'screenshot'
                  ? 'bg-accent text-white'
                  : 'bg-surface-muted text-text-secondary hover:bg-border',
              ]"
            >
              <Image :size="18" class="mt-0.5 shrink-0" />
              <span>
                <span class="block font-semibold">截图</span>
                <span class="block text-xs opacity-85 mt-1">仅上传截图，快速评估视觉质量</span>
              </span>
            </button>
            <button
              @click="reviewMode = 'code'"
              :class="[
                'flex items-start gap-2 px-4 py-3 rounded-lg transition-smooth text-left cursor-pointer focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
                reviewMode === 'code'
                  ? 'bg-accent text-white'
                  : 'bg-surface-muted text-text-secondary hover:bg-border',
              ]"
            >
              <Code :size="18" class="mt-0.5 shrink-0" />
              <span>
                <span class="block font-semibold">代码</span>
                <span class="block text-xs opacity-85 mt-1">仅粘贴前端代码，审查结构与样式</span>
              </span>
            </button>
            <button
              @click="reviewMode = 'screenshot_code'"
              :class="[
                'flex items-start gap-2 px-4 py-3 rounded-lg transition-smooth text-left cursor-pointer focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
                reviewMode === 'screenshot_code'
                  ? 'bg-accent text-white'
                  : 'bg-surface-muted text-text-secondary hover:bg-border',
              ]"
            >
              <Eye :size="18" class="mt-0.5 shrink-0" />
              <span>
                <span class="block font-semibold">两者</span>
                <span class="block text-xs opacity-85 mt-1">截图 + 代码，获得更完整建议</span>
              </span>
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
            class="border-2 border-dashed border-border/80 rounded-lg p-8 text-center hover:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:border-accent focus:outline-none transition-smooth cursor-pointer bg-surface-muted/40"
            @click="($refs.fileInput as HTMLInputElement).click()"
            @paste="handleScreenshotPaste"
          >
            <Upload :size="32" class="text-accent mx-auto mb-2" />
            <p class="text-text-primary font-medium">点击、拖拽或 Ctrl+V 粘贴截图</p>
            <p class="text-text-secondary text-sm mt-1">支持 PNG, JPG, WebP (最大 20MB)</p>
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
            class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none font-mono text-sm text-text-primary placeholder:text-text-muted"
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
            class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
            placeholder="如: 登录页、Dashboard、表单页..."
          />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">目标风格 (可选)</label>
          <input
            v-model="targetStyle"
            type="text"
            class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
            placeholder="如: 简洁现代、Material Design..."
          />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">补充说明 (可选)</label>
          <textarea
            v-model="description"
            class="w-full px-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
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
        <div class="flex flex-col gap-3 border-t border-border pt-4 sm:flex-row sm:items-center sm:justify-between">
          <p v-if="!canSubmit && !loading" class="text-sm text-text-secondary">
            请先填写标题并提供截图或代码。
          </p>
          <p v-else class="text-sm text-text-muted">
            准备好后开始分析，结果会显示在右侧。
          </p>
          <div class="flex gap-3 sm:justify-end">
            <button
              @click="handleSubmit"
              :disabled="loading || !canSubmit"
              :class="[
                'flex items-center gap-2 px-7 py-3 rounded-lg font-semibold shadow-sm transition-smooth focus-visible:ring-2 focus-visible:ring-accent focus:outline-none',
                loading || !canSubmit
                  ? 'bg-surface-muted text-text-muted cursor-not-allowed shadow-none'
                  : 'bg-accent text-white hover:bg-accent/80 cursor-pointer',
              ]"
            >
              <Loader2 v-if="loading" :size="18" class="animate-spin" />
              <Eye v-else :size="18" />
              <span>{{ loading ? '分析中...' : '开始分析' }}</span>
            </button>
          <button
            @click="resetForm"
              class="px-4 py-3 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth cursor-pointer focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
          >
            重置
          </button>
          </div>
        </div>
      </div>

      <!-- Result Panel -->
      <div class="bg-surface border border-border rounded-xl p-6 shadow-sm">
        <h2 class="text-lg font-semibold text-text-primary mb-4">2. 分析结果</h2>

        <!-- No Result -->
        <div v-if="!result && !loading" class="py-8">
          <div class="rounded-xl border border-accent/20 bg-accent-soft/60 p-5">
            <div class="flex items-center gap-3 mb-3">
              <div class="w-10 h-10 rounded-lg bg-accent text-white flex items-center justify-center">
                <Eye :size="20" />
              </div>
              <div>
                <h3 class="font-semibold text-text-primary">将输出哪些内容</h3>
                <p class="text-sm text-text-secondary">提交后会在这里生成可执行的 UI 审查结论。</p>
              </div>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div class="rounded-lg bg-surface border border-border p-3">
                <p class="font-medium text-text-primary">评分维度</p>
                <p class="text-sm text-text-secondary mt-1">视觉层级、一致性、可访问性、对比度等。</p>
              </div>
              <div class="rounded-lg bg-surface border border-border p-3">
                <p class="font-medium text-text-primary">问题优先级</p>
                <p class="text-sm text-text-secondary mt-1">按 high / medium / low 聚焦关键问题。</p>
              </div>
              <div class="rounded-lg bg-surface border border-border p-3">
                <p class="font-medium text-text-primary">改进建议</p>
                <p class="text-sm text-text-secondary mt-1">给出可落地的设计和实现建议。</p>
              </div>
              <div class="rounded-lg bg-surface border border-border p-3">
                <p class="font-medium text-text-primary">报告文件</p>
                <p class="text-sm text-text-secondary mt-1">分析完成后可导出 Markdown 报告。</p>
              </div>
            </div>
          </div>
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

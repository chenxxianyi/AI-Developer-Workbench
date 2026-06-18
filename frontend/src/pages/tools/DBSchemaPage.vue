<script setup lang="ts">
/**
 * DB Schema Review Page
 * Review and optimize database schema
 */

import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { runDBSchema } from '@/api/tools'
import type { DBSchemaInput } from '@/types/tool'
import type { Report, DBSchemaResult } from '@/types/report'
import {
  Database,
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
const schemaType = ref('')
const databaseType = ref('')
const businessContext = ref('')
const schemaContent = ref('')
const targetGoal = ref('')

// Result state
const loading = ref(false)
const error = ref<string | null>(null)
const result = ref<Report<DBSchemaResult> | null>(null)

// Submit
async function handleSubmit() {
  if (!title.value.trim()) {
    error.value = '请输入标题'
    return
  }
  if (!schemaContent.value.trim()) {
    error.value = '请输入数据库 Schema'
    return
  }

  loading.value = true
  error.value = null
  result.value = null

  try {
    const payload: DBSchemaInput = {
      title: title.value,
      schema_type: schemaType.value || 'SQL',
      schema_content: schemaContent.value,
      database_type: databaseType.value || undefined,
      business_context: businessContext.value || undefined,
      target_goal: targetGoal.value || undefined,
    }

    result.value = await runDBSchema(payload) as Report<DBSchemaResult>
  } catch (err: any) {
    error.value = err.message || '分析失败'
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
  schemaType.value = ''
  databaseType.value = ''
  businessContext.value = ''
  schemaContent.value = ''
  targetGoal.value = ''
  result.value = null
  error.value = null
}

const canSubmit = computed(() => title.value.trim() && schemaContent.value.trim())
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
        <div class="w-12 h-12 bg-accent rounded-xl flex items-center justify-center">
          <Database :size="24" class="text-white" />
        </div>
        <div>
          <h1 class="text-2xl font-bold text-text-primary">DB Schema Review</h1>
          <p class="text-text-secondary">Review and optimize database schema with AI analysis</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Input Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">输入参数</h2>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">标题 *</label>
          <input v-model="title" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="输入分析标题..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">Schema 类型</label>
          <input v-model="schemaType" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: SQL, GORM, Prisma..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">数据库类型</label>
          <input v-model="databaseType" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: MySQL, PostgreSQL..." />
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">数据库 Schema *</label>
          <textarea v-model="schemaContent" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none font-mono text-sm" rows="12" placeholder="粘贴 CREATE TABLE、GORM model、Prisma schema 等..."></textarea>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">业务背景</label>
          <textarea v-model="businessContext" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" rows="3" placeholder="描述业务场景和数据使用方式..."></textarea>
        </div>

        <div class="mb-4">
          <label class="block text-sm font-medium text-text-secondary mb-2">优化目标</label>
          <input v-model="targetGoal" type="text" class="w-full px-4 py-2 bg-surface-muted border border-border rounded-lg focus:border-accent focus:outline-none" placeholder="如: 性能优化、扩展性、数据完整性..." />
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
            <Database v-else :size="18" />
            <span>{{ loading ? '分析中...' : '开始分析' }}</span>
          </button>
          <button @click="resetForm" class="px-4 py-2 bg-surface-muted text-text-secondary rounded-lg hover:bg-border transition-smooth">重置</button>
        </div>
      </div>

      <!-- Result Panel -->
      <div class="bg-surface border border-border rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">分析结果</h2>

        <div v-if="!result && !loading" class="text-center py-12">
          <Database :size="48" class="text-text-muted mx-auto mb-4" />
          <p class="text-text-secondary">输入 Schema 后开始分析</p>
        </div>

        <div v-if="loading" class="text-center py-12">
          <Loader2 :size="48" class="text-accent mx-auto mb-4 animate-spin" />
          <p class="text-text-secondary">AI 正在分析 Schema...</p>
        </div>

        <div v-if="result && !loading">
          <!-- Summary -->
          <div class="mb-6 p-4 bg-surface-muted rounded-lg">
            <div class="flex items-center justify-between mb-2">
              <span class="text-text-secondary">评分</span>
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

          <!-- Optimized Schema -->
          <div v-if="result.report_data?.optimized_schema" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">优化后的 Schema</h3>
            <pre class="p-3 bg-surface-muted rounded text-sm text-text-secondary overflow-x-auto whitespace-pre-wrap">{{ result.report_data.optimized_schema }}</pre>
          </div>

          <!-- Migration Suggestions -->
          <div v-if="result.report_data?.migration_suggestions?.length" class="mb-6">
            <h3 class="text-md font-semibold text-text-primary mb-3">迁移建议</h3>
            <ul class="space-y-1">
              <li v-for="(sug, idx) in result.report_data.migration_suggestions" :key="idx" class="flex items-start gap-2 text-text-secondary">
                <CheckCircle2 :size="14" class="text-accent mt-1" />
                <span class="text-sm">{{ sug }}</span>
              </li>
            </ul>
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
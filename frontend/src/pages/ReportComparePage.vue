<script setup lang="ts">
/**
 * Report Compare Page
 * Shows the delta between two same-tool reports (baseline → target).
 */
import { onMounted, ref, computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { compareReports } from '@/api/reports'
import { getToolDisplayMeta } from '@/utils/toolDisplay'
import type { ReportCompare } from '@/types/report'
import { ArrowLeft, TrendingUp, TrendingDown, Minus, AlertTriangle } from '@lucide/vue'

const route = useRoute()

const loading = ref(false)
const error = ref<string | null>(null)
const data = ref<ReportCompare | null>(null)

const baselineId = computed(() => route.params.id as string)
const targetId = computed(() => route.params.targetId as string)
const toolMeta = computed(() => (data.value ? getToolDisplayMeta(data.value.tool_type) : null))

onMounted(async () => {
  if (!baselineId.value || !targetId.value) {
    error.value = '缺少报告 ID'
    return
  }
  loading.value = true
  try {
    data.value = await compareReports(baselineId.value, targetId.value)
  } catch (err: any) {
    error.value = err.message || '对比加载失败'
  } finally {
    loading.value = false
  }
})

function scoreTrend(delta: number | null | undefined): 'up' | 'down' | 'flat' | 'unknown' {
  if (delta === null || delta === undefined) return 'unknown'
  if (delta > 0) return 'up'
  if (delta < 0) return 'down'
  return 'flat'
}

function formatDelta(n: number): string {
  return n > 0 ? `+${n}` : `${n}`
}
</script>

<template>
  <div class="max-w-5xl mx-auto">
    <!-- Header -->
    <div class="mb-6">
      <RouterLink to="/reports" class="inline-flex items-center gap-1.5 text-sm text-text-muted hover:text-accent transition-smooth mb-4">
        <ArrowLeft :size="16" />
        返回列表
      </RouterLink>
      <h1 class="text-2xl font-bold text-text-primary">报告对比</h1>
      <p v-if="toolMeta" class="text-text-secondary mt-1">{{ toolMeta.name }} · 基线 → 目标</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="py-16 text-center">
      <div class="animate-spin inline-block w-8 h-8 border-2 border-accent border-t-transparent rounded-full mb-3" />
      <p class="text-text-muted">正在生成对比…</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="rounded-lg border border-danger/30 bg-danger/5 p-5">
      <div class="flex items-center gap-2 text-danger">
        <AlertTriangle :size="18" />
        <span class="font-medium">对比失败</span>
      </div>
      <p class="text-sm text-text-muted mt-1">{{ error }}</p>
    </div>

    <!-- Compare content -->
    <div v-else-if="data" class="space-y-6">
      <!-- Warnings -->
      <div v-if="data.warnings?.length" class="rounded-lg border border-amber-300/40 bg-amber-50 dark:bg-amber-900/10 p-4">
        <div v-for="w in data.warnings" :key="w" class="flex items-center gap-2 text-sm text-amber-700 dark:text-amber-300">
          <AlertTriangle :size="14" />
          {{ w }}
        </div>
      </div>

      <!-- Report cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="rounded-lg border border-border bg-surface p-5">
          <p class="text-xs text-text-muted mb-1">基线报告</p>
          <RouterLink :to="`/reports/${data.baseline_report.id}`" class="font-semibold text-text-primary hover:text-accent">
            {{ data.baseline_report.title }}
          </RouterLink>
          <p class="text-sm text-text-secondary mt-1">{{ data.baseline_report.summary }}</p>
          <p class="text-sm mt-2">
            <span class="text-text-muted">评分</span>
            <span class="ml-2 font-semibold">{{ data.baseline_report.total_score ?? '—' }}</span>
            <span v-if="data.baseline_report.grade" class="ml-2">{{ data.baseline_report.grade }}</span>
          </p>
        </div>
        <div class="rounded-lg border border-border bg-surface p-5">
          <p class="text-xs text-text-muted mb-1">目标报告</p>
          <RouterLink :to="`/reports/${data.target_report.id}`" class="font-semibold text-text-primary hover:text-accent">
            {{ data.target_report.title }}
          </RouterLink>
          <p class="text-sm text-text-secondary mt-1">{{ data.target_report.summary }}</p>
          <p class="text-sm mt-2">
            <span class="text-text-muted">评分</span>
            <span class="ml-2 font-semibold">{{ data.target_report.total_score ?? '—' }}</span>
            <span v-if="data.target_report.grade" class="ml-2">{{ data.target_report.grade }}</span>
          </p>
        </div>
      </div>

      <!-- Score delta -->
      <div class="rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-3">分数变化</h3>
        <div v-if="scoreTrend(data.score_delta) === 'unknown'" class="text-text-muted text-sm">不可比较（报告无评分）</div>
        <div v-else class="flex items-center gap-2">
          <TrendingUp v-if="scoreTrend(data.score_delta) === 'up'" :size="20" class="text-success" />
          <TrendingDown v-else-if="scoreTrend(data.score_delta) === 'down'" :size="20" class="text-danger" />
          <Minus v-else :size="20" class="text-text-muted" />
          <span class="text-xl font-bold" :class="scoreTrend(data.score_delta) === 'up' ? 'text-success' : scoreTrend(data.score_delta) === 'down' ? 'text-danger' : 'text-text-muted'">
            {{ formatDelta(data.score_delta ?? 0) }}
          </span>
          <span v-if="data.grade_delta" class="text-sm text-text-muted ml-2">{{ data.grade_delta }}</span>
        </div>
      </div>

      <!-- Issue count delta -->
      <div class="rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-3">问题数量变化</h3>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div class="rounded-md bg-surface-muted p-3">
            <p class="text-xs text-text-muted">高危</p>
            <p class="text-lg font-bold" :class="data.issue_count_delta.high < 0 ? 'text-success' : data.issue_count_delta.high > 0 ? 'text-danger' : ''">{{ formatDelta(data.issue_count_delta.high) }}</p>
          </div>
          <div class="rounded-md bg-surface-muted p-3">
            <p class="text-xs text-text-muted">中危</p>
            <p class="text-lg font-bold" :class="data.issue_count_delta.medium < 0 ? 'text-success' : data.issue_count_delta.medium > 0 ? 'text-danger' : ''">{{ formatDelta(data.issue_count_delta.medium) }}</p>
          </div>
          <div class="rounded-md bg-surface-muted p-3">
            <p class="text-xs text-text-muted">低危</p>
            <p class="text-lg font-bold">{{ formatDelta(data.issue_count_delta.low) }}</p>
          </div>
          <div class="rounded-md bg-surface-muted p-3">
            <p class="text-xs text-text-muted">总数</p>
            <p class="text-lg font-bold" :class="data.issue_count_delta.total < 0 ? 'text-success' : data.issue_count_delta.total > 0 ? 'text-danger' : ''">{{ formatDelta(data.issue_count_delta.total) }}</p>
          </div>
        </div>
      </div>

      <!-- Issues breakdown -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="rounded-lg border border-success/30 bg-success/5 p-4">
          <h4 class="font-semibold text-success mb-2">已解决 ({{ data.issues.resolved.length }})</h4>
          <ul class="space-y-1 text-sm">
            <li v-for="(it, i) in data.issues.resolved" :key="i" class="text-text-secondary">
              <span class="text-text-muted">[{{ it.severity }}]</span> {{ it.title }}
            </li>
            <li v-if="!data.issues.resolved.length" class="text-text-muted">无</li>
          </ul>
        </div>
        <div class="rounded-lg border border-danger/30 bg-danger/5 p-4">
          <h4 class="font-semibold text-danger mb-2">新增 ({{ data.issues.new.length }})</h4>
          <ul class="space-y-1 text-sm">
            <li v-for="(it, i) in data.issues.new" :key="i" class="text-text-secondary">
              <span class="text-text-muted">[{{ it.severity }}]</span> {{ it.title }}
            </li>
            <li v-if="!data.issues.new.length" class="text-text-muted">无</li>
          </ul>
        </div>
        <div class="rounded-lg border border-border bg-surface-muted p-4">
          <h4 class="font-semibold text-text-primary mb-2">持续 ({{ data.issues.persist.length }})</h4>
          <ul class="space-y-1 text-sm">
            <li v-for="(it, i) in data.issues.persist" :key="i" class="text-text-secondary">
              <span class="text-text-muted">[{{ it.severity }}]</span> {{ it.title }}
            </li>
            <li v-if="!data.issues.persist.length" class="text-text-muted">无</li>
          </ul>
        </div>
      </div>

      <!-- Action items delta -->
      <div class="rounded-lg border border-border bg-surface p-5">
        <h3 class="text-lg font-semibold text-text-primary mb-3">行动项变化</h3>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 text-sm">
          <div>
            <p class="text-success font-medium mb-1">已解决 ({{ data.action_items.resolved.length }})</p>
            <ul class="space-y-1">
              <li v-for="(a, i) in data.action_items.resolved" :key="i" class="text-text-secondary">{{ a.title }}</li>
            </ul>
          </div>
          <div>
            <p class="text-danger font-medium mb-1">新增 ({{ data.action_items.new.length }})</p>
            <ul class="space-y-1">
              <li v-for="(a, i) in data.action_items.new" :key="i" class="text-text-secondary">{{ a.title }}</li>
            </ul>
          </div>
          <div>
            <p class="text-text-muted font-medium mb-1">持续 ({{ data.action_items.persist.length }})</p>
            <ul class="space-y-1">
              <li v-for="(a, i) in data.action_items.persist" :key="i" class="text-text-secondary">{{ a.title }}</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * Reports Page
 * List all reports with filtering, sorting, and pagination.
 */
import { onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReportStore } from '@/stores/reportStore'
import type { ToolType } from '@/types/tool'
import type { ReportStatus } from '@/types/report'
import { getToolDisplayMeta } from '@/utils/toolDisplay'
import { Filter, ArrowUpDown, AlertCircle, FileText, RefreshCw } from '@lucide/vue'
import ReportListCard from '@/components/report/ReportListCard.vue'
import PaginationBar from '@/components/common/PaginationBar.vue'

const store = useReportStore()
const route = useRoute()
const router = useRouter()

// Sync URL query params to store.
function syncFromURL() {
  const q = route.query
  if (q.tool_type && typeof q.tool_type === 'string') {
    store.toolFilter = q.tool_type as ToolType
  }
  if (q.status && typeof q.status === 'string') {
    store.statusFilter = q.status as ReportStatus
  }
  if (q.sort && typeof q.sort === 'string') {
    store.sort = q.sort as 'newest' | 'oldest' | 'score_desc' | 'score_asc'
  }
  if (q.page && typeof q.page === 'string') {
    store.currentPage = parseInt(q.page, 10) || 1
  }
}

// Push current filters to URL.
function pushToURL() {
  const q: Record<string, string> = {}
  if (store.toolFilter) q.tool_type = store.toolFilter
  if (store.statusFilter) q.status = store.statusFilter
  if (store.sort !== 'newest') q.sort = store.sort
  if (store.currentPage > 1) q.page = String(store.currentPage)
  router.replace({ query: q })
}

// Watch store changes and push to URL.
watch(
  () => [store.toolFilter, store.statusFilter, store.sort, store.currentPage],
  () => pushToURL(),
  { deep: true },
)

onMounted(async () => {
  syncFromURL()
  await store.fetchReports()
  pushToURL()
})

const resetFilters = () => {
  store.setFilter('')
  store.setStatusFilter('')
}
</script>

<template>
  <div>
    <!-- Header -->
    <div class="mb-6 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <h2 class="text-2xl font-bold text-text-primary">报告列表</h2>
        <p class="mt-1 text-sm text-text-muted">管理所有分析报告，支持筛选、排序和分页。</p>
      </div>
      <button
        class="inline-flex items-center gap-2 self-start rounded-md border border-border bg-surface px-4 py-2 text-sm font-medium text-text-primary transition-smooth hover:bg-surface-hover"
        @click="store.fetchReports()"
        :disabled="store.loading"
        aria-label="刷新列表"
      >
        <RefreshCw :size="15" :class="{ 'animate-spin': store.loading }" />
        刷新
      </button>
    </div>

    <!-- Filter Bar -->
    <div class="mb-5 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div class="flex flex-wrap items-center gap-2">
        <Filter :size="16" class="text-text-muted" aria-hidden="true" />
        <!-- Tool Type Filter -->
        <select
          :value="store.toolFilter"
          @change="store.setFilter(($event.target as HTMLSelectElement).value as ToolType | '')"
          class="rounded-md border border-border bg-surface px-3 py-1.5 text-sm text-text-primary focus:border-accent focus:ring-1 focus:ring-accent focus:outline-none"
          aria-label="按工具类型筛选"
        >
          <option value="">全部工具</option>
          <option v-for="tt in (['ui_review','project_doctor','agent_config','api_doc','db_schema'] as ToolType[])" :key="tt" :value="tt">
            {{ getToolDisplayMeta(tt).name }}
          </option>
        </select>
        <!-- Status Filter -->
        <select
          :value="store.statusFilter"
          @change="store.setStatusFilter(($event.target as HTMLSelectElement).value as ReportStatus | '')"
          class="rounded-md border border-border bg-surface px-3 py-1.5 text-sm text-text-primary focus:border-accent focus:ring-1 focus:ring-accent focus:outline-none"
          aria-label="按状态筛选"
        >
          <option value="">全部状态</option>
          <option value="succeeded">成功</option>
          <option value="failed">失败</option>
          <option value="processing">处理中</option>
        </select>
        <!-- Reset filters -->
        <button
          v-if="store.toolFilter || store.statusFilter"
          class="text-xs text-text-muted hover:text-accent underline px-2"
          @click="resetFilters"
        >
          清除筛选
        </button>
      </div>

      <!-- Sort -->
      <div class="flex items-center gap-2">
        <ArrowUpDown :size="16" class="text-text-muted" aria-hidden="true" />
        <select
          :value="store.sort"
          @change="store.setSort(($event.target as HTMLSelectElement).value as 'newest' | 'oldest' | 'score_desc' | 'score_asc')"
          class="rounded-md border border-border bg-surface px-3 py-1.5 text-sm text-text-primary focus:border-accent focus:ring-1 focus:ring-accent focus:outline-none"
          aria-label="排序方式"
        >
          <option value="newest">最新优先</option>
          <option value="oldest">最早优先</option>
          <option value="score_desc">评分从高到低</option>
          <option value="score_asc">评分从低到高</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="store.loading" class="rounded-lg border border-dashed border-border py-16 text-center">
      <RefreshCw :size="28" class="mx-auto mb-3 text-text-muted animate-spin" />
      <p class="text-text-muted">加载中...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="store.error" class="rounded-lg border border-danger/30 bg-danger/5 px-5 py-8 text-center">
      <AlertCircle :size="28" class="mx-auto mb-3 text-danger" />
      <p class="text-danger font-medium">{{ store.error }}</p>
      <button
        class="mt-3 text-sm text-accent hover:underline"
        @click="store.fetchReports()"
      >
        重试
      </button>
    </div>

    <!-- Empty: No reports at all -->
    <div v-else-if="store.total === 0 && !store.toolFilter && !store.statusFilter" class="rounded-lg border border-dashed border-border py-16 text-center">
      <FileText :size="36" class="mx-auto mb-4 text-text-muted" />
      <h3 class="text-lg font-semibold text-text-primary mb-1">暂无报告</h3>
      <p class="text-sm text-text-muted">运行任一工具后，报告将出现在这里。</p>
    </div>

    <!-- Empty: No matching results -->
    <div v-else-if="store.reports.length === 0" class="rounded-lg border border-dashed border-border py-16 text-center">
      <FileText :size="36" class="mx-auto mb-4 text-text-muted" />
      <h3 class="text-lg font-semibold text-text-primary mb-1">没有匹配的报告</h3>
      <p class="text-sm text-text-muted mb-4">尝试调整筛选条件。</p>
      <button
        class="text-sm text-accent hover:underline"
        @click="resetFilters"
      >
        清除所有筛选
      </button>
    </div>

    <!-- Report List -->
    <div v-else class="space-y-3">
      <TransitionGroup name="list" tag="div" class="space-y-3">
        <ReportListCard
          v-for="report in store.reports"
          :key="report.id"
          :report="report"
        />
      </TransitionGroup>

      <!-- Pagination -->
      <PaginationBar
        :current-page="store.currentPage"
        :total-pages="store.totalPages"
        :has-prev="store.hasPrevPage"
        :has-next="store.hasNextPage"
        @prev="store.setPage(store.currentPage - 1)"
        @next="store.setPage(store.currentPage + 1)"
      />
    </div>
  </div>
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.2s ease;
}
.list-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.list-leave-to {
  opacity: 0;
}

@media (max-width: 640px) {
  .animate-spin {
    animation: spin 1s linear infinite;
  }
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>

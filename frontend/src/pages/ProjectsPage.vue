<script setup lang="ts">
/**
 * Projects Page
 * List, search, create entry, navigate to detail
 */
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/projectStore'
import PaginationBar from '@/components/common/PaginationBar.vue'
import { FolderPlus, Search, Folder, ArrowRight } from '@lucide/vue'

const router = useRouter()
const store = useProjectStore()

onMounted(async () => {
  await store.fetchProjects()
  store.recalcTotalPages()
})

const totalPages = computed(() => Math.ceil(store.total / store.pageSize) || 1)

function goCreate() {
  router.push('/projects/new')
}

function goDetail(id: string) {
  router.push(`/projects/${id}`)
}

function onSearch(e: Event) {
  const target = e.target as HTMLInputElement
  store.setSearch(target.value)
}

function onPage(p: number) {
  store.setPage(p)
  store.recalcTotalPages()
}
</script>

<template>
  <div class="max-w-5xl mx-auto">
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">项目</h1>
        <p class="text-text-secondary mt-1">围绕项目沉淀报告、查看趋势和最新产物</p>
      </div>
      <button
        class="inline-flex items-center gap-2 px-4 py-2 bg-accent text-white rounded-lg hover:bg-accent/80 transition-smooth focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="goCreate"
      >
        <FolderPlus :size="18" />
        新建项目
      </button>
    </div>

    <!-- Search -->
    <div class="mb-4 relative">
      <Search :size="16" class="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
      <input
        type="text"
        placeholder="搜索项目名称或描述…"
        class="w-full pl-10 pr-4 py-2 bg-surface-muted border border-border/80 rounded-lg focus-visible:ring-2 focus-visible:ring-accent focus-visible:border-accent focus:outline-none text-text-primary placeholder:text-text-muted"
        :value="store.search"
        @input="onSearch"
      />
    </div>

    <!-- Loading -->
    <div v-if="store.loading" class="py-12 text-center text-text-muted">加载中…</div>

    <!-- Error -->
    <div v-else-if="store.error" class="rounded-lg border border-danger/30 bg-danger/5 p-4 text-danger">
      {{ store.error }}
    </div>

    <!-- Empty -->
    <div v-else-if="!store.projects.length" class="py-16 text-center">
      <Folder :size="48" class="mx-auto mb-3 text-text-muted" />
      <p class="text-text-muted">暂无项目，点击「新建项目」开始。</p>
    </div>

    <!-- List -->
    <div v-else class="space-y-3">
      <button
        v-for="p in store.projects"
        :key="p.id"
        class="group block w-full text-left rounded-lg border border-border bg-surface px-5 py-4 transition-smooth hover:border-accent/35 hover:bg-surface-hover focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="goDetail(p.id)"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0 flex-1">
            <h3 class="font-semibold text-text-primary group-hover:text-accent">{{ p.name }}</h3>
            <p v-if="p.description" class="text-sm text-text-muted mt-1 truncate">{{ p.description }}</p>
            <div class="flex items-center gap-4 mt-2 text-xs text-text-muted">
              <span v-if="p.repo_url">{{ p.repo_url }}</span>
              <span>报告 {{ p.report_count }}</span>
              <span v-if="p.average_score !== null">平均 {{ p.average_score.toFixed(0) }}</span>
            </div>
          </div>
          <ArrowRight :size="18" class="text-text-muted group-hover:text-accent flex-shrink-0" />
        </div>
      </button>
    </div>

    <div v-if="totalPages > 1" class="mt-6">
      <PaginationBar
        :current-page="store.currentPage"
        :total-pages="totalPages"
        :has-prev="store.currentPage > 1"
        :has-next="store.currentPage < totalPages"
        @prev="onPage(store.currentPage - 1)"
        @next="onPage(store.currentPage + 1)"
      />
    </div>
  </div>
</template>

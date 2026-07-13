<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, RouterLink, RouterView } from 'vue-router'
import AppBadge from '@/components/common/AppBadge.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import AppSkeleton from '@/components/common/AppSkeleton.vue'
import { LayoutDashboard, FileText, PenTool, Play, Monitor, Files, FileStack, Download, ExternalLink } from '@lucide/vue'

const route = useRoute()
const projectId = computed(() => route.params.projectId as string)

interface Project { id: string; name: string; type: string; status: string; quality_score?: number }
const project = ref<Project | null>(null)
const loading = ref(true)
const error = ref('')
const notFound = ref(false)

const navItems = [
  { to: `/projects/${projectId.value}`, icon: LayoutDashboard, label: '概览' },
  { to: `/projects/${projectId.value}/requirements`, icon: FileText, label: '需求' },
  { to: `/projects/${projectId.value}/blueprint`, icon: PenTool, label: '蓝图' },
  { to: `/projects/${projectId.value}/generation`, icon: Play, label: '生成' },
  { to: `/projects/${projectId.value}/preview`, icon: Monitor, label: '预览' },
  { to: `/projects/${projectId.value}/files`, icon: Files, label: '文件' },
  { to: `/projects/${projectId.value}/reports`, icon: FileStack, label: '报告' },
]

const statusVariant = {
  draft: 'default', generating: 'info', building: 'info', completed: 'success', failed: 'danger', archived: 'default',
} as const

onMounted(async () => {
  try {
    // TODO: Replace with actual API call when project service is unified
    project.value = { id: projectId.value, name: '加载中...', type: 'website', status: 'draft' }
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <!-- Loading -->
  <div v-if="loading" class="space-y-4">
    <AppSkeleton height="48px" />
    <AppSkeleton height="32px" width="60%" />
    <AppSkeleton height="200px" />
  </div>

  <!-- Error -->
  <ErrorState v-else-if="error" :message="error" retry-label="重试" @retry="loading = true; error = ''; /* retry logic */" />

  <!-- Not Found -->
  <EmptyState v-else-if="notFound" title="项目不存在" description="该项目可能已被删除或您没有访问权限">
    <template #action>
      <RouterLink to="/projects" class="text-sm text-accent hover:underline">返回项目列表</RouterLink>
    </template>
  </EmptyState>

  <!-- Project Layout -->
  <template v-else>
    <!-- Header -->
    <div class="flex items-start justify-between mb-6">
      <div>
        <div class="flex items-center gap-3 mb-1">
          <h1 class="text-2xl font-bold text-[var(--color-text-primary)]">{{ project?.name }}</h1>
          <AppBadge v-if="project?.status" :variant="(statusVariant as Record<string, any>)[project.status] || 'default'">{{ project.status }}</AppBadge>
        </div>
        <p class="text-sm text-[var(--color-text-muted)]">{{ project?.type === 'website' ? '网站生成' : project?.type === 'analysis' ? '项目分析' : '导入项目' }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button class="flex items-center gap-1.5 px-3 py-2 text-sm rounded-lg border border-[var(--color-border)] hover:bg-[var(--color-surface-muted)] transition-smooth">
          <Download :size="16" /><span>导出</span>
        </button>
        <button class="flex items-center gap-1.5 px-3 py-2 text-sm rounded-lg border border-[var(--color-border)] hover:bg-[var(--color-surface-muted)] transition-smooth">
          <ExternalLink :size="16" /><span>预览</span>
        </button>
      </div>
    </div>

    <!-- Sub-nav -->
    <nav class="flex gap-1 mb-6 border-b border-[var(--color-border)] overflow-x-auto">
      <RouterLink
        v-for="item in navItems" :key="item.label"
        :to="item.to"
        :class="[
          'flex items-center gap-2 px-4 py-2.5 text-sm font-medium whitespace-nowrap border-b-2 -mb-[1px] transition-colors duration-200',
          route.path === item.to || route.path.startsWith(item.to + '/')
            ? 'border-[var(--color-accent)] text-[var(--color-accent)]'
            : 'border-transparent text-[var(--color-text-muted)] hover:text-[var(--color-text-secondary)]',
        ]"
      >
        <component :is="item.icon" :size="16" />
        {{ item.label }}
      </RouterLink>
    </nav>

    <!-- Page Content -->
    <RouterView />
  </template>
</template>

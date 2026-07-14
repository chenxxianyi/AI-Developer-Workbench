<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import {
  CircleCheck,
  CircleX,
  ExternalLink,
  FolderKanban,
  Plus,
  Search,
  User,
} from '@lucide/vue'
import AdminPageShell from '@/components/admin/AdminPageShell.vue'

interface ManagedProject {
  id: string
  name: string
  user_id: string
  status: 'draft' | 'generating' | 'completed' | 'failed'
  created_at: string
}

const search = ref('')
const statusFilter = ref<'' | ManagedProject['status']>('')
const projects = ref<ManagedProject[]>([
  {
    id: '1',
    name: '示例项目',
    user_id: '',
    status: 'draft',
    created_at: '2026-07-13',
  },
])

const statusLabels: Record<ManagedProject['status'], string> = {
  draft: '草稿',
  generating: '生成中',
  completed: '已完成',
  failed: '失败',
}

const filteredProjects = computed(() => {
  const keyword = search.value.trim().toLowerCase()
  return projects.value.filter((project) => {
    const matchesSearch = !keyword
      || project.name.toLowerCase().includes(keyword)
      || project.user_id.toLowerCase().includes(keyword)
    const matchesStatus = !statusFilter.value || project.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

function formatDate(value: string) {
  return new Date(value).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}
</script>

<template>
  <AdminPageShell
    :icon="FolderKanban"
    title="项目管理"
    description="查看所有项目的归属、生命周期状态和创建时间。"
    badge-text="项目治理"
  >
    <template #actions>
      <RouterLink
        to="/projects/new"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
      >
        <Plus :size="16" />
        创建项目
      </RouterLink>
    </template>

    <section class="mb-5 rounded-lg border border-border bg-surface p-4 shadow-sm">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div class="relative w-full lg:max-w-md">
          <Search :size="17" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
          <input
            v-model="search"
            type="search"
            class="min-h-10 w-full rounded-lg border border-border bg-surface-muted pl-10 pr-4 text-sm text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            placeholder="搜索项目名称或用户 ID"
            aria-label="搜索项目"
          />
        </div>

        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <select
            v-model="statusFilter"
            class="min-h-10 rounded-lg border border-border bg-surface px-3 text-sm text-text-primary focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            aria-label="按项目状态筛选"
          >
            <option value="">全部状态</option>
            <option value="draft">草稿</option>
            <option value="generating">生成中</option>
            <option value="completed">已完成</option>
            <option value="failed">失败</option>
          </select>
          <span class="text-sm text-text-muted">{{ projects.length }} 个项目</span>
        </div>
      </div>
    </section>

    <section class="overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
      <div class="flex items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 class="font-semibold text-text-primary">项目列表</h2>
          <p class="mt-0.5 text-xs text-text-muted">跨用户查看全部项目工作区</p>
        </div>
        <span class="rounded-full border border-border bg-surface-muted px-2.5 py-1 text-xs font-semibold text-text-secondary">
          {{ filteredProjects.length }} 项
        </span>
      </div>

      <div v-if="filteredProjects.length" class="overflow-x-auto">
        <table class="w-full min-w-[800px] text-sm">
          <thead class="bg-surface-muted/70 text-xs text-text-muted">
            <tr>
              <th class="px-5 py-3 text-left font-semibold">项目</th>
              <th class="px-5 py-3 text-left font-semibold">所属用户</th>
              <th class="px-5 py-3 text-left font-semibold">状态</th>
              <th class="px-5 py-3 text-left font-semibold">创建时间</th>
              <th class="px-5 py-3 text-right font-semibold">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="project in filteredProjects"
              :key="project.id"
              class="border-t border-border transition-colors duration-200 hover:bg-surface-muted/50"
            >
              <td class="px-5 py-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent">
                    <FolderKanban :size="18" />
                  </div>
                  <div>
                    <RouterLink
                      :to="`/projects/${project.id}`"
                      class="font-semibold text-text-primary transition-colors duration-200 hover:text-accent"
                    >
                      {{ project.name }}
                    </RouterLink>
                    <p class="mt-0.5 font-mono text-xs text-text-muted">ID {{ project.id }}</p>
                  </div>
                </div>
              </td>
              <td class="px-5 py-4">
                <div class="flex items-center gap-2 text-text-secondary">
                  <User :size="15" class="text-text-muted" />
                  <span>{{ project.user_id || '未分配' }}</span>
                </div>
              </td>
              <td class="px-5 py-4">
                <span
                  :class="[
                    'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-semibold',
                    project.status === 'completed'
                      ? 'border-success/20 bg-success/10 text-success'
                      : project.status === 'failed'
                        ? 'border-danger/20 bg-danger/10 text-danger'
                        : project.status === 'generating'
                          ? 'border-accent/20 bg-accent-soft text-accent'
                          : 'border-border bg-surface-muted text-text-secondary',
                  ]"
                >
                  <CircleCheck v-if="project.status === 'completed'" :size="13" />
                  <CircleX v-else-if="project.status === 'failed'" :size="13" />
                  <span v-else class="h-1.5 w-1.5 rounded-full bg-current" />
                  {{ statusLabels[project.status] }}
                </span>
              </td>
              <td class="px-5 py-4 text-text-secondary">{{ formatDate(project.created_at) }}</td>
              <td class="px-5 py-4">
                <div class="flex items-center justify-end">
                  <RouterLink
                    :to="`/projects/${project.id}`"
                    title="打开项目"
                    aria-label="打开项目"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-accent-soft hover:text-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                  >
                    <ExternalLink :size="16" />
                  </RouterLink>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-else class="flex min-h-64 flex-col items-center justify-center px-6 text-center">
        <Search :size="30" class="mb-3 text-text-muted" />
        <p class="text-sm font-medium text-text-primary">没有匹配的项目</p>
        <p class="mt-1 text-xs text-text-muted">调整搜索关键词或项目状态。</p>
      </div>
    </section>
  </AdminPageShell>
</template>

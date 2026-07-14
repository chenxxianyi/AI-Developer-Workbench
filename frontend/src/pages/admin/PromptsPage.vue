<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  CircleCheck,
  CircleX,
  Eye,
  FileText,
  Pencil,
  Plus,
  Search,
  Trash2,
} from '@lucide/vue'
import AdminPageShell from '@/components/admin/AdminPageShell.vue'

interface PromptTemplate {
  id: string
  name: string
  type: string
  status: 'active' | 'disabled'
  version: number
}

const search = ref('')
const typeFilter = ref('')
const prompts = ref<PromptTemplate[]>([
  { id: '1', name: 'UI Review Default', type: 'ui_review', status: 'active', version: 1 },
])

const promptTypeLabels: Record<string, string> = {
  ui_review: 'UI 审查',
  project_doctor: '项目诊断',
  agent_config: 'Agent 配置',
  api_doc: 'API 文档',
  db_schema: '数据库审查',
}

const filteredPrompts = computed(() => {
  const keyword = search.value.trim().toLowerCase()
  return prompts.value.filter((prompt) => {
    const matchesSearch = !keyword
      || prompt.name.toLowerCase().includes(keyword)
      || prompt.type.toLowerCase().includes(keyword)
    const matchesType = !typeFilter.value || prompt.type === typeFilter.value
    return matchesSearch && matchesType
  })
})
</script>

<template>
  <AdminPageShell
    :icon="FileText"
    title="Prompt 模板管理"
    description="集中维护各类 AI 工具使用的提示词模板和版本。"
    badge-text="模板配置"
  >
    <template #actions>
      <button
        type="button"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
      >
        <Plus :size="16" />
        新建模板
      </button>
    </template>

    <section class="mb-5 rounded-lg border border-border bg-surface p-4 shadow-sm">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div class="relative w-full lg:max-w-md">
          <Search :size="17" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
          <input
            v-model="search"
            type="search"
            class="min-h-10 w-full rounded-lg border border-border bg-surface-muted pl-10 pr-4 text-sm text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            placeholder="搜索模板名称或类型"
            aria-label="搜索 Prompt 模板"
          />
        </div>

        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <select
            v-model="typeFilter"
            class="min-h-10 rounded-lg border border-border bg-surface px-3 text-sm text-text-primary focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            aria-label="按模板类型筛选"
          >
            <option value="">全部类型</option>
            <option v-for="(label, type) in promptTypeLabels" :key="type" :value="type">
              {{ label }}
            </option>
          </select>
          <span class="text-sm text-text-muted">{{ prompts.length }} 个模板</span>
        </div>
      </div>
    </section>

    <section class="overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
      <div class="flex items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 class="font-semibold text-text-primary">模板列表</h2>
          <p class="mt-0.5 text-xs text-text-muted">按工具类型管理当前生效版本</p>
        </div>
        <span class="rounded-full border border-border bg-surface-muted px-2.5 py-1 text-xs font-semibold text-text-secondary">
          {{ filteredPrompts.length }} 项
        </span>
      </div>

      <div v-if="filteredPrompts.length" class="overflow-x-auto">
        <table class="w-full min-w-[780px] text-sm">
          <thead class="bg-surface-muted/70 text-xs text-text-muted">
            <tr>
              <th class="px-5 py-3 text-left font-semibold">模板</th>
              <th class="px-5 py-3 text-left font-semibold">工具类型</th>
              <th class="px-5 py-3 text-left font-semibold">版本</th>
              <th class="px-5 py-3 text-left font-semibold">状态</th>
              <th class="px-5 py-3 text-right font-semibold">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="prompt in filteredPrompts"
              :key="prompt.id"
              class="border-t border-border transition-colors duration-200 hover:bg-surface-muted/50"
            >
              <td class="px-5 py-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent">
                    <FileText :size="18" />
                  </div>
                  <div>
                    <p class="font-semibold text-text-primary">{{ prompt.name }}</p>
                    <p class="mt-0.5 font-mono text-xs text-text-muted">ID {{ prompt.id }}</p>
                  </div>
                </div>
              </td>
              <td class="px-5 py-4">
                <span class="rounded-md border border-border bg-surface-muted px-2 py-1 text-xs font-medium text-text-secondary">
                  {{ promptTypeLabels[prompt.type] || prompt.type }}
                </span>
              </td>
              <td class="px-5 py-4 font-medium text-text-primary">v{{ prompt.version }}</td>
              <td class="px-5 py-4">
                <span
                  :class="[
                    'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-semibold',
                    prompt.status === 'active'
                      ? 'border-success/20 bg-success/10 text-success'
                      : 'border-border bg-surface-muted text-text-muted',
                  ]"
                >
                  <CircleCheck v-if="prompt.status === 'active'" :size="13" />
                  <CircleX v-else :size="13" />
                  {{ prompt.status === 'active' ? '已启用' : '已禁用' }}
                </span>
              </td>
              <td class="px-5 py-4">
                <div class="flex items-center justify-end gap-1">
                  <button
                    type="button"
                    title="查看模板"
                    aria-label="查看模板"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-surface-muted hover:text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                  >
                    <Eye :size="16" />
                  </button>
                  <button
                    type="button"
                    title="编辑模板"
                    aria-label="编辑模板"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-accent-soft hover:text-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                  >
                    <Pencil :size="16" />
                  </button>
                  <button
                    type="button"
                    title="删除模板"
                    aria-label="删除模板"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-danger/10 hover:text-danger focus-visible:ring-2 focus-visible:ring-danger focus:outline-none"
                  >
                    <Trash2 :size="16" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-else class="flex min-h-64 flex-col items-center justify-center px-6 text-center">
        <Search :size="30" class="mb-3 text-text-muted" />
        <p class="text-sm font-medium text-text-primary">没有匹配的模板</p>
        <p class="mt-1 text-xs text-text-muted">调整搜索关键词或模板类型。</p>
      </div>
    </section>
  </AdminPageShell>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  Bot,
  CircleCheck,
  CircleX,
  Pencil,
  Plus,
  Power,
  Search,
  Trash2,
} from '@lucide/vue'
import AdminPageShell from '@/components/admin/AdminPageShell.vue'

interface AIModel {
  id: string
  name: string
  provider: string
  status: 'active' | 'disabled'
}

const search = ref('')
const statusFilter = ref<'' | AIModel['status']>('')
const models = ref<AIModel[]>([
  { id: '1', name: 'GPT-4.1', provider: 'OpenAI', status: 'active' },
  { id: '2', name: 'Claude 3.5', provider: 'Anthropic', status: 'active' },
])

const filteredModels = computed(() => {
  const keyword = search.value.trim().toLowerCase()
  return models.value.filter((model) => {
    const matchesSearch = !keyword
      || model.name.toLowerCase().includes(keyword)
      || model.provider.toLowerCase().includes(keyword)
    const matchesStatus = !statusFilter.value || model.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

const activeCount = computed(() => models.value.filter((model) => model.status === 'active').length)
</script>

<template>
  <AdminPageShell
    :icon="Bot"
    title="AI 模型管理"
    description="维护可用模型、服务提供商和启用状态。"
    badge-text="模型配置"
  >
    <template #actions>
      <button
        type="button"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
      >
        <Plus :size="16" />
        添加模型
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
            placeholder="搜索模型名称或 Provider"
            aria-label="搜索模型"
          />
        </div>

        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <select
            v-model="statusFilter"
            class="min-h-10 rounded-lg border border-border bg-surface px-3 text-sm text-text-primary focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            aria-label="按模型状态筛选"
          >
            <option value="">全部状态</option>
            <option value="active">已启用</option>
            <option value="disabled">已禁用</option>
          </select>
          <span class="text-sm text-text-muted">
            {{ models.length }} 个模型，{{ activeCount }} 个已启用
          </span>
        </div>
      </div>
    </section>

    <section class="overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
      <div class="flex items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 class="font-semibold text-text-primary">模型列表</h2>
          <p class="mt-0.5 text-xs text-text-muted">当前系统可调用的模型配置</p>
        </div>
        <span class="rounded-full border border-border bg-surface-muted px-2.5 py-1 text-xs font-semibold text-text-secondary">
          {{ filteredModels.length }} 项
        </span>
      </div>

      <div v-if="filteredModels.length" class="overflow-x-auto">
        <table class="w-full min-w-[760px] text-sm">
          <thead class="bg-surface-muted/70 text-xs text-text-muted">
            <tr>
              <th class="px-5 py-3 text-left font-semibold">模型</th>
              <th class="px-5 py-3 text-left font-semibold">Provider</th>
              <th class="px-5 py-3 text-left font-semibold">状态</th>
              <th class="px-5 py-3 text-left font-semibold">模型 ID</th>
              <th class="px-5 py-3 text-right font-semibold">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="model in filteredModels"
              :key="model.id"
              class="border-t border-border transition-colors duration-200 hover:bg-surface-muted/50"
            >
              <td class="px-5 py-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent">
                    <Bot :size="18" />
                  </div>
                  <span class="font-semibold text-text-primary">{{ model.name }}</span>
                </div>
              </td>
              <td class="px-5 py-4 text-text-secondary">{{ model.provider }}</td>
              <td class="px-5 py-4">
                <span
                  :class="[
                    'inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-semibold',
                    model.status === 'active'
                      ? 'border-success/20 bg-success/10 text-success'
                      : 'border-border bg-surface-muted text-text-muted',
                  ]"
                >
                  <CircleCheck v-if="model.status === 'active'" :size="13" />
                  <CircleX v-else :size="13" />
                  {{ model.status === 'active' ? '已启用' : '已禁用' }}
                </span>
              </td>
              <td class="px-5 py-4 font-mono text-xs text-text-muted">{{ model.id }}</td>
              <td class="px-5 py-4">
                <div class="flex items-center justify-end gap-1">
                  <button
                    type="button"
                    title="编辑模型"
                    aria-label="编辑模型"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-accent-soft hover:text-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                  >
                    <Pencil :size="16" />
                  </button>
                  <button
                    type="button"
                    :title="model.status === 'active' ? '禁用模型' : '启用模型'"
                    :aria-label="model.status === 'active' ? '禁用模型' : '启用模型'"
                    class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-surface-muted hover:text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
                  >
                    <Power :size="16" />
                  </button>
                  <button
                    type="button"
                    title="删除模型"
                    aria-label="删除模型"
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
        <p class="text-sm font-medium text-text-primary">没有匹配的模型</p>
        <p class="mt-1 text-xs text-text-muted">调整搜索关键词或状态筛选。</p>
      </div>
    </section>
  </AdminPageShell>
</template>

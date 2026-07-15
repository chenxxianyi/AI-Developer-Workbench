<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { Bot, CircleCheck, CircleX, Pencil, Plus, Power, Search, Trash2, X } from '@lucide/vue'
import AdminPageShell from '@/components/admin/AdminPageShell.vue'
import apiClient from '@/api/client'

interface AIModel {
  id: string
  name: string
  provider: string
  base_url: string
  model: string
  vision_model: string
  timeout_seconds: number
  max_retries: number
  status: 'active' | 'disabled'
  is_default: boolean
}

type ModelForm = Omit<AIModel, 'id'>

const search = ref('')
const statusFilter = ref<'' | AIModel['status']>('')
const models = ref<AIModel[]>([])
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const success = ref('')
const editingId = ref<string | null>(null)
const showForm = ref(false)

const emptyForm = (): ModelForm => ({
  name: '',
  provider: 'openai',
  base_url: 'https://api.openai.com/v1',
  model: '',
  vision_model: '',
  timeout_seconds: 180,
  max_retries: 1,
  status: 'active',
  is_default: false,
})

const form = reactive<ModelForm>(emptyForm())

function resetForm(values?: Partial<ModelForm>) {
  Object.assign(form, emptyForm(), values || {})
}

async function loadModels() {
  loading.value = true
  error.value = ''
  try {
    models.value = await apiClient.get('/admin/models') as unknown as AIModel[]
  } catch (err: any) {
    error.value = err?.message || '获取模型配置失败'
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingId.value = null
  resetForm()
  showForm.value = true
  error.value = ''
  success.value = ''
}

function openEdit(model: AIModel) {
  editingId.value = model.id
  resetForm({
    name: model.name,
    provider: model.provider,
    base_url: model.base_url,
    model: model.model,
    vision_model: model.vision_model,
    timeout_seconds: model.timeout_seconds,
    max_retries: model.max_retries,
    status: model.status,
    is_default: model.is_default,
  })
  showForm.value = true
  error.value = ''
  success.value = ''
}

async function submitForm() {
  saving.value = true
  error.value = ''
  success.value = ''
  const payload = { ...form, vision_model: form.vision_model || form.model }
  try {
    if (editingId.value) {
      await apiClient.put(`/admin/models/${editingId.value}`, payload)
      success.value = '模型配置已更新'
    } else {
      await apiClient.post('/admin/models', payload)
      success.value = '模型配置已添加'
    }
    showForm.value = false
    await loadModels()
  } catch (err: any) {
    error.value = err?.message || '保存模型配置失败'
  } finally {
    saving.value = false
  }
}

async function toggleModel(model: AIModel) {
  await updateModel(model, { status: model.status === 'active' ? 'disabled' : 'active' })
}

async function setDefault(model: AIModel) {
  await updateModel(model, { status: 'active', is_default: true })
}

async function updateModel(model: AIModel, patch: Partial<ModelForm>) {
  error.value = ''
  success.value = ''
  const payload: ModelForm = {
    name: model.name,
    provider: model.provider,
    base_url: model.base_url,
    model: model.model,
    vision_model: model.vision_model,
    timeout_seconds: model.timeout_seconds,
    max_retries: model.max_retries,
    status: model.status,
    is_default: model.is_default,
    ...patch,
  }
  try {
    await apiClient.put(`/admin/models/${model.id}`, payload)
    success.value = payload.is_default ? '默认模型已切换，后续 AI 调用将使用该配置' : '模型状态已更新'
    await loadModels()
  } catch (err: any) {
    error.value = err?.message || '更新模型配置失败'
  }
}

async function deleteModel(model: AIModel) {
  if (!window.confirm(`确定删除模型配置「${model.name}」吗？`)) return
  error.value = ''
  success.value = ''
  try {
    await apiClient.delete(`/admin/models/${model.id}`)
    success.value = '模型配置已删除'
    await loadModels()
  } catch (err: any) {
    error.value = err?.message || '删除模型配置失败'
  }
}

onMounted(loadModels)

const filteredModels = computed(() => {
  const keyword = search.value.trim().toLowerCase()
  return models.value.filter((model) => {
    const matchesSearch = !keyword
      || model.name.toLowerCase().includes(keyword)
      || model.provider.toLowerCase().includes(keyword)
      || model.model.toLowerCase().includes(keyword)
      || model.base_url.toLowerCase().includes(keyword)
    const matchesStatus = !statusFilter.value || model.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

const activeCount = computed(() => models.value.filter((model) => model.status === 'active').length)
const defaultModel = computed(() => models.value.find((model) => model.is_default))
</script>

<template>
  <AdminPageShell
    :icon="Bot"
    title="AI 模型管理"
    description="这里显示真实后端配置和可切换的模型预设；不会展示 API Key。"
    badge-text="真实配置"
  >
    <template #actions>
      <button
        type="button"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="openCreate"
      >
        <Plus :size="16" />
        添加模型
      </button>
    </template>

    <section class="mb-5 rounded-lg border border-border bg-surface p-4 shadow-sm">
      <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <p class="text-sm text-text-muted">当前默认模型</p>
          <p class="mt-1 font-semibold text-text-primary">
            {{ defaultModel ? `${defaultModel.provider} / ${defaultModel.model}` : '未设置' }}
          </p>
        </div>
        <div class="relative w-full lg:max-w-md">
          <Search :size="17" class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
          <input
            v-model="search"
            type="search"
            class="min-h-10 w-full rounded-lg border border-border bg-surface-muted pl-10 pr-4 text-sm text-text-primary placeholder:text-text-muted focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            placeholder="搜索模型名称、Provider、模型 ID 或 Base URL"
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

    <p v-if="error" role="alert" class="mb-4 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">{{ error }}</p>
    <p v-if="success" role="status" class="mb-4 rounded-lg border border-success/20 bg-success/10 px-4 py-3 text-sm text-success">{{ success }}</p>

    <section v-if="showForm" class="mb-5 rounded-lg border border-border bg-surface p-5 shadow-sm">
      <div class="mb-4 flex items-center justify-between gap-3">
        <div>
          <h2 class="font-semibold text-text-primary">{{ editingId ? '编辑模型配置' : '添加模型配置' }}</h2>
          <p class="mt-0.5 text-xs text-text-muted">API Key 仍从后端 .env 读取，页面不会保存密钥。</p>
        </div>
        <button type="button" class="rounded-lg p-2 text-text-muted hover:bg-surface-muted" aria-label="关闭表单" @click="showForm = false">
          <X :size="18" />
        </button>
      </div>

      <form class="grid gap-4 md:grid-cols-2" @submit.prevent="submitForm">
        <label class="grid gap-1.5 text-sm font-medium text-text-secondary">
          显示名称
          <input v-model.trim="form.name" required class="min-h-10 rounded-lg border border-border bg-surface-muted px-3 text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" placeholder="例如 gpt-5.4" />
        </label>
        <label class="grid gap-1.5 text-sm font-medium text-text-secondary">
          Provider
          <input v-model.trim="form.provider" required class="min-h-10 rounded-lg border border-border bg-surface-muted px-3 text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" placeholder="openai" />
        </label>
        <label class="grid gap-1.5 text-sm font-medium text-text-secondary md:col-span-2">
          Base URL
          <input v-model.trim="form.base_url" required class="min-h-10 rounded-lg border border-border bg-surface-muted px-3 font-mono text-sm text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" placeholder="https://api.openai.com/v1" />
        </label>
        <label class="grid gap-1.5 text-sm font-medium text-text-secondary">
          文本模型 ID
          <input v-model.trim="form.model" required class="min-h-10 rounded-lg border border-border bg-surface-muted px-3 font-mono text-sm text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" placeholder="gpt-4.1" />
        </label>
        <label class="grid gap-1.5 text-sm font-medium text-text-secondary">
          视觉模型 ID
          <input v-model.trim="form.vision_model" class="min-h-10 rounded-lg border border-border bg-surface-muted px-3 font-mono text-sm text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" placeholder="留空则使用文本模型" />
        </label>
        <label class="grid gap-1.5 text-sm font-medium text-text-secondary">
          超时秒数
          <input v-model.number="form.timeout_seconds" type="number" min="1" class="min-h-10 rounded-lg border border-border bg-surface-muted px-3 text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" />
        </label>
        <label class="grid gap-1.5 text-sm font-medium text-text-secondary">
          重试次数
          <input v-model.number="form.max_retries" type="number" min="0" class="min-h-10 rounded-lg border border-border bg-surface-muted px-3 text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" />
        </label>
        <label class="inline-flex items-center gap-2 text-sm text-text-secondary">
          <input v-model="form.is_default" type="checkbox" class="h-4 w-4 accent-accent" />
          设为默认模型（后续 AI 调用使用该配置）
        </label>
        <label class="inline-flex items-center gap-2 text-sm text-text-secondary">
          <input v-model="form.status" type="checkbox" true-value="active" false-value="disabled" class="h-4 w-4 accent-accent" />
          启用
        </label>
        <div class="flex justify-end gap-2 md:col-span-2">
          <button type="button" class="min-h-10 rounded-lg border border-border px-4 text-sm font-semibold text-text-secondary hover:bg-surface-muted" @click="showForm = false">取消</button>
          <button type="submit" :disabled="saving" class="min-h-10 rounded-lg bg-accent px-4 text-sm font-semibold text-white hover:bg-accent/80 disabled:opacity-60">
            {{ saving ? '保存中...' : '保存配置' }}
          </button>
        </div>
      </form>
    </section>

    <section class="overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
      <div class="flex items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 class="font-semibold text-text-primary">模型列表</h2>
          <p class="mt-0.5 text-xs text-text-muted">来自数据库和后端 .env 的真实模型配置</p>
        </div>
        <span class="rounded-full border border-border bg-surface-muted px-2.5 py-1 text-xs font-semibold text-text-secondary">
          {{ filteredModels.length }} 项
        </span>
      </div>

      <div v-if="loading" class="flex min-h-64 items-center justify-center text-sm text-text-muted">正在读取模型配置...</div>

      <div v-else-if="filteredModels.length" class="overflow-x-auto">
        <table class="w-full min-w-[980px] text-sm">
          <thead class="bg-surface-muted/70 text-xs text-text-muted">
            <tr>
              <th class="px-5 py-3 text-left font-semibold">模型</th>
              <th class="px-5 py-3 text-left font-semibold">Provider</th>
              <th class="px-5 py-3 text-left font-semibold">Base URL</th>
              <th class="px-5 py-3 text-left font-semibold">状态</th>
              <th class="px-5 py-3 text-left font-semibold">默认</th>
              <th class="px-5 py-3 text-right font-semibold">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="model in filteredModels" :key="model.id" class="border-t border-border transition-colors duration-200 hover:bg-surface-muted/50">
              <td class="px-5 py-4">
                <div class="flex items-center gap-3">
                  <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent">
                    <Bot :size="18" />
                  </div>
                  <div>
                    <span class="font-semibold text-text-primary">{{ model.name }}</span>
                    <p class="mt-0.5 font-mono text-xs text-text-muted">{{ model.model }} / vision: {{ model.vision_model }}</p>
                  </div>
                </div>
              </td>
              <td class="px-5 py-4 text-text-secondary">{{ model.provider }}</td>
              <td class="px-5 py-4 font-mono text-xs text-text-muted">{{ model.base_url }}</td>
              <td class="px-5 py-4">
                <span :class="['inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-semibold', model.status === 'active' ? 'border-success/20 bg-success/10 text-success' : 'border-border bg-surface-muted text-text-muted']">
                  <CircleCheck v-if="model.status === 'active'" :size="13" />
                  <CircleX v-else :size="13" />
                  {{ model.status === 'active' ? '已启用' : '已禁用' }}
                </span>
              </td>
              <td class="px-5 py-4">
                <span v-if="model.is_default" class="rounded-full border border-accent/20 bg-accent-soft px-2.5 py-1 text-xs font-semibold text-accent">默认</span>
                <button v-else type="button" class="text-xs font-semibold text-accent hover:underline" @click="setDefault(model)">设为默认</button>
              </td>
              <td class="px-5 py-4">
                <div class="flex items-center justify-end gap-1">
                  <button type="button" title="编辑模型" aria-label="编辑模型" class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-accent-soft hover:text-accent focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" @click="openEdit(model)">
                    <Pencil :size="16" />
                  </button>
                  <button type="button" :title="model.status === 'active' ? '禁用模型' : '启用模型'" :aria-label="model.status === 'active' ? '禁用模型' : '启用模型'" class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-surface-muted hover:text-text-primary focus-visible:ring-2 focus-visible:ring-accent focus:outline-none" @click="toggleModel(model)">
                    <Power :size="16" />
                  </button>
                  <button type="button" title="删除模型" aria-label="删除模型" :disabled="model.is_default" class="flex h-9 w-9 items-center justify-center rounded-lg text-text-muted transition-smooth hover:bg-danger/10 hover:text-danger focus-visible:ring-2 focus-visible:ring-danger focus:outline-none disabled:cursor-not-allowed disabled:opacity-40" @click="deleteModel(model)">
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

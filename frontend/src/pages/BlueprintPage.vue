<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  CheckCircle2,
  FileText,
  Files,
  Loader2,
  Package,
  PenTool,
  RefreshCw,
} from '@lucide/vue'
import apiClient from '@/api/client'
import ProjectStageShell from '@/components/project/ProjectStageShell.vue'

interface BlueprintPageItem {
  name: string
  route: string
}

interface BlueprintContent {
  product_positioning?: string
  tech_stack?: string
  pages?: BlueprintPageItem[]
}

interface BlueprintRecord extends BlueprintContent {
  id?: string
  content?: string | BlueprintContent
  status?: string
  version?: number
}

const route = useRoute()
const router = useRouter()
const projectId = route.params.projectId as string
const blueprintRecord = ref<BlueprintRecord | null>(null)
const loading = ref(true)
const generating = ref(false)
const confirming = ref(false)
const error = ref('')

const blueprint = computed<BlueprintContent | null>(() => {
  if (!blueprintRecord.value) return null
  const source = blueprintRecord.value.content ?? blueprintRecord.value
  if (typeof source !== 'string') return source

  try {
    return JSON.parse(source) as BlueprintContent
  } catch {
    return null
  }
})

const pages = computed(() => blueprint.value?.pages ?? [])

async function loadBlueprint() {
  loading.value = true
  error.value = ''
  try {
    blueprintRecord.value = await apiClient.get(`/projects/${projectId}/blueprint`)
  } catch {
    blueprintRecord.value = null
  } finally {
    loading.value = false
  }
}

async function generate() {
  if (generating.value) return
  generating.value = true
  error.value = ''
  try {
    blueprintRecord.value = await apiClient.post(`/projects/${projectId}/blueprint/generate`)
  } catch (err: any) {
    error.value = err.message || '蓝图生成失败，请稍后重试'
  } finally {
    generating.value = false
  }
}

async function confirm() {
  if (!blueprintRecord.value || confirming.value) return
  confirming.value = true
  error.value = ''
  try {
    await apiClient.post(`/projects/${projectId}/blueprint/confirm`)
    await router.push(`/projects/${projectId}/generation`)
  } catch (err: any) {
    error.value = err.message || '蓝图确认失败，请稍后重试'
  } finally {
    confirming.value = false
  }
}

onMounted(loadBlueprint)
</script>

<template>
  <ProjectStageShell
    :icon="PenTool"
    title="蓝图评审"
    description="确认产品定位、技术方案和页面结构后进入代码生成。"
    step-text="蓝图阶段"
  >
    <template #actions>
      <button
        v-if="blueprint"
        type="button"
        :disabled="generating"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-primary transition-smooth hover:bg-surface-muted disabled:cursor-not-allowed disabled:opacity-60 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="generate"
      >
        <Loader2 v-if="generating" :size="16" class="animate-spin" />
        <RefreshCw v-else :size="16" />
        {{ generating ? '生成中...' : '重新生成' }}
      </button>
      <button
        v-if="blueprint"
        type="button"
        :disabled="confirming"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="confirm"
      >
        <Loader2 v-if="confirming" :size="16" class="animate-spin" />
        <CheckCircle2 v-else :size="16" />
        {{ confirming ? '确认中...' : '确认蓝图' }}
      </button>
    </template>

    <div v-if="error" role="alert" class="mb-5 flex items-center gap-2 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
      <AlertCircle :size="18" class="shrink-0" />
      {{ error }}
    </div>

    <div v-if="loading" class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1.25fr)_minmax(280px,0.75fr)]">
      <div class="h-96 animate-pulse rounded-lg border border-border bg-surface-muted" />
      <div class="h-72 animate-pulse rounded-lg border border-border bg-surface-muted" />
    </div>

    <div
      v-else-if="blueprint"
      class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1.25fr)_minmax(280px,0.75fr)]"
    >
      <section class="rounded-lg border border-border bg-surface p-5 shadow-sm sm:p-6">
        <div class="mb-5 flex items-start gap-3">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent">
            <FileText :size="20" />
          </div>
          <div>
            <h2 class="text-lg font-semibold text-text-primary">产品定位</h2>
            <p class="mt-1 text-sm leading-6 text-text-secondary">
              {{ blueprint.product_positioning || 'AI 生成的企业网站' }}
            </p>
          </div>
        </div>

        <div class="border-t border-border pt-5">
          <div class="mb-3 flex items-center justify-between gap-3">
            <div class="flex items-center gap-2">
              <Files :size="18" class="text-accent" />
              <h2 class="font-semibold text-text-primary">页面结构</h2>
            </div>
            <span class="text-xs font-medium text-text-muted">{{ pages.length }} 个页面</span>
          </div>

          <div v-if="pages.length" class="divide-y divide-border">
            <div
              v-for="page in pages"
              :key="`${page.name}-${page.route}`"
              class="flex min-h-12 items-center justify-between gap-4 py-3"
            >
              <span class="text-sm font-medium text-text-primary">{{ page.name }}</span>
              <code class="max-w-[60%] truncate rounded-md bg-surface-muted px-2 py-1 text-xs text-text-secondary">
                {{ page.route }}
              </code>
            </div>
          </div>
          <p v-else class="rounded-lg bg-surface-muted px-4 py-6 text-center text-sm text-text-muted">
            蓝图中暂未定义页面。
          </p>
        </div>
      </section>

      <aside class="space-y-6">
        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm">
          <div class="mb-3 flex items-center gap-3">
            <div class="flex h-9 w-9 items-center justify-center rounded-lg bg-accent-soft text-accent">
              <Package :size="18" />
            </div>
            <h2 class="font-semibold text-text-primary">技术栈</h2>
          </div>
          <p class="text-sm leading-6 text-text-secondary">
            {{ blueprint.tech_stack || 'Vue 3 + Tailwind CSS' }}
          </p>
        </section>

        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm">
          <h2 class="font-semibold text-text-primary">蓝图状态</h2>
          <dl class="mt-4 space-y-3 text-sm">
            <div class="flex items-center justify-between gap-4">
              <dt class="text-text-muted">版本</dt>
              <dd class="font-medium text-text-primary">v{{ blueprintRecord?.version || 1 }}</dd>
            </div>
            <div class="flex items-center justify-between gap-4">
              <dt class="text-text-muted">状态</dt>
              <dd class="rounded-full border border-success/20 bg-success/10 px-2.5 py-1 text-xs font-semibold text-success">
                {{ blueprintRecord?.status || 'generated' }}
              </dd>
            </div>
            <div class="flex items-center justify-between gap-4">
              <dt class="text-text-muted">页面数量</dt>
              <dd class="font-medium text-text-primary">{{ pages.length }}</dd>
            </div>
          </dl>
        </section>
      </aside>
    </div>

    <section
      v-else
      class="flex min-h-96 flex-col items-center justify-center rounded-lg border border-border bg-surface px-6 py-12 text-center shadow-sm"
    >
      <div class="mb-4 flex h-14 w-14 items-center justify-center rounded-lg bg-accent-soft text-accent">
        <PenTool :size="26" />
      </div>
      <h2 class="text-lg font-semibold text-text-primary">尚未生成蓝图</h2>
      <p class="mt-2 max-w-md text-sm leading-6 text-text-secondary">
        生成后可在这里评审产品定位、技术栈和页面结构。
      </p>
      <button
        type="button"
        :disabled="generating"
        class="mt-6 inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-5 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="generate"
      >
        <Loader2 v-if="generating" :size="16" class="animate-spin" />
        <PenTool v-else :size="16" />
        {{ generating ? '生成中...' : '生成蓝图' }}
      </button>
    </section>
  </ProjectStageShell>
</template>

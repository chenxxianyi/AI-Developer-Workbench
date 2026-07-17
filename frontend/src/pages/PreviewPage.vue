<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  AlertCircle,
  ExternalLink,
  Loader2,
  Monitor,
  RefreshCw,
} from '@lucide/vue'
import apiClient from '@/api/client'
import ProjectStageShell from '@/components/project/ProjectStageShell.vue'

const route = useRoute()
const projectId = route.params.projectId as string
const previewUrl = ref('')
const rebuilding = ref(false)
const checking = ref(true)
const error = ref('')

async function rebuild() {
  if (rebuilding.value) return
  rebuilding.value = true
  error.value = ''
  try {
    const result = await apiClient.post(
      `/projects/${projectId}/build`,
      undefined,
      { timeout: 650_000 },
    ) as { preview_url?: string }
    previewUrl.value = result.preview_url || ''
    if (!previewUrl.value) error.value = '构建已完成，但没有返回可用的预览地址'
  } catch (err: any) {
    error.value = err.message || '预览构建失败，请稍后重试'
  } finally {
    rebuilding.value = false
  }
}

async function loadPreviewStatus() {
  checking.value = true
  try {
    const result = await apiClient.get(`/projects/${projectId}/build`) as { ready?: boolean; preview_url?: string }
    if (result.ready) previewUrl.value = result.preview_url || ''
  } catch {
    // A project without a successful build starts in the empty state.
  } finally {
    checking.value = false
  }
}

onMounted(loadPreviewStatus)
</script>

<template>
  <ProjectStageShell
    :icon="Monitor"
    title="在线预览"
    description="在工作区中检查生成结果，并可随时重新构建最新版本。"
    step-text="预览阶段"
  >
    <template #actions>
      <button
        v-if="previewUrl"
        type="button"
        :disabled="rebuilding"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg border border-border bg-surface px-4 text-sm font-medium text-text-primary transition-smooth hover:bg-surface-muted disabled:cursor-not-allowed disabled:opacity-60 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="rebuild"
      >
        <Loader2 v-if="rebuilding" :size="16" class="animate-spin" />
        <RefreshCw v-else :size="16" />
        {{ rebuilding ? '构建中...' : '重新构建' }}
      </button>
      <a
        v-if="previewUrl"
        :href="previewUrl"
        target="_blank"
        rel="noopener noreferrer"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
      >
        <ExternalLink :size="16" />
        新窗口打开
      </a>
    </template>

    <div v-if="error" role="alert" class="mb-5 flex items-center gap-2 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
      <AlertCircle :size="18" class="shrink-0" />
      {{ error }}
    </div>

    <section v-if="checking" class="h-[520px] animate-pulse rounded-lg border border-border bg-surface-muted" aria-label="正在检查预览状态" />

    <section v-else-if="previewUrl" class="overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-3 border-b border-border bg-surface-muted/70 px-4 py-3">
        <div class="flex min-w-0 items-center gap-3">
          <div class="flex items-center gap-1.5" aria-hidden="true">
            <span class="h-2.5 w-2.5 rounded-full bg-danger/70" />
            <span class="h-2.5 w-2.5 rounded-full bg-warning/70" />
            <span class="h-2.5 w-2.5 rounded-full bg-success/70" />
          </div>
          <div class="min-w-0 rounded-md border border-border bg-surface px-3 py-1.5">
            <p class="truncate font-mono text-xs text-text-muted">{{ previewUrl }}</p>
          </div>
        </div>
        <span class="text-xs font-medium text-text-muted">实时预览</span>
      </div>
      <iframe
        :src="previewUrl"
        title="项目在线预览"
        class="h-[70vh] min-h-[460px] w-full bg-white"
        sandbox="allow-forms allow-modals allow-popups allow-scripts allow-downloads"
      />
    </section>

    <section
      v-else
      class="flex min-h-[520px] flex-col items-center justify-center rounded-lg border border-border bg-surface px-6 py-12 text-center shadow-sm"
    >
      <div class="mb-4 flex h-14 w-14 items-center justify-center rounded-lg bg-accent-soft text-accent">
        <Monitor :size="26" />
      </div>
      <h2 class="text-lg font-semibold text-text-primary">尚未构建预览</h2>
      <p class="mt-2 max-w-md text-sm leading-6 text-text-secondary">
        构建完成后，生成的网站会直接显示在当前工作区。
      </p>
      <button
        type="button"
        :disabled="rebuilding"
        class="mt-6 inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-5 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="rebuild"
      >
        <Loader2 v-if="rebuilding" :size="16" class="animate-spin" />
        <RefreshCw v-else :size="16" />
        {{ rebuilding ? '构建中...' : '构建并预览' }}
      </button>
    </section>
  </ProjectStageShell>
</template>

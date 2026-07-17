<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import {
  CheckCircle2,
  Circle,
  Download,
  ExternalLink,
  FileStack,
  FileText,
  Files,
  FolderKanban,
  LayoutDashboard,
  LockKeyhole,
  Monitor,
  PenTool,
  Play,
} from '@lucide/vue'
import apiClient from '@/api/client'
import { getProject } from '@/api/projects'
import type { Project } from '@/types/project'

const route = useRoute()
const projectId = computed(() => route.params.projectId as string)
const project = ref<Project | null>(null)
const loading = ref(false)
const requirementsReady = ref(false)
const blueprintConfirmed = ref(false)

const workflowItems = computed(() => [
  { to: `/projects/${projectId.value}/requirements`, icon: FileText, label: '需求', complete: requirementsReady.value, enabled: true },
  { to: `/projects/${projectId.value}/blueprint`, icon: PenTool, label: '蓝图', complete: blueprintConfirmed.value, enabled: requirementsReady.value },
  { to: `/projects/${projectId.value}/generation`, icon: Play, label: '生成', complete: false, enabled: blueprintConfirmed.value },
  { to: `/projects/${projectId.value}/preview`, icon: Monitor, label: '预览', complete: false, enabled: blueprintConfirmed.value },
])
const resourceItems = computed(() => [
  { to: `/projects/${projectId.value}`, icon: LayoutDashboard, label: '概览' },
  { to: `/projects/${projectId.value}/files`, icon: Files, label: '文件' },
  { to: `/projects/${projectId.value}/reports`, icon: FileStack, label: '报告' },
])

const projectTitle = computed(() => project.value?.name || `项目 #${projectId.value}`)
const projectDescription = computed(() => project.value?.description || '从需求、蓝图到代码交付的项目工作区')

async function loadProject() {
  if (!projectId.value) return
  loading.value = true
  try {
    project.value = await getProject(projectId.value)
  } catch {
    project.value = null
  } finally {
    loading.value = false
  }
}

async function loadWorkflow() {
  requirementsReady.value = false
  blueprintConfirmed.value = false
  try {
    const requirement: any = await apiClient.get(`/projects/${projectId.value}/requirements`)
    const content = typeof requirement?.content === 'string' ? JSON.parse(requirement.content) : requirement?.content
    requirementsReady.value = Boolean(
      content?.goal
      && content?.target_users?.length
      && content?.must_have_features?.length
      && content?.acceptance_criteria?.length,
    )
  } catch { /* no valid requirement yet */ }
  try {
    const blueprint: any = await apiClient.get(`/projects/${projectId.value}/blueprint`)
    blueprintConfirmed.value = blueprint?.status === 'confirmed'
  } catch { /* no blueprint yet */ }
}

function exportProject() {
  const base = import.meta.env.VITE_API_BASE_URL || '/api'
  window.open(`${base}/projects/${projectId.value}/export`, '_blank')
}

watch(projectId, () => { void loadProject(); void loadWorkflow() }, { immediate: true })
watch(() => route.path, () => { void loadWorkflow() })
</script>

<template>
  <div class="mx-auto max-w-6xl">
    <section class="mb-6 overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
      <div class="flex flex-col gap-4 p-4 sm:flex-row sm:items-center sm:justify-between sm:p-5">
        <div class="flex min-w-0 items-center gap-3">
          <div class="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent">
            <FolderKanban :size="22" />
          </div>
          <div class="min-w-0">
            <div v-if="loading" class="space-y-2">
              <div class="h-5 w-40 animate-pulse rounded bg-surface-muted" />
              <div class="h-3 w-64 max-w-full animate-pulse rounded bg-surface-muted" />
            </div>
            <template v-else>
              <div class="flex flex-wrap items-center gap-2">
                <h2 class="truncate text-lg font-semibold text-text-primary">{{ projectTitle }}</h2>
                <span class="rounded-full border border-border bg-surface-muted px-2 py-0.5 text-xs font-medium text-text-muted">
                  ID {{ projectId }}
                </span>
              </div>
              <p class="mt-0.5 truncate text-sm text-text-secondary">{{ projectDescription }}</p>
            </template>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            class="inline-flex min-h-10 items-center gap-2 rounded-lg border border-border bg-surface px-3 text-sm font-medium text-text-primary transition-smooth hover:bg-surface-muted focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
            @click="exportProject"
          >
            <Download :size="16" />
            <span class="hidden sm:inline">导出</span>
          </button>
          <RouterLink
            :to="`/projects/${projectId}/preview`"
            class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-3 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
          >
            <ExternalLink :size="16" />
            <span class="hidden sm:inline">打开预览</span>
          </RouterLink>
        </div>
      </div>

      <nav class="border-t border-border px-3 py-3" aria-label="项目交付流程">
        <div class="mb-2 flex items-center justify-between gap-3">
          <span class="text-xs font-semibold uppercase tracking-wide text-text-muted">交付流程</span>
          <span class="text-xs text-text-muted">确认上一步后继续</span>
        </div>
        <ol class="flex min-w-max items-center overflow-x-auto">
          <li v-for="(item, index) in workflowItems" :key="item.label" class="flex items-center">
          <RouterLink
            v-if="item.enabled"
            :to="item.to"
            :aria-current="route.path === item.to ? 'page' : undefined"
            :class="[
              'flex min-h-10 items-center gap-2 rounded-lg px-3 text-sm font-medium transition-colors duration-200',
              route.path === item.to
                ? 'bg-accent-soft text-accent'
                : item.complete ? 'text-success hover:bg-success/5' : 'text-text-muted hover:bg-surface-muted hover:text-text-primary',
            ]"
          >
            <CheckCircle2 v-if="item.complete" :size="17" />
            <Circle v-else :size="17" />
            <span>{{ index + 1 }}. {{ item.label }}</span>
          </RouterLink>
          <span v-else class="flex min-h-10 cursor-not-allowed items-center gap-2 rounded-lg px-3 text-sm font-medium text-text-muted/60" :title="`请先完成${index === 1 ? '需求' : '蓝图确认'}`">
            <LockKeyhole :size="16" /><span>{{ index + 1 }}. {{ item.label }}</span>
          </span>
          <span v-if="index < workflowItems.length - 1" class="mx-1 h-px w-5 bg-border" aria-hidden="true" />
          </li>
        </ol>
        <div class="mt-3 flex flex-wrap items-center gap-1 border-t border-border pt-3">
          <span class="mr-2 text-xs font-semibold uppercase tracking-wide text-text-muted">项目资源</span>
          <RouterLink v-for="item in resourceItems" :key="item.label" :to="item.to" :class="['flex min-h-9 items-center gap-2 rounded-lg px-3 text-sm font-medium transition-colors', route.path === item.to ? 'bg-surface-muted text-text-primary' : 'text-text-muted hover:text-text-primary']">
            <component :is="item.icon" :size="16" />{{ item.label }}
          </RouterLink>
        </div>
      </nav>
    </section>

    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertCircle,
  CheckCircle2,
  FileCode,
  Files,
  Loader2,
  Monitor,
  Play,
  RotateCcw,
  Terminal,
  XCircle,
} from '@lucide/vue'
import apiClient from '@/api/client'
import ProjectStageShell from '@/components/project/ProjectStageShell.vue'

interface GenerationLog {
  text: string
  level: 'info' | 'warn' | 'error' | 'success'
}

const route = useRoute()
const router = useRouter()
const projectId = computed(() => route.params.projectId as string)
const taskId = ref('')
const progress = ref(0)
const currentStage = ref('')
const statusText = ref('等待启动')
const running = ref(false)
const completed = ref(false)
const failed = ref(false)
const starting = ref(false)
const error = ref('')
const logs = ref<GenerationLog[]>([])
const logContainer = ref<HTMLElement>()
let eventSource: EventSource | null = null

const safeProgress = computed(() => Math.min(100, Math.max(0, progress.value)))
const statusTone = computed(() => {
  if (completed.value) return 'border-success/20 bg-success/10 text-success'
  if (failed.value) return 'border-danger/20 bg-danger/10 text-danger'
  if (running.value) return 'border-accent/20 bg-accent-soft text-accent'
  return 'border-border bg-surface-muted text-text-secondary'
})

const generationStages = [
  { icon: FileCode, label: '生成项目结构与核心代码' },
  { icon: Files, label: '整理资源、配置和依赖文件' },
  { icon: Monitor, label: '完成后进入在线预览' },
]

function addLog(text: string, level: GenerationLog['level'] = 'info') {
  logs.value.push({ text: `[${new Date().toLocaleTimeString()}] ${text}`, level })
  nextTick(() => {
    if (logContainer.value) logContainer.value.scrollTop = logContainer.value.scrollHeight
  })
}

async function startGeneration() {
  if (starting.value || running.value) return

  eventSource?.close()
  starting.value = true
  error.value = ''
  completed.value = false
  failed.value = false
  progress.value = 0
  currentStage.value = ''
  logs.value = []
  statusText.value = '正在启动'
  addLog('正在启动代码生成任务')

  try {
    const task = await apiClient.post('/tasks', {
      project_id: projectId.value,
      type: 'generation',
    }) as { id: string }
    taskId.value = task.id
    running.value = true
    statusText.value = '生成中'
    addLog(`任务已创建：${task.id}`)
    connectSSE(task.id)
  } catch (err: any) {
    failed.value = true
    statusText.value = '启动失败'
    error.value = err.message || '代码生成启动失败，请稍后重试'
    addLog('启动生成失败', 'error')
  } finally {
    starting.value = false
  }
}

function handleTaskEvent(event: MessageEvent) {
  try {
    const data = JSON.parse(event.data)
    progress.value = Number(data.progress) || 0
    currentStage.value = data.message || data.stage || ''
    statusText.value = data.status === 'running' ? '生成中' : data.status || statusText.value

    if (data.type === 'task_completed') {
      progress.value = 100
      completed.value = true
      running.value = false
      statusText.value = '生成完成'
      addLog('代码生成已完成', 'success')
      eventSource?.close()
    } else if (data.type === 'task_failed') {
      failed.value = true
      running.value = false
      statusText.value = '生成失败'
      error.value = data.message || '代码生成失败'
      addLog(`生成失败：${data.message || '未知错误'}`, 'error')
      eventSource?.close()
    } else if (data.message) {
      addLog(data.message)
    }
  } catch {
    addLog('收到无法解析的任务消息', 'warn')
  }
}

function connectSSE(id: string) {
  const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'
  const token = localStorage.getItem('auth_token')
  const url = `${base}/tasks/${id}/stream${token ? `?token=${encodeURIComponent(token)}` : ''}`

  eventSource = new EventSource(url)
  eventSource.onmessage = handleTaskEvent
  eventSource.addEventListener('stage_progress', handleTaskEvent as EventListener)
  eventSource.addEventListener('task_completed', handleTaskEvent as EventListener)
  eventSource.addEventListener('task_failed', handleTaskEvent as EventListener)
  eventSource.onerror = () => {
    if (running.value) addLog('实时连接暂时中断，正在等待重连', 'warn')
  }
}

function cancelGeneration() {
  if (taskId.value) void apiClient.post(`/tasks/${taskId.value}/cancel`)
  eventSource?.close()
  running.value = false
  statusText.value = '已取消'
  addLog('生成任务已取消', 'warn')
}

function openPreview() {
  void router.push(`/projects/${projectId.value}/preview`)
}

onUnmounted(() => eventSource?.close())
</script>

<template>
  <ProjectStageShell
    :icon="Play"
    title="代码生成"
    description="根据已确认的蓝图生成项目结构、代码和配置文件。"
    step-text="生成阶段"
  >
    <template v-if="completed" #actions>
      <button
        type="button"
        class="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
        @click="openPreview"
      >
        <Monitor :size="16" />
        查看预览
      </button>
    </template>

    <div v-if="error" role="alert" class="mb-5 flex items-center gap-2 rounded-lg border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
      <AlertCircle :size="18" class="shrink-0" />
      {{ error }}
    </div>

    <div
      v-if="!taskId"
      class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1.2fr)_minmax(280px,0.8fr)]"
    >
      <section class="flex min-h-96 flex-col items-center justify-center rounded-lg border border-border bg-surface px-6 py-12 text-center shadow-sm">
        <div class="mb-4 flex h-14 w-14 items-center justify-center rounded-lg bg-accent-soft text-accent">
          <Play :size="26" />
        </div>
        <h2 class="text-lg font-semibold text-text-primary">准备开始生成</h2>
        <p class="mt-2 max-w-md text-sm leading-6 text-text-secondary">
          启动后可在本页查看实时进度、当前阶段和生成日志。
        </p>
        <button
          type="button"
          :disabled="starting"
          class="mt-6 inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-5 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
          @click="startGeneration"
        >
          <Loader2 v-if="starting" :size="16" class="animate-spin" />
          <Play v-else :size="16" />
          {{ starting ? '启动中...' : '开始生成' }}
        </button>
      </section>

      <aside class="rounded-lg border border-border bg-surface p-5 shadow-sm">
        <h2 class="font-semibold text-text-primary">生成范围</h2>
        <div class="mt-4 divide-y divide-border">
          <div
            v-for="stage in generationStages"
            :key="stage.label"
            class="flex items-center gap-3 py-4"
          >
            <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-surface-muted text-text-secondary">
              <component :is="stage.icon" :size="18" />
            </div>
            <span class="text-sm text-text-secondary">{{ stage.label }}</span>
          </div>
        </div>
      </aside>
    </div>

    <div v-else class="grid grid-cols-1 gap-6 lg:grid-cols-[minmax(0,1.35fr)_minmax(280px,0.65fr)]">
      <div class="space-y-6">
        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm sm:p-6">
          <div class="mb-4 flex flex-wrap items-start justify-between gap-3">
            <div>
              <p class="text-sm font-medium text-text-primary">{{ statusText }}</p>
              <p class="mt-1 text-sm text-text-muted">{{ currentStage || '等待任务进度更新' }}</p>
            </div>
            <span class="text-2xl font-bold text-accent">{{ safeProgress }}%</span>
          </div>
          <div
            class="h-2.5 w-full overflow-hidden rounded-full bg-surface-muted"
            role="progressbar"
            :aria-valuenow="safeProgress"
            aria-valuemin="0"
            aria-valuemax="100"
          >
            <div
              class="h-full rounded-full bg-accent transition-[width] duration-500"
              :style="{ width: `${safeProgress}%` }"
            />
          </div>
        </section>

        <section class="overflow-hidden rounded-lg border border-border bg-surface shadow-sm">
          <div class="flex items-center justify-between border-b border-border px-4 py-3">
            <div class="flex items-center gap-2">
              <Terminal :size="17" class="text-accent" />
              <h2 class="text-sm font-semibold text-text-primary">生成日志</h2>
            </div>
            <span class="text-xs text-text-muted">{{ logs.length }} 条</span>
          </div>
          <div
            ref="logContainer"
            class="h-80 overflow-y-auto bg-[#171717] p-4 font-mono text-xs leading-6 text-[#d4d4d4]"
            aria-live="polite"
          >
            <div
              v-for="(log, index) in logs"
              :key="`${index}-${log.text}`"
              :class="[
                log.level === 'error'
                  ? 'text-red-300'
                  : log.level === 'warn'
                    ? 'text-amber-300'
                    : log.level === 'success'
                      ? 'text-emerald-300'
                      : 'text-zinc-300',
              ]"
            >
              {{ log.text }}
            </div>
            <div v-if="running" class="mt-1 inline-block h-4 w-2 animate-pulse bg-accent" />
          </div>
        </section>
      </div>

      <aside class="space-y-6">
        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm">
          <div class="flex items-center justify-between gap-3">
            <h2 class="font-semibold text-text-primary">任务状态</h2>
            <span :class="['rounded-full border px-2.5 py-1 text-xs font-semibold', statusTone]">
              {{ statusText }}
            </span>
          </div>
          <dl class="mt-4 space-y-3 text-sm">
            <div class="flex items-center justify-between gap-3">
              <dt class="text-text-muted">任务 ID</dt>
              <dd class="max-w-40 truncate font-mono text-xs text-text-secondary">{{ taskId }}</dd>
            </div>
            <div class="flex items-center justify-between gap-3">
              <dt class="text-text-muted">当前进度</dt>
              <dd class="font-medium text-text-primary">{{ safeProgress }}%</dd>
            </div>
          </dl>
        </section>

        <section class="rounded-lg border border-border bg-surface p-5 shadow-sm">
          <h2 class="font-semibold text-text-primary">任务操作</h2>
          <div class="mt-4 space-y-2">
            <button
              v-if="running"
              type="button"
              class="inline-flex min-h-10 w-full items-center justify-center gap-2 rounded-lg border border-danger/30 bg-surface px-4 text-sm font-medium text-danger transition-smooth hover:bg-danger/5 focus-visible:ring-2 focus-visible:ring-danger focus:outline-none"
              @click="cancelGeneration"
            >
              <XCircle :size="16" />
              取消生成
            </button>
            <button
              v-if="failed || !running && !completed"
              type="button"
              :disabled="starting"
              class="inline-flex min-h-10 w-full items-center justify-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 disabled:opacity-60 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
              @click="startGeneration"
            >
              <Loader2 v-if="starting" :size="16" class="animate-spin" />
              <RotateCcw v-else :size="16" />
              重新生成
            </button>
            <button
              v-if="completed"
              type="button"
              class="inline-flex min-h-10 w-full items-center justify-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-white transition-smooth hover:bg-accent/80 focus-visible:ring-2 focus-visible:ring-accent focus:outline-none"
              @click="openPreview"
            >
              <CheckCircle2 :size="16" />
              查看生成结果
            </button>
          </div>
        </section>
      </aside>
    </div>
  </ProjectStageShell>
</template>


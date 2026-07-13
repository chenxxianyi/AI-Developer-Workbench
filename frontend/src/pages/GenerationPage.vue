<template>
  <div class="max-w-3xl mx-auto py-8">
    <h1 class="text-2xl font-bold text-[var(--color-text-primary)] mb-6">代码生成</h1>

    <!-- 启动生成 -->
    <div v-if="!taskId" class="text-center py-16">
      <div class="text-4xl mb-4">🚀</div>
      <h2 class="text-lg font-semibold mb-2">准备开始生成</h2>
      <p class="text-[var(--color-text-muted)] mb-6">确认蓝图后即可启动代码生成</p>
      <button @click="startGeneration" :disabled="starting" class="px-6 py-3 rounded-lg bg-[var(--color-accent)] text-white font-medium disabled:opacity-50">
        {{ starting ? '启动中...' : '开始生成' }}
      </button>
    </div>

    <!-- 生成进度 -->
    <div v-else class="space-y-6">
      <!-- 进度条 -->
      <div class="bg-[var(--color-surface)] rounded-[var(--radius-panel)] border border-[var(--color-border)] p-6">
        <div class="flex items-center justify-between mb-4">
          <span class="text-sm font-medium text-[var(--color-text-primary)]">{{ statusText }}</span>
          <span class="text-2xl font-bold text-[var(--color-accent)]">{{ progress }}%</span>
        </div>
        <div class="w-full bg-[var(--color-surface-muted)] rounded-full h-3 overflow-hidden">
          <div class="h-full bg-[var(--color-accent)] rounded-full transition-all duration-500" :style="{ width: progress + '%' }" />
        </div>
        <p v-if="currentStage" class="mt-3 text-sm text-[var(--color-text-muted)]">{{ currentStage }}</p>
      </div>

      <!-- 日志 -->
      <div class="bg-[var(--color-surface)] rounded-[var(--radius-panel)] border border-[var(--color-border)] p-4">
        <h3 class="text-sm font-semibold mb-3">生成日志</h3>
        <div ref="logContainer" class="bg-[#1e1e1e] text-[#d4d4d4] rounded-lg p-4 h-64 overflow-y-auto font-mono text-xs space-y-1">
          <div v-for="(log, i) in logs" :key="i" :class="log.level === 'error' ? 'text-red-400' : log.level === 'warn' ? 'text-yellow-400' : 'text-green-400'">
            {{ log.text }}
          </div>
          <div v-if="running" class="animate-pulse text-[var(--color-accent)]">▊</div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex gap-3">
        <button v-if="running" @click="cancelGeneration" class="px-4 py-2 rounded-lg border border-[var(--color-danger)] text-[var(--color-danger)] text-sm">取消</button>
        <button v-if="completed" @click="$router.push(`/projects/${projectId}/preview`)" class="px-6 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm">查看预览</button>
        <button v-if="failed" @click="startGeneration" class="px-6 py-2 rounded-lg bg-[var(--color-accent)] text-white text-sm">重新生成</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import apiClient from '@/api/client'

const route = useRoute()
const projectId = computed(() => route.params.projectId as string)

const taskId = ref('')
const progress = ref(0)
const currentStage = ref('')
const statusText = ref('初始化...')
const running = ref(false)
const completed = ref(false)
const failed = ref(false)
const starting = ref(false)
const logs = ref<{ text: string; level: string }[]>([])
const logContainer = ref<HTMLElement>()
let eventSource: EventSource | null = null

function addLog(text: string, level = 'info') {
  logs.value.push({ text: `[${new Date().toLocaleTimeString()}] ${text}`, level })
  nextTick(() => { if (logContainer.value) logContainer.value.scrollTop = logContainer.value.scrollHeight })
}

async function startGeneration() {
  starting.value = true
  addLog('正在启动代码生成...')
  try {
    const task = await apiClient.post('/tasks', { project_id: projectId.value, type: 'generation' }) as any
    taskId.value = task.id
    running.value = true
    completed.value = false
    failed.value = false
    addLog(`任务已创建: ${task.id}`)
    connectSSE(task.id)
  } catch {
    addLog('启动生成失败', 'error')
  } finally { starting.value = false }
}

function connectSSE(id: string) {
  const base = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  const token = localStorage.getItem('auth_token')
  const url = `${base}/tasks/${id}/stream${token ? '?token=' + token : ''}`
  
  eventSource = new EventSource(url)
  eventSource.onmessage = (e) => {
    try {
      const data = JSON.parse(e.data)
      progress.value = data.progress || 0
      currentStage.value = data.message || data.stage || ''
      statusText.value = data.status === 'running' ? '生成中...' : data.status
      
      if (data.type === 'task_completed') {
        completed.value = true; running.value = false
        addLog('✅ 生成完成！')
        eventSource?.close()
      } else if (data.type === 'task_failed') {
        failed.value = true; running.value = false
        addLog(`❌ 生成失败: ${data.message}`, 'error')
        eventSource?.close()
      } else if (data.message) {
        addLog(data.message)
      }
    } catch { /* ignore parse errors */ }
  }
  eventSource.onerror = () => { addLog('SSE 连接中断，尝试重连...', 'warn') }
}

function cancelGeneration() {
  if (taskId.value) apiClient.post(`/tasks/${taskId.value}/cancel`)
  eventSource?.close()
  running.value = false
  addLog('生成已取消', 'warn')
}

onUnmounted(() => eventSource?.close())
</script>

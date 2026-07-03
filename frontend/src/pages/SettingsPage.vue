<script setup lang="ts">
/**
 * Settings Page
 * Read-only system status display
 */

import { onMounted } from 'vue'
import { useSystemStore } from '@/stores/systemStore'
import { CheckCircle2, AlertCircle, Info, FlaskConical, Zap } from '@lucide/vue'

const systemStore = useSystemStore()

onMounted(async () => {
  await systemStore.fetchStatus()
})

const uploadLimits = systemStore.uploadLimits
</script>

<template>
  <div>
    <div v-if="systemStore.loading" class="text-center py-8">
      <div class="text-text-muted">加载中...</div>
    </div>

    <div v-else class="max-w-2xl">
      <!-- System Status Card -->
      <div class="p-6 bg-surface border border-border rounded-lg mb-6">
        <h2 class="text-xl font-semibold mb-4 text-text-primary">系统状态</h2>

        <div class="flex items-center gap-3 mb-4">
          <CheckCircle2
            v-if="systemStore.status?.healthy"
            :size="24"
            class="text-success"
          />
          <AlertCircle v-else :size="24" class="text-danger" />
          <span
            :class="[
              'font-medium',
              systemStore.status?.healthy ? 'text-success' : 'text-danger',
            ]"
          >
            {{ systemStore.status?.healthy ? '服务正常运行' : '服务异常' }}
          </span>
        </div>

        <!-- Mock Mode Indicator -->
        <div
          :class="[
            'flex items-center gap-3 px-4 py-3 rounded-md mb-4',
            systemStore.isMockMode
              ? 'bg-amber-50 border border-amber-200 dark:bg-amber-900/20 dark:border-amber-700/40'
              : 'bg-emerald-50 border border-emerald-200 dark:bg-emerald-900/20 dark:border-emerald-700/40',
          ]"
        >
          <FlaskConical
            v-if="systemStore.isMockMode"
            :size="20"
            class="text-amber-600 dark:text-amber-400"
          />
          <Zap
            v-else
            :size="20"
            class="text-emerald-600 dark:text-emerald-400"
          />
          <div>
            <div class="font-semibold text-sm" :class="systemStore.isMockMode ? 'text-amber-800 dark:text-amber-200' : 'text-emerald-800 dark:text-emerald-200'">
              {{ systemStore.isMockMode ? '演示模式' : '真实 AI 模式' }}
            </div>
            <div class="text-xs mt-0.5" :class="systemStore.isMockMode ? 'text-amber-600 dark:text-amber-400' : 'text-emerald-600 dark:text-emerald-400'">
              {{ systemStore.isMockMode ? '当前使用 Mock 数据，不会调用外部 AI 服务。适合演示和本地测试。' : '当前连接真实 AI 服务，分析结果由 AI 生成。' }}
            </div>
          </div>
        </div>

        <div class="space-y-2 text-sm">
          <div class="flex items-center justify-between">
            <span class="text-text-secondary">AI 服务商</span>
            <span class="font-medium text-text-primary">{{ systemStore.status?.provider }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-text-secondary">文本模型</span>
            <span class="font-medium text-text-primary">{{ systemStore.status?.text_model }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-text-secondary">视觉模型</span>
            <span class="font-medium text-text-primary">{{ systemStore.status?.vision_model }}</span>
          </div>
        </div>
      </div>

      <!-- Upload Limits Card -->
      <div class="p-6 bg-surface border border-border rounded-lg mb-6">
        <div class="flex items-center gap-2 mb-4">
          <Info :size="20" class="text-accent" />
          <h2 class="text-xl font-semibold text-text-primary">上传限制</h2>
        </div>

        <div class="space-y-2 text-sm">
          <div class="flex items-center justify-between">
            <span class="text-text-secondary">图片最大大小</span>
            <span class="font-medium text-text-primary">
              {{ (uploadLimits.image_max_bytes / 1_000_000).toFixed(0) }}MB
            </span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-text-secondary">ZIP 最大大小</span>
            <span class="font-medium text-text-primary">
              {{ (uploadLimits.zip_max_bytes / 1_000_000).toFixed(0) }}MB
            </span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-text-secondary">ZIP 最大文件数</span>
            <span class="font-medium text-text-primary">{{ uploadLimits.zip_max_files }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-text-secondary">ZIP 总大小限制</span>
            <span class="font-medium text-text-primary">
              {{ (uploadLimits.zip_max_total_bytes / 1_000_000).toFixed(0) }}MB
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

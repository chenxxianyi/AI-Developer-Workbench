<script setup lang="ts">
/**
 * Settings Page
 * Read-only system status display
 */

import { onMounted } from 'vue'
import { useSystemStore } from '@/stores/systemStore'
import { CheckCircle2, AlertCircle, Info } from '@lucide/vue'

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

        <div class="space-y-2 text-sm">
          <div class="flex items-center justify-between">
            <span class="text-text-secondary">AI Provider</span>
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
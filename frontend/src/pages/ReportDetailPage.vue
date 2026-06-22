<script setup lang="ts">
/**
 * Report Detail Page Placeholder
 */

import { RouterLink, useRoute } from 'vue-router'
import { FileStack, ArrowLeft } from '@lucide/vue'

const route = useRoute()
const reportId = route.params.id
</script>

<template>
  <div class="text-center py-12">
    <FileStack :size="48" class="mx-auto mb-4 text-text-muted" />
    <h2 class="text-2xl font-bold mb-2 text-text-primary">报告详情</h2>
    <p class="text-text-secondary mb-2">此页面将在第三阶段实现</p>
    <p class="text-sm text-text-muted mb-6">报告编号: {{ reportId }}</p>
    <RouterLink
      to="/reports"
      class="inline-flex items-center gap-2 px-6 py-3 bg-surface border border-border text-text-primary rounded-lg font-semibold hover:bg-surface-muted transition-smooth"
    >
      <ArrowLeft :size="20" />
      返回列表
    </RouterLink>
  </div>
</template>

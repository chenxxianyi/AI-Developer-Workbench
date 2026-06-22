<script setup lang="ts">
/**
 * Reports Page Placeholder
 */

import { RouterLink } from 'vue-router'
import { FileStack, ArrowRight } from '@lucide/vue'
</script>

<template>
  <div class="text-center py-12">
    <FileStack :size="48" class="mx-auto mb-4 text-text-muted" />
    <h2 class="text-2xl font-bold mb-2 text-text-primary">报告列表</h2>
    <p class="text-text-secondary mb-6">此页面将在第三阶段实现</p>
    <RouterLink
      to="/dashboard"
      class="inline-flex items-center gap-2 px-6 py-3 bg-accent text-white rounded-lg font-semibold hover:bg-accent/90 transition-smooth"
    >
      <ArrowRight :size="20" />
      返回工作台
    </RouterLink>
  </div>
</template>

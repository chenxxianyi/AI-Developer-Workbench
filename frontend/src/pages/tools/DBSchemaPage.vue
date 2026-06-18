<script setup lang="ts">
/**
 * Tool Page Placeholder - DB Schema Review
 */

import { RouterLink } from 'vue-router'
import { Database, ArrowLeft } from '@lucide/vue'
</script>

<template>
  <div class="text-center py-12">
    <div class="w-16 h-16 bg-accent-soft rounded-lg flex items-center justify-center mx-auto mb-4">
      <Database :size="32" class="text-accent" />
    </div>
    <h2 class="text-2xl font-bold mb-2 text-text-primary">DB Schema Review</h2>
    <p class="text-text-secondary mb-6">此工具页面将在 Phase 4 实现</p>
    <RouterLink
      to="/dashboard"
      class="inline-flex items-center gap-2 px-6 py-3 bg-surface border border-border text-text-primary rounded-lg font-semibold hover:bg-surface-muted transition-smooth"
    >
      <ArrowLeft :size="20" />
      返回 Dashboard
    </RouterLink>
  </div>
</template>
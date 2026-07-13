<template>
  <div class="space-y-4">
    <h1 class="text-2xl font-bold">项目管理</h1>
    <div class="bg-[var(--color-surface)] rounded-lg border overflow-hidden"><table class="w-full text-sm"><thead><tr class="border-b bg-[var(--color-surface-muted)]"><th class="px-4 py-3 text-left">名称</th><th class="px-4 py-3 text-left">用户</th><th class="px-4 py-3 text-left">状态</th><th class="px-4 py-3 text-left">创建时间</th><th class="px-4 py-3 text-right">操作</th></tr></thead><tbody><tr v-for="p in projects" :key="p.id" class="border-b hover:bg-[var(--color-surface-muted)]/50"><td class="px-4 py-3 font-medium">{{ p.name }}</td><td class="px-4 py-3 text-[var(--color-text-muted)]">{{ p.user_id || '—' }}</td><td class="px-4 py-3">{{ p.status || 'draft' }}</td><td class="px-4 py-3 text-[var(--color-text-muted)]">{{ p.created_at?.slice(0,10) }}</td><td class="px-4 py-3 text-right"><RouterLink :to="`/projects/${p.id}`" class="text-[var(--color-accent)] hover:underline text-xs">查看</RouterLink></td></tr></tbody></table></div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
const projects = ref([{ id: '1', name: '示例项目', status: 'draft', created_at: '2026-07-13' }])
</script>
